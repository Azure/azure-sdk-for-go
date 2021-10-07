// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"fmt"
	"io"
	"os"
)

func ParseConfig(path string) (*Config, error) {
	reader, err := getConfigReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %+v", err)
	}
	cfg, err := FromReader(reader).Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse configs: %+v", err)
	}
	return cfg, nil
}

func getConfigReader(config string) (io.Reader, error) {
	if config == "" {
		return os.Stdin, nil
	}
	return os.Open(config)
}
