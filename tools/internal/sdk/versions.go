package sdk

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

func ModifyVersionFile(absSDK, lastVersion, newVersion string) error {
	versionFile := VersionGoPath(absSDK)
	info, err := os.Stat(versionFile)
	if err != nil {
		return fmt.Errorf("failed to get stat of version file: %+v", err)
	}
	b, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return fmt.Errorf("failed to read version file: %+v", err)
	}

	content := strings.ReplaceAll(string(b), lastVersion, newVersion)

	// write everything back
	if err := ioutil.WriteFile(versionFile, []byte(content), info.Mode()); err != nil {
		return fmt.Errorf("failed to write version file: %+v", err)
	}
	return nil
}

func GetVersion(content string) (*semver.Version, error) {
	regex := regexp.MustCompile(`const Number = "(.+)"`)
	matches := regex.FindStringSubmatch(content)
	if len(matches) < 1 {
		return nil, fmt.Errorf("cannot find the version number in version.go")
	}
	return semver.NewVersion(matches[1])
}
