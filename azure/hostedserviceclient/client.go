package hostedserviceclient

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
	"github.com/MSOpenTech/azure-sdk-for-go/azure/locationclient"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	deleteAzureHostedServiceURL       = "services/hostedservices/%s?comp=media"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s?comp=media"

	paramNotSpecifiedError = "Parameter %s is not specified."
	invalidDnsLengthError  = "The DNS name must be between 3 and 25 characters."
)

//HostedServiceClient is used to manage operations on Azure Hosted Services
type HostedServiceClient struct {
	client *azure.Client
}

//NewClient is used to return a handle to the HostedService API
func NewClient(client *azure.Client) *HostedServiceClient {
	return &HostedServiceClient{client: client}
}

func (self *HostedServiceClient) CreateHostedService(dnsName, location string, reverseDnsFqdn string) (string, error) {
	if len(dnsName) == 0 {
		return "", fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}
	if len(location) == 0 {
		return "", fmt.Errorf(paramNotSpecifiedError, "location")
	}

	err := self.verifyDNSName(dnsName)
	if err != nil {
		return "", err
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

	hostedServiceDeployment := self.createHostedServiceDeploymentConfig(dnsName, location, reverseDnsFqdn)
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
		return false, "", fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := self.verifyDNSName(dnsName)
	if err != nil {
		return false, "", err
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
		return fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := self.verifyDNSName(dnsName)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(deleteAzureHostedServiceURL, dnsName)
	requestId, err := self.client.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	self.client.WaitAsyncOperation(requestId)
	return nil
}

func (self *HostedServiceClient) createHostedServiceDeploymentConfig(dnsName, location string, reverseDnsFqdn string) HostedServiceDeployment {
	deployment := HostedServiceDeployment{}
	deployment.ServiceName = dnsName
	label := base64.StdEncoding.EncodeToString([]byte(dnsName))
	deployment.Label = label
	deployment.Location = location
	deployment.ReverseDnsFqdn = reverseDnsFqdn
	deployment.Xmlns = azureXmlns

	return deployment
}

func (self *HostedServiceClient) verifyDNSName(dns string) error {
	if len(dns) < 3 || len(dns) > 25 {
		return fmt.Errorf(invalidDnsLengthError)
	}

	return nil
}
