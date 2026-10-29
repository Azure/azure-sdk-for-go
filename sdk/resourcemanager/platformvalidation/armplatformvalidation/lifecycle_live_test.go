// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armplatformvalidation_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	pv "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/platformvalidation/armplatformvalidation"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
	"github.com/stretchr/testify/require"
)

type ownedValidation struct {
	group, name, plan, run string
	mutationInProgress     bool
}

func completeOperation[T any](t *testing.T, ctx context.Context, owned *ownedValidation,
	begin func() (*runtime.Poller[T], error)) T {
	t.Helper()
	owned.mutationInProgress = true
	poller, err := begin()
	requireNoError(t, err)
	result, err := testutil.PollForTest(ctx, poller)
	owned.mutationInProgress = !poller.Done()
	requireNoError(t, err)
	return *result
}

func TestPlatformValidationBootLifecycle(t *testing.T) {
	clients, ctx := newTestClients(t)
	source := recording.GetEnvVariable("PLATFORMVALIDATION_SOURCE_VHD_URI", "https://example.invalid/image.vhd")
	parsedSource, err := url.Parse(source)
	requireNoError(t, err)
	require.True(t, parsedSource.Scheme == "https" && parsedSource.Host != "",
		"An approved HTTPS Linux Gen1 x64 VHD is required")
	if recording.GetRecordMode() != recording.PlaybackMode {
		require.NotEmpty(t, recording.GetEnvVariable("PLATFORMVALIDATION_SOURCE_VHD_URI", ""),
			"Set the approved VHD privately before starting live tests")
	}
	name, err := recording.GenerateAlphaNumericID(t, "gopv", 20, true)
	requireNoError(t, err)
	owned := &ownedValidation{group: "rg-" + name, name: name, plan: "go-boot-plan", run: "go-boot-run"}
	t.Logf("Owned resources: group=%s; cloudValidation=%s", owned.group, owned.name)
	groups, err := armresources.NewResourceGroupsClient(clients.subscriptionID, clients.credential, clients.options)
	requireNoError(t, err)
	_, err = groups.Get(ctx, owned.group, nil)
	requireNotFound(t, err)
	_, err = groups.Get(ctx, owned.name+"-mrg", nil)
	requireNotFound(t, err)
	tags := map[string]*string{"sdk-test": to.Ptr(owned.name)}
	_, err = groups.CreateOrUpdate(ctx, owned.group, armresources.ResourceGroup{
		Location: to.Ptr(clients.location), Tags: tags,
	}, nil)
	requireNoError(t, err)
	t.Cleanup(func() { cleanupValidation(t, clients, owned) })

	cvs := clients.factory.NewCloudValidationsClient()
	plans := clients.factory.NewValidationExecutionPlansClient()
	runs := clients.factory.NewExecutionPlanRunsClient()
	_, err = cvs.Get(ctx, owned.group, owned.name, nil)
	requireNotFound(t, err)
	cv := completeOperation(t, ctx, owned, func() (*runtime.Poller[pv.CloudValidationsClientCreateOrUpdateResponse], error) {
		return cvs.BeginCreateOrUpdate(ctx, owned.group, owned.name, pv.CloudValidation{
			Location: to.Ptr(clients.location), Tags: tags,
			Properties: &pv.CloudValidationProperties{Description: to.Ptr("Go SDK live validation")},
		}, nil)
	})
	require.NotNil(t, cv.Properties)
	require.Equal(t, "Succeeded", stringValue(cv.Properties.ProvisioningState))
	require.Equal(t, strings.ToLower("/subscriptions/"+clients.subscriptionID+"/resourceGroups/"+owned.group+
		"/providers/Microsoft.PlatformValidation/cloudValidations/"+owned.name), strings.ToLower(stringValue(cv.ID)))
	fetchedCV, err := cvs.Get(ctx, owned.group, owned.name, nil)
	requireNoError(t, err)
	require.Equal(t, cv.ID, fetchedCV.ID)
	found := false
	groupPager := cvs.NewListByResourceGroupPager(owned.group, nil)
	for groupPager.More() {
		page, err := groupPager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			found = found || stringValue(entry.ID) == stringValue(cv.ID)
		}
	}
	require.True(t, found, "Group LIST must include the owned CV")
	found = false
	subscriptionPager := cvs.NewListBySubscriptionPager(nil)
	for subscriptionPager.More() {
		page, err := subscriptionPager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			found = found || stringValue(entry.ID) == stringValue(cv.ID)
		}
	}
	require.True(t, found, "Subscription LIST must include the owned CV")
	completeOperation(t, ctx, owned, func() (*runtime.Poller[pv.CloudValidationsClientUpdateResponse], error) {
		return cvs.BeginUpdate(ctx, owned.group, owned.name, pv.CloudValidationUpdate{
			Tags:       map[string]*string{"sdk-test": to.Ptr(owned.name), "phase": to.Ptr("updated")},
			Properties: &pv.CloudValidationUpdateProperties{Description: to.Ptr("Go SDK updated validation")},
		}, nil)
	})
	fetchedCV, err = cvs.Get(ctx, owned.group, owned.name, nil)
	requireNoError(t, err)
	require.NotNil(t, fetchedCV.Properties)
	require.Equal(t, "Go SDK updated validation", stringValue(fetchedCV.Properties.Description))
	require.Equal(t, "updated", stringValue(fetchedCV.Tags["phase"]))
	t.Log("CloudValidation create/get/list/update passed")

	configuration, err := json.Marshal(map[string]any{
		"apiVersion": "microsoft.validate/executionPlan.v0", "kind": "ExecutionPlan",
		"metadata": map[string]any{"name": owned.plan},
		"parameters": map[string]any{"certificationPackageReference": map[string]any{
			"osType": "Linux", "vmGenerationType": "V1", "architectureType": "X64",
			"recommendedVMSizes": []string{"Standard_D4s_v3"},
			"storageProfile": map[string]any{
				"osDiskImage": map[string]any{"sourceVhdUri": source}, "dataDiskImages": []any{},
			},
			"additionalProperties": map[string]any{},
		}},
		"authoring": map[string]any{"steps": []any{map[string]any{
			"name": "vm-boot-test", "type": "test", "testRef": bootTestID,
		}}},
	})
	requireNoError(t, err)
	_, err = plans.Get(ctx, owned.group, owned.name, owned.plan, nil)
	requireNotFound(t, err)
	plan := completeOperation(t, ctx, owned, func() (*runtime.Poller[pv.ValidationExecutionPlansClientCreateOrUpdateResponse], error) {
		return plans.BeginCreateOrUpdate(ctx, owned.group, owned.name, owned.plan, pv.ValidationExecutionPlan{
			Location: to.Ptr(clients.location), Tags: tags,
			Properties: &pv.ValidationExecutionPlanProperties{
				Description: to.Ptr("Go SDK boot plan"), PlanConfigurationJSON: to.Ptr(string(configuration)),
			},
		}, nil)
	})
	require.NotNil(t, plan.Properties)
	require.Equal(t, "Succeeded", stringValue(plan.Properties.ProvisioningState))
	fetchedPlan, err := plans.Get(ctx, owned.group, owned.name, owned.plan, nil)
	requireNoError(t, err)
	require.Equal(t, plan.ID, fetchedPlan.ID)
	planPager := plans.NewListByResourceGroupPager(owned.group, owned.name, nil)
	count := 0
	for planPager.More() {
		page, err := planPager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			require.Equal(t, plan.ID, entry.ID)
			count++
		}
	}
	require.Equal(t, 1, count)
	completeOperation(t, ctx, owned, func() (*runtime.Poller[pv.ValidationExecutionPlansClientUpdateResponse], error) {
		return plans.BeginUpdate(ctx, owned.group, owned.name, owned.plan, pv.ValidationExecutionPlanUpdate{
			Tags: map[string]*string{"sdk-test": to.Ptr(owned.name), "phase": to.Ptr("updated")},
		}, nil)
	})
	fetchedPlan, err = plans.Get(ctx, owned.group, owned.name, owned.plan, nil)
	requireNoError(t, err)
	require.NotNil(t, fetchedPlan.Properties)
	require.Equal(t, "updated", stringValue(fetchedPlan.Tags["phase"]))
	require.Equal(t, "Go SDK boot plan", stringValue(fetchedPlan.Properties.Description))
	var original, returned any
	requireNoError(t, json.Unmarshal(configuration, &original))
	requireNoError(t, json.Unmarshal([]byte(stringValue(fetchedPlan.Properties.PlanConfigurationJSON)), &returned))
	// Comparing as a boolean keeps signed URLs out of assertion diagnostics.
	require.True(t, reflect.DeepEqual(original, returned), "Plan PATCH must preserve the image configuration")
	t.Log("Execution plan create/get/list/update passed")

	_, err = runs.Get(ctx, owned.group, owned.name, owned.plan, owned.run, nil)
	requireNotFound(t, err)
	run := completeOperation(t, ctx, owned, func() (*runtime.Poller[pv.ExecutionPlanRunsClientCreateOrUpdateResponse], error) {
		return runs.BeginCreateOrUpdate(ctx, owned.group, owned.name, owned.plan, owned.run, pv.ExecutionPlanRun{
			Properties: &pv.ExecutionPlanRunProperties{Description: to.Ptr("Go SDK boot execution " + owned.name)},
		}, nil)
	})
	for {
		require.NotNil(t, run.Properties)
		if terminalStatus(stringValue(run.Properties.Status)) {
			break
		}
		requireNoError(t, ctx.Err())
		recording.Sleep(10 * time.Second)
		fetched, err := runs.Get(ctx, owned.group, owned.name, owned.plan, owned.run, nil)
		requireNoError(t, err)
		run.ExecutionPlanRun = fetched.ExecutionPlanRun
	}
	require.Contains(t, []string{"Completed", "Succeeded"}, stringValue(run.Properties.Status))
	require.Equal(t, "Succeeded", stringValue(run.Properties.ProvisioningState))
	require.Nil(t, run.Properties.Error)
	summary := run.Properties.TestRunSummary
	require.NotNil(t, summary)
	require.Equal(t, to.Ptr(int32(1)), summary.TotalTests)
	require.Equal(t, to.Ptr(int32(1)), summary.PassedTests)
	require.Equal(t, to.Ptr(int32(0)), summary.FailedTests)
	require.Equal(t, to.Ptr(int32(0)), summary.SkippedTests)
	require.Equal(t, "Passed", stringValue(summary.OverallResult))
	fetchedRun, err := runs.Get(ctx, owned.group, owned.name, owned.plan, owned.run, nil)
	requireNoError(t, err)
	require.Equal(t, run.ID, fetchedRun.ID)
	runPager := runs.NewListByExecutionPlanPager(owned.group, owned.name, owned.plan, nil)
	count = 0
	for runPager.More() {
		page, err := runPager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			require.Equal(t, run.ID, entry.ID)
			count++
		}
	}
	require.Equal(t, 1, count)
	results := clients.factory.NewValidationTestRunsClient()
	resultPager := results.NewListByExecutionPlanRunPager(owned.group, owned.name, owned.plan, owned.run, nil)
	count = 0
	for resultPager.More() {
		page, err := resultPager.NextPage(ctx)
		requireNoError(t, err)
		for _, entry := range page.Value {
			result, err := results.Get(ctx, owned.group, owned.name, owned.plan, owned.run, stringValue(entry.Name), nil)
			requireNoError(t, err)
			require.Equal(t, entry.ID, result.ID)
			require.NotNil(t, result.Properties)
			require.Equal(t, "Completed", stringValue(result.Properties.Status))
			require.Equal(t, "Succeeded", stringValue(result.Properties.ProvisioningState))
			require.Equal(t, bootTestID, stringValue(result.Properties.TestID))
			require.NotEmpty(t, result.Properties.PassDetails)
			require.Empty(t, result.Properties.FailureDetails)
			require.Nil(t, result.Properties.Error)
			count++
		}
	}
	require.Equal(t, 1, count)
	t.Log("Boot result: total=1; passed=1; failed=0; skipped=0")
}

