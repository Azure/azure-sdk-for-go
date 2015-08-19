// Package osimage provides a client for Operating System Images.
package osimage

import (
	"encoding/xml"
	"errors"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureImageListURL    = "services/images"
	errInvalidImage      = "Can not find image %s in specified subscription, please specify another image name."
	errParamNotSpecified = "Parameter %s is not specified."
	errImageNotFound     = "Filter too restrictive, no images left?"
)

// NewClient is used to instantiate a new OSImageClient from an Azure client.
func NewClient(client management.Client) OSImageClient {
	return OSImageClient{client: client}
}

func (c OSImageClient) ListOSImages() (ListOSImagesResponse, error) {
	var l ListOSImagesResponse

	response, err := c.client.SendAzureGetRequest(azureImageListURL)
	if err != nil {
		return l, err
	}

	err = xml.Unmarshal(response, &l)
	return l, err
}

func (c OSImageClient) GetOSImage(filter func(OSImage) bool) (OSImage, error) {
	allimages, err := c.ListOSImages()
	if err != nil {
		return OSImage{}, err
	}
	filtered := []OSImage{}
	for _, im := range allimages.OSImages {
		if filter(im) {
			filtered = append(filtered, im)
		}
	}
	if len(filtered) == 0 {
		return OSImage{}, errors.New(errImageNotFound)
	}

	image := filtered[0]
	for _, im := range filtered {
		if im.PublishedDate > image.PublishedDate {
			image = im
		}
	}

	return image, nil
}
