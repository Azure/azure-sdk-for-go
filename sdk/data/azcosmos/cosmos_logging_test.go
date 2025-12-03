// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestLoggingSuccessResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusOK),
		mock.WithBody([]byte(`{}`)), // Minimal valid JSON
	)

	logger := &captureLogger{}
	azlog.SetListener(func(event azlog.Event, message string) {
		logger.log(event, message)
	})
	defer azlog.SetListener(nil)

	azlog.SetEvents(azlog.EventResponse)
	defer azlog.SetEvents()

	validKey := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	cred, _ := NewKeyCredential(validKey)
	// We don't care if NewClientWithKey fails due to invalid JSON body for account properties,
	// we just want to see the HTTP 200 log.
	client, _ := NewClientWithKey(srv.URL(), cred, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: srv,
			Retry: policy.RetryOptions{
				MaxRetries: -1,
			},
		},
	})

	client.CreateDatabase(context.Background(), DatabaseProperties{ID: "testdb"}, nil)

	logs := logger.filterByEvent(azlog.EventResponse)
	found := false
	for _, logMsg := range logs {
		if strings.Contains(logMsg, "RESPONSE Status: 200") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 200 OK response log, but not found. Logs:\n%s", logger.getAllLogs())
	}
}

func TestLoggingFailureResponse(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithStatusCode(http.StatusInternalServerError),
		mock.WithBody([]byte(`{"code":"InternalServerError"}`)),
	)

	logger := &captureLogger{}
	azlog.SetListener(func(event azlog.Event, message string) {
		logger.log(event, message)
	})
	defer azlog.SetListener(nil)

	azlog.SetEvents(azlog.EventResponse)
	defer azlog.SetEvents()

	validKey := "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
	cred, _ := NewKeyCredential(validKey)
	client, _ := NewClientWithKey(srv.URL(), cred, &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: srv,
			Retry: policy.RetryOptions{
				MaxRetries: -1,
			},
		},
	})

	client.CreateDatabase(context.Background(), DatabaseProperties{ID: "testdb"}, nil)

	logs := logger.filterByEvent(azlog.EventResponse)
	found := false
	for _, logMsg := range logs {
		if strings.Contains(logMsg, "RESPONSE Status: 500") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 500 Internal Server Error response log, but not found. Logs:\n%s", logger.getAllLogs())
	}
}
