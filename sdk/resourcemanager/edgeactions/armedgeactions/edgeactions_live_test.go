// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func await[T any](t *testing.T, ctx context.Context, poller *runtime.Poller[T]) *T {
	t.Helper()
	result, err := pollForScenario(ctx, poller)
	require.NoError(t, err)
	assert.True(t, poller.Done())
	return result
}

func requireStatus(t *testing.T, err error, status int) {
	t.Helper()
	var responseError *azcore.ResponseError
	require.ErrorAs(t, err, &responseError)
	require.Equal(t, status, responseError.StatusCode)
}

func newScenario(t *testing.T) (context.Context, *armedgeactions.ClientFactory, string, string) {
	t.Helper()
	if recording.GetRecordMode() != recording.PlaybackMode {
		require.NotEmpty(t, recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", ""), "set an authorized test subscription before a live/record run")
	}
	stop := testutil.StartRecording(t, pathToPackage)
	t.Cleanup(stop)
	require.NoError(t, recording.SetDefaultMatcher(t, &recording.SetDefaultMatcherOptions{}))
	ctx := t.Context()
	cred, options := testutil.GetCredAndClientOptions(t)
	options.Retry.TryTimeout = 2 * time.Minute
	spec, err := loadResponseSpec()
	require.NoError(t, err)
	validator := &responseValidator{spec: spec}
	options.PerCallPolicies = append(options.PerCallPolicies, validator)
	validationClient, err := arm.NewClient("edgeactions-tests", "v0.0.0", cred, options)
	require.NoError(t, err)
	validator.pipeline = validationClient.Pipeline()
	ctx = context.WithValue(ctx, validatorKey{}, validator)
	sub := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", subscriptionID)
	factory, err := armedgeactions.NewClientFactory(sub, cred, options)
	require.NoError(t, err)
	groups, err := armresources.NewResourceGroupsClient(sub, cred, options)
	require.NoError(t, err)
	group, err := recording.GenerateAlphaNumericID(t, "go-sdk-test-edgeactions-", 40, true)
	require.NoError(t, err)
	name, err := recording.GenerateAlphaNumericID(t, "eago", 24, true)
	require.NoError(t, err)
	_, err = groups.Get(ctx, group, nil)
	requireStatus(t, err, http.StatusNotFound)
	// Register before PUT to cover partial creation; never delete a pre-existing group.
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		poller, err := groups.BeginDelete(cleanupCtx, group, nil)
		if hasStatus(err, http.StatusNotFound) {
			return
		}
		require.NoError(t, err)
		await(t, cleanupCtx, poller)
	})
	_, err = groups.CreateOrUpdate(ctx, group, armresources.ResourceGroup{
		Location: to.Ptr(recording.GetEnvVariable("LOCATION", "eastus")),
	}, nil)
	require.NoError(t, err)
	cleanup, err := createOwnedAction(ctx, factory, group, name)
	if cleanup != nil {
		t.Cleanup(func() {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			require.NoError(t, cleanup(cleanupCtx), "clean up owned EdgeAction")
		})
	}
	require.NoError(t, err)
	return ctx, factory, group, name
}

func TestRecordedEdgeActionTags(t *testing.T) {
	ctx, factory, group, name := newScenario(t)
	client := factory.NewClient()
	before, err := client.Get(ctx, group, name, nil)
	require.NoError(t, err)
	assert.Equal(t, to.Ptr("global"), before.Location)
	assert.Equal(t, map[string]*string{"purpose": to.Ptr("initial"), "remove-me": to.Ptr("initial")}, before.Tags)
	pager := client.NewListByResourceGroupPager(group, nil)
	var names []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		for _, action := range page.Value {
			require.NotNil(t, action)
			require.NotNil(t, action.Name)
			names = append(names, *action.Name)
		}
	}
	assert.Contains(t, names, name)
	poller, err := client.BeginUpdate(ctx, group, name, armedgeactions.EdgeActionUpdate{}, nil)
	if err == nil {
		_, err = pollForScenario(ctx, poller)
	}
	requireStatus(t, err, http.StatusBadRequest)
	unchanged, err := client.Get(ctx, group, name, nil)
	require.NoError(t, err)
	assert.Equal(t, before.Tags, unchanged.Tags)

	for _, tags := range []map[string]*string{{"purpose": to.Ptr("updated")}, {}} {
		poller, err := client.BeginUpdate(ctx, group, name, armedgeactions.EdgeActionUpdate{Tags: tags}, nil)
		require.NoError(t, err)
		result := await(t, ctx, poller)
		assert.Equal(t, tags, result.Tags)
		assert.Equal(t, before.Location, result.Location)
		assert.Equal(t, before.SKU, result.SKU)
		updated, err := client.Get(ctx, group, name, nil)
		require.NoError(t, err)
		assert.Equal(t, tags, updated.Tags)
		assert.Equal(t, before.Location, updated.Location)
		assert.Equal(t, before.SKU, updated.SKU)
	}
	require.NotNil(t, before.SKU)
	for _, update := range []armedgeactions.EdgeActionUpdate{
		{Tags: azcore.NullValue[map[string]*string]()},
		{SKU: azcore.NullValue[*armedgeactions.SKUTypeUpdate]()},
		{SKU: &armedgeactions.SKUTypeUpdate{Name: before.SKU.Name, Tier: before.SKU.Tier}},
	} {
		poller, err := client.BeginUpdate(ctx, group, name, update, nil)
		if err == nil {
			_, err = pollForScenario(ctx, poller)
		}
		requireStatus(t, err, http.StatusBadRequest)
	}
	unchanged, err = client.Get(ctx, group, name, nil)
	require.NoError(t, err)
	assert.Empty(t, unchanged.Tags)
	require.NoError(t, cleanupAction(ctx, factory, group, name))
	_, err = client.Get(ctx, group, name, nil)
	requireStatus(t, err, http.StatusNotFound)
}

