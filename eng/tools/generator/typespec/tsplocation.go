// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

import (
	"os"

	"github.com/goccy/go-yaml"
)

// tsp-location.yaml
type TspLocation struct {
	Directory             string   `yaml:"directory"`
	Commit                string   `yaml:"commit"`
	Repo                  string   `yaml:"repo"`
	AdditionalDirectories []string `yaml:"additionalDirectories"`

	ModuleVersion string `yaml:"module-version"`
}

func ParseTspLocation(tspLocationPath string) (*TspLocation, error) {
	var tl TspLocation

	data, err := os.ReadFile(tspLocationPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &tl)
	if err != nil {
		return nil, err
	}

	return &tl, nil
}
