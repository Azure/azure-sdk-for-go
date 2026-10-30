// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/require"
)

const (
	contractSubscriptionID    = "00000000-0000-0000-0000-000000000000"
	contractMoverID           = "/subscriptions/" + contractSubscriptionID + "/resourceGroups/resource-group/providers/Microsoft.StorageMover/storageMovers/mover"
	contractStorageAccountID  = "/subscriptions/" + contractSubscriptionID + "/resourceGroups/resource-group/providers/Microsoft.Storage/storageAccounts/account"
	contractPartnerAccountID  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/partner/providers/Microsoft.Storage/storageAccounts/partneraccount"
	contractPartnerEndpointID = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/partner/providers/Microsoft.StorageMover/storageMovers/partner/endpoints/target"
	contractPartnerTenantID   = "22222222-2222-2222-2222-222222222222"
	contractConnectionID      = contractMoverID + "/connections/connection"
)

type apiContractCredential struct{}

func (apiContractCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

type apiContractTransport func(*http.Request) (*http.Response, error)

func (f apiContractTransport) Do(request *http.Request) (*http.Response, error) {
	return f(request)
}

func newAPIContractOptions(t *testing.T, method, path, requestBody, responseBody string) *arm.ClientOptions {
	t.Helper()
	requests := 0
	t.Cleanup(func() { require.Equal(t, 1, requests, "expected exactly one SDK request") })
	return &arm.ClientOptions{ClientOptions: policy.ClientOptions{
		Retry: policy.RetryOptions{MaxRetries: -1},
		Transport: apiContractTransport(func(request *http.Request) (*http.Response, error) {
			requests++
			require.Equal(t, method, request.Method)
			require.Equal(t, contractMoverID+"/"+path, request.URL.Path)
			require.Equal(t, "2026-05-01", request.URL.Query().Get("api-version"))
			if requestBody == "" {
				require.True(t, request.Body == nil || request.Body == http.NoBody, "expected no request body")
			} else {
				require.NotNil(t, request.Body)
				body, err := io.ReadAll(request.Body)
				require.NoError(t, err)
				require.JSONEq(t, requestBody, string(body))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Request:    request,
			}, nil
		}),
	}}
}

func TestAPIContractEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		properties armstoragemover.EndpointBasePropertiesClassification
		wantJSON   string
	}{
		{
			name: "blob-cross-tenant",
			properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
				StorageAccountResourceID:  to.Ptr(contractStorageAccountID),
				BlobContainerName:         to.Ptr("container"),
				EndpointKind:              to.Ptr(armstoragemover.EndpointKindTarget),
				EnableCrossTenantTransfer: to.Ptr(true),
				AllowedStorageAccounts:    []*string{to.Ptr(contractPartnerAccountID)},
			},
			wantJSON: fmt.Sprintf(`{"properties":{
				"endpointType":"AzureStorageBlobContainer","storageAccountResourceId":%q,
				"blobContainerName":"container","endpointKind":"Target",
				"enableCrossTenantTransfer":true,"allowedStorageAccounts":[%q]
			}}`, contractStorageAccountID, contractPartnerAccountID),
		},
		{
			name: "smb-cross-tenant",
			properties: &armstoragemover.AzureStorageSmbFileShareEndpointProperties{
				StorageAccountResourceID:  to.Ptr(contractStorageAccountID),
				FileShareName:             to.Ptr("share"),
				EndpointKind:              to.Ptr(armstoragemover.EndpointKindTarget),
				EnableCrossTenantTransfer: to.Ptr(true),
				AllowedStorageAccounts:    []*string{to.Ptr(contractPartnerAccountID)},
			},
			wantJSON: fmt.Sprintf(`{"properties":{
				"endpointType":"AzureStorageSmbFileShare","storageAccountResourceId":%q,
				"fileShareName":"share","endpointKind":"Target",
				"enableCrossTenantTransfer":true,"allowedStorageAccounts":[%q]
			}}`, contractStorageAccountID, contractPartnerAccountID),
		},
		{
			name: "nfs-fsx",
			properties: &armstoragemover.NfsMountEndpointProperties{
				Host:         to.Ptr("10.0.0.5"),
				Export:       to.Ptr("/export"),
				EndpointKind: to.Ptr(armstoragemover.EndpointKindSource),
				NfsVersion:   to.Ptr(armstoragemover.NfsVersionNFSv41),
				SourceType:   to.Ptr(armstoragemover.NfsMountSourceTypeFSXEFS),
			},
			wantJSON: `{"properties":{
				"endpointType":"NfsMount","host":"10.0.0.5","export":"/export",
				"endpointKind":"Source","nfsVersion":"NFSv4_1","sourceType":"FSX-EFS"
			}}`,
		},
		{
			name: "smb-fsx",
			properties: &armstoragemover.SmbMountEndpointProperties{
				Host:         to.Ptr("10.0.0.6"),
				ShareName:    to.Ptr("share"),
				EndpointKind: to.Ptr(armstoragemover.EndpointKindSource),
				SourceType:   to.Ptr(armstoragemover.SmbMountSourceTypeFSXSMB),
				Credentials: &armstoragemover.AzureKeyVaultSmbCredentials{
					UsernameURI: to.Ptr("https://vault.vault.azure.net/secrets/username"),
					PasswordURI: to.Ptr("https://vault.vault.azure.net/secrets/password"),
				},
			},
			wantJSON: `{"properties":{
				"endpointType":"SmbMount","host":"10.0.0.6","shareName":"share",
				"endpointKind":"Source","sourceType":"FSX-SMB",
				"credentials":{"type":"AzureKeyVaultSmb",
					"usernameUri":"https://vault.vault.azure.net/secrets/username",
					"passwordUri":"https://vault.vault.azure.net/secrets/password"}
			}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := newAPIContractOptions(t, http.MethodPut, "endpoints/"+tt.name, tt.wantJSON, tt.wantJSON)
			client, err := armstoragemover.NewEndpointsClient(contractSubscriptionID, apiContractCredential{}, options)
			require.NoError(t, err)
			response, err := client.CreateOrUpdate(context.Background(), "resource-group", "mover", tt.name,
				armstoragemover.Endpoint{Properties: tt.properties}, nil)
			require.NoError(t, err)
			require.IsType(t, tt.properties, response.Properties)
			body, err := json.Marshal(response.Endpoint)
			require.NoError(t, err)
			require.JSONEq(t, tt.wantJSON, string(body))
		})
	}
}

func TestAPIContractEndpointUpdates(t *testing.T) {
	endpoints := []struct {
		name       string
		properties func(*bool, []*string, *string) armstoragemover.EndpointBaseUpdatePropertiesClassification
	}{
		{
			name: "AzureStorageBlobContainer",
			properties: func(enabled *bool, accounts []*string, description *string) armstoragemover.EndpointBaseUpdatePropertiesClassification {
				return &armstoragemover.AzureStorageBlobContainerEndpointUpdateProperties{
					EnableCrossTenantTransfer: enabled, AllowedStorageAccounts: accounts, Description: description,
				}
			},
		},
		{
			name: "AzureStorageSmbFileShare",
			properties: func(enabled *bool, accounts []*string, description *string) armstoragemover.EndpointBaseUpdatePropertiesClassification {
				return &armstoragemover.AzureStorageSmbFileShareEndpointUpdateProperties{
					EnableCrossTenantTransfer: enabled, AllowedStorageAccounts: accounts, Description: description,
				}
			},
		},
	}
	updates := []struct {
		name        string
		enabled     *bool
		accounts    []*string
		description *string
		wantFields  string
	}{
		{
			name: "replace-allow-list", accounts: []*string{to.Ptr(contractPartnerAccountID)},
			wantFields: fmt.Sprintf(`"allowedStorageAccounts":[%q]`, contractPartnerAccountID),
		},
		{
			name: "omit-cross-tenant-settings", description: to.Ptr("Updated endpoint"),
			wantFields: `"description":"Updated endpoint"`,
		},
		{
			name: "disable-and-clear-allow-list", enabled: to.Ptr(false), accounts: []*string{},
			wantFields: `"enableCrossTenantTransfer":false,"allowedStorageAccounts":[]`,
		},
	}
	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			for _, update := range updates {
				t.Run(update.name, func(t *testing.T) {
					wantJSON := fmt.Sprintf(`{"identity":{"type":"None"},"properties":{"endpointType":%q,%s}}`,
						endpoint.name, update.wantFields)
					options := newAPIContractOptions(t, http.MethodPatch, "endpoints/endpoint", wantJSON, wantJSON)
					client, err := armstoragemover.NewEndpointsClient(contractSubscriptionID, apiContractCredential{}, options)
					require.NoError(t, err)
					response, err := client.Update(context.Background(), "resource-group", "mover", "endpoint",
						armstoragemover.EndpointBaseUpdateParameters{
							Identity:   &armstoragemover.ManagedServiceIdentity{Type: to.Ptr(armstoragemover.ManagedServiceIdentityTypeNone)},
							Properties: endpoint.properties(update.enabled, update.accounts, update.description),
						}, nil)
					require.NoError(t, err)
					body, err := json.Marshal(response.Endpoint)
					require.NoError(t, err)
					require.JSONEq(t, wantJSON, string(body))
				})
			}
		})
	}
}

func TestAPIContractJobDefinitions(t *testing.T) {
	syncedUntil := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		properties *armstoragemover.JobDefinitionProperties
		wantJSON   string
	}{
		{
			name: "cross-tenant-job",
			properties: &armstoragemover.JobDefinitionProperties{
				CopyMode:                      to.Ptr(armstoragemover.CopyModeAdditive),
				SourceName:                    to.Ptr("source"),
				TargetName:                    to.Ptr("target"),
				JobType:                       to.Ptr(armstoragemover.JobTypeCloudToCloud),
				IsCrossTenantJob:              to.Ptr(true),
				CrossTenantEndpointTenantID:   to.Ptr(contractPartnerTenantID),
				CrossTenantEndpointResourceID: to.Ptr(contractPartnerEndpointID),
				SyncMode:                      to.Ptr("SnapshotBased"),
				MoverSyncedUntil:              &syncedUntil,
				Schedule: &armstoragemover.ScheduleInfo{
					Frequency: to.Ptr(armstoragemover.FrequencyHourly), RepeatInterval: to.Ptr("PT1H"), IsActive: to.Ptr(true),
				},
			},
			wantJSON: fmt.Sprintf(`{"properties":{
				"copyMode":"Additive","sourceName":"source","targetName":"target","jobType":"CloudToCloud",
				"isCrossTenantJob":true,"crossTenantEndpointTenantId":%q,"crossTenantEndpointResourceId":%q,
				"syncMode":"SnapshotBased","moverSyncedUntil":"2026-05-01T00:00:00Z",
				"schedule":{"frequency":"Hourly","repeatInterval":"PT1H","isActive":true}
			}}`, contractPartnerTenantID, contractPartnerEndpointID),
		},
		{
			name: "agentless-job",
			properties: &armstoragemover.JobDefinitionProperties{
				CopyMode: to.Ptr(armstoragemover.CopyModeAdditive), SourceName: to.Ptr("source"), TargetName: to.Ptr("target"),
				JobType: to.Ptr(armstoragemover.JobTypeOnPremToCloudAgentLess), Connections: []*string{to.Ptr(contractConnectionID)},
			},
			wantJSON: fmt.Sprintf(`{"properties":{
				"copyMode":"Additive","sourceName":"source","targetName":"target",
				"jobType":"OnPremToCloudAgentLess","connections":[%q]
			}}`, contractConnectionID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := newAPIContractOptions(t, http.MethodPut, "projects/project/jobDefinitions/"+tt.name, tt.wantJSON, tt.wantJSON)
			client, err := armstoragemover.NewJobDefinitionsClient(contractSubscriptionID, apiContractCredential{}, options)
			require.NoError(t, err)
			response, err := client.CreateOrUpdate(context.Background(), "resource-group", "mover", "project", tt.name,
				armstoragemover.JobDefinition{Properties: tt.properties}, nil)
			require.NoError(t, err)
			require.Equal(t, tt.properties, response.Properties)
		})
	}
}

func TestAPIContractJobDefinitionUpdate(t *testing.T) {
	wantJSON := fmt.Sprintf(`{"properties":{
		"connections":[%q],"syncMode":"FullScan","moverSyncedUntil":"2026-05-02T00:00:00Z",
		"schedule":{"frequency":"Hourly","repeatInterval":"PT2H","isActive":false}
	}}`, contractConnectionID)
	options := newAPIContractOptions(t, http.MethodPatch, "projects/project/jobDefinitions/job", wantJSON, wantJSON)
	client, err := armstoragemover.NewJobDefinitionsClient(contractSubscriptionID, apiContractCredential{}, options)
	require.NoError(t, err)
	syncedUntil := time.Date(2026, time.May, 2, 0, 0, 0, 0, time.UTC)
	response, err := client.Update(context.Background(), "resource-group", "mover", "project", "job",
		armstoragemover.JobDefinitionUpdateParameters{Properties: &armstoragemover.JobDefinitionUpdateProperties{
			Connections: []*string{to.Ptr(contractConnectionID)}, SyncMode: to.Ptr("FullScan"), MoverSyncedUntil: &syncedUntil,
			Schedule: &armstoragemover.ScheduleInfo{
				Frequency: to.Ptr(armstoragemover.FrequencyHourly), RepeatInterval: to.Ptr("PT2H"), IsActive: to.Ptr(false),
			},
		}}, nil)
	require.NoError(t, err)
	body, err := json.Marshal(response.JobDefinition)
	require.NoError(t, err)
	require.JSONEq(t, wantJSON, string(body))
}

func TestAPIContractReconcileJob(t *testing.T) {
	jobRunID := contractMoverID + "/projects/project/jobDefinitions/job/jobRuns/run"
	options := newAPIContractOptions(t, http.MethodPost, "projects/project/jobDefinitions/job/reconcileJob", "",
		fmt.Sprintf(`{"jobRunResourceId":%q}`, jobRunID))
	client, err := armstoragemover.NewJobDefinitionsClient(contractSubscriptionID, apiContractCredential{}, options)
	require.NoError(t, err)
	response, err := client.ReconcileJob(context.Background(), "resource-group", "mover", "project", "job", nil)
	require.NoError(t, err)
	require.Equal(t, to.Ptr(jobRunID), response.JobRunResourceID.JobRunResourceID)
}
