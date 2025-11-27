// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/automation/pipeline"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	internalutils "github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/spf13/cobra"
)

// Command returns the automation v2 command. Note that this command is designed to run in the root directory of
// azure-sdk-for-go. It does not work if you are running this tool in somewhere else
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "automation-v2 <generate input filepath> <generate output filepath>",
		Args: cobra.RangeArgs(2, 3),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0)          // remove the time stamp prefix
			log.SetOutput(os.Stdout) // set the output to stdout
			cmd.SetErrPrefix("[ERROR]")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := execute(args[0], args[1]); err != nil {
				return errors.New(logError(err))
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("Using current directory as SDK root: %s", cwd)

	ctx := automationContext{
		sdkRoot:    internalutils.NormalizePath(cwd),
		specRoot:   input.SpecFolder,
		commitHash: input.HeadSha,
	}
	output, err := ctx.generate(input)
	if output != nil && len(output.Packages) != 0 {
		log.Printf("Writing output to file '%s'...", outputPath)
		if err := pipeline.WriteOutput(outputPath, output); err != nil {
			return fmt.Errorf("cannot write generate output: %+v", err)
		}
	}
	if err != nil {
		return err
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

	errorBuilder := generateErrorBuilder{}
	if input.RunMode == utils.AutomationRunModeLocal || input.RunMode == utils.AutomationRunModeRelease {
		if input.SdkReleaseType != "" && input.SdkReleaseType != utils.SDKReleaseTypeStable && input.SdkReleaseType != utils.SDKReleaseTypePreview {
			return nil, fmt.Errorf("invalid SDK release type:%s, only support 'stable' or 'beta'", input.SdkReleaseType)
		}
		if input.SdkReleaseType != "" && input.ApiVersion != "" {
			if strings.HasSuffix(input.ApiVersion, "-preview") && input.SdkReleaseType == utils.SDKReleaseTypeStable {
				return nil, fmt.Errorf("SDK release type is stable, but API version: %s is preview", input.ApiVersion)
			}
		}
		if (input.ApiVersion != "" && input.SdkReleaseType == "") || (input.ApiVersion == "" && input.SdkReleaseType != "") {
			return nil, fmt.Errorf("both APIVersion and SDKReleaseType parameters are required for self-serve SDK generation")
		}
	} else {
		// ignore sdk release type and api version for spec-pull-request and batch mode
		input.SdkReleaseType = ""
		input.ApiVersion = ""
	}

	// create sdk repo ref
	sdkRepo, err := repo.OpenSDKRepository(ctx.sdkRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get sdk repo: %+v", err)
	}

	generateCtx := common.GenerateContext{
		SDKPath:        sdkRepo.Root(),
		SDKRepo:        &sdkRepo,
		SpecPath:       ctx.specRoot,
		SpecCommitHash: ctx.commitHash,
		SpecRepoURL:    input.RepoHTTPSURL,
	}
	results := make([]pipeline.PackageResult, 0)

	// process all typespec projects
	typeSpecNamespaceResults := make(map[string]*common.GenerateResult)
	for _, tspProjectFolder := range input.RelatedTypeSpecProjectFolder {
		log.Printf("Start to process typespec project: %s", tspProjectFolder)
		result, err := generateCtx.GenerateFromTypeSpec(filepath.Join(input.SpecFolder, tspProjectFolder, "tspconfig.yaml"), &common.GenerateParam{
			TspClientOptions: []string{"--debug"},
			ApiVersion:       input.ApiVersion,
			SdkReleaseType:   input.SdkReleaseType,
		})
		if err != nil {
			errorBuilder.add(err)
			continue
		}
		typeSpecNamespaceResults[tspProjectFolder] = result
		log.Printf("Finish processing typespec project: %s", tspProjectFolder)
	}
	// process result
	for tspProjectFolder, namespaceResult := range typeSpecNamespaceResults {
		result := processNamespaceResult(generateCtx, namespaceResult)
		result.TypespecProject = []string{tspProjectFolder}
		results = append(results, result)
	}

	// process all autorest projects
	if input.RelatedReadmeMdFile != "" {
		input.RelatedReadmeMdFiles = append(input.RelatedReadmeMdFiles, input.RelatedReadmeMdFile)
	}
	swaggerNamespaceResults := make(map[string][]*common.GenerateResult)
	for _, readme := range input.RelatedReadmeMdFiles {
		log.Printf("Start to process swagger project: %s", readme)
		absReadme, err := filepath.Abs(filepath.Join(ctx.specRoot, readme))
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path for spec path '%s': %+v", ctx.specRoot, err)
		}
		absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")
		generateCtx.SpecReadmeFile = absReadme
		generateCtx.SpecReadmeGoFile = absReadmeGo
		generateCtx.SpecCommitHash = "" // always use local config for swagger automation
		rpMap, err := ctx.getRPMap(absReadmeGo)
		if err != nil {
			errorBuilder.add(err)
			continue
		}
		result, errs := generateCtx.GenerateFromSwagger(rpMap, &common.GenerateParam{
			RemoveTagSet:        true,
			SkipGenerateExample: true,
		})
		if len(errs) > 0 {
			errorBuilder.add(errs...)
			continue
		}
		swaggerNamespaceResults[readme] = result
		log.Printf("Finish processing swagger project: %s", readme)
	}
	// process result
	for readme, namespaceResults := range swaggerNamespaceResults {
		for _, namespaceResult := range namespaceResults {
			result := processNamespaceResult(generateCtx, namespaceResult)
			result.ReadmeMd = []string{readme}
			results = append(results, result)
		}
	}

	return &pipeline.GenerateOutput{
		Packages: results,
	}, errorBuilder.build()
}

