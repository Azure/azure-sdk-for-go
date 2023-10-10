//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/internal/generated"
)

const timeFormat = time.RFC3339Nano

// Client is the struct for interacting with an Azure App Configuration instance.
type Client struct {
	appConfigClient *generated.AzureAppConfigurationClient
	syncTokenPolicy *syncTokenPolicy
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient returns a pointer to a Client object affinitized to an endpoint.
func NewClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return newClient(endpoint, runtime.NewBearerTokenPolicy(cred, []string{
		fmt.Sprintf("%s://%s/.default", u.Scheme, u.Host),
	}, nil), options)
}

// NewClientFromConnectionString parses the connection string and returns a pointer to a Client object.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	endpoint, credential, secret, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return newClient(endpoint, newHmacAuthenticationPolicy(credential, secret), options)
}

func newClient(endpoint string, authPolicy policy.Policy, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	syncTokenPolicy := newSyncTokenPolicy()

	client, err := azcore.NewClient(moduleName+".Client", moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy, syncTokenPolicy},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		appConfigClient: generated.NewAzureAppConfigurationClient(endpoint, client),
		syncTokenPolicy: syncTokenPolicy,
	}, nil
}

// UpdateSyncToken sets an external synchronization token to ensure service requests receive up-to-date values.
func (c *Client) UpdateSyncToken(token string) {
	c.syncTokenPolicy.addToken(token)
}

// AddSetting creates a configuration setting only if the setting does not already exist in the configuration store.
//   - ctx controls the lifetime of the HTTP operation
//   - key is the name of the setting to create
//   - value is the value for the setting. pass nil if the setting doesn't have a value
//   - options contains the optional values. can be nil
func (c *Client) AddSetting(ctx context.Context, key string, value *string, options *AddSettingOptions) (AddSettingResponse, error) {
	if options == nil {
		options = &AddSettingOptions{}
	}

	setting := Setting{Key: &key, Value: value, Label: options.Label, ContentType: options.ContentType}

	etagAny := azcore.ETagAny
	kv, opts := setting.toGeneratedPutOptions(nil, &etagAny)
	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.Key, kv, &opts)
	if err != nil {
		return AddSettingResponse{}, err
	}

	return AddSettingResponse{
		Setting:   settingFromGenerated(resp.KeyValue),
		SyncToken: resp.SyncToken,
	}, nil
}

// DeleteSetting deletes a configuration setting from the configuration store.
func (c *Client) DeleteSetting(ctx context.Context, key string, options *DeleteSettingOptions) (DeleteSettingResponse, error) {
	if options == nil {
		options = &DeleteSettingOptions{}
	}

	setting := Setting{Key: &key, Label: options.Label}

	resp, err := c.appConfigClient.DeleteKeyValue(ctx, *setting.Key, setting.toGeneratedDeleteOptions(options.OnlyIfUnchanged))
	if err != nil {
		return DeleteSettingResponse{}, err
	}

	return DeleteSettingResponse{
		Setting:   settingFromGenerated(resp.KeyValue),
		SyncToken: resp.SyncToken,
	}, nil
}

// GetSetting retrieves an existing configuration setting from the configuration store.
func (c *Client) GetSetting(ctx context.Context, key string, options *GetSettingOptions) (GetSettingResponse, error) {
	if options == nil {
		options = &GetSettingOptions{}
	}

	setting := Setting{Key: &key, Label: options.Label}

	resp, err := c.appConfigClient.GetKeyValue(ctx, *setting.Key, setting.toGeneratedGetOptions(options.OnlyIfChanged, options.AcceptDateTime))
	if err != nil {
		return GetSettingResponse{}, err
	}

	var lastModified *time.Time
	if resp.LastModified != nil {
		tt, err := time.Parse(http.TimeFormat, *resp.LastModified)
		if err != nil {
			return GetSettingResponse{}, err
		}
		lastModified = &tt
	}

	return GetSettingResponse{
		Setting:      settingFromGenerated(resp.KeyValue),
		SyncToken:    resp.SyncToken,
		LastModified: lastModified,
	}, nil
}

