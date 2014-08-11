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

	requestURL :=  fmt.Sprintf("https://management.core.windows.net/%s/services/storageservices", azure.GetPublishSettings().SubscriptionID)
	output := azure.SendAzureGetRequest(requestURL)

	err := xml.Unmarshal(output, storageServiceList)
	if err != nil {
		return storageServiceList, err
	}

	return storageServiceList, nil
}

func GetStorageServiceByName(serviceName string) (*StorageService){
	storageService := new(StorageService)
	requestURL :=  fmt.Sprintf("https://management.core.windows.net/%s/services/storageservices/%s", azure.GetPublishSettings().SubscriptionID, serviceName)
	output := azure.SendAzureGetRequest(requestURL)

	err := xml.Unmarshal(output, storageService)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	return storageService
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

func CreateStorageService(name, location string) (*StorageService){
	storageDeploymentConfig := createStorageServiceDeploymentConf(name, location)
	deploymentBytes, err := xml.Marshal(storageDeploymentConfig)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	requestURL :=  fmt.Sprintf("https://management.core.windows.net/%s/services/storageservices", azure.GetPublishSettings().SubscriptionID)
	azure.SendAzurePostRequest(requestURL, deploymentBytes)

	storageService := GetStorageServiceByName(storageDeploymentConfig.ServiceName)

	return storageService
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
