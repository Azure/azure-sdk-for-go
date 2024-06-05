// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:generate pwsh ./testdata/gen.ps1
//go:generate goimports -w .
//go:generate go run ./internal/generate
//go:generate goimports -w .

package aznamespaces
