//go:build go1.18
// +build go1.18

//go:generate autorest ./autorest.md
//go:generate rm ./models_serde.go
//go:generate rm ./models.go
//go:generate gofmt -w .

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azingest
