// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// Constants mirrored from the .NET StorageMoverManagementTestBase. These are placeholder values that
// the resource provider accepts at the metadata level for endpoint/jobdefinition CRUD.
const (
	scenarioStorageAccountName    = "testsmstore24"
	scenarioContainerName         = "testsmcontainer"
	scenarioFileShareName         = "testfileshare"
	scenarioNfsFileShareName      = "testnfsfileshare"
	scenarioSmbHost               = "10.0.0.1"
	scenarioSmbShareName          = "testshare"
	scenarioNfsExport             = "/"
	scenarioMultiCloudConnectorID = "/subscriptions/b6b34ad8-ca89-4f85-beb7-c2ec13702dac/resourceGroups/E2E-Management-RGsyn/providers/Microsoft.HybridConnectivity/publicCloudConnectors/e2e-sm-rp-connector"
	scenarioAwsS3BucketID         = "/subscriptions/b6b34ad8-ca89-4f85-beb7-c2ec13702dac/resourceGroups/aws_640698235822/providers/Microsoft.AWSConnector/s3Buckets/e2e-sm-rp-bucket"
	scenarioKeyVaultUsernameURI   = "https://examples-azureKeyVault.vault.azure.net/secrets/examples-username"
	scenarioKeyVaultPasswordURI   = "https://examples-azureKeyVault.vault.azure.net/secrets/examples-password"
	scenarioKeyVaultAccessKeyURI  = "https://examples-azureKeyVault.vault.azure.net/secrets/examples-accesskey"
	scenarioKeyVaultSecretKeyURI  = "https://examples-azureKeyVault.vault.azure.net/secrets/examples-secretkey"
	scenarioS3SourceURI           = "https://s3.example.com/bucket"

	// Cross-subscription shared infrastructure for matrix rows #31 and #32. All of these live in
	// the XDataMove-Synthetics subscription (b6b34ad8-…) and must not be recreated.
	scenarioCrossSubID                       = "b6b34ad8-ca89-4f85-beb7-c2ec13702dac"
	scenarioPrivateLinkServiceID             = "/subscriptions/b6b34ad8-ca89-4f85-beb7-c2ec13702dac/resourceGroups/E2E-Management-RGsyn/providers/Microsoft.Network/privateLinkServices/test-pls-wcs"
	scenarioPrivateLinkServiceRG             = "E2E-Management-RGsyn"
	scenarioPrivateLinkServiceName           = "test-pls-wcs"
	scenarioAwsPrivateS3BucketID             = "/subscriptions/b6b34ad8-ca89-4f85-beb7-c2ec13702dac/resourceGroups/aws_640698235822/providers/Microsoft.AWSConnector/s3Buckets/e2e-sm-rp-private-bucket"
	scenarioTestStorageAccountID             = "/subscriptions/b6b34ad8-ca89-4f85-beb7-c2ec13702dac/resourceGroups/CP_Mover_IN_WCUS/providers/Microsoft.Storage/storageAccounts/cpmoveraccount"
	scenarioTestStorageAccountRG             = "CP_Mover_IN_WCUS"
	scenarioTestStorageAccountName           = "cpmoveraccount"
	scenarioWestCentralUSLocation            = "westcentralus"
	scenarioStorageBlobDataContributorRoleID = "ba92f5b4-2d11-453d-a403-e96b0029c9fe"
)

// scenarioBaseSuite holds the shared resource group, credentials, and client options used by every
// scenario test suite. Each suite creates its own resource group; tests within a suite create their
// own mutable child resources to keep tests order-independent.
type scenarioBaseSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

// setupBase initializes recording, credentials, and a fresh resource group for the suite in the
// default location (`LOCATION` env, or `eastus`). Call from SetupSuite.
func (b *scenarioBaseSuite) setupBase() {
	b.setupBaseInLocation(recording.GetEnvVariable("LOCATION", "eastus"))
}

// setupBaseInLocation is the location-explicit form of setupBase used by suites that must run in a
// specific region (e.g. matrix rows #31 and #32 need westcentralus because the shared
// cross-subscription storage account `cpmoveraccount` and PrivateLinkService `test-pls-wcs` live
// in WCUS).
func (b *scenarioBaseSuite) setupBaseInLocation(location string) {
	testutil.StartRecording(b.T(), pathToPackage)

	// AZSDK3493 is the default Body Key Sanitizer that replaces "$..name" with "Sanitized" in every
	// response body. Storage Mover scenario tests assert on resource names returned from
	// CreateOrUpdate/Get, so we disable that sanitizer for this module's recordings. Resource names
	// are already exposed in the URL path, so redacting them in the body provides no extra
	// protection. We also disable AZSDK3430 ($..id) because we assert get-equivalence on .ID.
	//
	// IMPORTANT: TestInstance must be set so the unregister applies to the already-opened recording
	// session — without it, the unregister only modifies the global pool, leaving this session's
	// already-snapshotted sanitizer set untouched.
	if err := recording.RemoveRegisteredSanitizers([]string{"AZSDK3493", "AZSDK3430"}, &recording.RecordingOptions{UseHTTPS: true, TestInstance: b.T()}); err != nil {
		b.T().Logf("warning: failed to remove name/id sanitizers: %v", err)
	}

	b.ctx = context.Background()
	b.cred, b.options = testutil.GetCredAndClientOptions(b.T())
	b.location = location
	b.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	rg, _, err := testutil.CreateResourceGroup(b.ctx, b.subscriptionID, b.cred, b.options, b.location)
	b.Require().NoError(err)
	b.resourceGroupName = *rg.Name
}

