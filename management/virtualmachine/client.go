// Package virtualmachine implements operations on Azure virtual machines using the Service Management REST API
package virtualmachine

import (
	"encoding/xml"
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
)

const (
	azureDeploymentListURL   = "services/hostedservices/%s/deployments"
	azureDeploymentURL       = "services/hostedservices/%s/deployments/%s"
	deleteAzureDeploymentURL = "services/hostedservices/%s/deployments/%s?comp=media"
	azureRoleURL             = "services/hostedservices/%s/deployments/%s/roles/%s"
	azureOperationsURL       = "services/hostedservices/%s/deployments/%s/roleinstances/%s/Operations"
	azureRoleSizeListURL     = "rolesizes"

	errParamNotSpecified = "Parameter %s is not specified."
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
