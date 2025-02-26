// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type replacement struct {
	regex   *regexp.Regexp
	replace string
}

type replacer struct {
	// replacements maps file paths to replacements to make in those files
	replacements map[string][]replacement
}

func must[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (r *replacer) Replace(paths []string, regex, replace string) {
	if r.replacements == nil {
		r.replacements = make(map[string][]replacement)
	}
	for _, p := range paths {
		p = must(filepath.Abs(p))
		r.replacements[p] = append(r.replacements[p], replacement{
			regex:   regexp.MustCompile(regex),
			replace: replace,
		})
	}
}

func (r *replacer) Do() error {
	for path, tasks := range r.replacements {
		f, err := os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		for _, t := range tasks {
			after := t.regex.ReplaceAll(b, []byte(t.replace))
			if bytes.Equal(b, after) {
				log.Printf(`replacement "%s -> %s" had no effect in %s`, t.regex, t.replace, filepath.Base(path))
			}
			b = after
		}
		if err = os.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	r := replacer{}
	for before, after := range map[string]string{
		"OcpBatchFileIsdirectory": "OCPBatchFileIsDirectory",
		"OcpBatchFileMode":        "OCPBatchFileMode",
		"OcpBatchFileURL":         "OCPBatchFileURL",
		"OcpCreationTime":         "OCPCreationTime",
	} {
		r.Replace([]string{"client.go", "responses.go"}, before, after)
	}
	for before, after := range map[string]string{
		"Ocpdate":  "OCPDate",
		"OcpRange": "OCPRange",
	} {
		r.Replace([]string{"client.go", "options.go"}, before, after)
	}
	if err := r.Do(); err != nil {
		log.Fatal(err)
	}
}
