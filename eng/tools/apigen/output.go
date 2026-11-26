// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func writeBundleFiles(directory string, document, metadata []byte) (result error) {
	info, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("inspect output directory %q: %w", directory, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("output path %q is not a directory", directory)
	}
	outputs := []struct {
		path      string
		content   []byte
		temporary string
	}{
		{path: filepath.Join(directory, "api.md"), content: document},
		{path: filepath.Join(directory, "api.metadata.yml"), content: metadata},
	}
	defer func() {
		for _, output := range outputs {
			if output.temporary != "" {
				if err := os.Remove(output.temporary); err != nil {
					result = errors.Join(result, fmt.Errorf("remove temporary output %q: %w", output.temporary, err))
				}
			}
		}
	}()

	// Stage both files before replacing either destination.
	for i := range outputs {
		temporary, err := prepareOutputFile(outputs[i].path, outputs[i].content)
		if err != nil {
			return err
		}
		outputs[i].temporary = temporary
	}
	for i := range outputs {
		if err := os.Rename(outputs[i].temporary, outputs[i].path); err != nil {
			return fmt.Errorf("replace output file %q: %w", outputs[i].path, err)
		}
		outputs[i].temporary = ""
	}
	return nil
}

func prepareOutputFile(path string, content []byte) (temporary string, result error) {
	info, err := os.Lstat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return "", fmt.Errorf("output path %q is not a regular file", path)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect output path %q: %w", path, err)
	}

	file, err := os.CreateTemp(filepath.Dir(path), ".apigen-*.tmp")
	if err != nil {
		return "", fmt.Errorf("create temporary output for %q: %w", path, err)
	}
	temporary = file.Name()
	name := temporary
	defer func() {
		if file != nil {
			if err := file.Close(); err != nil {
				result = errors.Join(result, fmt.Errorf("close temporary output: %w", err))
			}
		}
		if result != nil {
			if err := os.Remove(name); err != nil {
				result = errors.Join(result, fmt.Errorf("remove temporary output %q: %w", name, err))
			}
			temporary = ""
		}
	}()

	if _, err := file.Write(content); err != nil {
		return "", fmt.Errorf("write temporary output for %q: %w", path, err)
	}
	if info != nil {
		if err := file.Chmod(info.Mode().Perm()); err != nil {
			return "", fmt.Errorf("set output permissions: %w", err)
		}
	}
	if err := file.Sync(); err != nil {
		return "", fmt.Errorf("sync temporary output: %w", err)
	}
	err = file.Close()
	file = nil
	if err != nil {
		return "", fmt.Errorf("close temporary output: %w", err)
	}
	return temporary, nil
}
