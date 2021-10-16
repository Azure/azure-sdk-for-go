package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

func (s *azfileLiveTestSuite) TestFileCreateDeleteDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareName := generateShareName(testName)
	shareClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, shareClient)

	// Create and delete file in root directory.
	fileName := generateFileName(testName)
	file, err := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_assert.Nil(err)

	cResp, err := file.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.NotEqual(cResp.ETag, azcore.ETag(""))
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(cResp.RequestID, "")
	_assert.NotEqual(cResp.Version, "")
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotNil(cResp.IsServerEncrypted)

	delResp, err := file.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(delResp.RawResponse.StatusCode, 202)
	_assert.NotEqual(delResp.RequestID, "")
	_assert.NotEqual(delResp.Version, "")
	_assert.Equal(delResp.Date.IsZero(), false)

	dirName := generateDirectoryName(testName)
	dirClient := createNewDirectoryFromShare(_assert, dirName, shareClient)
	defer delDirectory(_assert, dirClient)

	//// Create and delete file in named directory.
	//file = dir.NewFileURL(generateFileName())
	//
	//cResp, err = file.Create(context.Background(), 0, azfile.FileHTTPHeaders{}, nil)
	//_assert.Nil(err)
	//_assert(cResp.Response().StatusCode, chk.Equals, 201)
	//_assert(cResp.ETag(), chk.Not(chk.Equals), azfile.ETagNone)
	//_assert(cResp.LastModified().IsZero(), chk.Equals, false)
	//_assert(cResp.RequestID(), chk.Not(chk.Equals), "")
	//_assert(cResp.Version(), chk.Not(chk.Equals), "")
	//_assert(cResp.Date().IsZero(), chk.Equals, false)
	//_assert(cResp.IsServerEncrypted(), chk.NotNil)
	//
	//delResp, err = file.Delete(context.Background())
	//_assert.Nil(err)
	//_assert(delResp.Response().StatusCode, chk.Equals, 202)
	//_assert(delResp.RequestID(), chk.Not(chk.Equals), "")
	//_assert(delResp.Version(), chk.Not(chk.Equals), "")
	//_assert(delResp.Date().IsZero(), chk.Equals, false)
}
