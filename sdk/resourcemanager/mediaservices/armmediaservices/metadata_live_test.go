//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmediaservices_test

import (
	"context"
	"fmt"
	"testing"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/suite"
)

type MetadataTestSuite struct {
	suite.Suite

	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	accountName          string
	assetName            string
	contentKeyPolicyName string
	filterName           string
	jobName              string
	liveOutputName       string
	storageAccountId     string
	storageAccountName   string
	streamingLocatorName string
	location             string
	resourceGroupName    string
	subscriptionId       string
}

func (testsuite *MetadataTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/mediaservices/armmediaservices/testdata")

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, true)
	testsuite.assetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "assetnam", 14, false)
	testsuite.contentKeyPolicyName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "contentk", 14, false)
	testsuite.filterName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "filterna", 14, false)
	testsuite.jobName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "jobname", 13, false)
	testsuite.liveOutputName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "liveoutp", 14, false)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storageaccount", 20, true)
	testsuite.streamingLocatorName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "streamin", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MetadataTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(MetadataTestSuite))
}

func (testsuite *MetadataTestSuite) Prepare() {
	var err error
	// From step Create_StorageAccount
	template := map[string]any{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]any{
			"storageAccountId": map[string]any{
				"type":  "string",
				"value": "[resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName'))]",
			},
		},
		"parameters": map[string]any{
			"location": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.location,
			},
			"storageAccountName": map[string]any{
				"type":         "string",
				"defaultValue": testsuite.storageAccountName,
			},
		},
		"resources": []any{
			map[string]any{
				"name":       "[parameters('storageAccountName')]",
				"type":       "Microsoft.Storage/storageAccounts",
				"apiVersion": "2022-05-01",
				"kind":       "StorageV2",
				"location":   "[parameters('location')]",
				"properties": map[string]any{
					"accessTier":                   "Hot",
					"allowBlobPublicAccess":        true,
					"allowCrossTenantReplication":  true,
					"allowSharedKeyAccess":         true,
					"defaultToOAuthAuthentication": false,
					"dnsEndpointType":              "Standard",
					"encryption": map[string]any{
						"keySource":                       "Microsoft.Storage",
						"requireInfrastructureEncryption": false,
						"services": map[string]any{
							"blob": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
							"file": map[string]any{
								"enabled": true,
								"keyType": "Account",
							},
						},
					},
					"minimumTlsVersion": "TLS1_2",
					"networkAcls": map[string]any{
						"bypass":              "AzureServices",
						"defaultAction":       "Allow",
						"ipRules":             []any{},
						"virtualNetworkRules": []any{},
					},
					"publicNetworkAccess":      "Enabled",
					"supportsHttpsTrafficOnly": true,
				},
				"sku": map[string]any{
					"name": "Standard_LRS",
					"tier": "Standard",
				},
			},
		},
		"variables": map[string]any{},
	}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template: template,
			Mode:     to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend, err := testutil.CreateDeployment(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName, "Create_StorageAccount", &deployment)
	testsuite.Require().NoError(err)
	testsuite.storageAccountId = deploymentExtend.Properties.Outputs.(map[string]interface{})["storageAccountId"].(map[string]interface{})["value"].(string)

	// From step Mediaservices_CreateOrUpdate
	fmt.Println("Call operation: Mediaservices_CreateOrUpdate")
	client, err := armmediaservices.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	clientCreateOrUpdateResponsePoller, err := client.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmediaservices.MediaService{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
		Properties: &armmediaservices.MediaServiceProperties{
			StorageAccounts: []*armmediaservices.StorageAccount{
				{
					Type: to.Ptr(armmediaservices.StorageAccountTypePrimary),
					ID:   to.Ptr(testsuite.storageAccountId),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, clientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

func (testsuite *MetadataTestSuite) TestMetadataTests() {
	var err error
	// From step Assets_CreateOrUpdate
	fmt.Println("Call operation: Assets_CreateOrUpdate")
	assetsClient, err := armmediaservices.NewAssetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = assetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, armmediaservices.Asset{
		Properties: &armmediaservices.AssetProperties{
			Description:        to.Ptr("A documentary showing the ascent of Mount Logan"),
			StorageAccountName: to.Ptr(testsuite.storageAccountName),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Assets_List
	fmt.Println("Call operation: Assets_List")
	assetsClientNewListPager := assetsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armmediaservices.AssetsClientListOptions{Filter: nil,
		Top:     nil,
		Orderby: nil,
	})
	for assetsClientNewListPager.More() {
		_, err := assetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Assets_Get
	fmt.Println("Call operation: Assets_Get")
	_, err = assetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, nil)
	testsuite.Require().NoError(err)

	// From step Assets_Update
	fmt.Println("Call operation: Assets_Update")
	_, err = assetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, armmediaservices.Asset{
		Properties: &armmediaservices.AssetProperties{
			Description: to.Ptr("A documentary showing the ascent of Mount Baker in HD"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Assets_ListContainerSas
	fmt.Println("Call operation: Assets_ListContainerSas")
	_, err = assetsClient.ListContainerSas(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, armmediaservices.ListContainerSasInput{
		ExpiryTime:  to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-01T10:00:00.007Z"); return t }()),
		Permissions: to.Ptr(armmediaservices.AssetContainerPermissionReadWrite),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Assets_ListStreamingLocators
	fmt.Println("Call operation: Assets_ListStreamingLocators")
	_, err = assetsClient.ListStreamingLocators(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingPolicies_Create
	fmt.Println("Call operation: StreamingPolicies_Create")
	streamingPoliciesClient, err := armmediaservices.NewStreamingPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = streamingPoliciesClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, "clearStreamingPolicy", armmediaservices.StreamingPolicy{
		Properties: &armmediaservices.StreamingPolicyProperties{
			NoEncryption: &armmediaservices.NoEncryption{
				EnabledProtocols: &armmediaservices.EnabledProtocols{
					Dash:            to.Ptr(true),
					Download:        to.Ptr(true),
					Hls:             to.Ptr(true),
					SmoothStreaming: to.Ptr(true),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StreamingPolicies_List
	fmt.Println("Call operation: StreamingPolicies_List")
	streamingPoliciesClientNewListPager := streamingPoliciesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armmediaservices.StreamingPoliciesClientListOptions{Filter: nil,
		Top:     nil,
		Orderby: nil,
	})
	for streamingPoliciesClientNewListPager.More() {
		_, err := streamingPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StreamingPolicies_Get
	fmt.Println("Call operation: StreamingPolicies_Get")
	_, err = streamingPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, "clearStreamingPolicy", nil)
	testsuite.Require().NoError(err)

	// From step StreamingLocators_Create
	fmt.Println("Call operation: StreamingLocators_Create")
	streamingLocatorsClient, err := armmediaservices.NewStreamingLocatorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = streamingLocatorsClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingLocatorName, armmediaservices.StreamingLocator{
		Properties: &armmediaservices.StreamingLocatorProperties{
			AssetName:           to.Ptr(testsuite.assetName),
			StreamingPolicyName: to.Ptr("clearStreamingPolicy"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StreamingLocators_List
	fmt.Println("Call operation: StreamingLocators_List")
	streamingLocatorsClientNewListPager := streamingLocatorsClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armmediaservices.StreamingLocatorsClientListOptions{Filter: nil,
		Top:     nil,
		Orderby: nil,
	})
	for streamingLocatorsClientNewListPager.More() {
		_, err := streamingLocatorsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step StreamingLocators_Get
	fmt.Println("Call operation: StreamingLocators_Get")
	_, err = streamingLocatorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingLocatorName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingLocators_ListContentKeys
	fmt.Println("Call operation: StreamingLocators_ListContentKeys")
	_, err = streamingLocatorsClient.ListContentKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingLocatorName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingLocators_ListPaths
	fmt.Println("Call operation: StreamingLocators_ListPaths")
	_, err = streamingLocatorsClient.ListPaths(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingLocatorName, nil)
	testsuite.Require().NoError(err)

	// From step ContentKeyPolicies_CreateOrUpdate
	fmt.Println("Call operation: ContentKeyPolicies_CreateOrUpdate")
	contentKeyPoliciesClient, err := armmediaservices.NewContentKeyPoliciesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = contentKeyPoliciesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.contentKeyPolicyName, armmediaservices.ContentKeyPolicy{
		Properties: &armmediaservices.ContentKeyPolicyProperties{
			Description: to.Ptr("ArmPolicyDescription"),
			Options: []*armmediaservices.ContentKeyPolicyOption{
				{
					Name: to.Ptr("ArmPolicyOptionName"),
					Configuration: &armmediaservices.ContentKeyPolicyPlayReadyConfiguration{
						ODataType: to.Ptr("#Microsoft.Media.ContentKeyPolicyPlayReadyConfiguration"),
						Licenses: []*armmediaservices.ContentKeyPolicyPlayReadyLicense{
							{
								AllowTestDevices: to.Ptr(true),
								BeginDate:        to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-10-16T18:22:53.46Z"); return t }()),
								ContentKeyLocation: &armmediaservices.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader{
									ODataType: to.Ptr("#Microsoft.Media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader"),
								},
								ContentType: to.Ptr(armmediaservices.ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload),
								LicenseType: to.Ptr(armmediaservices.ContentKeyPolicyPlayReadyLicenseTypePersistent),
								PlayRight: &armmediaservices.ContentKeyPolicyPlayReadyPlayRight{
									AllowPassingVideoContentToUnknownOutput:            to.Ptr(armmediaservices.ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed),
									DigitalVideoOnlyContentRestriction:                 to.Ptr(false),
									ImageConstraintForAnalogComponentVideoRestriction:  to.Ptr(true),
									ImageConstraintForAnalogComputerMonitorRestriction: to.Ptr(false),
									ScmsRestriction: to.Ptr[int32](2),
								},
								SecurityLevel: to.Ptr(armmediaservices.SecurityLevelSL150),
							}},
					},
					Restriction: &armmediaservices.ContentKeyPolicyOpenRestriction{
						ODataType: to.Ptr("#Microsoft.Media.ContentKeyPolicyOpenRestriction"),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ContentKeyPolicies_List
	fmt.Println("Call operation: ContentKeyPolicies_List")
	contentKeyPoliciesClientNewListPager := contentKeyPoliciesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armmediaservices.ContentKeyPoliciesClientListOptions{Filter: nil,
		Top:     nil,
		Orderby: nil,
	})
	for contentKeyPoliciesClientNewListPager.More() {
		_, err := contentKeyPoliciesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ContentKeyPolicies_Get
	fmt.Println("Call operation: ContentKeyPolicies_Get")
	_, err = contentKeyPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.contentKeyPolicyName, nil)
	testsuite.Require().NoError(err)

	// From step ContentKeyPolicies_Update
	fmt.Println("Call operation: ContentKeyPolicies_Update")
	_, err = contentKeyPoliciesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.contentKeyPolicyName, armmediaservices.ContentKeyPolicy{
		Properties: &armmediaservices.ContentKeyPolicyProperties{
			Description: to.Ptr("Updated Policy"),
			Options: []*armmediaservices.ContentKeyPolicyOption{
				{
					Name: to.Ptr("ClearKeyOption"),
					Configuration: &armmediaservices.ContentKeyPolicyClearKeyConfiguration{
						ODataType: to.Ptr("#Microsoft.Media.ContentKeyPolicyClearKeyConfiguration"),
					},
					Restriction: &armmediaservices.ContentKeyPolicyOpenRestriction{
						ODataType: to.Ptr("#Microsoft.Media.ContentKeyPolicyOpenRestriction"),
					},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step ContentKeyPolicies_GetPolicyPropertiesWithSecrets
	fmt.Println("Call operation: ContentKeyPolicies_GetPolicyPropertiesWithSecrets")
	_, err = contentKeyPoliciesClient.GetPolicyPropertiesWithSecrets(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.contentKeyPolicyName, nil)
	testsuite.Require().NoError(err)

	// From step AssetFilters_CreateOrUpdate
	fmt.Println("Call operation: AssetFilters_CreateOrUpdate")
	assetFiltersClient, err := armmediaservices.NewAssetFiltersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = assetFiltersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, testsuite.filterName, armmediaservices.AssetFilter{
		Properties: &armmediaservices.MediaFilterProperties{
			FirstQuality: &armmediaservices.FirstQuality{
				Bitrate: to.Ptr[int32](128000),
			},
			PresentationTimeRange: &armmediaservices.PresentationTimeRange{
				EndTimestamp:               to.Ptr[int64](170000000),
				ForceEndTimestamp:          to.Ptr(false),
				LiveBackoffDuration:        to.Ptr[int64](0),
				PresentationWindowDuration: to.Ptr[int64](9223372036854775000),
				StartTimestamp:             to.Ptr[int64](0),
				Timescale:                  to.Ptr[int64](10000000),
			},
			Tracks: []*armmediaservices.FilterTrackSelection{
				{
					TrackSelections: []*armmediaservices.FilterTrackPropertyCondition{
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeType),
							Value:     to.Ptr("Audio"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationNotEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeLanguage),
							Value:     to.Ptr("en"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationNotEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeFourCC),
							Value:     to.Ptr("EC-3"),
						}},
				},
				{
					TrackSelections: []*armmediaservices.FilterTrackPropertyCondition{
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeType),
							Value:     to.Ptr("Video"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeBitrate),
							Value:     to.Ptr("3000000-5000000"),
						}},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AssetFilters_List
	fmt.Println("Call operation: AssetFilters_List")
	assetFiltersClientNewListPager := assetFiltersClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, nil)
	for assetFiltersClientNewListPager.More() {
		_, err := assetFiltersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AssetFilters_Get
	fmt.Println("Call operation: AssetFilters_Get")
	_, err = assetFiltersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, testsuite.filterName, nil)
	testsuite.Require().NoError(err)

	// From step AssetFilters_Update
	fmt.Println("Call operation: AssetFilters_Update")
	_, err = assetFiltersClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, testsuite.filterName, armmediaservices.AssetFilter{
		Properties: &armmediaservices.MediaFilterProperties{
			FirstQuality: &armmediaservices.FirstQuality{
				Bitrate: to.Ptr[int32](128000),
			},
			PresentationTimeRange: &armmediaservices.PresentationTimeRange{
				EndTimestamp:               to.Ptr[int64](170000000),
				ForceEndTimestamp:          to.Ptr(false),
				LiveBackoffDuration:        to.Ptr[int64](0),
				PresentationWindowDuration: to.Ptr[int64](9223372036854775000),
				StartTimestamp:             to.Ptr[int64](10),
				Timescale:                  to.Ptr[int64](10000000),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AccountFilters_CreateOrUpdate
	fmt.Println("Call operation: AccountFilters_CreateOrUpdate")
	accountFiltersClient, err := armmediaservices.NewAccountFiltersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountFiltersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.filterName, armmediaservices.AccountFilter{
		Properties: &armmediaservices.MediaFilterProperties{
			FirstQuality: &armmediaservices.FirstQuality{
				Bitrate: to.Ptr[int32](128000),
			},
			PresentationTimeRange: &armmediaservices.PresentationTimeRange{
				EndTimestamp:               to.Ptr[int64](170000000),
				ForceEndTimestamp:          to.Ptr(false),
				LiveBackoffDuration:        to.Ptr[int64](0),
				PresentationWindowDuration: to.Ptr[int64](9223372036854775000),
				StartTimestamp:             to.Ptr[int64](0),
				Timescale:                  to.Ptr[int64](10000000),
			},
			Tracks: []*armmediaservices.FilterTrackSelection{
				{
					TrackSelections: []*armmediaservices.FilterTrackPropertyCondition{
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeType),
							Value:     to.Ptr("Audio"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationNotEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeLanguage),
							Value:     to.Ptr("en"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationNotEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeFourCC),
							Value:     to.Ptr("EC-3"),
						}},
				},
				{
					TrackSelections: []*armmediaservices.FilterTrackPropertyCondition{
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeType),
							Value:     to.Ptr("Video"),
						},
						{
							Operation: to.Ptr(armmediaservices.FilterTrackPropertyCompareOperationEqual),
							Property:  to.Ptr(armmediaservices.FilterTrackPropertyTypeBitrate),
							Value:     to.Ptr("3000000-5000000"),
						}},
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AccountFilters_List
	fmt.Println("Call operation: AccountFilters_List")
	accountFiltersClientNewListPager := accountFiltersClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for accountFiltersClientNewListPager.More() {
		_, err := accountFiltersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AccountFilters_Get
	fmt.Println("Call operation: AccountFilters_Get")
	_, err = accountFiltersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.filterName, nil)
	testsuite.Require().NoError(err)

	// From step AccountFilters_Update
	fmt.Println("Call operation: AccountFilters_Update")
	_, err = accountFiltersClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.filterName, armmediaservices.AccountFilter{
		Properties: &armmediaservices.MediaFilterProperties{
			FirstQuality: &armmediaservices.FirstQuality{
				Bitrate: to.Ptr[int32](128000),
			},
			PresentationTimeRange: &armmediaservices.PresentationTimeRange{
				EndTimestamp:               to.Ptr[int64](170000000),
				ForceEndTimestamp:          to.Ptr(false),
				LiveBackoffDuration:        to.Ptr[int64](0),
				PresentationWindowDuration: to.Ptr[int64](9223372036854775000),
				StartTimestamp:             to.Ptr[int64](10),
				Timescale:                  to.Ptr[int64](10000000),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step StreamingPolicies_Delete
	fmt.Println("Call operation: StreamingPolicies_Delete")
	_, err = streamingPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, "secureStreamingPolicyWithCommonEncryptionCbcsOnly", nil)
	testsuite.Require().NoError(err)

	// From step AccountFilters_Delete
	fmt.Println("Call operation: AccountFilters_Delete")
	_, err = accountFiltersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.filterName, nil)
	testsuite.Require().NoError(err)

	// From step AssetFilters_Delete
	fmt.Println("Call operation: AssetFilters_Delete")
	_, err = assetFiltersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, testsuite.filterName, nil)
	testsuite.Require().NoError(err)

	// From step ContentKeyPolicies_Delete
	fmt.Println("Call operation: ContentKeyPolicies_Delete")
	_, err = contentKeyPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.contentKeyPolicyName, nil)
	testsuite.Require().NoError(err)

	// From step StreamingLocators_Delete
	fmt.Println("Call operation: StreamingLocators_Delete")
	_, err = streamingLocatorsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.streamingLocatorName, nil)
	testsuite.Require().NoError(err)

	// From step Assets_Delete
	fmt.Println("Call operation: Assets_Delete")
	_, err = assetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.assetName, nil)
	testsuite.Require().NoError(err)
}
