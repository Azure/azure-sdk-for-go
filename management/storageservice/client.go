package storageservice

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureStorageServiceListURL         = "services/storageservices"
	azureStorageServiceURL             = "services/storageservices/%s"
	azureStorageAccountAvailabilityURL = "services/storageservices/operations/isavailable/%s"

	azureXmlns = "http://schemas.microsoft.com/windowsazure"

	errBlobEndpointNotFound = "Blob endpoint was not found in storage serice %s"
	errParamNotSpecified    = "Parameter %s is not specified."
)

//NewClient is used to instantiate a new StorageServiceClient from an Azure client
func NewClient(self management.Client) StorageServiceClient {
	return StorageServiceClient{client: self}
}

func (self StorageServiceClient) GetStorageServiceList() (*StorageServiceList, error) {
	storageServiceList := new(StorageServiceList)

	response, err := self.client.SendAzureGetRequest(azureStorageServiceListURL)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response, storageServiceList)
	if err != nil {
		return storageServiceList, err
	}

	return storageServiceList, nil
}

func (self StorageServiceClient) GetStorageServiceByName(serviceName string) (*StorageService, error) {
	if serviceName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "serviceName")
	}

	storageService := new(StorageService)
	requestURL := fmt.Sprintf(azureStorageServiceURL, serviceName)
	response, err := self.client.SendAzureGetRequest(requestURL)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response, storageService)
	if err != nil {
		return nil, err
	}

	return storageService, nil
}

func (self StorageServiceClient) GetStorageServiceByLocation(location string) (*StorageService, error) {
	if location == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "location")
	}

	storageService := new(StorageService)
	storageServiceList, err := self.GetStorageServiceList()
	if err != nil {
		return storageService, err
	}

	for _, storageService := range storageServiceList.StorageServices {
		if storageService.StorageServiceProperties.Location != location {
			continue
		}

		return &storageService, nil
	}

	return nil, nil
}

func (self StorageServiceClient) GetStorageServiceKeys(serviceName string) (GetStorageServiceKeysResponse, error) {
	if serviceName == "" {
		return GetStorageServiceKeysResponse{}, fmt.Errorf(errParamNotSpecified, "serviceName")
	}

	requestURL := fmt.Sprintf(azureStorageServiceURL, serviceName)
	data, err := self.client.SendAzureGetRequest(requestURL)
	if err != nil {
		return GetStorageServiceKeysResponse{}, err
	}

	response := GetStorageServiceKeysResponse{}
	err = xml.Unmarshal(data, response)

	return response, err
}

func (self StorageServiceClient) CreateAsync(parameters StorageAccountCreateParameters) (string, error) {
	data, err := xml.Marshal(CreateStorageServiceInput{
		StorageAccountCreateParameters: parameters})
	if err != nil {
		return "", err
	}

	return self.client.SendAzurePostRequest(azureStorageServiceListURL, data)
}

func (self StorageServiceClient) Create(parameters StorageAccountCreateParameters) (*StorageService, error) {
	requestId, err := self.CreateAsync(parameters)
	if err != nil {
		return nil, err
	}

	err = self.client.WaitAsyncOperation(requestId)
	if err != nil {
		return nil, err
	}

	storageService, err := self.GetStorageServiceByName(parameters.ServiceName)
	if err != nil {
		return nil, err
	}

	return storageService, nil
}

func (self StorageServiceClient) GetBlobEndpoint(storageService *StorageService) (string, error) {
	for _, endpoint := range storageService.StorageServiceProperties.Endpoints {
		if !strings.Contains(endpoint, ".blob.core") {
			continue
		}

		return endpoint, nil
	}

	return "", errors.New(fmt.Sprintf(errBlobEndpointNotFound, storageService.ServiceName))
}

// The Check Storage Account Name Availability operation checks to see if the specified storage account name is available, or if it has already been taken.
// See https://msdn.microsoft.com/en-us/library/azure/jj154125.aspx
func (self StorageServiceClient) IsAvailable(name string) (bool, string, error) {
	if name == "" {
		return false, "", fmt.Errorf(errParamNotSpecified, "name")
	}

	requestURL := fmt.Sprintf(azureStorageAccountAvailabilityURL, name)
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
