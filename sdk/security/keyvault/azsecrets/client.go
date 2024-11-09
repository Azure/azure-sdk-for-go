// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package azsecrets

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// Client - The key vault client performs cryptographic key operations and vault operations
// against the Key Vault service.
// Don't use this type directly, use a constructor function instead.
type Client struct {
	internal     *azcore.Client
	vaultBaseUrl string
}

// BackupSecret - Backs up the specified secret.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - options - BackupSecretOptions contains the optional parameters for the Client.BackupSecret method.
func (client *Client) BackupSecret(ctx context.Context, name string, options *BackupSecretOptions) (BackupSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.BackupSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.backupSecretCreateRequest(ctx, name, options)
	if err != nil {
		return BackupSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return BackupSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return BackupSecretResponse{}, err
	}
	resp, err := client.backupSecretHandleResponse(httpResp)
	return resp, err
}

// backupSecretCreateRequest creates the BackupSecret request.
func (client *Client) backupSecretCreateRequest(ctx context.Context, name string, _ *BackupSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}/backup"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// backupSecretHandleResponse handles the BackupSecret response.
func (client *Client) backupSecretHandleResponse(resp *http.Response) (BackupSecretResponse, error) {
	result := BackupSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackupSecretResult); err != nil {
		return BackupSecretResponse{}, err
	}
	return result, nil
}

// DeleteSecret - Deletes a secret from a specified key vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - options - DeleteSecretOptions contains the optional parameters for the Client.DeleteSecret method.
func (client *Client) DeleteSecret(ctx context.Context, name string, options *DeleteSecretOptions) (DeleteSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.DeleteSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteSecretCreateRequest(ctx, name, options)
	if err != nil {
		return DeleteSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeleteSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DeleteSecretResponse{}, err
	}
	resp, err := client.deleteSecretHandleResponse(httpResp)
	return resp, err
}

// deleteSecretCreateRequest creates the DeleteSecret request.
func (client *Client) deleteSecretCreateRequest(ctx context.Context, name string, _ *DeleteSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteSecretHandleResponse handles the DeleteSecret response.
func (client *Client) deleteSecretHandleResponse(resp *http.Response) (DeleteSecretResponse, error) {
	result := DeleteSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedSecret); err != nil {
		return DeleteSecretResponse{}, err
	}
	return result, nil
}

// GetDeletedSecret - Gets the specified deleted secret.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - options - GetDeletedSecretOptions contains the optional parameters for the Client.GetDeletedSecret method.
func (client *Client) GetDeletedSecret(ctx context.Context, name string, options *GetDeletedSecretOptions) (GetDeletedSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetDeletedSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getDeletedSecretCreateRequest(ctx, name, options)
	if err != nil {
		return GetDeletedSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetDeletedSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GetDeletedSecretResponse{}, err
	}
	resp, err := client.getDeletedSecretHandleResponse(httpResp)
	return resp, err
}

// getDeletedSecretCreateRequest creates the GetDeletedSecret request.
func (client *Client) getDeletedSecretCreateRequest(ctx context.Context, name string, _ *GetDeletedSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/deletedsecrets/{secret-name}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getDeletedSecretHandleResponse handles the GetDeletedSecret response.
func (client *Client) getDeletedSecretHandleResponse(resp *http.Response) (GetDeletedSecretResponse, error) {
	result := GetDeletedSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedSecret); err != nil {
		return GetDeletedSecretResponse{}, err
	}
	return result, nil
}

// GetSecret - Get a specified secret from a given key vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - version - The version of the secret. This URI fragment is optional. If not specified, the
//     latest version of the secret is returned.
//   - options - GetSecretOptions contains the optional parameters for the Client.GetSecret method.
func (client *Client) GetSecret(ctx context.Context, name string, version string, options *GetSecretOptions) (GetSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getSecretCreateRequest(ctx, name, version, options)
	if err != nil {
		return GetSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GetSecretResponse{}, err
	}
	resp, err := client.getSecretHandleResponse(httpResp)
	return resp, err
}

// getSecretCreateRequest creates the GetSecret request.
func (client *Client) getSecretCreateRequest(ctx context.Context, name string, version string, _ *GetSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}/{secret-version}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	urlPath = strings.ReplaceAll(urlPath, "{secret-version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getSecretHandleResponse handles the GetSecret response.
func (client *Client) getSecretHandleResponse(resp *http.Response) (GetSecretResponse, error) {
	result := GetSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return GetSecretResponse{}, err
	}
	return result, nil
}

