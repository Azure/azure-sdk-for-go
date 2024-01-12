//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestAuthenticationRecord_MarshalUnmarshal(t *testing.T) {
	for _, test := range []struct {
		desc, version string
		err           bool
	}{
		{desc: "no version", err: true},
		{desc: "supported version", version: supportedAuthRecordVersions[0]},
		{desc: "unsupported version", err: true, version: "42"},
	} {
		t.Run(test.desc, func(t *testing.T) {
			record := authenticationRecord{
				Authority:     "authority",
				ClientID:      "client",
				HomeAccountID: "oid.tid",
				TenantID:      "tenant",
				Username:      "user",
				Version:       test.version,
			}
			marshaled, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			var unmarshaled authenticationRecord
			err = json.Unmarshal(marshaled, &unmarshaled)
			if err != nil {
				if !test.err {
					t.Fatal(err)
				}
				if actual := err.Error(); !strings.Contains(actual, "version") {
					t.Fatalf("unexpected error %q", actual)
				}
				return
			} else if test.err {
				t.Fatal("expected an error")
			}
			if !reflect.DeepEqual(unmarshaled, record) {
				t.Fatal("records should be equal")
			}
		})
	}
}
