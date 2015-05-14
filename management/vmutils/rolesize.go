package vmutils

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
	locationclient "github.com/Azure/azure-sdk-for-go/management/location"
	vm "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

// IsRoleSizeValid retrieves the available rolesizes using
// vmclient.GetRoleSizeList() and returns whether that the provided roleSizeName
// is part of that list
func IsRoleSizeValid(vmclient vm.VirtualMachineClient, roleSizeName string) (bool, error) {
	if roleSizeName == "" {
		return false, fmt.Errorf(errParamNotSpecified, "roleSizeName")
	}

	roleSizeList, err := vmclient.GetRoleSizeList()
	if err != nil {
		return false, err
	}

	for _, roleSize := range roleSizeList.RoleSizes {
		if roleSize.Name == roleSizeName {
			return true, nil
		}
	}

	return false, nil
}

// IsRoleSizeAvailableInLocation retrieves all available sizes in the specified
// location using location.GetLocation() and returns whether that the provided
// roleSizeName is part of that list.
func IsRoleSizeAvailableInLocation(managementclient management.Client, location, roleSizeName string) (bool, error) {
	if location == "" {
		return false, fmt.Errorf(errParamNotSpecified, "location")
	}
	if roleSizeName == "" {
		return false, fmt.Errorf(errParamNotSpecified, "roleSizeName")
	}

	locationClient := locationclient.NewClient(managementclient)
	locationInfo, err := locationClient.GetLocation(location)
	if err != nil {
		return false, err
	}

	for _, availableRoleSize := range locationInfo.VirtualMachineRoleSizes {
		if availableRoleSize == roleSizeName {
			return true, nil
		}
	}

	return false, nil
}
