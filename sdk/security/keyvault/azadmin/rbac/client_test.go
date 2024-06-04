//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rbac_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRoleDefinition(t *testing.T) {
	client := startAccessControlTest(t)

	var name, roleName string
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		name, roleName = "bedc5fb2-8738-40a4-8b20-cedfb43a1922", "c023eb03-4e31-464c-84f7-001b7f23bd13"
	} else {
		name, roleName = uuid.New().String(), uuid.New().String()
	}
	scope := rbac.RoleScopeGlobal
	roleType := rbac.RoleTypeCustomRole
	permission := rbac.DataActionBackupHsmKeys
	parameters := rbac.RoleDefinitionCreateParameters{
		Properties: &rbac.RoleDefinitionProperties{
			AssignableScopes: []*rbac.RoleScope{to.Ptr(scope)},
			Description:      to.Ptr("test"),
			Permissions:      []*rbac.Permission{{DataActions: []*rbac.DataAction{to.Ptr(permission)}}},
			RoleName:         to.Ptr(roleName),
			RoleType:         to.Ptr(roleType),
		},
	}
	testSerde(t, &parameters)

	// test create definition
	createdDefinition, err := client.CreateOrUpdateRoleDefinition(context.Background(), scope, name, parameters, nil)
	require.NoError(t, err)
	require.Equal(t, name, *createdDefinition.Name)
	require.Len(t, createdDefinition.Properties.AssignableScopes, 1)
	require.Equal(t, scope, *createdDefinition.Properties.AssignableScopes[0])
	require.Equal(t, "test", *createdDefinition.Properties.Description)
	require.Equal(t, roleType, *createdDefinition.Properties.RoleType)
	require.Equal(t, roleName, *createdDefinition.Properties.RoleName)
	require.Len(t, createdDefinition.Properties.Permissions, 1)
	require.Equal(t, permission, *createdDefinition.Properties.Permissions[0].DataActions[0])
	testSerde(t, &createdDefinition)

	// update
	updatedPermission := rbac.DataActionCreateHsmKey
	parameters.Properties.Permissions[0].DataActions = []*rbac.DataAction{to.Ptr(updatedPermission)}
	updatedDefinition, err := client.CreateOrUpdateRoleDefinition(context.Background(), scope, name, parameters, nil)
	require.NoError(t, err)
	require.Equal(t, createdDefinition.ID, updatedDefinition.ID)
	require.Len(t, updatedDefinition.Properties.Permissions, 1)
	require.Equal(t, updatedPermission, *updatedDefinition.Properties.Permissions[0].DataActions[0])

	// get
	gotDefinition, err := client.GetRoleDefinition(context.Background(), scope, name, nil)
	require.NoError(t, err)
	require.Equal(t, *updatedDefinition.ID, *gotDefinition.ID)

	// test list role definitions and check if created definition is in list exactly once
	updatedDefinitionCount := 0
	pager := client.NewListRoleDefinitionsPager(scope, nil)
	require.True(t, pager.More())
	for pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NotNil(t, res.Value)
		for _, roleDefinition := range res.Value {
			require.NotNil(t, roleDefinition.Properties)
			require.NotNil(t, roleDefinition.ID)
			require.NotNil(t, roleDefinition.Name)
			require.NotNil(t, roleDefinition.Type)

			if *roleDefinition.ID == *updatedDefinition.ID {
				updatedDefinitionCount++
			}
		}

		testSerde(t, &res)
	}
	require.Equal(t, 1, updatedDefinitionCount)

	// test delete
	deletedDefinition, err := client.DeleteRoleDefinition(context.Background(), scope, name, nil)
	require.NoError(t, err)
	require.Equal(t, updatedDefinition.ID, deletedDefinition.ID)

	// verify role definition is deleted
	pager = client.NewListRoleDefinitionsPager(scope, nil)
	for pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		for _, roleDefinition := range res.Value {
			require.NotNil(t, roleDefinition.Properties)
			require.NotNil(t, roleDefinition.ID)
			require.NotNil(t, roleDefinition.Name)
			require.NotNil(t, roleDefinition.Type)

			if *roleDefinition.ID == *updatedDefinition.ID {
				t.Fatal("expected role definition to be deleted")
			}
		}
	}
}

func TestDeleteRoleDefinition_FailureInvalidRole(t *testing.T) {
	client := startAccessControlTest(t)
	var httpErr *azcore.ResponseError

	res, err := client.DeleteRoleDefinition(context.Background(), "", "invalidDefinition", nil)
	require.Error(t, err)
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "RoleDefinitionNotFound")
	require.Equal(t, httpErr.StatusCode, 404)
	require.Nil(t, res.Properties)
	require.Nil(t, res.ID)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)

	testSerde(t, &res)
}

