//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "snippet-generator <base-directory> [strict-mode]",
		Short: `This CLI tool will let you embed your sample code snippets from your test files to your markdown files`,
		Long: `This CLI tool will let you embed your sample code snippets from your test files to your markdown files.

You need to use magic comments "// Snippet: Name" and "// EndSnippet" to define code snippets with names in your go test files,
and then reference them in your markdown files using "Snippet:Name" right after your code block declarations. 
`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := args[0]
			strict := true
			if len(args) > 1 {
				var err error
				strict, err = strconv.ParseBool(args[1])
				if err != nil {
					return err
				}
			}
			return execute(baseDir, strict)
		},
	}

	return rootCmd
}

func execute(dir string, strict bool) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	p := NewDirectoryProcessor(abs, strict)

	return p.Process()
}
