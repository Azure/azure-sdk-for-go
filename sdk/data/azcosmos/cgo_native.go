// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

/*
#cgo CFLAGS: -I${SRCDIR}/internal/native

// The driver is linked as a static archive rather than a shared library, so a program built
// against it carries the driver with it: no library has to be installed on the machine that runs
// it, and nothing has to be on a search path at startup. That is what lets the eventual per-target
// driver modules be ordinary Go modules.
//
// The system libraries below are what the archive itself depends on. A shared library would carry
// those dependencies in its own header; a static archive leaves them to whoever links it.
#cgo linux LDFLAGS: ${SRCDIR}/internal/native/lib/libazurecosmosdriver.a -lm -ldl -lpthread
#cgo darwin LDFLAGS: ${SRCDIR}/internal/native/lib/libazurecosmosdriver.a -liconv -framework Security -framework CoreFoundation
#cgo windows LDFLAGS: ${SRCDIR}/internal/native/lib/libazurecosmosdriver.a -lntdll -lsecur32 -lcrypt32 -lbcrypt -lws2_32 -luserenv

#include "azurecosmosdriver.h"
*/
import "C"

// This file carries the cgo directives for the whole package. They are declared once here rather
// than repeated in each file that imports "C", because cgo unions the directives across the
// package: repeating them means every file has to be kept in step, and a file that drifts is a
// link error rather than a compile error.
//
// The archive is expected at internal/native/lib. It is not committed: the distribution design
// puts the per-target binaries in their own modules, in a separate repository, so this location is
// where a locally built driver goes during development. See the README.
