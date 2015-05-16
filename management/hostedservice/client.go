// Package hostedservice provides a client for Hosted Services.
package hostedservice

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
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

func (h HostedServiceClient) CreateHostedService(dnsName, location string, reverseDNSFqdn string, label string, description string) (management.OperationID, error) {
	if dnsName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}
	if location == "" {
		return "", fmt.Errorf(errParamNotSpecified, "location")
	}

	hostedServiceDeployment := h.createHostedServiceDeploymentConfig(dnsName, location, reverseDNSFqdn, label, description)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		return "", err
	}

	requestURL := azureHostedServiceListURL
	return h.client.SendAzurePostRequest(requestURL, hostedServiceBytes)
}

func (h HostedServiceClient) CheckHostedServiceNameAvailability(dnsName string) (bool, string, error) {
	if dnsName == "" {
		return false, "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	requestURL := fmt.Sprintf(azureHostedServiceAvailabilityURL, dnsName)
	response, err := h.client.SendAzureGetRequest(requestURL)
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

func (h HostedServiceClient) DeleteHostedService(dnsName string, deleteDisksAndBlobs bool) (management.OperationID, error) {
	if dnsName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "dnsName")
	}

	requestURL := fmt.Sprintf(getHostedServicePropertiesURL, dnsName)
	if deleteDisksAndBlobs {
		requestURL += "?comp=media"
	}
	return h.client.SendAzureDeleteRequest(requestURL)
}

func (h HostedServiceClient) GetHostedService(name string) (HostedService, error) {
	hostedService := HostedService{}
	if name == "" {
		return hostedService, fmt.Errorf(errParamNotSpecified, "name")
	}

	requestURL := fmt.Sprintf(getHostedServicePropertiesURL, name)
	response, err := h.client.SendAzureGetRequest(requestURL)
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

func (h HostedServiceClient) ListHostedServices() (ListHostedServiceResponse, error) {
	var response ListHostedServiceResponse

	data, err := h.client.SendAzureGetRequest(azureHostedServiceListURL)
	if err != nil {
		return response, err
	}

	err = xml.Unmarshal(data, &response)

	return response, err
}

func (h HostedServiceClient) createHostedServiceDeploymentConfig(dnsName, location string, reverseDNSFqdn string, label string, description string) CreateHostedService {
	encodedLabel := base64.StdEncoding.EncodeToString([]byte(label))
	deployment := CreateHostedService{
		ServiceName:    dnsName,
		Label:          encodedLabel,
		Description:    description,
		Location:       location,
		ReverseDNSFqdn: reverseDNSFqdn,
		Xmlns:          azureXmlns,
	}
	return deployment
}

func (h HostedServiceClient) AddCertificate(dnsName string, certData []byte, certificateFormat CertificateFormat, password string) (management.OperationID, error) {
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
	return h.client.SendAzurePostRequest(requestURL, buffer)
}
