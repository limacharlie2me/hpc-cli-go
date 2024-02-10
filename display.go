package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/rotisserie/eris"
)

func displayTenancyInformation(hpcInstances []Instance) error { //nolint:funlen,cyclop,gocognit
	path, err := os.Getwd()

	if err != nil {
		return eris.Wrap(err, "os.Getwd")
	}

	if cli.Display { //nolint:nestif
		t := tablewriter.NewWriter(os.Stdout)

		t.SetHeader([]string{
			"Rack_ID",
			"Island_ID",
			"Host_ID",
			"Shape",
			"Tenancy_Ocid",
			"Compartment_Ocid",
			"Instance_Ocid",
			"Cluster_Network_Display_Name",
			"CN_Ocid",
			"BMC_Ocid",
			"Instance_Pool",
			"Instance_Config",
			"Host_Reservation"})

		for _, instance := range hpcInstances {
			var (
				rackID             string
				clusterNetworkOCID string
				bmcOCID            string
				ipOCID             string
				icOCID             string
				hostReservation    string
				hostID             string
				tenancyID          string
				compartmentID      string
				instanceID         string
			)

			if instance.ClusterNetworkOCID != "" {
				clusterNetworkOCID = fmt.Sprintf("..%s", instance.ClusterNetworkOCID[len(instance.ClusterNetworkOCID)-6:])
				ipOCID = fmt.Sprintf("..%s", instance.InstancePoolOCID[len(instance.InstancePoolOCID)-6:])
				icOCID = fmt.Sprintf("..%s", instance.InstanceConfiguration[len(instance.InstanceConfiguration)-6:])
			} else {
				clusterNetworkOCID = ""
				ipOCID = ""
				icOCID = ""
			}

			if instance.BareMetalClusterID != "" {
				bmcOCID = fmt.Sprintf("..%s", instance.BareMetalClusterID[len(instance.BareMetalClusterID)-6:])
			}

			if instance.HostReservation != "" {
				hostReservation = fmt.Sprintf("..%s", instance.HostReservation[len(instance.HostReservation)-8:])
			}

			if instance.HostID != "" {
				hostID = fmt.Sprintf("..%s", instance.HostID[len(instance.HostID)-8:])
			}

			if instance.TenancyID != "" {
				tenancyID = fmt.Sprintf("..%s", instance.TenancyID[len(instance.TenancyID)-6:])
			}

			if instance.CompartmentID != "" {
				compartmentID = fmt.Sprintf("..%s", instance.CompartmentID[len(instance.CompartmentID)-6:])
			}

			if instance.RackID != "" {
				rackID = fmt.Sprintf("..%s", instance.RackID[len(instance.RackID)-8:])
			}

			if instance.InstanceID != "" {
				instanceID = fmt.Sprintf("..%s", instance.InstanceID[len(instance.InstanceID)-8:])
			}

			t.Append([]string{
				rackID,
				instance.IslandID,
				hostID,
				instance.Shape,
				tenancyID,
				compartmentID,
				instanceID,
				instance.ClusterNetworkDisplayName,
				clusterNetworkOCID,
				bmcOCID,
				ipOCID,
				icOCID,
				hostReservation,
			})
		}

		t.Render()
		fmt.Printf("[INFO] Displaying %d rows\n", len(hpcInstances))
	}

	file, err := json.MarshalIndent(hpcInstances, "", " ")

	if err != nil {
		return eris.Wrap(err, "jsonMarshalIndent")
	}

	if err := os.WriteFile("hpc_tenancies.json", file, 0600); err != nil {
		return eris.Wrap(err, "Error writing to JSON file")
	}

	fmt.Printf("\n[INFO] The JSON file can be found at %s/hpc_tenancies.json\n", path)

	return nil
}