func TestRoleAssignment(t *testing.T) {
	client := startAccessControlTest(t)

	var name, principalID, roleDefinitionID string
	scope := rbac.RoleScopeGlobal
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		name, principalID, roleDefinitionID = "bedc5fb2-8738-40a4-8b20-cedfb43a1922", "c023eb03-4e31-464c-84f7-001b7f23bd13", "Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/33413926-3206-4cdd-b39a-83574fe37a17"
	} else {
		name, principalID = uuid.New().String(), uuid.New().String()
		// get random role definition to use for their service principal
		pager := client.NewListRoleDefinitionsPager(scope, nil)
		require.True(t, pager.More())
		roleDefinitions, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, roleDefinitions.Value)
		roleDefinitionID = *roleDefinitions.Value[rand.Intn(len(roleDefinitions.Value))].ID
	}

	roleAssignment := rbac.RoleAssignmentCreateParameters{Properties: &rbac.RoleAssignmentProperties{PrincipalID: to.Ptr(principalID), RoleDefinitionID: to.Ptr(roleDefinitionID)}}
	testSerde(t, &roleAssignment)

	// create role assignment
	createdAssignment, err := client.CreateRoleAssignment(context.Background(), scope, name, roleAssignment, nil)
	require.NoError(t, err)
	require.Equal(t, name, *createdAssignment.Name)
	require.Equal(t, scope, *createdAssignment.Properties.Scope)
	require.Equal(t, roleDefinitionID, *createdAssignment.Properties.RoleDefinitionID)

	if recording.GetRecordMode() == recording.PlaybackMode {
		require.Equal(t, "00000000-0000-0000-0000-000000000000", *createdAssignment.Properties.PrincipalID)
	} else {
		require.Equal(t, principalID, *createdAssignment.Properties.PrincipalID)
	}

	// test if able to get role assignment
	gotAssignment, err := client.GetRoleAssignment(context.Background(), scope, name, nil)
	require.NoError(t, err)
	require.Equal(t, *createdAssignment.ID, *gotAssignment.ID)

	// test if new role assignment is in list of all assignments
	assignmentsPager := client.NewListRoleAssignmentsPager(scope, nil)
	require.True(t, assignmentsPager.More())
	var assignmentCheck bool

	for assignmentsPager.More() {
		res, err := assignmentsPager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NotNil(t, res.Value)
		for _, roleAssignment := range res.Value {
			require.NotNil(t, roleAssignment.Properties)
			require.NotNil(t, roleAssignment.ID)
			require.NotNil(t, roleAssignment.Name)
			require.NotNil(t, roleAssignment.Type)

			if *roleAssignment.ID == *createdAssignment.ID {
				assignmentCheck = true
			}
		}

		testSerde(t, &res)
	}
	require.True(t, assignmentCheck)

	// delete role assignment and check that role assignment is no longer in list
	deletedAssignment, err := client.DeleteRoleAssignment(context.Background(), scope, name, nil)
	require.NoError(t, err)
	require.Equal(t, *createdAssignment.ID, *deletedAssignment.ID)

	assignmentsPager = client.NewListRoleAssignmentsPager(scope, nil)
	require.True(t, assignmentsPager.More())

	for assignmentsPager.More() {
		res, err := assignmentsPager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NotNil(t, res.Value)
		for _, roleAssignment := range res.Value {
			require.NotNil(t, roleAssignment.Properties)
			require.NotNil(t, roleAssignment.ID)
			require.NotNil(t, roleAssignment.Name)
			require.NotNil(t, roleAssignment.Type)

			require.NotEqual(t, *roleAssignment.ID, *createdAssignment.ID)
		}

		testSerde(t, &res)
	}
}

func TestDeleteRoleAssignment_FailureInvalidRole(t *testing.T) {
	client := startAccessControlTest(t)
	var httpErr *azcore.ResponseError

	res, err := client.DeleteRoleAssignment(context.Background(), "", "invalidRole", nil)
	require.Error(t, err)
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "RoleAssignmentNotFound")
	require.Equal(t, httpErr.StatusCode, 404)
	require.Nil(t, res.Properties)
	require.Nil(t, res.ID)
	require.Nil(t, res.Name)
	require.Nil(t, res.Type)

	testSerde(t, &res)
}
