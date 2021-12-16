package main

type ConfigFile struct {
	Packages []Package
}

type Package struct {
	Name                 string
	CoverageGoal         float64
	EnvironmentVariables map[string]string
}
