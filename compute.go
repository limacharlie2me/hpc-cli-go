package main

import (
	"context"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	sashay "bitbucket.oci.oraclecorp.com/one-punch/sashay2"
	computeadmin "bitbucket.oci.oraclecorp.com/one-punch/sashay2/computeadmin"
)

func compute(ctx context.Context) ([]Instance, error) { //nolint:funlen
	log := log.With().Str("Function", "compute").Logger()
	defer log.Debug().Msg("Return")

	instances := make([]Instance, 0)
	hostOutChan := make(chan computeadmin.HostV2, 4096)
	instanceOutChan := make(chan Instance, 4096)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(hostOutChan)
		defer log.Debug().Msg("Close hostOutChan")

		if err := fetchComputeHosts(ctx, hostOutChan); err != nil {
			return eris.Wrapf(err, "fetchCompute.fetchComputeHosts")
		}

		return nil
	})

	g.Go(func() error {
		defer close(instanceOutChan)
		defer log.Debug().Msg("Close instanceOutChan")

		if err := filterRoutine(ctx, hostOutChan, instanceOutChan); err != nil {
			return eris.Wrap(err, "filterRoutine")
		}

		return nil
	})

	bar := p.Bar("[INFO] Filtering HPC Hosts: \t\t")
	bar.SetPriority(2)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Debug().Msg("Done")

				bar.SetTotal(-1, true)

				return
			case instance, ok := <-instanceOutChan:
				if !ok {
					bar.SetTotal(-1, true)

					return
				}

				log.Debug().Int("Instance OUT CHAN", len(instanceOutChan)).Msg("")

				instances = append(instances, instance)

				bar.IncrBy(1)
				bar.SetTotal(bar.Current()+100, false)
			}
		}
	}()

	wg.Wait()

	err := g.Wait()

	if err != nil {
		return nil, eris.Wrapf(err, "errGroup error")
	}

	p.Wait()

	return instances, nil
}

func fetchComputeHosts(ctx context.Context, hostOutChan chan computeadmin.HostV2) error {
	log := log.With().Str("Function", "fetchComputeHosts").Logger()
	defer log.Debug().Msg("Return")

	computeAdminV2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	bar := p.Bar("[INFO] Fetching Compute Hosts: \t\t")
	bar.SetPriority(1)

	listHostsRequest := computeAdminV2.ListHostsV2(ctx)

	walkFunc := func(host computeadmin.HostV2) error {
		hostOutChan <- host

		bar.IncrBy(1)
		bar.SetTotal(bar.Current()+100, false)

		return nil
	}

	if err := listHostsRequest.Walk(walkFunc); err != nil {
		return eris.Wrap(err, "listHostsRequest.Walk")
	}

	bar.SetTotal(-1, true)

	return nil
}

func filterRoutine(ctx context.Context, hostOutChan chan computeadmin.HostV2, instanceOutChan chan Instance) error {
	log := log.With().Str("Function", "filterRoutine").Logger()
	defer log.Debug().Msg("Return")

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < NumThreadsFilterRoutine; i++ {
		threadNum := i

		g.Go(func() error {
			log := log.With().Int("Thread", threadNum).Logger()

			for {
				select {
				case <-ctx.Done():
					log.Debug().Msg("Done")

					return nil
				case host, ok := <-hostOutChan:
					if !ok {
						return nil
					}

					log.Debug().Int("Host OUT CHAN", len(hostOutChan)).Msg("")

					filterHostsWithIsland(ctx, host, instanceOutChan)
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		return eris.Wrap(err, "g.Wait")
	}

	return nil
}

func filterHostsWithIsland(ctx context.Context, host computeadmin.HostV2, instanceOutChan chan Instance) {
	islandID, islandIDOK := host.GetHPCIslandIDOk()
	instanceID, instanceIDOK := host.GetInstanceIDOk()

	if islandIDOK && instanceIDOK {
		singleInstance, response, err := fetchInstance(ctx, *instanceID)
		if err != nil {
			singleInstance = &computeadmin.InstanceV2{}
		}

		instanceOutChan <- Instance{
			RackID:             host.GetRackID(),
			HostID:             host.GetID(),
			InstanceID:         *instanceID,
			Shape:              host.GetShape(),
			TenancyID:          singleInstance.GetTenantID(),
			CompartmentID:      singleInstance.GetCompartmentID(),
			IslandID:           *islandID,
			PoolName:           host.GetPoolName(),
			BareMetalClusterID: singleInstance.GetBareMetalClusterID(),
			HostReservation:    host.GetHostReservationID(),
			response:           response,
		}
	}
}

func fetchInstance(ctx context.Context, instanceID string) (*computeadmin.InstanceV2, int, error) {
	computeAdminAPIv2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return nil, 500, eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	instance, response, err := computeAdminAPIv2.GetInstanceV2(ctx, instanceID).Execute()
	if err != nil {
		if response == nil {
			return nil, 500, eris.Wrap(err, "fetchInstance.GetInstanceV2")
		}

		return nil, response.StatusCode, eris.Wrap(err, "fetchInstance.GetInstanceV2")
	}

	return instance, response.StatusCode, nil
}
