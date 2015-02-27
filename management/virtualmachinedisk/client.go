package virtualmachinedisk

import (
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureVMDiskURL       = "services/disks/%s"
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new DiskClient from an Azure client
func NewClient(client management.Client) DiskClient {
	return DiskClient{client: client}
}

func (self DiskClient) DeleteDisk(diskName string) error {
	if diskName == "" {
		return fmt.Errorf(errParamNotSpecified, "diskName")
	}

	requestURL := fmt.Sprintf(azureVMDiskURL, diskName)
	requestId, err := self.client.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	return self.client.WaitAsyncOperation(requestId)
}
