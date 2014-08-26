package storageServiceClient

import (
	"fmt"
	"strings"
	"errors"
	"encoding/xml"
	"encoding/base64"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)

func GetStorageServiceList() (*StorageServiceList, error){
	storageServiceList := new(StorageServiceList)

	requestURL := "services/storageservices"
	response, err := azure.SendAzureGetRequest(requestURL)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response, storageServiceList)
	if err != nil {
		return storageServiceList, err
	}

	return storageServiceList, nil
}

func GetStorageServiceByName(serviceName string) (*StorageService, error){
	storageService := new(StorageService)
	requestURL := fmt.Sprintf("services/storageservices/%s", serviceName)
	response, err := azure.SendAzureGetRequest(requestURL)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(response, storageService)
	if err != nil {
		return nil, err
	}

	return storageService, nil
}

func GetStorageServiceByLocation(location string) (*StorageService, error) {
	storageService := new(StorageService)
	storageServiceList, err := GetStorageServiceList()
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

func CreateStorageService(name, location string) (*StorageService, error){
	storageDeploymentConfig := createStorageServiceDeploymentConf(name, location)
	deploymentBytes, err := xml.Marshal(storageDeploymentConfig)
	if err != nil {
		return nil, err
	}

	requestURL := "services/storageservices"
	requestId, err := azure.SendAzurePostRequest(requestURL, deploymentBytes)
	if err != nil {
		return nil, err
	}

	azure.WaitAsyncOperation(requestId)
	storageService, err := GetStorageServiceByName(storageDeploymentConfig.ServiceName)
	if err != nil {
		return nil, err
	}

	return storageService, nil
}

func GetBlobEndpoint(storageService *StorageService) (string, error) {
	for _, endpoint := range storageService.StorageServiceProperties.Endpoints {
		if !strings.Contains(endpoint, ".blob.core") {
			continue
		}

		return endpoint, nil
	}

	return "", errors.New(fmt.Sprintf("Blob endpoint was not found in storage serice %s", storageService.ServiceName))
}

func createStorageServiceDeploymentConf(name, location string) (StorageServiceDeployment){
	storageServiceDeployment := StorageServiceDeployment{}

	storageServiceDeployment.ServiceName = name
	label := base64.StdEncoding.EncodeToString([]byte(name))
	storageServiceDeployment.Label = label
	storageServiceDeployment.Location = location
	storageServiceDeployment.Xmlns = "http://schemas.microsoft.com/windowsazure"

	return storageServiceDeployment
}
