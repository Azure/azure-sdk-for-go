//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmysqlflexibleservers

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ServersMigrationClient contains the methods for the ServersMigration group.
// Don't use this type directly, use NewServersMigrationClient() instead.
type ServersMigrationClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewServersMigrationClient creates a new instance of ServersMigrationClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewServersMigrationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ServersMigrationClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ServersMigrationClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCutoverMigration - Cutover migration for MySQL import, it will switch source elastic server DNS to flexible server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - options - ServersMigrationClientBeginCutoverMigrationOptions contains the optional parameters for the ServersMigrationClient.BeginCutoverMigration
//     method.
func (client *ServersMigrationClient) BeginCutoverMigration(ctx context.Context, resourceGroupName string, serverName string, options *ServersMigrationClientBeginCutoverMigrationOptions) (*runtime.Poller[ServersMigrationClientCutoverMigrationResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.cutoverMigration(ctx, resourceGroupName, serverName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ServersMigrationClientCutoverMigrationResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ServersMigrationClientCutoverMigrationResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CutoverMigration - Cutover migration for MySQL import, it will switch source elastic server DNS to flexible server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *ServersMigrationClient) cutoverMigration(ctx context.Context, resourceGroupName string, serverName string, options *ServersMigrationClientBeginCutoverMigrationOptions) (*http.Response, error) {
	var err error
	const operationName = "ServersMigrationClient.BeginCutoverMigration"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.cutoverMigrationCreateRequest(ctx, resourceGroupName, serverName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// cutoverMigrationCreateRequest creates the CutoverMigration request.
func (client *ServersMigrationClient) cutoverMigrationCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ServersMigrationClientBeginCutoverMigrationOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{serverName}/cutoverMigration"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}
