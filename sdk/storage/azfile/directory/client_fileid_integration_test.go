// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/stretchr/testify/require"
)

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesByFileID() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	pathProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(pathProps.FileName)

	dirIDClient := shareClient.NewDirectoryClientByFileID(*pathProps.ID)
	_require.Equal(*pathProps.ID, dirIDClient.FileID())

	resp, err := dirIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*pathProps.ID, *resp.ID)
	_require.Equal(*pathProps.ParentID, *resp.ParentID)
	_require.Equal(dirName, *resp.FileName)
	_require.Equal(*pathProps.FileAttributes, *resp.FileAttributes)

	// unsupported operations fail before a request is sent, so the directory is left untouched
	_, err = dirIDClient.Delete(context.Background(), nil)
	_require.ErrorIs(err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesByFileIDError() {
	_require := require.New(d.T())

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(d.T().Name()), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.NewDirectoryClientByFileID("11111111111111111111").GetProperties(context.Background(), nil)
	_require.Error(err)
}
