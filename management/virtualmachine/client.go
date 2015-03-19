// Package virtualmachine implements operations on Azure virtual machines using the Service Management REST API
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
	locationclient "github.com/MSOpenTech/azure-sdk-for-go/management/location"
	storageserviceclient "github.com/MSOpenTech/azure-sdk-for-go/management/storageservice"
)

const (
	azureDeploymentListURL   = "services/hostedservices/%s/deployments"
	azureDeploymentURL       = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL = "services/hostedservices/%s/deployments/%s?comp=media"
	azureRoleURL             = "services/hostedservices/%s/deployments/%s/roles/%s"
	azureOperationsURL       = "services/hostedservices/%s/deployments/%s/roleinstances/%s/Operations"
	azureRoleSizeListURL     = "rolesizes"

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
)

//NewClient is used to instantiate a new VmClient from an Azure client
func NewClient(client management.Client) VirtualMachineClient {
	return VirtualMachineClient{client: client}
}

func (self VirtualMachineClient) CreateDeployment(role *Role, cloudserviceName string) (requestId string, err error) {
	if role == nil {
		return "", fmt.Errorf(errParamNotSpecified, "role")
	}

	vMDeploymentBytes, err := xml.Marshal(DeploymentRequest{
		Name:           role.RoleName,
		DeploymentSlot: "Production",
		Label:          role.RoleName,
		RoleList:       []*Role{role}})

	requestURL := fmt.Sprintf(azureDeploymentListURL, cloudserviceName)
	return self.client.SendAzurePostRequest(requestURL, vMDeploymentBytes)
}

func (self VirtualMachineClient) CreateAzureVMConfiguration(dnsName, instanceSize, imageName, location string) (*Role, error) {
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

	role, err := self.createAzureVMRole(dnsName, instanceSize, imageName, location)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (self VirtualMachineClient) AddAzureLinuxProvisioningConfig(azureVMConfiguration *Role, userName, password, certPath string, sshPort int) (*Role, error) {
	if azureVMConfiguration == nil {
		return nil, fmt.Errorf(errParamNotSpecified, "azureVMConfiguration")
	}
	if userName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "userName")
	}

	provisioningConfig, err := self.createLinuxProvisioningConfig(azureVMConfiguration.RoleName, userName, password, certPath)
	if err != nil {
		return nil, err
	}

	networkConfig, networkErr := self.createNetworkConfig(osLinux, sshPort)
	if networkErr != nil {
		return nil, err
	}

	azureVMConfiguration.ConfigurationSets = &[]ConfigurationSet{provisioningConfig, networkConfig}

	if len(certPath) > 0 {
		azureVMConfiguration.UseCertAuth = true
		azureVMConfiguration.CertPath = certPath
	}

	return azureVMConfiguration, nil
}

func (self VirtualMachineClient) SetAzureVMExtension(azureVMConfiguration *Role, name string, publisher string, version string, referenceName string, state string, publicConfigurationValue string, privateConfigurationValue string) (*Role, error) {
	if azureVMConfiguration == nil {
		return nil, fmt.Errorf(errParamNotSpecified, "azureVMConfiguration")
	}
	if name == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "name")
	}
	if publisher == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "publisher")
	}
	if version == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "version")
	}
	if referenceName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "referenceName")
	}

	extension := ResourceExtensionReference{
		Name:          name,
		Publisher:     publisher,
		Version:       version,
		ReferenceName: referenceName,
		State:         state,
	}

	if privateConfigurationValue != "" {
		extension.ParameterValues = append(extension.ParameterValues, ResourceExtensionParameter{
			Key:   "ignored",
			Value: base64.StdEncoding.EncodeToString([]byte(privateConfigurationValue)),
			Type:  "Private",
		})
	}

	if publicConfigurationValue != "" {
		extension.ParameterValues = append(extension.ParameterValues, ResourceExtensionParameter{
			Key:   "ignored",
			Value: base64.StdEncoding.EncodeToString([]byte(publicConfigurationValue)),
			Type:  "Public",
		})
	}

	if azureVMConfiguration.ResourceExtensionReferences == nil {
		azureVMConfiguration.ResourceExtensionReferences = &[]ResourceExtensionReference{extension}
	} else {
		*azureVMConfiguration.ResourceExtensionReferences = append(*azureVMConfiguration.ResourceExtensionReferences, extension)
	}

	return azureVMConfiguration, nil
}

func (self VirtualMachineClient) SetAzureDockerVMExtension(azureVMConfiguration *Role, dockerPort int, version string) (*Role, error) {
	if azureVMConfiguration == nil {
		return nil, fmt.Errorf(errParamNotSpecified, "azureVMConfiguration")
	}

	if version == "" {
		version = "0.3"
	}

	err := self.addDockerPort(*azureVMConfiguration.ConfigurationSets, dockerPort)
	if err != nil {
		return nil, err
	}

	publicConfiguration, err := self.createDockerPublicConfig(dockerPort)
	if err != nil {
		return nil, err
	}

	privateConfiguration := "{}"

	azureVMConfiguration, err = self.SetAzureVMExtension(azureVMConfiguration, "DockerExtension", "MSOpenTech.Extensions", version, "DockerExtension", "enable", publicConfiguration, privateConfiguration)
	return azureVMConfiguration, nil
}

