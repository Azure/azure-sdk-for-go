// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"slices"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

func TestSpanForClient(t *testing.T) {
	endpoint, _ := url.Parse("https://localhost:8081/")
	aSpan, err := getSpanNameForClient(endpoint, operationTypeQuery, resourceTypeDatabase, "test")
	if err != nil {
		t.Fatalf("Failed to get span name: %v", err)
	}
	if aSpan.name != "query_databases test" {
		t.Fatalf("Expected span name to be 'query_databases test', but got %s", aSpan.name)
	}
	if len(aSpan.options.Attributes) == 0 {
		t.Fatalf("Expected span options to have attributes, but got none")
	}

	idx := slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.system" && a.Value == "cosmosdb" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.system' with value 'cosmosdb', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.cosmosdb.connection_mode" && a.Value == "gateway" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.cosmosdb.connection_mode' with value 'gateway', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.address" && a.Value == "localhost" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.address' with value 'localhost', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.port" && a.Value == "8081" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.port' with value '8081', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.operation.name" && a.Value == "query_databases" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.operation.name' with value 'query_databases', but got none")
	}

	aSpan, err = getSpanNameForClient(endpoint, operationTypeCreate, resourceTypeDatabase, "test")
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
}

func TestSpanForDatabases(t *testing.T) {
	endpoint, _ := url.Parse("https://localhost:8081/")
	aSpan, err := getSpanNameForDatabases(endpoint, operationTypeCreate, resourceTypeDatabase, "test")
	if err != nil {
		t.Fatalf("Failed to get span name: %v", err)
	}
	if aSpan.name != "create_database test" {
		t.Fatalf("Expected span name to be 'create_database test', but got %s", aSpan.name)
	}
	if len(aSpan.options.Attributes) == 0 {
		t.Fatalf("Expected span options to have attributes, but got none")
	}

	idx := slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.system" && a.Value == "cosmosdb" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.system' with value 'cosmosdb', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.cosmosdb.connection_mode" && a.Value == "gateway" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.cosmosdb.connection_mode' with value 'gateway', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.address" && a.Value == "localhost" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.address' with value 'localhost', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.port" && a.Value == "8081" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.port' with value '8081', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.operation.name" && a.Value == "create_database" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.operation.name' with value 'create_database', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.namespace" && a.Value == "test" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.namespace' with value 'test', but got none")
	}

	aSpan, err = getSpanNameForDatabases(endpoint, operationTypeCreate, resourceTypeCollection, "test")
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
}

func TestSpanForContainers(t *testing.T) {
	endpoint, _ := url.Parse("https://localhost:8081/")
	aSpan, err := getSpanNameForContainers(endpoint, operationTypeCreate, resourceTypeCollection, "db", "test")
	if err != nil {
		t.Fatalf("Failed to get span name: %v", err)
	}
	if aSpan.name != "create_container test" {
		t.Fatalf("Expected span name to be 'create_container test', but got %s", aSpan.name)
	}
	if len(aSpan.options.Attributes) == 0 {
		t.Fatalf("Expected span options to have attributes, but got none")
	}

	idx := slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.system" && a.Value == "cosmosdb" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.system' with value 'cosmosdb', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.cosmosdb.connection_mode" && a.Value == "gateway" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.cosmosdb.connection_mode' with value 'gateway', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.address" && a.Value == "localhost" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.address' with value 'localhost', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "server.port" && a.Value == "8081" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'server.port' with value '8081', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.operation.name" && a.Value == "create_container" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.operation.name' with value 'create_container', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.namespace" && a.Value == "db" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.namespace' with value 'db', but got none")
	}

	idx = slices.IndexFunc(aSpan.options.Attributes, func(a tracing.Attribute) bool { return a.Key == "db.collection.name" && a.Value == "test" })
	if idx == -1 {
		t.Fatalf("Expected attribute 'db.collection.name' with value 'test', but got none")
	}

	aSpan, err = getSpanNameForContainers(endpoint, operationTypeCreate, resourceTypeDatabase, "db", "test")
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}
}
