// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestThrottledCompletionCopiesAndClassifiesEveryField(t *testing.T) {
	for _, tt := range []struct {
		name            string
		packedSubStatus int
	}{
		{"header refines absent packed sub-status", 0},
		{"packed sub-status is authoritative", 3200},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := syntheticThrottledCompletion(tt.packedSubStatus)

			require.Equal(t, CodeThrottled, err.Code)
			require.Equal(t, http.StatusTooManyRequests, err.StatusCode)
			require.Equal(t, 3200, err.SubStatus)
			require.Equal(t, 4.5, err.RequestCharge)
			require.Equal(t, "activity-123", err.ActivityID)
			require.Equal(t, SessionToken("1:2"), err.SessionToken)
			require.Equal(t, azcore.ETag("\"etag\""), err.ETag)
			require.Equal(t, 125*time.Millisecond, err.RetryAfter)
			require.Equal(t, "Request rate is large.", err.Message)
			require.JSONEq(t, `{"code":"TooManyRequests","message":"slow down"}`, string(err.Body))
			require.True(t, err.FromWire)
		})
	}
}
