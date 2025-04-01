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
	replace []byte
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
			replace: []byte(replace),
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
			after := t.regex.ReplaceAll(b, t.replace)
			if bytes.Equal(b, after) {
				log.Printf(`replacement "%s -> %s" had no effect in %s`, t.regex, t.replace, filepath.Base(path))
			}
			b = after
		}
		if err := f.Truncate(0); err != nil {
			return err
		}
		if _, err = f.WriteAt(b, 0); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	r := replacer{}
	for before, after := range map[string]string{
		"OcpBatchFileIsdirectory":    "OCPBatchFileIsDirectory",
		"OcpBatchFile((?:Mode|URL))": "OCPBatchFile$1",
		"OcpCreationTime":            "OCPCreationTime",
	} {
		r.Replace([]string{"client.go", "responses.go"}, before, after)
	}
	for before, after := range map[string]string{
		"Ocpdate":  "OCPDate",
		"OcpRange": "OCPRange",
	} {
		r.Replace([]string{"client.go", "options.go"}, before, after)
	}
	// ETag fields should be azcore.ETag, not string
	r.Replace(
		[]string{"models.go", "options.go", "responses.go"},
		`((?:ETag|If(?:None)?Match) )\*string`,
		"$1*azcore.ETag",
	)
	for before, after := range map[string]string{
		`(\*\w+\.If(None)?Match)`: "string($1)",
		`(\w+\.ETag = )(&\w+)`:    "${1}(*azcore.ETag)($2)",
	} {
		r.Replace([]string{"client.go"}, before, after)
	}
	// add import for azcore.ETag. This would break if
	// the emitter added another import to these files
	r.Replace(
		[]string{"models.go", "options.go"},
		`import "time"`,
		"import (\n\t\"time\"\n\t\"github.com/Azure/azure-sdk-for-go/sdk/azcore\"\n)",
	)
	if err := r.Do(); err != nil {
		log.Fatal(err)
	}
}
