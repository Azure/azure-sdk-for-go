// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os/exec"
)

// execute `go generate` command and fetch result
func ExecuteGoGenerate(path string) error {
	cmd := exec.Command("go", "generate")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go generate` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute go generate '%s': %+v", string(output), err)
	}
	return nil
}
