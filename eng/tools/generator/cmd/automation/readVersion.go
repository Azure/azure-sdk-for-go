// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

// ReadVersion reads the version of azure-sdk-for-go
func ReadVersion(path string) (string, error) {
	c, err := exports.Get(path)
	if err != nil {
		return "", err
	}

	version := ""
	for k, v := range c.Consts {
		if k == "Number" {
			version = v.Value
		}
	}

	if version == "" {
		return "", fmt.Errorf("cannot get version number in package '%s'", path)
	}

	return version, nil
}
