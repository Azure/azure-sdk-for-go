package virtualmachine

import (
	"encoding/xml"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
	vmdisk "github.com/MSOpenTech/azure-sdk-for-go/management/virtualmachinedisk"
)

// VmClient is used to manage operations on Azure Virtual Machines
type VirtualMachineClient struct {
	client management.Client
}

// Type for creating a deployment and Virtual Machine in the deployment based on the specified configuration.
// See https://msdn.microsoft.com/en-us/library/azure/jj157194.aspx
type DeploymentRequest struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure Deployment"`
	// Required parameters:
	Name           string ``            // Specifies a name for the deployment. The deployment name must be unique among other deployments for the cloud service.
	DeploymentSlot string ``            // Specifies the environment in which the Virtual Machine is to be deployed. The only allowable value is Production.
	Label          string ``            // Specifies an identifier for the deployment. The label can be up to 100 characters long. The label can be used for tracking purposes.
	RoleList       []Role `xml:">Role"` // Contains information about the Virtual Machines that are to be deployed.
	// Optional parameters:
	VirtualNetworkName string          `xml:",omitempty"`                         // Specifies the name of an existing virtual network to which the deployment will belong.
	DnsServers         *[]DnsServer    `xml:"Dns>DnsServers>DnsServer,omitempty"` // Contains a list of DNS servers to associate with the Virtual Machine.
	ReservedIPName     string          `xml:",omitempty"`                         // Specifies the name of a reserved IP address that is to be assigned to the deployment.
	LoadBalancers      *[]LoadBalancer `xml:">LoadBalancer,omitempty"`            //Contains a list of internal load balancers that can be assigned to input endpoints.
}

// Type forreceiving deployment information
// See https://msdn.microsoft.com/en-us/library/azure/ee460804.aspx
type DeploymentResponse struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure Deployment"`

	Name                   string
	DeploymentSlot         string
	Status                 DeploymentStatus
	Label                  string
	Url                    string
	Configuration          string
	RoleInstanceList       []RoleInstance `xml:">RoleInstance"`
	UpgradeStatus          UpgradeStatus
	UpgradeDomainCount     int
	RoleList               []Role `xml:">Role"`
	SdkVersion             string
	Locked                 bool
	RollbackAllowed        bool
	CreatedTime            string
	LastModifiedTime       string
	VirtualNetworkName     string
	DnsServers             []DnsServer        `xml:"Dns>DnsServers>DnsServer"`
	LoadBalancers          []LoadBalancer     `xml:">LoadBalancer"`
	ExtendedProperties     []ExtendedProperty `xml:">ExtendedProperty"`
	PersistentVMDowntime   PersistentVMDowntime
	VirtualIPs             []VirtualIP `xml:">VirtualIP`
	ExtensionConfiguration string      // cloud service extensions not fully implemented
	ReservedIPName         string
	InternalDnsSuffix      string
}

type DeploymentStatus string

const (
	DeploymentStatusRunning                DeploymentStatus = "Running"
	DeploymentStatusSuspended              DeploymentStatus = "Suspended"
	DeploymentStatusRunningTransitioning   DeploymentStatus = "RunningTransitioning"
	DeploymentStatusSuspendedTransitioning DeploymentStatus = "SuspendedTransitioning"
	DeploymentStatusStarting               DeploymentStatus = "Starting"
	DeploymentStatusSuspending             DeploymentStatus = "Suspending"
	DeploymentStatusDeploying              DeploymentStatus = "Deploying"
	DeploymentStatusDeleting               DeploymentStatus = "Deleting"
)

type RoleInstance struct {
	RoleName                          string
	InstanceName                      string
	InstanceStatus                    InstanceStatus
	ExtendedInstanceStatus            string
	InstanceUpgradeDomain             int
	InstanceFaultDomain               int
	InstanceSize                      string
	InstanceStateDetails              string
	InstanceErrorCode                 string
	IpAddress                         string
	InstanceEndpoints                 []InstanceEndpoint `xml:">InstanceEndpoint"`
	PowerState                        PowerState
	HostName                          string
	RemoteAccessCertificateThumbprint string
	GuestAgentStatus                  string                    // todo: implement
	ResourceExtensionStatusList       []ResourceExtensionStatus `xml:">ResourceExtensionStatus"`
	PublicIPs                         []PublicIP                `xml:">PublicIP"`
}

