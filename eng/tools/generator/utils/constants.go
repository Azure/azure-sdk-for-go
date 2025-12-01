// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

const (
	SDKRemoteURL = "https://github.com/Azure/azure-sdk-for-go.git"

	ChangelogFileName     = "CHANGELOG.md"
	GoModFileName         = "go.mod"
	SdkRootPath           = "/sdk"
	ReadmeFileName        = "README.md"
	ClientFactoryFileName = "client_factory.go"

	SDKReleaseTypeStable  = "stable"
	SDKReleaseTypePreview = "preview"

	AutomationRunModeRelease = "release"
	AutomationRunModeLocal   = "local"
)

// PullRequestLabel represents a pull request label type
type PullRequestLabel string

const (
	StableLabel                PullRequestLabel = "Stable"
	BetaLabel                  PullRequestLabel = "Beta"
	FirstGALabel               PullRequestLabel = "FirstGA"
	FirstGABreakingChangeLabel PullRequestLabel = "FirstGA,BreakingChange"
	FirstBetaLabel             PullRequestLabel = "FirstBeta"
	StableBreakingChangeLabel  PullRequestLabel = "Stable,BreakingChange"
	BetaBreakingChangeLabel    PullRequestLabel = "Beta,BreakingChange"
)

// PackageStatus represents the type of package being processed
type PackageStatus int

const (
	// PackageStatusNew represents a new package that doesn't exist yet
	PackageStatusNew PackageStatus = iota
	// PackageStatusExisting represents an existing package
	PackageStatusExisting
)
