//go:build go1.18
// +build go1.18

//go:generate autorest ./settings/autorest.md
//go:generate autorest ./rbac/autorest.md
//go:generate autorest ./backup/autorest.md
//go:generate rm ./backup/constants.go
//go:generate gofmt -w .

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azadmin
