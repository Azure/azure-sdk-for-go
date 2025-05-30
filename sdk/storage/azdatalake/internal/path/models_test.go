//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatRenameOptions(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedSource string
	}{
		{
			name:           "Simple path",
			path:           "dir1/file1.txt",
			expectedSource: "dir1/file1.txt",
		},
		{
			name:           "Path with spaces",
			path:           "dir1/file with spaces.txt",
			expectedSource: "dir1/file+with+spaces.txt",
		},
		{
			name:           "Path with unicode characters",
			path:           "dir1/lör 006.jpg",
			expectedSource: "dir1/l%C3%B6r+006.jpg",
		},
		{
			name:           "Path with special characters",
			path:           "dir1/file+name@!&%.txt",
			expectedSource: "dir1/file%2Bname%40%21%26%25.txt",
		},
		{
			name:           "Path with query parameters",
			path:           "dir1/file1.txt?param1=value1&param2=value2",
			expectedSource: "dir1/file1.txt?param1=value1&param2=value2",
		},
		{
			name:           "Path with special characters and query parameters",
			path:           "dir1/file name+special@!&%.txt?param=value with spaces",
			expectedSource: "dir1/file+name%2Bspecial%40%21%26%25.txt?param=value+with+spaces",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, _, createOpts, _ := FormatRenameOptions(nil, test.path)
			assert.NotNil(t, createOpts)
			assert.NotNil(t, createOpts.RenameSource)
			assert.Equal(t, test.expectedSource, *createOpts.RenameSource)
		})
	}
}