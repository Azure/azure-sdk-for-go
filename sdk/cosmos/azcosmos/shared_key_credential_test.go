// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCanonicalizedAuthHeader(t *testing.T) {

	account := "someAccount"
	key := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	cred, err := NewSharedKeyCredential(account, key)

	assert.NoError(t, err)

	method := "GET"
	resourceType := "dbs"
	resourceId := "dbs/testdb"
	xmsDate := "Thu, 27 Apr 2017 00:51:12 GMT"
	tokenType := "master"
	version := "1.0"

	emptyAuthHeader := cred.buildCanonicalizedAuthHeader("", resourceType, resourceId, xmsDate, tokenType, version)
	assert.Equal(t, emptyAuthHeader, "")
	emptyAuthHeader = cred.buildCanonicalizedAuthHeader(method, "", resourceId, xmsDate, tokenType, version)
	assert.Equal(t, emptyAuthHeader, "")
	emptyAuthHeader = cred.buildCanonicalizedAuthHeader(method, resourceType, "", xmsDate, tokenType, version)
	assert.Equal(t, emptyAuthHeader, "")

	stringToSign := strings.ToLower(join(method, "\n", resourceType, "\n", resourceId, "\n", xmsDate, "\n", "", "\n"))
	signature := cred.computeHMACSHA256(stringToSign)
	expected := url.QueryEscape(fmt.Sprintf("type=%s&ver=%s&sig=%s", tokenType, version, signature))

	authHeader := cred.buildCanonicalizedAuthHeader(method, resourceType, resourceId, xmsDate, tokenType, version)

	assert.GreaterOrEqual(t, len(authHeader), 1)
	assert.Equal(t, authHeader, expected)
}
