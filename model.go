package main

import baremetalstats "bitbucket.oci.oraclecorp.com/one-punch/sashay2/baremetalstats"

type Instance struct {
	RackID                    string `json:"rackId"`
	HostID                    string `json:"hostId"`
	InstanceID                string `json:"instanceId"`
	CompartmentID             string `json:"compartmentId"`
	TenancyID                 string `json:"tenancyId"`
	IslandID                  string `json:"islandId"`
	Shape                     string `json:"shape"`
	PoolName                  string `json:"poolName"`
	BareMetalClusterID        string `json:"bareMetalClusterId,omitempty"`
	ClusterNetworkOCID        string `json:"clusterNetwordId,omitempty"`
	ClusterNetworkDisplayName string `json:"clusterNetworkDisplayName,omitempty"`
	InstancePoolOCID          string `json:"instancePoolId,omitempty"`
	InstanceConfiguration     string `json:"instanceConfigurationId,omitempty"`
	HostReservation           string `json:"hostReservationId,omitempty"`
	HealthCheck               Health `json:"healthChecks,omitempty"`
	response                  int
}

type Health struct {
	Smartnic *baremetalstats.HostSmartNicData `json:"hostSmartNicData,omitempty"`
	Ilom     *baremetalstats.HostData         `json:"ilomDetails,omitempty"`
}

type ClusterNetwork struct {
	clusterOCID           string
	clusterDisplayName    string
	compartmentOCID       string
	clusterBmOCID         string
	InstancePoolOCID      string
	InstanceConfiguration string
	HostReservationOCID   string
}

type InstancePool struct {
	instancePoolOCID      string
	clusterNetworkOCID    string
	InstanceConfiguration string
}

type SearchType uint

const (
	SearchTypeByClusterNetworkID SearchType = iota
	SearchTypeByHpcIslandID
	SearchTypeByBareMetalClusterID
	SearchTypeByRackID
	SearchTypeByInstanceID
	SearchTypeAll
	SearchTypeNone
)

func (search SearchType) String() string {
	switch search {
	case SearchTypeByClusterNetworkID:
		return "ByClusterNetworkID"
	case SearchTypeByHpcIslandID:
		return "ByHpcIslandID"
	case SearchTypeByBareMetalClusterID:
		return "ByBareMetalClusterID"
	case SearchTypeByRackID:
		return "ByRackID"
	case SearchTypeByInstanceID:
		return "ByInstanceID"
	case SearchTypeAll:
		return "All"
	case SearchTypeNone:
		return "None"
	}

	return "UNKNOWN"
}
