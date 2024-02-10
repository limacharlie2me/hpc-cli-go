package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	sashay "bitbucket.oci.oraclecorp.com/one-punch/sashay2"
	"bitbucket.oci.oraclecorp.com/one-punch/sashay2/computemanagementadmin"
)

var (
	clusterNetworks []ClusterNetwork
	instancePools   []InstancePool
)

func computeManagement(ctx context.Context) ([]ClusterNetwork, error) { //nolint:funlen,cyclop
	log := log.With().Str("Function", "computeManagement").Logger()
	defer log.Debug().Msg("Return")

	clusterNetworkChan := make(chan ClusterNetwork, 128)
	instancePoolChan := make(chan InstancePool, 128)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(clusterNetworkChan)
		defer log.Debug().Msg("Close clusterNetworkChan")

		if err := fetchClusterNetworks(ctx, clusterNetworkChan); err != nil {
			return eris.Wrapf(err, "fetchComputeManagement.fetchClusterNetworks")
		}

		return nil
	})

	for i := 0; i < NumThreadsComputeManagement; i++ {
		threadNum := i

		g.Go(func() error {
			log := log.With().Int("Thread", threadNum).Logger()

			for {
				select {
				case <-ctx.Done():
					log.Debug().Msg("Done")

					return nil
				case clusterNetwork, ok := <-clusterNetworkChan:
					if !ok {
						return nil
					}

					log.Debug().Int("Cluster OUT CHAN", len(clusterNetworkChan)).Msg("")

					instancePoolRoutine(ctx, clusterNetwork.clusterOCID, instancePoolChan)
				case instancePool, ok := <-instancePoolChan:
					if !ok {
						return nil
					}

					log.Debug().Int("IP OUT CHAN", len(instancePoolChan)).Msg("")

					instancePools = append(instancePools, instancePool)
				}
			}
		})
	}

	err := g.Wait()

	close(instancePoolChan)

	if err != nil {
		return nil, eris.Wrapf(err, "errGroup error")
	}

	return joinCim(clusterNetworks, instancePools), nil
}

func fetchClusterNetworks(ctx context.Context, clusterNetworkChan chan ClusterNetwork) error {
	log := log.With().Str("Function", "fetchClusterNetworks").Logger()
	defer log.Debug().Msg("Return")

	computeManagementAdminAPI, err := sashay.ComputeManagementAdminAPI(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeManagementAdminAPI: %s", cli.Region)
	}

	listClusterNetworksRequest := computeManagementAdminAPI.ListClusterNetworksInternal(ctx)
	bar := p.Bar("[INFO] Fetching Cluster Networks: \t")
	bar.SetPriority(0)

	walkFunc := func(network computemanagementadmin.ClusterNetworkSummary, response *http.Response, pageCount int, resourceCount int) error {
		clusterNetwork := ClusterNetwork{
			clusterOCID:        network.ID,
			clusterDisplayName: network.GetDisplayName(),
			compartmentOCID:    network.CompartmentID,
			clusterBmOCID:      network.BareMetalClusterID,
		}

		clusterNetworkChan <- clusterNetwork
		clusterNetworks = append(clusterNetworks, clusterNetwork)

		bar.IncrBy(1)
		bar.SetTotal(bar.Current()+100, false)

		return nil
	}

	if err := listClusterNetworksRequest.Walk(walkFunc); err != nil {
		return eris.Wrap(err, "listClusterNetworksRequest.Walk")
	}

	bar.SetTotal(-1, true)

	return nil
}

var warnings = sync.Map{}

func instancePoolRoutine(ctx context.Context, clusterNetworkID string, ipOutChan chan InstancePool) {
	computeManagementAdminAPI, err := sashay.ComputeManagementAdminAPI(ctx, cli.Region)
	if err != nil {
		warnings.Store(clusterNetworkID, "[WARN] GetClusterNetworkInternal: Error retrieving CN")

		return
	}

	clusterNetwork, response, err := computeManagementAdminAPI.GetClusterNetworkInternal(ctx, clusterNetworkID).Execute()

	if err != nil {
		switch {
		case response != nil && response.StatusCode == 404:
			warnings.Store(clusterNetworkID, "[WARN] GetClusterNetworkInternal: 404")
		case response != nil && response.StatusCode == 500:
			warnings.Store(clusterNetworkID, "[WARN] GetClusterNetworkInternal: 500")
		default:
			warnings.Store(clusterNetworkID, "[WARN] GetClusterNetworkInternal: Error retrieving CN")
		}

		return
	}

	for _, group := range clusterNetwork.ClusterNetworkPlacementGroups {
		for _, pool := range group.InstancePools {
			ipOutChan <- InstancePool{
				instancePoolOCID:      pool.ID,
				clusterNetworkOCID:    clusterNetworkID,
				InstanceConfiguration: pool.GetInstanceConfigurationID(),
			}
		}
	}
}

func joinCim(clusters []ClusterNetwork, instancePools []InstancePool) []ClusterNetwork {
	var clusterNetworks []ClusterNetwork

	for _, cn := range clusters {
		for _, ip := range instancePools {
			if cn.clusterOCID == ip.clusterNetworkOCID {
				cn.InstancePoolOCID = ip.instancePoolOCID
				cn.InstanceConfiguration = ip.InstanceConfiguration
				clusterNetworks = append(clusterNetworks, cn)
			}
		}
	}

	return clusterNetworks
}
