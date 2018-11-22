package easykey

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"net/url"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func mustParse(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

func TestKeyLookup(t *testing.T) {
	secrets := map[string][]versioner{
		"key0": []versioner{
			Secret{
				ID:    mustURLToID(mustParse("https://exp.vault.azure.net/secrets/key0/version0")),
				Value: "value0",
			}, Secret{
				ID:    mustURLToID(mustParse("https://exp.vault.azure.net/secrets/key0/version1")),
				Value: "value1",
			},
		},
		"key1": []versioner{
			Secret{
				ID:    mustURLToID(mustParse("https://exp.vault.azure.net/secrets/key1/version0")),
				Value: "value0",
			},
		},
	}

	tests := []struct {
		desc         string
		key, version string
		wantErr      bool
		want         Secret
	}{
		{
			desc:    "key doesn't exist",
			key:     "notexist",
			version: LatestVersion,
			wantErr: true,
		},
		{
			desc:    "version doesn't exist",
			key:     "key0",
			version: "version2",
			wantErr: true,
		},
		{
			desc:    "success getting key/version",
			key:     "key0",
			version: "version1",
			want:    secrets["key0"][1].(Secret),
		},
		{
			desc:    "success getting key/LatestVersion",
			key:     "key1",
			version: "version0",
			want:    secrets["key1"][0].(Secret),
		},
	}

	k := keyLookup(secrets)
	for _, test := range tests {
		got, err := k.value(test.key, test.version)
		switch {
		case err == nil && test.wantErr:
			t.Errorf("TestKeyLookup(%s): got err == nil, want err != nil", test.desc)
		case err != nil && !test.wantErr:
			t.Errorf("TestKeyLookup(%s): got err == %s, want err == nil", test.desc, err)
		case err != nil:
			continue
		}

		if diff := pretty.Compare(test.want, got); diff != "" {
			t.Errorf("TestKeyLookup(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}