type InstanceStatus string

const (
	InstanceStatusUnknown            = "Unknown"
	InstanceStatusCreatingVM         = "CreatingVM"
	InstanceStatusStartingVM         = "StartingVM"
	InstanceStatusCreatingRole       = "CreatingRole"
	InstanceStatusStartingRole       = "StartingRole"
	InstanceStatusReadyRole          = "ReadyRole"
	InstanceStatusBusyRole           = "BusyRole"
	InstanceStatusStoppingRole       = "StoppingRole"
	InstanceStatusStoppingVM         = "StoppingVM"
	InstanceStatusDeletingVM         = "DeletingVM"
	InstanceStatusStoppedVM          = "StoppedVM"
	InstanceStatusRestartingRole     = "RestartingRole"
	InstanceStatusCyclingRole        = "CyclingRole"
	InstanceStatusFailedStartingRole = "FailedStartingRole"
	InstanceStatusFailedStartingVM   = "FailedStartingVM"
	InstanceStatusUnresponsiveRole   = "UnresponsiveRole"
	InstanceStatusStoppedDeallocated = "StoppedDeallocated"
	InstanceStatusPreparing          = "Preparing"
)

type InstanceEndpoint struct {
	Name       string
	Vip        string
	PublicPort int
	LocalPort  int
	Protocol   InputEndpointProtocol
}

type PowerState string

const (
	PowerStateStarting PowerState = "Starting"
	PowerStateStarted  PowerState = "Started"
	PowerStateStopping PowerState = "Stopping"
	PowerStateStopped  PowerState = "Stopped"
	PowerStateUnknown  PowerState = "Unknown"
)

type ResourceExtensionStatus struct {
	HandlerName            string
	Version                string
	Status                 ResourceExtensionState
	Code                   string
	FormattedMessage       FormattedMessage
	ExtensionSettingStatus ExtensionSettingStatus
}

type ResourceExtensionState string

const (
	ResourceExtensionStateInstalling   ResourceExtensionState = "Installing"
	ResourceExtensionStateReady        ResourceExtensionState = "Ready"
	ResourceExtensionStateNotReady     ResourceExtensionState = "NotReady"
	ResourceExtensionStateUnresponsive ResourceExtensionState = "Unresponsive"
)

type FormattedMessage struct {
	Language string
	Message  string
}

type ExtensionSettingStatus struct {
	Timestamp        string
	Name             string
	Operation        string
	Status           ExtensionSettingState
	Code             string
	FormattedMessage FormattedMessage
	SubStatusList    []SubStatus `xml:">SubStatus"`
}

type ExtensionSettingState string

const (
	ExtensionSettingStateTransitioning ExtensionSettingState = "transitioning"
	ExtensionSettingStateError         ExtensionSettingState = "error"
	ExtensionSettingStateSuccess       ExtensionSettingState = "success"
	ExtensionSettingStateWarning       ExtensionSettingState = "warning"
)

type SubStatus struct {
	Name             string
	Status           ExtensionSettingState
	FormattedMessage FormattedMessage
}

type UpgradeStatus struct {
	UpgradeType               UpgradeType
	CurrentUpgradeDomainState CurrentUpgradeDomainState
	CurrentUpgradeDomain      int
}

type UpgradeType string

const (
	UpgradeTypeAuto         UpgradeType = "Auto"
	UpgradeTypeManual       UpgradeType = "Manual"
	UpgradeTypeSimultaneous UpgradeType = "Simultaneous"
)

type CurrentUpgradeDomainState string

const (
	CurrentUpgradeDomainStateBefore CurrentUpgradeDomainState = "Before"
	CurrentUpgradeDomainStateDuring CurrentUpgradeDomainState = "During"
)

type ExtendedProperty struct {
	Name  string
	Value string
}

type PersistentVMDowntime struct {
	StartTime string
	EndTime   string
	Status    string
}

type VirtualIP struct {
	Address        string
	IsReserved     bool
	ReservedIPName string
	Type           IPAddressType
}

