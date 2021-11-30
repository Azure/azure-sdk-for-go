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
		return fmt.Errorf("failed to execute `go generate` '%s': %+v", string(output), err)
	}
	return nil
}

// execute `pwsh Invoke-MgmtTestgen` command and fetch result
func ExecuteExampleGenerate(path, packagePath string) error {
	cmd := exec.Command("pwsh", "../../../../eng/scripts/Invoke-MgmtTestGen.ps1", "-skipBuild", "-cleanGenerated", "-format", "-tidy", "-generateExample", packagePath)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Invoke-MgmtTestgen` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `pwsh Invoke-MgmtTestgen` '%s': %+v", string(output), err)
	}
	return nil
}

// execute `goimports` command and fetch result
func ExecuteGoimports(path string) error {
	cmd := exec.Command("go", "get", "golang.org/x/tools/cmd/goimports")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go get golang.org/x/tools/cmd/goimports` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go get golang.org/x/tools/cmd/goimports` '%s': %+v", string(output), err)
	}
	cmd = exec.Command("goimports", "-w", ".")
	cmd.Dir = path
	output, err = cmd.CombinedOutput()
	log.Printf("Result of `goimports` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `goimports` '%s': %+v", string(output), err)
	}
	return nil
}
