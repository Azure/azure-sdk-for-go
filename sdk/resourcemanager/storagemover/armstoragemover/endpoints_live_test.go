// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstoragemover_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
	"github.com/stretchr/testify/suite"
)

// EndpointScenarioSuite mirrors .NET EndpointTests. All 16 endpoint scenario tests share a single
// resource group and Storage Mover; each test uses a uniquely-named endpoint so tests are
// order-independent. All names are pre-generated in SetupSuite (under the suite's test name) so
// playback can resolve them deterministically — `recording.GenerateAlphaNumericID` is keyed on
// `t.Name()`, and SetupSuite is the only context where `s.T()` returns the suite-level testing.T.
type EndpointScenarioSuite struct {
	scenarioBaseSuite

	storageMoverName string

	// Names for TestEndpointCreateUpdateGetDelete (creates 4 endpoints).
	crudContainerName string
	crudNfsName       string
	crudSmbName       string
	crudFsName        string

	mccCRUDName      string
	s3HmacName       string
	nfsSrcName       string
	smbSrcName       string
	mccSrcName       string
	blobSrcName      string
	blobTgtName      string
	smbFsTgtName     string
	nfsFsTgtName     string
	nfsTgtFailName   string
	smbTgtFailName   string
	mccTgtFailName   string
	smbFsSrcFailName string
	nfsFsSrcFailName string
	nfsFsCRUDName    string
}

func TestEndpointScenarioSuite(t *testing.T) {
	suite.Run(t, new(EndpointScenarioSuite))
}

func (s *EndpointScenarioSuite) SetupSuite() {
	s.setupBase()
	s.storageMoverName = s.generateName("stomover")
	s.crudContainerName = s.generateName("conep")
	s.crudNfsName = s.generateName("nfsep")
	s.crudSmbName = s.generateName("smbep")
	s.crudFsName = s.generateName("fsep")
	s.mccCRUDName = s.generateName("mcc")
	s.s3HmacName = s.generateName("s3hmac")
	s.nfsSrcName = s.generateName("nfssrc")
	s.smbSrcName = s.generateName("smbsrc")
	s.mccSrcName = s.generateName("mccsrc")
	s.blobSrcName = s.generateName("blobsrc")
	s.blobTgtName = s.generateName("blobtgt")
	s.smbFsTgtName = s.generateName("smbfstgt")
	s.nfsFsTgtName = s.generateName("nfsfstgt")
	s.nfsTgtFailName = s.generateName("nfstgt")
	s.smbTgtFailName = s.generateName("smbtgt")
	s.mccTgtFailName = s.generateName("mcctgt")
	s.smbFsSrcFailName = s.generateName("smbfssrc")
	s.nfsFsSrcFailName = s.generateName("nfsfssrc")
	s.nfsFsCRUDName = s.generateName("nfsfs")
	s.createStorageMover(s.storageMoverName, nil, "")
}

func (s *EndpointScenarioSuite) TearDownSuite() { s.teardownBase() }

func (s *EndpointScenarioSuite) endpointsClient() *armstoragemover.EndpointsClient {
	c, err := armstoragemover.NewEndpointsClient(s.subscriptionID, s.cred, s.options)
	s.Require().NoError(err)
	return c
}

func (s *EndpointScenarioSuite) deleteEndpoint(client *armstoragemover.EndpointsClient, name string) {
	poller, err := client.BeginDelete(s.ctx, s.resourceGroupName, s.storageMoverName, name, nil)
	s.Require().NoError(err)
	_, err = testutil.PollForTest(s.ctx, poller)
	s.Require().NoError(err)
}

