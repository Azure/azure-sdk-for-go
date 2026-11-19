//go:build !noattestation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import _ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"

const attestationImported = true
