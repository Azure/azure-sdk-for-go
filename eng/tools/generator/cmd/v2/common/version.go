package common

import (
	"github.com/Masterminds/semver"
)

func IsBetaVersion(v string) (bool, error) {

	newVersion, err := semver.NewVersion(v)
	if err != nil {
		return false, err
	}

	if newVersion.Major() == 0 || newVersion.Prerelease() != "" {
		return true, nil
	}

	return false, nil
}
