//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

// Index
// 0. AuthN/AuthZ
// 1. ServiceClient
// 2. ShareClient
// 3. DirectoryClient
// 4. FileClient
// 5. Data Generators
// 6. Utility functions

// 0. AuthN/AuthZ ------------------------------------------------------------------------------------------------------
type testAccountType string

const (
	testAccountDefault   testAccountType = ""
	testAccountSecondary testAccountType = "SECONDARY_"
	testAccountPremium   testAccountType = "PREMIUM_"
)

// getRequiredEnv gets an environment variable by name and returns an error if it is not found
func getRequiredEnv(name string) (string, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		return env, nil
	} else {
		return "", errors.New("Required environment variable not set: " + name)
	}
}

func getAccountInfo(recording *testframework.Recording, accountType testAccountType) (string, string, error) {
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	accountName, accountKey := "", ""
	var err error
	if recording == nil {
		accountName, err = getRequiredEnv(accountNameEnvVar)
		if err != nil {
			return "", "", err
		}
		accountKey, err = getRequiredEnv(accountKeyEnvVar)
		if err != nil {
			return "", "", err
		}

	} else {
		accountName, err = recording.GetEnvVar(accountNameEnvVar, testframework.NoSanitization)
		if err != nil {
			return "", "", err
		}
		accountKey, err = recording.GetEnvVar(accountKeyEnvVar, testframework.Secret_Base64String)
		if err != nil {
			return "", "", err
		}
	}
	return accountName, accountKey, nil
}

//nolint
func getConnectionString(recording *testframework.Recording, accountType testAccountType) (string, error) {
	accountName, accountKey, err := getAccountInfo(recording, accountType)
	if err != nil {
		return "", err
	}
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return connectionString, nil
}

func getGenericCredential(recording *testframework.Recording, accountType testAccountType) (*SharedKeyCredential, error) {
	accountName, accountKey, err := getAccountInfo(recording, accountType)
	if err != nil {
		return nil, err
	}
	if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

// 1. ServiceClient -----------------------------------------------------------------------------------------------------

func getServiceClient(_require *require.Assertions, recording *testframework.Recording, accountType testAccountType, options *ClientOptions) *ServiceClient {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transport: recording,
				Retry:     policy.RetryOptions{MaxRetries: -1},
			}
		}
	}

	cred, err := getGenericCredential(recording, accountType)
	if err != nil {
		_require.Nil("Unable to fetch service client because " + err.Error())
	}

	serviceURL, _ := url.Parse("https://" + cred.AccountName() + ".file.core.windows.net/")
	svcClient, err := NewServiceClientWithSharedKey(serviceURL.String(), cred, options)
	if err != nil {
		_require.Nil("Unable to fetch service client because " + err.Error())
	}
	return svcClient
}

//nolint
func getServiceClientFromConnectionString(recording *testframework.Recording, accountType testAccountType,
	options *ClientOptions) (*ServiceClient, error) {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transport: recording,
				Retry:     policy.RetryOptions{MaxRetries: -1},
			}
		}
	}

	connectionString, err := getConnectionString(recording, accountType)
	if err != nil {
		return nil, err
	}

	primaryURL, cred, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	svcClient, err := NewServiceClientWithSharedKey(primaryURL, cred, options)
	return svcClient, err
}

// 2. ShareClient  --------------------------------------------------------------------------------------------------------

func getShareClient(_require *require.Assertions, shareName string, s *ServiceClient) *ShareClient {
	srClient, err := s.NewShareClient(shareName)
	_require.Nil(err)
	return srClient
}

func createNewShare(_require *require.Assertions, shareName string, serviceClient *ServiceClient) *ShareClient {
	srClient := getShareClient(_require, shareName, serviceClient)
	_, err := srClient.Create(ctx, nil)
	_require.Nil(err)
	return srClient
}

func delShare(_require *require.Assertions, srClient *ShareClient, options *ShareDeleteOptions) {
	_, err := srClient.Delete(context.Background(), options)
	_require.Nil(err)
	//_require.Equal(deleteShareResp.RawResponse.StatusCode, 202)
}

// 3. DirectoryClient -----------------------------------------------------------------------------------------------------

func getDirectoryClientFromShare(_require *require.Assertions, dirName string, srClient *ShareClient) *DirectoryClient {
	dirClient, err := srClient.NewDirectoryClient(dirName)
	_require.Nil(err)
	return dirClient
}

