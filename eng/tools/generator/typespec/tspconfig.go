// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type TypeSpecConfig struct {
	Path string

	TypeSpecProjectSchema
}

type TypeSpecProjectSchema struct {
	Parameters *map[string]any      `yaml:"parameters,omitempty"`
	Options    *TypeSpecEmitOptions `yaml:"options,omitempty"`
}

type TypeSpecEmitOptions struct {
	GoConfig *GoEmitterOptions `yaml:"@azure-tools/typespec-go,omitempty"`
}

func ParseTypeSpecConfig(tspconfigPath string) (*TypeSpecConfig, error) {
	tspConfig := TypeSpecConfig{}
	tspConfig.Path = tspconfigPath

	var err error
	var data []byte
	if strings.HasPrefix(tspconfigPath, "http") {
		// http path
		resp, err := http.Get(tspconfigPath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// local path
		data, err = os.ReadFile(tspconfigPath)
		if err != nil {
			return nil, err
		}
	}

	err = yaml.Unmarshal(data, &(tspConfig.TypeSpecProjectSchema))
	if err != nil {
		return nil, err
	}

	_, err = tspConfig.ExistEmitOption()
	if err != nil {
		return nil, err
	}

	if tspConfig.Options.GoConfig.ServiceDir == "" && tspConfig.Parameters != nil {
		if serviceDirRaw, ok := (*tspConfig.Parameters)["service-dir"]; ok {
			if serviceDirMap, ok := serviceDirRaw.(map[string]any); ok {
				if defaultVal, ok := serviceDirMap["default"]; ok {
					if defaultStr, ok := defaultVal.(string); ok {
						tspConfig.Options.GoConfig.ServiceDir = defaultStr
					}
				}
			}
		}
	}

	tspConfig.Options.GoConfig.EmitterOutputDir = replacePlaceholder(tspConfig, tspConfig.Options.GoConfig.EmitterOutputDir)
	tspConfig.Options.GoConfig.Module = replacePlaceholder(tspConfig, tspConfig.Options.GoConfig.Module)
	tspConfig.Options.GoConfig.ContainingModule = replacePlaceholder(tspConfig, tspConfig.Options.GoConfig.ContainingModule)

	if err = tspConfig.Options.GoConfig.Validate(); err != nil {
		return nil, err
	}

	return &tspConfig, err
}

func replacePlaceholder(tspConfig TypeSpecConfig, str string) string {
	result := str
	if strings.Contains(result, "{service-dir}") {
		result = strings.ReplaceAll(result, "{service-dir}", tspConfig.Options.GoConfig.ServiceDir)
	}
	if strings.Contains(result, "{package-dir}") {
		result = strings.ReplaceAll(result, "{package-dir}", tspConfig.Options.GoConfig.PackageDir)
	}
	return result
}

func (tc *TypeSpecConfig) GetPackageRelativePath() string {
	if tc.Options.GoConfig.EmitterOutputDir != "" {
		re := regexp.MustCompile(`\{output-dir\}/`)
		return re.ReplaceAllString(tc.Options.GoConfig.EmitterOutputDir, "")
	} else if tc.Options.GoConfig.PackageDir == "" {
		return tc.GetModuleRelativePath()
	} else {
		return tc.Options.GoConfig.ServiceDir + "/" + tc.Options.GoConfig.PackageDir
	}
}

func (tc *TypeSpecConfig) GetModuleRelativePath() string {

	module := tc.Options.GoConfig.Module
	if tc.Options.GoConfig.ContainingModule != "" {
		module = tc.Options.GoConfig.ContainingModule
	}

	re := regexp.MustCompile(`github\.com/Azure/azure-sdk-for-go/|/v\d+$`)
	return re.ReplaceAllString(module, "")
}

func (tc TypeSpecConfig) ExistEmitOption() (bool, error) {
	if tc.Options == nil {
		return false, fmt.Errorf("no options found in %s", tc.Path)
	}
	if tc.Options.GoConfig == nil {
		return false, fmt.Errorf("Go emitter config not found in %s", tc.Path)
	}
	return true, nil
}

// GetRpAndPackageName return [rpName, packageName]
// module: github.com/Azure/azure-sdk-for-go/sdk/.../{rpName}/{packageName}
func (tc TypeSpecConfig) GetRpAndPackageName() ([2]string, error) {
	module := tc.Options.GoConfig.Module
	if module == "" {
		module = tc.Options.GoConfig.ContainingModule
	}
	return tc.GetRpAndPackageNameByModule(module)
}

// GetRpAndPackageName return [rpName, packageName]
// module: github.com/Azure/azure-sdk-for-go/sdk/.../{rpName}/{packageName}
func (tc TypeSpecConfig) GetRpAndPackageNameByModule(module string) ([2]string, error) {
	s := strings.Split(module, "/")
	l := len(s)
	if l < 2 {
		return [2]string{}, fmt.Errorf("module is invalid")
	}
	if !strings.Contains(s[l-1], "arm") && !strings.Contains(s[l-1], "az") {
		return [2]string{}, fmt.Errorf("packageName is invalid and must start with `arm` or `az`")
	}

	return [2]string{s[l-2], s[l-1]}, nil
}

func ExistGoConfigInTspConfig(tspconfig string) (bool, error) {
	if tspconfig == "" {
		return false, fmt.Errorf("tspconfig path is empty")
	}

	tsc, err := ParseTypeSpecConfig(tspconfig)
	if err != nil {
		return false, err
	}

	return tsc.ExistEmitOption()
}
