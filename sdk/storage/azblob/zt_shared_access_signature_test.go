//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//func (s *azblobUnrecordedTestSuite) TestSASServiceClient(t *testing.T) {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
//	accountKey := os.Getenv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
//	cred, err := NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := NewServiceClient(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName := generateContainerName(testName)
//
//	containerClient := createNewContainer(_require, containerName, serviceClient)
//	_require.Nil(err)
//	defer deleteContainer(_require, containerClient)
//
//	resources := AccountSASResourceTypes{
//		Object:    true,
//		Service:   true,
//		Container: true,
//	}
//	permissions := AccountSASPermissions{
//		Read:   true,
//		Add:    true,
//		Write:  true,
//		Create: true,
//		Update: true,
//		Delete: true,
//	}
//	services := AccountSASServices{
//		Blob: true,
//	}
//	start := time.Date(2021, time.August, 4, 1, 1, 0, 0, time.UTC)
//	expiry := time.Date(2022, time.August, 4, 1, 1, 0, 0, time.UTC)
//
//	sasUrl, err := serviceClient.GetSASURL(resources, permissions, services, start, expiry)
//	_require.Nil(err)
//
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	svcClient, err := (t, sasUrl, azcore.NewAnonymousCredential())
//	_require.Nil(err)
//
//	_, err = svcClient.CreateTable(context.Background(), containerName+"002", nil)
//	_require.Nil(err)
//
//	_, err = svcClient.DeleteTable(context.Background(), containerName+"002", nil)
//	_require.Nil(err)
//}
//
//func TestSASClient(t *testing.T) {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
//	accountKey := os.Getenv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
//	cred, err := NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := NewServiceClient(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName, err := createRandomName(t, containerNamePrefix)
//	_require.Nil(err)
//
//	delete := func() {
//		_, err := serviceClient.DeleteTable(context.Background(), containerName, nil)
//		_require.Nil(err)
//	}
//	defer delete()
//
//	_, err = serviceClient.CreateTable(context.Background(), containerName, nil)
//	_require.Nil(err)
//
//	permissions := SASPermissions{
//		Read: true,
//		Add:  true,
//	}
//	start := time.Date(2021, time.August, 4, 1, 1, 0, 0, time.UTC)
//	expiry := time.Date(2022, time.August, 4, 1, 1, 0, 0, time.UTC)
//
//	c := serviceClient.NewClient(containerName)
//	sasUrl, err := c.GetTableSASToken(permissions, start, expiry)
//	_require.Nil(err)
//
//	err = recording.StartRecording(t, pathToPackage, nil)
//	_require.Nil(err)
//	client, err := createClientForRecording(t, "", sasUrl, azcore.NewAnonymousCredential())
//	_require.Nil(err)
//	defer recording.StopRecording(t, nil) //nolint
//
//	entity := map[string]string{
//		"PartitionKey": "pk001",
//		"RowKey":       "rk001",
//		"Value":        "5",
//	}
//	marshalled, err := json.Marshal(entity)
//	_require.Nil(err)
//
//	_, err = client.AddEntity(context.Background(), marshalled, nil)
//	_require.Nil(err)
//}
//
//func TestSASClientReadOnly(t *testing.T) {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
//	accountKey := os.Getenv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
//	cred, err := NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := NewServiceClient(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName, err := createRandomName(t, containerNamePrefix)
//	_require.Nil(err)
//
//	delete := func() {
//		_, err := serviceClient.DeleteTable(context.Background(), containerName, nil)
//		_require.Nil(err)
//	}
//	defer delete()
//
//	_, err = serviceClient.CreateTable(context.Background(), containerName, nil)
//	_require.Nil(err)
//
//	client := serviceClient.NewClient(containerName)
//	err = insertNEntities("pk001", 4, client)
//	_require.Nil(err)
//
//	permissions := SASPermissions{
//		Read: true,
//	}
//	start := time.Date(2021, time.August, 4, 1, 1, 0, 0, time.UTC)
//	expiry := time.Date(2022, time.August, 4, 1, 1, 0, 0, time.UTC)
//
//	c := serviceClient.NewClient(containerName)
//	sasUrl, err := c.GetTableSASToken(permissions, start, expiry)
//	_require.Nil(err)
//
//	err = recording.StartRecording(t, pathToPackage, nil)
//	_require.Nil(err)
//	client, err = createClientForRecording(t, "", sasUrl, azcore.NewAnonymousCredential())
//	_require.Nil(err)
//	defer recording.StopRecording(t, nil) //nolint
//
//	entity := map[string]string{
//		"PartitionKey": "pk001",
//		"RowKey":       "rk001",
//		"Value":        "5",
//	}
//	marshalled, err := json.Marshal(entity)
//	_require.Nil(err)
//
//	// Failure on a read
//	_, err = client.AddEntity(context.Background(), marshalled, nil)
//	require.Error(t, err)
//
//	// Success on a list
//	pager := client.List(nil)
//	count := 0
//	for pager.NextPage(context.Background()) {
//		count += len(pager.PageResponse().Entities)
//	}
//
//	require.NoError(t, pager.Err())
//	require.Equal(t, 4, count)
//}
//
//func TestSASCosmosClientReadOnly(t *testing.T) {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	accountName := os.Getenv("TABLES_COSMOS_ACCOUNT_NAME")
//	accountKey := os.Getenv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
//	cred, err := NewSharedKeyCredential(accountName, accountKey)
//	_require.Nil(err)
//
//	serviceClient, err := NewServiceClient(fmt.Sprintf("https://%s.table.cosmos.azure.com/", accountName), cred, nil)
//	_require.Nil(err)
//
//	containerName, err := createRandomName(t, containerNamePrefix)
//	_require.Nil(err)
//
//	delete := func() {
//		_, err := serviceClient.DeleteTable(context.Background(), containerName, nil)
//		_require.Nil(err)
//	}
//	defer delete()
//
//	_, err = serviceClient.CreateTable(context.Background(), containerName, nil)
//	_require.Nil(err)
//
//	client := serviceClient.NewClient(containerName)
//	err = insertNEntities("pk001", 4, client)
//	_require.Nil(err)
//
//	permissions := SASPermissions{
//		Read: true,
//	}
//	start := time.Date(2021, time.August, 4, 1, 1, 0, 0, time.UTC)
//	expiry := time.Date(2022, time.August, 4, 1, 1, 0, 0, time.UTC)
//
//	c := serviceClient.NewClient(containerName)
//	sasUrl, err := c.GetTableSASToken(permissions, start, expiry)
//	_require.Nil(err)
//
//	err = recording.StartRecording(t, pathToPackage, nil)
//	_require.Nil(err)
//	client, err = createClientForRecording(t, "", sasUrl, azcore.NewAnonymousCredential())
//	_require.Nil(err)
//	defer recording.StopRecording(t, nil) //nolint
//
//	entity := map[string]string{
//		"PartitionKey": "pk001",
//		"RowKey":       "rk001",
//		"Value":        "5",
//	}
//	marshalled, err := json.Marshal(entity)
//	_require.Nil(err)
//
//	// Failure on a read
//	_, err = client.AddEntity(context.Background(), marshalled, nil)
//	require.Error(t, err)
//
//	// Success on a list
//	pager := client.List(nil)
//	count := 0
//	for pager.NextPage(context.Background()) {
//		count += len(pager.PageResponse().Entities)
//	}
//
//	require.NoError(t, pager.Err())
//	require.Equal(t, 4, count)
//}
