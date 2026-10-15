// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryPagerResolvesBinaryPreference(t *testing.T) {
	for _, tt := range []struct {
		name    string
		env     string
		client  *BinaryEncodingOptions
		request *BinaryEncodingOptions
		want    bool
	}{
		{"environment enabled", "true", nil, nil, true},
		{"environment disabled", "false", nil, nil, false},
		{"client enabled", "false", &BinaryEncodingOptions{}, nil, true},
		{"client disabled", "true", &BinaryEncodingOptions{Enabled: to(false)}, nil, false},
		{"request enabled", "false", &BinaryEncodingOptions{Enabled: to(false)}, &BinaryEncodingOptions{}, true},
		{"request disabled", "true", &BinaryEncodingOptions{}, &BinaryEncodingOptions{Enabled: to(false)}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", tt.env)
			client, err := NewClientWithKey("https://myaccount.documents.azure.com", mustKeyCredential(t),
				&ClientOptions{BinaryEncoding: tt.client})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			container, err := client.NewContainer("db", "items")
			require.NoError(t, err)
			pager := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"), NewFeedScopeForFullContainer(),
				&QueryOptions{Operation: OperationOptions{BinaryEncoding: tt.request}})
			defer func() { require.NoError(t, pager.Close()) }()
			require.NoError(t, pager.validationErr)
			require.NotNil(t, pager.req.options.Operation.BinaryEncoding)
			require.Equal(t, tt.want, pager.req.options.Operation.BinaryEncoding.enabled())
			if tt.request != nil && tt.request.Enabled != nil {
				*tt.request.Enabled = !tt.want
			}
			if tt.client != nil && tt.client.Enabled != nil {
				*tt.client.Enabled = !tt.want
			}
			require.Equal(t, tt.want, pager.req.options.Operation.BinaryEncoding.enabled(), "pager must own its preference")
		})
	}
}
