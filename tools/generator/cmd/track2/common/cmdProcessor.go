package common

import (
	"fmt"
	"os/exec"
)

func ExecuteGoGenerate(path string) error {
	cmd := exec.Command("go", "generate")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute go generate '%s': %+v", string(output), err)
	}
	return nil
}
