// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"os/exec"
)

// FormatPackage formats the given package using gofmt
func FormatPackage(dir string) error {
	c := exec.Command("gofmt", "-w", "-s", dir)
	b, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(b))
	}
	return nil
}
