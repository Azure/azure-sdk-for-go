//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armfeatures_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armfeatures"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFeaturesClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	featuresClient := armfeatures.NewClient(subscriptionID, cred, opt)

	pager := featuresClient.List("Microsoft.Network", nil)
	require.NoError(t, pager.Err())
}

func TestFeaturesClient_ListAll(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	featuresClient := armfeatures.NewClient(subscriptionID, cred, opt)

	pager := featuresClient.ListAll(nil)
	require.NoError(t, pager.Err())
}

func TestFeaturesClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	featuresClient := armfeatures.NewClient(subscriptionID, cred, opt)
	featureName, _ := createRandomName(t, "feature")
	_, err := featuresClient.Get(ctx, "Microsoft.Network", featureName, nil)
	require.Error(t, err)
}

func TestFeatureClient_ListOperations(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)

	featureClient := armfeatures.NewFeatureClient(cred, opt)
	pager := featureClient.ListOperations(nil)
	require.NoError(t, pager.Err())
}

//func TestFeaturesClient_Register(t *testing.T) {
//	stop := startTest(t)
//	defer stop()
//
//	cred, opt := authenticateTest(t)
//	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
//	ctx := context.Background()
//
//	featuresClient := armfeatures.NewFeaturesClient(subscriptionID, cred, opt)
//	featureName, _ := createRandomName(t, "feature")
//	register, err := featuresClient.Register(ctx, "Microsoft.Network", featureName, nil)
//	require.NoError(t, err)
//	require.Equal(t, featureName, register.Name)
//}
//
//func TestFeaturesClient_Unregister(t *testing.T) {
//	stop := startTest(t)
//	defer stop()
//
//	cred, opt := authenticateTest(t)
//	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
//	ctx := context.Background()
//
//	featuresClient := armfeatures.NewFeaturesClient(subscriptionID, cred, opt)
//	featureName, _ := createRandomName(t, "feature")
//	register, err := featuresClient.Unregister(ctx, "Microsoft.Network", featureName, nil)
//	require.NoError(t, err)
//	require.Equal(t, featureName, register.Name)
//}