func (ctx *automationContext) getRPMap(absReadmeGo string) (map[string][]common.PackageInfo, error) {
	log.Printf("Get all namespaces from readme file")
	rpMap, err := common.ReadV2ModuleNameToGetNamespace(absReadmeGo)
	if err != nil {
		return nil, fmt.Errorf("cannot get rp and namespaces from readme '%s': %+v", absReadmeGo, err)
	}
	return rpMap, nil
}

func processNamespaceResult(generateCtx common.GenerateContext, namespaceResult *common.GenerateResult) pipeline.PackageResult {
	content := namespaceResult.ChangelogMD
	breaking := namespaceResult.Changelog.HasBreakingChanges()
	if namespaceResult.PullRequestLabels == string(utils.FirstGABreakingChangeLabel) || namespaceResult.PullRequestLabels == string(utils.BetaBreakingChangeLabel) {
		// If the PR is first beta or first GA, it is not necessary to report SDK breaking change in spec PR
		breaking = false
	}
	breakingChangeItems := namespaceResult.Changelog.GetBreakingChangeItems()

	srcFolder := filepath.Join(generateCtx.SDKPath, namespaceResult.PackageRelativePath)
	goSourceArtifact := namespaceResult.PackageRelativePath + ".gosource"
	apiViewArtifact := filepath.Join(generateCtx.SDKPath, goSourceArtifact)
	if namespaceResult.ModuleRelativePath != "" {
		srcFolder = filepath.Join(generateCtx.SDKPath, namespaceResult.ModuleRelativePath)
		goSourceArtifact = namespaceResult.ModuleRelativePath + ".gosource"
		apiViewArtifact = filepath.Join(generateCtx.SDKPath, goSourceArtifact)
	}
	err := zipDirectory(srcFolder, apiViewArtifact)
	if err != nil {
		fmt.Println(err)
	}

	return pipeline.PackageResult{
		Version:       namespaceResult.Version,
		PackageName:   namespaceResult.PackageRelativePath,
		Path:          []string{namespaceResult.PackageRelativePath},
		PackageFolder: namespaceResult.PackageRelativePath,
		Changelog: &pipeline.Changelog{
			Content:             &content,
			HasBreakingChange:   &breaking,
			BreakingChangeItems: &breakingChangeItems,
		},
		APIViewArtifact: goSourceArtifact,
		Language:        "Go",
	}
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
	return fmt.Errorf("total %d error(s): \n%s\n%s", len(b.errors), strings.Join(messages, "\n"), `Refer to the detail errors for potential fixes.
If the issue persists, contact the Go language support channel at https://aka.ms/azsdk/go-lang-teams-channel and include this spec pull request.`)
}

func logError(err error) string {
	buidler := strings.Builder{}
	for i, line := range strings.Split(err.Error(), "\n") {
		if l := strings.TrimSpace(line); l != "" {
			if i == 0 {
				buidler.WriteString(fmt.Sprintf("%s\n", l))
				continue
			}
			buidler.WriteString(fmt.Sprintf("[ERROR] %s\n", l))
		}
	}

	return buidler.String()
}

func zipDirectory(srcFolder, dstZip string) error {
	outFile, err := os.Create(dstZip)
	if err != nil {
		return err
	}
	w := zip.NewWriter(outFile)
	srcFolder = strings.TrimSuffix(srcFolder, string(os.PathSeparator))
	err = filepath.Walk(srcFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name, err = filepath.Rel(filepath.Dir(srcFolder), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}
		hw, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(hw, f)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = outFile.Close()
	if err != nil {
		return err
	}
	return nil
}
