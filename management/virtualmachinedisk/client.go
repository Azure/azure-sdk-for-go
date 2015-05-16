// Package virtualmachinedisk provides a client for Virtual Machine Disks.
package virtualmachinedisk

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureVMDiskURL       = "services/disks/%s"
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new DiskClient from an Azure client
func NewClient(client management.Client) DiskClient {
	return DiskClient{client: client}
}

func (c DiskClient) DeleteDisk(diskName string, deleteVhdToo bool) (management.OperationID, error) {
	if diskName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "diskName")
	}

	requestURL := fmt.Sprintf(azureVMDiskURL, diskName)
	if deleteVhdToo {
		requestURL += "?comp=media"
	}

	return c.client.SendAzureDeleteRequest(requestURL)
}
