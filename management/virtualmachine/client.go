// Package virtualmachine implements operations on Azure virtual machines using the Service Management REST API
package virtualmachine

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureDeploymentListURL   = "services/hostedservices/%s/deployments"
	azureDeploymentURL       = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL = "services/hostedservices/%s/deployments/%s?comp=media"
	azureRoleURL             = "services/hostedservices/%s/deployments/%s/roles/%s"
	azureOperationsURL       = "services/hostedservices/%s/deployments/%s/roleinstances/%s/Operations"
	azureRoleSizeListURL     = "rolesizes"

	dockerPublicConfigVersion = 2

	errParamNotSpecified            = "Parameter %s is not specified."
	errProvisioningConfDoesNotExist = "You should set azure VM provisioning config first"
	errInvalidRoleSize              = "Invalid role size: %s. Available role sizes: %s"
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

func (self VirtualMachineClient) createEndpoint(name string, protocol InputEndpointProtocol, extertalPort int, internalPort int) InputEndpoint {
	return InputEndpoint{
		Name:      name,
		Protocol:  protocol,
		Port:      extertalPort,
		LocalPort: internalPort,
	}
}
