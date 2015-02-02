package vnetClient

import (
	"encoding/xml"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)

const (
	azureNetworkConfigurationURL = "services/networking/media"
)

//GetVirtualNetworkConfiguration retreives the current virtual network
//configuration for the currently active subscription. Note that the
//underlying Azure API means that network related operations are not safe
//for running concurrently.
func GetVirtualNetworkConfiguration() (NetworkConfiguration, error) {
	networkConfiguration := NewNetworkConfiguration()
	response, err := azure.SendAzureGetRequest(azureNetworkConfigurationURL)
	if err != nil {
		return networkConfiguration, err
	}

	err = xml.Unmarshal(response, &networkConfiguration)
	if err != nil {
		return networkConfiguration, err
	}

	return networkConfiguration, nil
}

//SetVirtualNetworkConfiguration configures the virtual networks for the
//currently active subscription according to the NetworkConfiguration given.
//Note that the underlying Azure API means that network related operations
//are not safe for running concurrently.
func SetVirtualNetworkConfiguration(networkConfiguration NetworkConfiguration) error {
	networkConfiguration.setXmlNamespaces()
	networkConfigurationBytes, err := xml.Marshal(networkConfiguration)
	if err != nil {
		return err
	}

	requestId, err := azure.SendAzurePutRequest(azureNetworkConfigurationURL, "text/plain", networkConfigurationBytes)
	if err != nil {
		return err
	}

	err = azure.WaitAsyncOperation(requestId)
	return err
}
