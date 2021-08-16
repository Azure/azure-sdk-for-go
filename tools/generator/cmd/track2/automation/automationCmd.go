// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/automation/pipeline"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/track2/common"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
	"github.com/spf13/cobra"
)

// Command returns the automation for track2 command. Note that this command is designed to run in the root directory of
// azure-sdk-for-go. It does not work if you are running this tool in somewhere else
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "automation-track2 <generate input filepath> <generate output filepath>",
		Args: cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := execute(args[0], args[1]); err != nil {
				logError(err)
				return err
			}
			return nil
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	return cmd
}

func execute(inputPath, outputPath string) error {
	log.Printf("Reading generate input file from '%s'...", inputPath)
	input, err := pipeline.ReadInput(inputPath)
	if err != nil {
		return fmt.Errorf("cannot read generate input: %+v", err)
	}
	log.Printf("Generating using the following GenerateInput...\n%s", input.String())
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("Using current directory as SDK root: %s", cwd)

	ctx := automationContext{
		sdkRoot:    utils.NormalizePath(cwd),
		specRoot:   input.SpecFolder,
		commitHash: input.HeadSha,
	}
	output, err := ctx.generate(input)
	if err != nil {
		return err
	}
	log.Printf("Output generated: \n%s", output.String())
	log.Printf("Writing output to file '%s'...", outputPath)
	if err := pipeline.WriteOutput(outputPath, output); err != nil {
		return fmt.Errorf("cannot write generate output: %+v", err)
	}
	return nil
}

type automationContext struct {
	sdkRoot    string
	specRoot   string
	commitHash string
}

// TODO -- support dry run
func (ctx *automationContext) generate(input *pipeline.GenerateInput) (*pipeline.GenerateOutput, error) {
	if input.DryRun {
		return nil, fmt.Errorf("dry run not supported yet")
	}

	// iterate over all the readme
	results := make([]pipeline.PackageResult, 0)
	errorBuilder := generateErrorBuilder{}
	for _, readme := range input.RelatedReadmeMdFiles {
		log.Printf("Start to process readme file: %s", readme)
		generateCtx := common.GenerateContext{
			SdkPath:    ctx.sdkRoot,
			SpecPath:   ctx.specRoot,
			CommitHash: ctx.commitHash,
		}

		namespaceResults, errors := generateCtx.GenerateForAutomation(readme)
		if len(errors) != 0 {
			errorBuilder.add(errors...)
			continue
		}

		for _, namespaceResult := range namespaceResults {
			content := namespaceResult.ChangelogMd
			breaking := namespaceResult.Changelog.HasBreakingChanges()
			breakingChangeItems := namespaceResult.Changelog.GetBreakingChangeItems()

			results = append(results, pipeline.PackageResult{
				Version:     namespaceResult.Version,
				PackageName: namespaceResult.PackageName,
				Path:        []string{fmt.Sprintf("sdk/%s/%s", namespaceResult.RpName, namespaceResult.PackageName)},
				ReadmeMd:    []string{readme},
				Changelog: &pipeline.Changelog{
					Content:             &content,
					HasBreakingChange:   &breaking,
					BreakingChangeItems: &breakingChangeItems,
				},
			})
		}
		log.Printf("Finish to process readme file: %s", readme)
	}

	return &pipeline.GenerateOutput{
		Packages: results,
	}, errorBuilder.build()
}

type generateErrorBuilder struct {
	errors []error
}

func (b *generateErrorBuilder) add(err ...error) {
	b.errors = append(b.errors, err...)
}

func (b *generateErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	var messages []string
	for _, err := range b.errors {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("total %d error(s): \n%s", len(b.errors), strings.Join(messages, "\n"))
}

func logError(err error) {
	for _, line := range strings.Split(err.Error(), "\n") {
		if l := strings.TrimSpace(line); l != "" {
			log.Printf("[ERROR] %s", l)
		}
	}
}
