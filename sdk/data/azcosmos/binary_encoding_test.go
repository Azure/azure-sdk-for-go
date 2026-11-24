// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryPagerPreservesBinaryInheritance(t *testing.T) {
	for _, tt := range []struct {
		name    string
		request *BinaryEncodingOptions
	}{
		{"inherit", nil},
		{"enabled", &BinaryEncodingOptions{}},
		{"disabled", &BinaryEncodingOptions{Disabled: true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			container := newTestContainer(t)
			pager := container.NewQueryItemsPager[json.RawMessage](NewQuery("SELECT * FROM c"), NewFeedScopeForFullContainer(),
				&QueryOptions{Operation: OperationOptions{BinaryEncoding: tt.request}})
			defer func() { require.NoError(t, pager.Close()) }()
			require.NoError(t, pager.validationErr)
			require.Equal(t, tt.request, pager.req.options.Operation.BinaryEncoding)
			if tt.request != nil {
				want := tt.request.Disabled
				tt.request.Disabled = !want
				require.Equal(t, want, pager.req.options.Operation.BinaryEncoding.Disabled, "pager must own its preference")
			}
		})
	}
}
