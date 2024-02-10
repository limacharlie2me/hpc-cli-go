package main

import (
	"context"
	"fmt"

	"github.com/rotisserie/eris"

	sashay "bitbucket.oci.oraclecorp.com/one-punch/sashay2"
)

func validateSearchID(ctx context.Context, id string, searchType SearchType) error {
	switch searchType {
	case SearchTypeByClusterNetworkID:
		fmt.Printf("[INFO] Validating OCID %s\n", id)

		return validateClusterNetworkOCID(ctx, id)
	case SearchTypeByHpcIslandID:
		fmt.Printf("[INFO] Validating IslandID %s\n", id)

		return validateIslandID(ctx, id)
	case SearchTypeByBareMetalClusterID:
		fmt.Printf("[INFO] Validating OCID %s\n", id)

		return validateBareMetalClusterOCID(ctx, id)
	case SearchTypeByRackID:
		fmt.Printf("[INFO] Validating RackID %s\n", id)

		return validateRackID(ctx, id)
	case SearchTypeByInstanceID:
		fmt.Printf("[INFO] Validating OCID %s\n", id)

		return validateInstanceOCID(ctx, id)
	case SearchTypeAll:
		fmt.Print("[INFO] Listing all HPC instances in region\n")

		return nil
	case SearchTypeNone:
		return nil
	default:
		return nil
	}
}

func checkArgs(params map[SearchType]string) (SearchType, string, error) {
	if cli.ListALL {
		return SearchTypeAll, "", nil
	}

	var (
		args          bool
		searchID      string
		setSearchType SearchType
	)

	argCount := 0

	// Checks for more than one search term
	// Sets the search ID and search type
	for key, value := range params {
		if len(value) > 0 {
			argCount++

			args = true

			searchID = value

			setSearchType = key
		}
	}

	// Exit if no search term has been entered
	if !args {
		return SearchTypeNone, "", eris.Wrap(eris.New("Please enter a search term. Example --hpc-island-id=bldg3-block3"), "No search term entered")
	}

	// Exit if more than one search term has been entered
	if argCount > 1 {
		return SearchTypeNone, "", eris.Wrap(eris.New("You've entered too many search terms. Enter --cluster-network OR -- bare-metal-cluster OR --rack-id"), "Too many terms have been entered")
	}

	return setSearchType, searchID, nil
}

func validateClusterNetworkOCID(ctx context.Context, eval string) error {
	computeManagementAdminAPI, err := sashay.ComputeManagementAdminAPI(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeManagementAdminAPI: %s", cli.Region)
	}

	if _, _, err := computeManagementAdminAPI.GetClusterNetworkInternal(ctx, eval).Execute(); err != nil {
		return eris.Wrapf(err, "Cluster Network OCID not found: %s", eval)
	}

	return nil
}

func validateIslandID(ctx context.Context, eval string) error {
	computeAdminAPIv2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	if _, _, err := computeAdminAPIv2.GetHPCIsland(ctx, eval).Execute(); err != nil {
		return eris.Wrap(err, "The HPC Island ID you've entered cannot be found")
	}

	return nil
}

func validateBareMetalClusterOCID(ctx context.Context, eval string) error {
	computeAdminAPIv2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	if _, _, err := computeAdminAPIv2.GetBareMetalCluster(ctx, eval).Execute(); err != nil {
		return eris.Wrapf(err, "Bare Metal Cluster OCID not found: %s", eval)
	}

	return nil
}

func validateRackID(ctx context.Context, eval string) error {
	computeAdminAPIv2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	if _, _, err := computeAdminAPIv2.GetRack(ctx, eval).Execute(); err != nil {
		return eris.Wrapf(err, "Rack ID not found: %s", eval)
	}

	return nil
}

func validateInstanceOCID(ctx context.Context, eval string) error {
	computeAdminAPIv2, err := sashay.ComputeAdminAPIv2(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "ComputeAdminAPIv2: %s", cli.Region)
	}

	if _, _, err := computeAdminAPIv2.GetInstanceV2(ctx, eval).Execute(); err != nil {
		return eris.Wrapf(err, "Instance OCID not found: %s", eval)
	}

	return nil
}
