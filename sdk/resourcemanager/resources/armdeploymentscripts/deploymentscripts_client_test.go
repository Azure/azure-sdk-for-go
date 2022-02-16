//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdeploymentscripts_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentscripts"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeploymentScriptsClient_BeginCreate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)
}

func TestDeploymentScriptsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)

	// get deployment scripts
	getResp, err := deploymentScriptsClient.Get(ctx, rgName, scriptName, nil)
	require.NoError(t, err)
	require.Equal(t, scriptName, *getResp.GetDeploymentScript().Name)
}

func TestDeploymentScriptsClient_GetLogs(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)

	// get log
	getResp, err := deploymentScriptsClient.GetLogs(ctx, rgName, scriptName, nil)
	require.NoError(t, err)
	require.Equal(t, 1, len(getResp.Value))
}

func TestDeploymentScriptsClient_GetLogsDefault(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)

	getResp, err := deploymentScriptsClient.GetLogsDefault(ctx, rgName, scriptName, nil)
	require.NoError(t, err)
	require.Equal(t, "default", *getResp.Name)
}

func TestDeploymentScriptsClient_Update(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)

	// update deployment script
	updateResp, err := deploymentScriptsClient.Update(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.ClientUpdateOptions{
			DeploymentScript: &armdeploymentscripts.DeploymentScriptUpdateParameter{
				Tags: map[string]*string{
					"test": to.StringPtr("recording"),
				},
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "recording", *updateResp.GetDeploymentScript().Tags["test"])
}

func TestDeploymentScriptsClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create identity
	userAssignedIdentitiesClient := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, opt)
	identityName, err := createRandomName(t, "identity")
	require.NoError(t, err)
	identityResp, err := userAssignedIdentitiesClient.CreateOrUpdate(
		ctx,
		rgName,
		identityName,
		armmsi.Identity{
			Location: to.StringPtr(location),
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, identityName, *identityResp.Name)

	// create deployment script
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	scriptName, err := createRandomName(t, "script")
	require.NoError(t, err)
	dsPollerResp, err := deploymentScriptsClient.BeginCreate(
		ctx,
		rgName,
		scriptName,
		&armdeploymentscripts.AzurePowerShellScript{
			Identity: &armdeploymentscripts.ManagedServiceIdentity{
				Type: armdeploymentscripts.ManagedServiceIdentityTypeUserAssigned.ToPtr(),
				UserAssignedIdentities: map[string]*armdeploymentscripts.UserAssignedIdentity{
					*identityResp.ID: {},
				},
			},
			Kind:     armdeploymentscripts.ScriptTypeAzurePowerShell.ToPtr(),
			Location: to.StringPtr(location),
			Properties: &armdeploymentscripts.AzurePowerShellScriptProperties{
				RetentionInterval:   to.StringPtr("PT26H"),
				PrimaryScriptURI:    to.StringPtr("https://raw.githubusercontent.com/Azure/azure-docs-json-samples/master/deployment-script/deploymentscript-helloworld.ps1"),
				Arguments:           to.StringPtr("-name \"John Dole\""),
				Timeout:             to.StringPtr("PT30M"),
				AzPowerShellVersion: to.StringPtr("3.0"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dsResp, err := dsPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scriptName, *dsResp.GetDeploymentScript().Name)

	// delete deployment script
	delResp, err := deploymentScriptsClient.Delete(ctx, rgName, scriptName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestDeploymentScriptsClient_ListBySubscription(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// list deployment script by subscription
	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)
	pager := deploymentScriptsClient.ListBySubscription(nil)
	require.NoError(t, pager.Err())
}

func TestDeploymentScriptsClient_ListByResourceGroup(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	deploymentScriptsClient := armdeploymentscripts.NewClient(subscriptionID, cred, opt)

	// list deployment script by resource group
	pager := deploymentScriptsClient.ListByResourceGroup(rgName, nil)
	require.NoError(t, pager.Err())
}
