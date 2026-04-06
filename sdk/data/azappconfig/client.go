// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/audience"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/synctoken"
)

const timeFormat = time.RFC3339Nano

// Client is the struct for interacting with an Azure App Configuration instance.
type Client struct {
	appConfigClient *generated.AzureAppConfigurationClient
	cache           *synctoken.Cache
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

	audience := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if options != nil && !reflect.ValueOf(options.Cloud).IsZero() {
		if cfg, ok := options.Cloud.Services[ServiceName]; ok && cfg.Audience != "" {
			audience = cfg.Audience
		}
	}

	return newClient(endpoint, runtime.NewBearerTokenPolicy(cred, []string{audience + "/.default"}, nil), options)
}

// NewClientFromConnectionString parses the connection string and returns a pointer to a Client object.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	endpoint, credential, secret, err := auth.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return newClient(endpoint, auth.NewHMACPolicy(credential, secret), options)
}

func newClient(endpoint string, authPolicy policy.Policy, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	audienceConfigured := false
	if !reflect.ValueOf(options.Cloud).IsZero() {
		if cfg, ok := options.Cloud.Services[ServiceName]; ok && cfg.Audience != "" {
			audienceConfigured = true
		}
	}

	cache := synctoken.NewCache()
	client, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy, synctoken.NewPolicy(cache), audience.NewAudienceErrorHandlingPolicy(audienceConfigured)},
		Tracing: runtime.TracingOptions{
			Namespace: "Microsoft.AppConfig",
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		appConfigClient: generated.NewAzureAppConfigurationClient(endpoint, client),
		cache:           cache,
	}, nil
}

// SetSyncToken is used to set a sync token from an external source.
// SyncTokens are required to be in the format "<id>=<value>;sn=<sn>".
// Multiple SyncTokens must be comma delimited.
func (c *Client) SetSyncToken(syncToken SyncToken) error {
	return c.cache.Set(syncToken)
}

// AddSetting creates a configuration setting only if the setting does not already exist in the configuration store.
//   - ctx controls the lifetime of the HTTP operation
//   - key is the name of the setting to create
//   - value is the value for the setting. pass nil if the setting doesn't have a value
//   - options contains the optional values. can be nil
func (c *Client) AddSetting(ctx context.Context, key string, value *string, options *AddSettingOptions) (AddSettingResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.AddSetting", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &AddSettingOptions{}
	}

	setting := Setting{Key: &key, Value: value, Label: options.Label, ContentType: options.ContentType, Tags: options.Tags}

	etagAny := azcore.ETagAny
	kv, opts := setting.toGeneratedPutOptions(nil, &etagAny)
	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.Key, kv, &opts)
	if err != nil {
		return AddSettingResponse{}, err
	}

	return AddSettingResponse{
		Setting:   settingFromGenerated(resp.KeyValue),
		SyncToken: SyncToken(*resp.SyncToken),
	}, nil
}

// DeleteSetting deletes a configuration setting from the configuration store.
func (c *Client) DeleteSetting(ctx context.Context, key string, options *DeleteSettingOptions) (DeleteSettingResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.DeleteSetting", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

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
		SyncToken: SyncToken(*resp.SyncToken),
	}, nil
}

// GetSetting retrieves an existing configuration setting from the configuration store.
func (c *Client) GetSetting(ctx context.Context, key string, options *GetSettingOptions) (GetSettingResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetSetting", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetSettingOptions{}
	}

	setting := Setting{Key: &key, Label: options.Label}

	resp, err := c.appConfigClient.GetKeyValue(ctx, *setting.Key, setting.toGeneratedGetOptions(options.OnlyIfChanged, options.AcceptDateTime))
	if err != nil {
		return GetSettingResponse{}, err
	}

	return GetSettingResponse{
		Setting:      settingFromGenerated(resp.KeyValue),
		SyncToken:    SyncToken(*resp.SyncToken),
		LastModified: resp.LastModified,
	}, nil
}