// Contains the configuration sets that are used to create virtual machines.
type Role struct {
	RoleName                    string                        `xml:",omitempty"` // Specifies the name for the Virtual Machine.
	RoleType                    string                        `xml:",omitempty"` // Specifies the type of role to use. For Virtual Machines, this must be PersistentVMRole.
	ConfigurationSets           *[]ConfigurationSet           `xml:"ConfigurationSets>ConfigurationSet",omitempty"`
	ResourceExtensionReferences *[]ResourceExtensionReference `xml:"ResourceExtensionReferences>ResourceExtensionReference,omitempty"`
	VMImageName                 string                        `xml:",omitempty"`                                         // Specifies the name of the VM Image that is to be used to create the Virtual Machine. If this element is used, the ConfigurationSets element is not used.
	MediaLocation               string                        `xml:",omitempty"`                                         // Required if the Virtual Machine is being created from a published VM Image. Specifies the location of the VHD file that is created when VMImageName specifies a published VM Image.
	AvailabilitySetName         string                        `xml:",omitempty"`                                         // Specifies the name of a collection of Virtual Machines. Virtual Machines specified in the same availability set are allocated to different nodes to maximize availability.
	DataVirtualHardDisks        *[]DataVirtualHardDisk        `xml:"DataVirtualHardDisks>DataVirtualHardDisk,omitempty"` // Contains the parameters that are used to add a data disk to a Virtual Machine. If you are creating a Virtual Machine by using a VM Image, this element is not used.
	OSVirtualHardDisk           *OSVirtualHardDisk            `xml:",omitempty"`                                         // Contains the parameters that are used to create the operating system disk for a Virtual Machine. If you are creating a Virtual Machine by using a VM Image, this element is not used.
	RoleSize                    string                        `xml:",omitempty"`                                         // Specifies the size of the Virtual Machine. The default size is Small.
	ProvisionGuestAgent         bool                          `xml:",omitempty"`                                         // Indicates whether the VM Agent is installed on the Virtual Machine. To run a resource extension in a Virtual Machine, this service must be installed.
	VMImageInput                *VMImageInput                 `xml:",omitempty"`                                         // When a VM Image is used to create a new PersistentVMRole, the DiskConfigurations in the VM Image are used to create new Disks for the new VM. This parameter can be used to resize the newly created Disks to a larger size than the underlying DiskConfigurations in the VM Image.

	UseCertAuth bool   `xml:"-"`
	CertPath    string `xml:"-"`
}

// When a VM Image is used to create a new PersistantVMRole, the DiskConfigurations in the VM Image are used to create new Disks for the new VM. This parameter can be used to resize the newly created Disks to a larger size than the underlying DiskConfigurations in the VM Image.
type VMImageInput struct {
	OSDiskConfiguration    *OSDiskConfiguration    `xml:",omitempty"`                                             // This corresponds to the OSDiskConfiguration of the VM Image used to create a new role. The OSDiskConfiguration element is only available using version 2014-10-01 or higher.
	DataDiskConfigurations []DataDiskConfiguration `xml:"DataDiskConfigurations>DataDiskConfiguration,omitempty"` // This corresponds to the DataDiskConfigurations of the VM Image used to create a new role. The DataDiskConfigurations element is only available using version 2014-10-01 or higher.
}

// Used to resize the OS disk of a new VM created from a previously saved VM image
type OSDiskConfiguration struct {
	ResizedSizeInGB float64
}

// Used to resize the data disks of a new VM created from a previously saved VM image
type DataDiskConfiguration struct {
	OSDiskConfiguration
	Name string // The Name of the DataDiskConfiguration being referenced to.

}

// Contains a collection of resource extensions that are to be installed on the Virtual Machine. The VM Agent must be installed on the Virtual Machine to install resource extensions.
// For more information, see Manage Extensions: https://msdn.microsoft.com/en-us/library/dn606311.aspx.
type ResourceExtensionReference struct {
	ReferenceName   string
	Publisher       string
	Name            string
	Version         string
	ParameterValues []ResourceExtensionParameter `xml:"ResourceExtensionParameterValues>ResourceExtensionParameterValue,omitempty"`
	State           string
}

// Specifies the key, value, and type of a parameter that is passed to the resource extension when it is installed.
type ResourceExtensionParameter struct {
	Key   string
	Value string
	Type  ResourceExtensionParameterType // If this value is set to Private, the parameter will not be returned by Get Deployment ().
}

type ResourceExtensionParameterType string

