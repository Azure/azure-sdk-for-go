//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:generate pwsh ./testdata/genopenapi.ps1
//go:generate autorest  ./autorest.md
//go:generate go mod tidy
//go:generate goimports -w .

//go:generate pwsh ./testdata/rename_custom_and_tests.ps1
//go:generate go run ./internal/transform/cmd -op rename-method -file client.go -name "(*Client).GetChatCompletions" -new-name getChatCompletions
//go:generate go run ./internal/transform/cmd -op rename-method -file client.go -name "(*Client).GetCompletions" -new-name getCompletions
//go:generate go run ./internal/transform/cmd -op rename-struct -file models.go -name ChatCompletionsOptions -new-name chatCompletionsOptions
//go:generate go run ./internal/transform/cmd -op rename-struct -file models.go -name CompletionsOptions -new-name completionsOptions
//go:generate go run ./internal/transform/cmd -op copy-struct -file models.go -name completionsOptions -new-name CompletionsOptions
//go:generate go run ./internal/transform/cmd -op copy-struct -file models.go -name completionsOptions -new-name CompletionsStreamOptions
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct CompletionsOptions -field Stream
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct CompletionsOptions -field StreamOptions
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct CompletionsStreamOptions -field Stream
//go:generate go run ./internal/transform/cmd -op copy-struct -file models.go -name chatCompletionsOptions -new-name ChatCompletionsOptions
//go:generate go run ./internal/transform/cmd -op copy-struct -file models.go -name chatCompletionsOptions -new-name ChatCompletionsStreamOptions
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct ChatCompletionsOptions -field Stream
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct ChatCompletionsOptions -field StreamOptions
//go:generate go run ./internal/transform/cmd -op remove-field -file models.go -struct ChatCompletionsStreamOptions -field Stream
//go:generate pwsh ./testdata/rename_custom_and_tests.ps1 -Reverse

//go:generate go mod tidy
//go:generate goimports -w .

// running the tests that check that generation went the way we expected to.
//go:go test -v ./internal

package azopenai
