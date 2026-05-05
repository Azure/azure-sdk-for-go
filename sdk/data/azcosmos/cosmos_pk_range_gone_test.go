// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func Test_isPartitionKeyRangeGoneError_410WithSplitSubstatus(t *testing.T) {
	require.True(t, isPartitionKeyRangeGoneError(http.StatusGone, subStatusPartitionKeyRangeGone))
	require.True(t, isPartitionKeyRangeGoneError(http.StatusGone, subStatusCompletingSplit))
	require.True(t, isPartitionKeyRangeGoneError(http.StatusGone, subStatusCompletingPartitionMigration))
}

func Test_isPartitionKeyRangeGoneError_410WithOtherSubstatus(t *testing.T) {
	require.False(t, isPartitionKeyRangeGoneError(http.StatusGone, "9999"))
	require.False(t, isPartitionKeyRangeGoneError(http.StatusGone, ""))
}

func Test_isPartitionKeyRangeGoneError_non410(t *testing.T) {
	require.False(t, isPartitionKeyRangeGoneError(http.StatusNotFound, subStatusPartitionKeyRangeGone))
	require.False(t, isPartitionKeyRangeGoneError(http.StatusOK, subStatusCompletingSplit))
}

func Test_isPKRangeGoneResponseError_nonResponseError(t *testing.T) {
	err := errors.New("some random error")
	require.False(t, isPKRangeGoneResponseError(err))
}

func Test_isPKRangeGoneResponseError_non410ResponseError(t *testing.T) {
	err := &azcore.ResponseError{
		StatusCode: http.StatusNotFound,
	}
	require.False(t, isPKRangeGoneResponseError(err))
}

func Test_isPKRangeGoneResponseError_410WithSubstatus(t *testing.T) {
	header := http.Header{}
	header.Set(cosmosHeaderSubstatus, subStatusCompletingSplit)
	resp := &http.Response{
		StatusCode: http.StatusGone,
		Header:     header,
	}
	err := &azcore.ResponseError{
		StatusCode:  http.StatusGone,
		RawResponse: resp,
	}
	require.True(t, isPKRangeGoneResponseError(err))
}

func Test_isPKRangeGoneResponseError_410WithoutSubstatus(t *testing.T) {
	header := http.Header{}
	header.Set(cosmosHeaderSubstatus, "9999")
	resp := &http.Response{
		StatusCode: http.StatusGone,
		Header:     header,
	}
	err := &azcore.ResponseError{
		StatusCode:  http.StatusGone,
		RawResponse: resp,
	}
	require.False(t, isPKRangeGoneResponseError(err))
}

func Test_isPKRangeGoneResponseError_410WithNilRawResponse(t *testing.T) {
	err := &azcore.ResponseError{
		StatusCode:  http.StatusGone,
		RawResponse: nil,
	}
	// Conservative: any 410 without a raw response is treated as PKRange gone
	require.True(t, isPKRangeGoneResponseError(err))
}
