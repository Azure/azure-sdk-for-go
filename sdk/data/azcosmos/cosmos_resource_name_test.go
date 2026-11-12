// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

func TestResourceNameValidation(t *testing.T) {
	client := &Client{enableResourceNameValidation: true}
	for _, name := range []string{
		"", ".", "..", "other/item", `other\item`, "item?query", "item#fragment",
		"\t", "\n", "\r", ".\t.", ".\n.", ".\r.", "tar\tget", "tar\nget", "tar\rget",
		"../../other-db/colls/other-container/docs/target", string([]byte{0xff}),
	} {
		t.Run(name, func(t *testing.T) {
			for _, resourceType := range []string{"database", "container", "item"} {
				require.ErrorContains(t, client.validateResourceName(resourceType, name), "invalid "+resourceType+" ID")
				disabled := &Client{}
				require.NoError(t, disabled.validateResourceName(resourceType, name))
			}
		})
	}
	for _, name := range []string{
		"item", "item with spaces", "item+value=1", "item%41", "item%2Fname",
		"item..backup", "...", `item"quoted`, `allowed","id":"victim`, "\u4e16\u754c", "\ufffd",
	} {
		require.NoError(t, client.validateResourceName("item", name), "name: %q", name)
	}
}

func TestResourceNameValidationItemBody(t *testing.T) {
	client := &Client{enableResourceNameValidation: true}
	for _, body := range []string{
		`{"id":"safe","nested":{"id":"../opaque"},"value":[1,true,null]}`,
		`{"ID":"../opaque"}`,
		`{"value":"no id; leave required-field validation to the service"}`,
		`{"id":"item\"quoted"}`, `{"id":"safe","other":"` + string([]byte{0xff}) + `"}`,
	} {
		require.NoError(t, client.validateItemBodyID([]byte(body)), "body: %s", body)
	}
	for _, body := range []string{
		``, `not JSON`, `null`, `[]`, `"item"`, `{"id":"safe"} {}`,
		`{"id":""}`, `{"id":null}`, `{"id":42}`, `{"id":{}}`,
		`{"id":"../item"}`, `{"\u0069d":"item\nname"}`,
		`{"id":"safe","id":"other"}`, `{"id":"safe","\u0069d":"other"}`,
		`{"id":"` + string([]byte{0xff}) + `"}`,
	} {
		require.Error(t, client.validateItemBodyID([]byte(body)), "body: %s", body)
		disabled := &Client{}
		require.NoError(t, disabled.validateItemBodyID([]byte(body)))
	}
}

func TestResourceNameIDsJSONEscaped(t *testing.T) {
	for _, id := range []string{`allowed","id":"victim`, `item"name`, `item\name`, "item\n\t\r", "item\x00", "\u4e16\u754c"} {
		for _, properties := range []any{DatabaseProperties{ID: id}, ContainerProperties{ID: id}} {
			body, err := json.Marshal(properties)
			require.NoError(t, err)
			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(body, &fields))
			var actualID string
			require.NoError(t, json.Unmarshal(fields["id"], &actualID))
			require.Equal(t, id, actualID)
			require.Equal(t, 1, bytes.Count(body, []byte(`"id":`)))
		}
	}
}

func TestResourceNameValidationConstructors(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		client := &Client{enableResourceNameValidation: enabled}
		database, err := client.NewDatabase("database")
		require.NoError(t, err)
		for _, name := range []string{".", "..", "other/name", `other\name`, "name?", "name#", "name\t", "name\n", "name\r", string([]byte{0xff})} {
			constructors := map[string]func() error{
				"database":           func() error { _, err := client.NewDatabase(name); return err },
				"container database": func() error { _, err := client.NewContainer(name, "container"); return err },
				"container":          func() error { _, err := client.NewContainer("database", name); return err },
				"database container": func() error { _, err := database.NewContainer(name); return err },
			}
			for operation, construct := range constructors {
				err := construct()
				if enabled {
					require.Error(t, err, "%s: %q", operation, name)
				} else {
					require.NoError(t, err, "%s: %q", operation, name)
				}
			}
		}
		_, err = client.NewDatabase("")
		require.Error(t, err)
		_, err = client.NewContainer("", "container")
		require.Error(t, err)
		_, err = client.NewContainer("database", "")
		require.Error(t, err)
		_, err = database.NewContainer("")
		require.Error(t, err)
	}
}

