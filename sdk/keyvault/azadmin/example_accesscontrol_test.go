// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
	"github.com/google/uuid"
)

var accessControlClient azadmin.AccessControlClient

func ExampleNewAccessControlClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azadmin.NewAccessControlClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleAccessControlClient_CreateOrUpdateRoleDefinition() {
	scope := azadmin.RoleScopeGlobal
	name := uuid.New().String()
	roleType := azadmin.RoleTypeCustomRole
	roleName := "ExampleRoleName"
	parameters := azadmin.RoleDefinitionCreateParameters{
		Properties: &azadmin.RoleDefinitionProperties{
			AssignableScopes: []*azadmin.RoleScope{to.Ptr(scope)},
			Description:      to.Ptr("Example description"),
			Permissions:      []*azadmin.Permission{{DataActions: []*azadmin.DataAction{to.Ptr(azadmin.DataActionBackupHsmKeys), to.Ptr(azadmin.DataActionCreateHsmKey)}}},
			RoleName:         to.Ptr(roleName),
			RoleType:         to.Ptr(roleType),
		},
	}

	roleDefinition, err := accessControlClient.CreateOrUpdateRoleDefinition(context.TODO(), scope, name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Definition Name: %s", *roleDefinition.Name)
}

func ExampleAccessControlClient_CreateRoleAssignment() {
	scope := azadmin.RoleScopeGlobal
	name := uuid.New().String()
	parameters := azadmin.RoleAssignmentCreateParameters{
		Properties: &azadmin.RoleAssignmentProperties{
			PrincipalID:      to.Ptr("d26e28bc-991f-11ed-a8fc-0242ac120002"),                                                                      // example principal ID
			RoleDefinitionID: to.Ptr("Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/c368d8da-991f-11ed-a8fc-0242ac120002"), // example role definition ID
		},
	}

	roleAssignment, err := accessControlClient.CreateRoleAssignment(context.TODO(), scope, name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Assignment Name: %s", *roleAssignment.Name)
}
