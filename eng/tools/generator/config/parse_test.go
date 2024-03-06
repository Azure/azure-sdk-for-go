// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
)

func TestParser_Parse(t *testing.T) {
	testCase := []struct {
		input    string
		expected *config.Config
	}{
		{
			input: `{
  "track1Requests": {
    "readme1": {
      "tag1": [{"requestLink": "link1" }]
    }
  }
}`,
			expected: &config.Config{
				Track1Requests: config.Track1ReleaseRequests{
					"readme1": {
						"tag1": []config.ReleaseRequestInfo{
							{
								RequestLink: "link1",
							},
						},
					},
				},
			},
		},
		{
			input: `{
  "track1Requests": {
    "readme1": {
      "tag1": [{"requestLink": "link1"}, {"requestLink": "link2"} ]
    }
  }
}`,
			expected: &config.Config{
				Track1Requests: config.Track1ReleaseRequests{
					"readme1": {
						"tag1": []config.ReleaseRequestInfo{
							{
								RequestLink: "link1",
							},
							{
								RequestLink: "link2",
							},
						},
					},
				},
			},
		},
		{
			input: `{
  "track1Requests": {
    "readme1": {
      "tag1": [{"requestLink": "link1"}, {"requestLink": "link2" }],
      "tag2": [{"requestLink": "link3"}]
    }
  }
}`,
			expected: &config.Config{
				Track1Requests: config.Track1ReleaseRequests{
					"readme1": {
						"tag1": []config.ReleaseRequestInfo{
							{
								RequestLink: "link1",
							},
							{
								RequestLink: "link2",
							},
						},
						"tag2": []config.ReleaseRequestInfo{
							{
								RequestLink: "link3",
							},
						},
					},
				},
			},
		},
		{
			input: `{
  "track1Requests": {
    "readme1": {
      "tag1": [{"requestLink": "link1"}, {"requestLink": "link2"}],
      "tag2": [{"requestLink": "link3"}]
    }
  },
  "refresh": {
  	"additionalOptions": ["--use=somethingNew"]
  }
}`,
			expected: &config.Config{
				Track1Requests: config.Track1ReleaseRequests{
					"readme1": {
						"tag1": []config.ReleaseRequestInfo{
							{
								RequestLink: "link1",
							},
							{
								RequestLink: "link2",
							},
						},
						"tag2": []config.ReleaseRequestInfo{
							{
								RequestLink: "link3",
							},
						},
					},
				},
				RefreshInfo: config.RefreshInfo{
					AdditionalFlags: []string{
						"--use=somethingNew",
					},
				},
			},
		},
	}

	for i, c := range testCase {
		t.Logf("testing %d", i)
		cfg, err := config.FromReader(strings.NewReader(c.input)).Parse()
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(cfg, c.expected) {
			t.Fatalf("expected '%s' but got '%s'", cfg.String(), c.expected.String())
		}
	}
}
