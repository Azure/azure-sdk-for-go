//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

// TODO: Add tests for UserDelegationCredential
/*func (s *azblobUnrecordedTestSuite) TestUserDelegationCredential() {
	_require := require.New(s.T())
	testName := s.T().Name()

	clientOptions := azcore.ClientOptions{}
	opts1 := azidentity.ManagedIdentityCredentialOptions{ClientOptions: clientOptions, ID: nil}
	cred, err := azidentity.NewManagedIdentityCredential(&opts1)
	if err != nil {
		log.Fatal(err)
	}
	_require.Nil(err)

	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	// Set current and past time
	currentTime := time.Now().UTC().Add(-10 * time.Second).Format(SASTimeFormat)
	pastTime := time.Now().UTC().Add(-10 * time.Second).Add(48 * time.Hour).Format(SASTimeFormat)
	info := KeyInfo{
		Start:  &currentTime,
		Expiry: &pastTime,
	}

	ctx := context.Background()
	udkResp, err := svcClient.NewServiceClientWithUserDelegationCredential(ctx, info, nil, nil)
	if err != nil {
		s.Fail("Unable to create user delegation credential because " + err.Error())
	}

	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)
	src, err := containerClient.NewBlockBlobClient("src")
	if err != nil {
		s.Fail("Unable to fetch block blob client because " + err.Error())
	}

	// Get source blob url with OAuth Cred
	srcBlobParts, _ := NewBlobURLParts(src.URL())
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParametersWithUserDelegation(&udkResp)
	if err != nil {
		s.T().Fatal(err)
	}
}*/