// SetReadOnly sets an existing configuration setting to read only or read write state in the configuration store.
func (c *Client) SetReadOnly(ctx context.Context, key string, isReadOnly bool, options *SetReadOnlyOptions) (SetReadOnlyResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.SetReadOnly", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &SetReadOnlyOptions{}
	}

	setting := Setting{Key: &key, Label: options.Label}

	if isReadOnly {
		var resp generated.AzureAppConfigurationClientPutLockResponse
		resp, err = c.appConfigClient.PutLock(ctx, *setting.Key, setting.toGeneratedPutLockOptions(options.OnlyIfUnchanged))
		if err == nil {
			return SetReadOnlyResponse{
				Setting:   settingFromGenerated(resp.KeyValue),
				SyncToken: SyncToken(*resp.SyncToken),
			}, nil
		}
	} else {
		var resp generated.AzureAppConfigurationClientDeleteLockResponse
		resp, err = c.appConfigClient.DeleteLock(ctx, *setting.Key, setting.toGeneratedDeleteLockOptions(options.OnlyIfUnchanged))
		if err == nil {
			return SetReadOnlyResponse{
				Setting:   settingFromGenerated(resp.KeyValue),
				SyncToken: SyncToken(*resp.SyncToken),
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
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.SetSetting", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &SetSettingOptions{}
	}

	setting := Setting{Key: &key, Value: value, Label: options.Label, ContentType: options.ContentType, Tags: options.Tags}

	kv, opts := setting.toGeneratedPutOptions(options.OnlyIfUnchanged, nil)
	resp, err := c.appConfigClient.PutKeyValue(ctx, *setting.Key, kv, &opts)
	if err != nil {
		return SetSettingResponse{}, err
	}

	return SetSettingResponse{
		Setting:   settingFromGenerated(resp.KeyValue),
		SyncToken: SyncToken(*resp.SyncToken),
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

			css := make([]Setting, len(page.Items))
			for i := range page.Items {
				css[i] = settingFromGenerated(page.Items[i])
			}

			return ListRevisionsPageResponse{
				Settings:  css,
				SyncToken: SyncToken(*page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

// NewListSettingsPager creates a pager that retrieves setting entities that match the specified setting selector.
func (c *Client) NewListSettingsPager(selector SettingSelector, options *ListSettingsOptions) *runtime.Pager[ListSettingsPageResponse] {
	if options == nil {
		options = &ListSettingsOptions{}
	}
	pagerInternal := c.appConfigClient.NewGetKeyValuesPagerWithMatchConditions(options.MatchConditions, selector.toGeneratedGetKeyValues())
	return runtime.NewPager(runtime.PagingHandler[ListSettingsPageResponse]{
		More: func(ListSettingsPageResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListSettingsPageResponse) (ListSettingsPageResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListSettingsPageResponse{}, err
			}
			css := make([]Setting, len(page.Items))
			for i := range page.Items {
				css[i] = settingFromGenerated(page.Items[i])
			}

			return ListSettingsPageResponse{
				Settings:  css,
				ETag:      (*azcore.ETag)(page.ETag),
				SyncToken: SyncToken(*page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

// NewListSnapshotsPager - Gets a list of key-value snapshots.
//
//   - options - NewListSnapshotsPagerOptions contains the optional parameters to retrieve a snapshot
//     method.
func (c *Client) NewListSnapshotsPager(options *ListSnapshotsOptions) *runtime.Pager[ListSnapshotsResponse] {
	opts := (*generated.AzureAppConfigurationClientGetSnapshotsOptions)(options)
	ssRespPager := c.appConfigClient.NewGetSnapshotsPager(opts)

	return runtime.NewPager(runtime.PagingHandler[ListSnapshotsResponse]{
		More: func(ListSnapshotsResponse) bool {
			return ssRespPager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListSnapshotsResponse) (ListSnapshotsResponse, error) {
			page, err := ssRespPager.NextPage(ctx)
			if err != nil {
				return ListSnapshotsResponse{}, err
			}

			snapshots := make([]Snapshot, len(page.Items))

			for i := range page.Items {
				snapshot := page.Items[i]

				convertedFilters := make([]SettingFilter, len(snapshot.Filters))
				for j := range snapshot.Filters {
					convertedFilters[j] = SettingFilter{
						KeyFilter:   snapshot.Filters[j].Key,
						LabelFilter: snapshot.Filters[j].Label,
					}
				}

				snapshots[i] = Snapshot{
					Filters:         convertedFilters,
					CompositionType: snapshot.CompositionType,
					RetentionPeriod: snapshot.RetentionPeriod,
					Tags:            snapshot.Tags,
					Created:         snapshot.Created,
					ETag:            (*azcore.ETag)(snapshot.Etag),
					Expires:         snapshot.Expires,
					ItemsCount:      snapshot.ItemsCount,
					Name:            snapshot.Name,
					Size:            snapshot.Size,
					Status:          snapshot.Status,
				}
			}

			return ListSnapshotsResponse{
				Snapshots: snapshots,
				SyncToken: SyncToken(*page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

// NewListSettingsForSnapshotPager
//
// - snapshotName - The name of the snapshot to list configuration settings for
// - options - ListSettingsForSnapshotOptions contains the optional parameters to retrieve Snapshot configuration settings
func (c *Client) NewListSettingsForSnapshotPager(snapshotName string, options *ListSettingsForSnapshotOptions) *runtime.Pager[ListSettingsForSnapshotResponse] {
	if options == nil {
		options = &ListSettingsForSnapshotOptions{}
	}

	ssRespPager := c.appConfigClient.NewGetKeyValuesPager(&generated.AzureAppConfigurationClientGetKeyValuesOptions{
		AcceptDatetime: options.AcceptDatetime,
		After:          options.After,
		IfMatch:        (*string)(options.IfMatch),
		IfNoneMatch:    (*string)(options.IfNoneMatch),
		Select:         options.Select,
		Snapshot:       &snapshotName,
		Key:            &options.Key,
		Label:          &options.Label,
	})

	return runtime.NewPager(runtime.PagingHandler[ListSettingsForSnapshotResponse]{
		More: func(ListSettingsForSnapshotResponse) bool {
			return ssRespPager.More()
		},
		Fetcher: func(ctx context.Context, cur *ListSettingsForSnapshotResponse) (ListSettingsForSnapshotResponse, error) {
			page, err := ssRespPager.NextPage(ctx)
			if err != nil {
				return ListSettingsForSnapshotResponse{}, err
			}

			settings := make([]Setting, len(page.Items))
			for i := range page.Items {
				setting := page.Items[i]

				settings[i] = settingFromGenerated(setting)
			}

			return ListSettingsForSnapshotResponse{
				Settings:  settings,
				SyncToken: SyncToken(*page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

// BeginCreateSnapshot creates a snapshot of the configuration store.
//
// - snapshotName - The name of the snapshot to create.
// - settingFilter - The filters to apply on the key-values.
// - options - BeginCreateSnapshotOptions contains the optional parameters to create a Snapshot
func (c *Client) BeginCreateSnapshot(ctx context.Context, snapshotName string, settingFilter []SettingFilter, options *BeginCreateSnapshotOptions) (*runtime.Poller[CreateSnapshotResponse], error) {
	if options == nil {
		options = &BeginCreateSnapshotOptions{}
	}

	filter := make([]generated.KeyValueFilter, len(settingFilter))
	for i := range settingFilter {
		filter[i] = generated.KeyValueFilter{
			Key:   settingFilter[i].KeyFilter,
			Label: settingFilter[i].LabelFilter,
		}
	}

	// if no filters were specified, add an empty filter to mean "all the things"
	if len(filter) == 0 {
		filter = append(filter, generated.KeyValueFilter{})
	}

	entity := generated.Snapshot{
		CompositionType: options.CompositionType,
		Filters:         filter,
		Name:            &snapshotName,
		RetentionPeriod: options.RetentionPeriod,
		Tags:            options.Tags,
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, c.appConfigClient.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CreateSnapshotResponse]{
			Tracer: c.appConfigClient.Tracer(),
		})
	}

	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.BeginCreateSnapshot", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	resp, err := c.appConfigClient.CreateSnapshot(ctx, snapshotName, entity, nil)
	if err != nil {
		return nil, err
	}
	poller, err := runtime.NewPoller(resp, c.appConfigClient.Pipeline(), &runtime.NewPollerOptions[CreateSnapshotResponse]{
		Tracer: c.appConfigClient.Tracer(),
	})
	return poller, err
}

// GetSnapshot gets a snapshot
//
// - snapshotName - The name of the snapshot to get.
// - options - GetSnapshotOptions contains the optional parameters to get a snapshot
func (c *Client) GetSnapshot(ctx context.Context, snapshotName string, options *GetSnapshotOptions) (GetSnapshotResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetSnapshot", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetSnapshotOptions{}
	}

	getResp, err := c.appConfigClient.GetSnapshot(ctx, snapshotName, &generated.AzureAppConfigurationClientGetSnapshotOptions{
		IfMatch:     (*string)(options.IfMatch),
		IfNoneMatch: (*string)(options.IfNoneMatch),
		Select:      options.Select,
	})

	if err != nil {
		return GetSnapshotResponse{}, err
	}

	convertedFilters := make([]SettingFilter, len(getResp.Filters))
	for i := range getResp.Filters {
		convertedFilters[i] = SettingFilter{
			KeyFilter:   getResp.Filters[i].Key,
			LabelFilter: getResp.Filters[i].Label,
		}
	}

	resp := GetSnapshotResponse{
		Snapshot: Snapshot{
			Filters:         convertedFilters,
			CompositionType: getResp.CompositionType,
			RetentionPeriod: getResp.RetentionPeriod,
			Tags:            getResp.Tags,
			Created:         getResp.Created,
			ETag:            (*azcore.ETag)(getResp.Etag),
			Expires:         getResp.Expires,
			ItemsCount:      getResp.ItemsCount,
			Name:            getResp.Name,
			Size:            getResp.Size,
			Status:          getResp.Status,
		},
		SyncToken: SyncToken(*getResp.SyncToken),
		Link:      getResp.Link,
	}

	return resp, nil
}

// ArchiveSnapshot archives a snapshot
//
// - snapshotName - The name of the snapshot to archive.
// - options - ArchiveSnapshotOptions contains the optional parameters to archive a snapshot
func (c *Client) ArchiveSnapshot(ctx context.Context, snapshotName string, options *ArchiveSnapshotOptions) (ArchiveSnapshotResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.ArchiveSnapshot", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &ArchiveSnapshotOptions{}
	}

	opts := updateSnapshotStatusOptions{
		IfMatch:     options.IfMatch,
		IfNoneMatch: options.IfNoneMatch,
	}
	resp, err := c.updateSnapshotStatus(ctx, snapshotName, generated.SnapshotStatusArchived, &opts)

	if err != nil {
		return ArchiveSnapshotResponse{}, err
	}

	return (ArchiveSnapshotResponse)(resp), nil
}

// RecoverSnapshot recovers a snapshot
//
// - snapshotName - The name of the snapshot to recover.
// - options - RecoverSnapshotOptions contains the optional parameters to recover a snapshot
func (c *Client) RecoverSnapshot(ctx context.Context, snapshotName string, options *RecoverSnapshotOptions) (RecoverSnapshotResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.RecoverSnapshot", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &RecoverSnapshotOptions{}
	}

	opts := updateSnapshotStatusOptions{
		IfMatch:     options.IfMatch,
		IfNoneMatch: options.IfNoneMatch,
	}
	resp, err := c.updateSnapshotStatus(ctx, snapshotName, generated.SnapshotStatusReady, &opts)

	if err != nil {
		return RecoverSnapshotResponse{}, err
	}

	return (RecoverSnapshotResponse)(resp), nil
}

func (c *Client) updateSnapshotStatus(ctx context.Context, snapshotName string, status SnapshotStatus, options *updateSnapshotStatusOptions) (updateSnapshotStatusResponse, error) {
	entity := generated.SnapshotUpdateParameters{
		Status: &status,
	}

	updateResp, err := c.appConfigClient.UpdateSnapshot(ctx, snapshotName, entity, &generated.AzureAppConfigurationClientUpdateSnapshotOptions{
		IfMatch:     (*string)(options.IfMatch),
		IfNoneMatch: (*string)(options.IfNoneMatch),
	})

	if err != nil {
		return updateSnapshotStatusResponse{}, err
	}

	convertedFilters := make([]SettingFilter, len(updateResp.Filters))
	for i := range updateResp.Filters {
		convertedFilters[i] = SettingFilter{
			KeyFilter:   updateResp.Filters[i].Key,
			LabelFilter: updateResp.Filters[i].Label,
		}
	}

	resp := updateSnapshotStatusResponse{
		Snapshot: Snapshot{
			Filters:         convertedFilters,
			CompositionType: updateResp.CompositionType,
			RetentionPeriod: updateResp.RetentionPeriod,
			Tags:            updateResp.Tags,
			Created:         updateResp.Created,
			ETag:            (*azcore.ETag)(updateResp.Etag),
			Expires:         updateResp.Expires,
			ItemsCount:      updateResp.ItemsCount,
			Name:            updateResp.Name,
			Size:            updateResp.Size,
			Status:          updateResp.Status,
		},
		SyncToken: SyncToken(*updateResp.SyncToken),
		Link:      updateResp.Link,
	}

	return resp, nil
}
