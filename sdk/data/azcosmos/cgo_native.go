// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

// This file carries the cgo directives for the whole package. They are declared once here rather
// than repeated in each file that imports "C", because cgo unions the directives across the
// package: repeating them means every file has to be kept in step, and a file that drifts is a
// link error rather than a compile error.
//
// Target-specific internal packages carry the .syso archives and system-linker flags. Keeping each
// archive behind a conditionally imported package matters because Go treats ios as darwin and
// android as linux for build tags and .syso filename selection.
