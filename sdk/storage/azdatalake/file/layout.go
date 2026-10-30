// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

// LayoutAwareRouting determines whether the parallel range reads of a DownloadBuffer or
// DownloadFile are routed to the endpoint that serves each range, based on the file's layout. It
// is a performance optimization only: the data downloaded is the same whatever the mode.
type LayoutAwareRouting = blob.LayoutAwareRouting

const (
	// LayoutAwareRoutingAuto is the zero value, and therefore the default when no value is
	// specified. The client library decides whether layout aware routing is used, and that
	// decision may change in a future release. Currently, LayoutAwareRoutingAuto resolves to
	// LayoutAwareRoutingDisabled; use LayoutAwareRoutingEnabled to opt in.
	LayoutAwareRoutingAuto LayoutAwareRouting = blob.LayoutAwareRoutingAuto

	// LayoutAwareRoutingEnabled opts in to layout aware routing. When the initial read of a download
	// carries the layout download hint and data remains, the layout of the remaining range is fetched
	// and cached (with automatic background refresh), and each remaining range is read from the
	// endpoint that serves it.
	LayoutAwareRoutingEnabled LayoutAwareRouting = blob.LayoutAwareRoutingEnabled

	// LayoutAwareRoutingDisabled never routes by layout; every read goes to the client's configured
	// endpoint.
	LayoutAwareRoutingDisabled LayoutAwareRouting = blob.LayoutAwareRoutingDisabled
)

// PossibleLayoutAwareRoutingValues returns the possible values for the LayoutAwareRouting const type.
func PossibleLayoutAwareRoutingValues() []LayoutAwareRouting {
	return blob.PossibleLayoutAwareRoutingValues()
}
