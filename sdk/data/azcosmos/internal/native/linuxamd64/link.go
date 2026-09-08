// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && linux && !android && amd64

// Package linuxamd64 links the native Cosmos driver for glibc Linux on amd64.
package linuxamd64

/*
#include <features.h>

#ifndef __GLIBC__
#error "azcosmos: the bundled linux/amd64 driver requires glibc; musl is not supported yet"
#endif

#cgo LDFLAGS: -lm -ldl -lpthread
*/
import "C"
