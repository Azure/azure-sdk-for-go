// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rbac_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
	"github.com/google/uuid"
)

var client rbac.Client

func ExampleClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := rbac.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_CreateOrUpdateRoleDefinition() {
	scope := rbac.RoleScopeGlobal
	name := uuid.New().String()
	roleType := rbac.RoleTypeCustomRole
	roleName := "ExampleRoleName"
	parameters := rbac.RoleDefinitionCreateParameters{
		Properties: &rbac.RoleDefinitionProperties{
			AssignableScopes: []*rbac.RoleScope{to.Ptr(scope)},
			Description:      to.Ptr("Example description"),
			Permissions:      []*rbac.Permission{{DataActions: []*rbac.DataAction{to.Ptr(rbac.DataActionBackupHsmKeys), to.Ptr(rbac.DataActionCreateHsmKey)}}},
			RoleName:         to.Ptr(roleName),
			RoleType:         to.Ptr(roleType),
		},
	}

	roleDefinition, err := client.CreateOrUpdateRoleDefinition(context.TODO(), scope, name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Definition Name: %s", *roleDefinition.Name)
}

func ExampleClient_CreateRoleAssignment() {
	scope := rbac.RoleScopeGlobal
	name := uuid.New().String()
	parameters := rbac.RoleAssignmentCreateParameters{
		Properties: &rbac.RoleAssignmentProperties{
			PrincipalID:      to.Ptr("d26e28bc-991f-11ed-a8fc-0242ac120002"),                                                                      // example principal ID
			RoleDefinitionID: to.Ptr("Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/c368d8da-991f-11ed-a8fc-0242ac120002"), // example role definition ID
		},
	}

	roleAssignment, err := client.CreateRoleAssignment(context.TODO(), scope, name, parameters, nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Assignment Name: %s", *roleAssignment.Name)
}
