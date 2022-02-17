// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetEmptyAccessPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, &SetAccessPolicyOptions{})
	require.NoError(t, err)
}

func TestSetAccessPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: to.StringPtr("empty"),
	})

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: to.StringPtr("partial"),
		AccessPolicy: &AccessPolicy{
			Permission: to.StringPtr("r"),
		},
	})

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: to.StringPtr("full"),
		AccessPolicy: &AccessPolicy{
			Start:      to.TimePtr(time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)),
			Expiry:     to.TimePtr(time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)),
			Permission: to.StringPtr("f"),
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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: to.StringPtr("null"),
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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	require.NoError(t, err)
	defer deleteContainer(_assert, containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, *resp.BlobPublicAccess, PublicAccessTypeBlob)
}

func TestContainerSetPermissionsPublicAccessTypeBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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

func TestContainerSetPermissionsACLMoreThanFive(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

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

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func TestContainerSetPermissionsDeleteAndModifyACL(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	require.NoError(t, err)
}

func TestContainerSetPermissionsSignedIdentifierTooLong(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(_assert, containerClient)

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
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func TestContainerSetPermissionsIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestContainerSetPermissionsIfUnModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

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

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}
