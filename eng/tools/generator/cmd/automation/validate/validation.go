// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
)

// MetadataValidateContext describes the context needed in validation of the metadata
type MetadataValidateContext struct {
	Readme  string
	SDKRoot string
}

func (ctx *MetadataValidateContext) getRelPackagePath(metadata model.Metadata) (string, error) {
	if err := ctx.rootCheck(metadata); err != nil {
		return "", err
	}
	rel, err := filepath.Rel(ctx.SDKRoot, metadata.PackagePath())
	if err != nil {
		return "", fmt.Errorf("cannot get relative path from output-folder '%s' to the root directory '%s': %+v", metadata.PackagePath(), ctx.SDKRoot, err)
	}
	return utils.NormalizePath(rel), nil
}

func (ctx *MetadataValidateContext) rootCheck(metadata model.Metadata) error {
	r := filepath.Clean(ctx.SDKRoot)
	o := filepath.Clean(metadata.PackagePath())
	if !strings.HasPrefix(o, r) {
		return fmt.Errorf("the output-folder '%s' is not under root directory '%s', the output-folder is either not configured or not correctly configured", metadata.PackagePath(), ctx.SDKRoot)
	}
	return nil
}

// PreviewCheck ensures the output-folder of a preview package is under the preview sub-directory
func (ctx *MetadataValidateContext) PreviewCheck(tag string, metadata model.Metadata) error {
	log.Printf("Executing PreviewCheck...")
	rel, err := ctx.getRelPackagePath(metadata)
	if err != nil {
		return err
	}
	if isPreviewPackage(metadata.SwaggerFiles()) {
		if !previewOutputRegex.MatchString(rel) {
			return fmt.Errorf("the output-folder of a preview package '%s' must be under the `preview` subdirectory", tag)
		}
	} else {
		if previewOutputRegex.MatchString(rel) {
			return fmt.Errorf("the output-folder of a stable package '%s' must NOT be under the `preview` subdirectory", tag)
		}
	}
	return nil
}

// MgmtCheck ensures that the management plane package has the correct output-folder
func (ctx *MetadataValidateContext) MgmtCheck(tag string, metadata model.Metadata) error {
	log.Printf("Executing MgmtCheck...")
	if isMgmtPackage(ctx.Readme) {
		rel, err := ctx.getRelPackagePath(metadata)
		if err != nil {
			return err
		}
		if !mgmtOutputRegex.MatchString(rel) {
			return fmt.Errorf("the output-folder of a management plane package '%s' is expected to have this pattern: 'services/(preview)?/{RPname}/mgmt/{packageVersion}/{namespace}'", tag)
		}
	}
	return nil
}

// NamespaceCheck ensures that the namespace only contains lower case letters, numbers and underscores
func (ctx *MetadataValidateContext) NamespaceCheck(tag string, metadata model.Metadata) error {
	log.Printf("Executing NamespaceCheck...")
	if len(metadata.Namespace()) == 0 {
		return fmt.Errorf("the namespace in readme.go.md cannot be empty")
	}
	if !namespaceRegex.MatchString(metadata.Namespace()) {
		return fmt.Errorf("the namespace can only contain lower case letters, numbers and underscores")
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
	mgmtOutputRegex     = regexp.MustCompile(`/mgmt/`)
	namespaceRegex      = regexp.MustCompile(`^[a-z][a-z0-9_]*[a-z0-9]?$`)
)
