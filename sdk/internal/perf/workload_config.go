// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type runInvocation struct {
	TestName    string
	CLIArgs     []string
	ConfigArgs  []string
	ConfigPath  string
	Workload    string
	UsesConfig  bool
	HasTestName bool
}

type workloadFile struct {
	DefaultWorkload string                      `json:"defaultWorkload"`
	Workloads       map[string]workloadSettings `json:"workloads"`
}

type workloadSettings struct {
	Test       string         `json:"test"`
	Parameters map[string]any `json:"parameters"`
	Args       []string       `json:"args"`
}

func resolveRunInvocation(args []string) (runInvocation, error) {
	invocation := runInvocation{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--config":
			if i+1 >= len(args) {
				return runInvocation{}, fmt.Errorf("missing value for --config")
			}
			invocation.ConfigPath = args[i+1]
			i++
			continue
		case strings.HasPrefix(arg, "--config="):
			invocation.ConfigPath = strings.TrimPrefix(arg, "--config=")
			continue
		case arg == "--workload":
			if i+1 >= len(args) {
				return runInvocation{}, fmt.Errorf("missing value for --workload")
			}
			invocation.Workload = args[i+1]
			i++
			continue
		case strings.HasPrefix(arg, "--workload="):
			invocation.Workload = strings.TrimPrefix(arg, "--workload=")
			continue
		}

		if i == 0 && !strings.HasPrefix(arg, "-") {
			invocation.TestName = arg
			invocation.HasTestName = true
			continue
		}

		invocation.CLIArgs = append(invocation.CLIArgs, arg)
	}

	if invocation.ConfigPath == "" {
		return invocation, nil
	}

	invocation.UsesConfig = true
	cfg, err := loadWorkloadFile(invocation.ConfigPath)
	if err != nil {
		return runInvocation{}, err
	}

	workloadName := invocation.Workload
	if workloadName == "" {
		workloadName = cfg.DefaultWorkload
	}
	if workloadName == "" {
		return runInvocation{}, fmt.Errorf("no workload specified. provide --workload or set defaultWorkload in %s", invocation.ConfigPath)
	}

	workload, ok := cfg.Workloads[workloadName]
	if !ok {
		return runInvocation{}, fmt.Errorf("workload %q not found in %s", workloadName, invocation.ConfigPath)
	}

	if invocation.TestName == "" {
		invocation.TestName = workload.Test
	}
	if invocation.TestName == "" {
		return runInvocation{}, fmt.Errorf("workload %q in %s is missing test", workloadName, invocation.ConfigPath)
	}

	invocation.Workload = workloadName
	invocation.ConfigArgs = workloadToArgs(workload)
	return invocation, nil
}

func loadWorkloadFile(path string) (workloadFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return workloadFile{}, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	var cfg workloadFile
	if err = dec.Decode(&cfg); err != nil {
		return workloadFile{}, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return cfg, nil
}

func workloadToArgs(workload workloadSettings) []string {
	args := append([]string{}, workload.Args...)

	if len(workload.Parameters) == 0 {
		return args
	}

	keys := make([]string, 0, len(workload.Parameters))
	for k := range workload.Parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := parameterValueToString(workload.Parameters[key])
		if value == "" {
			continue
		}
		args = append(args, fmt.Sprintf("--%s=%s", key, value))
	}

	return args
}

func parameterValueToString(v any) string {
	switch tv := v.(type) {
	case nil:
		return ""
	case string:
		return tv
	case bool:
		return strconv.FormatBool(tv)
	case json.Number:
		return tv.String()
	case float64:
		if tv == float64(int64(tv)) {
			return strconv.FormatInt(int64(tv), 10)
		}
		return strconv.FormatFloat(tv, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", tv)
	}
}
