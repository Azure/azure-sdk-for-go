// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package readme

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Command returns the generate-go-readme command.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-go-readme <rp readme filepath>",
		Short: "Generate a go readme file or add go track2 part to go readme file according to base swagger readme file",
		Args:  cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := execute(args[0]); err != nil {
				logError(err)
				return err
			}
			return nil
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	return cmd
}

func execute(rpReadmeFilepath string) error {
	if _, err := os.Stat(rpReadmeFilepath); errors.Is(err, os.ErrNotExist) {
		return err
	}
	basePath := filepath.Dir(rpReadmeFilepath)
	rpName := filepath.Base(filepath.Dir(basePath))

	readmeGoFile := filepath.Join(basePath, "readme.go.md")
	content := ""
	if _, err := os.Stat(readmeGoFile); err == nil {
		b, err := os.ReadFile(readmeGoFile)
		if err != nil {
			return err
		}
		content = string(b)
	}

	if !strings.Contains(content, "$(go) && $(track2)") {
		content += fmt.Sprintf("\n\n``` yaml $(go) && $(track2)\nlicense-header: MICROSOFT_MIT_NO_VERSION\nmodule-name: sdk/resourcemanager/%s/arm%s\nmodule: github.com/Azure/azure-sdk-for-go/$(module-name)\noutput-folder: $(go-sdk-folder)/$(module-name)\nazure-arm: true\n```\n", rpName, rpName)
		if err := os.WriteFile(readmeGoFile, []byte(content), 0644); err != nil {
			return err
		}
	}

	log.Printf("Succeed to generate readme.go.me file from '%s'...", rpReadmeFilepath)

	return nil
}

func logError(err error) {
	for _, line := range strings.Split(err.Error(), "\n") {
		if l := strings.TrimSpace(line); l != "" {
			log.Printf("[ERROR] %s", l)
		}
	}
}
