package testcommon

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

const (
	FilesystemPrefix            = "gofs"
	FilePrefix                  = "gotestfile"
	DefaultData                 = "Godatalakedata"
	InvalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

func GenerateFilesystemName(testName string) string {
	return FilesystemPrefix + GenerateEntityName(testName)
}

func GenerateEntityName(testName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(testName), "/", ""), "test", "")
}

func BeforeTest(t *testing.T, suite string, test string) {
	const blobURLRegex = `https://\S+\.blob\.core\.windows\.net`
	const dfsURLRegex = `https://\S+\.dfs\.core\.windows\.net`
	const tokenRegex = `(?:Bearer\s).*`

	require.NoError(t, recording.AddURISanitizer(FakeBlobStorageURL, blobURLRegex, nil))
	require.NoError(t, recording.AddURISanitizer(FakeDFSStorageURL, dfsURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeBlobStorageURL, blobURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source", FakeDFSStorageURL, dfsURLRegex, nil))
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-copy-source-authorization", FakeToken, tokenRegex, nil))
	// we freeze request IDs and timestamps to avoid creating noisy diffs
	// NOTE: we can't freeze time stamps as that breaks some tests that use if-modified-since etc (maybe it can be fixed?)
	//testframework.AddHeaderRegexSanitizer("X-Ms-Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	require.NoError(t, recording.AddHeaderRegexSanitizer("x-ms-request-id", "00000000-0000-0000-0000-000000000000", "", nil))
	//testframework.AddHeaderRegexSanitizer("Date", "Wed, 10 Aug 2022 23:34:14 GMT", "", nil)
	// TODO: more freezing
	//testframework.AddBodyRegexSanitizer("RequestId:00000000-0000-0000-0000-000000000000", `RequestId:\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`, nil)
	//testframework.AddBodyRegexSanitizer("Time:2022-08-11T00:21:56.4562741Z", `Time:\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d*)?Z`, nil)
	require.NoError(t, recording.Start(t, "sdk/storage/azdatalake/testdata", nil))
}

func AfterTest(t *testing.T, suite string, test string) {
	require.NoError(t, recording.Stop(t, nil))
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
