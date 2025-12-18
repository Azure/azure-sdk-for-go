//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure File Storage service allowing you to manipulate file shares.
type Client base.Client[generated.ServiceClient]

// NewClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
//
// Note that service-level operations do not support token credential authentication.
// This constructor exists to allow the construction of a share.Client that has token credential authentication.
// Also note that ClientOptions.FileRequestIntent is currently required for token authentication.
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	audience := base.GetAudience((*base.ClientOptions)(options))
	conOptions := shared.GetClientOptions(options)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{audience}, &policy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: conOptions.InsecureAllowCredentialWithHTTP,
	})
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewServiceClient(serviceURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewServiceClient(serviceURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewServiceClient(serviceURL, azClient, cred, (*base.ClientOptions)(conOptions))), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (s *Client) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) getClientOptions() *base.ClientOptions {
	return base.GetClientOptions((*base.Client[generated.ServiceClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return s.generated().Endpoint()
}

// GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceURL object.
// OAuth is required for this call, as well as any role that can delegate access to the storage account.
func (s *Client) GetUserDelegationCredential(ctx context.Context, info KeyInfo, o *GetUserDelegationCredentialOptions) (*UserDelegationCredential, error) {
	url := s.URL()
	parts := strings.Split(strings.TrimPrefix(url, "https://"), ".")
	account := parts[0]

	getUserDelegationKeyOptions := o.format()
	udk, err := s.generated().GetUserDelegationKey(ctx, info, getUserDelegationKeyOptions)
	if err != nil {
		return nil, err
	}

	return exported.NewUserDelegationCredential(account, udk.UserDelegationKey), nil
}

// NewShareClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewShareClient(shareName string) *share.Client {
	shareURL := runtime.JoinPaths(s.generated().Endpoint(), shareName)
	return (*share.Client)(base.NewShareClient(shareURL, s.generated().InternalClient().WithClientName(exported.ModuleName), s.sharedKey(), s.getClientOptions()))
}

// CreateShare is a lifecycle method to creates a new share under the specified account.
// If the share with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created share.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/create-share.
func (s *Client) CreateShare(ctx context.Context, shareName string, options *CreateShareOptions) (CreateShareResponse, error) {
	shareClient := s.NewShareClient(shareName)
	createShareResp, err := shareClient.Create(ctx, options)
	return createShareResp, err
}

// DeleteShare is a lifecycle method that marks the specified share for deletion.
// The share and any files contained within it are later deleted during garbage collection.
// If the share is not found, a ResourceNotFoundError will be raised.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/delete-share.
func (s *Client) DeleteShare(ctx context.Context, shareName string, options *DeleteShareOptions) (DeleteShareResponse, error) {
	shareClient := s.NewShareClient(shareName)
	deleteShareResp, err := shareClient.Delete(ctx, options)
	return deleteShareResp, err
}

// RestoreShare restores soft-deleted share.
// Operation will only be successful if used within the specified number of days set in the delete retention policy.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/restore-share.
func (s *Client) RestoreShare(ctx context.Context, deletedShareName string, deletedShareVersion string, options *RestoreShareOptions) (RestoreShareResponse, error) {
	shareClient := s.NewShareClient(deletedShareName)
	createShareResp, err := shareClient.Restore(ctx, deletedShareVersion, options)
	return createShareResp, err
}

// GetProperties operation gets the properties of a storage account's File service.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file-service-properties.
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	resp, err := s.generated().GetProperties(ctx, opts)
	return resp, err
}

// SetProperties operation sets properties for a storage account's File service endpoint.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-service-properties.
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	svcProperties, o := options.format()
	resp, err := s.generated().SetProperties(ctx, svcProperties, o)
	return resp, err
}

// NewListSharesPager operation returns a pager of the shares under the specified account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (s *Client) NewListSharesPager(options *ListSharesOptions) *runtime.Pager[ListSharesSegmentResponse] {
	listOptions := generated.ServiceClientListSharesSegmentOptions{}
	if options != nil {
		if options.Include.Deleted {
			listOptions.Include = append(listOptions.Include, ListSharesIncludeTypeDeleted)
		}
		if options.Include.Metadata {
			listOptions.Include = append(listOptions.Include, ListSharesIncludeTypeMetadata)
		}
		if options.Include.Snapshots {
			listOptions.Include = append(listOptions.Include, ListSharesIncludeTypeSnapshots)
		}
		listOptions.Marker = options.Marker
		listOptions.Maxresults = options.MaxResults
		listOptions.Prefix = options.Prefix
	}

	return runtime.NewPager(runtime.PagingHandler[ListSharesSegmentResponse]{
		More: func(page ListSharesSegmentResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListSharesSegmentResponse) (ListSharesSegmentResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = s.generated().ListSharesSegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.NextMarker
				req, err = s.generated().ListSharesSegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListSharesSegmentResponse{}, err
			}
			resp, err := s.generated().InternalClient().Pipeline().Do(req)
			if err != nil {
				return ListSharesSegmentResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListSharesSegmentResponse{}, runtime.NewResponseError(resp)
			}
			return s.generated().ListSharesSegmentHandleResponse(resp)
		},
	})
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (s *Client) GetSASURL(resources sas.AccountResourceTypes, permissions sas.AccountPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if s.sharedKey() == nil {
		return "", fileerror.MissingSharedKeyCredential
	}
	st := o.format()
	qps, err := sas.AccountSignatureValues{
		Version:       sas.Version,
		Permissions:   permissions.String(),
		ResourceTypes: resources.String(),
		StartTime:     st,
		ExpiryTime:    expiry.UTC(),
	}.SignWithSharedKey(s.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := s.URL()
	if !strings.HasSuffix(endpoint, "/") {
		// add a trailing slash to be consistent with the portal
		endpoint += "/"
	}
	endpoint += "?" + qps.Encode()

	return endpoint, nil
}
