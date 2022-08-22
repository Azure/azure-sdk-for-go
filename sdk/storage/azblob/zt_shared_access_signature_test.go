//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//func (s *AZBlobUnrecordedTestsSuite) TestSASContainerClient() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
//	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
//	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := serviceClient.NewContainerClient(containerName)
//
//	permissions := container.SASPermissions{
//		Read: true,
//		Add:  true,
//	}
//	start := time.Now().Add(-5 * time.Minute).UTC()
//	expiry := time.Now().Add(time.Hour)
//
//	sasUrl, err := containerClient.GetSASURL(permissions, start, expiry)
//	_require.Nil(err)
//
//	containerClient2, err := container.NewClientWithNoCredential(sasUrl, nil)
//	_require.Nil(err)
//
//	_, err = containerClient2.Create(ctx, &container.CreateOptions{Metadata: basicMetadata})
//	_require.NotNil(err)
//	validateBlobErrorCode(_require, err, bloberror.AuthorizationFailure)
//}
////
//func (s *AZBlobUnrecordedTestsSuite) TestSASContainerClient2() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
//	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
//	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := serviceClient.NewContainerClient(containerName)
//
//	sasUrlReadAdd, err := containerClient.GetSASURL(container.SASPermissions{Read: true, Add: true},
//		time.Now().Add(-5*time.Minute).UTC(), time.Now().Add(time.Hour))
//	_require.Nil(err)
//	_, err = containerClient.Create(ctx, &container.CreateOptions{Metadata: basicMetadata})
//	_require.Nil(err)
//
//	containerClient1, err := container.NewClientWithNoCredential(sasUrlReadAdd, nil)
//	_require.Nil(err)
//
//	_, err = containerClient1.GetProperties(ctx, nil)
//	_require.Nil(err)
//	//validateBlobErrorCode(_require, err, bloberror.AuthorizationFailure)
//	//
//	//sasUrlRCWL, err := containerClient.GetSASURL(container.SASPermissions{Add: true, Create: true, Delete: true, List: true},
//	//	time.Now().Add(-5*time.Minute).UTC(), time.Now().Add(time.Hour))
//	//_require.Nil(err)
//	//
//	//containerClient2, err := container.NewClientWithNoCredential(sasUrlRCWL, nil)
//	//_require.Nil(err)
//	//
//	//_, err = containerClient2.Create(ctx, nil)
//	//_require.Nil(err)
//}
