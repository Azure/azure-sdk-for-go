// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// newTestFeatureFlagClient constructs an [*azappconfig.FeatureFlagClient] authenticated with Entra
// ID and wired to the test recording transport. Behaves like [newTestClient] but returns the
// feature flag client.
func newTestFeatureFlagClient(t *testing.T) *azappconfig.FeatureFlagClient {
	if recording.GetRecordMode() != recording.PlaybackMode && os.Getenv("APPCONFIGURATION_ENDPOINT_STRING") == "" {
		t.Skip("set APPCONFIGURATION_ENDPOINT_STRING to run this test")
	}

	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	err = recording.SetDefaultMatcher(t, &recording.SetDefaultMatcherOptions{
		IgnoredQueryParameters: []string{"api-version"},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	client, err := azappconfig.NewFeatureFlagClient(endpoint, credential, &azappconfig.FeatureFlagClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	return client
}

func TestFeatureFlagClient(t *testing.T) {
	const (
		name  = "flag-TestFeatureFlagClient"
		label = "flag-label"
	)

	client := newTestFeatureFlagClient(t)

	// Clean up any leftover flag from a previous test run
	_, _ = client.DeleteFeatureFlag(context.Background(), name, &azappconfig.DeleteFeatureFlagOptions{
		Label: to.Ptr(label),
	})

	// AddFeatureFlag creates the flag
	addResp, err := client.AddFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Name:        to.Ptr(name),
		Label:       to.Ptr(label),
		Enabled:     to.Ptr(true),
		Description: to.Ptr("initial"),
	}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, addResp)
	require.NotNil(t, addResp.Name)
	require.NotNil(t, addResp.Label)
	require.NotNil(t, addResp.Enabled)
	require.NotNil(t, addResp.Description)
	require.Equal(t, name, *addResp.Name)
	require.Equal(t, label, *addResp.Label)
	require.True(t, *addResp.Enabled)
	require.Equal(t, "initial", *addResp.Description)

	// Adding again should conflict
	addResp2, err := client.AddFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Name:    to.Ptr(name),
		Label:   to.Ptr(label),
		Enabled: to.Ptr(true),
	}, nil)
	require.Error(t, err)
	require.Empty(t, addResp2)

	// GetFeatureFlag retrieves the created flag
	getResp, err := client.GetFeatureFlag(context.Background(), name, &azappconfig.GetFeatureFlagOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
	require.NotEmpty(t, getResp)
	require.Equal(t, name, *getResp.Name)
	require.Equal(t, label, *getResp.Label)
	require.True(t, *getResp.Enabled)
	require.NotNil(t, getResp.ETag)

	// GetFeatureFlag with OnlyIfChanged and the current ETag should return 304 -> error
	etag := getResp.ETag
	getResp2, err := client.GetFeatureFlag(context.Background(), name, &azappconfig.GetFeatureFlagOptions{
		Label:         to.Ptr(label),
		OnlyIfChanged: etag,
	})
	require.Error(t, err)
	require.Empty(t, getResp2)

	// SetFeatureFlag overwrites the flag
	setResp, err := client.SetFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Name:        to.Ptr(name),
		Label:       to.Ptr(label),
		Enabled:     to.Ptr(false),
		Description: to.Ptr("updated"),
	}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, setResp)
	require.False(t, *setResp.Enabled)
	require.Equal(t, "updated", *setResp.Description)
	require.NotEqual(t, "", string(*setResp.ETag))

	// GetFeatureFlag with OnlyIfChanged and the stale ETag should now succeed
	getResp3, err := client.GetFeatureFlag(context.Background(), name, &azappconfig.GetFeatureFlagOptions{
		Label:         to.Ptr(label),
		OnlyIfChanged: etag,
	})
	require.NoError(t, err)
	require.NotEmpty(t, getResp3)
	require.False(t, *getResp3.Enabled)
	require.Equal(t, "updated", *getResp3.Description)

	// SetFeatureFlag with OnlyIfUnchanged and a stale ETag should fail
	setResp2, err := client.SetFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Name:    to.Ptr(name),
		Label:   to.Ptr(label),
		Enabled: to.Ptr(true),
	}, &azappconfig.SetFeatureFlagOptions{
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err)
	require.Empty(t, setResp2)

	// SetFeatureFlag with OnlyIfUnchanged and the current ETag should succeed
	currentETag := getResp3.ETag
	setResp3, err := client.SetFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Name:    to.Ptr(name),
		Label:   to.Ptr(label),
		Enabled: to.Ptr(true),
	}, &azappconfig.SetFeatureFlagOptions{
		OnlyIfUnchanged: currentETag,
	})
	require.NoError(t, err)
	require.NotEmpty(t, setResp3)
	require.True(t, *setResp3.Enabled)

	// NewListFeatureFlagsPager returns the flag
	listPager := client.NewListFeatureFlagsPager(azappconfig.FeatureFlagSelector{
		NameFilter:  to.Ptr(name),
		LabelFilter: to.Ptr(label),
		Fields:      azappconfig.AllFeatureFlagFields(),
	}, nil)
	require.NotNil(t, listPager)
	require.True(t, listPager.More())

	listResp, err := listPager.NextPage(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, listResp.FeatureFlags)
	require.Equal(t, name, *listResp.FeatureFlags[0].Name)
	require.Equal(t, label, *listResp.FeatureFlags[0].Label)

	// NewListFeatureFlagRevisionsPager returns at least one revision
	revPager := client.NewListFeatureFlagRevisionsPager(azappconfig.FeatureFlagSelector{
		NameFilter:  to.Ptr(name),
		LabelFilter: to.Ptr(label),
		Fields:      azappconfig.AllFeatureFlagFields(),
	}, nil)
	require.NotNil(t, revPager)
	require.True(t, revPager.More())

	revResp, err := revPager.NextPage(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, revResp.FeatureFlags)
	require.Equal(t, name, *revResp.FeatureFlags[0].Name)

	// DeleteFeatureFlag with a stale ETag should fail
	delResp, err := client.DeleteFeatureFlag(context.Background(), name, &azappconfig.DeleteFeatureFlagOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err)
	require.Empty(t, delResp)

	// DeleteFeatureFlag with the current ETag should succeed
	currentETag = setResp3.ETag
	delResp2, err := client.DeleteFeatureFlag(context.Background(), name, &azappconfig.DeleteFeatureFlagOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: currentETag,
	})
	require.NoError(t, err)
	require.NotEmpty(t, delResp2)
	require.Equal(t, name, *delResp2.Name)
}

