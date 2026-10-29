// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

func TestSerdePreservesModel(t *testing.T) {
	for _, nextLink := range []*string{nil, to.Ptr(""), to.Ptr("https://example.vault.azure.net/deletedkeys?skiptoken=next")} {
		deleted := &azkeys.DeletedKeyPropertiesListResult{NextLink: nextLink}
		keys := &azkeys.KeyPropertiesListResult{NextLink: nextLink}
		for i := 0; i < 4; i++ {
			id := to.Ptr(azkeys.ID(fmt.Sprintf("https://example.vault.azure.net/keys/key%d/version", i)))
			deleted.Value = append(deleted.Value, &azkeys.DeletedKeyProperties{KID: id})
			keys.Value = append(keys.Value, &azkeys.KeyProperties{KID: id})
		}
		for _, test := range []struct {
			name  string
			model serdeModel
		}{
			{"deletedKeys", deleted},
			{"keys", keys},
		} {
			page := "lastPage"
			if nextLink != nil {
				page = "emptyNextLink"
				if *nextLink != "" {
					page = "nextPage"
				}
			}
			t.Run(test.name+"/"+page, func(t *testing.T) {
				before, err := test.model.MarshalJSON()
				require.NoError(t, err)

				testSerde(t, test.model)

				after, err := test.model.MarshalJSON()
				require.NoError(t, err)
				require.JSONEq(t, string(before), string(after))
			})
		}
	}
}
