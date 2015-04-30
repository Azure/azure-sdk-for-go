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

func (self DiskClient) DeleteDisk(diskName string, deleteVhdToo bool) (management.OperationId, error) {
	if diskName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "diskName")
	}

	requestURL := fmt.Sprintf(azureVMDiskURL, diskName)
	if deleteVhdToo {
		requestURL += "?comp=media"
	}

	return self.client.SendAzureDeleteRequest(requestURL)
}
