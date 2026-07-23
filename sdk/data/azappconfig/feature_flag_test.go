// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
	"github.com/stretchr/testify/require"
)

func TestFeatureFlagFromGenerated(t *testing.T) {
	lastModified := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	gen := generated.FeatureFlag{
		Enabled:      to.Ptr(true),
		Name:         to.Ptr("flag"),
		Label:        to.Ptr("label"),
		Description:  to.Ptr("desc"),
		Etag:         to.Ptr("etag"),
		Tags:         map[string]*string{"env": to.Ptr("test")},
		Allocation:   &generated.FeatureFlagAllocation{DefaultWhenEnabled: to.Ptr("A")},
		Conditions:   &generated.FeatureFlagConditions{RequirementType: to.Ptr(generated.RequirementTypeAll)},
		Telemetry:    &generated.FeatureFlagTelemetryConfiguration{Enabled: to.Ptr(true)},
		Variants:     []generated.FeatureFlagVariantDefinition{{Name: to.Ptr("A"), Value: to.Ptr("true")}},
		LastModified: &lastModified,
	}

	pub := featureFlagFromGenerated(gen)

	require.Equal(t, gen.Enabled, pub.Enabled)
	require.Equal(t, gen.Name, pub.Name)
	require.Equal(t, gen.Label, pub.Label)
	require.Equal(t, gen.Description, pub.Description)
	require.Equal(t, gen.Tags, pub.Tags)
	require.Equal(t, gen.Allocation, pub.Allocation)
	require.Equal(t, gen.Conditions, pub.Conditions)
	require.Equal(t, gen.Telemetry, pub.Telemetry)
	require.Equal(t, gen.Variants, pub.Variants)
	require.Equal(t, gen.LastModified, pub.LastModified)
	require.NotNil(t, pub.ETag)
	require.Equal(t, azcore.ETag("etag"), *pub.ETag)
}

func TestFeatureFlagToGenerated(t *testing.T) {
	etag := azcore.ETag("etag")
	lastModified := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	pub := FeatureFlag{
		Enabled:      to.Ptr(true),
		Name:         to.Ptr("flag"),
		Label:        to.Ptr("label"),
		Description:  to.Ptr("desc"),
		ETag:         &etag,
		Tags:         map[string]*string{"env": to.Ptr("test")},
		Allocation:   &FeatureFlagAllocation{DefaultWhenEnabled: to.Ptr("A")},
		Conditions:   &FeatureFlagConditions{RequirementType: to.Ptr(RequirementTypeAll)},
		Telemetry:    &FeatureFlagTelemetryConfiguration{Enabled: to.Ptr(true)},
		Variants:     []FeatureFlagVariantDefinition{{Name: to.Ptr("A"), Value: to.Ptr("true")}},
		LastModified: &lastModified,
	}

	gen := pub.toGenerated()

	require.Equal(t, pub.Enabled, gen.Enabled)
	require.Equal(t, pub.Name, gen.Name)
	require.Equal(t, pub.Label, gen.Label)
	require.Equal(t, pub.Description, gen.Description)
	require.Equal(t, pub.Tags, gen.Tags)
	require.Equal(t, pub.Allocation, gen.Allocation)
	require.Equal(t, pub.Conditions, gen.Conditions)
	require.Equal(t, pub.Telemetry, gen.Telemetry)
	require.Equal(t, pub.Variants, gen.Variants)
	require.Equal(t, pub.LastModified, gen.LastModified)
	require.NotNil(t, gen.Etag)
	require.Equal(t, "etag", *gen.Etag)
}

func TestFeatureFlagRoundTrip(t *testing.T) {
	etag := azcore.ETag("abc")
	pub := FeatureFlag{
		Enabled:     to.Ptr(true),
		Name:        to.Ptr("flag"),
		Label:       to.Ptr("label"),
		Description: to.Ptr("desc"),
		ETag:        &etag,
	}

	got := featureFlagFromGenerated(pub.toGenerated())
	require.Equal(t, pub, got)
}

func TestAllFeatureFlagFields(t *testing.T) {
	fields := AllFeatureFlagFields()
	require.ElementsMatch(t, generated.PossibleFeatureFlagFieldsValues(), fields)
}
