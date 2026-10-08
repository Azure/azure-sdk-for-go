// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionalBatchRequiresCompletePartitionKeys(t *testing.T) {
	for _, test := range []struct {
		name     string
		metadata string
		key      PartitionKey
		message  string
	}{
		{"scalar", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`, NewPartitionKeyString("tenant"), ""},
		{"legacy hash", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":1}}`, NewPartitionKeyNull(), ""},
		{"hash without version", `{"partitionKey":{"paths":["/pk"],"kind":"Hash"}}`, NewPartitionKeyUndefined(), ""},
		{"two levels", `{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("tenant").AppendBool(false), ""},
		{"three levels", `{"partitionKey":{"paths":["/pk","/child","/leaf"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("tenant").AppendNull().AppendUndefined(), ""},
		{"two level prefix", `{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("tenant"), "exactly 2"},
		{"three level prefix", `{"partitionKey":{"paths":["/pk","/child","/leaf"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("tenant").AppendString("region"), "exactly 3"},
		{"too many", `{"partitionKey":{"paths":["/pk"],"kind":"Hash"}}`, NewPartitionKeyString("tenant").AppendString("region"), "exactly 1"},
		{"zero", `{"partitionKey":{"paths":["/pk"],"kind":"Hash"}}`, PartitionKey{}, "exactly 1"},
		{"invalid JSON", `{`, NewPartitionKeyNull(), "decoding"},
		{"missing definition", `{}`, NewPartitionKeyNull(), "no partition key paths"},
		{"missing paths", `{"partitionKey":{"kind":"Hash"}}`, NewPartitionKeyNull(), "no partition key paths"},
		{"range kind", `{"partitionKey":{"paths":["/pk"],"kind":"Range"}}`, NewPartitionKeyNull(), "Hash or MultiHash"},
		{"invalid hash paths", `{"partitionKey":{"paths":["/pk","/child"],"kind":"Hash","version":2}}`, NewPartitionKeyNull(), "unsupported"},
		{"invalid version", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":3}}`, NewPartitionKeyNull(), "unsupported"},
		{"legacy multihash", `{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash","version":1}}`, NewPartitionKeyNull(), "unsupported"},
		{"versionless multihash", `{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash"}}`, NewPartitionKeyNull(), "unsupported"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateCompleteBatchPartitionKey([]byte(test.metadata), test.key)
			if test.message == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, test.message)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Contains(t, []Code{CodeBadRequest, CodeSerializationFailed}, cosmosErr.Code)
		})
	}
	metadata := []byte(`{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash","version":2}}`)
	require.NoError(t, validateQueryPartitionKey(metadata, NewPartitionKeyString("tenant")),
		"extracting shared metadata validation must preserve query prefix support")
}
