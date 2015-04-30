package virtualnetwork

import (
	"encoding/xml"

	"github.com/Azure/azure-sdk-for-go/management"
)

const (
	azureNetworkConfigurationURL = "services/networking/media"
)

//VnetClient is used to return a handle to the VnetClient API
func NewClient(client management.Client) VirtualNetworkClient {
	return VirtualNetworkClient{client: client}
}

//GetVirtualNetworkConfiguration retreives the current virtual network
//configuration for the currently active subscription. Note that the
//underlying Azure API means that network related operations are not safe
//for running concurrently.
func (self VirtualNetworkClient) GetVirtualNetworkConfiguration() (NetworkConfiguration, error) {
	networkConfiguration := self.NewNetworkConfiguration()
	response, err := self.client.SendAzureGetRequest(azureNetworkConfigurationURL)
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
func (self VirtualNetworkClient) SetVirtualNetworkConfiguration(networkConfiguration NetworkConfiguration) error {
	networkConfiguration.setXmlNamespaces()
	networkConfigurationBytes, err := xml.Marshal(networkConfiguration)
	if err != nil {
		return err
	}

	operationId, err := self.client.SendAzurePutRequest(azureNetworkConfigurationURL, "text/plain", networkConfigurationBytes)
	if err != nil {
		return err
	}

	err = self.client.WaitAsyncOperation(operationId)
	return err
}
