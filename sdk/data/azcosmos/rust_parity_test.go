// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncodingRustDefaultsAndEnvironment(t *testing.T) {
	for _, value := range []string{"1", "true", "YES", " on ", "false", "", "nope"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", value)
			want := value == "1" || value == "true" || value == "YES" || value == " on "
			require.Equal(t, want, resolveClientBinaryEncoding(nil).enabled())
			require.True(t, resolveClientBinaryEncoding(&BinaryEncodingOptions{}).enabled())
			require.False(t, resolveClientBinaryEncoding(&BinaryEncodingOptions{Enabled: to(false)}).enabled())
		})
	}
}

func TestClientCopiesDedicatedBinaryOptions(t *testing.T) {
	options := &BinaryEncodingOptions{Enabled: to(false), RequestTextResponse: true}
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", mustKeyCredential(t), &ClientOptions{BinaryEncoding: options})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	*options.Enabled = true
	options.RequestTextResponse = false
	require.False(t, client.binaryEncoding.enabled())
	require.True(t, client.binaryEncoding.RequestTextResponse)
}

func TestPatchTrackingRustRetentionAndUUID(t *testing.T) {
	for _, id := range []PatchTrackingID{
		"00112233-4455-6677-8899-AABBCCDDEEFF",
		"00112233445566778899aabbccddeeff",
		"{00112233-4455-6677-8899-aabbccddeeff}",
		"urn:uuid:00112233-4455-6677-8899-aabbccddeeff",
	} {
		normalized, err := id.normalized()
		require.NoError(t, err)
		require.Equal(t, PatchTrackingID("00112233-4455-6677-8899-aabbccddeeff"), normalized)
	}
	require.NoError(t, (PatchItemOptions{TrackingRetention: to(time.Duration(math.MaxInt64))}).validateTracking())
}

func TestFaultInjectionValidationAndCopies(t *testing.T) {
	rules := []FaultInjectionRule{{
		ID: "fault", Condition: FaultInjectionCondition{Operation: FaultInjectionReadItem},
		Result: FaultInjectionResult{CustomStatusCode: 429, Probability: to(float32(1)),
			Headers: map[string]string{"x-test": "one"}, Body: []byte("body"), RetryAfter: to(time.Second)},
		HitLimit: to(uint32(0)), ExpireAfter: to(time.Second),
	}}
	require.NoError(t, validateFaultInjectionRules(rules))
	copy := cloneFaultInjectionRules(rules)
	rules[0].Result.Headers["x-test"] = "two"
	rules[0].Result.Body[0] = 'x'
	*rules[0].HitLimit = 9
	*rules[0].Result.Probability = 0
	require.Equal(t, "one", copy[0].Result.Headers["x-test"])
	require.Equal(t, []byte("body"), copy[0].Result.Body)
	require.Zero(t, *copy[0].HitLimit)
	require.Equal(t, float32(1), *copy[0].Result.Probability)
	for _, mutate := range []func(*FaultInjectionRule){
		func(r *FaultInjectionRule) { r.ID = "" },
		func(r *FaultInjectionRule) { r.Condition.Operation = 99 },
		func(r *FaultInjectionRule) { r.Result.Error = FaultInjectionTimeout },
		func(r *FaultInjectionRule) { r.Result.Probability = to(float32(math.NaN())) },
		func(r *FaultInjectionRule) { r.Result.RetryAfter = to(-time.Second) },
		func(r *FaultInjectionRule) { r.Result.Headers = map[string]string{"bad name": "value"} },
	} {
		invalid := cloneFaultInjectionRules(copy)
		mutate(&invalid[0])
		require.Error(t, validateFaultInjectionRules(invalid))
	}
	require.Error(t, validateFaultInjectionRules(append(copy, copy[0])))
}
