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
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/tools/apidiff/repo"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <package searching dir>",
	Short: "Generates a report of all the export content for every package under the specified directory.",
	Long: `The export command generates a report of all the export content for all of the packages under the directory 
specified in <package searching dir>.

The export content will be categorized by consts, funcs, interfaces, structs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return exportCommand(args[0])
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportCommand(searchingDir string) error {
	wt, err := repo.Get(searchingDir)
	if err != nil {
		return fmt.Errorf("failed to get repo: %+v", err)
	}
	r, err := getRepoContent(&wt, searchingDir)
	if err != nil {
		return err
	}
	return r.Print(os.Stdout)
}
