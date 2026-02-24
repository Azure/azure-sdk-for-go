// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:generate pwsh ./testdata/genopenapi.ps1
//go:generate autorest  ./autorest.md
//go:generate rm -f options.go openai_client.go responses.go
//go:generate go mod tidy
//go:generate gofmt -w .

// running the tests that check that generation went the way we expected to.
//go:go test -v ./internal

package azopenai
