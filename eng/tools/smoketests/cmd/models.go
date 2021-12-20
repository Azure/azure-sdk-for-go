package cmd

import "fmt"

type ConfigFile struct {
	Packages []Package
}

type Package struct {
	Name                 string
	CoverageGoal         float64
	EnvironmentVariables map[string]string
}

type Module struct {
	Name    string
	Version string
	Replace string
}

type SemVer struct {
	Major, Minor, Patch int
}

func (s SemVer) Newer(s2 SemVer) bool {
	if s.Major > s2.Major {
		return true
	} else if s.Major == s2.Major && s.Minor > s2.Minor {
		return true
	} else if s.Major == s2.Major && s.Minor == s2.Minor && s.Patch > s2.Patch {
		return true
	}
	return false
}

func (s SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}