// SetReadOnly sets an existing configuration setting to read only or read write state in the configuration store.
func (c *Client) SetReadOnly(ctx context.Context, key string, isReadOnly bool, options *SetReadOnlyOptions) (SetReadOnlyResponse, error) {
	if options == nil {
		options = &SetReadOnlyOptions{}
	}

	setting := Setting{Key: &key, Label: options.Label}

	var err error
	if isReadOnly {
		var resp generated.AzureAppConfigurationClientPutLockResponse
		resp, err = c.appConfigClient.PutLock(ctx, *setting.Key, setting.toGeneratedPutLockOptions(options.OnlyIfUnchanged))
		if err == nil {
			return SetReadOnlyResponse{
				Setting:   settingFromGenerated(resp.KeyValue),
				SyncToken: resp.SyncToken,
			}, nil
		}
	} else {
		var resp generated.AzureAppConfigurationClientDeleteLockResponse
		resp, err = c.appConfigClient.DeleteLock(ctx, *setting.Key, setting.toGeneratedDeleteLockOptions(options.OnlyIfUnchanged))
		if err == nil {
			return SetReadOnlyResponse{
				Setting:   settingFromGenerated(resp.KeyValue),
				SyncToken: resp.SyncToken,
			}, nil
		}
	}

	return SetReadOnlyResponse{}, err
}

// SetSetting creates a configuration setting if it doesn't exist or overwrites the existing setting in the configuration store.
//   - ctx controls the lifetime of the HTTP operation
//   - key is the name of the setting to create
//   - value is the value for the setting. pass nil if the setting doesn't have a value
//   - options contains the optional values. can be nil
func (c *Client) SetSetting(ctx context.Context, key string, value *string, options *SetSettingOptions) (SetSettingResponse, error) {
	if options == nil {
		options = &SetSettingOptions{}
	}

	setting := Setting{Key: &key, Value: value, Label: options.Label, ContentType: options.ContentType}

	kv, opts := setting.toGeneratedPutOptions(options.OnlyIfUnchanged, nil)
	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.Key, kv, &opts)
	if err != nil {
		return SetSettingResponse{}, err
	}

	return SetSettingResponse{
		Setting:   settingFromGenerated(resp.KeyValue),
		SyncToken: resp.SyncToken,
	}, nil
}

// NewListRevisionsPager creates a pager that retrieves the revisions of one or more
// configuration setting entities that match the specified setting selector.
func (c *Client) NewListRevisionsPager(selector SettingSelector, options *ListRevisionsOptions) *runtime.Pager[ListRevisionsPageResponse] {
	pagerInternal := c.appConfigClient.NewGetRevisionsPager(selector.toGeneratedGetRevisions())
	return runtime.NewPager(runtime.PagingHandler[ListRevisionsPageResponse]{
		More: func(ListRevisionsPageResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListRevisionsPageResponse) (ListRevisionsPageResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListRevisionsPageResponse{}, err
			}
			var css []Setting
			for _, cs := range page.Items {
				if cs != nil {
					css = append(css, settingFromGenerated(*cs))
				}
			}

			return ListRevisionsPageResponse{
				Settings:  css,
				SyncToken: page.SyncToken,
			}, nil
		},
	})
}

// NewListSettingsPager creates a pager that retrieves setting entities that match the specified setting selector.
func (c *Client) NewListSettingsPager(selector SettingSelector, options *ListSettingsOptions) *runtime.Pager[ListSettingsPageResponse] {
	pagerInternal := c.appConfigClient.NewGetKeyValuesPager(selector.toGeneratedGetKeyValues())
	return runtime.NewPager(runtime.PagingHandler[ListSettingsPageResponse]{
		More: func(ListSettingsPageResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListSettingsPageResponse) (ListSettingsPageResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListSettingsPageResponse{}, err
			}
			var css []Setting
			for _, cs := range page.Items {
				if cs != nil {
					css = append(css, settingFromGenerated(*cs))
				}
			}

			return ListSettingsPageResponse{
				Settings:  css,
				SyncToken: page.SyncToken,
			}, nil
		},
	})
}
