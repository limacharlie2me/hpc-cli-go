package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"golang.org/x/sync/errgroup"
)

var (
	NumThreadsFilterRoutine     = runtime.NumCPU()
	NumThreadsProcessRoutine    = runtime.NumCPU()
	NumThreadsComputeManagement = runtime.NumCPU()

	p *progressbar
)

func hpc(ctx context.Context) error { //nolint:funlen
	p = newProgressbar()

	start := time.Now()

	params := map[SearchType]string{
		SearchTypeByClusterNetworkID:   cli.ClusterNetworkOCID,
		SearchTypeByHpcIslandID:        cli.HpcIslandID,
		SearchTypeByBareMetalClusterID: cli.BareMetalClusterOCID,
		SearchTypeByRackID:             cli.RackSKID,
		SearchTypeByInstanceID:         cli.InstanceOCID,
	}

	searchType, id, err := checkArgs(params)
	if err != nil {
		fmt.Println("[ERROR]", err)

		os.Exit(1)
	}

	if err := validateSearchID(ctx, id, searchType); err != nil {
		return eris.Wrapf(err, "validateSearchID: %s %s", id, searchType)
	}

	clusterNetworks, hpcInstances, err := fetch(ctx)
	if err != nil {
		return eris.Wrap(err, "fetch")
	}

	// Prints any warnings that may have been generated
	warnings.Range(func(key, value any) bool {
		fmt.Printf("%v %v\n", value, key)

		return true
	})

	// Returns broken HPC hosts, if --list-broken-hosts flag set
	if cli.ListBroken {
		hpcInstances = filterBrokenHosts(hpcInstances)
	}

	hpcInstances = combineData(hpcInstances, clusterNetworks)
	hpcInstances = filter(hpcInstances, id, searchType)

	if cli.HealthCheck {
		result, err := process(ctx, hpcInstances)
		if err != nil {
			return eris.Wrap(err, "process")
		}

		hpcInstances = result
	}

	if err := displayTenancyInformation(hpcInstances); err != nil {
		return eris.Wrap(err, "displayTenancyInformation")
	}

	duration := time.Since(start)
	fmt.Printf("\nExecution time of %s\n", duration)

	return nil
}

func fetch(ctx context.Context) ([]ClusterNetwork, []Instance, error) {
	var (
		clusterNetworks []ClusterNetwork
		hpcInstances    []Instance
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error

		hpcInstances, err = compute(ctx)
		if err != nil {
			return eris.Wrapf(err, "[ERROR]")
		}

		return nil
	})

	g.Go(func() error {
		if !cli.ListBroken {
			var err error

			clusterNetworks, err = computeManagement(ctx)
			if err != nil {
				return eris.Wrapf(err, "[ERROR]")
			}

			return nil
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, nil, eris.Wrapf(err, "errGroup error")
	}

	return clusterNetworks, hpcInstances, nil
}

func combineData(instances []Instance, clusterNetworks []ClusterNetwork) []Instance {
	var hpcInstances, networkNoInstance []Instance

	lonelyNetworks := make(map[string]bool)

	for _, instance := range instances {
		if cli.ListBroken {
			hpcInstances = append(hpcInstances, instance)
		} else {
			for _, clusterNetwork := range clusterNetworks {
				if strings.Compare(instance.BareMetalClusterID, clusterNetwork.clusterBmOCID) == 0 {
					lonelyNetworks[instance.BareMetalClusterID] = true
					instance.ClusterNetworkDisplayName = clusterNetwork.clusterDisplayName
					instance.ClusterNetworkOCID = clusterNetwork.clusterOCID
					instance.InstancePoolOCID = clusterNetwork.InstancePoolOCID
					instance.InstanceConfiguration = clusterNetwork.InstanceConfiguration
				}
			}
			hpcInstances = append(hpcInstances, instance)
		}
	}

	for _, network := range clusterNetworks {
		if !lonelyNetworks[network.clusterBmOCID] {
			networkNoInstance = append(networkNoInstance, Instance{
				BareMetalClusterID:        network.clusterBmOCID,
				ClusterNetworkDisplayName: network.clusterDisplayName,
				ClusterNetworkOCID:        network.clusterOCID,
				InstancePoolOCID:          network.InstancePoolOCID,
				InstanceConfiguration:     network.InstanceConfiguration,
				HostReservation:           network.HostReservationOCID,
			})
		}
	}

	hpcInstances = append(hpcInstances, networkNoInstance...)

	return hpcInstances
}
