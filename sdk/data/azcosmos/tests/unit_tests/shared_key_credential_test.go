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

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/assert"
)

func Test_buildCanonicalizedAuthHeader(t *testing.T) {
	key := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	cred, err := NewKeyCredential(key)

	assert.NoError(t, err)

	method := "GET"
	resourceType := "dbs"
	resourceId := "dbs/testdb"
	xmsDate := "Thu, 27 Apr 2017 00:51:12 GMT"
	tokenType := "master"
	version := "1.0"

	emptyAuthHeader := cred.buildCanonicalizedAuthHeader(false, "", resourceType, resourceId, xmsDate, tokenType, version)
	assert.Equal(t, emptyAuthHeader, "")
	emptyAuthHeader = cred.buildCanonicalizedAuthHeader(false, method, "", resourceId, xmsDate, tokenType, version)
	assert.Equal(t, emptyAuthHeader, "")

	stringToSign := join(strings.ToLower(method), "\n", strings.ToLower(resourceType), "\n", resourceId, "\n", strings.ToLower(xmsDate), "\n", "", "\n")
	signature := cred.computeHMACSHA256(stringToSign)
	expected := url.QueryEscape(fmt.Sprintf("type=%s&ver=%s&sig=%s", tokenType, version, signature))

	authHeader := cred.buildCanonicalizedAuthHeader(false, method, resourceType, resourceId, xmsDate, tokenType, version)

	assert.GreaterOrEqual(t, len(authHeader), 1)
	assert.Equal(t, expected, authHeader)
}

func Test_buildCanonicalizedAuthHeaderFromRequest(t *testing.T) {
	key := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	cred, err := NewKeyCredential(key)

	assert.NoError(t, err)

	method := "GET"
	resourceType := "dbs"
	resourceId := "dbs/testdb"
	xmsDate := "Thu, 27 Apr 2017 00:51:12 GMT"
	tokenType := "master"
	version := "1.0"

	stringToSign := join(strings.ToLower(method), "\n", strings.ToLower(resourceType), "\n", resourceId, "\n", strings.ToLower(xmsDate), "\n", "", "\n")
	signature := cred.computeHMACSHA256(stringToSign)
	expected := url.QueryEscape(fmt.Sprintf("type=%s&ver=%s&sig=%s", tokenType, version, signature))

	req, _ := azruntime.NewRequest(context.TODO(), http.MethodGet, "http://localhost")
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "dbs/testdb",
	}

	req.Raw().Header.Set(headerXmsDate, xmsDate)
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.SetOperationValue(operationContext)
	authHeader, _ := cred.buildCanonicalizedAuthHeaderFromRequest(req)

	assert.Equal(t, expected, authHeader)
}

func Test_buildCanonicalizedAuthHeaderFromRequestWithRid(t *testing.T) {
	key := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	cred, err := NewKeyCredential(key)

	assert.NoError(t, err)

	method := "GET"
	resourceType := "dbs"
	resourceId := "dbs/rid"
	xmsDate := "Thu, 27 Apr 2017 00:51:12 GMT"
	tokenType := "master"
	version := "1.0"

	stringToSign := join(strings.ToLower(method), "\n", strings.ToLower(resourceType), "\n", resourceId, "\n", strings.ToLower(xmsDate), "\n", "", "\n")
	signature := cred.computeHMACSHA256(stringToSign)
	expected := url.QueryEscape(fmt.Sprintf("type=%s&ver=%s&sig=%s", tokenType, version, signature))

	req, _ := azruntime.NewRequest(context.TODO(), http.MethodGet, "http://localhost")
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "dbs/Rid",
		isRidBased:      true,
	}

	req.Raw().Header.Set(headerXmsDate, xmsDate)
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.SetOperationValue(operationContext)
	authHeader, _ := cred.buildCanonicalizedAuthHeaderFromRequest(req)

	assert.Equal(t, expected, authHeader)
}

func Test_buildCanonicalizedAuthHeaderFromRequestWithEscapedCharacters(t *testing.T) {
	key := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

	cred, err := NewKeyCredential(key)

	assert.NoError(t, err)

	method := "GET"
	resourceType := "dbs"
	originalResourceId := "dbs/name with spaces"
	resourceId := url.PathEscape(originalResourceId)
	xmsDate := "Thu, 27 Apr 2017 00:51:12 GMT"
	tokenType := "master"
	version := "1.0"

	stringToSign := join(strings.ToLower(method), "\n", strings.ToLower(resourceType), "\n", originalResourceId, "\n", strings.ToLower(xmsDate), "\n", "", "\n")
	signature := cred.computeHMACSHA256(stringToSign)
	expected := url.QueryEscape(fmt.Sprintf("type=%s&ver=%s&sig=%s", tokenType, version, signature))

	req, _ := azruntime.NewRequest(context.TODO(), http.MethodGet, "http://localhost")
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: resourceId,
	}

	req.Raw().Header.Set(headerXmsDate, xmsDate)
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.SetOperationValue(operationContext)
	authHeader, _ := cred.buildCanonicalizedAuthHeaderFromRequest(req)

	assert.Equal(t, expected, authHeader)
}
