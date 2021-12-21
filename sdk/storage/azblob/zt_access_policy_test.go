// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestSetEmptyAccessPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, &SetAccessPolicyOptions{})
	require.NoError(t, err)
}

func TestSetAccessPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		AccessPolicy: &AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	require.NoError(t, err)
}

func TestSetMultipleAccessPolicies(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	id := "empty"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id2,
		AccessPolicy: &AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id3,
		AccessPolicy: &AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	require.NoError(t, err)

	// Make a Get to assert two access policies
	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp.SignedIdentifiers, 3)
}

func TestSetNullAccessPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	id := "null"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, len(resp.SignedIdentifiers), 1)
}

func TestContainerGetSetPermissionsMultiplePolicies(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := AccessPolicyPermission{Read: true, Write: true}.String()
	readOnly := AccessPolicyPermission{Read: true}.String()
	id1, id2 := "0000", "0001"
	permissions := []*SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)

	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp.SignedIdentifiers, permissions)
}

func TestContainerGetPermissionsPublicAccessNotNone(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	require.NoError(t, err)
	defer deleteContainer(t, containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, *resp.BlobPublicAccess, PublicAccessTypeBlob)
}

//func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessNone() {
//	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
//	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
//	svcClient := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, svcClient)
//	defer deleteContainer(assert.New(s.T()), containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	// Container is created with PublicAccessTypeBlob, so setting it to None will actually test that it is changed through this method
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert.NoError(err)
//
//	_assert.NoError(err)
//	bsu2, err := NewServiceClient(svcClient.URL(), azcore.AnonymousCredential(), nil)
//	_assert.NoError(err)
//
//	containerClient2 := bsu2.NewContainerClient(containerName)
//	blobURL2 := containerClient2.NewBlockBlobClient(blobName)
//
//	// Get permissions via the original container URL so the request succeeds
//	resp, err := containerClient.GetAccessPolicy(ctx, nil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//	_assert.NoError(err)
//
//	// If we cannot access a blob's data, we will also not be able to enumerate blobs
//	p := containerClient2.ListBlobsFlat(nil)
//	p.NextPage(ctx)
//	err = p.Err() // grab the next page
//	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
//
//	_, err = blobURL2.Download(ctx, nil)
//	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
//}

func TestContainerSetPermissionsPublicAccessTypeBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.BlobPublicAccess, PublicAccessTypeBlob)
}

func TestContainerSetPermissionsPublicAccessContainer(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	access := PublicAccessTypeContainer
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.BlobPublicAccess, PublicAccessTypeContainer)
}

////// TODO: After Pacer is ready
////func (s *azblobTestSuite) TestContainerSetPermissionsACLSinglePolicy() {
////	svcClient := getServiceClient()
////	credential, err := getGenericCredential(t)
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////	containerClient, containerName := createNewContainer(c, svcClient)
////	defer deleteContainer(assert.New(s.T()), containerClient)
////	_, blobName := createNewBlockBlob(c, containerClient)
////
////	start := time.Now().UTC().Add(-15 * time.Second)
////	expiry := start.Add(5 * time.Minute).UTC()
////	listOnly := AccessPolicyPermission{List: true}.String()
////	id := "0000"
////	permissions := []SignedIdentifier{{
////		ID: &id,
////		AccessPolicy: &AccessPolicy{
////			Start:      &start,
////			Expiry:     &expiry,
////			Permission: &listOnly,
////		},
////	}}
////
////	setAccessPolicyOptions := SetAccessPolicyOptions{
////		ContainerAcquireLeaseOptions: ContainerAcquireLeaseOptions{
////			ContainerACL: permissions,
////		},
////	}
////	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
////	_assert.NoError(err)
////
////	serviceSASValues := BlobSASSignatureValues{Identifier: "0000", ContainerName: containerName}
////	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	sasURL := svcClient.URL()
////	sasURL.RawQuery = queryParams.Encode()
////	sasPipeline := (NewAnonymousCredential(), PipelineOptions{})
////	sasBlobServiceURL := NewServiceURL(sasURL, sasPipeline)
////
////	// Verifies that the SAS can access the resource
////	sasContainer := sasBlobServiceURL.NewContainerClient(containerName)
////	resp, err := sasContainer.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.NoError(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////
////	// Verifies that successful sas access is not just because it's public
////	anonymousBlobService := NewServiceURL(svcClient.URL(), sasPipeline)
////	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
////	_, err = anonymousContainer.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
////}

func TestContainerSetPermissionsACLMoreThanFive(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	// failing in playback
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)
	permissions := make([]*SignedIdentifier, 6)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidXMLDocument)
}

func TestContainerSetPermissionsDeleteAndModifyACL(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)
	listOnly := AccessPolicyPermission{Read: true}.String()
	permissions := make([]*SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions1)
	require.NoError(t, err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp.SignedIdentifiers, 1)
	require.EqualValues(t, resp.SignedIdentifiers, permissions)
}

func TestContainerSetPermissionsDeleteAllPolicies(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Len(t, resp.SignedIdentifiers, len(permissions))
	require.EqualValues(t, resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: []*SignedIdentifier{},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.SignedIdentifiers)
}

func TestContainerSetPermissionsInvalidPolicyTimes(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)
}

func TestContainerSetPermissionsNilPolicySlice(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	require.NoError(t, err)
}

func TestContainerSetPermissionsSignedIdentifierTooLong(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	require.NoError(t, err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NotNil(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidXMLDocument)
}

func TestContainerSetPermissionsIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.BlobPublicAccess)
}

func TestContainerSetPermissionsIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NotNil(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestContainerSetPermissionsIfUnModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.BlobPublicAccess)
}

func TestContainerSetPermissionsIfUnModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()
	err := recording.SetBodilessMatcher(t, nil)
	require.NoError(t, err)

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.NotNil(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}
