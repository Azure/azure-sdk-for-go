//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatTimesForSigning(t *testing.T) {
	testdata := []struct {
		inputStart       time.Time
		inputEnd         time.Time
		inputSnapshot    time.Time
		expectedStart    string
		expectedEnd      string
		expectedSnapshot string
	}{
		{expectedStart: "", expectedEnd: "", expectedSnapshot: ""},
		{inputStart: time.Date(1955, 6, 25, 22, 15, 56, 345456, time.UTC), expectedStart: "1955-06-25T22:15:56Z", expectedEnd: "", expectedSnapshot: ""},
		{inputEnd: time.Date(2023, 4, 5, 8, 50, 27, 4500, time.UTC), expectedStart: "", expectedEnd: "2023-04-05T08:50:27Z", expectedSnapshot: ""},
		{inputSnapshot: time.Date(2021, 1, 5, 22, 15, 33, 1234879, time.UTC), expectedStart: "", expectedEnd: "", expectedSnapshot: "2021-01-05T22:15:33.0012348Z"},
		{
			inputStart:       time.Date(1955, 6, 25, 22, 15, 56, 345456, time.UTC),
			inputEnd:         time.Date(2023, 4, 5, 8, 50, 27, 4500, time.UTC),
			inputSnapshot:    time.Date(2021, 1, 5, 22, 15, 33, 1234879, time.UTC),
			expectedStart:    "1955-06-25T22:15:56Z",
			expectedEnd:      "2023-04-05T08:50:27Z",
			expectedSnapshot: "2021-01-05T22:15:33.0012348Z",
		},
	}
	for _, c := range testdata {
		start, end, ss := formatTimesForSigning(c.inputStart, c.inputEnd, c.inputSnapshot)
		require.Equal(t, c.expectedStart, start)
		require.Equal(t, c.expectedEnd, end)
		require.Equal(t, c.expectedSnapshot, ss)
	}
}

func TestFormatTimeWithDefaultFormat(t *testing.T) {
	testdata := []struct {
		input        time.Time
		expectedTime string
	}{
		{input: time.Date(1955, 4, 5, 8, 50, 27, 4500, time.UTC), expectedTime: "1955-04-05T08:50:27Z"},
		{input: time.Date(1917, 3, 9, 16, 22, 56, 0, time.UTC), expectedTime: "1917-03-09T16:22:56Z"},
		{input: time.Date(2021, 1, 5, 22, 15, 0, 0, time.UTC), expectedTime: "2021-01-05T22:15:00Z"},
		{input: time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC), expectedTime: "2023-06-25T00:00:00Z"},
	}
	for _, c := range testdata {
		formattedTime := formatTimeWithDefaultFormat(&c.input)
		require.Equal(t, c.expectedTime, formattedTime)
	}
}

func TestFormatTime(t *testing.T) {
	testdata := []struct {
		input        time.Time
		format       string
		expectedTime string
	}{
		{input: time.Date(1955, 4, 5, 8, 50, 27, 4500, time.UTC), format: "2006-01-02T15:04:05.0000000Z", expectedTime: "1955-04-05T08:50:27.0000045Z"},
		{input: time.Date(1955, 4, 5, 8, 50, 27, 4500, time.UTC), format: "", expectedTime: "1955-04-05T08:50:27Z"},
		{input: time.Date(1917, 3, 9, 16, 22, 56, 0, time.UTC), format: "2006-01-02T15:04:05Z", expectedTime: "1917-03-09T16:22:56Z"},
		{input: time.Date(1917, 3, 9, 16, 22, 56, 0, time.UTC), format: "", expectedTime: "1917-03-09T16:22:56Z"},
		{input: time.Date(2021, 1, 5, 22, 15, 0, 0, time.UTC), format: "2006-01-02T15:04Z", expectedTime: "2021-01-05T22:15Z"},
		{input: time.Date(2021, 1, 5, 22, 15, 0, 0, time.UTC), format: "", expectedTime: "2021-01-05T22:15:00Z"},
		{input: time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC), format: "2006-01-02", expectedTime: "2023-06-25"},
		{input: time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC), format: "", expectedTime: "2023-06-25T00:00:00Z"},
	}
	for _, c := range testdata {
		formattedTime := formatTime(&c.input, c.format)
		require.Equal(t, c.expectedTime, formattedTime)
	}
}

