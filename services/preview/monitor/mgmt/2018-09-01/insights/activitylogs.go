package insights

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

// ActivityLogsClient is the monitor Management Client
type ActivityLogsClient struct {
    BaseClient
}
// NewActivityLogsClient creates an instance of the ActivityLogsClient client.
func NewActivityLogsClient(subscriptionID string) ActivityLogsClient {
    return NewActivityLogsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewActivityLogsClientWithBaseURI creates an instance of the ActivityLogsClient client using a custom endpoint.  Use
// this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
    func NewActivityLogsClientWithBaseURI(baseURI string, subscriptionID string) ActivityLogsClient {
        return ActivityLogsClient{ NewWithBaseURI(baseURI, subscriptionID)}
    }

// List provides the list of records from the activity logs.
    // Parameters:
        // filter - reduces the set of data collected.<br>This argument is required and it also requires at least the
        // start date/time.<br>The **$filter** argument is very restricted and allows only the following patterns.<br>-
        // *List events for a resource group*: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and
        // eventTimestamp le '2014-07-20T04:36:37.6407898Z' and resourceGroupName eq 'resourceGroupName'.<br>- *List
        // events for resource*: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le
        // '2014-07-20T04:36:37.6407898Z' and resourceUri eq 'resourceURI'.<br>- *List events for a subscription in a
        // time range*: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le
        // '2014-07-20T04:36:37.6407898Z'.<br>- *List events for a resource provider*: $filter=eventTimestamp ge
        // '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z' and resourceProvider eq
        // 'resourceProviderName'.<br>- *List events for a correlation Id*: $filter=eventTimestamp ge
        // '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z' and correlationId eq
        // 'correlationID'.<br><br>**NOTE**: No other syntax is allowed.
        // selectParameter - used to fetch events with only the given properties.<br>The **$select** argument is a
        // comma separated list of property names to be returned. Possible values are: *authorization*, *claims*,
        // *correlationId*, *description*, *eventDataId*, *eventName*, *eventTimestamp*, *httpRequest*, *level*,
        // *operationId*, *operationName*, *properties*, *resourceGroupName*, *resourceProviderName*, *resourceId*,
        // *status*, *submissionTimestamp*, *subStatus*, *subscriptionId*
func (client ActivityLogsClient) List(ctx context.Context, filter string, selectParameter string) (result EventDataCollectionPage, err error) {
    if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ActivityLogsClient.List")
        defer func() {
            sc := -1
            if result.edc.Response.Response != nil {
                sc = result.edc.Response.Response.StatusCode
            }
            tracing.EndSpan(ctx, sc, err)
        }()
    }
                result.fn = client.listNextResults
    req, err := client.ListPreparer(ctx, filter, selectParameter)
    if err != nil {
    err = autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "List", nil , "Failure preparing request")
    return
    }

            resp, err := client.ListSender(req)
            if err != nil {
            result.edc.Response = autorest.Response{Response: resp}
            err = autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "List", resp, "Failure sending request")
            return
            }

            result.edc, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "List", resp, "Failure responding to request")
            }

    return
    }

    // ListPreparer prepares the List request.
    func (client ActivityLogsClient) ListPreparer(ctx context.Context, filter string, selectParameter string) (*http.Request, error) {
            pathParameters := map[string]interface{} {
            "subscriptionId": autorest.Encode("path",client.SubscriptionID),
            }

                        const APIVersion = "2015-04-01"
        queryParameters := map[string]interface{} {
        "$filter": autorest.Encode("query",filter),
        "api-version": APIVersion,
        }
            if len(selectParameter) > 0 {
            queryParameters["$select"] = autorest.Encode("query",selectParameter)
            }

        preparer := autorest.CreatePreparer(
    autorest.AsGet(),
    autorest.WithBaseURL(client.BaseURI),
    autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/microsoft.insights/eventtypes/management/values",pathParameters),
    autorest.WithQueryParameters(queryParameters))
    return preparer.Prepare((&http.Request{}).WithContext(ctx))
    }

    // ListSender sends the List request. The method will close the
    // http.Response Body if it receives an error.
    func (client ActivityLogsClient) ListSender(req *http.Request) (*http.Response, error) {
            return client.Send(req, azure.DoRetryWithRegistration(client.Client))
            }

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ActivityLogsClient) ListResponder(resp *http.Response) (result EventDataCollection, err error) {
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
            func (client ActivityLogsClient) listNextResults(ctx context.Context, lastResults EventDataCollection) (result EventDataCollection, err error) {
            req, err := lastResults.eventDataCollectionPreparer(ctx)
            if err != nil {
            return result, autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "listNextResults", nil , "Failure preparing next results request")
            }
            if req == nil {
            return
            }
            resp, err := client.ListSender(req)
            if err != nil {
            result.Response = autorest.Response{Response: resp}
            return result, autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "listNextResults", resp, "Failure sending next results request")
            }
            result, err = client.ListResponder(resp)
            if err != nil {
            err = autorest.NewErrorWithError(err, "insights.ActivityLogsClient", "listNextResults", resp, "Failure responding to next results request")
            }
            return
                    }

    // ListComplete enumerates all values, automatically crossing page boundaries as required.
    func (client ActivityLogsClient) ListComplete(ctx context.Context, filter string, selectParameter string) (result EventDataCollectionIterator, err error) {
        if tracing.IsEnabled() {
            ctx = tracing.StartSpan(ctx, fqdn + "/ActivityLogsClient.List")
            defer func() {
                sc := -1
                if result.Response().Response.Response != nil {
                    sc = result.page.Response().Response.Response.StatusCode
                }
                tracing.EndSpan(ctx, sc, err)
            }()
     }
        result.page, err = client.List(ctx, filter, selectParameter)
                return
        }

