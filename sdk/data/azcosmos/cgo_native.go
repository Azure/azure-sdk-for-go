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
// TEMPORARY: the -L and -rpath entries pointing into internal/native/lib go away with the
// committed libraries themselves, once azure-cosmos-driver carries a module for every platform
// this builds on. See internal/native/lib/README.md for the trigger and the steps.
//
// Because it is dynamic, the library has to be resolvable at run time and not just at build time.
// Only executable-relative rpaths are recorded, so a library shipped alongside a binary is found
// and nothing else is searched.
//
// The build directory is deliberately not an rpath. ${SRCDIR} is an absolute path on whoever built
// the binary — a module cache path for a consumer — and baking it in makes the loader search a
// directory the running host may let someone else create, which is a library-injection foothold in
// every binary built against this package. Point the loader at the library explicitly instead when
// running from the tree, which is what the test invocations below do.
//
// Building a package that imports this one needs cgo's rpath allowlist opened, because cgo rejects
// the executable-relative forms by default. Running the tests from the tree additionally needs the
// library directory as an rpath, because the test binary is built into a temporary directory and
// not next to the library. That rpath is passed at the invocation rather than declared here, so it
// reaches the test binary and never a consumer's:
//
//	CGO_LDFLAGS_ALLOW='^-Wl,-rpath,(@(executable_path|loader_path)|\$ORIGIN)$' \
//	  go test -tags azcosmos_driver \
//	    -ldflags "-extldflags=-Wl,-rpath,$PWD/internal/native/lib/darwin_arm64" ./...
//
// and on linux the same with linux_amd64. DYLD_LIBRARY_PATH is not an option on darwin: macOS
// strips DYLD_* from the environment of anything a protected process spawns, so it does not
// survive as far as the test binary.
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/internal/native/lib/darwin_arm64 -lazurecosmosdriver
#cgo darwin,arm64 LDFLAGS: -Wl,-rpath,@executable_path -Wl,-rpath,@loader_path
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/internal/native/lib/linux_amd64 -lazurecosmosdriver
#cgo linux,amd64 LDFLAGS: -Wl,-rpath,$ORIGIN

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