func TestResourceNameValidationInvalidUTF8BeforeRequests(t *testing.T) {
	client := &Client{enableResourceNameValidation: true}
	database, err := client.NewDatabase("database")
	require.NoError(t, err)
	container, err := database.NewContainer("container")
	require.NoError(t, err)
	pk := NewPartitionKeyString("pk")
	for _, id := range []string{string([]byte{0xff}), "id" + string([]byte{0xe2, 0x82})} {
		body := []byte(`{"id":"` + id + `"}`)
		operations := map[string]func() error{
			"database constructor":           func() error { _, err := client.NewDatabase(id); return err },
			"container database constructor": func() error { _, err := client.NewContainer(id, "container"); return err },
			"container constructor":          func() error { _, err := client.NewContainer("database", id); return err },
			"create database": func() error {
				_, err := client.CreateDatabase(context.Background(), DatabaseProperties{ID: id}, nil)
				return err
			},
			"create container": func() error {
				_, err := database.CreateContainer(context.Background(), ContainerProperties{ID: id}, nil)
				return err
			},
			"replace container": func() error {
				_, err := container.Replace(context.Background(), ContainerProperties{ID: id}, nil)
				return err
			},
			"create item": func() error { _, err := container.CreateItem(context.Background(), pk, body, nil); return err },
			"upsert item": func() error { _, err := container.UpsertItem(context.Background(), pk, body, nil); return err },
			"replace item": func() error {
				_, err := container.ReplaceItem(context.Background(), pk, id, []byte(`{}`), nil)
				return err
			},
			"replace item body": func() error {
				_, err := container.ReplaceItem(context.Background(), pk, "item", body, nil)
				return err
			},
			"read item":   func() error { _, err := container.ReadItem(context.Background(), pk, id, nil); return err },
			"delete item": func() error { _, err := container.DeleteItem(context.Background(), pk, id, nil); return err },
			"patch item": func() error {
				_, err := container.PatchItem(context.Background(), pk, id, PatchOperations{}, nil)
				return err
			},
			"read many": func() error {
				_, err := container.ReadManyItems(context.Background(), []ItemIdentity{{ID: id, PartitionKey: pk}}, nil)
				return err
			},
		}
		for operation, execute := range operations {
			t.Run(operation, func(t *testing.T) {
				require.ErrorContains(t, execute(), "valid UTF-8")
			})
		}
	}
}

func TestResourceNameValidationClientOptions(t *testing.T) {
	key, err := NewKeyCredential("dGVzdA==")
	require.NoError(t, err)
	constructors := map[string]func(*ClientOptions) (*Client, error){
		"key": func(options *ClientOptions) (*Client, error) {
			return NewClientWithKey("https://localhost", key, options)
		},
		"token": func(options *ClientOptions) (*Client, error) {
			return NewClient("https://localhost", mockCredential{}, options)
		},
		"connection string": func(options *ClientOptions) (*Client, error) {
			return NewClientFromConnectionString("AccountEndpoint=https://localhost;AccountKey=dGVzdA==;", options)
		},
	}
	for name, construct := range constructors {
		t.Run(name, func(t *testing.T) {
			for _, options := range []*ClientOptions{nil, {}, {EnableResourceNameValidation: true}} {
				client, err := construct(options)
				require.NoError(t, err)
				t.Cleanup(client.Close)
				enabled := options != nil && options.EnableResourceNameValidation
				require.Equal(t, enabled, client.enableResourceNameValidation)
				if options != nil {
					options.EnableResourceNameValidation = !enabled
					require.Equal(t, enabled, client.enableResourceNameValidation, "options must be captured at construction")
				}
			}
		})
	}
}

