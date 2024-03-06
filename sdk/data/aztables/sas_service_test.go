// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSASSignatureValues(t *testing.T) {
	fakeKey := base64.StdEncoding.EncodeToString([]byte("fake-key"))
	cred, err := NewSharedKeyCredential("fake-account", fakeKey)
	require.NoError(t, err)

	startTime, err := time.Parse(time.RFC3339, "2023-11-08T15:04:05Z")
	require.NoError(t, err)

	sasValues := SASSignatureValues{
		Protocol:   SASProtocolHTTPS,
		StartTime:  startTime,
		ExpiryTime: startTime.Add(time.Hour),
		Permissions: SASPermissions{
			Read: true,
		}.String(),
		TableName: "fake-table",
	}
	sig, err := sasValues.Sign(cred)
	require.NoError(t, err)
	const expected = "se=2023-11-08T16%3A04%3A05Z&sig=WLeRe04Jnm2q7wuetbkWgFDtdWg%2BiE7RKwUSLbecPjE%3D&sp=r&spr=https&st=2023-11-08T15%3A04%3A05Z&sv=2019-02-02&tn=fake-table"
	require.EqualValues(t, expected, sig)
}
