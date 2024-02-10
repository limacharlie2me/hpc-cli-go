package main

import (
	"fmt"
	"strings"
)

func filterBrokenHosts(instances []Instance) []Instance {
	var brokenHosts []Instance

	fmt.Printf("[INFO] Filtering on hosts in the broken pool\n")

	for _, instance := range instances {
		if strings.Compare(instance.PoolName, "broken_pool") == 0 {
			brokenHosts = append(brokenHosts, instance)
		}
	}

	fmt.Printf("[INFO] %d HPC hosts in the broken pool\n", len(brokenHosts))

	if len(brokenHosts) == 0 {
		fmt.Printf("[INFO] The number of Broken Hosts is 0. Exiting...\n")
	}

	return brokenHosts
}

func filter(hpcInstances []Instance, searchID string, searchType SearchType) []Instance {
	switch searchType {
	case SearchTypeByClusterNetworkID:
		return filterByClusterNetworkOCID(hpcInstances, searchID)
	case SearchTypeByHpcIslandID:
		return filterByIslandID(hpcInstances, searchID)
	case SearchTypeByBareMetalClusterID:
		return filterByBMC(hpcInstances, searchID)
	case SearchTypeByRackID:
		return filterByRackID(hpcInstances, searchID)
	case SearchTypeByInstanceID:
		return filterByInstance(hpcInstances, searchID)
	case SearchTypeAll:
		return hpcInstances
	case SearchTypeNone:
		return nil
	default:
		return nil
	}
}

func filterByClusterNetworkOCID(instances []Instance, clusterNetworkOCID string) []Instance {
	var hpcInstances []Instance

	fmt.Println("[INFO] Filtering by Cluster Network OCID")

	for _, i := range instances {
		if i.ClusterNetworkOCID == clusterNetworkOCID {
			hpcInstances = append(hpcInstances, i)
		}
	}

	fmt.Printf("[INFO] Found %d Instances in Cluster Network %s\n", len(hpcInstances), clusterNetworkOCID)

	return hpcInstances
}

func filterByIslandID(instances []Instance, islandID string) []Instance {
	var hpcInstances []Instance

	fmt.Printf("[INFO] Filtering by HPC Island ID\n")

	for _, i := range instances {
		if i.IslandID == islandID {
			hpcInstances = append(hpcInstances, i)
		}
	}

	fmt.Printf("[INFO] Found %d Instances in HPC Island %s\n", len(hpcInstances), islandID)

	return hpcInstances
}

func filterByBMC(instances []Instance, bmcOCID string) []Instance {
	var hpcInstances []Instance

	fmt.Printf("[INFO] Filtering by Bare Metal Cluster OCID\n")

	for _, i := range instances {
		if i.BareMetalClusterID == bmcOCID {
			hpcInstances = append(hpcInstances, i)
		}
	}

	fmt.Printf("[INFO] Found %d Instances in Bare Metal Cluster %s\n", len(hpcInstances), bmcOCID)

	return hpcInstances
}

func filterByRackID(instances []Instance, rackID string) []Instance {
	var hpcInstances []Instance

	fmt.Printf("[INFO] Filtering by Rack ID\n")

	for _, i := range instances {
		if i.RackID == rackID {
			hpcInstances = append(hpcInstances, i)
		}
	}

	fmt.Printf("[INFO] Found %d Instances on Rack %s\n", len(hpcInstances), rackID)

	return hpcInstances
}

func filterByInstance(instances []Instance, instanceOCID string) []Instance {
	var hpcInstances []Instance

	fmt.Printf("[INFO] Filtering by Instance OCID\n")

	for _, i := range instances {
		if i.InstanceID == instanceOCID {
			hpcInstances = append(hpcInstances, i)
		}
	}

	fmt.Printf("[INFO] Found %d instance\n", len(hpcInstances))

	return hpcInstances
}