// TestEndpointCreateUpdateGetDelete mirrors .NET EndpointTests.CreateUpdateGetDeleteTest.
// Exercises 4 endpoint types (blob container, NFS mount, SMB mount, SMB file share), patches the
// SMB endpoint, and lists.
//
// Note on the SMB PATCH workaround: in api-version 2025-12-01 the RP requires a top-level Identity
// field on Endpoint update for non-NfsMount / non-MultiCloudConnector endpoints whose existing
// resource has no identity, even though the spec marks Identity Optional. We pass
// `Identity: {Type: ManagedServiceIdentityTypeNone}` as a workaround. See the cross-language
// playbook for the RP source bug location.
func (s *EndpointScenarioSuite) TestEndpointCreateUpdateGetDelete() {
	client := s.endpointsClient()

	// Container endpoint.
	cResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudContainerName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			BlobContainerName:        to.Ptr(scenarioContainerName),
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			Description:              to.Ptr("New container endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.crudContainerName, *cResp.Name)
	cProps, ok := cResp.Properties.(*armstoragemover.AzureStorageBlobContainerEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeAzureStorageBlobContainer, *cProps.EndpointType)

	cGet, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudContainerName, nil)
	s.Require().NoError(err)
	s.Equal(s.crudContainerName, *cGet.Name)

	// NFS endpoint.
	nfsResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudNfsName, armstoragemover.Endpoint{
		Properties: &armstoragemover.NfsMountEndpointProperties{
			Host:        to.Ptr(scenarioSmbHost),
			Export:      to.Ptr(scenarioNfsExport),
			Description: to.Ptr("New NFS endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.crudNfsName, *nfsResp.Name)
	nfsProps, ok := nfsResp.Properties.(*armstoragemover.NfsMountEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeNfsMount, *nfsProps.EndpointType)
	s.Equal(scenarioNfsExport, *nfsProps.Export)
	s.Equal(scenarioSmbHost, *nfsProps.Host)

	// SMB endpoint with credentials.
	smbResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudSmbName, armstoragemover.Endpoint{
		Properties: &armstoragemover.SmbMountEndpointProperties{
			Host:      to.Ptr(scenarioSmbHost),
			ShareName: to.Ptr(scenarioSmbShareName),
			Credentials: &armstoragemover.AzureKeyVaultSmbCredentials{
				UsernameURI: to.Ptr(scenarioKeyVaultUsernameURI),
				PasswordURI: to.Ptr(scenarioKeyVaultPasswordURI),
			},
			Description: to.Ptr("New Smb mount endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	smbProps, ok := smbResp.Properties.(*armstoragemover.SmbMountEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeSmbMount, *smbProps.EndpointType)
	s.Equal(scenarioKeyVaultUsernameURI, *smbProps.Credentials.UsernameURI)
	s.Equal(scenarioKeyVaultPasswordURI, *smbProps.Credentials.PasswordURI)
	s.Equal(scenarioSmbHost, *smbProps.Host)
	s.Equal(scenarioSmbShareName, *smbProps.ShareName)

	// PATCH the SMB endpoint to clear credentials. The Identity workaround is required.
	patchResp, err := client.Update(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudSmbName, armstoragemover.EndpointBaseUpdateParameters{
		Identity: &armstoragemover.ManagedServiceIdentity{
			Type: to.Ptr(armstoragemover.ManagedServiceIdentityTypeNone),
		},
		Properties: &armstoragemover.SmbMountEndpointUpdateProperties{
			Credentials: &armstoragemover.AzureKeyVaultSmbCredentials{
				UsernameURI: to.Ptr(""),
				PasswordURI: to.Ptr(""),
			},
			Description: to.Ptr("Update endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	patchProps, ok := patchResp.Properties.(*armstoragemover.SmbMountEndpointProperties)
	s.Require().True(ok)
	s.Equal("", *patchProps.Credentials.UsernameURI)
	s.Equal("", *patchProps.Credentials.PasswordURI)
	s.Equal(scenarioSmbHost, *patchProps.Host)
	s.Equal(scenarioSmbShareName, *patchProps.ShareName)

	// Delete SMB endpoint.
	s.deleteEndpoint(client, s.crudSmbName)

	// SMB file share endpoint.
	fsResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudFsName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageSmbFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioFileShareName),
			Description:              to.Ptr("new file share endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	fsProps, ok := fsResp.Properties.(*armstoragemover.AzureStorageSmbFileShareEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeAzureStorageSmbFileShare, *fsProps.EndpointType)
	s.Equal(scenarioFileShareName, *fsProps.FileShareName)
	s.Equal("new file share endpoint", *fsProps.Description)

	count := 0
	pager := client.NewListPager(s.resourceGroupName, s.storageMoverName, nil)
	for pager.More() {
		page, err := pager.NextPage(s.ctx)
		s.Require().NoError(err)
		count += len(page.Value)
	}
	s.Greater(count, 1)

	// Verify the SMB endpoint we deleted is gone.
	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.crudSmbName, nil)
	s.expectResponseError(err)
}

// TestMultiCloudConnectorEndpointCreateGetDelete mirrors .NET MultiCloudConnectorEndpointCreateGetDeleteTest.
func (s *EndpointScenarioSuite) TestMultiCloudConnectorEndpointCreateGetDelete() {
	client := s.endpointsClient()

	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.mccCRUDName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureMultiCloudConnectorEndpointProperties{
			MultiCloudConnectorID: to.Ptr(scenarioMultiCloudConnectorID),
			AwsS3BucketID:         to.Ptr(scenarioAwsS3BucketID),
			Description:           to.Ptr("Test multi-cloud connector endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.mccCRUDName, *createResp.Name)
	props, ok := createResp.Properties.(*armstoragemover.AzureMultiCloudConnectorEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeAzureMultiCloudConnector, *props.EndpointType)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.mccCRUDName, nil)
	s.Require().NoError(err)
	getProps := getResp.Properties.(*armstoragemover.AzureMultiCloudConnectorEndpointProperties)
	s.Equal("Test multi-cloud connector endpoint", *getProps.Description)
	s.NotNil(getProps.MultiCloudConnectorID)
	s.NotNil(getProps.AwsS3BucketID)

	s.deleteEndpoint(client, s.mccCRUDName)
	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.mccCRUDName, nil)
	s.expectResponseError(err)
}

// TestS3WithHmacEndpointCreateGetDelete mirrors .NET S3WithHmacEndpointCreateGetDeleteTest. The .NET
// test is `[Ignore]` but Python and Go can run it: placeholder Key Vault URIs and a synthetic source
// URI are accepted at the metadata level (real S3 only matters when running an actual job).
func (s *EndpointScenarioSuite) TestS3WithHmacEndpointCreateGetDelete() {
	client := s.endpointsClient()

	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.s3HmacName, armstoragemover.Endpoint{
		Properties: &armstoragemover.S3WithHmacEndpointProperties{
			SourceURI:   to.Ptr(scenarioS3SourceURI),
			SourceType:  to.Ptr(armstoragemover.S3WithHmacSourceTypeMINIO),
			Description: to.Ptr("Test S3 with HMAC endpoint"),
			Credentials: &armstoragemover.AzureKeyVaultS3WithHmacCredentials{
				AccessKeyURI: to.Ptr(scenarioKeyVaultAccessKeyURI),
				SecretKeyURI: to.Ptr(scenarioKeyVaultSecretKeyURI),
			},
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.s3HmacName, *createResp.Name)
	createProps, ok := createResp.Properties.(*armstoragemover.S3WithHmacEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeS3WithHmac, *createProps.EndpointType)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.s3HmacName, nil)
	s.Require().NoError(err)
	getProps := getResp.Properties.(*armstoragemover.S3WithHmacEndpointProperties)
	s.Equal(scenarioS3SourceURI, *getProps.SourceURI)
	s.Equal(armstoragemover.S3WithHmacSourceTypeMINIO, *getProps.SourceType)
	s.Equal("Test S3 with HMAC endpoint", *getProps.Description)
	s.Require().NotNil(getProps.Credentials)
	s.Equal(scenarioKeyVaultAccessKeyURI, *getProps.Credentials.AccessKeyURI)
	s.Equal(scenarioKeyVaultSecretKeyURI, *getProps.Credentials.SecretKeyURI)

	s.deleteEndpoint(client, s.s3HmacName)
	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.s3HmacName, nil)
	s.expectResponseError(err)
}

// ------------------------- EndpointKind valid (7 tests) ---------------------------

// TestNfsMountEndpointKindSourceSucceeds mirrors .NET NfsMountEndpointKindSource_Succeeds.
func (s *EndpointScenarioSuite) TestNfsMountEndpointKindSourceSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsSrcName, armstoragemover.Endpoint{
		Properties: &armstoragemover.NfsMountEndpointProperties{
			Host:         to.Ptr(scenarioSmbHost),
			Export:       to.Ptr(scenarioNfsExport),
			EndpointKind: to.Ptr(armstoragemover.EndpointKindSource),
			Description:  to.Ptr("NFS source endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.NfsMountEndpointProperties)
	s.Equal(armstoragemover.EndpointKindSource, *props.EndpointKind)
	s.deleteEndpoint(client, s.nfsSrcName)
}

// TestSmbMountEndpointKindSourceSucceeds mirrors .NET SmbMountEndpointKindSource_Succeeds.
func (s *EndpointScenarioSuite) TestSmbMountEndpointKindSourceSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.smbSrcName, armstoragemover.Endpoint{
		Properties: &armstoragemover.SmbMountEndpointProperties{
			Host:         to.Ptr(scenarioSmbHost),
			ShareName:    to.Ptr(scenarioSmbShareName),
			EndpointKind: to.Ptr(armstoragemover.EndpointKindSource),
			Description:  to.Ptr("SMB source endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.SmbMountEndpointProperties)
	s.Equal(armstoragemover.EndpointKindSource, *props.EndpointKind)
	s.deleteEndpoint(client, s.smbSrcName)
}

// TestMultiCloudConnectorEndpointKindSourceSucceeds mirrors .NET MultiCloudConnectorEndpointKindSource_Succeeds.
func (s *EndpointScenarioSuite) TestMultiCloudConnectorEndpointKindSourceSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.mccSrcName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureMultiCloudConnectorEndpointProperties{
			MultiCloudConnectorID: to.Ptr(scenarioMultiCloudConnectorID),
			AwsS3BucketID:         to.Ptr(scenarioAwsS3BucketID),
			EndpointKind:          to.Ptr(armstoragemover.EndpointKindSource),
			Description:           to.Ptr("Multi-cloud connector source endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.AzureMultiCloudConnectorEndpointProperties)
	s.Equal(armstoragemover.EndpointKindSource, *props.EndpointKind)
	s.deleteEndpoint(client, s.mccSrcName)
}

// TestBlobContainerEndpointKindSourceSucceeds mirrors .NET BlobContainerEndpointKindSource_Succeeds.
func (s *EndpointScenarioSuite) TestBlobContainerEndpointKindSourceSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.blobSrcName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			BlobContainerName:        to.Ptr(scenarioContainerName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindSource),
			Description:              to.Ptr("Blob container source endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.AzureStorageBlobContainerEndpointProperties)
	s.Equal(armstoragemover.EndpointKindSource, *props.EndpointKind)
	s.deleteEndpoint(client, s.blobSrcName)
}

// TestBlobContainerEndpointKindTargetSucceeds mirrors .NET BlobContainerEndpointKindTarget_Succeeds.
func (s *EndpointScenarioSuite) TestBlobContainerEndpointKindTargetSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.blobTgtName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageBlobContainerEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			BlobContainerName:        to.Ptr(scenarioContainerName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindTarget),
			Description:              to.Ptr("Blob container target endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.AzureStorageBlobContainerEndpointProperties)
	s.Equal(armstoragemover.EndpointKindTarget, *props.EndpointKind)
	s.deleteEndpoint(client, s.blobTgtName)
}

// TestSmbFileShareEndpointKindTargetSucceeds mirrors .NET SmbFileShareEndpointKindTarget_Succeeds.
func (s *EndpointScenarioSuite) TestSmbFileShareEndpointKindTargetSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.smbFsTgtName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageSmbFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioFileShareName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindTarget),
			Description:              to.Ptr("SMB file share target endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.AzureStorageSmbFileShareEndpointProperties)
	s.Equal(armstoragemover.EndpointKindTarget, *props.EndpointKind)
	s.deleteEndpoint(client, s.smbFsTgtName)
}

// TestNfsFileShareEndpointKindTargetSucceeds mirrors .NET NfsFileShareEndpointKindTarget_Succeeds.
func (s *EndpointScenarioSuite) TestNfsFileShareEndpointKindTargetSucceeds() {
	client := s.endpointsClient()
	resp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsFsTgtName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageNfsFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioNfsFileShareName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindTarget),
			Description:              to.Ptr("NFS file share target endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	props := resp.Properties.(*armstoragemover.AzureStorageNfsFileShareEndpointProperties)
	s.Equal(armstoragemover.EndpointKindTarget, *props.EndpointKind)
	s.deleteEndpoint(client, s.nfsFsTgtName)
}

// ------------------------- EndpointKind invalid (5 tests) ---------------------------

// TestNfsMountEndpointKindTargetFails mirrors .NET NfsMountEndpointKindTarget_Fails.
func (s *EndpointScenarioSuite) TestNfsMountEndpointKindTargetFails() {
	client := s.endpointsClient()
	_, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsTgtFailName, armstoragemover.Endpoint{
		Properties: &armstoragemover.NfsMountEndpointProperties{
			Host:         to.Ptr(scenarioSmbHost),
			Export:       to.Ptr(scenarioNfsExport),
			EndpointKind: to.Ptr(armstoragemover.EndpointKindTarget),
		},
	}, nil)
	s.expectResponseError(err)
}

// TestSmbMountEndpointKindTargetFails mirrors .NET SmbMountEndpointKindTarget_Fails.
func (s *EndpointScenarioSuite) TestSmbMountEndpointKindTargetFails() {
	client := s.endpointsClient()
	_, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.smbTgtFailName, armstoragemover.Endpoint{
		Properties: &armstoragemover.SmbMountEndpointProperties{
			Host:         to.Ptr(scenarioSmbHost),
			ShareName:    to.Ptr(scenarioSmbShareName),
			EndpointKind: to.Ptr(armstoragemover.EndpointKindTarget),
		},
	}, nil)
	s.expectResponseError(err)
}

// TestMultiCloudConnectorEndpointKindTargetFails mirrors .NET MultiCloudConnectorEndpointKindTarget_Fails.
func (s *EndpointScenarioSuite) TestMultiCloudConnectorEndpointKindTargetFails() {
	client := s.endpointsClient()
	_, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.mccTgtFailName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureMultiCloudConnectorEndpointProperties{
			MultiCloudConnectorID: to.Ptr(scenarioMultiCloudConnectorID),
			AwsS3BucketID:         to.Ptr(scenarioAwsS3BucketID),
			EndpointKind:          to.Ptr(armstoragemover.EndpointKindTarget),
		},
	}, nil)
	s.expectResponseError(err)
}

// TestSmbFileShareEndpointKindSourceFails mirrors .NET SmbFileShareEndpointKindSource_Fails.
func (s *EndpointScenarioSuite) TestSmbFileShareEndpointKindSourceFails() {
	client := s.endpointsClient()
	_, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.smbFsSrcFailName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageSmbFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioFileShareName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindSource),
		},
	}, nil)
	s.expectResponseError(err)
}

// TestNfsFileShareEndpointKindSourceFails mirrors .NET NfsFileShareEndpointKindSource_Fails.
func (s *EndpointScenarioSuite) TestNfsFileShareEndpointKindSourceFails() {
	client := s.endpointsClient()
	_, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsFsSrcFailName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageNfsFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioNfsFileShareName),
			EndpointKind:             to.Ptr(armstoragemover.EndpointKindSource),
		},
	}, nil)
	s.expectResponseError(err)
}