// NewListDeletedSecretPropertiesPager - Lists deleted secrets for the specified vault.
//
// Generated from API version 7.5
//   - options - ListDeletedSecretPropertiesOptions contains the optional parameters for the Client.NewListDeletedSecretPropertiesPager
//     method.
func (client *Client) NewListDeletedSecretPropertiesPager(options *ListDeletedSecretPropertiesOptions) *runtime.Pager[ListDeletedSecretPropertiesResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListDeletedSecretPropertiesResponse]{
		More: func(page ListDeletedSecretPropertiesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListDeletedSecretPropertiesResponse) (ListDeletedSecretPropertiesResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listDeletedSecretPropertiesCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return ListDeletedSecretPropertiesResponse{}, err
			}
			return client.listDeletedSecretPropertiesHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listDeletedSecretPropertiesCreateRequest creates the ListDeletedSecretProperties request.
func (client *Client) listDeletedSecretPropertiesCreateRequest(ctx context.Context, options *ListDeletedSecretPropertiesOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/deletedsecrets"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listDeletedSecretPropertiesHandleResponse handles the ListDeletedSecretProperties response.
func (client *Client) listDeletedSecretPropertiesHandleResponse(resp *http.Response) (ListDeletedSecretPropertiesResponse, error) {
	result := ListDeletedSecretPropertiesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedSecretPropertiesListResult); err != nil {
		return ListDeletedSecretPropertiesResponse{}, err
	}
	return result, nil
}

// NewListSecretPropertiesPager - List secrets in a specified key vault.
//
// Generated from API version 7.5
//   - options - ListSecretPropertiesOptions contains the optional parameters for the Client.NewListSecretPropertiesPager method.
func (client *Client) NewListSecretPropertiesPager(options *ListSecretPropertiesOptions) *runtime.Pager[ListSecretPropertiesResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListSecretPropertiesResponse]{
		More: func(page ListSecretPropertiesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListSecretPropertiesResponse) (ListSecretPropertiesResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listSecretPropertiesCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return ListSecretPropertiesResponse{}, err
			}
			return client.listSecretPropertiesHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listSecretPropertiesCreateRequest creates the ListSecretProperties request.
func (client *Client) listSecretPropertiesCreateRequest(ctx context.Context, options *ListSecretPropertiesOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSecretPropertiesHandleResponse handles the ListSecretProperties response.
func (client *Client) listSecretPropertiesHandleResponse(resp *http.Response) (ListSecretPropertiesResponse, error) {
	result := ListSecretPropertiesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretPropertiesListResult); err != nil {
		return ListSecretPropertiesResponse{}, err
	}
	return result, nil
}

// NewListSecretPropertiesVersionsPager - List all versions of the specified secret.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - options - ListSecretPropertiesVersionsOptions contains the optional parameters for the Client.NewListSecretPropertiesVersionsPager
//     method.
func (client *Client) NewListSecretPropertiesVersionsPager(name string, options *ListSecretPropertiesVersionsOptions) *runtime.Pager[ListSecretPropertiesVersionsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListSecretPropertiesVersionsResponse]{
		More: func(page ListSecretPropertiesVersionsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListSecretPropertiesVersionsResponse) (ListSecretPropertiesVersionsResponse, error) {
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listSecretPropertiesVersionsCreateRequest(ctx, name, options)
			}, nil)
			if err != nil {
				return ListSecretPropertiesVersionsResponse{}, err
			}
			return client.listSecretPropertiesVersionsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listSecretPropertiesVersionsCreateRequest creates the ListSecretPropertiesVersions request.
func (client *Client) listSecretPropertiesVersionsCreateRequest(ctx context.Context, name string, _ *ListSecretPropertiesVersionsOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}/versions"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSecretPropertiesVersionsHandleResponse handles the ListSecretPropertiesVersions response.
func (client *Client) listSecretPropertiesVersionsHandleResponse(resp *http.Response) (ListSecretPropertiesVersionsResponse, error) {
	result := ListSecretPropertiesVersionsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SecretPropertiesListResult); err != nil {
		return ListSecretPropertiesVersionsResponse{}, err
	}
	return result, nil
}

// PurgeDeletedSecret - Permanently deletes the specified secret.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - options - PurgeDeletedSecretOptions contains the optional parameters for the Client.PurgeDeletedSecret method.
func (client *Client) PurgeDeletedSecret(ctx context.Context, name string, options *PurgeDeletedSecretOptions) (PurgeDeletedSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.PurgeDeletedSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.purgeDeletedSecretCreateRequest(ctx, name, options)
	if err != nil {
		return PurgeDeletedSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PurgeDeletedSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return PurgeDeletedSecretResponse{}, err
	}
	return PurgeDeletedSecretResponse{}, nil
}

// purgeDeletedSecretCreateRequest creates the PurgeDeletedSecret request.
func (client *Client) purgeDeletedSecretCreateRequest(ctx context.Context, name string, _ *PurgeDeletedSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/deletedsecrets/{secret-name}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// RecoverDeletedSecret - Recovers the deleted secret to the latest version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the deleted secret.
//   - options - RecoverDeletedSecretOptions contains the optional parameters for the Client.RecoverDeletedSecret method.
func (client *Client) RecoverDeletedSecret(ctx context.Context, name string, options *RecoverDeletedSecretOptions) (RecoverDeletedSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.RecoverDeletedSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.recoverDeletedSecretCreateRequest(ctx, name, options)
	if err != nil {
		return RecoverDeletedSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RecoverDeletedSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RecoverDeletedSecretResponse{}, err
	}
	resp, err := client.recoverDeletedSecretHandleResponse(httpResp)
	return resp, err
}

// recoverDeletedSecretCreateRequest creates the RecoverDeletedSecret request.
func (client *Client) recoverDeletedSecretCreateRequest(ctx context.Context, name string, _ *RecoverDeletedSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/deletedsecrets/{secret-name}/recover"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// recoverDeletedSecretHandleResponse handles the RecoverDeletedSecret response.
func (client *Client) recoverDeletedSecretHandleResponse(resp *http.Response) (RecoverDeletedSecretResponse, error) {
	result := RecoverDeletedSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return RecoverDeletedSecretResponse{}, err
	}
	return result, nil
}

// RestoreSecret - Restores a backed up secret to a vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - parameters - The parameters to restore the secret.
//   - options - RestoreSecretOptions contains the optional parameters for the Client.RestoreSecret method.
func (client *Client) RestoreSecret(ctx context.Context, parameters RestoreSecretParameters, options *RestoreSecretOptions) (RestoreSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.RestoreSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.restoreSecretCreateRequest(ctx, parameters, options)
	if err != nil {
		return RestoreSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RestoreSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RestoreSecretResponse{}, err
	}
	resp, err := client.restoreSecretHandleResponse(httpResp)
	return resp, err
}

// restoreSecretCreateRequest creates the RestoreSecret request.
func (client *Client) restoreSecretCreateRequest(ctx context.Context, parameters RestoreSecretParameters, _ *RestoreSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/restore"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// restoreSecretHandleResponse handles the RestoreSecret response.
func (client *Client) restoreSecretHandleResponse(resp *http.Response) (RestoreSecretResponse, error) {
	result := RestoreSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return RestoreSecretResponse{}, err
	}
	return result, nil
}

// SetSecret - Sets a secret in a specified key vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret. The value you provide may be copied globally for the
//     purpose of running the service. The value provided should not include
//     personally identifiable or sensitive information.
//   - parameters - The parameters for setting the secret.
//   - options - SetSecretOptions contains the optional parameters for the Client.SetSecret method.
func (client *Client) SetSecret(ctx context.Context, name string, parameters SetSecretParameters, options *SetSecretOptions) (SetSecretResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.SetSecret", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.setSecretCreateRequest(ctx, name, parameters, options)
	if err != nil {
		return SetSecretResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SetSecretResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SetSecretResponse{}, err
	}
	resp, err := client.setSecretHandleResponse(httpResp)
	return resp, err
}

// setSecretCreateRequest creates the SetSecret request.
func (client *Client) setSecretCreateRequest(ctx context.Context, name string, parameters SetSecretParameters, _ *SetSecretOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// setSecretHandleResponse handles the SetSecret response.
func (client *Client) setSecretHandleResponse(resp *http.Response) (SetSecretResponse, error) {
	result := SetSecretResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return SetSecretResponse{}, err
	}
	return result, nil
}

// UpdateSecretProperties - Updates the attributes associated with a specified secret in a given key vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 7.5
//   - name - The name of the secret.
//   - version - The version of the secret.
//   - parameters - The parameters for update secret operation.
//   - options - UpdateSecretPropertiesOptions contains the optional parameters for the Client.UpdateSecretProperties method.
func (client *Client) UpdateSecretProperties(ctx context.Context, name string, version string, parameters UpdateSecretPropertiesParameters, options *UpdateSecretPropertiesOptions) (UpdateSecretPropertiesResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.UpdateSecretProperties", client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateSecretPropertiesCreateRequest(ctx, name, version, parameters, options)
	if err != nil {
		return UpdateSecretPropertiesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return UpdateSecretPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return UpdateSecretPropertiesResponse{}, err
	}
	resp, err := client.updateSecretPropertiesHandleResponse(httpResp)
	return resp, err
}

// updateSecretPropertiesCreateRequest creates the UpdateSecretProperties request.
func (client *Client) updateSecretPropertiesCreateRequest(ctx context.Context, name string, version string, parameters UpdateSecretPropertiesParameters, _ *UpdateSecretPropertiesOptions) (*policy.Request, error) {
	host := "{vaultBaseUrl}"
	host = strings.ReplaceAll(host, "{vaultBaseUrl}", client.vaultBaseUrl)
	urlPath := "/secrets/{secret-name}/{secret-version}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{secret-name}", url.PathEscape(name))
	urlPath = strings.ReplaceAll(urlPath, "{secret-version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.5")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateSecretPropertiesHandleResponse handles the UpdateSecretProperties response.
func (client *Client) updateSecretPropertiesHandleResponse(resp *http.Response) (UpdateSecretPropertiesResponse, error) {
	result := UpdateSecretPropertiesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Secret); err != nil {
		return UpdateSecretPropertiesResponse{}, err
	}
	return result, nil
}
