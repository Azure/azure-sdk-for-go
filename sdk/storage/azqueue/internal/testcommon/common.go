// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testcommon

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2/queueerror"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	RecordingDirectory = "sdk/storage/azqueue/testdata"
	QueuePrefix        = "goq"
	QueueDefaultData   = "this is some default data"
)

func GenerateQueueName(testName string) string {
	return QueuePrefix + GenerateEntityName(testName)
}

func GenerateName(prefix string) string {
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

func GenerateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func ValidateHTTPErrorCode(_require *require.Assertions, err error, code int) {
	_require.Error(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(responseErr.StatusCode, code)
	} else {
		_require.Equal(strings.Contains(err.Error(), strconv.Itoa(code)), true)
	}
}

// GetRequiredEnv gets an environment variable by name and returns an error if it is not found
func GetRequiredEnv(name string) (string, error) {
	env, ok := os.LookupEnv(name)
	if ok {
		return env, nil
	} else {
		return "", errors.New("Required environment variable not set: " + name)
	}
}

func SetupSuite(suite *suite.Suite) *recording.TestProxyInstance {
	proxy, err := recording.StartTestProxy(RecordingDirectory, nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	return proxy
}

func TearDownSuite(suite *suite.Suite, proxy *recording.TestProxyInstance) {
	err := recording.StopTestProxy(proxy)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func BeforeTest(t *testing.T, suite string, test string) {
	const urlRegex = `https://\S+\.queue\.core\.windows\.net`
	require.NoError(t, recording.AddURISanitizer(FakeStorageURL, urlRegex, nil))
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	// testframework.AddHeaderRegexSanitizer("X-Ms-Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	// TODO: more freezing
	require.NoError(t, recording.Start(t, RecordingDirectory, nil))
}

func AfterTest(t *testing.T, suite string, test string) {
	require.NoError(t, recording.Stop(t, nil))
}

func ValidateQueueErrorCode(_require *require.Assertions, err error, code queueerror.Code) {
	_require.Error(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(string(code), responseErr.ErrorCode)
	} else {
		_require.Contains(err.Error(), code)
	}
}
