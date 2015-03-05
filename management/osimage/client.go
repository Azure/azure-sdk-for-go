package osimage

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureImageListURL    = "services/images"
	errInvalidImage      = "Can not find image %s in specified subscription, please specify another image name."
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new OsImageClient from an Azure client
func NewClient(client management.Client) OsImageClient {
	return OsImageClient{client: client}
}

func (self OsImageClient) GetImageList() ([]OSImage, error) {
	imageList := imageList{}

	response, err := self.client.SendAzureGetRequest(azureImageListURL)
	if err != nil {
		return imageList.OSImages, err
	}

	if err = xml.Unmarshal(response, &imageList); err != nil {
		return imageList.OSImages, err
	}

	return imageList.OSImages, err
}
