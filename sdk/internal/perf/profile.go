// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func startCPUProfile(enabled bool, path string) (func() error, error) {
	if !enabled {
		return nil, nil
	}
	if path == "" {
		return nil, fmt.Errorf("CPU profile path cannot be empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating CPU profile directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creating CPU profile: %w", err)
	}
	if err = pprof.StartCPUProfile(file); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("starting CPU profile: %w", err)
	}

	return func() error {
		pprof.StopCPUProfile()
		if err := file.Close(); err != nil {
			return fmt.Errorf("closing CPU profile: %w", err)
		}
		return nil
	}, nil
}
