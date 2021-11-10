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

func TestFileServicesClient_SetServiceProperties(t *testing.T) {
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

	// put file services
	fileServicesClient := armstorage.NewFileServicesClient(subscriptionID,cred,opt)
	putResp,err := fileServicesClient.SetServiceProperties(
		ctx,
		rgName,
		scName,
		armstorage.FileServiceProperties{
			FileServiceProperties: &armstorage.FileServicePropertiesProperties{
				Cors: &armstorage.CorsRules{
					CorsRules: []*armstorage.CorsRule{
						{
							AllowedOrigins: []*string{
								to.StringPtr("http://www.contoso.com"),
								to.StringPtr("http://www.fabrikam.com"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemHEAD.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemMERGE.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemOPTIONS.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPOST.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPUT.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(100),
							ExposedHeaders: []*string{
								to.StringPtr("x-ms-meta-*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("x-ms-meta-abc"),
								to.StringPtr("x-ms-meta-data*"),
								to.StringPtr("x-ms-meta-target*"),
							},
						},
						{
							AllowedOrigins: []*string{
								to.StringPtr("*"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(2),
							ExposedHeaders: []*string{
								to.StringPtr("*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("*"),
							},
						},
						{
							AllowedOrigins: []*string{
								to.StringPtr("http://www.abc23.com"),
								to.StringPtr("https://www.fabrikam.com/*"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPUT.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(2000),
							ExposedHeaders: []*string{
								to.StringPtr("x-ms-meta-12345675754564*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("x-ms-meta-abc"),
								to.StringPtr("x-ms-meta-data*"),
								to.StringPtr("x-ms-meta-target*"),
							},
						},
					},
				},
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)
}

func TestFileServicesClient_GetServiceProperties(t *testing.T) {
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

	// put file services
	fileServicesClient := armstorage.NewFileServicesClient(subscriptionID,cred,opt)
	putResp,err := fileServicesClient.SetServiceProperties(
		ctx,
		rgName,
		scName,
		armstorage.FileServiceProperties{
			FileServiceProperties: &armstorage.FileServicePropertiesProperties{
				Cors: &armstorage.CorsRules{
					CorsRules: []*armstorage.CorsRule{
						{
							AllowedOrigins: []*string{
								to.StringPtr("http://www.contoso.com"),
								to.StringPtr("http://www.fabrikam.com"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemHEAD.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemMERGE.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemOPTIONS.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPOST.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPUT.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(100),
							ExposedHeaders: []*string{
								to.StringPtr("x-ms-meta-*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("x-ms-meta-abc"),
								to.StringPtr("x-ms-meta-data*"),
								to.StringPtr("x-ms-meta-target*"),
							},
						},
						{
							AllowedOrigins: []*string{
								to.StringPtr("*"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(2),
							ExposedHeaders: []*string{
								to.StringPtr("*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("*"),
							},
						},
						{
							AllowedOrigins: []*string{
								to.StringPtr("http://www.abc23.com"),
								to.StringPtr("https://www.fabrikam.com/*"),
							},
							AllowedMethods: []*armstorage.CorsRuleAllowedMethodsItem{
								armstorage.CorsRuleAllowedMethodsItemGET.ToPtr(),
								armstorage.CorsRuleAllowedMethodsItemPUT.ToPtr(),
							},
							MaxAgeInSeconds: to.Int32Ptr(2000),
							ExposedHeaders: []*string{
								to.StringPtr("x-ms-meta-12345675754564*"),
							},
							AllowedHeaders: []*string{
								to.StringPtr("x-ms-meta-abc"),
								to.StringPtr("x-ms-meta-data*"),
								to.StringPtr("x-ms-meta-target*"),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)

	// get file service
	getResp,err:=fileServicesClient.GetServiceProperties(ctx,rgName,scName,nil)
	require.NoError(t, err)
	require.Equal(t, scName,*getResp.Name)
}

func TestFileServicesClient_List(t *testing.T) {
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

	// put file services
	fileServicesClient := armstorage.NewFileServicesClient(subscriptionID,cred,opt)
	listPager,err := fileServicesClient.List(ctx,rgName,scName,nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(listPager.Value),0)
}