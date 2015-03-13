package virtualmachine

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
	"unicode"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
	hostedserviceclient "github.com/MSOpenTech/azure-sdk-for-go/management/hostedservice"
	locationclient "github.com/MSOpenTech/azure-sdk-for-go/management/location"
	storageserviceclient "github.com/MSOpenTech/azure-sdk-for-go/management/storageservice"
)

const (
	azureXmlns                        = "http://schemas.microsoft.com/windowsazure"
	azureDeploymentListURL            = "services/hostedservices/%s/deployments"
	azureHostedServiceListURL         = "services/hostedservices"
	deleteAzureHostedServiceURL       = "services/hostedservices/%s?comp=media"
	azureHostedServiceAvailabilityURL = "services/hostedservices/operations/isavailable/%s"
	azureDeploymentURL                = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL          = "services/hostedservices/%s/deployments/%s?comp=media"
	azureRoleURL                      = "services/hostedservices/%s/deployments/%s/roles/%s"
	azureOperationsURL                = "services/hostedservices/%s/deployments/%s/roleinstances/%s/Operations"
	azureCertificatListURL            = "services/hostedservices/%s/certificates"
	azureRoleSizeListURL              = "rolesizes"

	osLinux                   = "Linux"
	osWindows                 = "Windows"
	dockerPublicConfigVersion = 2

	errParamNotSpecified            = "Parameter %s is not specified."
	errProvisioningConfDoesNotExist = "You should set azure VM provisioning config first"
	errInvalidCertExtension         = "Certificate %s is invalid. Please specify %s certificate."
	errInvalidOS                    = "You must specify correct OS param. Valid values are 'Linux' and 'Windows'"
	errInvalidPasswordLength        = "Password must be between 4 and 30 characters."
	errInvalidPassword              = "Password must have at least one upper case, lower case and numeric character."
	errInvalidRoleSize              = "Invalid role size: %s. Available role sizes: %s"
	errInvalidRoleSizeInLocation    = "Role size: %s not available in location: %s."
	errInvalidDnsLength             = "The DNS name must be between 3 and 25 characters."
)

//NewClient is used to instantiate a new VmClient from an Azure client
func NewClient(client management.Client) VirtualMachineClient {
	return VirtualMachineClient{client: client}
}

func (self VirtualMachineClient) CreateAzureVM(config *VMCreateConfiguration) error {
	if config == nil {
		return fmt.Errorf(errParamNotSpecified, "azureVMConfiguration")
	}

	hardDisk, err := self.createOSVirtualHardDisk(config)
	if err != nil {
		return err
	}
	config.role.OSVirtualHardDisk = hardDisk

	hostedServiceClient := hostedserviceclient.NewClient(self.client)
	requestId, err := hostedServiceClient.CreateHostedService(config.DnsName, config.Location, "", config.DnsName, "")
	if err != nil {
		return err
	}

	err = self.client.WaitAsyncOperation(requestId)
	if err != nil {
		return err
	}

	if config.role.UseCertAuth {
		err = self.uploadServiceCert(config.DnsName, config.role.CertPath)
		if err != nil {
			hostedServiceClient.DeleteHostedService(config.DnsName)
			return err
		}
	}

	vMDeployment := self.createVMDeploymentConfig(config.role)
	vMDeploymentBytes, err := xml.Marshal(vMDeployment)
	if err != nil {
		hostedServiceClient.DeleteHostedService(config.DnsName)
		return err
	}

	requestURL := fmt.Sprintf(azureDeploymentListURL, config.role.RoleName)
	requestId, err = self.client.SendAzurePostRequest(requestURL, vMDeploymentBytes)
	if err != nil {
		hostedServiceClient.DeleteHostedService(config.DnsName)
		return err
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) NewVMCreateConfiguration(dnsName, instanceSize, imageName, location string) (*VMCreateConfiguration, error) {

	if dnsName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "dnsName")
	}
	if instanceSize == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "instanceSize")
	}
	if imageName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "imageName")
	}
	if location == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "location")
	}

	locationClient := locationclient.NewClient(self.client)
	locationInfo, err := locationClient.GetLocation(location)
	if err != nil {
		return nil, err
	}

	sizeAvailable, err := self.isInstanceSizeAvailableInLocation(locationInfo, instanceSize)
	if err != nil {
		return nil, err
	}

	if sizeAvailable == false {
		return nil, fmt.Errorf(errInvalidRoleSizeInLocation, instanceSize, location)
	}

	config := &VMCreateConfiguration{
		ImageName: imageName,
		Location:  location,
		client:    self,
		role: &Role{
			RoleName:            dnsName,
			RoleSize:            instanceSize,
			RoleType:            "PersistentVMRole",
			ProvisionGuestAgent: true,
		},
	}

	return config, nil
}

