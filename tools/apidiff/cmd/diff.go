// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var diffCmd = &cobra.Command{
	Use:   "diff <base export filepath> <target export filepath>",
	Short: "Generate a diff report between the two export report files",
	Long:  `The diff command consumes two JSON files with the export reports, and generates a diff report between them.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return diffCommand(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().BoolVarP(&asMarkdown, "markdown", "m", false, "emits the report in markdown format")
}

func diffCommand(basePath, targetPath string) error {
	base, err := ioutil.ReadFile(basePath)
	if err != nil {
		return fmt.Errorf("failed to read base export file %s: %+v", basePath, err)
	}
	target, err := ioutil.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target export file %s: %+v", targetPath, err)
	}
	var baseExport, targetExport repoContent
	if err := json.Unmarshal(base, &baseExport); err != nil {
		return fmt.Errorf("failed to unmarshal base export: %+v", err)
	}
	if err := json.Unmarshal(target, &targetExport); err != nil {
		return fmt.Errorf("failed to unmarshal target export: %+v", err)
	}

	r := getPkgsReport(baseExport, targetExport)

	if asMarkdown {
		fmt.Println(r.toMarkdown())
	} else {
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	}
	return nil
}