func TestFeatureFlagClient_GetNotFound(t *testing.T) {
	client := newTestFeatureFlagClient(t)

	resp, err := client.GetFeatureFlag(context.Background(), "flag-TestFeatureFlagClient_GetNotFound-nonexistent", nil)
	require.Error(t, err)
	require.Empty(t, resp)

	var respErr *azcore.ResponseError
	require.True(t, errors.As(err, &respErr), "expected an *azcore.ResponseError, got: %v", err)
	require.Equal(t, 404, respErr.StatusCode)
}

func TestFeatureFlagClient_NameRequired(t *testing.T) {
	// Purely local validation — no HTTP is performed when flag.Name is nil.
	client, err := azappconfig.NewFeatureFlagClient("https://fake.azconfig.io", credential, nil)
	require.NoError(t, err)

	_, err = client.AddFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Enabled: to.Ptr(true),
	}, nil)
	require.Error(t, err)

	_, err = client.SetFeatureFlag(context.Background(), azappconfig.FeatureFlag{
		Enabled: to.Ptr(true),
	}, nil)
	require.Error(t, err)
}

func TestFeatureFlagClient_FromClient(t *testing.T) {
	client := newTestClient(t)

	ffClient := client.NewFeatureFlagClient()
	require.NotNil(t, ffClient)

	// The FF client shares the pipeline of the source Client, so getting an unknown
	// flag should surface as a normal 404.
	_, err := ffClient.GetFeatureFlag(context.Background(), "flag-TestFeatureFlagClient_FromClient-nonexistent", nil)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.True(t, errors.As(err, &respErr), "expected an *azcore.ResponseError, got: %v", err)
	require.Equal(t, 404, respErr.StatusCode)
}

