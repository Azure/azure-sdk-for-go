package autorest

import (
	"fmt"
	"os/exec"
)

func FormatPackage(dir string) error {
	c := exec.Command("gofmt", "-w", "-s", dir)
	b, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(b))
	}
	return nil
}
