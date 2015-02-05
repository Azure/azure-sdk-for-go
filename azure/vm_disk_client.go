package azure

import (
	"fmt"
)

const (
	azureVMDiskURL = "services/disks/%s"
)

type DiskClient struct {
	client *Client
}

func (client *Client) Disk() *DiskClient {
	return &DiskClient{client: client}
}

//Region public methods starts

func (self *DiskClient) DeleteDisk(diskName string) error {
	if len(diskName) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "diskName")
	}

	requestURL := fmt.Sprintf(azureVMDiskURL, diskName)
	requestId, err := self.client.sendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	self.client.waitAsyncOperation(requestId)
	return nil
}

//Region public methods ends
