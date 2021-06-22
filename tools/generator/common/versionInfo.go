package common

import "fmt"

type VersionInfo struct {
	LatestVersion string
	NewVersion    string
}

func (i VersionInfo) String() string {
	return fmt.Sprintf("Latest version: %s, new version: %s", i.LatestVersion, i.NewVersion)
}
