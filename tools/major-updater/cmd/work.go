package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	autorestArgsPattern = "--use=@microsoft.azure/autorest.go@~2.1.99 %s --go --multiapi --go-sdk-folder=%s --use-onever"
)

func autorestCommand(file string, sdk string) *exec.Cmd {
	autorestArgs := fmt.Sprintf(autorestArgsPattern, file, sdk)
	c := exec.Command("autorest", strings.Split(autorestArgs, " ")...)
	return c
}
