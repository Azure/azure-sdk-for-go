package virtualmachine

import (
	"encoding/xml"
)

const (
	azureResourceExtensionsURL = "services/resourceextensions"
)

// Lists the resource extensions that are available to add to a virtual machine
// See https://msdn.microsoft.com/en-us/library/azure/dn495441.aspx
func (self VirtualMachineClient) GetResourceExtensions() (extensions []ResourceExtension, err error) {
	data, err := self.client.SendAzureGetRequest(azureResourceExtensionsURL)
	if err != nil {
		return extensions, err
	}

	var response ResourceExtensions
	err = xml.Unmarshal(data, &response)
	extensions = response.List
	return
}
