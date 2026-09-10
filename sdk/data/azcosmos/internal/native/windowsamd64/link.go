// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && windows && amd64

// Package windowsamd64 compiles the Go token-provider callback shim for Windows on amd64.
package windowsamd64

/*
#include "../../../azurecosmosdriver.h"
*/
import "C"
