package vmutils

import (
	"fmt"

	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

// Configures vm role to deploy from a remote image source. "remoteImageSourceUrl" can be any publically
// accessible URL to a VHD file, including but not limited to a SAS Azure Storage blob url. "os" needs to be
// either "Linux" or "Windows". "label" is optional.
func ConfigureDeploymentFromRemoteImage(role *Role, remoteImageSourceUrl, os, newDiskName, destinationVhdStorageUrl, label string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.OSVirtualHardDisk = &OSVirtualHardDisk{
		RemoteSourceImageLink: remoteImageSourceUrl,
		MediaLink:             destinationVhdStorageUrl,
		DiskName:              newDiskName,
		OS:                    os,
		DiskLabel:             label,
	}
	return nil
}

// Configures vm role to deploy from a platform image. See osimage package for methods to retrieve a list of
// the available platform images. "label" is optional.
func ConfigureDeploymentFromPlatformImage(role *Role, imageName, destinationVhdStorageUrl, label string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.OSVirtualHardDisk = &OSVirtualHardDisk{
		SourceImageName: imageName,
		MediaLink:       destinationVhdStorageUrl,
	}
	return nil
}

// Configures vm role to deploy from a previously captured VM image.
func ConfigureDeploymentFromVMImage(role *Role, vmImageName, destinationContainerStorageUrl string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.VMImageName = vmImageName
	role.MediaLocation = destinationContainerStorageUrl
	role.ProvisionGuestAgent = false
	return nil
}

// Configures vm role to deploy from an existing disk. 'label' is optional.
func ConfigureDeploymentFromExistingOSDisk(role *Role, osDiskName, label string) error {
	role.OSVirtualHardDisk = &OSVirtualHardDisk{
		DiskName:  osDiskName,
		DiskLabel: label,
	}
	return nil
}
