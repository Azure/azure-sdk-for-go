//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToPtrMethods(t *testing.T) {
	d := DeletionRecoveryLevelCustomizedRecoverable
	require.Equal(t, d.ToPtr(), &d)

	j := JSONWebKeyCurveNameP256
	require.Equal(t, j.ToPtr(), &j)

	o := JSONWebKeyOperationDecrypt
	require.Equal(t, o.ToPtr(), &o)
}

func TestToGeneratedMethods(t *testing.T) {
	var l *ListKeyVersionsOptions //nolint
	l = nil
	_ = l.toGenerated()
}
