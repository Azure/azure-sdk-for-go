// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcognitiveservices_test

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/require"
)

const (
	recordedSubscriptionID = "00000000-0000-0000-0000-000000000000"
	recordedResourceGroup  = "sanitized-resource-group"
	recordedAccountName    = "sanitized-cognitive-account"
	lifecycleControlName   = "sdk-cost-control-lifecycle"
	sanitizedETag          = "sanitized-etag"
)

func TestCostControlLifecycle(t *testing.T) {
	stopRecording := testutil.StartRecording(t, pathToPackage)
	defer stopRecording()

	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", recordedSubscriptionID)
	resourceGroupName := recording.GetEnvVariable("AZURE_RESOURCE_GROUP", recordedResourceGroup)
	accountName := recording.GetEnvVariable("AZURE_COGNITIVE_SERVICES_ACCOUNT", recordedAccountName)
	addCostControlSanitizers(t, subscriptionID, resourceGroupName, accountName)

	credential, options := testutil.GetCredAndClientOptions(t)
	client, err := armcognitiveservices.NewCostControlsClient(subscriptionID, credential, options)
	require.NoError(t, err)

	ctx := context.Background()
	created := false
	t.Cleanup(func() {
		if !created {
			return
		}
		_, cleanupErr := client.Delete(ctx, resourceGroupName, accountName, lifecycleControlName, nil)
		if cleanupErr != nil && !isResponseStatus(cleanupErr, http.StatusNotFound) {
			t.Logf("cleanup failed: %v", cleanupErr)
		}
	})

	createdResponse, err := client.CreateOrUpdate(ctx, resourceGroupName, accountName, lifecycleControlName,
		newCostControl("SDK Cost Control lifecycle", 25),
		&armcognitiveservices.CostControlsClientCreateOrUpdateOptions{IfNoneMatch: ptr("*")})
	require.NoError(t, err)
	created = true
	require.Equal(t, lifecycleControlName, *createdResponse.Name)
	require.Equal(t, armcognitiveservices.CostControlThresholdActionAudit,
		*createdResponse.Properties.Rules[0].Thresholds[0].Action)

	fetched, err := client.Get(ctx, resourceGroupName, accountName, lifecycleControlName, nil)
	require.NoError(t, err)
	require.NotNil(t, fetched.Etag)
	require.Equal(t, float64(25), *fetched.Properties.Rules[0].Amount)

	found := false
	pager := client.NewListPager(resourceGroupName, accountName, nil)
	for pager.More() {
		page, pageErr := pager.NextPage(ctx)
		require.NoError(t, pageErr)
		for _, item := range page.Value {
			if item.Name != nil && *item.Name == lifecycleControlName {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	require.True(t, found)

	updated, err := client.Update(ctx, resourceGroupName, accountName, lifecycleControlName,
		armcognitiveservices.CostControlPatch{Properties: &armcognitiveservices.CostControlPatchProperties{
			DisplayName: ptr("SDK Cost Control lifecycle updated"),
			Rules:       []*armcognitiveservices.CostControlRule{newRule(50)},
		}}, &armcognitiveservices.CostControlsClientUpdateOptions{IfMatch: fetched.Etag})
	require.NoError(t, err)
	require.Equal(t, "SDK Cost Control lifecycle updated", *updated.Properties.DisplayName)
	require.Equal(t, float64(50), *updated.Properties.Rules[0].Amount)

	fetchedAfterUpdate, err := client.Get(ctx, resourceGroupName, accountName, lifecycleControlName, nil)
	require.NoError(t, err)
	require.NotNil(t, fetchedAfterUpdate.Etag)
	require.Equal(t, float64(50), *fetchedAfterUpdate.Properties.Rules[0].Amount)

	_, err = client.Delete(ctx, resourceGroupName, accountName, lifecycleControlName,
		&armcognitiveservices.CostControlsClientDeleteOptions{IfMatch: fetchedAfterUpdate.Etag})
	require.NoError(t, err)
	created = false

	_, err = client.Get(ctx, resourceGroupName, accountName, lifecycleControlName, nil)
	require.True(t, isResponseStatus(err, http.StatusNotFound), "expected 404, got %v", err)
}

func addCostControlSanitizers(t *testing.T, subscriptionID, resourceGroupName, accountName string) {
	options := &recording.RecordingOptions{UseHTTPS: true, TestInstance: t}
	require.NoError(t, recording.RemoveRegisteredSanitizers([]string{"AZSDK3493", "AZSDK3430"}, options))
	require.NoError(t, recording.AddGeneralRegexSanitizer(recordedSubscriptionID,
		regexp.QuoteMeta(subscriptionID), options))
	require.NoError(t, recording.AddGeneralRegexSanitizer(recordedResourceGroup,
		regexp.QuoteMeta(resourceGroupName), options))
	require.NoError(t, recording.AddGeneralRegexSanitizer(recordedAccountName,
		regexp.QuoteMeta(accountName), options))
	require.NoError(t, recording.AddBodyKeySanitizer("$..etag", sanitizedETag, "", options))
	require.NoError(t, recording.AddHeaderRegexSanitizer("etag", sanitizedETag, ".+", options))
	require.NoError(t, recording.AddHeaderRegexSanitizer("if-match", sanitizedETag, ".+", options))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-correlation-request-id",
		"sanitized-correlation-request-id", ".+", options))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-operation-identifier",
		"sanitized-operation-identifier", ".+", options))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-routing-request-id",
		"sanitized-routing-request-id", ".+", options))
}

func isResponseStatus(err error, status int) bool {
	var responseError *azcore.ResponseError
	return errors.As(err, &responseError) && responseError.StatusCode == status
}
