//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"io"
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
	f := false
	return SetDefaultMatcher(t, &SetDefaultMatcherOptions{
		CompareBodies: &f,
	})
}

type SetDefaultMatcherOptions struct {
	CompareBodies       *bool
	ExcludedHeaders     []string
	IgnoredHeaders      []string
	IgnoreQueryOrdering *bool
}

func (s *SetDefaultMatcherOptions) fillOptions() {
	f := false
	t := true
	if s == nil {
		s = &SetDefaultMatcherOptions{
			CompareBodies:       &t,
			IgnoreQueryOrdering: &f,
		}
		return
	}

	if s.CompareBodies == nil {
		s.CompareBodies = &t
	}
	if s.IgnoreQueryOrdering == nil {
		s.IgnoreQueryOrdering = &f
	}
}

func addDefaults(added []string) []string {
	if added == nil {
		return nil
	}
	needToAdd := []string{":path", ":authority", ":method", ":scheme"}
	for _, a := range added {
		for idx, n := range needToAdd {
			if a == n {
				needToAdd = append(needToAdd[:idx], needToAdd[idx+1:]...)
			}
		}
	}
	return append(added, needToAdd...)
}

// SetDefaultMatcher adjusts the "match" operation to exclude the body when matching a request to a recording's entries.
// Pass in `nil` for `t` if you want the bodiless matcher to apply everywhere
func SetDefaultMatcher(t *testing.T, options *SetDefaultMatcherOptions) error {
	if recordMode != PlaybackMode {
		return nil
	}
	options.fillOptions()
	req, err := http.NewRequest("POST", "http://localhost:5000/Admin/SetMatcher", http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header["x-abstraction-identifier"] = []string{"CustomDefaultMatcher"}
	if t != nil {
		req.Header[IDHeader] = []string{GetRecordingId(t)}
	}

	if !(*options.CompareBodies) {
		options.ExcludedHeaders = append(options.ExcludedHeaders, "Content-Length")
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

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))

	return handleProxyResponse(client.Do(req))
}
