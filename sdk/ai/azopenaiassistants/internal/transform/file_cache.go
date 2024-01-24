//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import "os"

type FileCache struct {
	files map[string]string
}

func NewFileCache() *FileCache {
	return &FileCache{files: map[string]string{}}
}

func (fc *FileCache) LoadFile(fileName string) (string, error) {
	if fc.files[fileName] == "" {
		buff, err := os.ReadFile(fileName)

		if err != nil {
			return "", err
		}

		fc.files[fileName] = string(buff)
	}

	return fc.files[fileName], nil
}

func (fc *FileCache) UpdateFile(fileName string, text string) {
	fc.files[fileName] = text
}

func (fc *FileCache) WriteAll() error {
	for name, contents := range fc.files {
		if err := os.WriteFile(name, []byte(contents), 0500); err != nil {
			return err
		}
	}
	return nil
}
