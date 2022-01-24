//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armpolicy_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/stretchr/testify/require"
	"testing"
)

var policyDefinitionOptions = armpolicy.Definition{
	Properties: &armpolicy.DefinitionProperties{
		PolicyType:  armpolicy.PolicyTypeCustom.ToPtr(),
		Description: to.StringPtr("test case"),
		Parameters: map[string]*armpolicy.ParameterDefinitionsValue{
			"prefix": {
				Type: armpolicy.ParameterTypeString.ToPtr(),
				Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
					Description: to.StringPtr("prefix description"),
					DisplayName: to.StringPtr("test case prefix"),
				},
			},
			"suffix": {
				Type: armpolicy.ParameterTypeString.ToPtr(),
				Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
					Description: to.StringPtr("suffix description"),
					DisplayName: to.StringPtr("test case suffix"),
				},
			},
		},
		PolicyRule: map[string]interface{}{
			"if": map[string]interface{}{
				"not": map[string]interface{}{
					"field": "name",
					"like":  "[concat(parameters('prefix'), '*', parameters('suffix'))]",
				},
			},
			"then": map[string]interface{}{
				"effect": "deny",
			},
		},
	},
}

func TestPolicyDefinitionsClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create policy definition
	policyDefinitionName, err := createRandomName(t, "definition")
	require.NoError(t, err)
	policyDefinitionsClient := armpolicy.NewDefinitionsClient(subscriptionID, cred, opt)
	cResp, err := policyDefinitionsClient.CreateOrUpdate(ctx, policyDefinitionName, policyDefinitionOptions, nil)
	require.NoError(t, err)
	require.Equal(t, policyDefinitionName, *cResp.Name)
}

func TestPolicyDefinitionsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create policy definition
	policyDefinitionName, err := createRandomName(t, "definition")
	require.NoError(t, err)
	policyDefinitionsClient := armpolicy.NewDefinitionsClient(subscriptionID, cred, opt)
	cResp, err := policyDefinitionsClient.CreateOrUpdate(ctx, policyDefinitionName, policyDefinitionOptions, nil)
	require.NoError(t, err)
	require.Equal(t, policyDefinitionName, *cResp.Name)

	// get policy definition
	getResp, err := policyDefinitionsClient.Get(ctx, policyDefinitionName, nil)
	require.NoError(t, err)
	require.Equal(t, policyDefinitionName, *getResp.Name)
}

func TestPolicyDefinitionsClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create policy definition
	policyDefinitionName, err := createRandomName(t, "definition")
	require.NoError(t, err)
	policyDefinitionsClient := armpolicy.NewDefinitionsClient(subscriptionID, cred, opt)
	cResp, err := policyDefinitionsClient.CreateOrUpdate(ctx, policyDefinitionName, policyDefinitionOptions, nil)
	require.NoError(t, err)
	require.Equal(t, policyDefinitionName, *cResp.Name)

	// delete policy definition
	delResp, err := policyDefinitionsClient.Delete(ctx, policyDefinitionName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestPolicyDefinitionsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	policyDefinitionsClient := armpolicy.NewDefinitionsClient(subscriptionID, cred, opt)

	// list policy definition
	pager := policyDefinitionsClient.List(nil)
	require.NoError(t, pager.Err())
	require.False(t, pager.NextPage(context.Background()))
}

func TestPolicyDefinitionsClient_ListBuiltIn(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	policyDefinitionsClient := armpolicy.NewDefinitionsClient(subscriptionID, cred, opt)

	// list policy definition
	pager := policyDefinitionsClient.ListBuiltIn(nil)
	require.NoError(t, pager.Err())
	require.False(t, pager.NextPage(context.Background()))
}
