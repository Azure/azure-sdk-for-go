package attestation

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
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// Client is the describes the interface for the per-tenant enclave service.
type Client struct {
	BaseClient
}

// NewClient creates an instance of the Client client.
func NewClient() Client {
	return Client{New()}
}

// AttestOpenEnclave processes an OpenEnclave report , producing an artifact. The type of artifact produced is
// dependent upon attestation policy.
// Parameters:
// instanceURL - the attestation instance base URI, for example https://mytenant.attest.azure.net.
// request - request object containing the quote
func (client Client) AttestOpenEnclave(ctx context.Context, instanceURL string, request AttestOpenEnclaveRequest) (result Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/Client.AttestOpenEnclave")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AttestOpenEnclavePreparer(ctx, instanceURL, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestOpenEnclave", nil, "Failure preparing request")
		return
	}

	resp, err := client.AttestOpenEnclaveSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestOpenEnclave", resp, "Failure sending request")
		return
	}

	result, err = client.AttestOpenEnclaveResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestOpenEnclave", resp, "Failure responding to request")
	}

	return
}

// AttestOpenEnclavePreparer prepares the AttestOpenEnclave request.
func (client Client) AttestOpenEnclavePreparer(ctx context.Context, instanceURL string, request AttestOpenEnclaveRequest) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"instanceUrl": instanceURL,
	}

	const APIVersion = "2020-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{instanceUrl}", urlParameters),
		autorest.WithPath("/attest/OpenEnclave"),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AttestOpenEnclaveSender sends the AttestOpenEnclave request. The method will close the
// http.Response Body if it receives an error.
func (client Client) AttestOpenEnclaveSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// AttestOpenEnclaveResponder handles the response to the AttestOpenEnclave request. The method always
// closes the http.Response Body.
func (client Client) AttestOpenEnclaveResponder(resp *http.Response) (result Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// AttestSgxEnclave processes an SGX enclave quote, producing an artifact. The type of artifact produced is dependent
// upon attestation policy.
// Parameters:
// instanceURL - the attestation instance base URI, for example https://mytenant.attest.azure.net.
// request - request object containing the quote
func (client Client) AttestSgxEnclave(ctx context.Context, instanceURL string, request AttestSgxEnclaveRequest) (result Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/Client.AttestSgxEnclave")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AttestSgxEnclavePreparer(ctx, instanceURL, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestSgxEnclave", nil, "Failure preparing request")
		return
	}

	resp, err := client.AttestSgxEnclaveSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestSgxEnclave", resp, "Failure sending request")
		return
	}

	result, err = client.AttestSgxEnclaveResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestSgxEnclave", resp, "Failure responding to request")
	}

	return
}

// AttestSgxEnclavePreparer prepares the AttestSgxEnclave request.
func (client Client) AttestSgxEnclavePreparer(ctx context.Context, instanceURL string, request AttestSgxEnclaveRequest) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"instanceUrl": instanceURL,
	}

	const APIVersion = "2020-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{instanceUrl}", urlParameters),
		autorest.WithPath("/attest/SgxEnclave"),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AttestSgxEnclaveSender sends the AttestSgxEnclave request. The method will close the
// http.Response Body if it receives an error.
func (client Client) AttestSgxEnclaveSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// AttestSgxEnclaveResponder handles the response to the AttestSgxEnclave request. The method always
// closes the http.Response Body.
func (client Client) AttestSgxEnclaveResponder(resp *http.Response) (result Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// AttestTpm processes attestation evidence from a VBS enclave, producing an attestation result. The attestation result
// produced is dependent upon the attestation policy.
// Parameters:
// instanceURL - the attestation instance base URI, for example https://mytenant.attest.azure.net.
// request - request object
func (client Client) AttestTpm(ctx context.Context, instanceURL string, request TpmAttestationRequest) (result TpmAttestationResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/Client.AttestTpm")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.AttestTpmPreparer(ctx, instanceURL, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestTpm", nil, "Failure preparing request")
		return
	}

	resp, err := client.AttestTpmSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestTpm", resp, "Failure sending request")
		return
	}

	result, err = client.AttestTpmResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "attestation.Client", "AttestTpm", resp, "Failure responding to request")
	}

	return
}

// AttestTpmPreparer prepares the AttestTpm request.
func (client Client) AttestTpmPreparer(ctx context.Context, instanceURL string, request TpmAttestationRequest) (*http.Request, error) {
	urlParameters := map[string]interface{}{
		"instanceUrl": instanceURL,
	}

	const APIVersion = "2020-10-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithCustomBaseURL("{instanceUrl}", urlParameters),
		autorest.WithPath("/attest/Tpm"),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AttestTpmSender sends the AttestTpm request. The method will close the
// http.Response Body if it receives an error.
func (client Client) AttestTpmSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}

// AttestTpmResponder handles the response to the AttestTpm request. The method always
// closes the http.Response Body.
func (client Client) AttestTpmResponder(resp *http.Response) (result TpmAttestationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
