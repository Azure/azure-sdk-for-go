//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azacr

import generated "github.com/Azure/azure-sdk-for-go/sdk/acr/azacr/internal"

const (
	moduleName    = "azacr"
	moduleVersion = "v0.1.0"
)

// ArtifactArchitecture - The artifact platform's architecture.
type ArtifactArchitecture = generated.ArtifactArchitecture

const (
	// ArtifactArchitectureAmd64 - AMD64
	ArtifactArchitectureAmd64 = generated.ArtifactArchitectureAmd64
	// ArtifactArchitectureArm - ARM
	ArtifactArchitectureArm = generated.ArtifactArchitectureArm
	// ArtifactArchitectureArm64 - ARM64
	ArtifactArchitectureArm64 = generated.ArtifactArchitectureArm64
	// ArtifactArchitectureI386 - i386
	ArtifactArchitectureI386 = generated.ArtifactArchitectureI386
	// ArtifactArchitectureMips - MIPS
	ArtifactArchitectureMips = generated.ArtifactArchitectureMips
	// ArtifactArchitectureMips64 - MIPS64
	ArtifactArchitectureMips64 = generated.ArtifactArchitectureMips64
	// ArtifactArchitectureMips64Le - MIPS64LE
	ArtifactArchitectureMips64Le = generated.ArtifactArchitectureMips64Le
	// ArtifactArchitectureMipsLe - MIPSLE
	ArtifactArchitectureMipsLe = generated.ArtifactArchitectureMipsLe
	// ArtifactArchitecturePpc64 - PPC64
	ArtifactArchitecturePpc64 = generated.ArtifactArchitecturePpc64
	// ArtifactArchitecturePpc64Le - PPC64LE
	ArtifactArchitecturePpc64Le = generated.ArtifactArchitecturePpc64Le
	// ArtifactArchitectureRiscV64 - RISCv64
	ArtifactArchitectureRiscV64 = generated.ArtifactArchitectureRiscV64
	// ArtifactArchitectureS390X - s390x
	ArtifactArchitectureS390X = generated.ArtifactArchitectureS390X
	// ArtifactArchitectureWasm - Wasm
	ArtifactArchitectureWasm = generated.ArtifactArchitectureWasm
)

// PossibleArtifactArchitectureValues returns the possible values for the ArtifactArchitecture const type.
func PossibleArtifactArchitectureValues() []ArtifactArchitecture {
	return generated.PossibleArtifactArchitectureValues()
}

type ArtifactOperatingSystem = generated.ArtifactOperatingSystem

const (
	ArtifactOperatingSystemAix       = generated.ArtifactOperatingSystemAix
	ArtifactOperatingSystemAndroid   = generated.ArtifactOperatingSystemAndroid
	ArtifactOperatingSystemDarwin    = generated.ArtifactOperatingSystemDarwin
	ArtifactOperatingSystemDragonfly = generated.ArtifactOperatingSystemDragonfly
	ArtifactOperatingSystemFreeBsd   = generated.ArtifactOperatingSystemFreeBsd
	ArtifactOperatingSystemIOS       = generated.ArtifactOperatingSystemIOS
	ArtifactOperatingSystemIllumos   = generated.ArtifactOperatingSystemIllumos
	ArtifactOperatingSystemJS        = generated.ArtifactOperatingSystemJS
	ArtifactOperatingSystemLinux     = generated.ArtifactOperatingSystemLinux
	ArtifactOperatingSystemNetBsd    = generated.ArtifactOperatingSystemNetBsd
	ArtifactOperatingSystemOpenBsd   = generated.ArtifactOperatingSystemOpenBsd
	ArtifactOperatingSystemPlan9     = generated.ArtifactOperatingSystemPlan9
	ArtifactOperatingSystemSolaris   = generated.ArtifactOperatingSystemSolaris
	ArtifactOperatingSystemWindows   = generated.ArtifactOperatingSystemWindows
)

// PossibleArtifactOperatingSystemValues returns the possible values for the ArtifactOperatingSystem const type.
func PossibleArtifactOperatingSystemValues() []ArtifactOperatingSystem {
	return generated.PossibleArtifactOperatingSystemValues()
}

// ManifestOrderBy - Sort options for ordering manifests in a collection.
type ManifestOrderBy = generated.ArtifactManifestOrderBy

const (
	// ManifestOrderByNone - Do not provide an orderby value in the request.
	ManifestOrderByNone = generated.ArtifactManifestOrderByNone
	// ManifestOrderByLastUpdatedOnDescending - Order manifests by LastUpdatedOn field, from most recently updated to least recently updated.
	ManifestOrderByLastUpdatedOnDescending = generated.ArtifactManifestOrderByLastUpdatedOnDescending
	// ManifestOrderByLastUpdatedOnAscending - Order manifest by LastUpdatedOn field, from least recently updated to most recently updated.
	ManifestOrderByLastUpdatedOnAscending = generated.ArtifactManifestOrderByLastUpdatedOnAscending
)

// PossibleManifestOrderByValues returns the possible values for the ManifestOrderBy const type.
func PossibleManifestOrderByValues() []ManifestOrderBy {
	return generated.PossibleArtifactManifestOrderByValues()
}

// TagOrderBy - Sort options for ordering tags in a collection.
type TagOrderBy = generated.ArtifactTagOrderBy

const (
	// TagOrderByNone - Do not provide an orderby value in the request.
	TagOrderByNone = generated.ArtifactTagOrderByNone
	// TagOrderByLastUpdatedOnDescending - Order tags by LastUpdatedOn field, from most recently updated to least recently updated.
	TagOrderByLastUpdatedOnDescending = generated.ArtifactTagOrderByLastUpdatedOnDescending
	// TagOrderByLastUpdatedOnAscending - Order tags by LastUpdatedOn field, from least recently updated to most recently updated.
	TagOrderByLastUpdatedOnAscending = generated.ArtifactTagOrderByLastUpdatedOnAscending
)

// PossibleTagOrderByValues returns the possible values for the TagOrderBy const type.
func PossibleTagOrderByValues() []TagOrderBy {
	return generated.PossibleArtifactTagOrderByValues()
}
