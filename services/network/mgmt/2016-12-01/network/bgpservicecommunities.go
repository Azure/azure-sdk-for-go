package network

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

// BgpServiceCommunitiesClient is the network Client
type BgpServiceCommunitiesClient struct {
    BaseClient
}
// NewBgpServiceCommunitiesClient creates an instance of the BgpServiceCommunitiesClient client.
func NewBgpServiceCommunitiesClient(subscriptionID string) BgpServiceCommunitiesClient {
    return NewBgpServiceCommunitiesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBgpServiceCommunitiesClientWithBaseURI creates an instance of the BgpServiceCommunitiesClient client.
    func NewBgpServiceCommunitiesClientWithBaseURI(baseURI string, subscriptionID string) BgpServiceCommunitiesClient {
        return BgpServiceCommunitiesClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// List gets all the available bgp service communities.
func (client BgpServiceCommunitiesClient) List(ctx context.Context) (result BgpServiceCommunityListResultPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/BgpServiceCommunitiesClient.List")
        defer func() {
            sc := -1
            if result.bsclr.Response.Response != nil {
                sc = result.bsclr.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx)
    if err != nil {
    err = autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.bsclr.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "List", resp, "Failure sending request")
            return
            }

            result.bsclr, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client BgpServiceCommunitiesClient) ListPreparer(ctx context.Context) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2016-12-01"
        queryParameters := map[string]interface{} {
        "api-version": APIVersion,
        }

    preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Network/bgpServiceCommunities",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client BgpServiceCommunitiesClient) ListSender(req *http.Request) (*http.Response, error) {
            return autorest.SendWithSender(client, req,
            azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client BgpServiceCommunitiesClient) ListResponder(resp *http.Response) (result BgpServiceCommunityListResult, err error) {
    err = autorest.Respond(
    resp,
    client.ByInspecting(),
    azure.WithErrorUnlessStatusCode(http.StatusOK),
    autorest.ByUnmarshallingJSON(&result),
    autorest.ByClosing())
    result.Response = autorest.Response{Response: resp}
        return
    }

            // listNextResults retrieves the next set of results, if any.
            func (client BgpServiceCommunitiesClient) listNextResults(ctx context.Context, lastResults BgpServiceCommunityListResult) (result BgpServiceCommunityListResult, err error) {
            req, err := lastResults.bgpServiceCommunityListResultPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "network.BgpServiceCommunitiesClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client BgpServiceCommunitiesClient) ListComplete(ctx context.Context) (result BgpServiceCommunityListResultIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/BgpServiceCommunitiesClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx)
                return
        }

