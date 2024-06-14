package typespec

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
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

	return &tspConfig, err
}

func (tc *TypeSpecConfig) EditEmit(emits []string) {
	if tc.Emit == nil {
		tc.Emit = emits
		return
	}

	tc.Emit = append(tc.Emit, emits...)
	tc.Emit = slices.Compact(tc.Emit)
}

func (tc *TypeSpecConfig) OnlyEmit(emit string) {
	tc.Emit = []string{emit}
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

func (tc *TypeSpecConfig) Write() error {
	data, err := yaml.MarshalWithOptions(tc.TypeSpecProjectSchema, yaml.IndentSequence(true))
	if err != nil {
		return err
	}

	return os.WriteFile(tc.Path, data, 0666)
}

func (tc *TypeSpecConfig) WriteTo(w io.Writer) (int64, error) {
	data, err := yaml.MarshalWithOptions(tc.TypeSpecProjectSchema, yaml.IndentSequence(true))
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	return int64(n), err
}

func (tc *TypeSpecConfig) EmitOption(emit string) (any, error) {
	if v, ok := tc.Options[emit]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("emit %s not found in %s", emit, tc.Path)
}

func (tc TypeSpecConfig) ExistEmitOption(emit string) bool {
	_, ok := tc.Options[emit]
	return ok
}
