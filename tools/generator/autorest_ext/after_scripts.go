package autorest_ext

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	ProfileGeneration = "go generate ./profiles"
	ProfileFormat     = "gofmt -w ./profiles"
)

func RegenerateProfiles(sdkRoot string) error {
	if err := executeScript(ProfileGeneration, sdkRoot); err != nil {
		return err
	}

	if err := executeScript(ProfileFormat, sdkRoot); err != nil {
		return err
	}

	return nil
}

func executeScript(script, dir string) error {
	argument := strings.Split(script, " ")
	c := exec.Command(argument[0], argument[1:]...)
	c.Dir = dir
	b, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(b))
	}
	return nil
}
