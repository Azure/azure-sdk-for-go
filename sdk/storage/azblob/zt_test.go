//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type azblobTestSuite struct {
	suite.Suite
}

// nolint
type azblobUnrecordedTestSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &azblobTestSuite{})
	//suite.Run(t, &azblobUnrecordedTestSuite{})
}

const (
	fakeStorageAccount = "fakestorage"
	fakeStorageURL     = "https://fakestorage.blob.core.windows.net"
)

// nolint
func (s *azblobTestSuite) BeforeTest(suite string, test string) {
	const urlRegex = `https://\S+\.blob\.core\.windows\.net`
	recording.AddURISanitizer(fakeStorageURL, urlRegex, nil)
	recording.AddHeaderRegexSanitizer("x-ms-copy-source", fakeStorageURL, urlRegex, nil)
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	//testframework.AddHeaderRegexSanitizer("X-Ms-Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil)
	//testframework.AddHeaderRegexSanitizer("Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	// TODO: more freezing
	//testframework.AddBodyRegexSanitizer("RequestId:00000000-0000-0000-0000-000000000000", `RequestId:\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`, nil)
	//testframework.AddBodyRegexSanitizer("Time:2022-08-11T00:21:56.4562741Z", `Time:\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d*)?Z`, nil)
	require.NoError(s.T(), recording.Start(s.T(), "sdk/storage/azblob/testdata", nil))
}

// nolint
func (s *azblobTestSuite) AfterTest(suite string, test string) {
	require.NoError(s.T(), recording.Stop(s.T(), nil))
}

// nolint
func (s *azblobUnrecordedTestSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *azblobUnrecordedTestSuite) AfterTest(suite string, test string) {

}

// Some tests require setting service properties. It can take up to 30 seconds for the new properties to be reflected across all FEs.
// We will enable the necessary property and try to run the test implementation. If it fails with an error that should be due to
// those changes not being reflected yet, we will wait 30 seconds and try the test again. If it fails this time for any reason,
// we fail the test. It is the responsibility of the testImplFunc to determine which error string indicates the test should be retried.
// There can only be one such string. All errors that cannot be due to this detail should be asserted and not returned as an error string.
func runTestRequiringServiceProperties(_require *require.Assertions, svcClient *service.Client, code string,
	enableServicePropertyFunc func(*require.Assertions, *service.Client),
	testImplFunc func(*require.Assertions, *service.Client) error,
	disableServicePropertyFunc func(*require.Assertions, *service.Client)) {

	enableServicePropertyFunc(_require, svcClient)
	defer disableServicePropertyFunc(_require, svcClient)

	err := testImplFunc(_require, svcClient)
	// We cannot assume that the error indicative of slow update will necessarily be a StorageError. As in ListBlobs.
	if err != nil && err.Error() == code {
		time.Sleep(time.Second * 30)
		err = testImplFunc(_require, svcClient)
		_require.Nil(err)
	}
}

func enableSoftDelete(_require *require.Assertions, client *service.Client) {
	days := int32(1)
	_, err := client.SetProperties(ctx, &service.SetPropertiesOptions{
		DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: to.Ptr(true), Days: &days}})
	_require.Nil(err)
}

func disableSoftDelete(_require *require.Assertions, client *service.Client) {
	_, err := client.SetProperties(ctx, &service.SetPropertiesOptions{DeleteRetentionPolicy: &service.RetentionPolicy{Enabled: to.Ptr(false)}})
	_require.Nil(err)
}

// nolint
func validateUpload(_require *require.Assertions, blobClient *blockblob.Client) {
	resp, err := blobClient.DownloadStream(ctx, nil)
	_require.Nil(err)
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Len(data, 0)
}

func validateHTTPErrorCode(_require *require.Assertions, err error, code int) {
	_require.NotNil(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(responseErr.StatusCode, code)
	} else {
		_require.Equal(strings.Contains(err.Error(), strconv.Itoa(code)), true)
	}
}

func validateBlobErrorCode(_require *require.Assertions, err error, code bloberror.Code) {
	_require.NotNil(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(string(code), responseErr.ErrorCode)
	} else {
		_require.Contains(err.Error(), code)
	}
}
