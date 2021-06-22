package autorest_ext

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"

	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	sdkutils "github.com/Azure/azure-sdk-for-go/tools/internal/utils"
)

type GenerateContext interface {
	SDKRoot() string
	SpecRoot() string
	RepoContent() map[string]exports.Content
}

type GenerateInput struct {
	// Readme is the relative path of the readme file to the root directory of azure-sdk-for-go
	Readme string
	// Tag is the readme tag to be generated
	Tag string
	// SDKRoot is the root directory of azure-sdk-for-go
	SDKRoot string
	// CommitHash is the head commit hash of azure-rest-api-specs
	CommitHash string
	// Options specifies the options that this generation task will be using
	Options model.Options
}

type GenerateOptions struct {
	// MetadataOutputRoot specifies the root directory of all the metadata goes.
	// Metadata will be generated to a temp directory if not specified.
	// The metadataOutput directory will not be removed after the generation succeeded
	MetadataOutputRoot string
	// Stderr ...
	Stderr io.Writer
	// Stdout ...
	Stdout io.Writer
}

type GenerateResult struct {
	MetadataOutputRoot string
	Metadata           autorest.GenerationMetadata
	Package            autorest.ChangelogResult
}

func GeneratePackage(ctx GenerateContext, input GenerateInput, options GenerateOptions) (*GenerateResult, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	absReadme := filepath.Join(ctx.SpecRoot(), input.Readme)
	metadataOutput := filepath.Join(options.MetadataOutputRoot, input.Tag)
	g := autorest.NewGeneratorFromOptions(input.Options).WithTag(input.Tag).WithMetadataOutput(metadataOutput).WithReadme(absReadme)

	// generate
	if err := generate(g, options.Stdout, options.Stderr); err != nil {
		return nil, fmt.Errorf("failed to execute autorest: %+v", err)
	}
	// write the changelog and metadata file
	result, metadata, err := changelogAndMetadata(ctx, input, metadataOutput, g.Arguments())
	if err != nil {
		return nil, err
	}

	return &GenerateResult{
		MetadataOutputRoot: options.MetadataOutputRoot,
		Metadata:           *metadata,
		Package:            *result,
	}, nil
}

func generate(generator *autorest.Generator, stdout, stderr io.Writer) error {
	stdoutPipe, _ := generator.StdoutPipe()
	stderrPipe, _ := generator.StderrPipe()
	defer stdoutPipe.Close()
	defer stderrPipe.Close()
	var arguments []string
	linq.From(generator.Arguments()).Select(func(item interface{}) interface{} {
		return item.(model.Option).Format()
	}).ToSlice(&arguments)
	log.Printf("Generation parameters: %s", strings.Join(arguments, ", "))
	_ = generator.Start()
	// we put all the output from autorest to stderr since those are logs in order not to interrupt the proper output of the release command
	go utils.ScannerPrint(bufio.NewScanner(stdoutPipe), stdout)
	go utils.ScannerPrint(bufio.NewScanner(stderrPipe), stderr)
	return generator.Wait()
}

func changelogAndMetadata(ctx GenerateContext, input GenerateInput, metadataOutput string, argument []model.Option) (*autorest.ChangelogResult, *autorest.GenerationMetadata, error) {
	result, err := changelog(ctx, metadataOutput)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write changelog file: %+v", err)
	}

	// write the metadata file
	metadata, err := metadata(input, *result, argument)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write metadata file: %+v", err)
	}

	return result, metadata, nil
}

func changelog(ctx GenerateContext, metadataOutput string) (*autorest.ChangelogResult, error) {
	// parse the metadata from autorest
	metadataMap, err := autorest.NewMetadataProcessorFromLocation(metadataOutput).Process()
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata in '%s': %+v", metadataOutput, err)
	}
	// process the changelog
	changelogResults, err := autorest.NewChangelogProcessorFromContext(ctx).Process(metadataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to process the changelog: %+v", err)
	}
	// we should only have one changelog
	if len(changelogResults) != 1 {
		return nil, fmt.Errorf("expecting 1 changelog result, but got %d", len(changelogResults))
	}

	changelogPath, err := sdk.WriteChangelogFile(changelogResults[0])
	if err != nil {
		return nil, fmt.Errorf("failed to write changelog file: %+v", err)
	}
	log.Printf("changelog file writes to '%s'", changelogPath)
	return &changelogResults[0], nil
}

func metadata(input GenerateInput, result autorest.ChangelogResult, arguments []model.Option) (*autorest.GenerationMetadata, error) {
	metadata := getMetadata(input, result, arguments)
	metadataPath, err := sdk.WriteMetadataFile(result.PackageFullPath, metadata)
	if err != nil {
		return nil, err
	}
	log.Printf("metadata file writes to '%s'", metadataPath)
	return &metadata, nil
}

func getMetadata(input GenerateInput, result autorest.ChangelogResult, arguments []model.Option) autorest.GenerationMetadata {
	options := autorest.AdditionalOptionsToString(arguments)
	codeGenVersion := input.Options.CodeGeneratorVersion()
	return autorest.GenerationMetadata{
		CommitHash:     input.CommitHash,
		Readme:         autorest.NormalizedSpecRoot + sdkutils.NormalizePath(input.Readme),
		Tag:            input.Tag,
		CodeGenVersion: codeGenVersion,
		RepositoryURL:  "https://github.com/Azure/azure-rest-api-specs.git",
		AutorestCommand: fmt.Sprintf("autorest --use=%s --tag=%s --go-sdk-folder=/_/azure-sdk-for-go %s /_/azure-rest-api-specs/%s",
			codeGenVersion, result.Tag, strings.Join(options, " "), sdkutils.NormalizePath(input.Readme)),
		AdditionalProperties: autorest.GenerationMetadataAdditionalProperties{
			AdditionalOptions: strings.Join(options, " "),
		},
	}
}

func (input GenerateInput) validate() error {
	if input.Readme == "" {
		return fmt.Errorf("`Readme` cannot be empty in input")
	}
	if filepath.IsAbs(input.Readme) {
		return fmt.Errorf("`Readme` must be a relative path")
	}
	if input.Tag == "" {
		return fmt.Errorf("`Tag` cannot be empty in input")
	}
	if input.Options == nil {
		return fmt.Errorf("`Options` cannot be nil")
	}
	return nil
}

func (options *GenerateOptions) validate() error {
	if options.MetadataOutputRoot == "" {
		options.MetadataOutputRoot = filepath.Join(os.TempDir(), fmt.Sprintf("generation-metadata-%v", time.Now().Unix()))
	}
	return nil
}
