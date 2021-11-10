//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManagementPoliciesClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// create management policy
	managementPoliciesClient := armstorage.NewManagementPoliciesClient(subscriptionID,cred,opt)
	mpResp,err := managementPoliciesClient.CreateOrUpdate(
		ctx,
		rgName,scName,
		armstorage.ManagementPolicyNameDefault,
		armstorage.ManagementPolicy{
			Properties: &armstorage.ManagementPolicyProperties{
				Policy: &armstorage.ManagementPolicySchema{
					Rules: []*armstorage.ManagementPolicyRule{
						{
							Enabled: to.BoolPtr(true),
							Name: to.StringPtr("olcmtest"),
							Type: armstorage.RuleTypeLifecycle.ToPtr(),
							Definition: &armstorage.ManagementPolicyDefinition{
								Filters: &armstorage.ManagementPolicyFilter{
									BlobTypes: []*string{
										to.StringPtr("blockBlob"),
									},
									PrefixMatch: []*string{
										to.StringPtr("olcmtestcontainer"),
									},
								},
								Actions: &armstorage.ManagementPolicyAction{
									BaseBlob: &armstorage.ManagementPolicyBaseBlob{
										TierToCool: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(30),
										},
										TierToArchive: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(90),
										},
										Delete: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(1000),
										},
									},
									Snapshot: &armstorage.ManagementPolicySnapShot{
										Delete: &armstorage.DateAfterCreation{
											DaysAfterCreationGreaterThan: to.Float32Ptr(30),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, "olcmtest",mpResp.Name)
}

func TestManagementPoliciesClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// create management policy
	managementPoliciesClient := armstorage.NewManagementPoliciesClient(subscriptionID,cred,opt)
	mpResp,err := managementPoliciesClient.CreateOrUpdate(
		ctx,
		rgName,scName,
		armstorage.ManagementPolicyNameDefault,
		armstorage.ManagementPolicy{
			Properties: &armstorage.ManagementPolicyProperties{
				Policy: &armstorage.ManagementPolicySchema{
					Rules: []*armstorage.ManagementPolicyRule{
						{
							Enabled: to.BoolPtr(true),
							Name: to.StringPtr("olcmtest"),
							Type: armstorage.RuleTypeLifecycle.ToPtr(),
							Definition: &armstorage.ManagementPolicyDefinition{
								Filters: &armstorage.ManagementPolicyFilter{
									BlobTypes: []*string{
										to.StringPtr("blockBlob"),
									},
									PrefixMatch: []*string{
										to.StringPtr("olcmtestcontainer"),
									},
								},
								Actions: &armstorage.ManagementPolicyAction{
									BaseBlob: &armstorage.ManagementPolicyBaseBlob{
										TierToCool: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(30),
										},
										TierToArchive: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(90),
										},
										Delete: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(1000),
										},
									},
									Snapshot: &armstorage.ManagementPolicySnapShot{
										Delete: &armstorage.DateAfterCreation{
											DaysAfterCreationGreaterThan: to.Float32Ptr(30),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "olcmtest",mpResp.Name)

	// get management policy
	getResp,err := managementPoliciesClient.Get(ctx,rgName,scName,"default",nil)
	require.NoError(t, err)
	require.Equal(t, "default",*getResp.Name)
}

func TestManagementPoliciesClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// create management policy
	managementPoliciesClient := armstorage.NewManagementPoliciesClient(subscriptionID,cred,opt)
	mpResp,err := managementPoliciesClient.CreateOrUpdate(
		ctx,
		rgName,scName,
		armstorage.ManagementPolicyNameDefault,
		armstorage.ManagementPolicy{
			Properties: &armstorage.ManagementPolicyProperties{
				Policy: &armstorage.ManagementPolicySchema{
					Rules: []*armstorage.ManagementPolicyRule{
						{
							Enabled: to.BoolPtr(true),
							Name: to.StringPtr("olcmtest"),
							Type: armstorage.RuleTypeLifecycle.ToPtr(),
							Definition: &armstorage.ManagementPolicyDefinition{
								Filters: &armstorage.ManagementPolicyFilter{
									BlobTypes: []*string{
										to.StringPtr("blockBlob"),
									},
									PrefixMatch: []*string{
										to.StringPtr("olcmtestcontainer"),
									},
								},
								Actions: &armstorage.ManagementPolicyAction{
									BaseBlob: &armstorage.ManagementPolicyBaseBlob{
										TierToCool: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(30),
										},
										TierToArchive: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(90),
										},
										Delete: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(1000),
										},
									},
									Snapshot: &armstorage.ManagementPolicySnapShot{
										Delete: &armstorage.DateAfterCreation{
											DaysAfterCreationGreaterThan: to.Float32Ptr(30),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "olcmtest",mpResp.Name)

	// delete management policy
	delResp,err := managementPoliciesClient.Delete(ctx,rgName,scName,"default",nil)
	require.NoError(t, err)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}