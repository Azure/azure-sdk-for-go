package virtualmachine

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/azure"
)

//VmClient is used to manage operations on Azure Virtual Machines
type VmClient struct {
	client *azure.Client
}

type VMDeployment struct {
	XMLName          xml.Name `xml:"Deployment"`
	Xmlns            string   `xml:"xmlns,attr"`
	Name             string
	DeploymentSlot   string
	Status           string `xml:",omitempty"`
	Label            string
	Url              string `xml:",omitempty"`
	RoleList         RoleList
	RoleInstanceList RoleInstanceList `xml:",omitempty"`
	VirtualIPs       VirtualIPs       `xml:",omitempty"`
}

type RoleList struct {
	Role []*Role
}

type RoleInstanceList struct {
	RoleInstance []*RoleInstance
}

type RoleInstance struct {
	RoleName          string
	InstanceName      string
	InstanceStatus    string
	InstanceSize      string
	PowerState        string
	IpAddress         string
	InstanceEndpoints InstanceEndpoints `xml:",omitempty"`
}

type InstanceEndpoints struct {
	InstanceEndpoint []InstanceEndpoint
}

type InstanceEndpoint struct {
	Name       string
	Vip        string
	PublicPort int
	LocalPort  int
	Protocol   string
}

type Role struct {
	RoleName                    string
	RoleType                    string
	ConfigurationSets           ConfigurationSets
	ResourceExtensionReferences ResourceExtensionReferences `xml:",omitempty"`
	OSVirtualHardDisk           OSVirtualHardDisk
	RoleSize                    string
	ProvisionGuestAgent         bool
	UseCertAuth                 bool   `xml:"-"`
	CertPath                    string `xml:"-"`
}

type ConfigurationSets struct {
	ConfigurationSet []ConfigurationSet
}

type ResourceExtensionReferences struct {
	ResourceExtensionReference []ResourceExtensionReference
}

type InputEndpoints struct {
	InputEndpoint []InputEndpoint
}

type ResourceExtensionReference struct {
	ReferenceName                    string
	Publisher                        string
	Name                             string
	Version                          string
	ResourceExtensionParameterValues ResourceExtensionParameterValues `xml:",omitempty"`
	State                            string
}

type ResourceExtensionParameterValues struct {
	ResourceExtensionParameterValue []ResourceExtensionParameter
}

type ResourceExtensionParameter struct {
	Key   string
	Value string
	Type  string
}

type OSVirtualHardDisk struct {
	MediaLink       string
	SourceImageName string
	HostCaching     string `xml:",omitempty"`
	DiskName        string `xml:",omitempty"`
	OS              string `xml:",omitempty"`
}

type ConfigurationSet struct {
	ConfigurationSetType             string
	HostName                         string `xml:",omitempty"`
	UserName                         string `xml:",omitempty"`
	UserPassword                     string `xml:",omitempty"`
	DisableSshPasswordAuthentication bool
	InputEndpoints                   InputEndpoints `xml:",omitempty"`
	SSH                              SSH            `xml:",omitempty"`
	CustomData                       string         `xml:",omitempty"`
}

type SSH struct {
	PublicKeys PublicKeyList
}

type PublicKeyList struct {
	PublicKey []PublicKey
}

type PublicKey struct {
	Fingerprint string
	Path        string
}

type InputEndpoint struct {
	LocalPort int
	Name      string
	Port      int
	Protocol  string
	Vip       string
}

type ServiceCertificate struct {
	XMLName           xml.Name `xml:"CertificateFile"`
	Xmlns             string   `xml:"xmlns,attr"`
	Data              string
	CertificateFormat string
	Password          string `xml:",omitempty"`
}

type StartRoleOperation struct {
	Xmlns         string `xml:"xmlns,attr"`
	OperationType string
}

type ShutdownRoleOperation struct {
	Xmlns         string `xml:"xmlns,attr"`
	OperationType string
}

type RestartRoleOperation struct {
	Xmlns         string `xml:"xmlns,attr"`
	OperationType string
}

type RoleSizeList struct {
	XMLName   xml.Name   `xml:"RoleSizes"`
	Xmlns     string     `xml:"xmlns,attr"`
	RoleSizes []RoleSize `xml:"RoleSize"`
}

type RoleSize struct {
	Name                               string
	Label                              string
	Cores                              int
	MemoryInMb                         int
	SupportedByWebWorkerRoles          bool
	SupportedByVirtualMachines         bool
	MaxDataDiskCount                   int
	WebWorkerResourceDiskSizeInMb      int
	VirtualMachineResourceDiskSizeInMb int
}

type VirtualIPs struct {
	VirtualIP []VirtualIP
}

type VirtualIP struct {
	Address string
}

type dockerPublicConfig struct {
	DockerPort int `json:"dockerport"`
	Version    int `json:"version"`
}
