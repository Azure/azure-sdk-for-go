package hostedServiceClient

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"

	azure "github.com/MSOpenTech/azure-sdk-for-go"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/locationClient"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	deleteAzureHostedServiceURL       = "services/hostedservices/%s?comp=media"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s?comp=media"

	invalidDnsLengthError  = "The DNS name must be between 3 and 25 characters."
	paramNotSpecifiedError = "Parameter %s is not specified."
)

func CreateHostedService(dnsName, location string, reverseDnsFqdn string) (string, error) {
	if len(dnsName) == 0 {
		return "", fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}
	if len(location) == 0 {
		return "", fmt.Errorf(paramNotSpecifiedError, "location")
	}

	err := verifyDNSName(dnsName)
	if err != nil {
		return "", err
	}

	result, reason, err := CheckHostedServiceNameAvailability(dnsName)
	if err != nil {
		return "", err
	}
	if !result {
		return "", fmt.Errorf("%s Hosted service name: %s", reason, dnsName)
	}

	err = locationClient.ResolveLocation(location)
	if err != nil {
		return "", err
	}

	hostedServiceDeployment := createHostedServiceDeploymentConfig(dnsName, location, reverseDnsFqdn)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		return "", err
	}

	requestURL := azureHostedServiceListURL
	requestId, err := azure.SendAzurePostRequest(requestURL, hostedServiceBytes)
	if err != nil {
		return "", err
	}

	return requestId, nil
}

func CheckHostedServiceNameAvailability(dnsName string) (bool, string, error) {
	if len(dnsName) == 0 {
		return false, "", fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := verifyDNSName(dnsName)
	if err != nil {
		return false, "", err
	}

	requestURL := fmt.Sprintf(azureHostedServiceAvailabilityURL, dnsName)
	response, err := azure.SendAzureGetRequest(requestURL)
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

func DeleteHostedService(dnsName string) error {
	if len(dnsName) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "dnsName")
	}

	err := verifyDNSName(dnsName)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(deleteAzureHostedServiceURL, dnsName)
	requestId, err := azure.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	azure.WaitAsyncOperation(requestId)
	return nil
}

func createHostedServiceDeploymentConfig(dnsName, location string, reverseDnsFqdn string) HostedServiceDeployment {
	deployment := HostedServiceDeployment{}
	deployment.ServiceName = dnsName
	label := base64.StdEncoding.EncodeToString([]byte(dnsName))
	deployment.Label = label
	deployment.Location = location
	deployment.ReverseDnsFqdn = reverseDnsFqdn
	deployment.Xmlns = azureXmlns

	return deployment
}

func verifyDNSName(dns string) error {
	if len(dns) < 3 || len(dns) > 25 {
		return fmt.Errorf(invalidDnsLengthError)
	}

	return nil
}
