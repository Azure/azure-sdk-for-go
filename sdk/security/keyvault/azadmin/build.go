//go:generate tsp-client update --output-dir ./settings
//go:generate tsp-client update --output-dir ./rbac
//go:generate tsp-client update --output-dir ./backup
//go:generate go run ./internal/generate/transforms.go
//go:generate goimports -w .
//go:generate rm ./backup/constants.go
//go:generate rm ./backup/go.mod
//go:generate rm ./rbac/go.mod
//go:generate rm ./settings/go.mod

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin
