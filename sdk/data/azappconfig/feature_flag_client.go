// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/audience"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/querynormalization"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/synctoken"
)

// FeatureFlagClient interacts with the Azure App Configuration feature flag endpoint.
//
// Use FeatureFlagClient to manage feature flags. To manage key-value configuration
// settings, use [Client] instead.
type FeatureFlagClient struct {
	appConfigClient *generated.AzureAppConfigurationClient
	ffClient        *generated.AzureAppConfigurationFeatureFlagClient
	cache           *synctoken.Cache
}

// FeatureFlagClientOptions are the configurable options on a [FeatureFlagClient].
type FeatureFlagClientOptions struct {
	azcore.ClientOptions
}

// NewFeatureFlagClient returns a pointer to a FeatureFlagClient object affinitized to the given endpoint.
func NewFeatureFlagClient(endpoint string, cred azcore.TokenCredential, options *FeatureFlagClientOptions) (*FeatureFlagClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	audienceStr := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if options != nil && !reflect.ValueOf(options.Cloud).IsZero() {
		if cfg, ok := options.Cloud.Services[ServiceName]; ok && cfg.Audience != "" {
			audienceStr = cfg.Audience
		}
	}

	return newFeatureFlagClient(endpoint, runtime.NewBearerTokenPolicy(cred, []string{audienceStr + "/.default"}, nil), options)
}

// NewFeatureFlagClientFromConnectionString parses the connection string and returns a pointer to a FeatureFlagClient object.
func NewFeatureFlagClientFromConnectionString(connectionString string, options *FeatureFlagClientOptions) (*FeatureFlagClient, error) {
	endpoint, credential, secret, err := auth.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return newFeatureFlagClient(endpoint, auth.NewHMACPolicy(credential, secret), options)
}

func newFeatureFlagClient(endpoint string, authPolicy policy.Policy, options *FeatureFlagClientOptions) (*FeatureFlagClient, error) {
	if options == nil {
		options = &FeatureFlagClientOptions{}
	}

	audienceConfigured := false
	if !reflect.ValueOf(options.Cloud).IsZero() {
		if cfg, ok := options.Cloud.Services[ServiceName]; ok && cfg.Audience != "" {
			audienceConfigured = true
		}
	}

	cache := synctoken.NewCache()
	client, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{querynormalization.NewPolicy(), authPolicy, synctoken.NewPolicy(cache), audience.NewAudienceErrorHandlingPolicy(audienceConfigured)},
		Tracing: runtime.TracingOptions{
			Namespace: "Microsoft.AppConfig",
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	appConfigClient := generated.NewAzureAppConfigurationClient(endpoint, client)
	return &FeatureFlagClient{
		appConfigClient: appConfigClient,
		ffClient:        appConfigClient.NewAzureAppConfigurationFeatureFlagClient(),
		cache:           cache,
	}, nil
}

// NewFeatureFlagClient returns a FeatureFlagClient that shares the pipeline and sync-token
// cache of this [Client]. Use this when a Client has already been constructed and you also
// need to manage feature flags.
func (c *Client) NewFeatureFlagClient() *FeatureFlagClient {
	return &FeatureFlagClient{
		appConfigClient: c.appConfigClient,
		ffClient:        c.appConfigClient.NewAzureAppConfigurationFeatureFlagClient(),
		cache:           c.cache,
	}
}

// SetSyncToken is used to set a sync token from an external source.
// SyncTokens are required to be in the format "<id>=<value>;sn=<sn>".
// Multiple SyncTokens must be comma delimited.
func (c *FeatureFlagClient) SetSyncToken(syncToken SyncToken) error {
	return c.cache.Set(syncToken)
}

// AddFeatureFlag creates a feature flag only if the flag does not already exist in the configuration store.
//   - ctx controls the lifetime of the HTTP operation
//   - flag is the feature flag to create. flag.Name is required.
//   - options contains the optional values. can be nil
func (c *FeatureFlagClient) AddFeatureFlag(ctx context.Context, flag FeatureFlag, options *AddFeatureFlagOptions) (AddFeatureFlagResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "FeatureFlagClient.AddFeatureFlag", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if flag.Name == nil {
		err = errors.New("flag.Name is required")
		return AddFeatureFlagResponse{}, err
	}

	etagAny := azcore.ETagAny
	opts := &generated.AzureAppConfigurationFeatureFlagClientPutFeatureFlagOptions{
		Label:       flag.Label,
		IfNoneMatch: toGeneratedETagString(&etagAny),
	}

	resp, err := c.ffClient.PutFeatureFlag(ctx, *flag.Name, flag.toGenerated(), opts)
	if err != nil {
		return AddFeatureFlagResponse{}, err
	}

	return AddFeatureFlagResponse{
		FeatureFlag: featureFlagFromGenerated(resp.FeatureFlag),
		SyncToken:   syncTokenValue(resp.SyncToken),
	}, nil
}

// SetFeatureFlag creates a feature flag if it doesn't exist or overwrites the existing flag in the configuration store.
//   - ctx controls the lifetime of the HTTP operation
//   - flag is the feature flag to create or update. flag.Name is required.
//   - options contains the optional values. can be nil
func (c *FeatureFlagClient) SetFeatureFlag(ctx context.Context, flag FeatureFlag, options *SetFeatureFlagOptions) (SetFeatureFlagResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "FeatureFlagClient.SetFeatureFlag", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if flag.Name == nil {
		err = errors.New("flag.Name is required")
		return SetFeatureFlagResponse{}, err
	}

	if options == nil {
		options = &SetFeatureFlagOptions{}
	}

	opts := &generated.AzureAppConfigurationFeatureFlagClientPutFeatureFlagOptions{
		Label:   flag.Label,
		IfMatch: toGeneratedETagString(options.OnlyIfUnchanged),
	}

	resp, err := c.ffClient.PutFeatureFlag(ctx, *flag.Name, flag.toGenerated(), opts)
	if err != nil {
		return SetFeatureFlagResponse{}, err
	}

	return SetFeatureFlagResponse{
		FeatureFlag: featureFlagFromGenerated(resp.FeatureFlag),
		SyncToken:   syncTokenValue(resp.SyncToken),
	}, nil
}

// GetFeatureFlag retrieves an existing feature flag from the configuration store.
func (c *FeatureFlagClient) GetFeatureFlag(ctx context.Context, name string, options *GetFeatureFlagOptions) (GetFeatureFlagResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "FeatureFlagClient.GetFeatureFlag", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetFeatureFlagOptions{}
	}

	var dt *string
	if options.AcceptDateTime != nil {
		str := options.AcceptDateTime.Format(timeFormat)
		dt = &str
	}

	opts := &generated.AzureAppConfigurationFeatureFlagClientGetFeatureFlagOptions{
		AcceptDatetime: dt,
		Label:          options.Label,
		IfNoneMatch:    toGeneratedETagString(options.OnlyIfChanged),
	}

	resp, err := c.ffClient.GetFeatureFlag(ctx, name, opts)
	if err != nil {
		return GetFeatureFlagResponse{}, err
	}

	return GetFeatureFlagResponse{
		FeatureFlag: featureFlagFromGenerated(resp.FeatureFlag),
		SyncToken:   syncTokenValue(resp.SyncToken),
	}, nil
}

