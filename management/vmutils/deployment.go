package vmutils

import (
	"fmt"

	. "github.com/MSOpenTech/azure-sdk-for-go/management/virtualmachine"
)

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
