// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && linux && !android && amd64

// Package linuxamd64 provides the Go token callback shim for glibc Linux on amd64.
package linuxamd64

/*
#include <features.h>

#ifndef __GLIBC__
#error "azcosmos: the linux/amd64 driver requires glibc; musl is not supported yet"
#endif
*/
import "C"