func TestFeatureFlagClient_ComplexModel(t *testing.T) {
	const (
		name  = "flag-TestFeatureFlagClient_ComplexModel"
		label = "flag-label"
	)

	client := newTestFeatureFlagClient(t)

	// Clean up any leftover flag from a previous test run
	_, _ = client.DeleteFeatureFlag(context.Background(), name, &azappconfig.DeleteFeatureFlagOptions{
		Label: to.Ptr(label),
	})

	flag := azappconfig.FeatureFlag{
		Name:    to.Ptr(name),
		Label:   to.Ptr(label),
		Enabled: to.Ptr(true),
		Conditions: &azappconfig.FeatureFlagConditions{
			RequirementType: to.Ptr(azappconfig.RequirementTypeAll),
			Filters: []azappconfig.FeatureFlagFilter{
				{
					Name: to.Ptr("Microsoft.Percentage"),
					Parameters: map[string]*string{
						"Value": to.Ptr("50"),
					},
				},
			},
		},
		Variants: []azappconfig.FeatureFlagVariantDefinition{
			{
				Name:           to.Ptr("On"),
				Value:          to.Ptr("true"),
				StatusOverride: to.Ptr(azappconfig.StatusOverrideEnabled),
			},
			{
				Name:           to.Ptr("Off"),
				Value:          to.Ptr("false"),
				StatusOverride: to.Ptr(azappconfig.StatusOverrideDisabled),
			},
		},
		Allocation: &azappconfig.FeatureFlagAllocation{
			DefaultWhenEnabled:  to.Ptr("On"),
			DefaultWhenDisabled: to.Ptr("Off"),
			Percentile: []azappconfig.PercentileAllocation{
				{
					Variant: to.Ptr("On"),
					From:    to.Ptr(0.0),
					To:      to.Ptr(50.0),
				},
				{
					Variant: to.Ptr("Off"),
					From:    to.Ptr(50.0),
					To:      to.Ptr(100.0),
				},
			},
		},
		Telemetry: &azappconfig.FeatureFlagTelemetryConfiguration{
			Enabled: to.Ptr(true),
			Metadata: map[string]*string{
				"source": to.Ptr("azappconfig-tests"),
			},
		},
		Tags: map[string]*string{
			"environment": to.Ptr("test"),
		},
	}

	setResp, err := client.SetFeatureFlag(context.Background(), flag, nil)
	require.NoError(t, err)
	require.NotEmpty(t, setResp)

	getResp, err := client.GetFeatureFlag(context.Background(), name, &azappconfig.GetFeatureFlagOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
	require.NotNil(t, getResp.Conditions)
	require.NotNil(t, getResp.Conditions.RequirementType)
	require.Equal(t, azappconfig.RequirementTypeAll, *getResp.Conditions.RequirementType)
	require.Len(t, getResp.Conditions.Filters, 1)
	require.Equal(t, "Microsoft.Percentage", *getResp.Conditions.Filters[0].Name)
	require.Len(t, getResp.Variants, 2)
	require.NotNil(t, getResp.Allocation)
	require.Len(t, getResp.Allocation.Percentile, 2)
	require.NotNil(t, getResp.Telemetry)
	require.True(t, *getResp.Telemetry.Enabled)
	require.Equal(t, flag.Tags, getResp.Tags)

	// Clean up
	_, err = client.DeleteFeatureFlag(context.Background(), name, &azappconfig.DeleteFeatureFlagOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
}
