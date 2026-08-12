// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// nativeDriverVersion pins the azure_data_cosmos_driver_native build this module binds to. The
// C ABI is only guaranteed compatible within a pinned version, so this must be updated
// deliberately and in lockstep with the vendored driver binaries.
//
// This deliberately lives outside version.go: the release tooling greps *version*.go for the
// first semver-shaped literal to determine the module version, so a second version constant
// there would eventually be misreported as the module version.
//
// TODO: confirm the GA pin with the Rust driver crew before API sign-off.
//
//nolint:unused // consumed once client construction lands.
const nativeDriverVersion = "0.1.0"
