//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
)

func validateSAS(_require *require.Assertions, sas string, parameters blob.SASQueryParameters) {
	sasCompMap := make(map[string]string)
	for _, sasComp := range strings.Split(sas, "&") {
		comp := strings.Split(sasComp, "=")
		sasCompMap[comp[0]] = comp[1]
	}

	_require.Equal(parameters.Version(), sasCompMap["sv"])
	_require.Equal(parameters.Services(), sasCompMap["ss"])
	_require.Equal(parameters.ResourceTypes(), sasCompMap["srt"])
	_require.Equal(string(parameters.Protocol()), sasCompMap["spr"])
	if _, ok := sasCompMap["st"]; ok {
		startTime, _, err := blob.ParseSASTimeString(sasCompMap["st"])
		_require.Nil(err)
		_require.Equal(parameters.StartTime(), startTime)
	}
	if _, ok := sasCompMap["se"]; ok {
		endTime, _, err := blob.ParseSASTimeString(sasCompMap["se"])
		_require.Nil(err)
		_require.Equal(parameters.ExpiryTime(), endTime)
	}

	if _, ok := sasCompMap["snapshot"]; ok {
		snapshotTime, _, err := blob.ParseSASTimeString(sasCompMap["snapshot"])
		_require.Nil(err)
		_require.Equal(parameters.SnapshotTime(), snapshotTime)
	}
	ipRange := parameters.IPRange()
	_require.Equal(ipRange.String(), sasCompMap["sip"])
	_require.Equal(parameters.Identifier(), sasCompMap["si"])
	_require.Equal(parameters.Resource(), sasCompMap["sr"])
	_require.Equal(parameters.Permissions(), sasCompMap["sp"])

	sign, err := url.QueryUnescape(sasCompMap["sig"])
	_require.Nil(err)

	_require.Equal(parameters.Signature(), sign)
	_require.Equal(parameters.CacheControl(), sasCompMap["rscc"])
	_require.Equal(parameters.ContentDisposition(), sasCompMap["rscd"])
	_require.Equal(parameters.ContentEncoding(), sasCompMap["rsce"])
	_require.Equal(parameters.ContentLanguage(), sasCompMap["rscl"])
	_require.Equal(parameters.ContentType(), sasCompMap["rsct"])
	_require.Equal(parameters.SignedOid(), sasCompMap["skoid"])
	_require.Equal(parameters.SignedTid(), sasCompMap["sktid"])

	if _, ok := sasCompMap["skt"]; ok {
		signedStart, _, err := blob.ParseSASTimeString(sasCompMap["skt"])
		_require.Nil(err)
		_require.Equal(parameters.SignedStart(), signedStart)
	}
	_require.Equal(parameters.SignedService(), sasCompMap["sks"])

	if _, ok := sasCompMap["ske"]; ok {
		signedExpiry, _, err := blob.ParseSASTimeString(sasCompMap["ske"])
		_require.Nil(err)
		_require.Equal(parameters.SignedExpiry(), signedExpiry)
	}

	_require.Equal(parameters.SignedVersion(), sasCompMap["skv"])
	_require.Equal(parameters.SignedDirectoryDepth(), sasCompMap["sdd"])
	_require.Equal(parameters.PreauthorizedAgentObjectId(), sasCompMap["saoid"])
	_require.Equal(parameters.AgentObjectId(), sasCompMap["suoid"])
	_require.Equal(parameters.SignedCorrelationId(), sasCompMap["scid"])
}

func (s *azblobTestSuite) TestSASGeneration() {
	_require := require.New(s.T())
	sas := "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"
	_url := fmt.Sprintf("https://teststorageaccount.blob.core.windows.net/testcontainer/testpath?%s", sas)
	_uri, err := url.Parse(_url)
	_require.Nil(err)
	sasQueryParams := exported.NewSASQueryParameters(_uri.Query(), true)
	validateSAS(_require, sas, sasQueryParams)
}
