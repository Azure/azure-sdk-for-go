//go:build go1.18
// +build go1.18

package azadmin_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var servicePrincipalId string = "aed295e0-2ae7-4c2a-9abc-813f0ca233d3"

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

func TestSetRoleDefinition(t *testing.T) {
	client := startAccessControlTest(t)
	roleDefinitionName := uuid.New()
	roleName := uuid.New()

	parameters := azadmin.RoleDefinitionCreateParameters{
		Properties: &azadmin.RoleDefinitionProperties{
			AssignableScopes: []*azadmin.RoleScope{to.Ptr(azadmin.RoleScopeGlobal)},
			Description:      to.Ptr(""),
			Permissions:      []*azadmin.Permission{},
			RoleName:         to.Ptr(roleName.String()),
			RoleType:         to.Ptr(azadmin.RoleTypeCustomRole),
		},
	}
	testSerde(t, &parameters)

	res, err := client.SetRoleDefinition(context.Background(), "/", roleDefinitionName.String(), parameters, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Properties.AssignableScopes)
	require.Equal(t, "", *res.Properties.Description)
	require.Equal(t, []*azadmin.Permission{}, res.Properties.Permissions)
	require.Equal(t, roleName.String(), *res.Properties.RoleName)
	require.NotNil(t, res.Properties.RoleType)
	require.Equal(t, "Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/"+roleDefinitionName.String(), *res.ID)
	require.Equal(t, roleDefinitionName.String(), *res.Name)

	testSerde(t, &res)

}
func TestCreateRoleAssignment(t *testing.T) {
	client := startAccessControlTest(t)
	roleDefinitionID := "Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/7b127d3c-77bd-4e3e-bbe0-dbb8971fa7f8"
	roleAssignmentID := uuid.New()

	roleAssignment := azadmin.RoleAssignmentCreateParameters{Properties: &azadmin.RoleAssignmentProperties{PrincipalID: to.Ptr(servicePrincipalId), RoleDefinitionID: to.Ptr(roleDefinitionID)}}

	res, err := client.CreateRoleAssignment(context.Background(), "/", roleAssignmentID.String(), roleAssignment, nil)
	require.NoError(t, err)
	_ = res

}

func TestGetRoleAssignment(t *testing.T) {
	//client := startAccessControlTest(t)

	//res, err := client.GetRoleAssignment(context.Background(), "", )
}

func TestNewListRoleDefinitionsPager(t *testing.T) {
	client := startAccessControlTest(t)

	pager := client.NewListRoleDefinitionsPager("/", nil)
	require.True(t, pager.More())

	for pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NotNil(t, res.Value)
		for _, roleDef := range res.Value {
			require.NotNil(t, roleDef.Properties)
			require.NotNil(t, roleDef.ID)
			require.NotNil(t, roleDef.Name)
			require.NotNil(t, roleDef.Type)
		}

		testSerde(t, &res)

	}
}

func TestNewListRoleAssignmentsPager(t *testing.T) {
	client := startAccessControlTest(t)

	pager := client.NewListRoleAssignmentsPager("", nil)
	require.True(t, pager.More())

	for pager.More() {
		res, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)

		require.NotNil(t, res.Value)
		for _, roleAssignment := range res.Value {
			require.NotNil(t, roleAssignment.Properties)
			require.NotNil(t, roleAssignment.ID)
			require.NotNil(t, roleAssignment.Name)
			require.NotNil(t, roleAssignment.Type)
		}

		testSerde(t, &res)
	}
}
