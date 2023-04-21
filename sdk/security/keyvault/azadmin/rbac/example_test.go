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

func ExampleNewClient() {
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
	roleName := "<role name>"
	parameters := rbac.RoleDefinitionCreateParameters{
		Properties: &rbac.RoleDefinitionProperties{
			AssignableScopes: []*rbac.RoleScope{to.Ptr(scope)},
			Description:      to.Ptr("<description>"),
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

func ExampleClient_DeleteRoleAssignment() {
	deletedRoleAssignment, err := client.DeleteRoleAssignment(context.Background(), rbac.RoleScopeGlobal, "<role assignment name>", nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Deleted Role Assignment Name: %s", *deletedRoleAssignment.Name)
}

func ExampleClient_DeleteRoleDefinition() {
	deletedRoleDefinition, err := client.DeleteRoleDefinition(context.Background(), rbac.RoleScopeGlobal, "<role definition name>", nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Deleted Role Definition Name: %s", *deletedRoleDefinition.Name)
}

func ExampleClient_GetRoleAssignment() {
	roleAssignment, err := client.GetRoleAssignment(context.Background(), rbac.RoleScopeGlobal, "<role assignment name>", nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Assignment Name: %s", *roleAssignment.Name)
}

func ExampleClient_GetRoleDefinition() {
	roleDefinition, err := client.GetRoleDefinition(context.Background(), rbac.RoleScopeGlobal, "<role definition name>", nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Role Definition Name: %s", *roleDefinition.Name)
}

func ExampleClient_NewListRoleAssignmentsPager() {
	pager := client.NewListRoleAssignmentsPager(rbac.RoleScopeGlobal, nil)

	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			//TODO: handle error
		}
		fmt.Println("Role Assignment Name List")
		for index, roleAssignment := range nextResult.Value {
			fmt.Printf("%d) %s\n", index, *roleAssignment.Name)
		}
	}
}

func ExampleClient_NewListRoleDefinitionsPager() {
	pager := client.NewListRoleAssignmentsPager(rbac.RoleScopeGlobal, nil)

	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			//TODO: handle error
		}
		fmt.Println("Role Definition Name List")
		for index, roleDefinition := range nextResult.Value {
			fmt.Printf("%d) %s\n", index, *roleDefinition.Name)
		}
	}
}
