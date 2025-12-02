// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type emulatorTests struct {
	host string
	key  string
}

func newEmulatorTests(t *testing.T) *emulatorTests {
	return newEmulatorTestsWithEndpoint(t, "https://localhost:8081/")
}

func newEmulatorTestsWithComputeGateway(t *testing.T) *emulatorTests {
	return newEmulatorTestsWithEndpoint(t, "https://localhost:8903/")
}

func newEmulatorTestsWithEndpoint(t *testing.T, e string) *emulatorTests {
	envCheck := os.Getenv("EMULATOR")
	if envCheck == "" {
		t.Skip("set EMULATOR environment variable to run this test")
	}

	return &emulatorTests{
		host: e,
		key:  "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
	}
}

func (e *emulatorTests) getClient(t *testing.T, tp tracing.Provider) *Client {
	cred, _ := NewKeyCredential(e.key)

	// Create a client with a custom transport that skips TLS verification
	// Since there's a self-signed certificate in the emulator, we need to skip verification
	transport := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
		Transport:       transport,
	}}

	client, err := NewClientWithKey(e.host, cred, options)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func (e *emulatorTests) getAadClient(t *testing.T, tp tracing.Provider) *Client {
	cred := &emulatorTokenCredential{}
	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
	}}
	client, err := NewClient(e.host, cred, options)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func (e *emulatorTests) createDatabase(
	t *testing.T,
	ctx context.Context,
	client *Client,
	dbName string) *DatabaseClient {
	database := DatabaseProperties{ID: dbName}
	resp, err := client.CreateDatabase(ctx, database, nil)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if resp.DatabaseProperties.ID != database.ID {
		t.Errorf("Unexpected id match: %v", resp.DatabaseProperties)
	}

	db, _ := client.NewDatabase(dbName)
	return db
}

func (e *emulatorTests) deleteDatabase(
	t *testing.T,
	ctx context.Context,
	database *DatabaseClient) {
	_, err := database.Delete(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to delete database: %v", err)
	}
}

func (e *emulatorTests) marshallItem(id string, pk string) []byte {
	item := map[string]string{
		"id": id,
		"pk": pk,
	}

	marshalled, _ := json.Marshal(item)
	return marshalled
}

type emulatorTokenCredential struct {
}

func (c *emulatorTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	header := `{"typ":"JWT","alg":"RS256","x5t":"CosmosEmulatorPrimaryMaster","kid":"CosmosEmulatorPrimaryMaster"}`
	unixNow := time.Now().Unix()
	expiration := unixNow + 7200
	payload := `{ 
		"appid":"localhost", 
		"aio":"",
		"appidacr":"1",
		"idp": "https://localhost:8081/",
		"oid": "96313034-4739-43cb-93cd-74193adbe5b6",
		"rh": "",
		"sub": "localhost",
		"tid": "EmulatorFederation",
		"uti": "",
		"ver": "1.0",
		"scp": "user_impersonation",
		"groups":[ 
			"7ce1d003-4cb3-4879-b7c5-74062a35c66e",
			"e99ff30c-c229-4c67-ab29-30a6aebc3e58",
			"5549bb62-c77b-4305-bda9-9ec66b85d9e4",
			"c44fd685-5c58-452c-aaf7-13ce75184f65",
			"be895215-eab5-43b7-9536-9ef8fe130330"], 
		"nbf":` + strconv.FormatInt(unixNow, 10) + `, 
		"exp":` + strconv.FormatInt(expiration, 10) + `, 
		"iat":` + strconv.FormatInt(unixNow, 10) + `,
		"iss":"https://sts.fake-issuer.net/7b1999a1-dfd7-440e-8204-00170979b984",
		"aud":"https://localhost.localhost" 
	}`

	headerBase64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	payloadBase64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	masterKeyBase64 := base64.RawURLEncoding.EncodeToString([]byte("C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="))

	token := headerBase64 + "." + payloadBase64 + "." + masterKeyBase64

	return azcore.AccessToken{
		Token:     token,
		ExpiresOn: time.Unix(expiration, 0),
	}, nil
}
