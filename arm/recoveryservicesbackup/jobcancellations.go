package recoveryservicesbackup

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
// Code generated by Microsoft (R) AutoRest Code Generator 2.2.21.0
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// JobCancellationsClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type JobCancellationsClient struct {
	ManagementClient
}

// NewJobCancellationsClient creates an instance of the JobCancellationsClient client.
func NewJobCancellationsClient(subscriptionID string) JobCancellationsClient {
	return NewJobCancellationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewJobCancellationsClientWithBaseURI creates an instance of the JobCancellationsClient client.
func NewJobCancellationsClientWithBaseURI(baseURI string, subscriptionID string) JobCancellationsClient {
	return JobCancellationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Trigger cancels a job. This is an asynchronous operation. To know the status of the cancellation, call
// GetCancelOperationResult API.
//
// vaultName is the name of the recovery services vault. resourceGroupName is the name of the resource group where the
// recovery services vault is present. jobName is name of the job to cancel.
func (client JobCancellationsClient) Trigger(vaultName string, resourceGroupName string, jobName string) (result autorest.Response, err error) {
	req, err := client.TriggerPreparer(vaultName, resourceGroupName, jobName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobCancellationsClient", "Trigger", nil, "Failure preparing request")
		return
	}

	resp, err := client.TriggerSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobCancellationsClient", "Trigger", resp, "Failure sending request")
		return
	}

	result, err = client.TriggerResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.JobCancellationsClient", "Trigger", resp, "Failure responding to request")
	}

	return
}

// TriggerPreparer prepares the Trigger request.
func (client JobCancellationsClient) TriggerPreparer(vaultName string, resourceGroupName string, jobName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"jobName":           autorest.Encode("path", jobName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupJobs/{jobName}/cancel", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// TriggerSender sends the Trigger request. The method will close the
// http.Response Body if it receives an error.
func (client JobCancellationsClient) TriggerSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// TriggerResponder handles the response to the Trigger request. The method always
// closes the http.Response Body.
func (client JobCancellationsClient) TriggerResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}
