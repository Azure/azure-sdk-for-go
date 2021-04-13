// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/pkgchk/track1"
)

// GeneratedFrom gives the information of the generation metadata, including the commit hash that this package is generated from,
// the readme path, and the tag
func GeneratedFrom(commitHash, readme, tag string) string {
	return fmt.Sprintf("Generated from https://github.com/Azure/azure-rest-api-specs/tree/%s/%s tag: `%s`", commitHash, readme, tag)
}

// GenerationMetadata contains all the metadata that has been used when generating a track 1 package
type GenerationMetadata struct {
	AutorestVersion      string                 `json:"autorest,omitempty"`
	CommitHash           string                 `json:"commit,omitempty"`
	Readme               string                 `json:"readme,omitempty"`
	Tag                  string                 `json:"tag,omitempty"`
	CodeGenVersion       string                 `json:"use,omitempty"`
	RepositoryURL        string                 `json:"repository_url,omitempty"`
	AutorestCommand      string                 `json:"autorest_command,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additional_properties,omitempty"`
}

// String ...
func (m GenerationMetadata) String() string {
	return fmt.Sprintf(`%s

Code generator %s
`, GeneratedFrom(m.CommitHash, m.Readme, m.Tag), m.CodeGenVersion)
}

// Parse parses the metadata info stored in a changelog with certain format into the GenerationMetadata struct
func Parse(reader io.Reader) (*GenerationMetadata, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("expecting at least 3 lines from changelog, but only get %d line(s)", len(lines))
	}
	// parse the first line to get readme, tag and commit hash
	m, err := parseFirstLine(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, err
	}
	m.CodeGenVersion, err = parseThirdLine(strings.TrimSpace(lines[2]))
	if err != nil {
		return nil, err
	}
	return m, nil
}

// CollectGenerationMetadata iterates every track 1 go sdk package under root, and collect all the GenerationMetadata into a map
// using relative path of the package as keys
func CollectGenerationMetadata(root string) (map[string]GenerationMetadata, error) {
	pkgs, err := track1.List(root)
	if err != nil {
		return nil, fmt.Errorf("failed to get track 1 package list under root '%s': %+v", root, err)
	}
	result := make(map[string]GenerationMetadata)
	for _, pkg := range pkgs {
		m, err := GetGenerationMetadata(pkg)
		if err != nil {
			return nil, err
		}
		result[pkg.FullPath()] = *m
	}
	return result, nil
}

// GetGenerationMetadata gets the GenerationMetadata in one specific package
func GetGenerationMetadata(pkg track1.Package) (*GenerationMetadata, error) {
	changelogPath := filepath.Join(pkg.FullPath(), ChangelogFilename)
	file, err := os.Open(changelogPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %+v", changelogPath, err)
	}
	defer file.Close()
	return Parse(file)
}

func parseFirstLine(line string) (*GenerationMetadata, error) {
	matches := firstLineRegex.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil, fmt.Errorf("expecting 4 matches for line '%s', but only get the following matches: [%s]", line, strings.Join(matches, ", "))
	}
	return &GenerationMetadata{
		CommitHash: matches[1],
		Readme:     matches[2],
		Tag:        matches[3],
	}, nil
}

func parseThirdLine(line string) (string, error) {
	matches := thirdLineRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", fmt.Errorf("expecting 2 matches for line '%s', but only get the following matches: [%s]", line, strings.Join(matches, ", "))
	}
	return matches[1], nil
}

var (
	firstLineRegex = regexp.MustCompile("^Generated from https://github\\.com/Azure/azure-rest-api-specs/tree/([0-9a-f]+)/(.+) tag: `(.+)`$")
	thirdLineRegex = regexp.MustCompile(`^Code generator (\S+)$`)
)

const (
	// ChangelogFilename ...
	ChangelogFilename = "CHANGELOG.md"
)
