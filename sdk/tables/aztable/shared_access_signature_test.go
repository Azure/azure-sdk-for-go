// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSasCreateTable(t *testing.T) {

	tempSAS := "blahblahblah"
	sasCredential, err := NewAzureSasCredential(tempSAS)
	require.NoError(t, err)

	client, err := NewTableClient("sastable", "https://seankaneprim.table.core.windows.net", sasCredential, nil)
	require.NoError(t, err)

	_, err = client.Create(context.Background())
	require.NoError(t, err)
}
