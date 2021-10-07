// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// ReadBatchTags reads from a io.Reader of readme.go.md, parses the `multiapi` section and produces a slice of tags
func ReadBatchTags(reader io.Reader) ([]string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	start := -1
	end := -1
	for i, line := range lines {
		if multiAPISectionBeginRegex.MatchString(line) {
			if start >= 0 {
				return nil, fmt.Errorf("multiple multiapi section found on line %d and %d, we should only have one", start+1, i+1)
			}
			start = i
		}
		if start >= 0 && end < 0 && multiAPISectionEndRegex.MatchString(line) {
			end = i
		}
	}

	if start < 0 {
		return nil, fmt.Errorf("cannot find multiapi section")
	}
	if end < 0 {
		return nil, fmt.Errorf("multiapi section does not properly end")
	}

	// get the content of the mutliapi section
	multiAPISection := lines[start+1 : end]
	// the multiapi section should at least have two lines
	if len(multiAPISection) < 2 {
		return nil, fmt.Errorf("multiapi section cannot be parsed")
	}
	// verify the first line of the section should be "batch:"
	if !batchRegex.MatchString(multiAPISection[0]) {
		return nil, fmt.Errorf("multiapi section should begin with `batch:`")
	}
	// iterate over the rest lines, should one for one tag
	var tags []string
	for _, line := range multiAPISection[1:] {
		matches := tagRegex.FindStringSubmatch(line)
		if len(matches) < 2 {
			log.Printf("[WARNING] line in batch '%s' does not starts with 'tags', ignore", line)
			continue
		}
		tags = append(tags, matches[1])
	}
	return tags, nil
}

var (
	multiAPISectionBeginRegex = regexp.MustCompile("^```\\s*yaml\\s*\\$\\(go\\)\\s*&&\\s*\\$\\(multiapi\\)")
	multiAPISectionEndRegex   = regexp.MustCompile("^\\s*```\\s*$")
	batchRegex                = regexp.MustCompile(`^\s*batch:\s*$`)
	tagRegex                  = regexp.MustCompile(`^\s*-\s+tag: (\S+)\s*$`)
)
