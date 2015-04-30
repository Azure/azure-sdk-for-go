package vmutils

import (
	"fmt"

	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
	vmdisk "github.com/Azure/azure-sdk-for-go/management/virtualmachinedisk"
)

// Adds configuration for a new (empty) data disk
func ConfigureWithNewDataDisk(role *Role, label, destinationVhdStorageUrl string, sizeInGB float64, cachingType vmdisk.HostCachingType) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	appendDataDisk(role, DataVirtualHardDisk{
		DiskLabel:           label,
		HostCaching:         cachingType,
		LogicalDiskSizeInGB: sizeInGB,
		MediaLink:           destinationVhdStorageUrl,
	})
	return nil
}

// Adds configuration for an existing data disk
func ConfigureWithExistingDataDisk(role *Role, diskname string, cachingType vmdisk.HostCachingType) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	appendDataDisk(role, DataVirtualHardDisk{
		DiskName:    diskname,
		HostCaching: cachingType,
	})
	return nil
}

// Adds configuration for adding a vhd in a storage account as a data disk
func ConfigureWithVhdDataDisk(role *Role, sourceVhdStorageUrl string, cachingType vmdisk.HostCachingType) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	appendDataDisk(role, DataVirtualHardDisk{
		SourceMediaLink: sourceVhdStorageUrl,
		HostCaching:     cachingType,
	})
	return nil
}

func appendDataDisk(role *Role, disk DataVirtualHardDisk) {
	if role.DataVirtualHardDisks == nil {
		role.DataVirtualHardDisks = &[]DataVirtualHardDisk{disk}
	} else {
		disk.Lun = len(*role.DataVirtualHardDisks)
		newDisks := append(*role.DataVirtualHardDisks, disk)
		role.DataVirtualHardDisks = &newDisks
	}
}
