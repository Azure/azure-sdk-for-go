package vmClient

import (
	"encoding/xml"
)

type VMDeployment struct {
	XMLName   xml.Name `xml:"Deployment"`
	Xmlns	  	string `xml:"xmlns,attr"`
	Name	  	string
	DeploymentSlot string
	Status 		string	`xml:",omitempty"`
	Label 		string
	Url 		string	`xml:",omitempty"`
	RoleList RoleList
}

type HostedServiceDeployment struct {
	XMLName   xml.Name `xml:"CreateHostedService"`
	Xmlns	  	string `xml:"xmlns,attr"`
	ServiceName		string
	Label			string
	Description		string
	Location		string
}

type RoleList struct {
	Role	[]*Role
}

type Role struct {
	RoleName      string
	RoleType	  string
	ConfigurationSets	ConfigurationSets
	ResourceExtensionReferences ResourceExtensionReferences `xml:",omitempty"`
	OSVirtualHardDisk OSVirtualHardDisk
	RoleSize	  string
	ProvisionGuestAgent bool
	UseCertAuth bool `xml:"-"`
	CertPath string `xml:"-"`
}

type ConfigurationSets struct {
	ConfigurationSet []ConfigurationSet
}

type ResourceExtensionReferences struct {
	ResourceExtensionReference	[]ResourceExtensionReference
}

type InputEndpoints struct {
	InputEndpoint []InputEndpoint
}

type SubnetNames struct {
	SubnetName []SubnetName
}

type ResourceExtensionReference struct {
	ReferenceName		string
	Publisher			string
	Name				string
	Version				string
	ResourceExtensionParameterValues ResourceExtensionParameterValues `xml:",omitempty"`
	State				string
}

type ResourceExtensionParameterValues struct {
	ResourceExtensionParameterValue	[]ResourceExtensionParameter
}

type ResourceExtensionParameter struct {
	Key		    string
	Value		string
	Type		string
}

type OSVirtualHardDisk struct {
	MediaLink			string
	SourceImageName		string
	HostCaching			string	`xml:",omitempty"`
	DiskName			string	`xml:",omitempty"`
	OS					string	`xml:",omitempty"`
}

type ConfigurationSet struct {
	ConfigurationSetType	string
	HostName				string	`xml:",omitempty"`
	UserName				string	`xml:",omitempty"`
	UserPassword			string	`xml:",omitempty"`
	DisableSshPasswordAuthentication bool
	InputEndpoints			InputEndpoints	`xml:",omitempty"`
	SubnetNames				SubnetNames		`xml:",omitempty"`
	SSH						SSH		`xml:",omitempty"`
}

type SSH struct {
	PublicKeys PublicKeyList
}

type PublicKeyList struct {
	PublicKey []PublicKey
}

type PublicKey struct {
	Fingerprint			string
	Path				string
}

type InputEndpoint struct {
	LocalPort	int
	Name			string
	Port	int
	Protocol		string
}

//!TODO
type SubnetName struct {
}

type ServiceCertificate struct {
	XMLName   xml.Name `xml:"CertificateFile"`
	Xmlns	  	string `xml:"xmlns,attr"`
	Data	string
	CertificateFormat	string
	Password	string	`xml:",omitempty"`
}

type StartRoleOperation struct {
	Xmlns	  		string `xml:"xmlns,attr"`
	OperationType 	string
}
