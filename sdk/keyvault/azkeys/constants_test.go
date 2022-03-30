//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestToPtrMethods(t *testing.T) {
	d := DeletionRecoveryLevelCustomizedRecoverable
	require.Equal(t, to.Ptr(d), &d)

	j := CurveNameP256
	require.Equal(t, to.Ptr(j), &j)

	o := OperationDecrypt
	require.Equal(t, to.Ptr(o), &o)

	a := ExportEncryptionAlgorithmRSAAESKEYWRAP256
	require.Equal(t, to.Ptr(a), &a)
}

//nolint
func TestToGeneratedMethods(t *testing.T) {
	// If done incorrectly, this will have a nil pointer reference
	var l *ListPropertiesOfKeyVersionsOptions = nil
	_ = l.toGenerated()
}
