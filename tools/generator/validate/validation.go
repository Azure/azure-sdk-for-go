package validate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
)

// MetadataValidateFunc is a function that validates a metadata is legal or not
type MetadataValidateFunc func(ctx *MetadataValidateContext, metadata model.Metadata) error

// MetadataValidateContext describes the context needed in validation of the metadata
type MetadataValidateContext struct {
	Readme     string
	SDKRoot    string
	Validators []MetadataValidateFunc
}

// Validate validates the metadata
func (ctx *MetadataValidateContext) Validate(metadata model.Metadata) error {
	builder := metadataValidationErrorBuilder{}
	for _, validator := range ctx.Validators {
		err := validator(ctx, metadata)
		builder.add(err)
	}
	return builder.build()
}

type metadataValidationErrorBuilder struct {
	errors []error
}

func (b *metadataValidationErrorBuilder) add(err error) {
	if err != nil {
		b.errors = append(b.errors, err)
	}
}

func (b *metadataValidationErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	var messages []string
	for _, e := range b.errors {
		messages = append(messages, e.Error())
	}
	return &MetadataValidationError{
		errors:  b.errors,
		message: strings.Join(messages, "\n"),
	}
}

// MetadataValidationError ...
type MetadataValidationError struct {
	errors  []error
	message string
}

// Error ...
func (e *MetadataValidationError) Error() string {
	return fmt.Sprintf("metadata validation failed with %d errors: \n%s", len(e.errors), e.message)
}

func (ctx *MetadataValidateContext) getRelPackagePath(pkgPath string) (string, error) {
	rel, err := filepath.Rel(ctx.SDKRoot, pkgPath)
	if err != nil {
		return "", fmt.Errorf("cannot get relative path from output-folder '%s' to the root directory '%s': %+v", pkgPath, ctx.SDKRoot, err)
	}
	return utils.NormalizePath(rel), nil
}

func rootCheck(ctx *MetadataValidateContext, metadata model.Metadata) error {
	r := filepath.Clean(ctx.SDKRoot)
	o := filepath.Clean(metadata.PackagePath())
	if !strings.HasPrefix(o, r) {
		return fmt.Errorf("the output-folder '%s' is not under root directory '%s', the output-folder is either not configured or not correctly configured", metadata.PackagePath(), ctx.SDKRoot)
	}
	return nil
}

// PreviewCheck ensures the output-folder of a preview package is under the preview sub-directory
func PreviewCheck(ctx *MetadataValidateContext, metadata model.Metadata) error {
	if err := rootCheck(ctx, metadata); err != nil {
		return err
	}
	if isPreviewPackage(metadata.SwaggerFiles()) {
		rel, err := ctx.getRelPackagePath(metadata.PackagePath())
		if err != nil {
			return err
		}
		if !previewOutputRegex.MatchString(rel) {
			return fmt.Errorf("the output-folder of a preview package must be under the `services/preview` subdirectory")
		}
	}
	return nil
}

// MgmtCheck ensures that the management plane package has the correct output-folder
func MgmtCheck(ctx *MetadataValidateContext, metadata model.Metadata) error {
	if isMgmtPackage(ctx.Readme) {
		rel, err := ctx.getRelPackagePath(metadata.PackagePath())
		if err != nil {
			return err
		}
		if !mgmtOutputRegex.MatchString(rel) {
			return fmt.Errorf("the output-folder of a management plane package must be in this pattern: '^services(/preview)?/[^/]+/mgmt/[^/]+/[^/]+$'")
		}
	}
	return nil
}

func isPreviewPackage(inputFiles []string) bool {
	for _, inputFile := range inputFiles {
		if isPreviewSwagger(inputFile) {
			return true
		}
	}
	return false
}

func isPreviewSwagger(inputFile string) bool {
	return previewSwaggerRegex.MatchString(inputFile)
}

func isMgmtPackage(readme string) bool {
	return mgmtReadmeRegex.MatchString(readme)
}

var (
	previewSwaggerRegex = regexp.MustCompile(`^preview|.+[/\\]preview[/\\]`)
	previewOutputRegex  = regexp.MustCompile(`^services/preview/`)
	mgmtReadmeRegex     = regexp.MustCompile(`[/\\]resource-manager[/\\]`)
	mgmtOutputRegex     = regexp.MustCompile(`^services(/preview)?/[^/]+/mgmt/[^/]+/[^/]+$`)
)
