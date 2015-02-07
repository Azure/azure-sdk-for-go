package disk

import (
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
)

const (
	azureVMDiskURL       = "services/disks/%s"
	errParamNotSpecified = "Parameter %s is not specified."
)

//DiskClient is used to manage operations on Azure Disks
type DiskClient struct {
	client *azure.Client
}

//NewClient is used to instantiate a new DiskClient from an Azure client
func Disk(client *azure.Client) *DiskClient {
	return &DiskClient{client: client}
}

//Region public methods starts

func (self *DiskClient) DeleteDisk(diskName string) error {
	if len(diskName) == 0 {
		return fmt.Errorf(errParamNotSpecified, "diskName")
	}

	requestURL := fmt.Sprintf(azureVMDiskURL, diskName)
	requestId, err := self.client.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	self.client.WaitAsyncOperation(requestId)
	return nil
}

//Region public methods ends
