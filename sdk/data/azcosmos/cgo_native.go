// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

/*
#cgo CFLAGS: -I${SRCDIR}/internal/native

// The library is linked dynamically, matching the driver modules in
// github.com/Azure/azure-cosmos-driver. A static archive would avoid the run-time search below,
// but it is six times the size, which matters while these binaries are committed here.
//
// Because it is dynamic, the library has to be resolvable at run time and not just at build time.
// The loader searches rpath entries in order, so the executable-relative ones come first: a library
// shipped alongside a binary has to win over the path it happened to be built from. The build path
// is kept last so that `go test` works without staging a copy, but it is only a fallback.
//
// Building a package that imports this one therefore needs cgo's rpath allowlist opened, because
// cgo rejects the executable-relative forms by default:
//
//	CGO_LDFLAGS_ALLOW='^-Wl,-rpath,(@(executable_path|loader_path)|\$ORIGIN)$' go test ./...
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/internal/native/lib/darwin_arm64 -lazurecosmosdriver
#cgo darwin,arm64 LDFLAGS: -Wl,-rpath,@executable_path -Wl,-rpath,@loader_path -Wl,-rpath,${SRCDIR}/internal/native/lib/darwin_arm64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/internal/native/lib/linux_amd64 -lazurecosmosdriver
#cgo linux,amd64 LDFLAGS: -Wl,-rpath,$ORIGIN -Wl,-rpath,${SRCDIR}/internal/native/lib/linux_amd64

#include "azurecosmosdriver.h"
*/
import "C"

// This file carries the cgo directives for the whole package. They are declared once here rather
// than repeated in each file that imports "C", because cgo unions the directives across the
// package: repeating them means every file has to be kept in step, and a file that drifts is a
// link error rather than a compile error.
//
// The libraries under internal/native/lib are committed, which is deliberately temporary. The
// distribution design puts the per-target binaries in their own modules in a separate repository,
// consumed as ordinary Go modules; that repository exists and already carries a darwin/arm64
// module, but not yet one for the platform CI runs on. Committing them here unblocks that in the
// meantime and is expected to be reverted, which is why only the two platforms that are actually
// built and tested are present rather than a full matrix.
