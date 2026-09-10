// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && windows && amd64

package azcosmos

import (
	_ "github.com/Azure/azure-cosmos-driver/windows/amd64"
	_ "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2/internal/native/windowsamd64"
)