func TestEncodingScenarioResourceNameCompatibility(t *testing.T) {
	emulator := &emulatorTests{host: "https://localhost", key: "dGVzdA=="}
	tracer := newSpanValidator(t, &spanMatcher{})
	validatingClient := emulator.getClientWithOptions(t, tracer, ClientOptions{EnableResourceNameValidation: true})
	defer validatingClient.Close()
	defaultClient := emulator.getClient(t, tracer)
	defer defaultClient.Close()

	statuses := []int{http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent}
	transport := &countingTransport{}
	transport.respFunc = func() (int, []byte) {
		return statuses[transport.count.Load()-1], []byte(`{}`)
	}
	internal, err := azcore.NewClient("azcosmostest", "v1.0.0", runtime.PipelineOptions{}, &policy.ClientOptions{Transport: transport})
	require.NoError(t, err)
	validatingClient.internal = internal
	defaultClient.internal = internal

	validatingContainer, err := validatingClient.NewContainer("database", "container")
	require.NoError(t, err)
	id := "Disallowed?Chars"
	_, err = validatingContainer.CreateItem(context.Background(), NewPartitionKeyString(id), []byte(`{"id":"Disallowed?Chars"}`), nil)
	require.ErrorContains(t, err, "invalid item ID")
	require.Zero(t, transport.count.Load())

	defaultContainer, err := defaultClient.NewContainer("database", "container")
	require.NoError(t, err)
	verifyEncodingScenario(t, defaultContainer, "default", id, http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusNoContent)
	require.EqualValues(t, 4, transport.count.Load())
}

func TestResourceNameValidationBeforeRequests(t *testing.T) {
	ctx := context.Background()
	pk := NewPartitionKeyString("pk")
	reachedPipeline := errors.New("request reached pipeline")
	for _, enabled := range []bool{false, true} {
		requests := 0
		internal, err := azcore.NewClient("azcosmostest", "v1.0.0", runtime.PipelineOptions{
			PerCall: []policy.Policy{policyFunc(func(*policy.Request) (*http.Response, error) {
				requests++
				return nil, reachedPipeline
			})},
		}, nil)
		require.NoError(t, err)
		client := &Client{
			endpoint: "https://localhost", internal: internal,
			endpointUrl: &url.URL{Scheme: "https", Host: "localhost"},
			gem:         &globalEndpointManager{}, enableResourceNameValidation: enabled,
		}
		database, err := client.NewDatabase("database")
		require.NoError(t, err)
		container, err := database.NewContainer("container")
		require.NoError(t, err)
		for _, id := range []string{".", "..", "other/item", `other\item`, "item?", "item#", "item\t", "item\n", "item\r"} {
			body, err := json.Marshal(map[string]string{"id": id})
			require.NoError(t, err)
			operations := map[string]func() error{
				"create database": func() error {
					_, err := client.CreateDatabase(ctx, DatabaseProperties{ID: id}, nil)
					return err
				},
				"create container": func() error {
					_, err := database.CreateContainer(ctx, ContainerProperties{ID: id}, nil)
					return err
				},
				"replace container": func() error {
					_, err := container.Replace(ctx, ContainerProperties{ID: id}, nil)
					return err
				},
				"create item":       func() error { _, err := container.CreateItem(ctx, pk, body, nil); return err },
				"upsert item":       func() error { _, err := container.UpsertItem(ctx, pk, body, nil); return err },
				"replace item":      func() error { _, err := container.ReplaceItem(ctx, pk, id, []byte(`{}`), nil); return err },
				"replace item body": func() error { _, err := container.ReplaceItem(ctx, pk, "item", body, nil); return err },
				"read item":         func() error { _, err := container.ReadItem(ctx, pk, id, nil); return err },
				"delete item":       func() error { _, err := container.DeleteItem(ctx, pk, id, nil); return err },
				"patch item":        func() error { _, err := container.PatchItem(ctx, pk, id, PatchOperations{}, nil); return err },
				"read many": func() error {
					_, err := container.ReadManyItems(ctx, []ItemIdentity{{ID: "safe", PartitionKey: pk}, {ID: id, PartitionKey: pk}}, nil)
					return err
				},
			}
			for operation, execute := range operations {
				t.Run(operation+"/"+id+map[bool]string{false: "/default", true: "/enabled"}[enabled], func(t *testing.T) {
					before := requests
					err := execute()
					if enabled {
						require.ErrorContains(t, err, "invalid ")
						require.Equal(t, before, requests, "validation must precede all network I/O")
					} else {
						require.ErrorIs(t, err, reachedPipeline)
						require.Greater(t, requests, before)
					}
				})
			}
		}
	}
}
