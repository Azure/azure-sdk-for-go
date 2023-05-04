//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob_test

// This test checks the storage challenge policy.
func (s *PageBlobUnrecordedTestsSuite) TestManagedDiskOAuth() {
	//_require := require.New(s.T())
	//
	//// Set up for this test.
	//// In Azure Portal create a managed disk.
	//// Under Access Control (IAM), ensure the "Data Operator for managed disks" role is added.
	//// Under Disk Export, check Enable data access authentication mode.
	//// Click on Generate URL and paste that URL below as urlWithSas.
	//cred, err := azidentity.NewDefaultAzureCredential(nil)
	//_require.NoError(err)
	//
	//urlWithSas := "https://md-XXXXX.blob.core.windows.net/XXXX/XXXX?sv=2018-03-28&sr=b&si=XXXXXXX&sig=XXXXXXXXXX"
	//
	//// Create a page blob client with OAuth
	//blobClient, err := pageblob.NewClient(urlWithSas, cred, nil)
	//_require.NoError(err)
	//
	//_, err = blobClient.GetProperties(context.TODO(), nil)
	//_require.NoError(err)
}
