// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !azcosmos_native_preview || !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import "context"

func feedRangeSupportError() error {
	return &Error{Code: CodeClientError, Message: "azcosmos: feed ranges require the unpublished native feed-range extension and an azcosmos_native_preview build linked to matching artifacts; native v0.3.0 does not expose this contract"}
}

func newFeedRange(string, string) (FeedRange, error) {
	return FeedRange{}, feedRangeSupportError()
}

func parseFeedRange(string) (FeedRange, error) {
	return FeedRange{}, feedRangeSupportError()
}

func compareFeedRanges(FeedRange, FeedRange, bool) (bool, error) {
	return false, feedRangeSupportError()
}

func (c *ContainerClient) getFeedRanges(context.Context, bool) (GetFeedRangesResponse, error) {
	return GetFeedRangesResponse{}, feedRangeSupportError()
}

func (c *ContainerClient) getFeedRangeForPartitionKey(context.Context, PartitionKey) (FeedRangeResponse, error) {
	return FeedRangeResponse{}, feedRangeSupportError()
}
