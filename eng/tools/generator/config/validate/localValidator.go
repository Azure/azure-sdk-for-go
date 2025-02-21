// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/ahmetb/go-linq/v3"
)

type localValidator struct {
	specRoot string
}

func (v *localValidator) Validate(cfg config.Config) error {
	var errResult error
	for readme, infoMap := range cfg.Track2Requests {
		if err := v.validateReadmeExistence(readme); err != nil {
			errResult = errors.Join(errResult, err)
			continue // readme file cannot pass validation, we just skip the validations
		}
		// get content of the readme
		contentOfReadme, err := getReadmeContent(v.specRoot, readme)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("cannot get readme.md content: %+v", err))
			continue
		}
		// validate the existence of readme.go.md
		if err := v.validateReadmeExistence(getReadmeGoFromReadme(readme)); err != nil {
			errResult = errors.Join(errResult, err)
			continue // readme.go.md is mandatory
		}
		// get content of the readme.go.md
		contentOfReadmeGo, err := getReadmeContent(v.specRoot, getReadmeGoFromReadme(readme))
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("cannot get readme.go.md content: %+v", err))
			continue
		}
		// get the keys from infoMap, which is the tags
		var tags []string
		linq.From(infoMap).Select(func(item interface{}) interface{} {
			return item.(config.Track2Request).PackageFlag
		}).ToSlice(&tags)
		// check the tags one by one
		if err := validateTagsInReadme(contentOfReadme, readme, tags...); err != nil {
			errResult = errors.Join(errResult, err)
		}
		// check module-name exist
		if err := validateModuleNameInReadmeGo(contentOfReadmeGo, readme); err != nil {
			errResult = errors.Join(errResult, err)
		}
	}
	return errResult
}

func (v *localValidator) validateReadmeExistence(readme string) error {
	full := filepath.Join(v.specRoot, readme)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return fmt.Errorf("readme file %q does not exist", readme)
	}

	return nil
}

func getReadmeContent(specRoot, readme string) ([]byte, error) {
	full := filepath.Join(specRoot, readme)
	return os.ReadFile(full)
}

func findTagInReadme(content []byte, tag string) bool {
	return regexp.MustCompile(fmt.Sprintf(tagDefinedInReadmeRegex, tag)).Match(content)
}

func getReadmeGoFromReadme(readme string) string {
	return strings.ReplaceAll(readme, readmeFilename, goReadmeFilename)
}

func GetRpAndPackageName(content []byte) (string, string) {
	moduleExist := regexp.MustCompile(goReadmeModuleName).Match(content)
	if moduleExist {
		moduleName := regexp.MustCompile(goReadmeModuleName).FindString(string(content))
		s := strings.Split(moduleName, "/")
		if len(s) == 4 {
			return s[2], s[3]
		}
		return "", ""
	}
	return "", ""
}

const (
	tagDefinedInReadmeRegex = `\$\(tag\)\s*==\s*'%s'`
	tagInBatchRegex         = `-\s*tag\s*:\s*`
	readmeFilename          = "readme.md"
	goReadmeFilename        = "readme.go.md"
	goReadmeModuleName      = `module-name: \S*`
)
