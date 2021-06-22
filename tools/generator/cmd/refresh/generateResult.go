// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package refresh

import "github.com/Azure/azure-sdk-for-go/tools/generator/autorest"

type GenerateResult struct {
	Readme     string
	Tag        string
	CommitHash string
	Package    autorest.ChangelogResult
}
