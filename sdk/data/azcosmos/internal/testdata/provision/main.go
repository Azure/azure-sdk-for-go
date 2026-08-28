// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Command provision creates the database and container the emulator tests operate on.
//
// Container management is not bound in v2 yet, so those tests take an existing container rather
// than creating one. This talks to the REST API directly instead of going through either SDK, so
// that setting up the tests does not depend on the code they are meant to exercise.
//
// It is idempotent: a resource that already exists is treated as success, so it can run before
// every job without needing a teardown.
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// emulatorKey is the emulator's well-known account key, which is published in its documentation
// and grants nothing anywhere else.
const emulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

// partitionKeyPath is the path the emulator tests build their partition key values against.
const partitionKeyPath = "/pk"

func main() {
	endpoint := envOr("AZCOSMOS_ENDPOINT", "http://localhost:8081")
	database := envOr("AZCOSMOS_DATABASE", "itemdb")
	container := envOr("AZCOSMOS_CONTAINER", "items")

	key, err := base64.StdEncoding.DecodeString(emulatorKey)
	if err != nil {
		fail("decoding the emulator key: %v", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		// The emulator serves a self-signed certificate over https. This is a local provisioning
		// helper for a well-known key, not part of the SDK, and it never talks to a real account.
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, //nolint:gosec // emulator only
	}

	if err := create(client, key, endpoint, "dbs", "", "dbs",
		fmt.Sprintf(`{"id":%q}`, database)); err != nil {
		fail("creating database %q: %v", database, err)
	}
	fmt.Printf("database %q ready\n", database)

	body := fmt.Sprintf(
		`{"id":%q,"partitionKey":{"paths":[%q],"kind":"Hash","version":2}}`,
		container, partitionKeyPath)
	if err := create(client, key, endpoint, "colls", "dbs/"+database,
		"dbs/"+database+"/colls", body); err != nil {
		fail("creating container %q: %v", container, err)
	}
	fmt.Printf("container %q ready\n", container)
}

// create issues one resource-creation request, treating an existing resource as success.
func create(client *http.Client, key []byte, endpoint, resourceType, resourceLink, path, body string) error {
	date := time.Now().UTC().Format(http.TimeFormat)

	request, err := http.NewRequest(http.MethodPost,
		strings.TrimRight(endpoint, "/")+"/"+path, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", authorization(key, "post", resourceType, resourceLink, date))
	request.Header.Set("x-ms-date", date)
	request.Header.Set("x-ms-version", "2018-12-31")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	payload, _ := io.ReadAll(response.Body)
	switch response.StatusCode {
	case http.StatusCreated, http.StatusOK, http.StatusConflict:
		// Conflict means an earlier run already provisioned it.
		return nil
	default:
		return fmt.Errorf("status %d: %s", response.StatusCode, payload)
	}
}

// authorization builds the master-key signature the REST API expects: an HMAC over the lowercased
// verb, resource type, resource link and date.
func authorization(key []byte, method, resourceType, resourceLink, date string) string {
	payload := method + "\n" + resourceType + "\n" + resourceLink + "\n" + strings.ToLower(date) + "\n\n"

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return url.QueryEscape("type=master&ver=1.0&sig=" + signature)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
