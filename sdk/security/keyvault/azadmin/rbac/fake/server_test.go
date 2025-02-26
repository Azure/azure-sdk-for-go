// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package fake_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/rbac/fake"
	"github.com/stretchr/testify/require"
)

var (
	name      = "name"
	roleScope = rbac.RoleScopeGlobal
)

func getServer() fake.Server {
	return fake.Server{
		CreateOrUpdateRoleDefinition: func(ctx context.Context, scope rbac.RoleScope, roleDefinitionName string, parameters rbac.RoleDefinitionCreateParameters, options *rbac.CreateOrUpdateRoleDefinitionOptions) (resp azfake.Responder[rbac.CreateOrUpdateRoleDefinitionResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.CreateOrUpdateRoleDefinitionResponse{
				RoleDefinition: rbac.RoleDefinition{
					Name: &roleDefinitionName,
				},
			}
			resp.SetResponse(http.StatusCreated, kvResp, nil)
			return
		},
		CreateRoleAssignment: func(ctx context.Context, scope rbac.RoleScope, roleAssignmentName string, parameters rbac.RoleAssignmentCreateParameters, options *rbac.CreateRoleAssignmentOptions) (resp azfake.Responder[rbac.CreateRoleAssignmentResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.CreateRoleAssignmentResponse{
				RoleAssignment: rbac.RoleAssignment{
					Name: &roleAssignmentName,
				},
			}
			resp.SetResponse(http.StatusCreated, kvResp, nil)
			return
		},
		DeleteRoleAssignment: func(ctx context.Context, scope rbac.RoleScope, roleAssignmentName string, options *rbac.DeleteRoleAssignmentOptions) (resp azfake.Responder[rbac.DeleteRoleAssignmentResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.DeleteRoleAssignmentResponse{
				RoleAssignment: rbac.RoleAssignment{
					Name: &roleAssignmentName,
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		DeleteRoleDefinition: func(ctx context.Context, scope rbac.RoleScope, roleDefinitionName string, options *rbac.DeleteRoleDefinitionOptions) (resp azfake.Responder[rbac.DeleteRoleDefinitionResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.DeleteRoleDefinitionResponse{
				RoleDefinition: rbac.RoleDefinition{
					Name: &roleDefinitionName,
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetRoleAssignment: func(ctx context.Context, scope rbac.RoleScope, roleAssignmentName string, options *rbac.GetRoleAssignmentOptions) (resp azfake.Responder[rbac.GetRoleAssignmentResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.GetRoleAssignmentResponse{
				RoleAssignment: rbac.RoleAssignment{
					Name: &roleAssignmentName,
					Properties: &rbac.RoleAssignmentPropertiesWithScope{
						Scope: &scope,
					},
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		GetRoleDefinition: func(ctx context.Context, scope rbac.RoleScope, roleDefinitionName string, options *rbac.GetRoleDefinitionOptions) (resp azfake.Responder[rbac.GetRoleDefinitionResponse], errResp azfake.ErrorResponder) {
			kvResp := rbac.GetRoleDefinitionResponse{
				RoleDefinition: rbac.RoleDefinition{
					Name: &roleDefinitionName,
				},
			}
			resp.SetResponse(http.StatusOK, kvResp, nil)
			return
		},
		NewListRoleAssignmentsPager: func(scope rbac.RoleScope, options *rbac.ListRoleAssignmentsOptions) (resp azfake.PagerResponder[rbac.ListRoleAssignmentsResponse]) {
			page1 := rbac.ListRoleAssignmentsResponse{
				RoleAssignmentListResult: rbac.RoleAssignmentListResult{
					Value: []*rbac.RoleAssignment{
						{
							Name: to.Ptr(name + "1"),
						},
						{
							Name: to.Ptr(name + "2"),
						},
					},
				},
			}
			page2 := rbac.ListRoleAssignmentsResponse{
				RoleAssignmentListResult: rbac.RoleAssignmentListResult{
					Value: []*rbac.RoleAssignment{
						{
							Name: to.Ptr(name + "3"),
						},
						{
							Name: to.Ptr(name + "4"),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
		NewListRoleDefinitionsPager: func(scope rbac.RoleScope, options *rbac.ListRoleDefinitionsOptions) (resp azfake.PagerResponder[rbac.ListRoleDefinitionsResponse]) {
			page1 := rbac.ListRoleDefinitionsResponse{
				RoleDefinitionListResult: rbac.RoleDefinitionListResult{
					Value: []*rbac.RoleDefinition{
						{
							Name: to.Ptr(name + "1"),
						},
						{
							Name: to.Ptr(name + "2"),
						},
					},
				},
			}
			page2 := rbac.ListRoleDefinitionsResponse{
				RoleDefinitionListResult: rbac.RoleDefinitionListResult{
					Value: []*rbac.RoleDefinition{
						{
							Name: to.Ptr(name + "3"),
						},
						{
							Name: to.Ptr(name + "4"),
						},
					},
				},
			}

			resp.AddPage(http.StatusOK, page1, nil)
			resp.AddPage(http.StatusOK, page2, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := rbac.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &rbac.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	createResp, err := client.CreateOrUpdateRoleDefinition(context.Background(), roleScope, name, rbac.RoleDefinitionCreateParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, name, *createResp.Name)

	createRoleAssignResp, err := client.CreateRoleAssignment(context.Background(), roleScope, name, rbac.RoleAssignmentCreateParameters{}, nil)
	require.NoError(t, err)
	require.Equal(t, name, *createRoleAssignResp.Name)

	deleteRoleAssignResp, err := client.DeleteRoleAssignment(context.Background(), roleScope, name, nil)
	require.NoError(t, err)
	require.Equal(t, name, *deleteRoleAssignResp.Name)

	deleteRoleDefResp, err := client.DeleteRoleDefinition(context.Background(), roleScope, name, nil)
	require.NoError(t, err)
	require.Equal(t, name, *deleteRoleDefResp.Name)

	getRoleAssignResp, err := client.GetRoleAssignment(context.Background(), roleScope, name, nil)
	require.NoError(t, err)
	require.Equal(t, name, *getRoleAssignResp.Name)
	require.Equal(t, roleScope, *getRoleAssignResp.Properties.Scope)

	getRoleDefResp, err := client.GetRoleDefinition(context.Background(), roleScope, name, nil)
	require.NoError(t, err)
	require.Equal(t, name, *getRoleDefResp.Name)

	roleAssignPager := client.NewListRoleAssignmentsPager(roleScope, nil)
	for roleAssignPager.More() {
		page, err := roleAssignPager.NextPage(context.Background())
		require.NoError(t, err)

		for _, roleAssign := range page.Value {
			require.Contains(t, *roleAssign.Name, name)
		}
	}

	roleDefPager := client.NewListRoleDefinitionsPager(roleScope, nil)
	for roleDefPager.More() {
		page, err := roleDefPager.NextPage(context.Background())
		require.NoError(t, err)

		for _, roleAssign := range page.Value {
			require.Contains(t, *roleAssign.Name, name)
		}
	}
}
