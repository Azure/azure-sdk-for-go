// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v9"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// JobDefinitionJobRunC2CScenarioSuite mirrors .NET JobDefinitionJobRunTests with the additional
// matrix row #31 (StartC2CJobWithPrivateSourceTest, promoted to source-of-truth 2026-05-20). This
// is a full private-bucket end-to-end test that exercises the Storage Mover Connection / Private
// Endpoint approval flow + target endpoint MSI + cross-sub RBAC + C2C job execution.
//
// Lives in a separate suite (vs the existing JobDefinitionJobRunScenarioSuite for matrix row #10)
// because all infrastructure must be provisioned in `westcentralus` (shared `cpmoveraccount` and
// `test-pls-wcs` live there), whereas row #10 only needs placeholder endpoints in the default
// location.
//
// See `storage-mover-scenario-tests-cross-language.md`'s "Porter's reference" callout for the
// 12-step flow and shared infrastructure IDs.
type JobDefinitionJobRunC2CScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string
	projectName      string
	connectionName   string
	sourceEndpoint   string
	targetEndpoint   string
	jobDefName       string
	containerName    string
}

func TestJobDefinitionJobRunC2CScenarioSuite(t *testing.T) {
	suite.Run(t, new(JobDefinitionJobRunC2CScenarioSuite))
}

func (s *JobDefinitionJobRunC2CScenarioSuite) SetupSuite() {
	s.setupBaseInLocation(scenarioWestCentralUSLocation)
	s.storageMoverName = s.generateName("stomover")
	s.projectName = s.generateName("project")
	// Connection name capped at 20 chars by the RP (JS port finding 2026-05-26).
	s.connectionName = s.generateName("conn")
	s.sourceEndpoint = s.generateName("mccsrc")
	s.targetEndpoint = s.generateName("blobtgt")
	s.jobDefName = s.generateName("jobdef")
	// Storage container names must be lowercase. Use `tc` + 8-char lowercase suffix to keep names
	// per-test (the .NET port surfaced cross-test container collision on a shared name).
	s.containerName = s.generateName("tc")
	s.createStorageMover(s.storageMoverName, nil, "")
	s.createProject(s.storageMoverName, s.projectName, "")
}

// TearDownSuite cleans up the resource group and stops the recording.
func (s *JobDefinitionJobRunC2CScenarioSuite) TearDownSuite() { s.teardownBase() }

// Fixed GUID for the Storage Blob Data Contributor role assignment. The role-assignment name must
// be a valid GUID; we use a stable constant rather than `uuid.New()` to keep recordings
// deterministic without adding a non-approved direct dependency on github.com/google/uuid.
// Idempotency: if a previous run did not clean up, the second create returns RoleAssignmentExists
// which we treat as success.
const scenarioRoleAssignmentName = "8a3b0f70-1c44-4f8e-a3d8-7b9e5d8a2f10"

