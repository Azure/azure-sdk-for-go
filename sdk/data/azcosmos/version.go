// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// cSpell:ignore azsdk

import "strings"

const (
	// moduleName intentionally omits the major version suffix; it is reported in telemetry.
	//
	//nolint:unused // consumed once client construction lands.
	moduleName = "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

	// serviceLibVersion is the semantic version (see http://semver.org) of this module.
	serviceLibVersion = "v2.0.0-beta.1"
)

func wrappingSDKIdentifier() string {
	return "azsdk-go-azcosmos/" + strings.TrimPrefix(serviceLibVersion, "v")
}