const (
	// Enum values for ResourceExtensionParameterType
	ResourceExtensionParameterTypePublic  ResourceExtensionParameterType = "Public"
	ResourceExtensionParameterTypePrivate ResourceExtensionParameterType = "Private"
)

// Specifies the properties that are used to create a data disk.
type DataVirtualHardDisk struct {
	DiskName    string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription, this element is used to identify the disk to add. If a new disk and the associated VHD are being created by Azure, this element is not used and Azure assigns a unique name that is a combination of the deployment name, role name, and identifying number. The name of the disk must contain only alphanumeric characters, underscores, periods, or dashes. The name must not be longer than 256 characters. The name must not end with period or dash.
	DiskLabel   string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription, this element is ignored. If a new disk is being created, this element is used to provide a description of the disk. The value of this element is only obtained programmatically and does not appear in the Management Portal.
	HostCaching vmdisk.HostCachingType `xml:",omitempty"` // Specifies the caching mode of the data disk. The default value is None.
	MediaLink   string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription or the VHD for the disk already exists in blob storage, this element is ignored. If a VHD file does not exist in blob storage, this element defines the location of the new VHD that is created when the new disk is added.

	Lun                 int     `xml:",omitempty"` // Specifies the Logical Unit Number (LUN) for the data disk. If the disk is the first disk that is added, this element is optional and the default value of 0 is used. If more than one disk is being added, this element is required. Valid LUN values are 0 through 31.
	LogicalDiskSizeInGB float64 `xml:",omitempty"` // Specifies the size, in GB, of an empty disk to be attached to the Virtual Machine. If the disk that is being added is already registered in the subscription, this element is ignored. If the disk and VHD is being created by Azure as it is added, this element defines the size of the new disk.
	SourceMediaLink     string  `xml:",omitempty"` // If the disk that is being added is already registered in the subscription or the VHD for the disk does not exist in blob storage, this element is ignored. If the VHD file exists in blob storage, this element defines the path to the VHD and a disk is registered from it and attached to the virtual machine.
}

// Specifies the properties that are used to create an OS disk.
type OSVirtualHardDisk struct {
	DiskName    string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription, this element is used to identify the disk to add. If a new disk and the associated VHD are being created by Azure, this element is not used and Azure assigns a unique name that is a combination of the deployment name, role name, and identifying number. The name of the disk must contain only alphanumeric characters, underscores, periods, or dashes. The name must not be longer than 256 characters. The name must not end with period or dash.
	DiskLabel   string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription, this element is ignored. If a new disk is being created, this element is used to provide a description of the disk. The value of this element is only obtained programmatically and does not appear in the Management Portal.
	HostCaching vmdisk.HostCachingType `xml:",omitempty"` // Specifies the caching mode of the data disk. The default value is None.
	MediaLink   string                 `xml:",omitempty"` // If the disk that is being added is already registered in the subscription or the VHD for the disk already exists in blob storage, this element is ignored. If a VHD file does not exist in blob storage, this element defines the location of the new VHD that is created when the new disk is added.

	SourceImageName       string  `xml:",omitempty"`
	OS                    string  `xml:",omitempty"`
	RemoteSourceImageLink string  `xml:",omitempty"` // Specifies a publicly accessible URI or a SAS URI to the location where an OS image is stored that is used to create the Virtual Machine. This location can be a different location than the user or platform image repositories in Azure. An image is always associated with a VHD, which is a .vhd file stored as a page blob in a storage account in Azure. If you specify the path to an image with this element, an associated VHD is created and you must use the MediaLink element to specify the location in storage where the VHD will be located. If this element is used, SourceImageName is not used.
	ResizedSizeInGB       float64 `xml:",omitempty"`
}

