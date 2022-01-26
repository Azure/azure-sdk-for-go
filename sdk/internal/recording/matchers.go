//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Optional parameters for the SetBodilessMatcher operation
type MatcherOptions struct {
}

// SetBodilessMatcher adjusts the "match" operation to exclude the body when matching a request to a recording's entries.
// Pass in `nil` for `t` if you want the bodiless matcher to apply everywhere
func SetBodilessMatcher(t *testing.T, options *MatcherOptions) error {
	if recordMode != PlaybackMode {
		return nil
	}
	req, err := http.NewRequest("POST", "http://localhost:5000/Admin/SetMatcher", http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header["x-abstraction-identifier"] = []string{"BodilessMatcher"}
	if t != nil {
		req.Header["x-recording-id"] = []string{GetRecordingId(t)}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to enable BodilessMatcher: %v", res)
	}
	return nil
}

type SetDefaultMatcherOptions struct {
	CompareBodies       *bool
	ExcludedHeaders     []string
	IgnoredHeaders      []string
	IgnoreQueryOrdering *bool
}

func addDefaults(added []string) []string {
	if added == nil {
		return nil
	}
	added = append(added, ":path", ":authority", ":path", ":scheme")
	return added
}

// SetDefaultMatcher adjusts the "match" operation to exclude the body when matching a request to a recording's entries.
// Pass in `nil` for `t` if you want the bodiless matcher to apply everywhere
func SetDefaultMatcher(t *testing.T, options *SetDefaultMatcherOptions) error {
	if recordMode != PlaybackMode {
		return nil
	}
	if options == nil {
		options = &SetDefaultMatcherOptions{}
	}
	req, err := http.NewRequest("POST", "http://localhost:5000/Admin/SetMatcher", http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header["x-abstraction-identifier"] = []string{"CustomDefaultMatcher"}
	if t != nil {
		req.Header["x-recording-id"] = []string{GetRecordingId(t)}
	}

	marshalled, err := json.MarshalIndent(struct {
		CompareBodies       *bool  `json:"compareBodies,omitempty"`
		ExcludedHeaders     string `json:"excludedHeaders,omitempty"`
		IncludedHeaders     string `json:"includedHeaders,omitempty"`
		IgnoreQueryOrdering *bool  `json:"ignoreQueryOrdering,omitempty"`
	}{
		CompareBodies:       options.CompareBodies,
		ExcludedHeaders:     strings.Join(addDefaults(options.ExcludedHeaders), ","),
		IncludedHeaders:     strings.Join(options.IgnoredHeaders, ","),
		IgnoreQueryOrdering: options.IgnoreQueryOrdering,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))

	return handleProxyResponse(client.Do(req))
}
