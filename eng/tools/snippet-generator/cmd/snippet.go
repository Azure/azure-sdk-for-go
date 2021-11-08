//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"math"
	"strings"
)

type Snippet struct {
	Name       string
	IsUsed     bool
	SourceText []string
	FilePath   string
}

func NewSnippet(name string, content []string, filepath string) *Snippet {
	return &Snippet{
		Name:       name,
		IsUsed:     false,
		SourceText: trimCommonIndent(trimEmptyLines(content)),
		FilePath:   filepath,
	}
}

func trimEmptyLines(content []string) []string {
	// remove leading empty lines
	start := 0
	end := len(content)-1
	for i := 0; i < len(content); i++ {
		if strings.TrimSpace(content[i]) != "" {
			start = i
			break
		}
	}

	for i := len(content)-1; i >= 0; i-- {
		if strings.TrimSpace(content[i]) != "" {
			end = i
			break
		}
	}

	return content[start:end+1]
}

func trimCommonIndent(content []string) []string {
	if len(content) == 0 {
		return content
	}
	if len(content) == 1 {
		return []string{strings.TrimSpace(replaceTab(content[0]))}
	}

	// find the maximum count of spaces
	max := math.MaxUint32
	for i := 0; i < len(content)-1; i++ {
		common := commonIndent(replaceTab(content[i]), replaceTab(content[i+1]))
		if common < max {
			max = common
		}
	}

	// iterate the array to trim spaces with the count of max
	var result []string
	prefixToTrim := strings.Repeat(" ", max)
	for _, line := range content {
		result = append(result, strings.TrimPrefix(replaceTab(line), prefixToTrim))
	}

	return result
}

func replaceTab(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

func commonIndent(left, right string) int {
	leftCount := countLeadingSpaces(left)
	rightCount := countLeadingSpaces(right)
	if leftCount < rightCount {
		return leftCount
	}
	return rightCount
}

func countLeadingSpaces(s string) int {
	if s == "" {
		return math.MaxUint32 // use this to ignore empty lines
	}
	return len(s) - len(strings.TrimLeft(s, " "))
}