// Specifies the configuration elements of the Virtual Machine. The type attribute is required to prevent the administrator password from being written to the operation history file.
type ConfigurationSet struct {
	ConfigurationSetType ConfigurationSetType

	// Windows provisioning:
	ComputerName              string                `xml:",omitempty"`                                                   // Optional. Specifies the computer name for the Virtual Machine. If you do not specify a computer name, one is assigned that is a combination of the deployment name, role name, and identifying number. Computer names must be 1 to 15 characters long.
	AdminPassword             string                `xml:",omitempty"`                                                   // Optional. Specifies the password to use for an administrator account on the Virtual Machine that is being created. If you are creating a Virtual Machine using an image, you must specify a name of an administrator account to be created on the machine using the AdminUsername element. You must use the AdminPassword element to specify the password of the administrator account that is being created. If you are creating a Virtual Machine using an existing specialized disk, this element is not used because the account should already exist on the disk.
	EnableAutomaticUpdates    bool                  `xml:",omitempty"`                                                   // Optional. Specifies whether automatic updates are enabled for the Virtual Machine. The default value is true.
	TimeZone                  string                `xml:",omitempty"`                                                   // Optional. Specifies the time zone for the Virtual Machine.
	DomainJoin                *DomainJoin           `xml:",omitempty"`                                                   // Optional. Contains properties that define a domain to which the Virtual Machine will be joined.
	StoredCertificateSettings *[]CertificateSetting `xml:"StoredCertificateSettings>StoredCertificateSetting,omitempty"` // Optional. Contains a list of service certificates with which to provision to the new Virtual Machine.
	WinRMListeners            *WinRMListener        `xml:"WinRM>Listeners>Listener,omitempty"`                           // Optional. Contains configuration settings for the Windows Remote Management service on the Virtual Machine. This enables remote Windows PowerShell.
	AdminUsername             string                `xml:",omitempty"`                                                   // Optional. Specifies the name of the administrator account that is created to access the Virtual Machine. If you are creating a Virtual Machine using an image, you must specify a name of an administrator account to be created by using this element. You must use the AdminPassword element to specify the password of the administrator account that is being created. If you are creating a Virtual Machine using an existing specialized disk, this element is not used because the account should already exist on the disk.
	AdditionalUnattendContent string                `xml:",omitempty"`                                                   // Specifies additional base-64 encoded XML formatted information that can be included in the Unattend.xml file, which is used by Windows Setup.

	// Linux provisioning:
	HostName                         string `xml:",omitempty"` // Required. Specifies the host name for the Virtual Machine. Host names must be 1 to 64 characters long.
	UserName                         string `xml:",omitempty"` // Required. Specifies the name of a user account to be created in the sudoer group of the Virtual Machine. User account names must be 1 to 32 characters long.
	UserPassword                     string `xml:",omitempty"` // Required. Specifies the password for the user account. Passwords must be 6 to 72 characters long.
	DisableSshPasswordAuthentication string `xml:",omitempty"` // Optional. Specifies whether SSH password authentication is disabled. By default this value is set to true.
	SSH                              *SSH   `xml:",omitempty"` // Optional. Specifies the SSH public keys and key pairs to use with the Virtual Machine.

	// In WindowsProvisioningConfiguration: The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes. The file is saved to %SYSTEMDRIVE%\AzureData\CustomData.bin. If the file exists, it is overwritten. The security on directory is set to System:Full Control and Administrators:Full Control.
	// In LinuxProvisioningConfiguration: The base-64 encoded string is located in the ovf-env.xml file on the ISO of the Virtual Machine. The file is copied to /var/lib/waagent/ovf-env.xml by the Azure Linux Agent. The Azure Linux Agent will also place the base-64 encoded data in /var/lib/waagent/CustomData during provisioning. The maximum length of the binary array is 65535 bytes.
	CustomData string `xml:",omitempty"` // Specifies a base-64 encoded string of custom data.

	// Network configuration:
	InputEndpoints                *[]InputEndpoint `xml:"InputEndpoints>InputEndpoint,omitempty"` // Optional in NetworkConfiguration. Contains a collection of external endpoints for the Virtual Machine.
	SubnetNames                   *[]string        `xml:"SubnetNames>SubnetName,omitempty"`       // Required if StaticVirtualNetworkIPAddress is specified; otherwise, optional in NetworkConfiguration. Contains a list of subnets to which the Virtual Machine will belong.
	StaticVirtualNetworkIPAddress string           `xml:",omitempty"`                             // Specifies the internal IP address for the Virtual Machine in a Virtual Network. If you specify this element, you must also specify the SubnetNames element with only one subnet defined. The IP address specified in this element must belong to the subnet that is defined in SubnetNames and it should not be the one of the first four IP addresses or the last IP address in the subnet. Deploying web roles or worker roles into a subnet that has Virtual Machines with StaticVirtualNetworkIPAddress defined is not supported.
	PublicIPs                     *[]PublicIP      `xml:"PublicIPs>PublicIP,omitempty"`           // Contains a public IP address that can be used in addition to the default virtual IP address for the Virtual Machine.
}

