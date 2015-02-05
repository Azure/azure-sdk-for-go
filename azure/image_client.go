package azure

import (
	"encoding/xml"
	"errors"
	"fmt"
)

const (
	azureImageListURL = "services/images"
	invalidImageError = "Can not find image %s in specified subscription, please specify another image name."
)

//ImageClient is used to manage operations on Azure Locations
type ImageClient struct {
	client *Client
}

//Image is used to return a handle to the Image API
func (client *Client) Image() *ImageClient {
	return &ImageClient{client: client}
}

func (self *ImageClient) GetImageList() (ImageList, error) {
	imageList := ImageList{}

	response, err := self.client.sendAzureGetRequest(azureImageListURL)
	if err != nil {
		return imageList, err
	}

	err = xml.Unmarshal(response, &imageList)
	if err != nil {
		return imageList, err
	}

	return imageList, err
}

func (self *ImageClient) ResolveImageName(imageName string) error {
	if len(imageName) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "imageName")
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

	return errors.New(fmt.Sprintf(invalidImageError, imageName))
}
