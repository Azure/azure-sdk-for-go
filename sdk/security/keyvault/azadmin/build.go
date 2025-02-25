//go:generate tsp-client update --output-dir ./settings --local-spec-repo /home/grace/code/azure-rest-api-specs/specification/keyvault/Security.KeyVault.Settings
//go:generate tsp-client update --output-dir ./rbac --local-spec-repo /home/grace/code/azure-rest-api-specs/specification/keyvault/Security.KeyVault.RBAC
//go:generate tsp-client update --output-dir ./backup --local-spec-repo /home/grace/code/azure-rest-api-specs/specification/keyvault/Security.KeyVault.BackupRestore
//go:generate go run ./internal/generate/transforms.go
//go:generate goimports -w .
//go:generate rm ./backup/constants.go

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin
