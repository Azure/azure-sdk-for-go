// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/internal/exported"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
	"time"
)

func TestQueuePermissions_String(t *testing.T) {
	testdata := []struct {
		input    QueuePermissions
		expected string
	}{
		{input: QueuePermissions{Read: true}, expected: "r"},
		{input: QueuePermissions{Add: true}, expected: "a"},
		{input: QueuePermissions{Update: true}, expected: "u"},
		{input: QueuePermissions{Process: true}, expected: "p"},
		{input: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, expected: "raup"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestQueuePermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected QueuePermissions
	}{
		{expected: QueuePermissions{Read: true}, input: "r"},
		{expected: QueuePermissions{Add: true}, input: "a"},
		{expected: QueuePermissions{Update: true}, input: "u"},
		{expected: QueuePermissions{Process: true}, input: "p"},
		{expected: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, input: "raup"},
		{expected: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, input: "puar"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseQueuePermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestQueuePermissions_ParseNegative(t *testing.T) {
	_, err := parseQueuePermissions("puatr") // Here 't' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "116")
}

func TestGetCanonicalName(t *testing.T) {
	testdata := []struct {
		inputAccount string
		inputQueue   string
		expected     string
	}{
		{inputAccount: "fakestorageaccount", inputQueue: "fakestoragequeue", expected: "/queue/fakestorageaccount/fakestoragequeue"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getCanonicalName(c.inputAccount, c.inputQueue))
	}
}

func TestQueueSignatureValues_SignWithSharedKey(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err, "error creating valid shared key credentials.")

	expiryDate, err := time.Parse("2006-01-02", "2023-07-20")
	require.Nil(t, err, "error creating valid expiry date.")

	testdata := []struct {
		object        QueueSignatureValues
		expected      QueryParameters
		expectedError error
	}{
		{
			object:        QueueSignatureValues{QueueName: "fakestoragequeue", Permissions: "r", ExpiryTime: expiryDate},
			expected:      QueryParameters{version: Version, permissions: "r", expiryTime: expiryDate},
			expectedError: nil,
		},
		{
			object:        QueueSignatureValues{QueueName: "fakestoragequeue", Permissions: "", ExpiryTime: expiryDate},
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        QueueSignatureValues{QueueName: "fakestoragequeue", Permissions: "r", ExpiryTime: time.Time{}},
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        QueueSignatureValues{QueueName: "fakestoragequeue", Permissions: "", ExpiryTime: time.Time{}, Identifier: "fakepolicyname"},
			expected:      QueryParameters{version: Version, identifier: "fakepolicyname"},
			expectedError: nil,
		},
	}
	for _, c := range testdata {
		act, err := c.object.SignWithSharedKey(cred)
		// ignore signature value
		act.signature = ""
		require.Equal(t, c.expected, act)
		require.Equal(t, c.expectedError, err)
	}
}

// TestQueueSignatureValues_SignWithSharedKeyExpiredNonUTCExpiry covers the reported scenario: StartTime is unset and ExpiryTime is an
// instant that has already passed, given in UTC+09:00. The SAS must expire at that instant, not 9 hours later.
func TestQueueSignatureValues_SignWithSharedKeyExpiredNonUTCExpiry(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.NoError(t, err)

	sign := func(expiry time.Time) QueryParameters {
		qp, err := QueueSignatureValues{
			Version:     "2020-10-02",
			Protocol:    ProtocolHTTPS,
			ExpiryTime:  expiry,
			Permissions: "r",
			QueueName:   "fakestoragequeue",
		}.SignWithSharedKey(cred)
		require.NoError(t, err)
		return qp
	}

	expiry := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := sign(expiry)
	actual := sign(expiry.In(time.FixedZone("UTC+09:00", 9*60*60)))

	values, err := url.ParseQuery(actual.Encode())
	require.NoError(t, err)
	require.False(t, values.Has("st"))
	require.Equal(t, "2024-01-01T00:00:00Z", values.Get("se"))
	require.Equal(t, expected.Signature(), actual.Signature())
}

func TestQueueSignatureValues_SignWithSharedKeyNonUTCTimes(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.NoError(t, err)

	sign := func(t *testing.T, start, expiry time.Time) QueryParameters {
		t.Helper()
		qp, err := QueueSignatureValues{
			Version:     "2020-10-02",
			Protocol:    ProtocolHTTPS,
			StartTime:   start,
			ExpiryTime:  expiry,
			Permissions: "r",
			QueueName:   "fakestoragequeue",
		}.SignWithSharedKey(cred)
		require.NoError(t, err)
		return qp
	}

	start := time.Date(2026, 10, 15, 11, 0, 0, 0, time.UTC)
	expiry := start.Add(time.Hour)
	expected := sign(t, start, expiry)
	for _, loc := range []*time.Location{time.FixedZone("UTC+09:00", 9*60*60), time.FixedZone("UTC-05:00", -5*60*60)} {
		t.Run(loc.String(), func(t *testing.T) {
			actual := sign(t, start.In(loc), expiry.In(loc))
			require.Equal(t, expected.Encode(), actual.Encode())
		})
	}
}

func TestQueueSignatureValues_SignWithUserDelegationNonUTCTimes(t *testing.T) {
	sign := func(t *testing.T, loc *time.Location) QueryParameters {
		t.Helper()
		start := time.Date(2026, 10, 15, 11, 0, 0, 0, time.UTC).In(loc)
		expiry := start.Add(time.Hour)
		oid, tid, svc, ver, val := "fakeoid", "faketid", "q", "2020-02-10", "AKIAIOSFODNN7EXAMPLE"
		udc := exported.NewUserDelegationCredential("fakeaccountname", exported.UserDelegationKey{
			SignedOID:     &oid,
			SignedTID:     &tid,
			SignedStart:   &start,
			SignedExpiry:  &expiry,
			SignedService: &svc,
			SignedVersion: &ver,
			Value:         &val,
		})
		qp, err := QueueSignatureValues{
			Version:     "2020-10-02",
			Protocol:    ProtocolHTTPS,
			StartTime:   start,
			ExpiryTime:  expiry,
			Permissions: "r",
			QueueName:   "fakestoragequeue",
		}.SignWithUserDelegation(udc)
		require.NoError(t, err)
		return qp
	}

	expected := sign(t, time.UTC)
	for _, loc := range []*time.Location{time.FixedZone("UTC+09:00", 9*60*60), time.FixedZone("UTC-05:00", -5*60*60)} {
		t.Run(loc.String(), func(t *testing.T) {
			actual := sign(t, loc)
			require.Equal(t, expected.Encode(), actual.Encode())
		})
	}
}
