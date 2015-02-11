package hostedservice

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
	locationclient "github.com/MSOpenTech/azure-sdk-for-go/azure/location"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	deleteAzureHostedServiceURL       = "services/hostedservices/%s?comp=media"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s?comp=media"
	getHostedServicePropertiesURL     = "services/hostedservices/%s"

	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to return a handle to the HostedService API
func NewClient(client *azure.Client) *HostedServiceClient {
	return &HostedServiceClient{client: client}
}

func (self *HostedServiceClient) CreateHostedService(dnsName, location string, reverseDnsFqdn string, label string, description string) (string, error) {
	if len(dnsName) == 0 {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}
	if len(location) == 0 {
		return "", fmt.Errorf(errParamNotSpecified, "location")
	}

	result, reason, err := self.CheckHostedServiceNameAvailability(dnsName)
	if err != nil {
		return "", err
	}
	if !result {
		return "", fmt.Errorf("%s Hosted service name: %s", reason, dnsName)
	}

	locationClient := locationclient.NewClient(self.client)
	err = locationClient.ResolveLocation(location)
	if err != nil {
		return "", err
	}

	hostedServiceDeployment := self.createHostedServiceDeploymentConfig(dnsName, location, reverseDnsFqdn, label, description)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		return "", err
	}

	requestURL := azureHostedServiceListURL
	requestId, err := self.client.SendAzurePostRequest(requestURL, hostedServiceBytes)
	if err != nil {
		return "", err
	}

	return requestId, nil
}

func (self *HostedServiceClient) CheckHostedServiceNameAvailability(dnsName string) (bool, string, error) {
	if len(dnsName) == 0 {
		return false, "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	requestURL := fmt.Sprintf(azureHostedServiceAvailabilityURL, dnsName)
	response, err := self.client.SendAzureGetRequest(requestURL)
	if err != nil {
		return false, "", err
	}

	availabilityResponse := new(AvailabilityResponse)
	err = xml.Unmarshal(response, availabilityResponse)
	if err != nil {
		return false, "", err
	}

	return availabilityResponse.Result, availabilityResponse.Reason, nil
}

func (self *HostedServiceClient) DeleteHostedService(dnsName string) error {
	if len(dnsName) == 0 {
		return fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	requestURL := fmt.Sprintf(deleteAzureHostedServiceURL, dnsName)
	requestId, err := self.client.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self *HostedServiceClient) GetHostedService(name string) (HostedService, error) {
	hostedService := HostedService{}

	requestURL := fmt.Sprintf(getHostedServicePropertiesURL, name)
	response, err := self.client.SendAzureGetRequest(requestURL)
	if err != nil {
		return hostedService, err
	}

	err = xml.Unmarshal(response, &hostedService)
	if err != nil {
		return hostedService, err
	}

	decodedLabel, err := base64.StdEncoding.DecodeString(hostedService.LabelBase64)
	if err != nil {
		return hostedService, err
	}
	hostedService.Label = string(decodedLabel)
	return hostedService, nil
}

func (self *HostedServiceClient) createHostedServiceDeploymentConfig(dnsName, location string, reverseDnsFqdn string, label string, description string) CreateHostedService {
	encodedLabel := base64.StdEncoding.EncodeToString([]byte(label))
	deployment := CreateHostedService{
		ServiceName:    dnsName,
		Label:          encodedLabel,
		Description:    description,
		Location:       location,
		ReverseDnsFqdn: reverseDnsFqdn,
		Xmlns:          azureXmlns,
	}
	return deployment
}
