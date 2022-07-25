//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type azblobTestSuite struct {
	suite.Suite
	mode testframework.RecordMode
}

//nolint
type azblobUnrecordedTestSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &azblobTestSuite{mode: testframework.Playback})
	//suite.Run(t, &azblobUnrecordedTestSuite{})
}

type testContext struct {
	recording *testframework.Recording
	context   *testframework.TestContext
}

// a map to store our created test contexts
var clientsMap = make(map[string]*testContext)

// recordedTestSetup is called before each test execution by the test suite's BeforeTest method
func recordedTestSetup(t *testing.T, mode testframework.RecordMode) {
	testName := t.Name()
	_require := require.New(t)

	// init the test framework
	_testContext := testframework.NewTestContext(
		func(msg string) { _require.FailNow(msg) },
		func(msg string) { t.Log(msg) },
		func() string { return testName })

	// mode should be test_framework.Playback.
	// This will automatically record if no test recording is available and playback if it is.
	recording, err := testframework.NewRecording(_testContext, mode)
	_require.Nil(err)

	_, err = recording.GetEnvVar(AccountNameEnvVar, testframework.NoSanitization)
	if err != nil {
		log.Fatal(err)
	}
	_, err = recording.GetEnvVar(AccountKeyEnvVar, testframework.Secret_Base64String)
	if err != nil {
		log.Fatal(err)
	}
	_ = recording.GetOptionalEnvVar(DefaultEndpointSuffixEnvVar, DefaultEndpointSuffix, testframework.NoSanitization)

	clientsMap[testName] = &testContext{recording: recording, context: &_testContext}
}

func getTestContext(key string) *testContext {
	return clientsMap[key]
}

func recordedTestTeardown(key string) {
	_context, ok := clientsMap[key]
	if ok && !(*_context.context).IsFailed() {
		_ = _context.recording.Stop()
	}
}

//nolint
func (s *azblobTestSuite) BeforeTest(suite string, test string) {
	// set up the test environment
	recordedTestSetup(s.T(), s.mode)
}

//nolint
func (s *azblobTestSuite) AfterTest(suite string, test string) {
	// teardown the test context
	recordedTestTeardown(s.T().Name())
}

//nolint
func (s *azblobUnrecordedTestSuite) BeforeTest(suite string, test string) {

}

//nolint
func (s *azblobUnrecordedTestSuite) AfterTest(suite string, test string) {

}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
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

////nolint
//func validateUpload(_require *require.Assertions, blobClient *blob.Client) {
//	resp, err := blobClient.Download(ctx, nil)
//	_require.Nil(err)
//	data, err := ioutil.ReadAll(resp.BodyReader(nil))
//	_require.Nil(err)
//	_require.Len(data, 0)
//}

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
		_require.Equal(responseErr.ErrorCode, string(code))
	} else {
		_require.Contains(err.Error(), code)
	}
}
