package virtualnetwork

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
)

const xmlNamespace = "http://schemas.microsoft.com/ServiceHosting/2011/07/NetworkConfiguration"
const xmlNamespaceXsd = "http://www.w3.org/2001/XMLSchema"
const xmlNamespaceXsi = "http://www.w3.org/2001/XMLSchema-instance"

//VnetClient is used to manage operations on Azure Virtual Networks
type VnetClient struct {
	client *azure.Client
}

//NetworkConfiguration represents the network configuration for an entire Azure
//subscription. TODO: Nicer builder methods for these that abstract away the
//underlying structure
type NetworkConfiguration struct {
	XMLName         xml.Name                    `xml:"NetworkConfiguration"`
	XmlNamespaceXsd string                      `xml:"xmlns:xsd,attr"`
	XmlNamespaceXsi string                      `xml:"xmlns:xsi,attr"`
	Xmlns           string                      `xml:"xmlns,attr"`
	Configuration   VirtualNetworkConfiguration `xml:"VirtualNetworkConfiguration"`
}

//NewNetworkConfiguration creates a new empty NetworkConfiguration structure for
//further configuration. The XML namespaces are set correctly.
func (client *VnetClient) NewNetworkConfiguration() NetworkConfiguration {
	networkConfiguration := NetworkConfiguration{}
	networkConfiguration.setXmlNamespaces()
	return networkConfiguration
}

//setXmlNamespaces ensure that all of the required namespaces are set. It should
//be called prior to marshalling the structure to XML for use with the Azure REST
//endpoint. It is used internally prior to submitting requests, but since it is
//idempotent there is no harm in repeat calls.
func (self *NetworkConfiguration) setXmlNamespaces() {
	self.XmlNamespaceXsd = xmlNamespaceXsd
	self.XmlNamespaceXsi = xmlNamespaceXsi
	self.Xmlns = xmlNamespace
}

type VirtualNetworkConfiguration struct {
	Dns                 Dns                  `xml:"Dns,omitempty"`
	LocalNetworkSites   []LocalNetworkSite   `xml:"LocalNetworkSites>LocalNetworkSite"`
	VirtualNetworkSites []VirtualNetworkSite `xml:"VirtualNetworkSites>VirtualNetworkSite"`
}

type Dns struct {
	DnsServers []DnsServer `xml:"DnsServers,omitempty>DnsServer,omitempty"`
}

type DnsServer struct {
	XMLName   xml.Name `xml:"DnsServer"`
	Name      string   `xml:"name,attr"`
	IPAddress string   `xml:"IPAddress,attr"`
}

type DnsServerRef struct {
	Name string `xml:"name,attr"`
}

type VirtualNetworkSite struct {
	Name          string         `xml:"name,attr"`
	Location      string         `xml:"Location,attr"`
	AddressSpace  AddressSpace   `xml:"AddressSpace"`
	Subnets       []Subnet       `xml:"Subnets>Subnet"`
	DnsServersRef []DnsServerRef `xml:"DnsServersRef,omitempty>DnsServerRef"`
}

type LocalNetworkSite struct {
	Name              string `xml:"name,attr"`
	VPNGatewayAddress string
	AddressSpace      AddressSpace
}

type AddressSpace struct {
	AddressPrefix []string
}

type Subnet struct {
	Name          string `xml:"name,attr"`
	AddressPrefix string
}
