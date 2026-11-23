// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build azcosmos_native_preview && cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"context"
	"errors"
)

func feedRangeSupportError() error { return nil }

func newFeedRange(minInclusive, maxExclusive string) (FeedRange, error) {
	min, freeMin := toNativeString(minInclusive)
	defer C.free(freeMin)
	max, freeMax := toNativeString(maxExclusive)
	defer C.free(freeMax)
	var handle *C.cosmos_feed_range_t
	if status := C.cosmos_feed_range_from_epk(min, max, &handle); status != 0 {
		return FeedRange{}, statusError(status, nil, "constructing feed range")
	}
	defer C.cosmos_feed_range_free(handle)
	return copyNativeFeedRange(handle)
}

func parseNativeFeedRange(serialized string) (*C.cosmos_feed_range_t, C.cosmos_status_code_t) {
	value, allocation := toNativeString(serialized)
	defer C.free(allocation)
	var handle *C.cosmos_feed_range_t
	status := C.cosmos_feed_range_parse(value, &handle)
	return handle, status
}

func parseFeedRange(serialized string) (FeedRange, error) {
	handle, status := parseNativeFeedRange(serialized)
	if status != 0 {
		return FeedRange{}, statusError(status, nil, "parsing feed range")
	}
	defer C.cosmos_feed_range_free(handle)
	return copyNativeFeedRange(handle)
}

func copyNativeFeedRange(handle *C.cosmos_feed_range_t) (FeedRange, error) {
	if handle == nil {
		return FeedRange{}, errors.New("azcosmos: native completion has no feed range")
	}
	var serialized *C.char
	if status := C.cosmos_feed_range_serialize(handle, &serialized); status != 0 {
		return FeedRange{}, statusError(status, nil, "serializing feed range")
	}
	defer C.cosmos_string_free(serialized)
	canonical, status := parseNativeFeedRange(C.GoString(serialized))
	if status != 0 {
		return FeedRange{}, statusError(status, nil, "normalizing feed range")
	}
	defer C.cosmos_feed_range_free(canonical)
	handle = canonical
	var info C.cosmos_feed_range_info_t
	if status := C.cosmos_feed_range_inspect(handle, &info); status != 0 {
		return FeedRange{}, statusError(status, nil, "inspecting feed range")
	}
	defer C.cosmos_string_free(info.min_inclusive)
	defer C.cosmos_string_free(info.max_exclusive)
	return FeedRange{
		serialized: C.GoString(serialized), minInclusive: C.GoString(info.min_inclusive),
		maxExclusive: C.GoString(info.max_exclusive), logical: info.is_logical_partition != 0,
	}, nil
}

func compareFeedRanges(left, right FeedRange, overlaps bool) (bool, error) {
	l, status := parseNativeFeedRange(left.serialized)
	if status != 0 {
		return false, statusError(status, nil, "parsing first feed range")
	}
	defer C.cosmos_feed_range_free(l)
	r, status := parseNativeFeedRange(right.serialized)
	if status != 0 {
		return false, statusError(status, nil, "parsing second feed range")
	}
	defer C.cosmos_feed_range_free(r)
	var result C.cosmos_feed_range_relation_t
	status = C.cosmos_feed_range_relate(l, r, &result)
	if status != 0 {
		return false, statusError(status, nil, "comparing feed ranges")
	}
	if overlaps {
		return result.overlaps != 0, nil
	}
	return result.is_subset != 0, nil
}

func (c *ContainerClient) getFeedRanges(ctx context.Context, forceRefresh bool) (GetFeedRangesResponse, error) {
	client := c.database.client
	ctx, release, err := client.queryContext(ctx, OperationOptions{})
	if err != nil {
		return GetFeedRangesResponse{}, err
	}
	defer release()
	driver, container, _, err := client.resolveFeedScope(ctx, c.database.id, c.id, PartitionKey{}, true, OperationOptions{})
	if err != nil {
		return GetFeedRangesResponse{}, err
	}
	defer C.cosmos_container_ref_free(container)
	var refresh C.uint8_t
	if forceRefresh {
		refresh = 1
	}
	result, err := client.driver.awaitCompletion(ctx, "discovering feed ranges",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, status *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			return C.cosmos_container_read_feed_ranges_submit(driver, container, refresh, queue, cookie, status)
		})
	if err != nil {
		return GetFeedRangesResponse{}, err
	}
	if result.err != nil {
		return GetFeedRangesResponse{}, result.err
	}
	if result.feedRanges == nil {
		return GetFeedRangesResponse{}, setupChargeError(errors.New("azcosmos: discovery completion has no feed range list"), result.response.Response, "discovering feed ranges")
	}
	return GetFeedRangesResponse{Response: result.response.Response, FeedRanges: result.feedRanges}, nil
}

func (c *ContainerClient) getFeedRangeForPartitionKey(ctx context.Context, partitionKey PartitionKey) (FeedRangeResponse, error) {
	_, container, setup, err := c.database.client.resolveFeedScope(ctx, c.database.id, c.id, partitionKey, false, OperationOptions{})
	if err != nil {
		return FeedRangeResponse{}, addSetupCharge(err, setup, "constructing partition key feed range")
	}
	defer C.cosmos_container_ref_free(container)
	components, freeComponents := partitionKey.toNative()
	defer freeComponents()
	var pk *C.cosmos_partition_key_t
	if status := C.cosmos_partition_key_create(components, partitionKey.partitionKeyLen(), &pk); status != 0 {
		return FeedRangeResponse{}, addSetupCharge(statusError(status, nil, "constructing partition key"), setup, "constructing partition key feed range")
	}
	defer C.cosmos_partition_key_free(pk)
	var handle *C.cosmos_feed_range_t
	if status := C.cosmos_feed_range_for_partition_key(container, pk, &handle); status != 0 {
		return FeedRangeResponse{}, addSetupCharge(statusError(status, nil, "constructing partition key feed range"), setup, "constructing partition key feed range")
	}
	defer C.cosmos_feed_range_free(handle)
	result, err := copyNativeFeedRange(handle)
	if err != nil {
		return FeedRangeResponse{}, addSetupCharge(err, setup, "copying partition key feed range")
	}
	return FeedRangeResponse{Response: setup, FeedRange: result}, nil
}

func translateFeedRanges(completion *C.cosmos_completion_t, result *completionResult) {
	list := C.cosmos_completion_take_feed_ranges(completion)
	if list == nil {
		return
	}
	defer C.cosmos_feed_range_list_free(list)
	if result.err != nil {
		return
	}
	length := C.cosmos_feed_range_list_len(list)
	if uint64(length) > uint64(int(^uint(0)>>1)) {
		result.err = setupChargeError(errors.New("azcosmos: invalid native feed range count"), result.response.Response, "copying feed ranges")
		return
	}
	ranges := make([]FeedRange, int(length))
	for i := range ranges {
		var handle *C.cosmos_feed_range_t
		status := C.cosmos_feed_range_list_get(list, C.uintptr_t(i), &handle)
		if status != 0 {
			result.err = setupChargeError(statusError(status, nil, "reading feed range list"), result.response.Response, "copying feed ranges")
			return
		}
		value, err := copyNativeFeedRange(handle)
		C.cosmos_feed_range_free(handle)
		if err != nil {
			result.err = setupChargeError(err, result.response.Response, "copying feed ranges")
			return
		}
		ranges[i] = value
	}
	result.feedRanges = ranges
}
