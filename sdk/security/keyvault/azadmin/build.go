//go:generate tsp-client update --output-dir ./settings
//go:generate tsp-client update --output-dir ./rbac
//go:generate tsp-client update --output-dir ./backup
//go:generate go run ./internal/generate/transforms.go
//go:generate gofmt -w .
//go:generate rm ./backup/constants.go

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin
