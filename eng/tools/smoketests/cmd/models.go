// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

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
	Package string
	Version string
	Replace string
}