type ConfigurationSetType string

const (
	// Enum values for ConfigurationSetType
	ConfigurationSetTypeWindowsProvisioning ConfigurationSetType = "WindowsProvisioningConfiguration"
	ConfigurationSetTypeLinuxProvisioning   ConfigurationSetType = "LinuxProvisioningConfiguration"
	ConfigurationSetTypeNetwork             ConfigurationSetType = "NetworkConfiguration"
)

// Contains properties that define a domain to which the Virtual Machine will be joined.
type DomainJoin struct {
	Credentials     Credentials `xml:",omitempty"` // Specifies the credentials to use to join the Virtual Machine to the domain.
	JoinDomain      string      `xml:",omitempty"` // Specifies the domain to join.
	MachineObjectOU string      `xml:",omitempty"` // Specifies the Lightweight Directory Access Protocol (LDAP) X 500-distinguished name of the organizational unit (OU) in which the computer account is created. This account is in Active Directory on a domain controller in the domain to which the computer is being joined.
}

// Specifies the credentials to use to join the Virtual Machine to the domain.
// If Domain is not specified, Username must specify the user principal name (UPN) format (user@fully-qualified-DNS-domain) or the fully-qualified-DNS-domain\username format.
type Credentials struct {
	Domain   string // Specifies the name of the domain used to authenticate an account. The value is a fully qualified DNS domain.
	Username string // Specifies a user name in the domain that can be used to join the domain.
	Password string // Specifies the password to use to join the domain.
}

// Specifies the parameters for the certificate which to provision to the new Virtual Machine.
type CertificateSetting struct {
	StoreLocation string // Required. Specifies the certificate store location on the Virtual Machine. The only supported value is "LocalMachine".
	StoreName     string // Required. Specifies the name of the certificate store from which the certificate is retrieved. For example, "My".
	Thumbprint    string // Required. Specifies the thumbprint of the certificate. The thumbprint must specify an existing service certificate.
}

// Specifies the protocol and certificate information for a WinRM listener.
type WinRMListener struct {
	Protocol              WinRMProtocol // Specifies the protocol of listener.
	CertificateThumbprint string        `xml:",omitempty"` // Specifies the certificate thumbprint for the secure connection. If this value is not specified, a self-signed certificate is generated and used for the Virtual Machine.
}

type WinRMProtocol string

const (
	// Enum values for WinRMProtocol
	WinRMProtocolHttp  WinRMProtocol = "Http"
	WinRMProtocolHttps WinRMProtocol = "Https"
)

// Specifies the SSH public keys and key pairs to use with the Virtual Machine.
type SSH struct {
	PublicKeys *[]PublicKey `xml:"PublicKeys>PublicKey"`
	KeyPairs   *[]KeyPair   `xml:"KeyPairs>KeyPair"`
}

// Specifies a public SSH key.
type PublicKey struct {
	Fingerprint string // Specifies the SHA1 fingerprint of an X509 certificate associated with the cloud service and includes the SSH public key.
	// Specifies the full path of a file, on the Virtual Machine, where the SSH public key is stored. If
	// the file already exists, the specified key is appended to the file.
	Path string // Usually /home/username/.ssh/authorized_keys
}

// Specifies an SSH keypair.
type KeyPair struct {
	Fingerprint string // Specifies the SHA1 fingerprint of an X509 certificate that is associated with the cloud service and includes the SSH keypair.
	// Specifies the full path of a file, on the virtual machine, which stores the SSH private key. The
	// file is overwritten when multiple keys are written to it. The SSH public key is stored in the same
	// directory and has the same name as the private key file with .pub suffix.
	Path string // Usually /home/username/.ssh/id_rsa
}

// Specifies the properties that define an external endpoint for the Virtual Machine.
type InputEndpoint struct {
	LocalPort int                   // Specifies the internal port on which the Virtual Machine is listening.
	Name      string                // Specifies the name of the external endpoint.
	Port      int                   // Specifies the external port to use for the endpoint.
	Protocol  InputEndpointProtocol //Specifies the transport protocol for the endpoint.
	Vip       string                `xml:",omitempty"`
}

