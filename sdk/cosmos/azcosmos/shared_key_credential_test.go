// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

func getTestSharedKeyCredentialInfo() (endpoint string, key string) {
	return "https://localhost:8081", "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
}

func Test_buildCanonicalizedAuthHeader(t *testing.T) {
	_, key := getTestSharedKeyCredentialInfo()

	cred, err := NewSharedKeyCredential(key)

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
	assert.Equal(t, expected, authHeader)
}

func Test_buildCanonicalizedAuthHeaderFromRequest(t *testing.T) {
	endpoint, key := getTestSharedKeyCredentialInfo()
	cred, _ := NewSharedKeyCredential(key)
	req, _ := azcore.NewRequest(context.TODO(), http.MethodPut, endpoint)
	req.SetOperationValue(cosmosOperationContext{resourceType: "dbs", resourceId: "dbs/testdb"})
	authHeader, _ := cred.buildCanonicalizedAuthHeaderFromRequest(req)

	assert.NotEqual(t, "", authHeader)
}