func createNewDirectoryFromShare(_require *require.Assertions, dirName string, srClient *ShareClient) *DirectoryClient {
	dirClient := getDirectoryClientFromShare(_require, dirName, srClient)

	_, err := dirClient.Create(ctx, nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	return dirClient
}

func delDirectory(_require *require.Assertions, dirClient *DirectoryClient) {
	_, err := dirClient.Delete(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(resp.RawResponse.StatusCode, 202)
}

// 4. FileClient -------------------------------------------------------------------------------------------------------

func getFileClientFromDirectory(_require *require.Assertions, fileName string, dirClient *DirectoryClient) *FileClient {
	fClient, err := dirClient.NewFileClient(fileName)
	_require.Nil(err)
	return fClient
}

// This is a convenience method, No public API to create file URL from share now. This method uses share's root directory.
func getFileClientFromShare(_require *require.Assertions, fileName string, srClient *ShareClient) *FileClient {
	dirClient, err := srClient.NewRootDirectoryClient()
	_require.Nil(err)
	fClient, err := dirClient.NewFileClient(fileName)
	_require.Nil(err)
	return fClient
}

func createNewFileFromShare(_require *require.Assertions, fileName string, fileSize int64,
	srClient *ShareClient) *FileClient {
	dirClient, err := srClient.NewRootDirectoryClient()
	_require.Nil(err)

	fClient := getFileClientFromDirectory(_require, fileName, dirClient)

	_, err = fClient.Create(ctx, &FileCreateOptions{
		FileContentLength: to.Ptr(fileSize),
	})
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	return fClient
}

func createNewFileFromShareWithPermissions(_require *require.Assertions, fileName string,
	fileSize int64, srClient *ShareClient) *FileClient {
	fClient := getFileClientFromShare(_require, fileName, srClient)

	_, err := fClient.Create(ctx, &FileCreateOptions{
		FileContentLength: to.Ptr(fileSize),
		FilePermissions: &FilePermissions{
			PermissionStr: &sampleSDDL,
		},
	})
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	return fClient
}

// This is a convenience method, No public API to create file URL from share now. This method uses share's root directory.
func createNewFileFromShareWithGivenData(_require *require.Assertions, fileName string,
	fileData string, srClient *ShareClient) *FileClient {
	fClient := getFileClientFromShare(_require, fileName, srClient)

	_, err := fClient.Create(ctx, &FileCreateOptions{
		FileContentLength: to.Ptr(int64(len(fileData))),
		FilePermissions: &FilePermissions{
			PermissionStr: &sampleSDDL,
		},
	})
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	putResp, err := fClient.UploadRange(ctx, 0, internal.NopCloser(strings.NewReader(fileDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotEqual(putResp.ETag, "")
	_require.NotEqual(putResp.RequestID, "")
	_require.NotEqual(putResp.Version, "")
	_require.Equal(putResp.Date.IsZero(), false)

	return fClient
}

func delFile(_require *require.Assertions, fileClient *FileClient) {
	_, err := fileClient.Delete(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(resp.RawResponse.StatusCode, 202)
}

// 5. Data Generators --------------------------------------------------------------------------------------------------

func getReaderToGeneratedBytes(n int) io.ReadSeekCloser {
	r, _ := generateData(n)
	return internal.NopCloser(r)
}

//nolint
func getRandomDataAndReader(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return nil, nil
	}
	return bytes.NewReader(data), data
}

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"

func generateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
	data := make([]byte, sizeInBytes)
	_len := len(random64BString)
	if sizeInBytes > _len {
		count := sizeInBytes / _len
		if sizeInBytes%_len != 0 {
			count = count + 1
		}
		copy(data[:], strings.Repeat(random64BString, count))
	} else {
		copy(data[:], random64BString)
	}
	return internal.NopCloser(bytes.NewReader(data)), data
}

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
//nolint
func generateName(prefix string) string {
	// These next lines up through the for loop are obtaining and walking up the stack
	// trace to extract the test name, which is stored in name
	pc := make([]uintptr, 10)
	runtime.Callers(0, pc)
	frames := runtime.CallersFrames(pc)
	name := ""
	for f, next := frames.Next(); next; f, next = frames.Next() {
		name = f.Function
		if strings.Contains(name, "Suite") {
			break
		}
	}
	funcNameStart := strings.Index(name, "Test")
	name = name[funcNameStart+len("Test"):] // Just get the name of the test and not any of the garbage at the beginning
	name = strings.ToLower(name)            // Ensure it is a valid resource name
	currentTime := time.Now()
	name = fmt.Sprintf("%s%s%d%d%d", prefix, strings.ToLower(name), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond())
	return name
}

func generateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func generateShareName(testName string) string {
	return sharePrefix + generateEntityName(testName)
}

func generateDirectoryName(testName string) string {
	return directoryPrefix + generateEntityName(testName)
}

func generateFileName(testName string) string {
	return filePrefix + generateEntityName(testName)
}

func generateEntityNameWithPrefix(prefix, testName string) string {
	return prefix + generateEntityName(testName)
}

// 6. Utility Functions ---------------------------------------------------------------------------------------------------

//nolint
func getRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func getRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}