func (self VirtualMachineClient) GetDeployment(cloudserviceName, deploymentName string) (*DeploymentResponse, error) {
	if cloudserviceName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return nil, fmt.Errorf(errParamNotSpecified, "deploymentName")
	}

	deployment := new(DeploymentResponse)

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

func (self VirtualMachineClient) DeleteDeployment(cloudserviceName, deploymentName string) (requestId string, err error) {
	if cloudserviceName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "cloudserviceName")
	}
	if deploymentName == "" {
		return "", fmt.Errorf(errParamNotSpecified, "deploymentName")
	}

	requestURL := fmt.Sprintf(deleteAzureDeploymentURL, cloudserviceName, deploymentName)
	return self.client.SendAzureDeleteRequest(requestURL)
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

	startRoleOperationBytes, err := xml.Marshal(StartRoleOperation{
		OperationType: "StartRoleOperation",
	})
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

	shutdownRoleOperationBytes, err := xml.Marshal(ShutdownRoleOperation{
		OperationType: "ShutdownRoleOperation",
	})
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

	restartRoleOperationBytes, err := xml.Marshal(RestartRoleOperation{
		OperationType: "RestartRoleOperation",
	})
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
		*configurationSets[i].InputEndpoints = append(*configurationSets[i].InputEndpoints, dockerEndpoint)
	}

	return nil
}
func (self VirtualMachineClient) createAzureVMRole(name, instanceSize, imageName, location string) (*Role, error) {
	vhd, err := self.createOSVirtualHardDisk(name, imageName, location)
	if err != nil {
		return nil, err
	}

	return &Role{
		RoleName:            name,
		RoleSize:            instanceSize,
		RoleType:            "PersistentVMRole",
		ProvisionGuestAgent: true,
		OSVirtualHardDisk:   vhd,
	}, nil
}

func (self VirtualMachineClient) createOSVirtualHardDisk(dnsName, imageName, location string) (*OSVirtualHardDisk, error) {
	link, err := self.getVHDMediaLink(dnsName, location)
	if err != nil {
		return &OSVirtualHardDisk{}, err
	}

	return &OSVirtualHardDisk{
		SourceImageName: imageName,
		MediaLink:       link,
	}, nil
}

func (self VirtualMachineClient) getVHDMediaLink(dnsName, location string) (string, error) {
	storageServiceClient := storageserviceclient.NewClient(self.client)

	storageService, err := storageServiceClient.GetStorageServiceByLocation(location)
	if err != nil {
		return "", err
	}

	if storageService == nil {
		uuid, err := newUUID()
		if err != nil {
			return "", err
		}

		serviceName := "portalvhds" + uuid
		storageService, err = storageServiceClient.Create(storageserviceclient.StorageAccountCreateParameters{
			ServiceName: serviceName,
			Location:    location,
			Label:       base64.StdEncoding.EncodeToString([]byte(serviceName)),
		})
		if err != nil {
			return "", err
		}
	}

	blobEndpoint, err := storageServiceClient.GetBlobEndpoint(storageService)
	if err != nil {
		return "", err
	}

	vhdMediaLink := blobEndpoint + "vhds/" + dnsName + "-" + time.Now().Local().Format("20060102150405") + ".vhd"
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

	if !disableSshPasswordAuthentication {
		provisioningConfig.DisableSshPasswordAuthentication = "false"
	}
	provisioningConfig.ConfigurationSetType = "LinuxProvisioningConfiguration"
	provisioningConfig.HostName = dnsName
	provisioningConfig.UserName = userName
	provisioningConfig.UserPassword = userPassword

	if len(certPath) > 0 {
		var err error
		ssh, err := self.createSshConfig(certPath, userName)
		if err != nil {
			return provisioningConfig, err
		}
		provisioningConfig.SSH = &ssh
	}

	return provisioningConfig, nil
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

	sshConfig.PublicKeys = &[]PublicKey{publicKey}
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
	networkConfig := ConfigurationSet{
		ConfigurationSetType: "NetworkConfiguration",
	}

	var endpoint InputEndpoint
	if os == osLinux {
		endpoint = self.createEndpoint("ssh", "tcp", sshPort, 22)
	} else if os == osWindows {
		//!TODO add rdp endpoint
	} else {
		return networkConfig, errors.New(fmt.Sprintf(errInvalidOS))
	}

	*networkConfig.InputEndpoints = append(*networkConfig.InputEndpoints, endpoint)

	return networkConfig, nil
}

func (self VirtualMachineClient) createEndpoint(name string, protocol InputEndpointProtocol, extertalPort int, internalPort int) InputEndpoint {
	return InputEndpoint{
		Name:      name,
		Protocol:  protocol,
		Port:      extertalPort,
		LocalPort: internalPort,
	}
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
