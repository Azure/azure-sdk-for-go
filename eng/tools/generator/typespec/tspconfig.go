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

/*
tspconfig schema: https://typespec.io/docs/handbook/configuration#schema
*/
type TypeSpecConfig struct {
	Path string

	TypeSpecProjectSchema
}

// https://typespec.io/docs/handbook/configuration#schema
type TypeSpecProjectSchema struct {
	Extends              string         `yaml:"extends,omitempty"`
	Parameters           map[string]any `yaml:"parameters,omitempty"`
	EnvironmentVariables map[string]any `yaml:"environment-variables,omitempty"`
	WarnAsError          bool           `yaml:"warn-as-error,omitempty"`
	OutPutDir            string         `yaml:"output-dir,omitempty"`
	Trace                []string       `yaml:"trace,omitempty"`
	Imports              string         `yaml:"imports,omitempty"`
	Emit                 []string       `yaml:"emit,omitempty"`
	Options              map[string]any `yaml:"options,omitempty"`
	Linter               LinterConfig   `yaml:"linter,omitempty"`
}

// <library name>:<rule/ruleset name>
type LinterConfig struct {
	Extends []RuleRef          `yaml:"extends,omitempty"`
	Enable  map[RuleRef]bool   `yaml:"enable,omitempty"`
	Disable map[RuleRef]string `yaml:"disable,omitempty"`
}

type RuleRef string

func (r RuleRef) Validate() bool {
	return regexp.MustCompile(`.*/.*`).MatchString(string(r))
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

	// replace {service-dir} and {package-dir}
	exist := tspConfig.ExistEmitOption(string(TypeSpec_GO))
	if exist {
		emitOption, err := tspConfig.EmitOption(string(TypeSpec_GO))
		if err != nil {
			return nil, err
		}

		goOption := emitOption.(map[string]any)
		module, ok := goOption["module"].(string)
		if !ok {
			return nil, fmt.Errorf("the module must be set in %s option", TypeSpec_GO)
		}

		if strings.Contains(module, "{service-dir}") {
			module = strings.ReplaceAll(module, "{service-dir}", goOption["service-dir"].(string))
		}

		if strings.Contains(module, "{package-dir}") {
			module = strings.ReplaceAll(module, "{package-dir}", goOption["package-dir"].(string))
		}

		goOption["module"] = module
		tspConfig.EditOptions(string(TypeSpec_GO), goOption, false)

		typespecGoOption, err := NewGoEmitterOptions(goOption)
		if err != nil {
			return nil, err
		}
		if err = typespecGoOption.Validate(); err != nil {
			return nil, err
		}
	}

	return &tspConfig, err
}

func (tc *TypeSpecConfig) EditOptions(emit string, option map[string]any, append bool) {
	if tc.Options == nil {
		tc.Options = make(map[string]any)
	}

	if _, ok := tc.Options[emit]; ok {
		if append {
			op1 := tc.Options[emit].(map[string]any)
			for k, v := range option {
				op1[k] = v
			}
			tc.Options[emit] = op1
		} else {
			tc.Options[emit] = option
		}
	} else {
		tc.Options[emit] = option
	}
}

func (tc *TypeSpecConfig) Save() error {
	data, err := yaml.MarshalWithOptions(tc.TypeSpecProjectSchema, yaml.IndentSequence(true))
	if err != nil {
		return err
	}

	return os.WriteFile(tc.Path, data, 0666)
}

func (tc *TypeSpecConfig) EmitOption(emit string) (any, error) {
	if tc.Options == nil {
		return nil, fmt.Errorf("no options found in %s", tc.Path)
	}

	if v, ok := tc.Options[emit]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("emit %s option not found in %s", emit, tc.Path)
}

func (tc TypeSpecConfig) ExistEmitOption(emit string) bool {
	_, err := tc.EmitOption(emit)
	return err == nil
}

// GetModuleName return [rpName, packageName]
// module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/{rpName}/{packageName}
func (tc TypeSpecConfig) GetModuleName() ([2]string, error) {
	option, err := tc.EmitOption(string(TypeSpec_GO))
	if err != nil {
		return [2]string{}, err
	}

	module := (option.(map[string]any))["module"].(string)
	s := strings.Split(module, "/")
	l := len(s)
	if l != 7 {
		return [2]string{}, fmt.Errorf("module is invalid and must be in the format of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/{rpName}/{packageName}`")
	}
	if !strings.Contains(s[l-1], "arm") {
		return [2]string{}, fmt.Errorf("packageName is invalid and must start with `arm`")
	}

	return [2]string{s[l-2], s[l-1]}, nil
}

func TspConfigExistEmitOption(tspconfig string, emit string) (bool, error) {
	if tspconfig == "" {
		return false, fmt.Errorf("tspconfig path is empty")
	}

	tsc, err := ParseTypeSpecConfig(tspconfig)
	if err != nil {
		return false, err
	}

	return tsc.ExistEmitOption(string(TypeSpec_GO)), nil
}
