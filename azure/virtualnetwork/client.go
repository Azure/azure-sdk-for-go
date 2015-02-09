package virtualnetwork

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
)

const (
	azureNetworkConfigurationURL = "services/networking/media"
)

//VnetClient is used to return a handle to the VnetClient API
func NewClient(client *azure.Client) *VnetClient {
	return &VnetClient{client: client}
}

//GetVirtualNetworkConfiguration retreives the current virtual network
//configuration for the currently active subscription. Note that the
//underlying Azure API means that network related operations are not safe
//for running concurrently.
func (self *VnetClient) GetVirtualNetworkConfiguration() (NetworkConfiguration, error) {
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
func (self *VnetClient) SetVirtualNetworkConfiguration(networkConfiguration NetworkConfiguration) error {
	networkConfiguration.setXmlNamespaces()
	networkConfigurationBytes, err := xml.Marshal(networkConfiguration)
	if err != nil {
		return err
	}

	requestId, err := self.client.SendAzurePutRequest(azureNetworkConfigurationURL, "text/plain", networkConfigurationBytes)
	if err != nil {
		return err
	}

	err = self.client.WaitAsyncOperation(requestId)
	return err
}