func TestRecordedEdgeActionVersionsAndFilters(t *testing.T) {
	for _, deploymentType := range []armedgeactions.EdgeActionVersionDeploymentType{
		armedgeactions.EdgeActionVersionDeploymentTypeFile, armedgeactions.EdgeActionVersionDeploymentTypeZip,
	} {
		t.Run(string(deploymentType), func(t *testing.T) {
			ctx, factory, group, name := newScenario(t)
			versions := factory.NewEdgeActionVersionsClient()
			content, err := versionCode(deploymentType)
			require.NoError(t, err)
			for _, version := range []string{"v1", "v2"} {
				poller, err := versions.BeginCreate(ctx, group, name, version, armedgeactions.EdgeActionVersion{
					Location: to.Ptr("global"),
					Properties: &armedgeactions.EdgeActionVersionProperties{
						DeploymentType:   to.Ptr(deploymentType),
						IsDefaultVersion: to.Ptr(armedgeactions.EdgeActionIsDefaultVersionFalse),
					},
				}, nil)
				require.NoError(t, err)
				created := await(t, ctx, poller)
				require.NotNil(t, created.Properties)
				require.NotNil(t, created.ID)
				assert.Equal(t, to.Ptr(version), created.Name)
				assert.Equal(t, to.Ptr(deploymentType), created.Properties.DeploymentType)
				assert.Equal(t, to.Ptr(armedgeactions.ProvisioningStateSucceeded), created.Properties.ProvisioningState)
				deploy, err := versions.BeginDeployVersionCode(ctx, group, name, version, armedgeactions.VersionCode{
					Name: to.Ptr(version), Content: to.Ptr(content),
				}, nil)
				require.NoError(t, err)
				result := await(t, ctx, deploy)
				require.NotNil(t, result.Properties, "deployVersionCode must deserialize the full version envelope")
				assert.Equal(t, created.ID, result.ID)
				assert.Equal(t, created.Name, result.Name)
				assert.Equal(t, created.Location, result.Location)
				assert.Equal(t, to.Ptr(deploymentType), result.Properties.DeploymentType)
				assert.Equal(t, to.Ptr(armedgeactions.ProvisioningStateSucceeded), result.Properties.ProvisioningState)
				deployed, err := versions.Get(ctx, group, name, version, nil)
				require.NoError(t, err)
				require.NotNil(t, deployed.Properties)
				assert.Equal(t, to.Ptr(armedgeactions.ProvisioningStateSucceeded), deployed.Properties.ProvisioningState)
				assert.Equal(t, to.Ptr(deploymentType), deployed.Properties.DeploymentType)
			}
			first, err := versions.Get(ctx, group, name, "v1", nil)
			require.NoError(t, err)
			second, err := versions.Get(ctx, group, name, "v2", nil)
			require.NoError(t, err)
			require.NotNil(t, first.Properties)
			require.NotNil(t, second.Properties)
			assert.Equal(t, to.Ptr(armedgeactions.EdgeActionIsDefaultVersionTrue), first.Properties.IsDefaultVersion)
			assert.Equal(t, to.Ptr(armedgeactions.EdgeActionIsDefaultVersionFalse), second.Properties.IsDefaultVersion)
			pager := versions.NewListByEdgeActionPager(group, name, nil)
			var versionNames []string
			for pager.More() {
				page, err := pager.NextPage(ctx)
				require.NoError(t, err)
				for _, version := range page.Value {
					require.NotNil(t, version)
					require.NotNil(t, version.Name)
					versionNames = append(versionNames, *version.Name)
				}
			}
			assert.ElementsMatch(t, []string{"v1", "v2"}, versionNames)
			getCode, err := versions.BeginGetVersionCode(ctx, group, name, "v1", nil)
			require.NoError(t, err)
			code := await(t, ctx, getCode)
			require.NotNil(t, code.Content)
			assertHandler(t, *code.Content)

			update, err := versions.BeginUpdate(ctx, group, name, "v1", armedgeactions.EdgeActionVersionUpdate{
				Tags: map[string]*string{"purpose": to.Ptr("updated")},
			}, nil)
			require.NoError(t, err)
			updateResult := await(t, ctx, update)
			assert.Equal(t, map[string]*string{"purpose": to.Ptr("updated")}, updateResult.Tags)
			require.NotNil(t, updateResult.Properties)
			assert.Equal(t, first.Properties.DeploymentType, updateResult.Properties.DeploymentType)
			assert.Equal(t, first.Properties.IsDefaultVersion, updateResult.Properties.IsDefaultVersion)
			updated, err := versions.Get(ctx, group, name, "v1", nil)
			require.NoError(t, err)
			require.NotNil(t, updated.Properties)
			assert.Equal(t, map[string]*string{"purpose": to.Ptr("updated")}, updated.Tags)
			assert.Equal(t, first.Properties.DeploymentType, updated.Properties.DeploymentType)
			assert.Equal(t, first.Properties.IsDefaultVersion, updated.Properties.IsDefaultVersion)
			update, err = versions.BeginUpdate(ctx, group, name, "v1", armedgeactions.EdgeActionVersionUpdate{
				Properties: &armedgeactions.EdgeActionVersionUpdateProperties{
					DeploymentType: to.Ptr(deploymentType), IsDefaultVersion: to.Ptr(armedgeactions.EdgeActionIsDefaultVersionTrue),
				},
			}, nil)
			require.NoError(t, err)
			await(t, ctx, update)
			otherType := armedgeactions.EdgeActionVersionDeploymentTypeZip
			if deploymentType == otherType {
				otherType = armedgeactions.EdgeActionVersionDeploymentTypeFile
			}
			for _, properties := range []*armedgeactions.EdgeActionVersionUpdateProperties{
				{IsDefaultVersion: to.Ptr(armedgeactions.EdgeActionIsDefaultVersionFalse)},
				{DeploymentType: to.Ptr(otherType)},
			} {
				poller, err := versions.BeginUpdate(ctx, group, name, "v1", armedgeactions.EdgeActionVersionUpdate{Properties: properties}, nil)
				if err == nil {
					_, err = pollForScenario(ctx, poller)
				}
				requireStatus(t, err, http.StatusBadRequest)
			}
			swap, err := versions.BeginSwapDefault(ctx, group, name, "v2", nil)
			require.NoError(t, err)
			await(t, ctx, swap)
			for _, version := range []struct {
				name     string
				expected armedgeactions.EdgeActionIsDefaultVersion
			}{{"v1", armedgeactions.EdgeActionIsDefaultVersionFalse}, {"v2", armedgeactions.EdgeActionIsDefaultVersionTrue}} {
				result, err := versions.Get(ctx, group, name, version.name, nil)
				require.NoError(t, err)
				require.NotNil(t, result.Properties)
				assert.Equal(t, to.Ptr(version.expected), result.Properties.IsDefaultVersion)
			}

			filters := factory.NewEdgeActionExecutionFiltersClient()
			require.NotNil(t, first.ID)
			require.NotNil(t, second.ID)
			createFilter, err := filters.BeginCreate(ctx, group, name, "sdkfilter", armedgeactions.EdgeActionExecutionFilter{
				Location: to.Ptr("global"),
				Properties: &armedgeactions.EdgeActionExecutionFilterProperties{
					VersionID: first.ID, ExecutionFilterIdentifierHeaderName: to.Ptr("x-sdk-test"),
					ExecutionFilterIdentifierHeaderValue: to.Ptr("initial"),
				},
			}, nil)
			require.NoError(t, err)
			filterResult := await(t, ctx, createFilter)
			require.NotNil(t, filterResult.Properties)
			require.NotNil(t, filterResult.Properties.VersionID)
			assert.True(t, strings.EqualFold(*first.ID, *filterResult.Properties.VersionID))
			assert.Equal(t, to.Ptr(armedgeactions.ProvisioningStateSucceeded), filterResult.Properties.ProvisioningState)
			initialFilter, err := filters.Get(ctx, group, name, "sdkfilter", nil)
			require.NoError(t, err)
			require.NotNil(t, initialFilter.Properties)
			require.NotNil(t, initialFilter.Properties.VersionID)
			assert.True(t, strings.EqualFold(*first.ID, *initialFilter.Properties.VersionID))
			filterPager := filters.NewListByEdgeActionPager(group, name, nil)
			var filterNames []string
			for filterPager.More() {
				page, err := filterPager.NextPage(ctx)
				require.NoError(t, err)
				for _, filter := range page.Value {
					require.NotNil(t, filter)
					require.NotNil(t, filter.Name)
					filterNames = append(filterNames, *filter.Name)
				}
			}
			assert.Equal(t, []string{"sdkfilter"}, filterNames)
			updateFilter, err := filters.BeginUpdate(ctx, group, name, "sdkfilter", armedgeactions.EdgeActionExecutionFilterUpdate{
				Tags: map[string]*string{"purpose": to.Ptr("updated")},
				Properties: &armedgeactions.EdgeActionExecutionFilterUpdateProperties{
					VersionID: second.ID, ExecutionFilterIdentifierHeaderName: to.Ptr("x-sdk-test-updated"),
					ExecutionFilterIdentifierHeaderValue: to.Ptr("updated"),
				},
			}, nil)
			require.NoError(t, err)
			filterUpdate := await(t, ctx, updateFilter)
			require.NotNil(t, filterUpdate.Properties)
			require.NotNil(t, filterUpdate.Properties.VersionID)
			assert.True(t, strings.EqualFold(*second.ID, *filterUpdate.Properties.VersionID))
			assert.Equal(t, map[string]*string{"purpose": to.Ptr("updated")}, filterUpdate.Tags)
			assert.Equal(t, to.Ptr("x-sdk-test-updated"), filterUpdate.Properties.ExecutionFilterIdentifierHeaderName)
			assert.Equal(t, to.Ptr("updated"), filterUpdate.Properties.ExecutionFilterIdentifierHeaderValue)
			updatedFilter, err := filters.Get(ctx, group, name, "sdkfilter", nil)
			require.NoError(t, err)
			assert.Equal(t, map[string]*string{"purpose": to.Ptr("updated")}, updatedFilter.Tags)
			require.NotNil(t, updatedFilter.Properties)
			require.NotNil(t, updatedFilter.Properties.VersionID)
			assert.True(t, strings.EqualFold(*second.ID, *updatedFilter.Properties.VersionID))
			assert.Equal(t, to.Ptr("x-sdk-test-updated"), updatedFilter.Properties.ExecutionFilterIdentifierHeaderName)
			assert.Equal(t, to.Ptr("updated"), updatedFilter.Properties.ExecutionFilterIdentifierHeaderValue)
			deleteFilter, err := filters.BeginDelete(ctx, group, name, "sdkfilter", nil)
			require.NoError(t, err)
			await(t, ctx, deleteFilter)
			filterPager = filters.NewListByEdgeActionPager(group, name, nil)
			for filterPager.More() {
				page, err := filterPager.NextPage(ctx)
				require.NoError(t, err)
				assert.Empty(t, page.Value)
			}
			for _, version := range []string{"v1", "v2"} {
				poller, err := versions.BeginDelete(ctx, group, name, version, nil)
				require.NoError(t, err)
				await(t, ctx, poller)
			}
			pager = versions.NewListByEdgeActionPager(group, name, nil)
			for pager.More() {
				page, err := pager.NextPage(ctx)
				require.NoError(t, err)
				assert.Empty(t, page.Value)
			}
		})
	}
}

func assertHandler(t *testing.T, content string) {
	t.Helper()
	decoded, err := base64.StdEncoding.DecodeString(content)
	require.NoError(t, err)
	if bytes.Equal(decoded, []byte(handlerCode)) {
		return
	}
	// The service may return a ZIP package even for a single-file deployment.
	archive, err := zip.NewReader(bytes.NewReader(decoded), int64(len(decoded)))
	require.NoError(t, err)
	var contents []string
	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}
		reader, err := file.Open()
		require.NoError(t, err)
		data, readErr := io.ReadAll(reader)
		closeErr := reader.Close()
		require.NoError(t, readErr)
		require.NoError(t, closeErr)
		contents = append(contents, string(data))
	}
	assert.Contains(t, contents, handlerCode)
}
