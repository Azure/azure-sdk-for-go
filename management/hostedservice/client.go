package hostedservice

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
	locationclient "github.com/Azure/azure-sdk-for-go/management/location"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s"
	getHostedServicePropertiesURL     = "services/hostedservices/%s"
	azureServiceCertificateURL        = "services/hostedservices/%s/certificates"

	errParamNotSpecified = "Parameter %s is not specified."
)

//NewClient is used to return a handle to the HostedService API
func NewClient(client management.Client) HostedServiceClient {
	return HostedServiceClient{client: client}
}

func (self HostedServiceClient) CreateHostedService(dnsName, location string, reverseDnsFqdn string, label string, description string) (management.OperationId, error) {
	if dnsName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}
	if location == "" {
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
	return self.client.SendAzurePostRequest(requestURL, hostedServiceBytes)
}

func (self HostedServiceClient) CheckHostedServiceNameAvailability(dnsName string) (bool, string, error) {
	if dnsName == "" {
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

func (self HostedServiceClient) DeleteHostedService(dnsName string, deleteDisksAndBlobsToo bool) (management.OperationId, error) {
	if dnsName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	requestURL := fmt.Sprintf(getHostedServicePropertiesURL, dnsName)
	if deleteDisksAndBlobsToo {
		requestURL += "?comp=media"
	}
	return self.client.SendAzureDeleteRequest(requestURL)
}

func (self HostedServiceClient) GetHostedService(name string) (HostedService, error) {
	hostedService := HostedService{}
	if name == "" {
		return hostedService, fmt.Errorf(errParamNotSpecified, "name")
	}

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

func (self HostedServiceClient) ListHostedServices() (ListHostedServiceResponse, error) {
	var response ListHostedServiceResponse

	data, err := self.client.SendAzureGetRequest(azureHostedServiceListURL)
	if err != nil {
		return response, err
	}

	err = xml.Unmarshal(data, &response)

	return response, err
}

func (self HostedServiceClient) createHostedServiceDeploymentConfig(dnsName, location string, reverseDnsFqdn string, label string, description string) CreateHostedService {
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

func (self HostedServiceClient) AddCertificate(dnsName string, certData []byte, certificateFormat CertificateFormat, password string) (management.OperationId, error) {
	if dnsName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	certBase64 := base64.StdEncoding.EncodeToString(certData)

	addCertificate := CertificateFile{
		Data:              certBase64,
		CertificateFormat: certificateFormat,
		Password:          password,
		Xmlns:             azureXmlns,
	}
	buffer, err := xml.Marshal(addCertificate)
	if err != nil {
		return "", err
	}

	requestURL := fmt.Sprintf(azureServiceCertificateURL, dnsName)
	return self.client.SendAzurePostRequest(requestURL, buffer)
}