// TestStartC2CJobWithPrivateSource mirrors .NET
// JobDefinitionJobRunTests.StartC2CJobWithPrivateSourceTest (matrix row #31). Full 12-step flow:
//  1. RG/mover/project already provisioned in SetupSuite (westcentralus).
//  2. Create per-test target storage container on the shared cross-sub `cpmoveraccount`.
//  3. Create Storage Mover Connection against the shared PrivateLinkService `test-pls-wcs`.
//  4. Look up the corresponding PrivateEndpointConnection on the PLS (cross-sub, with backoff —
//     the PE-connection name is server-generated and takes a few seconds to appear).
//  5. Approve the PE-connection (PATCH connectionState = Approved).
//  6. Poll the Storage Mover Connection until ConnectionStatusApproved (≤5 min lag).
//  7. Create target Blob endpoint with SystemAssigned MSI; capture principal-id.
//  8. Grant Storage Blob Data Contributor on the container scope to the MSI (retry on
//     PrincipalNotFound while AAD propagates).
//  9. Create MCC source endpoint pointing at the private S3 bucket.
//  10. Create a CloudToCloud JobDefinition with all required fields plus Connections.
//  11. StartJob → returns JobRunResourceID; extract the job-run basename.
//  12. Poll JobRun.Get on 30s cadence (capped at 30 min) until terminal; assert
//     JobRunStatusSucceeded.
//
// Cleanup happens in `defer`s in reverse order, each tolerating already-gone resources.
func (s *JobDefinitionJobRunC2CScenarioSuite) TestStartC2CJobWithPrivateSource() {
	// --- Step 2: per-test target container on the shared `cpmoveraccount`. ---
	storageClient, err := armstorage.NewBlobContainersClient(scenarioCrossSubID, s.cred, s.options)
	s.Require().NoError(err)
	_, err = storageClient.Create(s.ctx, scenarioTestStorageAccountRG, scenarioTestStorageAccountName, s.containerName, armstorage.BlobContainer{}, nil)
	s.Require().NoError(err)
	defer func() {
		_, _ = storageClient.Delete(s.ctx, scenarioTestStorageAccountRG, scenarioTestStorageAccountName, s.containerName, nil)
	}()

	// --- Step 3: Storage Mover Connection against the shared PLS. ---
	connClient, err := armstoragemover.NewConnectionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	connResp, err := connClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, armstoragemover.Connection{
		Properties: &armstoragemover.ConnectionProperties{
			PrivateLinkServiceID: to.Ptr(scenarioPrivateLinkServiceID),
			Description:          to.Ptr("Storage Mover Connection for C2C private-bucket test"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Require().NotNil(connResp.Properties)
	s.Require().NotNil(connResp.Properties.PrivateEndpointResourceID, "RP should return PrivateEndpointResourceID on create")
	peResourceID := *connResp.Properties.PrivateEndpointResourceID
	defer func() {
		poller, perr := connClient.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
		if perr == nil {
			_, _ = testutil.PollForTest(s.ctx, poller)
		}
		// 60s grace after Connection delete to let the PLS slot fully release before the next test
		// in any future co-run touches the same PLS. .NET port found cross-variant contamination
		// without this (finding (d), 2026-05-25). No-op in playback.
		recording.Sleep(60 * time.Second)
	}()

	// --- Step 4: find the PE-connection on the PLS. Retry up to 10× × 15s. ---
	plsClient, err := armnetwork.NewPrivateLinkServicesClient(scenarioCrossSubID, s.cred, s.options)
	s.Require().NoError(err)
	peConnectionName := s.findPrivateEndpointConnection(plsClient, peResourceID)
	s.Require().NotEmpty(peConnectionName, "private endpoint connection for %q did not appear within 150s", peResourceID)

	// --- Step 5: approve the PE-connection. ---
	_, err = plsClient.UpdatePrivateEndpointConnection(s.ctx, scenarioPrivateLinkServiceRG, scenarioPrivateLinkServiceName, peConnectionName, armnetwork.PrivateEndpointConnection{
		Name: to.Ptr(peConnectionName),
		Properties: &armnetwork.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armnetwork.PrivateLinkServiceConnectionState{
				Status:          to.Ptr("Approved"),
				Description:     to.Ptr("Approved by storagemover scenario test (row #31)"),
				ActionsRequired: to.Ptr("None"),
			},
		},
	}, nil)
	s.Require().NoError(err)

	// --- Step 6: poll Connection until Approved. ≤5 min lag (10 × 30s). ---
	s.waitForConnectionApproved(connClient)

	// --- Step 7: target Blob endpoint with SystemAssigned MSI. ---
	endpointsClient, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	targetResp, err := endpointsClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.targetEndpoint, armstoragemover.Endpoint{
		Identity: &armstoragemover.ManagedServiceIdentity{
			Type: to.Ptr(armstoragemover.ManagedServiceIdentityTypeSystemAssigned),
		},
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			BlobContainerName:        to.Ptr(s.containerName),
			StorageAccountResourceID: to.Ptr(scenarioTestStorageAccountID),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindTarget),
			Description:              to.Ptr("C2C private-bucket target endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Require().NotNil(targetResp.Identity, "target endpoint should have SystemAssigned identity")
	s.Require().NotNil(targetResp.Identity.PrincipalID, "SystemAssigned MSI principalID missing")
	principalID := *targetResp.Identity.PrincipalID
	defer func() {
		poller, perr := endpointsClient.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.targetEndpoint, nil)
		if perr == nil {
			_, _ = testutil.PollForTest(s.ctx, poller)
		}
	}()

	// --- Step 8: assign Storage Blob Data Contributor on container scope to the MSI. ---
	scope := scenarioTestStorageAccountID + "/blobServices/default/containers/" + s.containerName
	roleDefinitionID := "/subscriptions/" + scenarioCrossSubID + "/providers/Microsoft.Authorization/roleDefinitions/" + scenarioStorageBlobDataContributorRoleID
	raClient, err := armauthorization.NewRoleAssignmentsClient(scenarioCrossSubID, s.cred, s.options)
	s.Require().NoError(err)
	s.createRoleAssignmentWithRetry(raClient, scope, principalID, roleDefinitionID)
	defer func() {
		_, _ = raClient.Delete(s.ctx, scope, scenarioRoleAssignmentName, nil)
	}()

	// --- Step 9: MCC source endpoint against private bucket. ---
	_, err = endpointsClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.sourceEndpoint, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureMultiCloudConnectorEndpointProperties{
			MultiCloudConnectorID: to.Ptr(scenarioMultiCloudConnectorID),
			AwsS3BucketID:         to.Ptr(scenarioAwsPrivateS3BucketID),
			EndpointKind:          to.Ptr(armstoragemover.EndpointKindSource),
			Description:           to.Ptr("C2C private-bucket source endpoint (MCC)"),
		},
	}, nil)
	s.Require().NoError(err)
	defer func() {
		poller, perr := endpointsClient.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.sourceEndpoint, nil)
		if perr == nil {
			_, _ = testutil.PollForTest(s.ctx, poller)
		}
	}()

	// --- Step 10: CloudToCloud JobDefinition with Connections + all required fields. ---
	// .NET port finding (b): JobDefinition needs all 4 of Description/JobType/SourceSubpath/
	// TargetSubpath explicitly set in addition to the required CopyMode/SourceName/TargetName.
	jobDefClient, err := armstoragemover.NewJobDefinitionsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	connectionResourceID := *connResp.ID
	_, err = jobDefClient.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, armstoragemover.JobDefinition{
		Properties: &armstoragemover.JobDefinitionProperties{
			CopyMode:      to.Ptr(armstoragemover.CopyModeAdditive),
			SourceName:    to.Ptr(s.sourceEndpoint),
			TargetName:    to.Ptr(s.targetEndpoint),
			Connections:   []*string{to.Ptr(connectionResourceID)},
			Description:   to.Ptr("C2C private-bucket job definition (row #31)"),
			JobType:       to.Ptr(armstoragemover.JobTypeCloudToCloud),
			SourceSubpath: to.Ptr("/"),
			TargetSubpath: to.Ptr("/"),
		},
	}, nil)
	s.Require().NoError(err)
	defer func() {
		poller, perr := jobDefClient.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
		if perr == nil {
			_, _ = testutil.PollForTest(s.ctx, poller)
		}
	}()

	// --- Step 11: StartJob → JobRunResourceID basename. ---
	startResp, err := jobDefClient.StartJob(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, nil)
	s.Require().NoError(err)
	s.Require().NotNil(startResp.JobRunResourceID.JobRunResourceID)
	jobRunName := path.Base(*startResp.JobRunResourceID.JobRunResourceID)
	s.Require().NotEmpty(jobRunName)

	// --- Step 12: poll JobRun until terminal; assert Succeeded. ---
	s.waitForJobRunSucceeded(jobRunName)
}

