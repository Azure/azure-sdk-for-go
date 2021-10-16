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
	"github.com/stretchr/testify/assert"
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
	} else if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

// 1. ServiceClient -----------------------------------------------------------------------------------------------------

func getServiceClient(recording *testframework.Recording, accountType testAccountType, options *ClientOptions) (ServiceClient, error) {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transporter: recording,
				Retry:       policy.RetryOptions{MaxRetries: -1}}
		}
	}

	cred, err := getGenericCredential(recording, accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	serviceURL, _ := url.Parse("https://" + cred.AccountName() + ".file.core.windows.net/")
	serviceClient, err := NewServiceClient(serviceURL.String(), cred, options)

	return serviceClient, err
}

//nolint
func getServiceClientFromConnectionString(recording *testframework.Recording, accountType testAccountType, options *ClientOptions) (ServiceClient, error) {
	if recording != nil {
		if options == nil {
			options = &ClientOptions{
				Transporter: recording,
				Retry:       policy.RetryOptions{MaxRetries: -1}}
		}
	}

	connectionString, err := getConnectionString(recording, accountType)
	if err != nil {
		return ServiceClient{}, nil
	}
	primaryURL, cred, err := parseConnectionString(connectionString)
	if err != nil {
		return ServiceClient{}, nil
	}

	svcClient, err := NewServiceClient(primaryURL, cred, options)
	return svcClient, err
}

// 2. ShareClient  --------------------------------------------------------------------------------------------------------

func getShareClient(shareName string, s ServiceClient) (ShareClient, error) {
	return s.NewShareClient(shareName)
}

func createNewShare(_assert *assert.Assertions, shareName string, serviceClient ServiceClient) ShareClient {
	shareClient, err := getShareClient(shareName, serviceClient)
	_assert.Nil(err)

	cResp, err := shareClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	return shareClient
}

func delShare(_assert *assert.Assertions, shareClient ShareClient) {
	deleteShareResp, err := shareClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(deleteShareResp.RawResponse.StatusCode, 202)
}

// 3. DirectoryClient -----------------------------------------------------------------------------------------------------

func getDirectoryClientFromShare(_assert *assert.Assertions, directoryName string, shareClient ShareClient) DirectoryClient {
	dirClient := shareClient.NewDirectoryClient(directoryName)
	return dirClient
}

func createNewDirectoryFromShare(_assert *assert.Assertions, directoryName string, shareClient ShareClient) DirectoryClient {
	dirClient := getDirectoryClientFromShare(_assert, directoryName, shareClient)

	cResp, err := dirClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	return dirClient
}

func delDirectory(_assert *assert.Assertions, dirClient DirectoryClient) {
	resp, err := dirClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
}

// 4. FileClient -------------------------------------------------------------------------------------------------------

func getFileClientFromDirectory(_assert *assert.Assertions, fileName string, dirClient DirectoryClient) FileClient {
	fClient, err := dirClient.NewFileClient(fileName)
	_assert.Nil(err)
	return fClient
}

// This is a convenience method, No public API to create file URL from share now. This method uses share's root directory.
func getFileClientFromShare(_assert *assert.Assertions, fileName string, srClient ShareClient) FileClient {
	fClient, err := srClient.NewRootDirectoryClient().NewFileClient(fileName)
	_assert.Nil(err)
	return fClient
}

func createNewFileFromShare(_assert *assert.Assertions, fileName string, fileSize int64, srClient ShareClient) FileClient {
	dirClient := srClient.NewRootDirectoryClient()

	fClient := getFileClientFromDirectory(_assert, fileName, dirClient)

	cResp, err := fClient.Create(ctx, &CreateFileOptions{
		FileContentLength: to.Int64Ptr(fileSize),
	})
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	return fClient
}

func delFile(_assert *assert.Assertions, fileClient FileClient) {
	resp, err := fileClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
}

// 5. Data Generators --------------------------------------------------------------------------------------------------

func getReaderToGeneratedBytes(n int) io.ReadSeekCloser {
	r, _ := generateData(n)
	return internal.NopCloser(r)
}

//nolint
func getRandomDataAndReader(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n)
	rand.Read(data)
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
