package recoveryservices

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "net/http"
    "context"
    "github.com/Azure/go-autorest/tracing"
)

// UsagesClient is the recovery Services Client
type UsagesClient struct {
    BaseClient
}
// NewUsagesClient creates an instance of the UsagesClient client.
func NewUsagesClient(subscriptionID string) UsagesClient {
    return NewUsagesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewUsagesClientWithBaseURI creates an instance of the UsagesClient client.
    func NewUsagesClientWithBaseURI(baseURI string, subscriptionID string) UsagesClient {
        return UsagesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// ListByVaults fetches the usages of the vault.
    // Parameters:
        // resourceGroupName - the name of the resource group where the recovery services vault is present.
        // vaultName - the name of the recovery services vault.
func (client UsagesClient) ListByVaults(ctx context.Context, resourceGroupName string, vaultName string) (result VaultUsageList, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/UsagesClient.ListByVaults")
        defer func() {
            sc := -1
            if result.Response.Response != nil {
                sc = result.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
        req, err := client.ListByVaultsPreparer(ctx, resourceGroupName, vaultName)
    if err != nil {
    err = autorest.NewErrorWithError(err, "recoveryservices.UsagesClient", "ListByVaults", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListByVaultsSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "recoveryservices.UsagesClient", "ListByVaults", resp, "Failure sending request")
            return
            }

            result, err = client.ListByVaultsResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "recoveryservices.UsagesClient", "ListByVaults", resp, "Failure responding to request")
            }

    return
    }

    // ListByVaultsPreparer prepares the ListByVaults request.
    func (client UsagesClient) ListByVaultsPreparer(ctx context.Context, resourceGroupName string, vaultName string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "resourceGroupName": autorest.Encode("path",resourceGroupName),
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            "vaultName": autorest.Encode("path",vaultName),
            }

                        const APIVersion = "2016-06-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/usages",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListByVaultsSender sends the ListByVaults request. The method will close the
    // http.Response Body if it receives an error.
    func (client UsagesClient) ListByVaultsSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListByVaultsResponder handles the response to the ListByVaults request. The method always
// closes the http.Response Body.
func (client UsagesClient) ListByVaultsResponder(resp *http.Response) (result VaultUsageList, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

