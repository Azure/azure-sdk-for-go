// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	chk "gopkg.in/check.v1"
)

// For testing docs, see: https://labix.org/gocheck
// To test a specific test: go test -check.f MyTestSuite

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

type aztestsSuite struct{}

var _ = chk.Suite(&aztestsSuite{})

const (
	storageAccountNameEnvVar    = "TABLES_STORAGE_ACCOUNT_NAME"
	cosmosAccountNameEnnVar     = "TABLES_COSMOS_ACCOUNT_NAME"
	accountKeyEnvVar            = "TABLES_PRIMARY_STORAGE_ACCOUNT_KEY"
	storageEndpointSuffixEnvVar = "STORAGE_ENDPOINT_SUFFIX"
	cosmosEndpointSuffixEnvVar  = "COSMOS_TABLES_ENDPOINT_SUFFIX"
	storageAccountKeyEnvVar     = "TABLES_PRIMARY_STORAGE_ACCOUNT_KEY"
	cosmosAccountKeyEnvVar      = "TABLES_PRIMARY_COSMOS_ACCOUNT_KEY"
	tableNamePrefix             = "gotable"
	DefaultStorageSuffix        = "core.windows.net"
	DefaultCosmosSuffix         = "cosmos.azure.com"
)

type EndpointType string

const (
	StorageEndpoint EndpointType = "storage"
	CosmosEndpoint  EndpointType = "cosmos"
)

var ctx = context.Background()

func getRequiredEnv(name string) string {
	env, ok := os.LookupEnv(name)
	if ok {
		return env
	} else {
		panic("Required environment variable not set: " + name)
	}
}

func storageURI() string {
	return "https://" + storageAccountName() + ".table." + storageEndpointSuffix()
}

func cosmosURI() string {
	return "https://" + cosmosAccountName() + ".table" + cosmosAccountName()
}

func storageAccountName() string {
	return getRequiredEnv(storageAccountNameEnvVar)
}

func cosmosAccountName() string {
	return getRequiredEnv(cosmosAccountNameEnnVar)
}

func cosmosAccountKey() string {
	return getRequiredEnv(cosmosAccountKeyEnvVar)
}

func storageAccountKey() string {
	return getRequiredEnv(storageAccountKeyEnvVar)
}

func storageEndpointSuffix() string {
	suffix, ok := os.LookupEnv(storageEndpointSuffixEnvVar)
	if ok {
		return suffix
	} else {
		return DefaultStorageSuffix
	}
}

func cosmosEndpointSuffix() string {
	suffix, ok := os.LookupEnv(cosmosEndpointSuffix())
	if ok {
		return suffix
	} else {
		return DefaultCosmosSuffix
	}
}

func createTableClient(endpointType EndpointType) (TableClient, error) {
	if endpointType == StorageEndpoint {
		storageCred, _ := NewSharedKeyCredential(storageAccountName(), storageAccountKey())
		return NewTableClient(storageURI(), storageCred, nil)
	} else {
		cosmosCred, _ := NewSharedKeyCredential(cosmosAccountName(), cosmosAccountKey())
		return NewTableClient(cosmosURI(), cosmosCred, nil)
	}
}

func getGenericCredential(accountType string) (*SharedKeyCredential, error) {

	accountName, accountKey := getRequiredEnv(storageAccountNameEnvVar), getRequiredEnv(accountKeyEnvVar)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(storageAccountNameEnvVar + " and/or " + accountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

func generateName() string {
	currentTime := time.Now()
	name := fmt.Sprintf("%s%d%d%d", tableNamePrefix, currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond())
	return name
}