func terminalStatus(status string) bool {
	switch strings.ToLower(status) {
	case "completed", "succeeded", "failed", "timedout", "canceled", "cancelled":
		return true
	default:
		return false
	}
}

func deleteAndPoll[T any](t *testing.T, ctx context.Context, begin func() (*runtime.Poller[T], error)) error {
	t.Helper()
	poller, err := begin()
	if err != nil {
		return err
	}
	_, err = testutil.PollForTest(ctx, poller)
	return err
}

func cleanupValidation(t *testing.T, clients *testClients, owned *ownedValidation) {
	t.Helper()
	if owned.mutationInProgress {
		t.Errorf("Retaining %s/%s: a mutation has an unknown outcome", owned.group, owned.name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	cvs := clients.factory.NewCloudValidationsClient()
	plans := clients.factory.NewValidationExecutionPlansClient()
	runs := clients.factory.NewExecutionPlanRunsClient()
	cv, err := cvs.Get(ctx, owned.group, owned.name, nil)
	if !isNotFound(err) {
		requireNoError(t, err)
		require.Equal(t, owned.name, stringValue(cv.Tags["sdk-test"]), "Refusing to delete an unowned CV")
		require.NotNil(t, cv.Properties)
		require.Contains(t, []string{"Succeeded", "Failed", "Canceled"}, stringValue(cv.Properties.ProvisioningState))
		planPager := plans.NewListByResourceGroupPager(owned.group, owned.name, nil)
		for planPager.More() {
			page, err := planPager.NextPage(ctx)
			requireNoError(t, err)
			for _, plan := range page.Value {
				require.Equal(t, owned.plan, stringValue(plan.Name))
				require.Equal(t, owned.name, stringValue(plan.Tags["sdk-test"]))
				require.NotNil(t, plan.Properties)
				require.Contains(t, []string{"Succeeded", "Failed", "Canceled"}, stringValue(plan.Properties.ProvisioningState))
				runPager := runs.NewListByExecutionPlanPager(owned.group, owned.name, owned.plan, nil)
				for runPager.More() {
					runPage, err := runPager.NextPage(ctx)
					requireNoError(t, err)
					for _, run := range runPage.Value {
						require.Equal(t, owned.run, stringValue(run.Name))
						require.NotNil(t, run.Properties)
						require.Equal(t, "Go SDK boot execution "+owned.name, stringValue(run.Properties.Description))
						require.True(t, terminalStatus(stringValue(run.Properties.Status)), "Retaining an active workload")
						require.Contains(t, []string{"Succeeded", "Failed", "Canceled"}, stringValue(run.Properties.ProvisioningState))
						requireNoError(t, deleteAndPoll(t, ctx, func() (*runtime.Poller[pv.ExecutionPlanRunsClientDeleteResponse], error) {
							return runs.BeginDelete(ctx, owned.group, owned.name, owned.plan, owned.run, nil)
						}))
						_, err = runs.Get(ctx, owned.group, owned.name, owned.plan, owned.run, nil)
						requireNotFound(t, err)
					}
				}
				remainingRuns := runs.NewListByExecutionPlanPager(owned.group, owned.name, owned.plan, nil)
				for remainingRuns.More() {
					page, err := remainingRuns.NextPage(ctx)
					requireNoError(t, err)
					require.Empty(t, page.Value, "Refusing parent deletion while executions remain")
				}
				requireNoError(t, deleteAndPoll(t, ctx, func() (*runtime.Poller[pv.ValidationExecutionPlansClientDeleteResponse], error) {
					return plans.BeginDelete(ctx, owned.group, owned.name, owned.plan, nil)
				}))
				_, err = plans.Get(ctx, owned.group, owned.name, owned.plan, nil)
				requireNotFound(t, err)
			}
		}
		for attempt := 0; ; attempt++ {
			remainingPlans := plans.NewListByResourceGroupPager(owned.group, owned.name, nil)
			for remainingPlans.More() {
				page, err := remainingPlans.NextPage(ctx)
				requireNoError(t, err)
				require.Empty(t, page.Value, "Refusing parent deletion while plans remain")
			}
			err := deleteAndPoll(t, ctx, func() (*runtime.Poller[pv.CloudValidationsClientDeleteResponse], error) {
				return cvs.BeginDelete(ctx, owned.group, owned.name, nil)
			})
			var response *azcore.ResponseError
			if errors.As(err, &response) && response.StatusCode == http.StatusConflict &&
				response.ErrorCode == "CannotDeleteResource" && attempt < 29 {
				requireNoError(t, ctx.Err())
				recording.Sleep(10 * time.Second)
				continue
			}
			requireNoError(t, err)
			break
		}
		_, err = cvs.Get(ctx, owned.group, owned.name, nil)
		requireNotFound(t, err)
	}
	groups, err := armresources.NewResourceGroupsClient(clients.subscriptionID, clients.credential, clients.options)
	requireNoError(t, err)
	for attempt := 0; attempt < 60; attempt++ {
		_, err = groups.Get(ctx, owned.name+"-mrg", nil)
		if isNotFound(err) {
			break
		}
		requireNoError(t, err)
		recording.Sleep(10 * time.Second)
	}
	requireNotFound(t, err)
	group, err := groups.Get(ctx, owned.group, nil)
	requireNoError(t, err)
	require.Equal(t, owned.name, stringValue(group.Tags["sdk-test"]))
	resources, err := armresources.NewClient(clients.subscriptionID, clients.credential, clients.options)
	requireNoError(t, err)
	empty := false
	for attempt := 0; attempt < 30; attempt++ {
		empty = true
		pager := resources.NewListByResourceGroupPager(owned.group, nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			requireNoError(t, err)
			empty = empty && len(page.Value) == 0
		}
		if empty {
			break
		}
		recording.Sleep(10 * time.Second)
	}
	require.True(t, empty, "Retaining a nonempty resource group")
	requireNoError(t, deleteAndPoll(t, ctx, func() (*runtime.Poller[armresources.ResourceGroupsClientDeleteResponse], error) {
		return groups.BeginDelete(ctx, owned.group, nil)
	}))
	_, err = groups.Get(ctx, owned.group, nil)
	requireNotFound(t, err)
	t.Log("Cleanup verified: execution, plan, CV, managed group and owned test group are absent")
}

func isNotFound(err error) bool {
	var response *azcore.ResponseError
	return errors.As(err, &response) && response.StatusCode == http.StatusNotFound
}
