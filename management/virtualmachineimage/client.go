package virtualmachineimage

import (
	"encoding/xml"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureImageListURL      = "services/vmimages"
	azureRoleOperationsURL = "services/hostedservices/%s/deployments/%s/roleinstances/%s/operations"
	errParamNotSpecified   = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new Client from an Azure client
func NewClient(client management.Client) Client {
	return Client{client}
}

func (self Client) GetImageList() ([]VMImage, error) {
	imageList := vmImageList{}

	response, err := self.SendAzureGetRequest(azureImageListURL)
	if err != nil {
		return imageList.VMImages, err
	}

	err = xml.Unmarshal(response, &imageList)

	return imageList.VMImages, err
}

func (self Client) Capture(cloudServiceName, deploymentName, roleName string,
	name, label string, osState OSState, parameters CaptureParameters) (management.OperationId, error) {
	if cloudServiceName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "cloudServiceName")
	}
	if deploymentName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "roleName")
	}

	request := CaptureRoleAsVMImageOperation{
		VMImageName:       name,
		VMImageLabel:      label,
		OSState:           osState,
		CaptureParameters: parameters,
	}
	data, err := xml.Marshal(request)
	if err != nil {
		return "", err
	}

	return self.SendAzurePostRequest(fmt.Sprintf(azureRoleOperationsURL,
		cloudServiceName, deploymentName, roleName), data)
}