// findPrivateEndpointConnection lists PE-connections on the shared PLS until it finds one whose
// PrivateEndpoint.ID matches peResourceID. Retries up to 10× with 15s backoff (~150s total) to
// accommodate the few-second server-side delay before the PE-connection appears. Returns the
// server-generated PE-connection name, or empty if not found in time.
func (s *JobDefinitionJobRunC2CScenarioSuite) findPrivateEndpointConnection(client *armnetwork.PrivateLinkServicesClient, peResourceID string) string {
	const maxAttempts = 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		pager := client.NewListPrivateEndpointConnectionsPager(scenarioPrivateLinkServiceRG, scenarioPrivateLinkServiceName, nil)
		for pager.More() {
			page, err := pager.NextPage(s.ctx)
			s.Require().NoError(err)
			for _, pec := range page.Value {
				if pec.Properties != nil && pec.Properties.PrivateEndpoint != nil && pec.Properties.PrivateEndpoint.ID != nil &&
					strings.EqualFold(*pec.Properties.PrivateEndpoint.ID, peResourceID) && pec.Name != nil {
					return *pec.Name
				}
			}
		}
		recording.Sleep(15 * time.Second)
	}
	return ""
}

// waitForConnectionApproved polls the Storage Mover Connection's Get endpoint until
// ConnectionStatus == Approved. Up to 10 × 30s = 5 min.
func (s *JobDefinitionJobRunC2CScenarioSuite) waitForConnectionApproved(client *armstoragemover.ConnectionsClient) {
	const maxAttempts = 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.connectionName, nil)
		s.Require().NoError(err)
		if resp.Properties != nil && resp.Properties.ConnectionStatus != nil && *resp.Properties.ConnectionStatus == armstoragemover.ConnectionStatusApproved {
			return
		}
		recording.Sleep(30 * time.Second)
	}
	s.FailNow("Storage Mover Connection did not reach Approved status within 5 minutes")
}

