// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pagingSyncToken azappconfig.SyncToken = "test=MQ==;sn=1"

func newPagingClient(t *testing.T) (*azappconfig.Client, *mock.Server) {
	t.Helper()
	srv, close := mock.NewServer()
	t.Cleanup(close)
	client, err := azappconfig.NewClientFromConnectionString(
		fmt.Sprintf("Endpoint=%s;Id=fake;Secret=ZmFrZQ==", srv.URL()),
		&azappconfig.ClientOptions{ClientOptions: azcore.ClientOptions{
			Transport: srv,
			Retry:     policy.RetryOptions{MaxRetries: -1},
		}},
	)
	require.NoError(t, err)
	return client, srv
}

func assertPages[T any](t *testing.T, pager *runtime.Pager[T], expected []T) {
	t.Helper()
	for i, want := range expected {
		require.True(t, pager.More(), "page %d should be available", i)
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err, "fetching page %d", i)
		assert.Equal(t, want, page, "page %d", i)
	}
	assert.False(t, pager.More(), "the final page should end pagination")
}

func TestPagersNextLinks(t *testing.T) {
	settings := []azappconfig.Setting{
		{Key: to.Ptr("first"), Value: to.Ptr("value")},
		{Key: to.Ptr("second"), Value: to.Ptr("value")},
	}
	snapshots := []azappconfig.Snapshot{
		{Name: to.Ptr("first"), Filters: []azappconfig.SettingFilter{{KeyFilter: to.Ptr("key*")}}},
		{Name: to.Ptr("second"), Filters: []azappconfig.SettingFilter{{KeyFilter: to.Ptr("key*")}}},
	}
	etags := []azcore.ETag{`"first"`, `"second"`}
	tests := []struct {
		name       string
		path       string
		method     string
		itemFormat string
		run        func(*testing.T, *azappconfig.Client)
	}{
		{
			name: "Revisions", path: "/revisions", method: http.MethodGet,
			itemFormat: `{"key":%q,"value":"value"}`,
			run: func(t *testing.T, client *azappconfig.Client) {
				assertPages(t, client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil), []azappconfig.ListRevisionsPageResponse{
					{Settings: settings[:1], SyncToken: pagingSyncToken},
					{Settings: settings[1:], SyncToken: pagingSyncToken},
				})
			},
		},
		{
			name: "Settings", path: "/kv", method: http.MethodGet,
			itemFormat: `{"key":%q,"value":"value"}`,
			run: func(t *testing.T, client *azappconfig.Client) {
				assertPages(t, client.NewListSettingsPager(azappconfig.SettingSelector{}, nil), []azappconfig.ListSettingsPageResponse{
					{Settings: settings[:1], ETag: &etags[0], SyncToken: pagingSyncToken},
					{Settings: settings[1:], ETag: &etags[1], SyncToken: pagingSyncToken},
				})
			},
		},
		{
			name: "CheckSettings", path: "/kv", method: http.MethodHead,
			run: func(t *testing.T, client *azappconfig.Client) {
				assertPages(t, client.NewCheckSettingsPager(azappconfig.SettingSelector{}, nil), []azappconfig.CheckSettingsPageResponse{
					{ETag: &etags[0], SyncToken: pagingSyncToken},
					{ETag: &etags[1], SyncToken: pagingSyncToken},
				})
			},
		},
		{
			name: "Snapshots", path: "/snapshots", method: http.MethodGet,
			itemFormat: `{"name":%q,"filters":[{"key":"key*"}]}`,
			run: func(t *testing.T, client *azappconfig.Client) {
				assertPages(t, client.NewListSnapshotsPager(nil), []azappconfig.ListSnapshotsResponse{
					{Snapshots: snapshots[:1], SyncToken: pagingSyncToken},
					{Snapshots: snapshots[1:], SyncToken: pagingSyncToken},
				})
			},
		},
		{
			name: "SnapshotSettings", path: "/kv", method: http.MethodGet,
			itemFormat: `{"key":%q,"value":"value"}`,
			run: func(t *testing.T, client *azappconfig.Client) {
				assertPages(t, client.NewListSettingsForSnapshotPager("snapshot", nil), []azappconfig.ListSettingsForSnapshotResponse{
					{Settings: settings[:1], SyncToken: pagingSyncToken},
					{Settings: settings[1:], SyncToken: pagingSyncToken},
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, linkType := range []string{"Relative", "Absolute"} {
				t.Run(linkType, func(t *testing.T) {
					client, srv := newPagingClient(t)
					nextLink := tt.path + "?after=cursor%2B%2F%3D&api-version=next-version"
					if linkType == "Absolute" {
						nextLink = srv.URL() + nextLink
					}
					for i, name := range []string{"first", "second"} {
						srv.AppendResponse(mock.WithPredicate(func(req *http.Request) bool {
							assert.Equal(t, tt.method, req.Method)
							assert.Equal(t, tt.path, req.URL.Path)
							if i == 1 {
								assert.Equal(t, "cursor+/=", req.URL.Query().Get("after"))
								assert.Equal(t, "next-version", req.URL.Query().Get("api-version"))
							} else if tt.name == "SnapshotSettings" {
								assert.Equal(t, "snapshot", req.URL.Query().Get("snapshot"))
							}
							return false
						}))
						response := []mock.ResponseOption{
							mock.WithHeader("Sync-Token", string(pagingSyncToken)),
							mock.WithHeader("ETag", string(etags[i])),
						}
						if tt.method == http.MethodHead {
							if i == 0 {
								response = append(response, mock.WithHeader("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextLink)))
							}
						} else {
							body := fmt.Sprintf(`{"items":[%s]}`, fmt.Sprintf(tt.itemFormat, name))
							if i == 0 {
								body = fmt.Sprintf(`{"items":[%s],"@nextLink":%q}`, fmt.Sprintf(tt.itemFormat, name), nextLink)
							}
							response = append(response, mock.WithHeader("Content-Type", "application/json"), mock.WithBody([]byte(body)))
						}
						srv.AppendResponse(response...)
					}
					tt.run(t, client)
					assert.Equal(t, 2, srv.Requests())
				})
			}
		})
	}
}

func TestSettingsPagersMatchConditions(t *testing.T) {
	conditions := []azcore.MatchConditions{
		{IfNoneMatch: to.Ptr(azcore.ETag(`"first"`))},
		{IfMatch: to.Ptr(azcore.ETag(`"second"`))},
	}
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		t.Run(method, func(t *testing.T) {
			client, srv := newPagingClient(t)
			for i := range 3 {
				srv.AppendResponse(mock.WithPredicate(func(req *http.Request) bool {
					assert.Equal(t, method, req.Method)
					assert.Equal(t, "/kv", req.URL.Path)
					switch i {
					case 0:
						assert.Equal(t, `"first"`, req.Header.Get("If-None-Match"))
						assert.Empty(t, req.Header.Get("If-Match"))
					case 1:
						assert.Equal(t, `"second"`, req.Header.Get("If-Match"))
						assert.Empty(t, req.Header.Get("If-None-Match"))
					case 2:
						assert.Empty(t, req.Header.Get("If-Match"))
						assert.Empty(t, req.Header.Get("If-None-Match"))
					}
					return false
				}))
				response := []mock.ResponseOption{
					mock.WithHeader("ETag", fmt.Sprintf("\"page-%d\"", i)),
					mock.WithHeader("Sync-Token", string(pagingSyncToken)),
				}
				if i == 0 {
					response = append(response, mock.WithStatusCode(http.StatusNotModified))
				} else if method == http.MethodGet {
					response = append(response, mock.WithHeader("Content-Type", "application/json"), mock.WithBody([]byte(`{"items":[]}`)))
				}
				if i < 2 {
					response = append(response, mock.WithHeader("Link", fmt.Sprintf("</kv?after=%d>; rel=\"next\"", i+1)))
				}
				srv.AppendResponse(response...)
			}
			if method == http.MethodGet {
				expected := make([]azappconfig.ListSettingsPageResponse, 3)
				for i := range expected {
					expected[i] = azappconfig.ListSettingsPageResponse{
						Settings:  []azappconfig.Setting{},
						ETag:      to.Ptr(azcore.ETag(fmt.Sprintf("\"page-%d\"", i))),
						SyncToken: pagingSyncToken,
					}
				}
				assertPages(t, client.NewListSettingsPager(azappconfig.SettingSelector{}, &azappconfig.ListSettingsOptions{MatchConditions: conditions}), expected)
			} else {
				expected := make([]azappconfig.CheckSettingsPageResponse, 3)
				for i := range expected {
					expected[i] = azappconfig.CheckSettingsPageResponse{
						ETag:      to.Ptr(azcore.ETag(fmt.Sprintf("\"page-%d\"", i))),
						SyncToken: pagingSyncToken,
					}
				}
				assertPages(t, client.NewCheckSettingsPager(azappconfig.SettingSelector{}, &azappconfig.CheckSettingsOptions{MatchConditions: conditions}), expected)
			}
			assert.Equal(t, 3, srv.Requests())
		})
	}
}

func TestRevisionsPagerErrors(t *testing.T) {
	t.Run("Transport", func(t *testing.T) {
		client, srv := newPagingClient(t)
		expected := errors.New("transport failed")
		srv.AppendError(expected)
		pager := client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil)
		page, err := pager.NextPage(context.Background())
		require.ErrorIs(t, err, expected)
		assert.Empty(t, page)
		assert.Equal(t, 1, srv.Requests())
	})
	t.Run("Service", func(t *testing.T) {
		client, srv := newPagingClient(t)
		srv.AppendResponse(mock.WithStatusCode(http.StatusForbidden))
		pager := client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil)
		page, err := pager.NextPage(context.Background())
		var responseError *azcore.ResponseError
		require.ErrorAs(t, err, &responseError)
		assert.Equal(t, http.StatusForbidden, responseError.StatusCode)
		assert.Empty(t, page)
	})
	t.Run("InvalidJSON", func(t *testing.T) {
		client, srv := newPagingClient(t)
		srv.AppendResponse(mock.WithBody([]byte(`{"items":`)))
		pager := client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil)
		page, err := pager.NextPage(context.Background())
		require.Error(t, err)
		assert.Empty(t, page)
	})
	t.Run("InvalidNextLink", func(t *testing.T) {
		client, srv := newPagingClient(t)
		srv.AppendResponse(
			mock.WithHeader("Sync-Token", string(pagingSyncToken)),
			mock.WithBody([]byte(`{"items":[],"@nextLink":"http://[invalid"}`)),
		)
		pager := client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil)
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.True(t, pager.More())
		page, err := pager.NextPage(context.Background())
		require.Error(t, err)
		assert.Empty(t, page)
		assert.Equal(t, 1, srv.Requests())
	})
	t.Run("NextPageTransport", func(t *testing.T) {
		client, srv := newPagingClient(t)
		srv.AppendResponse(
			mock.WithHeader("Sync-Token", string(pagingSyncToken)),
			mock.WithBody([]byte(`{"items":[],"@nextLink":"/revisions?after=cursor"}`)),
		)
		expected := errors.New("next page transport failed")
		srv.AppendError(expected)
		pager := client.NewListRevisionsPager(azappconfig.SettingSelector{}, nil)
		_, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.True(t, pager.More())
		page, err := pager.NextPage(context.Background())
		require.ErrorIs(t, err, expected)
		assert.Empty(t, page)
		assert.Equal(t, 2, srv.Requests())
	})
}
