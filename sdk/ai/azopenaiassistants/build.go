//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

//go:generate pwsh ./testdata/genopenapi.ps1
//go:generate go run ./internal/transform
//go:generate goimports -w ./..
//go:generate go mod tidy
