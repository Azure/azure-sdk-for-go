//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package policy

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientOptions_Copy(t *testing.T) {
	var option *ClientOptions
	require.Nil(t, option.Clone())

	option = &ClientOptions{ClientOptions: policy.ClientOptions{
		Cloud: cloud.AzurePublic,
		Logging: policy.LogOptions{
			AllowedHeaders:     []string{"test1", "test2"},
			AllowedQueryParams: []string{"test1", "test2"},
		},
		Retry:            policy.RetryOptions{StatusCodes: []int{1, 2}},
		PerRetryPolicies: []policy.Policy{runtime.NewLogPolicy(nil)},
		PerCallPolicies:  []policy.Policy{runtime.NewLogPolicy(nil)},
	}}
	copiedOption := option.Clone()
	require.Equal(t, option.APIVersion, copiedOption.APIVersion)
	require.NotEqual(t, fmt.Sprintf("%p", &option.APIVersion), fmt.Sprintf("%p", &copiedOption.APIVersion))
	require.Equal(t, option.Cloud.Services, copiedOption.Cloud.Services)
	require.NotEqual(t, fmt.Sprintf("%p", option.Cloud.Services), fmt.Sprintf("%p", copiedOption.Cloud.Services))
	require.Equal(t, option.Logging.AllowedHeaders, copiedOption.Logging.AllowedHeaders)
	require.NotEqual(t, fmt.Sprintf("%p", option.Logging.AllowedHeaders), fmt.Sprintf("%p", copiedOption.Logging.AllowedHeaders))
	require.Equal(t, option.Logging.AllowedQueryParams, copiedOption.Logging.AllowedQueryParams)
	require.NotEqual(t, fmt.Sprintf("%p", option.Logging.AllowedQueryParams), fmt.Sprintf("%p", copiedOption.Logging.AllowedQueryParams))
	require.Equal(t, option.Retry.StatusCodes, copiedOption.Retry.StatusCodes)
	require.NotEqual(t, fmt.Sprintf("%p", option.Retry.StatusCodes), fmt.Sprintf("%p", copiedOption.Retry.StatusCodes))
	require.Equal(t, option.PerRetryPolicies, copiedOption.PerRetryPolicies)
	require.NotEqual(t, fmt.Sprintf("%p", option.PerRetryPolicies), fmt.Sprintf("%p", copiedOption.PerRetryPolicies))
	require.Equal(t, option.PerCallPolicies, copiedOption.PerCallPolicies)
	require.NotEqual(t, fmt.Sprintf("%p", option.PerCallPolicies), fmt.Sprintf("%p", copiedOption.PerCallPolicies))
}
