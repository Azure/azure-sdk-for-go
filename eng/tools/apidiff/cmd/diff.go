// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <base export filepath> <target export filepath>",
	Short: "Generate a diff report between the two export report files.",
	Long:  `The diff command consumes two JSON files with the export reports. The command generates a diff report between them.`,
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
	var baseExport, targetExport RepoContent
	if err := json.Unmarshal(base, &baseExport); err != nil {
		return fmt.Errorf("failed to unmarshal base export: %+v", err)
	}
	if err := json.Unmarshal(target, &targetExport); err != nil {
		return fmt.Errorf("failed to unmarshal target export: %+v", err)
	}

	r := getPkgsReport(baseExport, targetExport)

	if asMarkdown {
		fmt.Println(r.ToMarkdown(""))
	} else {
		b, _ := json.Marshal(r)
		fmt.Println(string(b))
	}
	return nil
}
