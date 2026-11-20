// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"

	"apigen/internal/loader"
	"apigen/internal/markdown"
	"apigen/internal/report"

	"github.com/spf13/cobra"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	command := newRootCommand()
	if err := command.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(command.ErrOrStderr(), "Error:", err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	var output, repositoryRoot string
	command := &cobra.Command{
		Use:           "apigen <module-directory>",
		Short:         "Generate a Markdown review report of a Go module's public API",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {
			if output == "" || output == "-" {
				return fmt.Errorf("output must be a directory path; stdout output is not supported")
			}
			loaded, err := loader.Load(command.Context(), args[0])
			if err != nil {
				return err
			}
			packages, err := report.ExtractPackages(loaded)
			if err != nil {
				return err
			}
			module := report.ExtractModule(loaded)
			var document bytes.Buffer
			if err := markdown.Write(&document, module, packages); err != nil {
				return err
			}
			metadata, err := reviewMetadata(loaded, repositoryRoot, document.Bytes())
			if err != nil {
				return err
			}
			return writeBundleFiles(output, document.Bytes(), metadata)
		},
	}
	command.Flags().StringVarP(&output, "output", "o", ".", "Output directory for api.md and api.metadata.yml")
	command.Flags().StringVar(&repositoryRoot, "repo-root", "", "Repository root for metadata paths (default: detect Git root)")
	command.AddCommand(newHashCommand())
	return command
}

func newHashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "hash <api.md>",
		Short: "Print the SHA-256 hash of an existing API report for release checks",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			document, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read API report: %w", err)
			}
			_, err = fmt.Fprintln(command.OutOrStdout(), apiHash(document))
			return err
		},
	}
}
