package vmutils

import (
	"fmt"

	vm "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

// ConfigureDeploymentFromRemoteImage configures VM Role to deploy from a remote
// image source. "remoteImageSourceURL" can be any publically accessible URL to
// a VHD file, including but not limited to a SAS Azure Storage blob url. "os"
// needs to be either "Linux" or "Windows". "label" is optional.
func ConfigureDeploymentFromRemoteImage(role *vm.Role, remoteImageSourceURL, os, newDiskName, destinationVhdStorageURL, label string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.OSVirtualHardDisk = &vm.OSVirtualHardDisk{
		RemoteSourceImageLink: remoteImageSourceURL,
		MediaLink:             destinationVhdStorageURL,
		DiskName:              newDiskName,
		OS:                    os,
		DiskLabel:             label,
	}
	return nil
}

// ConfigureDeploymentFromPlatformImage configures VM Role to deploy from a
// platform image. See osimage package for methods to retrieve a list of the
// available platform images. "label" is optional.
func ConfigureDeploymentFromPlatformImage(role *vm.Role, imageName, destinationVhdStorageURL, label string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.OSVirtualHardDisk = &vm.OSVirtualHardDisk{
		SourceImageName: imageName,
		MediaLink:       destinationVhdStorageURL,
	}
	return nil
}

// ConfigureDeploymentFromVMImage configures VM Role to deploy from a previously
// captured VM image.
func ConfigureDeploymentFromVMImage(role *vm.Role, vmImageName, destinationContainerStorageURL string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.VMImageName = vmImageName
	role.MediaLocation = destinationContainerStorageURL
	role.ProvisionGuestAgent = false
	return nil
}

// ConfigureDeploymentFromExistingOSDisk configures VM Role to deploy from an
// existing disk. 'label' is optional.
func ConfigureDeploymentFromExistingOSDisk(role *vm.Role, osDiskName, label string) error {
	role.OSVirtualHardDisk = &vm.OSVirtualHardDisk{
		DiskName:  osDiskName,
		DiskLabel: label,
	}
	return nil
}