// createRoleAssignmentWithRetry creates the Storage Blob Data Contributor role assignment on the
// target container scope for the freshly-minted MSI. Retries up to 10× on PrincipalNotFound while
// AAD propagates the new service principal. Idempotent on RoleAssignmentExists (treat as success
// when re-running against a not-cleaned-up environment).
func (s *JobDefinitionJobRunC2CScenarioSuite) createRoleAssignmentWithRetry(client *armauthorization.RoleAssignmentsClient, scope, principalID, roleDefinitionID string) {
	const maxAttempts = 10
	params := armauthorization.RoleAssignmentCreateParameters{
		Properties: &armauthorization.RoleAssignmentProperties{
			PrincipalID:      to.Ptr(principalID),
			RoleDefinitionID: to.Ptr(roleDefinitionID),
			PrincipalType:    to.Ptr(armauthorization.PrincipalTypeServicePrincipal),
		},
	}
	for attempt := 0; attempt < maxAttempts; attempt++ {
		_, err := client.Create(s.ctx, scope, scenarioRoleAssignmentName, params, nil)
		if err == nil {
			return
		}
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			if respErr.ErrorCode == "RoleAssignmentExists" {
				return
			}
			if respErr.ErrorCode == "PrincipalNotFound" {
				recording.Sleep(15 * time.Second)
				continue
			}
		}
		s.Require().NoError(err, "role assignment create failed (attempt %d)", attempt+1)
	}
	s.FailNow("Storage Blob Data Contributor role assignment did not succeed within retry budget")
}

// waitForJobRunSucceeded polls the JobRun on 30s cadence, capped at 30 minutes, until the run
// reaches a terminal status. Asserts Succeeded.
func (s *JobDefinitionJobRunC2CScenarioSuite) waitForJobRunSucceeded(jobRunName string) {
	const maxAttempts = 60
	jobRunsClient, err := armstoragemover.NewJobRunsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	var lastStatus armstoragemover.JobRunStatus
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := jobRunsClient.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.projectName, s.jobDefName, jobRunName, nil)
		s.Require().NoError(err)
		if resp.Properties != nil && resp.Properties.Status != nil {
			lastStatus = *resp.Properties.Status
			if isTerminalJobRunStatus(lastStatus) {
				s.Equal(armstoragemover.JobRunStatusSucceeded, lastStatus, "JobRun finished in non-Succeeded terminal status %q", lastStatus)
				return
			}
		}
		recording.Sleep(30 * time.Second)
	}
	s.FailNow(fmt.Sprintf("JobRun did not reach a terminal status within 30 minutes (last observed status: %q)", lastStatus))
}

func isTerminalJobRunStatus(status armstoragemover.JobRunStatus) bool {
	switch status {
	case armstoragemover.JobRunStatusSucceeded,
		armstoragemover.JobRunStatusFailed,
		armstoragemover.JobRunStatusCanceled:
		return true
	}
	return false
}