type InputEndpointProtocol string

const (
	// Enum values for InputEndpointProtocol
	InputEndpointProtocolTcp InputEndpointProtocol = "TCP"
	InputEndpointProtocolUdp InputEndpointProtocol = "UDP"
)

// Contains a public IP address that can be used in addition to default virtual IP address for the Virtual Machine.
type PublicIP struct {
	Name                 string // Specifies the name of the public IP address.
	IdleTimeoutInMinutes int    `xml:",omitempty"` // Specifies the timeout for the TCP idle connection. The value can be set between 4 and 30 minutes. The default value is 4 minutes. This element is only used when the protocol is set to TCP.
}

// Contains a certificate for adding it to a hosted service
type ServiceCertificate struct {
	XMLName           xml.Name `xml:"CertificateFile"`
	Data              string
	CertificateFormat string
	Password          string `xml:",omitempty"`
}

// Contains the information for starting a Role.
type StartRoleOperation struct {
	XMLName       xml.Name `xml:"http://schemas.microsoft.com/windowsazure StartRoleOperation"`
	OperationType string
}

// Contains the information for shutting down a Role.
type ShutdownRoleOperation struct {
	XMLName       xml.Name `xml:"http://schemas.microsoft.com/windowsazure ShutdownRoleOperation"`
	OperationType string
}

// Contains the information for restarting a Role.
type RestartRoleOperation struct {
	XMLName       xml.Name `xml:"http://schemas.microsoft.com/windowsazure RestartRoleOperation"`
	OperationType string
}

// Contains the information for capturing a Role
type CaptureRoleOperation struct {
	XMLName                   xml.Name `xml:"http://schemas.microsoft.com/windowsazure CaptureRoleOperation"`
	OperationType             string
	PostCaptureAction         PostCaptureAction
	ProvisioningConfiguration *ConfigurationSet `xml:",omitempty"`
	TargetImageLabel          string
	TargetImageName           string
}

type PostCaptureAction string

const (
	// Enum values for PostCaptureAction
	PostCaptureActionDelete      PostCaptureAction = "Delete"
	PostCaptureActionReprovision PostCaptureAction = "Reprovision"
)

// Contains a list of the available role sizes
type RoleSizeList struct {
	XMLName   xml.Name   `xml:"RoleSizes"`
	RoleSizes []RoleSize `xml:"RoleSize"`
}

// Contains a detailed explanation of a role size
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

// Contains the definition of a DNS server for virtual machine deployment
type DnsServer struct {
	Name    string
	Address string
}

// Contains the definition of a load balancer for virtual machine deployment
type LoadBalancer struct {
	Name                          string        // Specifies the name of the internal load balancer.
	Type                          IPAddressType `xml:"FrontendIpConfiguration>Type"`                                    // Specifies the type of virtual IP address that is provided by the load balancer. The only allowable value is Private.
	SubnetName                    string        `xml:"FrontendIpConfiguration>SubnetName,omitempty"`                    // Required if the deployment exists in a virtual network and a StaticVirtualNetworkIPAddress is assigned. Specifies the subnet of the virtual network that the load balancer uses. The virtual IP address that is managed by the load balancer is contained in this subnet.
	StaticVirtualNetworkIPAddress string        `xml:"FrontendIpConfiguration>StaticVirtualNetworkIPAddress,omitempty"` // Specifies a specific virtual IP address that the load balancer uses from the subnet in the virtual network.
}

type IPAddressType string

const (
	IPAddressTypePrivate IPAddressType = "Private" // Only allowed value (currently) for IPAddressType
)

type ResourceExtensions struct {
	List []ResourceExtension `xml:"ResourceExtension"`
}

type ResourceExtension struct {
	Publisher                   string
	Name                        string
	Version                     string
	Label                       string
	Description                 string
	PublicConfigurationSchema   string
	PrivateConfigurationSchema  string
	SampleConfig                string
	ReplicationCompleted        string
	Eula                        string
	PrivacyUri                  string
	HomepageUri                 string
	IsJsonExtension             bool
	IsInternalExtension         bool
	DisallowMajorVersionUpgrade bool
	CompanyName                 string
	SupportedOS                 string
	PublishedDate               string
}

type PersistentVMRole struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windowsazure PersistentVMRole"`
	Role
}
