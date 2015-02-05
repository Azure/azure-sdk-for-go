package azure

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	deleteAzureHostedServiceURL       = "services/hostedservices/%s?comp=media"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s?comp=media"

	invalidDnsLengthError = "The DNS name must be between 3 and 25 characters."
)

//HostedService is used to manage operations on Azure Hosted Services
type HostedService struct {
	client *Client
}

//HostedService is used to return a handle to the HostedService API
func (client *Client) HostedService() *HostedService {
	return &HostedService{client: client}
}

func (self *HostedService) CreateHostedService(dnsName, location string, reverseDnsFqdn string) (string, error) {
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

	err = self.client.Location().ResolveLocation(location)
	if err != nil {
		return "", err
	}

	hostedServiceDeployment := self.createHostedServiceDeploymentConfig(dnsName, location, reverseDnsFqdn)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		return "", err
	}

	requestURL := azureHostedServiceListURL
	requestId, err := self.client.sendAzurePostRequest(requestURL, hostedServiceBytes)
	if err != nil {
		return "", err
	}

	return requestId, nil
}

func (self *HostedService) CheckHostedServiceNameAvailability(dnsName string) (bool, string, error) {
	if len(dnsName) == 0 {
		return false, "", fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := self.verifyDNSName(dnsName)
	if err != nil {
		return false, "", err
	}

	requestURL := fmt.Sprintf(azureHostedServiceAvailabilityURL, dnsName)
	response, err := self.client.sendAzureGetRequest(requestURL)
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

func (self *HostedService) DeleteHostedService(dnsName string) error {
	if len(dnsName) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := self.verifyDNSName(dnsName)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(deleteAzureHostedServiceURL, dnsName)
	requestId, err := self.client.sendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	self.client.waitAsyncOperation(requestId)
	return nil
}

func (self *HostedService) createHostedServiceDeploymentConfig(dnsName, location string, reverseDnsFqdn string) HostedServiceDeployment {
	deployment := HostedServiceDeployment{}
	deployment.ServiceName = dnsName
	label := base64.StdEncoding.EncodeToString([]byte(dnsName))
	deployment.Label = label
	deployment.Location = location
	deployment.ReverseDnsFqdn = reverseDnsFqdn
	deployment.Xmlns = azureXmlns

	return deployment
}

func (self *HostedService) verifyDNSName(dns string) error {
	if len(dns) < 3 || len(dns) > 25 {
		return fmt.Errorf(invalidDnsLengthError)
	}

	return nil
}
