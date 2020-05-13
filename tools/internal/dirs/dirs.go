// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dirs

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/internal/files"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetSubdirs returns all of the subdirectories under current.
// Returns an empty slice if current contains no subdirectories.
func GetSubdirs(current string) ([]string, error) {
	children, err := ioutil.ReadDir(current)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, info := range children {
		if info.IsDir() {
			dirs = append(dirs, info.Name())
		}
	}
	return dirs, nil
}

// DeleteChildDirs deletes all child directories in the specified directory.
func DeleteChildDirs(dir string) error {
	children, err := GetSubdirs(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, child := range children {
		childPath := filepath.Join(dir, child)
		err = os.RemoveAll(childPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeepCompare compares the two directories to determine if they are identical recursively.
// note: will take a significant long time when invoking on large directories
func DeepCompare(base, target string) (bool, error) {
	identical := true
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		index := strings.Index(path, base)
		if index < 0 {
			return fmt.Errorf("failed to get base directory '%s' in path '%s'", base, path)
		}
		relativePath := path[index + len(base):]
		targetPath := filepath.Join(target, relativePath)
		if exists, err := files.Exists(targetPath); err != nil {
			return err
		} else if !exists {
			identical = false
			return notIdenticalError{}
		}

		return nil
	})

	if _, ok := err.(notIdenticalError); ok {
		return identical, nil
	}
	return identical, err
}

type notIdenticalError struct {
}

func (e notIdenticalError) Error() string {
	return "Not identical"
}
