// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

func Test_buildCanonicalizedAuthHeaderFromRequest(t *testing.T) {
	emulatorTests := newEmulatorTests()
	cred, _ := NewSharedKeyCredential(emulatorTests.key)
	req, _ := azcore.NewRequest(context.TODO(), http.MethodPost, emulatorTests.host+"dbs")
	req.SetOperationValue(cosmosOperationContext{resourceType: resourceTypeDatabase, resourceAddress: ""})

	req.Request.Header.Set(azcore.HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))

	authHeader, _ := cred.buildCanonicalizedAuthHeaderFromRequest(req)

	req.Request.Header.Set(azcore.HeaderAuthorization, authHeader)
	req.Request.Header.Set(azcore.HeaderXmsVersion, "2020-11-05")

	db := &CosmosDatabaseProperties{Id: "testdb"}
	req.MarshalAsJSON(db)

	assert.NotEqual(t, "", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req.Request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.LessOrEqual(t, resp.StatusCode, 400, string(body))
}
