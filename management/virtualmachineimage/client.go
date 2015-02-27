package virtualmachineimage

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureImageListURL    = "services/images"
	errInvalidImage      = "Can not find image %s in specified subscription, please specify another image name."
	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new ImageClient from an Azure client
func NewClient(client management.Client) ImageClient {
	return ImageClient{client: client}
}

func (self ImageClient) GetImageList() (ImageList, error) {
	imageList := ImageList{}

	response, err := self.client.SendAzureGetRequest(azureImageListURL)
	if err != nil {
		return imageList, err
	}

	err = xml.Unmarshal(response, &imageList)
	if err != nil {
		return imageList, err
	}

	return imageList, err
}

func (self ImageClient) ResolveImageName(imageName string) error {
	if imageName == "" {
		return fmt.Errorf(errParamNotSpecified, "imageName")
	}

	imageList, err := self.GetImageList()
	if err != nil {
		return err
	}

	for _, image := range imageList.OSImages {
		if image.Name != imageName && image.Label != imageName {
			continue
		}

		return nil
	}

	return errors.New(fmt.Sprintf(errInvalidImage, imageName))
}