func TestParseTime(t *testing.T) {
	testdata := []struct {
		input          string
		expectedTime   time.Time
		expectedFormat string
	}{
		{input: "1955-04-05T08:50:27.0000045Z", expectedTime: time.Date(1955, 4, 5, 8, 50, 27, 4500, time.UTC), expectedFormat: "2006-01-02T15:04:05.0000000Z"},
		{input: "1917-03-09T16:22:56Z", expectedTime: time.Date(1917, 3, 9, 16, 22, 56, 0, time.UTC), expectedFormat: "2006-01-02T15:04:05Z"},
		{input: "2021-01-05T22:15Z", expectedTime: time.Date(2021, 1, 5, 22, 15, 0, 0, time.UTC), expectedFormat: "2006-01-02T15:04Z"},
		{input: "2023-06-25", expectedTime: time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC), expectedFormat: "2006-01-02"},
	}
	for _, c := range testdata {
		parsedTime, format, err := parseTime(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expectedTime, parsedTime)
		require.Equal(t, c.expectedFormat, format)
	}
}

func TestParseTimeNegative(t *testing.T) {
	_, _, err := parseTime("notatime")
	require.Error(t, err, "fail to parse time with IOS 8601 formats, please refer to https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-a-service-sas for more details")
}

func TestIPRange_String(t *testing.T) {
	testdata := []struct {
		inputStart net.IP
		inputEnd   net.IP
		expected   string
	}{
		{expected: ""},
		{inputStart: net.IPv4(10, 255, 0, 0), expected: "10.255.0.0"},
		{inputStart: net.IPv4(10, 255, 0, 0), inputEnd: net.IPv4(10, 255, 0, 50), expected: "10.255.0.0-10.255.0.50"},
	}
	for _, c := range testdata {
		var ipRange IPRange
		if c.inputStart != nil {
			ipRange.Start = c.inputStart
		}
		if c.inputEnd != nil {
			ipRange.End = c.inputEnd
		}
		require.Equal(t, c.expected, ipRange.String())
	}
}

func TestSAS(t *testing.T) {
	// Note: This is a totally invalid fake SAS, this is just testing our ability to parse different query parameters on a SAS
	const sas = "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&rscc=cc&rscd=cd&rsce=ce&rscl=cl&rsct=ct&skoid=oid&sktid=tid&skt=2111-01-09T01:42:34.936Z&ske=2222-03-09T01:42:34.936Z&sks=s&skv=v&sdd=3&saoid=oid&suoid=oid&scid=cid&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"
	_url := fmt.Sprintf("https://teststorageaccount.blob.core.windows.net/testcontainer/testpath?%s", sas)
	_uri, err := url.Parse(_url)
	require.NoError(t, err)
	sasQueryParams := NewQueryParameters(_uri.Query(), true)
	validateSAS(t, sas, sasQueryParams)
}

