package storageservice

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
)

const (
	azureStorageServiceListURL = "services/storageservices"
	azureStorageServiceURL     = "services/storageservices/%s"

	azureXmlns = "http://schemas.microsoft.com/windowsazure"

	errBlobEndpointNotFound = "Blob endpoint was not found in storage serice %s"
	errParamNotSpecified    = "Parameter %s is not specified."
)

//StorageServiceClient is used to manage operations on Azure Storage
type StorageServiceClient struct {
	client *azure.Client
}

//NewClient is used to instantiate a new StorageServiceClient from an Azure client
func NewClient(self *azure.Client) *StorageServiceClient {
	return &StorageServiceClient{client: self}
}

func (self *StorageServiceClient) GetStorageServiceList() (*StorageServiceList, error) {
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

func (self *StorageServiceClient) GetStorageServiceByName(serviceName string) (*StorageService, error) {
	if len(serviceName) == 0 {
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

func (self *StorageServiceClient) GetStorageServiceByLocation(location string) (*StorageService, error) {
	if len(location) == 0 {
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

func (self *StorageServiceClient) CreateStorageService(name, location string) (*StorageService, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf(errParamNotSpecified, "name")
	}
	if len(location) == 0 {
		return nil, fmt.Errorf(errParamNotSpecified, "location")
	}

	storageDeploymentConfig := self.createStorageServiceDeploymentConf(name, location)
	deploymentBytes, err := xml.Marshal(storageDeploymentConfig)
	if err != nil {
		return nil, err
	}

	requestId, err := self.client.SendAzurePostRequest(azureStorageServiceListURL, deploymentBytes)
	if err != nil {
		return nil, err
	}

	self.client.WaitAsyncOperation(requestId)
	storageService, err := self.GetStorageServiceByName(storageDeploymentConfig.ServiceName)
	if err != nil {
		return nil, err
	}

	return storageService, nil
}

func (self *StorageServiceClient) GetBlobEndpoint(storageService *StorageService) (string, error) {
	for _, endpoint := range storageService.StorageServiceProperties.Endpoints {
		if !strings.Contains(endpoint, ".blob.core") {
			continue
		}

		return endpoint, nil
	}

	return "", errors.New(fmt.Sprintf(errBlobEndpointNotFound, storageService.ServiceName))
}

func (self *StorageServiceClient) createStorageServiceDeploymentConf(name, location string) StorageServiceDeployment {
	storageServiceDeployment := StorageServiceDeployment{}

	storageServiceDeployment.ServiceName = name
	label := base64.StdEncoding.EncodeToString([]byte(name))
	storageServiceDeployment.Label = label
	storageServiceDeployment.Location = location
	storageServiceDeployment.Xmlns = azureXmlns

	return storageServiceDeployment
}
