package changelog

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

// GeneratedFrom gives the information of the generation metadata, including the commit hash that this package is generated from,
// the readme path, and the tag
func GeneratedFrom(commitHash, readme, tag string) string {
	return fmt.Sprintf("Generated from https://github.com/Azure/azure-rest-api-specs/tree/%s/%s tag: `%s`", commitHash, readme, tag)
}

// GenerationMetadata contains all the metadata that has been used when generating a track 1 package
type GenerationMetadata struct {
	CommitHash     string
	Readme         string
	Tag            string
	CodeGenVersion string
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