// teardownBase cleans up the resource group and stops recording. Call from TearDownSuite.
func (b *scenarioBaseSuite) teardownBase() {
	if b.resourceGroupName != "" {
		_, err := testutil.DeleteResourceGroup(b.ctx, b.subscriptionID, b.cred, b.options, b.resourceGroupName)
		b.Require().NoError(err)
	}
	testutil.StopRecording(b.T())
}

// generateName returns a recording-stable random alphanumeric identifier with the given prefix.
func (b *scenarioBaseSuite) generateName(prefix string) string {
	name, err := recording.GenerateAlphaNumericID(b.T(), prefix, len(prefix)+8, true)
	b.Require().NoError(err)
	return name
}

// createStorageMover creates a Storage Mover with the given name in the suite's resource group and
// returns the created resource. Tests that share a single mover should call this once from SetupSuite;
// tests that mutate or delete the mover should call this from the test body with a unique name.
func (b *scenarioBaseSuite) createStorageMover(name string, tags map[string]*string, description string) *armstoragemover.StorageMover {
	client, err := armstoragemover.NewStorageMoversClient(b.subscriptionID, b.cred, b.options)
	b.Require().NoError(err)
	resp, err := client.CreateOrUpdate(b.ctx, b.resourceGroupName, name, armstoragemover.StorageMover{
		Location: to.Ptr(b.location),
		Tags:     tags,
		Properties: &armstoragemover.Properties{
			Description: to.Ptr(description),
		},
	}, nil)
	b.Require().NoError(err)
	return &resp.StorageMover
}

// createProject creates a Project under the given Storage Mover. An empty description is sent as nil
// so the resulting resource has an unset description (matching the .NET tests that pass an empty
// `StorageMoverProjectData()`).
func (b *scenarioBaseSuite) createProject(storageMoverName, projectName, description string) *armstoragemover.Project {
	client, err := armstoragemover.NewProjectsClient(b.subscriptionID, b.cred, b.options)
	b.Require().NoError(err)
	props := &armstoragemover.ProjectProperties{}
	if description != "" {
		props.Description = to.Ptr(description)
	}
	resp, err := client.CreateOrUpdate(b.ctx, b.resourceGroupName, storageMoverName, projectName, armstoragemover.Project{
		Properties: props,
	}, nil)
	b.Require().NoError(err)
	return &resp.Project
}

// createBlobEndpoint creates an AzureStorageBlobContainer endpoint with placeholder storage account
// resource ID. Returns the created endpoint.
func (b *scenarioBaseSuite) createBlobEndpoint(storageMoverName, endpointName, description string) *armstoragemover.Endpoint {
	client, err := armstoragemover.NewEndpointsClient(b.subscriptionID, b.cred, b.options)
	b.Require().NoError(err)
	resp, err := client.CreateOrUpdate(b.ctx, b.resourceGroupName, storageMoverName, endpointName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			BlobContainerName:        to.Ptr(scenarioContainerName),
			StorageAccountResourceID: to.Ptr(b.placeholderStorageAccountResourceID()),
			Description:              to.Ptr(description),
		},
	}, nil)
	b.Require().NoError(err)
	return &resp.Endpoint
}

// createNfsEndpoint creates an NfsMount endpoint with placeholder host/export.
func (b *scenarioBaseSuite) createNfsEndpoint(storageMoverName, endpointName, description string) *armstoragemover.Endpoint {
	client, err := armstoragemover.NewEndpointsClient(b.subscriptionID, b.cred, b.options)
	b.Require().NoError(err)
	resp, err := client.CreateOrUpdate(b.ctx, b.resourceGroupName, storageMoverName, endpointName, armstoragemover.Endpoint{
		Properties: &armstoragemover.NfsMountEndpointProperties{
			Host:        to.Ptr(scenarioSmbHost),
			Export:      to.Ptr(scenarioNfsExport),
			Description: to.Ptr(description),
		},
	}, nil)
	b.Require().NoError(err)
	return &resp.Endpoint
}

// placeholderStorageAccountResourceID returns a synthetic ARM resource ID. The Storage Mover RP only
// validates the format at metadata level; a non-existent storage account is accepted for endpoint
// CRUD scenarios that do not run an actual job.
func (b *scenarioBaseSuite) placeholderStorageAccountResourceID() string {
	return "/subscriptions/" + b.subscriptionID + "/resourceGroups/" + b.resourceGroupName +
		"/providers/Microsoft.Storage/storageAccounts/" + scenarioStorageAccountName
}

// expectResponseError unwraps err as *azcore.ResponseError and asserts the HTTP status is non-2xx.
// We deliberately do not lock to specific RP error codes since they drift across api-versions.
func (b *scenarioBaseSuite) expectResponseError(err error) {
	b.Require().Error(err)
	var respErr *azcore.ResponseError
	b.Require().True(errors.As(err, &respErr), "expected azcore.ResponseError, got %T: %v", err, err)
	b.Require().GreaterOrEqual(respErr.StatusCode, http.StatusBadRequest, "expected non-2xx status, got %d", respErr.StatusCode)
}