func (self *VMCreateConfiguration) AddLinuxConfig(userName, password, certPath string, sshPort int) error {
	if userName == "" {
		return fmt.Errorf(errParamNotSpecified, "userName")
	}

	configurationSets := ConfigurationSets{}
	provisioningConfig, err := self.client.createLinuxProvisioningConfig(self.DnsName, userName, password, certPath)
	if err != nil {
		return nil
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, provisioningConfig)

	networkConfig, networkErr := self.client.createNetworkConfig(osLinux, sshPort)
	if networkErr != nil {
		return nil
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, networkConfig)

	self.role.ConfigurationSets = configurationSets

	if len(certPath) > 0 {
		self.role.UseCertAuth = true
		self.role.CertPath = certPath
	}

	return nil
}

func (self *VMCreateConfiguration) SetExtension(name string, publisher string, version string, referenceName string, state string, publicConfigurationValue string, privateConfigurationValue string) error {
	if name == "" {
		return fmt.Errorf(errParamNotSpecified, "name")
	}
	if publisher == "" {
		return fmt.Errorf(errParamNotSpecified, "publisher")
	}
	if version == "" {
		return fmt.Errorf(errParamNotSpecified, "version")
	}
	if referenceName == "" {
		return fmt.Errorf(errParamNotSpecified, "referenceName")
	}

	extension := ResourceExtensionReference{}
	extension.Name = name
	extension.Publisher = publisher
	extension.Version = version
	extension.ReferenceName = referenceName
	extension.State = state

	if len(privateConfigurationValue) > 0 {
		privateConfig := ResourceExtensionParameter{}
		privateConfig.Key = "ignored"
		privateConfig.Value = base64.StdEncoding.EncodeToString([]byte(privateConfigurationValue))
		privateConfig.Type = "Private"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, privateConfig)
	}

	if len(publicConfigurationValue) > 0 {
		publicConfig := ResourceExtensionParameter{}
		publicConfig.Key = "ignored"
		publicConfig.Value = base64.StdEncoding.EncodeToString([]byte(publicConfigurationValue))
		publicConfig.Type = "Public"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, publicConfig)
	}

	self.role.ResourceExtensionReferences.ResourceExtensionReference = append(self.role.ResourceExtensionReferences.ResourceExtensionReference, extension)

	return nil
}

func (self *VMCreateConfiguration) SetDockerExtension(dockerPort int, version string) error {
	if version == "" {
		version = "0.3"
	}

	err := self.client.addDockerPort(self.role.ConfigurationSets.ConfigurationSet, dockerPort)
	if err != nil {
		return err
	}

	publicConfiguration, err := self.client.createDockerPublicConfig(dockerPort)
	if err != nil {
		return err
	}

	privateConfiguration := "{}"

	err = self.SetExtension("DockerExtension", "MSOpenTech.Extensions", version, "DockerExtension", "enable", publicConfiguration, privateConfiguration)
	return nil
}

