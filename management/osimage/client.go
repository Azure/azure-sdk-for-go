// Package osimage provides a client for Operating System Images.
package osimage

import (
	"encoding/xml"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureImageListURL        = "services/images"
	azureImageReplicateURL   = "services/images/%s/replicate"
	azureImageUnreplicateURL = "services/images/%s/unreplicate"
	errInvalidImage          = "Can not find image %s in specified subscription, please specify another image name."
	errParamNotSpecified     = "Parameter %s is not specified."
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

func (c OSImageClient) ReplicateImage(image, offer, sku, version string, regions ...string) (management.OperationID, error) {
	ri := ReplicationInput{
		TargetLocations: regions,
		Offer:           offer,
		Sku:             sku,
		Version:         version,
	}

	data, err := xml.Marshal(&ri)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(azureImageReplicateURL, image)

	return c.client.SendAzurePutRequest(url, "", data)
}

func (c OSImageClient) UnreplicateImage(image string) (management.OperationID, error) {
	url := fmt.Sprintf(azureImageUnreplicateURL, image)
	return c.client.SendAzurePutRequest(url, "", []byte{})
}
