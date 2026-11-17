// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package blob_test

import (
	"encoding/base64"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/stretchr/testify/require"
)

const sasScopeTestBlobURL = "https://myaccount.blob.core.windows.net/mycontainer/myblob"

func newSASScopeTestClient(t *testing.T) *blob.Client {
	t.Helper()
	cred, err := blob.NewSharedKeyCredential("myaccount", base64.StdEncoding.EncodeToString([]byte("fake-account-key")))
	require.NoError(t, err)
	client, err := blob.NewClientWithSharedKeyCredential(sasScopeTestBlobURL, cred, nil)
	require.NoError(t, err)
	return client
}

func getSASQuery(t *testing.T, sasURL string) url.Values {
	t.Helper()
	u, err := url.Parse(sasURL)
	require.NoError(t, err)
	return u.Query()
}

var sasScopeTestPerms = sas.BlobPermissions{Read: true, Write: true, Delete: true}

func TestGetSASURLValidSnapshotIsSnapshotScoped(t *testing.T) {
	snapshot := "2011-03-09T01:42:34.9360000Z"
	client, err := newSASScopeTestClient(t).WithSnapshot(snapshot)
	require.NoError(t, err)

	sasURL, err := client.GetSASURL(sasScopeTestPerms, time.Now().Add(time.Hour), nil)
	require.NoError(t, err)

	q := getSASQuery(t, sasURL)
	require.Equal(t, "bs", q.Get("sr"))
	require.Equal(t, snapshot, q.Get("snapshot"))
}

func TestGetSASURLInvalidSnapshotReturnsError(t *testing.T) {
	for _, snapshot := range []string{
		"x",
		"2011-03-09T01:42:34Z",         // RFC3339 without the 7-digit fractional seconds
		"0001-01-01T00:00:00.0000000Z", // well-formed, but parses to the zero time
	} {
		t.Run(snapshot, func(t *testing.T) {
			client, err := newSASScopeTestClient(t).WithSnapshot(snapshot)
			require.NoError(t, err)

			sasURL, err := client.GetSASURL(sasScopeTestPerms, time.Now().Add(time.Hour), nil)
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid snapshot")
			require.Empty(t, sasURL)
		})
	}
}

func TestGetSASURLInvalidSnapshotReturnsErrorForAllBlobClients(t *testing.T) {
	cred, err := blob.NewSharedKeyCredential("myaccount", base64.StdEncoding.EncodeToString([]byte("fake-account-key")))
	require.NoError(t, err)
	u := sasScopeTestBlobURL + "?snapshot=x"
	expiry := time.Now().Add(time.Hour)

	bb, err := blockblob.NewClientWithSharedKeyCredential(u, cred, nil)
	require.NoError(t, err)
	_, err = bb.GetSASURL(sasScopeTestPerms, expiry, nil)
	require.Error(t, err)

	ab, err := appendblob.NewClientWithSharedKeyCredential(u, cred, nil)
	require.NoError(t, err)
	_, err = ab.GetSASURL(sasScopeTestPerms, expiry, nil)
	require.Error(t, err)

	pb, err := pageblob.NewClientWithSharedKeyCredential(u, cred, nil)
	require.NoError(t, err)
	_, err = pb.GetSASURL(sasScopeTestPerms, expiry, nil)
	require.Error(t, err)
}

func TestGetSASURLEmptySnapshotIsBaseBlobScoped(t *testing.T) {
	sasURL, err := newSASScopeTestClient(t).GetSASURL(sasScopeTestPerms, time.Now().Add(time.Hour), nil)
	require.NoError(t, err)

	q := getSASQuery(t, sasURL)
	require.Equal(t, "b", q.Get("sr"))
	require.Empty(t, q.Get("snapshot"))
}

func TestGetSASURLVersionIDIsVersionScoped(t *testing.T) {
	for _, versionID := range []string{
		"2011-03-09T01:42:34.9360000Z", // valid version ID
		"x",                            // malformed version ID
	} {
		t.Run(versionID, func(t *testing.T) {
			client, err := newSASScopeTestClient(t).WithVersionID(versionID)
			require.NoError(t, err)

			sasURL, err := client.GetSASURL(sasScopeTestPerms, time.Now().Add(time.Hour), nil)
			require.NoError(t, err)

			// A version ID, well-formed or not, must never fall back to a base-blob SAS: it is
			// always scoped to sr=bv and the version ID is bound into the signature, so the
			// service rejects any request for a different (or no) version.
			q := getSASQuery(t, sasURL)
			require.Equal(t, "bv", q.Get("sr"))
			require.Equal(t, versionID, q.Get("versionid"))
		})
	}
}

func TestGetSASURLVersionIDBoundIntoSignature(t *testing.T) {
	base := newSASScopeTestClient(t)
	expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	sig := func(versionID string) string {
		c, err := base.WithVersionID(versionID)
		require.NoError(t, err)
		sasURL, err := c.GetSASURL(sasScopeTestPerms, expiry, nil)
		require.NoError(t, err)
		return getSASQuery(t, sasURL).Get("sig")
	}

	baseURL, err := base.GetSASURL(sasScopeTestPerms, expiry, nil)
	require.NoError(t, err)
	baseSig := getSASQuery(t, baseURL).Get("sig")

	v1 := sig("2011-03-09T01:42:34.9360000Z")
	v2 := sig("2011-03-09T01:42:35.9360000Z")
	require.NotEqual(t, v1, v2)
	require.NotEqual(t, baseSig, v1)
	require.NotEqual(t, baseSig, sig("x"))
}
