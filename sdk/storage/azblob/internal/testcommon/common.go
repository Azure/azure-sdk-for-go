//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"hash/crc64"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/require"
)

const (
	ContainerPrefix             = "goc"
	BlobPrefix                  = "gotestblob"
	BlockBlobDefaultData        = "GoBlockBlobData"
	InvalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

func GenerateContainerName(testName string) string {
	return ContainerPrefix + GenerateEntityName(testName)
}

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
func GenerateName(prefix string) string {
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

func GenerateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func GenerateBlobName(testName string) string {
	return BlobPrefix + GenerateEntityName(testName)
}

func GetReaderToGeneratedBytes(n int) io.ReadSeekCloser {
	r, _ := GenerateData(n)
	return streaming.NopCloser(r)
}

func GetDataAndReader(testName string, n int) (*bytes.Reader, []byte) {
	// Random seed for data generation
	seed := int64(crc64.Checksum([]byte(testName), shared.CRC64Table))
	random := rand.New(rand.NewSource(seed))
	data := make([]byte, n)
	_, _ = random.Read(data)
	return bytes.NewReader(data), data
}

const random64BString string = "2SDgZj6RkKYzJpu04sweQek4uWHO8ndPnYlZ0tnFS61hjnFZ5IkvIGGY44eKABov"

func GenerateData(sizeInBytes int) (io.ReadSeekCloser, []byte) {
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
	return streaming.NopCloser(bytes.NewReader(data)), data
}

func GetRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func GetRelativeTimeFromAnchor(anchorTime *time.Time, amount time.Duration) time.Time {
	return anchorTime.Add(amount * time.Second)
}

func GenerateBlockIDsList(count int) []string {
	blockIDs := make([]string, count)
	for i := 0; i < count; i++ {
		blockIDs[i] = BlockIDIntToBase64(i)
	}
	return blockIDs
}

// BlockIDIntToBase64 functions convert an int block ID to a base-64 string and vice versa
func BlockIDIntToBase64(blockID int) string {
	binaryBlockID := (&[4]byte{})[:]
	binary.LittleEndian.PutUint32(binaryBlockID, uint32(blockID))
	return base64.StdEncoding.EncodeToString(binaryBlockID)
}

func BlobListToMap(list []string) map[string]bool {
	out := make(map[string]bool)

	for _, v := range list {
		out[v] = true
	}

	return out
}

func ValidateHTTPErrorCode(_require *require.Assertions, err error, code int) {
	_require.NotNil(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(responseErr.StatusCode, code)
	} else {
		_require.Equal(strings.Contains(err.Error(), strconv.Itoa(code)), true)
	}
}

func ValidateBlobErrorCode(_require *require.Assertions, err error, code bloberror.Code) {
	_require.NotNil(err)
	var responseErr *azcore.ResponseError
	errors.As(err, &responseErr)
	if responseErr != nil {
		_require.Equal(string(code), responseErr.ErrorCode)
	} else {
		_require.Contains(err.Error(), code)
	}
}

func ValidateUpload(ctx context.Context, _require *require.Assertions, blobClient *blockblob.Client) {
	resp, err := blobClient.DownloadStream(ctx, nil)
	_require.Nil(err)
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Len(data, 0)
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

func BeforeTest(t *testing.T, suite string, test string) {
	const urlRegex = `https://\S+\.blob\.core\.windows\.net`
	const tokenRegex = `(?:Bearer\s).*`
	//const queryParamRegex = `=([^&|\n|\t\s]+)` // Note: Add query param name before this
	require.NoError(t, recording.AddURISanitizer(FakeStorageURL, urlRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeStorageURL, urlRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source-authorization", FakeToken, tokenRegex, nil))
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	//testframework.AddHeaderRegexSanitizer("X-Ms-Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	//testframework.AddHeaderRegexSanitizer("Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	// TODO: more freezing
	//testframework.AddBodyRegexSanitizer("RequestId:00000000-0000-0000-0000-000000000000", `RequestId:\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`, nil)
	//testframework.AddBodyRegexSanitizer("Time:2022-08-11T00:21:56.4562741Z", `Time:\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d*)?Z`, nil)
	require.NoError(t, recording.Start(t, "sdk/storage/azblob/testdata", nil))
}

func AfterTest(t *testing.T, suite string, test string) {
	require.NoError(t, recording.Stop(t, nil))
}
