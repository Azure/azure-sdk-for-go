//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmediaservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mediaservices/armmediaservices/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/mediaservices/resource-manager/Microsoft.Media/stable/2021-11-01/examples/accountFilters-list-all.json
func ExampleAccountFiltersClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armmediaservices.NewAccountFiltersClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager("contoso",
		"contosomedia",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/mediaservices/resource-manager/Microsoft.Media/stable/2021-11-01/examples/accountFilters-get-by-name.json
func ExampleAccountFiltersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armmediaservices.NewAccountFiltersClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"contoso",
		"contosomedia",
		"accountFilterWithTrack",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/mediaservices/resource-manager/Microsoft.Media/stable/2021-11-01/examples/accountFilters-create.json
func ExampleAccountFiltersClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armmediaservices.NewAccountFiltersClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"contoso",
		"contosomedia",
		"newAccountFilter",
		armmediaservices.AccountFilter{
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
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/mediaservices/resource-manager/Microsoft.Media/stable/2021-11-01/examples/accountFilters-delete.json
func ExampleAccountFiltersClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armmediaservices.NewAccountFiltersClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"contoso",
		"contosomedia",
		"accountFilterWithTimeWindowAndTrack",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/mediaservices/resource-manager/Microsoft.Media/stable/2021-11-01/examples/accountFilters-update.json
func ExampleAccountFiltersClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armmediaservices.NewAccountFiltersClient("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"contoso",
		"contosomedia",
		"accountFilterWithTimeWindowAndTrack",
		armmediaservices.AccountFilter{
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
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
