//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rbac_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"

var (
	credential azcore.TokenCredential
	hsmURL     string
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		hsmURL = os.Getenv("AZURE_MANAGEDHSM_URL")
	}
	if hsmURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_MANAGEDHSM_URL")
		}
		hsmURL = fakeHsmURL
	}

	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("KEYVAULT_TENANT_ID")
		clientID := lookupEnvVar("KEYVAULT_CLIENT_ID")
		secret := lookupEnvVar("KEYVAULT_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/security/keyvault/azadmin/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startAccessControlTest(t *testing.T) *rbac.Client {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &rbac.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := rbac.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}

func TestRoleDefinition(t *testing.T) {
	client := startAccessControlTest(t)

	name := uuid.New().String()
	scope := rbac.RoleScopeGlobal
	roleType := rbac.RoleTypeCustomRole
	roleName := uuid.New().String()
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

	scope := rbac.RoleScopeGlobal
	name := uuid.New().String()
	principalID := uuid.New().String()

	// get random role definition to use for their service principal
	pager := client.NewListRoleDefinitionsPager(scope, nil)
	require.True(t, pager.More())
	roleDefinitions, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	require.NotNil(t, roleDefinitions.Value)
	roleDefinition := roleDefinitions.Value[rand.Intn(len(roleDefinitions.Value))]
	roleAssignment := rbac.RoleAssignmentCreateParameters{Properties: &rbac.RoleAssignmentProperties{PrincipalID: to.Ptr(principalID), RoleDefinitionID: roleDefinition.ID}}
	testSerde(t, &roleAssignment)

	// create role assignment
	createdAssignment, err := client.CreateRoleAssignment(context.Background(), scope, name, roleAssignment, nil)
	require.NoError(t, err)
	require.Equal(t, name, *createdAssignment.Name)
	require.Equal(t, scope, *createdAssignment.Properties.Scope)
	require.Equal(t, principalID, *createdAssignment.Properties.PrincipalID)
	require.Equal(t, *roleDefinition.ID, *createdAssignment.Properties.RoleDefinitionID)

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