func validateSAS(t *testing.T, sas string, parameters QueryParameters) {
	sasCompMap := make(map[string]string)
	for _, sasComp := range strings.Split(sas, "&") {
		comp := strings.Split(sasComp, "=")
		sasCompMap[comp[0]] = comp[1]
	}

	require.Equal(t, parameters.Version(), sasCompMap["sv"])
	require.Equal(t, parameters.Services(), sasCompMap["ss"])
	require.Equal(t, parameters.ResourceTypes(), sasCompMap["srt"])
	require.Equal(t, string(parameters.Protocol()), sasCompMap["spr"])
	if _, ok := sasCompMap["st"]; ok {
		startTime, _, err := parseTime(sasCompMap["st"])
		require.NoError(t, err)
		require.Equal(t, parameters.StartTime(), startTime)
	}
	if _, ok := sasCompMap["se"]; ok {
		endTime, _, err := parseTime(sasCompMap["se"])
		require.NoError(t, err)
		require.Equal(t, parameters.ExpiryTime(), endTime)
	}

	if _, ok := sasCompMap["snapshot"]; ok {
		snapshotTime, _, err := parseTime(sasCompMap["snapshot"])
		require.NoError(t, err)
		require.Equal(t, parameters.SnapshotTime(), snapshotTime)
	}
	ipRange := parameters.IPRange()
	require.Equal(t, ipRange.String(), sasCompMap["sip"])
	require.Equal(t, parameters.Identifier(), sasCompMap["si"])
	require.Equal(t, parameters.Resource(), sasCompMap["sr"])
	require.Equal(t, parameters.Permissions(), sasCompMap["sp"])

	sign, err := url.QueryUnescape(sasCompMap["sig"])
	require.NoError(t, err)

	require.Equal(t, parameters.Signature(), sign)
	require.Equal(t, parameters.CacheControl(), sasCompMap["rscc"])
	require.Equal(t, parameters.ContentDisposition(), sasCompMap["rscd"])
	require.Equal(t, parameters.ContentEncoding(), sasCompMap["rsce"])
	require.Equal(t, parameters.ContentLanguage(), sasCompMap["rscl"])
	require.Equal(t, parameters.ContentType(), sasCompMap["rsct"])
	require.Equal(t, parameters.SignedOID(), sasCompMap["skoid"])
	require.Equal(t, parameters.SignedTID(), sasCompMap["sktid"])

	if _, ok := sasCompMap["skt"]; ok {
		signedStart, _, err := parseTime(sasCompMap["skt"])
		require.NoError(t, err)
		require.Equal(t, parameters.SignedStart(), signedStart)
	}
	require.Equal(t, parameters.SignedService(), sasCompMap["sks"])

	if _, ok := sasCompMap["ske"]; ok {
		signedExpiry, _, err := parseTime(sasCompMap["ske"])
		require.NoError(t, err)
		require.Equal(t, parameters.SignedExpiry(), signedExpiry)
	}

	require.Equal(t, parameters.SignedVersion(), sasCompMap["skv"])
	require.Equal(t, parameters.SignedDirectoryDepth(), sasCompMap["sdd"])
	require.Equal(t, parameters.AuthorizedObjectID(), sasCompMap["saoid"])
	require.Equal(t, parameters.UnauthorizedObjectID(), sasCompMap["suoid"])
	require.Equal(t, parameters.SignedCorrelationID(), sasCompMap["scid"])
}

func TestSASInvalidQueryParameter(t *testing.T) {
	// Signature is invalid below
	const sas = "sv=2019-12-12&signature=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D&sr=b"
	_url := fmt.Sprintf("https://teststorageaccount.blob.core.windows.net/testcontainer/testpath?%s", sas)
	_uri, err := url.Parse(_url)
	require.NoError(t, err)
	NewQueryParameters(_uri.Query(), true)
	// NewQueryParameters should not delete signature
	require.Contains(t, _uri.Query(), "signature")
}

func TestEncode(t *testing.T) {
	// Note: This is a totally invalid fake SAS, this is just testing our ability to parse different query parameters on a SAS
	expected := "rscc=cc&rscd=cd&rsce=ce&rscl=cl&rsct=ct&saoid=oid&scid=cid&sdd=3&se=2222-03-09T01%3A42%3A34Z&si=myIdentifier&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D&sip=168.1.5.60-168.1.5.70&ske=2222-03-09T01%3A42%3A34Z&skoid=oid&sks=s&skt=2111-01-09T01%3A42%3A34Z&sktid=tid&skv=v&sp=rw&spr=https%2Chttp&sr=b&srt=sco&ss=bf&st=2111-01-09T01%3A42%3A34Z&suoid=oid&sv=2019-12-12"
	randomOrder := "sdd=3&scid=cid&se=2222-03-09T01:42:34.936Z&rsce=ce&ss=bf&skoid=oid&si=myIdentifier&ske=2222-03-09T01:42:34.936Z&saoid=oid&sip=168.1.5.60-168.1.5.70&rscc=cc&srt=sco&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D&rsct=ct&skt=2111-01-09T01:42:34.936Z&rscl=cl&suoid=oid&sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&rscd=cd&sp=rw&sktid=tid&spr=https,http&sks=s&skv=v"
	testdata := []string{expected, randomOrder}

	for _, sas := range testdata {
		_url := fmt.Sprintf("https://teststorageaccount.blob.core.windows.net/testcontainer/testpath?%s", sas)
		_uri, err := url.Parse(_url)
		require.NoError(t, err)
		queryParams := NewQueryParameters(_uri.Query(), true)
		require.Equal(t, expected, queryParams.Encode())
	}
}
