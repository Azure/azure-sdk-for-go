//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type userAgentValidatingPolicy struct {
	t     *testing.T
	appID string
}

func (p userAgentValidatingPolicy) Do(req *policy.Request) (*http.Response, error) {
	expected := "azsdk-go-" + component + "/" + version
	if p.appID != "" {
		expected = p.appID + " " + expected
	}
	if ua := req.Raw().Header.Get("User-Agent"); !strings.HasPrefix(ua, expected) {
		p.t.Fatalf("unexpected User-Agent %s", ua)
	}
	return req.Next()
}

func TestIMDSEndpointParse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}

func TestManagedIdentityClient_UserAgent(t *testing.T) {
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t}},
		},
	}
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagedIdentityClient_ApplicationID(t *testing.T) {
	appID := "customvalue"
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t, appID: appID}},
		},
	}
	options.Telemetry.ApplicationID = appID
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagedIdentityClient_IMDSErrors(t *testing.T) {
	for _, test := range []struct {
		body, desc string
		code       int
	}{
		{
			desc: "No identity assigned",
			code: http.StatusBadRequest,
			body: `{"error":"invalid_request","error_description":"Identity not found"}`,
		},
		{
			desc: "Docker Desktop",
			code: http.StatusForbidden,
			body: "connecting to 169.254.169.254:80: connecting to 169.254.169.254:80: dial tcp 169.254.169.254:80: connectex: A socket operation was attempted to an unreachable network.",
		},
		{
			desc: "Docker Desktop",
			code: http.StatusForbidden,
			body: "connecting to 169.254.169.254:80: connecting to 169.254.169.254:80: dial tcp 169.254.169.254:80: connectex: A socket operation was attempted to an unreachable host.",
		},
	} {
		t.Run(fmt.Sprint(test.code), func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.SetResponse(mock.WithBody([]byte(test.body)), mock.WithStatusCode(test.code))
			client, err := newManagedIdentityClient(&ManagedIdentityCredentialOptions{
				ClientOptions: azcore.ClientOptions{Transport: srv},
			})
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.authenticate(context.Background(), nil, testTRO.Scopes)
			if err == nil {
				t.Fatal("expected an error")
			}
			if actual := err.Error(); !strings.Contains(actual, test.body) {
				t.Fatalf("expected response body in error, got %q", actual)
			}
			var unavailableErr credentialUnavailable
			if !errors.As(err, &unavailableErr) {
				t.Fatalf("expected %T, got %T", unavailableErr, err)
			}
		})
	}
}

func TestManagedIdentityClient_UserAssignedIDWarning(t *testing.T) {
	for _, test := range []struct {
		name          string
		createRequest func(*managedIdentityClient) error
	}{
		{
			name: "Azure Arc",
			createRequest: func(client *managedIdentityClient) error {
				_, err := client.createAzureArcAuthRequest(context.Background(), client.id, []string{liveTestScope}, "key")
				return err
			},
		},
		{
			name: "Cloud Shell",
			createRequest: func(client *managedIdentityClient) error {
				_, err := client.createCloudShellAuthRequest(context.Background(), client.id, []string{liveTestScope})
				return err
			},
		},
		{
			name: "Service Fabric",
			createRequest: func(client *managedIdentityClient) error {
				_, err := client.createServiceFabricAuthRequest(context.Background(), client.id, []string{liveTestScope})
				return err
			},
		},
	} {
		for _, id := range []ManagedIDKind{ClientID(fakeClientID), ResourceID(fakeResourceID)} {
			s := "-ClientID"
			if id.String() == fakeResourceID {
				s = "-ResourceID"
			}
			t.Run(test.name+s, func(t *testing.T) {
				msgs := []string{}
				log.SetListener(func(event log.Event, msg string) {
					if event == EventAuthentication {
						msgs = append(msgs, msg)
					}
				})
				client, err := newManagedIdentityClient(&ManagedIdentityCredentialOptions{
					ID: id,
				})
				if err != nil {
					t.Fatal(err)
				}
				err = test.createRequest(client)
				if err != nil {
					t.Fatal(err)
				}
				for _, msg := range msgs {
					if strings.Contains(msg, test.name) && strings.Contains(msg, "user-assigned") {
						return
					}
				}
				t.Fatalf("expected warning about user-assigned ID, got:\n%s", strings.Join(msgs, "\n"))
			})
		}
	}
}