// DeleteFeatureFlag deletes a feature flag from the configuration store.
func (c *FeatureFlagClient) DeleteFeatureFlag(ctx context.Context, name string, options *DeleteFeatureFlagOptions) (DeleteFeatureFlagResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "FeatureFlagClient.DeleteFeatureFlag", c.appConfigClient.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &DeleteFeatureFlagOptions{}
	}

	opts := &generated.AzureAppConfigurationFeatureFlagClientDeleteFeatureFlagOptions{
		Label:   options.Label,
		IfMatch: toGeneratedETagString(options.OnlyIfUnchanged),
	}

	resp, err := c.ffClient.DeleteFeatureFlag(ctx, name, opts)
	if err != nil {
		return DeleteFeatureFlagResponse{}, err
	}

	return DeleteFeatureFlagResponse{
		FeatureFlag: featureFlagFromGenerated(resp.FeatureFlag),
		SyncToken:   syncTokenValue(resp.SyncToken),
	}, nil
}

// NewListFeatureFlagsPager creates a pager that retrieves feature flags that match the specified selector.
func (c *FeatureFlagClient) NewListFeatureFlagsPager(selector FeatureFlagSelector, options *ListFeatureFlagsOptions) *runtime.Pager[ListFeatureFlagsPageResponse] {
	if options == nil {
		options = &ListFeatureFlagsOptions{}
	}
	_ = options // reserved for future match-conditions support
	pagerInternal := c.ffClient.NewGetFeatureFlagsPager(selector.toGeneratedGetFeatureFlags())
	return runtime.NewPager(runtime.PagingHandler[ListFeatureFlagsPageResponse]{
		More: func(ListFeatureFlagsPageResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListFeatureFlagsPageResponse) (ListFeatureFlagsPageResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListFeatureFlagsPageResponse{}, err
			}
			flags := make([]FeatureFlag, len(page.Items))
			for i := range page.Items {
				flags[i] = featureFlagFromGenerated(page.Items[i])
			}
			return ListFeatureFlagsPageResponse{
				FeatureFlags: flags,
				ETag:         (*azcore.ETag)(page.Etag),
				SyncToken:    syncTokenValue(page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

// NewListFeatureFlagRevisionsPager creates a pager that retrieves revisions of feature flags that match the specified selector.
func (c *FeatureFlagClient) NewListFeatureFlagRevisionsPager(selector FeatureFlagSelector, options *ListFeatureFlagRevisionsOptions) *runtime.Pager[ListFeatureFlagRevisionsPageResponse] {
	pagerInternal := c.ffClient.NewGetFeatureFlagRevisionsPager(selector.toGeneratedGetFeatureFlagRevisions())
	return runtime.NewPager(runtime.PagingHandler[ListFeatureFlagRevisionsPageResponse]{
		More: func(ListFeatureFlagRevisionsPageResponse) bool {
			return pagerInternal.More()
		},
		Fetcher: func(ctx context.Context, cur *ListFeatureFlagRevisionsPageResponse) (ListFeatureFlagRevisionsPageResponse, error) {
			page, err := pagerInternal.NextPage(ctx)
			if err != nil {
				return ListFeatureFlagRevisionsPageResponse{}, err
			}
			flags := make([]FeatureFlag, len(page.Items))
			for i := range page.Items {
				flags[i] = featureFlagFromGenerated(page.Items[i])
			}
			return ListFeatureFlagRevisionsPageResponse{
				FeatureFlags: flags,
				SyncToken:    syncTokenValue(page.SyncToken),
			}, nil
		},
		Tracer: c.appConfigClient.Tracer(),
	})
}

func syncTokenValue(v *string) SyncToken {
	if v == nil {
		return ""
	}
	return SyncToken(*v)
}