func (self VirtualMachineClient) GetVMDeployment(cloudserviceName, deploymentName string) (*VMDeployment, error) {
	if cloudserviceName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "deploymentName")
	}

	deployment := new(VMDeployment)

	requestURL := fmt.Sprintf(azureDeploymentURL, cloudserviceName, deploymentName)
	response, azureErr := self.client.SendAzureGetRequest(requestURL)
	if azureErr != nil {
		return nil, azureErr
	}

	err := xml.Unmarshal(response, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (self VirtualMachineClient) DeleteVMDeployment(cloudserviceName, deploymentName string) error {
	if cloudserviceName == "" {
		return fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return fmt.Errorf(errParamNotSpecified, "deploymentName")
	}

	requestURL := fmt.Sprintf(deleteAzureDeploymentURL, cloudserviceName, deploymentName)
	requestId, err := self.client.SendAzureDeleteRequest(requestURL)
	if err != nil {
		return err
	}

	err = self.client.WaitAsyncOperation(requestId)
	if err != nil {
		return err
	}

	return nil
}

func (self VirtualMachineClient) GetRole(cloudserviceName, deploymentName, roleName string) (*Role, error) {
	if cloudserviceName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "roleName")
	}

	role := new(Role)

	requestURL := fmt.Sprintf(azureRoleURL, cloudserviceName, deploymentName, roleName)
	response, azureErr := self.client.SendAzureGetRequest(requestURL)
	if azureErr != nil {
		return nil, azureErr
	}

	err := xml.Unmarshal(response, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (self VirtualMachineClient) StartRole(cloudserviceName, deploymentName, roleName string) error {
	if cloudserviceName == "" {
		return fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return fmt.Errorf(errParamNotSpecified, "roleName")
	}

	startRoleOperation := self.createStartRoleOperation()

	startRoleOperationBytes, err := xml.Marshal(startRoleOperation)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(azureOperationsURL, cloudserviceName, deploymentName, roleName)
	requestId, azureErr := self.client.SendAzurePostRequest(requestURL, startRoleOperationBytes)
	if azureErr != nil {
		return azureErr
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) ShutdownRole(cloudserviceName, deploymentName, roleName string) error {
	if cloudserviceName == "" {
		return fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return fmt.Errorf(errParamNotSpecified, "roleName")
	}

	shutdownRoleOperation := self.createShutdowRoleOperation()

	shutdownRoleOperationBytes, err := xml.Marshal(shutdownRoleOperation)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(azureOperationsURL, cloudserviceName, deploymentName, roleName)
	requestId, azureErr := self.client.SendAzurePostRequest(requestURL, shutdownRoleOperationBytes)
	if azureErr != nil {
		return azureErr
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) RestartRole(cloudserviceName, deploymentName, roleName string) error {
	if cloudserviceName == "" {
		return fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return fmt.Errorf(errParamNotSpecified, "roleName")
	}

	restartRoleOperation := self.createRestartRoleOperation()

	restartRoleOperationBytes, err := xml.Marshal(restartRoleOperation)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(azureOperationsURL, cloudserviceName, deploymentName, roleName)
	requestId, azureErr := self.client.SendAzurePostRequest(requestURL, restartRoleOperationBytes)
	if azureErr != nil {
		return azureErr
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) DeleteRole(cloudserviceName, deploymentName, roleName string) error {
	if cloudserviceName == "" {
		return fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return fmt.Errorf(errParamNotSpecified, "deploymentName")
	}
	if roleName == "" {
		return fmt.Errorf(errParamNotSpecified, "roleName")
	}

	requestURL := fmt.Sprintf(azureRoleURL, cloudserviceName, deploymentName, roleName)
	requestId, azureErr := self.client.SendAzureDeleteRequest(requestURL)
	if azureErr != nil {
		return azureErr
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) GetRoleSizeList() (RoleSizeList, error) {
	roleSizeList := RoleSizeList{}

	response, err := self.client.SendAzureGetRequest(azureRoleSizeListURL)
	if err != nil {
		return roleSizeList, err
	}

	err = xml.Unmarshal(response, &roleSizeList)
	if err != nil {
		return roleSizeList, err
	}

	return roleSizeList, err
}

func (self VirtualMachineClient) ResolveRoleSize(roleSizeName string) error {
	if roleSizeName == "" {
		return fmt.Errorf(errParamNotSpecified, "roleSizeName")
	}

	roleSizeList, err := self.GetRoleSizeList()
	if err != nil {
		return err
	}

	for _, roleSize := range roleSizeList.RoleSizes {
		if roleSize.Name != roleSizeName {
			continue
		}

		return nil
	}

	var availableSizes bytes.Buffer
	for _, existingSize := range roleSizeList.RoleSizes {
		availableSizes.WriteString(existingSize.Name + ", ")
	}

	return errors.New(fmt.Sprintf(errInvalidRoleSize, roleSizeName, strings.Trim(availableSizes.String(), ", ")))
}

func (self VirtualMachineClient) createStartRoleOperation() StartRoleOperation {
	startRoleOperation := StartRoleOperation{}
	startRoleOperation.OperationType = "StartRoleOperation"
	startRoleOperation.Xmlns = azureXmlns

	return startRoleOperation
}

func (self VirtualMachineClient) createShutdowRoleOperation() ShutdownRoleOperation {
	shutdownRoleOperation := ShutdownRoleOperation{}
	shutdownRoleOperation.OperationType = "ShutdownRoleOperation"
	shutdownRoleOperation.Xmlns = azureXmlns

	return shutdownRoleOperation
}

func (self VirtualMachineClient) createRestartRoleOperation() RestartRoleOperation {
	startRoleOperation := RestartRoleOperation{}
	startRoleOperation.OperationType = "RestartRoleOperation"
	startRoleOperation.Xmlns = azureXmlns

	return startRoleOperation
}

func (self VirtualMachineClient) createDockerPublicConfig(dockerPort int) (string, error) {
	config := dockerPublicConfig{DockerPort: dockerPort, Version: dockerPublicConfigVersion}
	configJson, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(configJson), nil
}

func (self VirtualMachineClient) addDockerPort(configurationSets []ConfigurationSet, dockerPort int) error {
	if len(configurationSets) == 0 {
		return errors.New(errProvisioningConfDoesNotExist)
	}

	for i := 0; i < len(configurationSets); i++ {
		if configurationSets[i].ConfigurationSetType != "NetworkConfiguration" {
			continue
		}

		dockerEndpoint := self.createEndpoint("docker", "tcp", dockerPort, dockerPort)
		configurationSets[i].InputEndpoints.InputEndpoint = append(configurationSets[i].InputEndpoints.InputEndpoint, dockerEndpoint)
	}

	return nil
}

func (self VirtualMachineClient) createVMDeploymentConfig(role *Role) VMDeployment {
	deployment := VMDeployment{}
	deployment.Name = role.RoleName
	deployment.Xmlns = azureXmlns
	deployment.DeploymentSlot = "Production"
	deployment.Label = role.RoleName
	deployment.RoleList.Role = append(deployment.RoleList.Role, role)

	return deployment
}

func (self VirtualMachineClient) createOSVirtualHardDisk(config *VMCreateConfiguration) (OSVirtualHardDisk, error) {
	oSVirtualHardDisk := OSVirtualHardDisk{}

	oSVirtualHardDisk.SourceImageName = config.ImageName

	var err error
	oSVirtualHardDisk.MediaLink, err = self.getVHDMediaLink(config)
	if err != nil {
		return oSVirtualHardDisk, err
	}

	return oSVirtualHardDisk, nil
}

func (self VirtualMachineClient) getVHDMediaLink(config *VMCreateConfiguration) (string, error) {
	storageServiceClient := storageserviceclient.NewClient(self.client)

	storageService, err := storageServiceClient.GetStorageServiceByLocation(config.Location)
	if err != nil {
		return "", err
	}

	if storageService == nil {
		uuid, err := newUUID()
		if err != nil {
			return "", err
		}

		serviceName := "portalvhds" + uuid
		storageService, err = storageServiceClient.CreateStorageService(serviceName, config.Location, config.StorageAccountType)
		if err != nil {
			return "", err
		}
	}

	blobEndpoint, err := storageServiceClient.GetBlobEndpoint(storageService)
	if err != nil {
		return "", err
	}

	vhdMediaLink := blobEndpoint + "vhds/" + config.DnsName + "-" + time.Now().Local().Format("20060102150405") + ".vhd"
	return vhdMediaLink, nil
}

func (self VirtualMachineClient) createLinuxProvisioningConfig(dnsName, userName, userPassword, certPath string) (ConfigurationSet, error) {
	provisioningConfig := ConfigurationSet{}

	disableSshPasswordAuthentication := false
	if userPassword == "" {
		disableSshPasswordAuthentication = true
		// We need to set dummy password otherwise azure API will throw an error
		userPassword = "P@ssword1"
	} else {
		err := self.verifyPassword(userPassword)
		if err != nil {
			return provisioningConfig, err
		}
	}

	provisioningConfig.DisableSshPasswordAuthentication = disableSshPasswordAuthentication
	provisioningConfig.ConfigurationSetType = "LinuxProvisioningConfiguration"
	provisioningConfig.HostName = dnsName
	provisioningConfig.UserName = userName
	provisioningConfig.UserPassword = userPassword

	if len(certPath) > 0 {
		var err error
		provisioningConfig.SSH, err = self.createSshConfig(certPath, userName)
		if err != nil {
			return provisioningConfig, err
		}
	}

	return provisioningConfig, nil
}

func (self VirtualMachineClient) uploadServiceCert(dnsName, certPath string) error {
	certificateConfig, err := self.createServiceCertDeploymentConf(certPath)
	if err != nil {
		return err
	}

	certificateConfigBytes, err := xml.Marshal(certificateConfig)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf(azureCertificatListURL, dnsName)
	requestId, azureErr := self.client.SendAzurePostRequest(requestURL, certificateConfigBytes)
	if azureErr != nil {
		return azureErr
	}

	return self.client.WaitAsyncOperation(requestId)
}

func (self VirtualMachineClient) createServiceCertDeploymentConf(certPath string) (ServiceCertificate, error) {
	certConfig := ServiceCertificate{}
	certConfig.Xmlns = azureXmlns
	data, err := ioutil.ReadFile(certPath)
	if err != nil {
		return certConfig, err
	}

	certData := base64.StdEncoding.EncodeToString(data)
	certConfig.Data = certData
	certConfig.CertificateFormat = "pfx"

	return certConfig, nil
}

func (self VirtualMachineClient) createSshConfig(certPath, userName string) (SSH, error) {
	sshConfig := SSH{}
	publicKey := PublicKey{}

	err := self.checkServiceCertExtension(certPath)
	if err != nil {
		return sshConfig, err
	}

	fingerprint, err := self.getServiceCertFingerprint(certPath)
	if err != nil {
		return sshConfig, err
	}

	publicKey.Fingerprint = fingerprint
	publicKey.Path = "/home/" + userName + "/.ssh/authorized_keys"

	sshConfig.PublicKeys.PublicKey = append(sshConfig.PublicKeys.PublicKey, publicKey)
	return sshConfig, nil
}

func (self VirtualMachineClient) getServiceCertFingerprint(certPath string) (string, error) {
	certData, readErr := ioutil.ReadFile(certPath)
	if readErr != nil {
		return "", readErr
	}

	block, rest := pem.Decode(certData)
	if block == nil {
		return "", errors.New(string(rest))
	}

	sha1sum := sha1.Sum(block.Bytes)
	fingerprint := fmt.Sprintf("%X", sha1sum)
	return fingerprint, nil
}

func (self VirtualMachineClient) checkServiceCertExtension(certPath string) error {
	certParts := strings.Split(certPath, ".")
	certExt := certParts[len(certParts)-1]

	acceptedExtension := "pem"
	if certExt != acceptedExtension {
		return errors.New(fmt.Sprintf(errInvalidCertExtension, certPath, acceptedExtension))
	}

	return nil
}

func (self VirtualMachineClient) createNetworkConfig(os string, sshPort int) (ConfigurationSet, error) {
	networkConfig := ConfigurationSet{}
	networkConfig.ConfigurationSetType = "NetworkConfiguration"

	var endpoint InputEndpoint
	if os == osLinux {
		endpoint = self.createEndpoint("ssh", "tcp", sshPort, 22)
	} else if os == osWindows {
		//!TODO add rdp endpoint
	} else {
		return networkConfig, errors.New(fmt.Sprintf(errInvalidOS))
	}

	networkConfig.InputEndpoints.InputEndpoint = append(networkConfig.InputEndpoints.InputEndpoint, endpoint)

	return networkConfig, nil
}

func (self VirtualMachineClient) createEndpoint(name string, protocol string, extertalPort int, internalPort int) InputEndpoint {
	endpoint := InputEndpoint{}
	endpoint.Name = name
	endpoint.Protocol = protocol
	endpoint.Port = extertalPort
	endpoint.LocalPort = internalPort

	return endpoint
}

func (self VirtualMachineClient) verifyPassword(password string) error {
	if len(password) < 4 || len(password) > 30 {
		return fmt.Errorf(errInvalidPasswordLength)
	}

next:
	for _, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
	} {
		for _, r := range password {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf(errInvalidPassword)
	}
	return nil
}

func (self VirtualMachineClient) isInstanceSizeAvailableInLocation(location *locationclient.Location, instanceSize string) (bool, error) {
	if instanceSize == "" {
		return false, fmt.Errorf(errParamNotSpecified, "vmSize")
	}

	for _, availableRoleSize := range location.VirtualMachineRoleSizes {
		if availableRoleSize == instanceSize {
			return true, nil
		}
	}

	return false, nil
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x", uuid[10:]), nil
}
