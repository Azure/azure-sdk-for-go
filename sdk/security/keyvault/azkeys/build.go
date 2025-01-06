//go:build go1.18
// +build go1.18

//go:generate tsp-client update --local-spec-repo /home/grace/code/azure-rest-api-specs/specification/keyvault/Security.KeyVault.Keys
//go:generate go run ./internal/transforms.go
//go:generate gofmt -w .

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys
