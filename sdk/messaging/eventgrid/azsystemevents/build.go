//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:generate autorest ./autorest.md
//go:generate goimports -w ./..
//go:generate go run ./internal/generate generate
//go:generate goimports -w ./..

package azsystemevents
