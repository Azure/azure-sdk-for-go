//go:build go1.18
// +build go1.18

//go:generate tsp-client sync --local-spec-repo ./keyvault
//go:generate tsp-client generate
//go:generate go run ./testdata/transforms.go
//go:generate rm ./constants.go
//go:generate gofmt -w .

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets
