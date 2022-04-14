//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func validateSAS(_assert *assert.Assertions, sas string, parameters SASQueryParameters) {
	sasCompMap := make(map[string]string)
	for _, sasComp := range strings.Split(sas, "&") {
		comp := strings.Split(sasComp, "=")
		sasCompMap[comp[0]] = comp[1]
	}

	_assert.Equal(parameters.Version(), sasCompMap["sv"])
	_assert.Equal(parameters.Services(), sasCompMap["ss"])
	_assert.Equal(parameters.ResourceTypes(), sasCompMap["srt"])
	_assert.Equal(string(parameters.Protocol()), sasCompMap["spr"])
	if _, ok := sasCompMap["st"]; ok {
		startTime, _, err := parseSASTimeString(sasCompMap["st"])
		_assert.Nil(err)
		_assert.Equal(parameters.StartTime(), startTime)
	}
	if _, ok := sasCompMap["se"]; ok {
		endTime, _, err := parseSASTimeString(sasCompMap["se"])
		_assert.Nil(err)
		_assert.Equal(parameters.ExpiryTime(), endTime)
	}

	if _, ok := sasCompMap["snapshot"]; ok {
		snapshotTime, _, err := parseSASTimeString(sasCompMap["snapshot"])
		_assert.Nil(err)
		_assert.Equal(parameters.SnapshotTime(), snapshotTime)
	}
	ipRange := parameters.IPRange()
	_assert.Equal(ipRange.String(), sasCompMap["sip"])
	_assert.Equal(parameters.Identifier(), sasCompMap["si"])
	_assert.Equal(parameters.Resource(), sasCompMap["sr"])
	_assert.Equal(parameters.Permissions(), sasCompMap["sp"])

	sign, err := url.QueryUnescape(sasCompMap["sig"])
	_assert.Nil(err)

	_assert.Equal(parameters.Signature(), sign)
	_assert.Equal(parameters.CacheControl(), sasCompMap["rscc"])
	_assert.Equal(parameters.ContentDisposition(), sasCompMap["rscd"])
	_assert.Equal(parameters.ContentEncoding(), sasCompMap["rsce"])
	_assert.Equal(parameters.ContentLanguage(), sasCompMap["rscl"])
	_assert.Equal(parameters.ContentType(), sasCompMap["rsct"])
	_assert.Equal(parameters.signedOid, sasCompMap["skoid"])
	_assert.Equal(parameters.SignedTid(), sasCompMap["sktid"])

	if _, ok := sasCompMap["skt"]; ok {
		signedStart, _, err := parseSASTimeString(sasCompMap["skt"])
		_assert.Nil(err)
		_assert.Equal(parameters.SignedStart(), signedStart)
	}
	_assert.Equal(parameters.SignedService(), sasCompMap["sks"])

	if _, ok := sasCompMap["ske"]; ok {
		signedExpiry, _, err := parseSASTimeString(sasCompMap["ske"])
		_assert.Nil(err)
		_assert.Equal(parameters.SignedExpiry(), signedExpiry)
	}

	_assert.Equal(parameters.SignedVersion(), sasCompMap["skv"])
	_assert.Equal(parameters.SignedDirectoryDepth(), sasCompMap["sdd"])
	_assert.Equal(parameters.PreauthorizedAgentObjectId(), sasCompMap["saoid"])
	_assert.Equal(parameters.AgentObjectId(), sasCompMap["suoid"])
	_assert.Equal(parameters.SignedCorrelationId(), sasCompMap["scid"])
}

func (s *azblobTestSuite) TestSASGeneration() {
	_assert := assert.New(s.T())
	sas := "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"
	_url := fmt.Sprintf("https://teststorageaccount.blob.core.windows.net/testcontainer/testpath?%s", sas)
	_uri, err := url.Parse(_url)
	_assert.Nil(err)
	sasQueryParams := newSASQueryParameters(_uri.Query(), true)
	validateSAS(_assert, sas, sasQueryParams)
}
