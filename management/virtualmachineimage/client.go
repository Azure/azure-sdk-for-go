package virtualmachineimage

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureImageListURL    = "services/vmimages"
	errInvalidImage      = "Can not find image %s in specified subscription, please specify another image name."
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new VmImageClient from an Azure client
func NewClient(client management.Client) VMImageClient {
	return VMImageClient{client: client}
}

func (self VMImageClient) GetImageList() ([]VMImage, error) {
	imageList := vmImageList{}

	response, err := self.client.SendAzureGetRequest(azureImageListURL)
	if err != nil {
		return imageList.VMImages, err
	}

	err = xml.Unmarshal(response, &imageList)

	return imageList.VMImages, err
}
