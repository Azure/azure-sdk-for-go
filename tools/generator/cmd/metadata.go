package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/validate"
)

type metadataContext struct {
	sdkRoot string
	readme  string
}

func (ctx metadataContext) processMetadata(metadataOutput string) ([]string, error) {
	// get the metadata
	m := autorest.NewMetadataProcessorFromLocation(metadataOutput)
	metadataMap, err := m.Process()
	if err != nil {
		return nil, err
	}
	var packages []string
	builder := validationErrorBuilder{
		readme: ctx.readme,
	}
	mCtx := validate.MetadataValidateContext{
		Readme:  ctx.readme,
		SDKRoot: ctx.sdkRoot,
		Validators: []validate.MetadataValidateFunc{
			validate.PreviewCheck,
			validate.MgmtCheck,
		},
	}
	for tag, metadata := range metadataMap {
		// first validate the output folder is valid
		if errors := mCtx.Validate(tag, metadata); len(errors) != 0 {
			builder.addMultiple(errors)
			continue
		}
		outputFolder := filepath.Clean(metadata.PackagePath())
		// first format the package
		if err := autorest.FormatPackage(outputFolder); err != nil {
			return nil, err
		}
		// get the package path - which is a relative path to the sdk root
		packagePath, err := filepath.Rel(ctx.sdkRoot, outputFolder)
		if err != nil {
			builder.add(err)
			continue
		}
		packages = append(packages, packagePath)
	}
	return packages, builder.build()
}

type validationErrorBuilder struct {
	readme string
	errors []error
}

func (b *validationErrorBuilder) addMultiple(errors []error) {
	b.errors = append(b.errors, errors...)
}

func (b *validationErrorBuilder) add(err error) {
	b.errors = append(b.errors, err)
}

func (b *validationErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	var messages []string
	for _, e := range b.errors {
		messages = append(messages, e.Error())
	}
	return fmt.Errorf("validation failed in readme '%s' with %d error(s): \n%s", b.readme, len(b.errors), strings.Join(messages, "\n"))
}