// TestNfsFileShareEndpointCreateGetDelete mirrors .NET NfsFileShareEndpointCreateGetDeleteTest.
func (s *EndpointScenarioSuite) TestNfsFileShareEndpointCreateGetDelete() {
	client := s.endpointsClient()
	createResp, err := client.CreateOrUpdate(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsFsCRUDName, armstoragemover.Endpoint{
		Properties: &armstoragemover.AzureStorageNfsFileShareEndpointProperties{
			StorageAccountResourceID: to.Ptr(s.placeholderStorageAccountResourceID()),
			FileShareName:            to.Ptr(scenarioNfsFileShareName),
			Description:              to.Ptr("Test NFS file share endpoint"),
		},
	}, nil)
	s.Require().NoError(err)
	s.Equal(s.nfsFsCRUDName, *createResp.Name)
	props, ok := createResp.Properties.(*armstoragemover.AzureStorageNfsFileShareEndpointProperties)
	s.Require().True(ok)
	s.Equal(armstoragemover.EndpointTypeAzureStorageNfsFileShare, *props.EndpointType)

	getResp, err := client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsFsCRUDName, nil)
	s.Require().NoError(err)
	getProps := getResp.Properties.(*armstoragemover.AzureStorageNfsFileShareEndpointProperties)
	s.Equal(scenarioNfsFileShareName, *getProps.FileShareName)
	s.Equal("Test NFS file share endpoint", *getProps.Description)
	s.NotNil(getProps.StorageAccountResourceID)

	s.deleteEndpoint(client, s.nfsFsCRUDName)
	_, err = client.Get(s.ctx, s.resourceGroupName, s.storageMoverName, s.nfsFsCRUDName, nil)
	s.expectResponseError(err)
}
