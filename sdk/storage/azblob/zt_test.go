package azblob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/shared"

	chk "gopkg.in/check.v1"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// For testing docs, see: https://labix.org/gocheck
// To test a specific test: go test -check.f MyTestSuite

// Hookup to the testing framework
func Test(t *testing.T) { chk.TestingT(t) }

type aztestsSuite struct{}

var _ = chk.Suite(&aztestsSuite{})

//func (s *aztestsSuite) TestRetryPolicyRetryReadsFromSecondaryHostField(c *chk.C) {
//	_, found := reflect.TypeOf(RetryOptions{}).FieldByName("RetryReadsFromSecondaryHost")
//	if !found {
//		// Make sure the RetryOption was not erroneously overwritten
//		c.Fatal("RetryOption's RetryReadsFromSecondaryHost field must exist in the Blob SDK - uncomment it and make sure the field is returned from the retryReadsFromSecondaryHost() method too!")
//	}
//}

const (
	containerPrefix             = "go"
	blobPrefix                  = "gotestblob"
	blockBlobDefaultData        = "GoBlockBlobData"
	validationErrorSubstring    = "validation failed"
	invalidHeaderErrorSubstring = "invalid header field" // error thrown by the http client
)

var ctx = context.Background()

//var basicHeaders = BlobHTTPHeaders{
//	ContentType:        "my_type",
//	ContentDisposition: "my_disposition",
//	CacheControl:       "control",
//	ContentMD5:         nil,
//	ContentLanguage:    "my_language",
//	ContentEncoding:    "my_encoding",
//}
//
//var basicMetadata = Metadata{"foo": "bar"}

type testPipeline struct{}

const testPipelineMessage string = "Test factory invoked"

func (tm testPipeline) Do(ctx context.Context, methodFactory pipeline.Factory, request pipeline.Request) (pipeline.Response, error) {
	return nil, errors.New(testPipelineMessage)
}

// This function generates an entity name by concatenating the passed prefix,
// the name of the test requesting the entity name, and the minute, second, and nanoseconds of the call.
// This should make it easy to associate the entities with their test, uniquely identify
// them, and determine the order in which they were created.
// Note that this imposes a restriction on the length of test names
func generateName(prefix string) string {
	// These next lines up through the for loop are obtaining and walking up the stack
	// trace to extrat the test name, which is stored in name
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

func generateContainerName() string {
	return generateName(containerPrefix)
}

func generateBlobName() string {
	return generateName(blobPrefix)
}

func getContainerClient(c *chk.C, s ServiceClient) (container ContainerClient, name string) {
	name = generateContainerName()
	container = s.NewContainerClient(name)

	return container, name
}

//func getBlockBlobURL(c *chk.C, container ContainerClient) (blob BlockBlobClient, name string) {
//	name = generateBlobName()
//	blob = container.NewBlockBlobClient(name)
//
//	return blob, name
//}
//
//func getAppendBlobURL(c *chk.C, container ContainerClient) (blob AppendBlobURL, name string) {
//	name = generateBlobName()
//	blob = container.NewAppendBlobURL(name)
//
//	return blob, name
//}
//
//func getPageBlobURL(c *chk.C, container ContainerClient) (blob PageBlobURL, name string) {
//	name = generateBlobName()
//	blob = container.NewPageBlobURL(name)
//
//	return
//}

func getReaderToRandomBytes(n int) *bytes.Reader {
	r, _ := getRandomDataAndReader(n)
	return r
}

func getRandomDataAndReader(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n, n)
	rand.Read(data)
	return bytes.NewReader(data), data
}

func createNewContainer(c *chk.C, bsu ServiceClient) (container ContainerClient, name string) {
	container, name = getContainerClient(c, bsu)

	cResp, err := container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return container, name
}

func deleteContainer(c *chk.C, container ContainerClient) {
	resp, err := container.Delete(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
}

func createNewContainerWithSuffix(c *chk.C, bsu ServiceClient, suffix string) (container ContainerClient, name string) {
	// The goal of adding the suffix is to be able to predetermine what order the containers will be in when listed.
	// We still need the container prefix to come first, though, to ensure only containers as a part of this test
	// are listed at all.
	name = generateName(containerPrefix + suffix)
	container = bsu.NewContainerClient(name)

	cResp, err := container.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	return container, name
}

//func createNewBlockBlob(c *chk.C, container ContainerClient) (blob BlockBlobClient, name string) {
//	blob, name = getBlockBlobURL(c, container)
//
//	cResp, err := blob.Upload(ctx, strings.NewReader(blockBlobDefaultData), BlobHTTPHeaders{},
//		nil, BlobAccessConditions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(cResp.StatusCode(), chk.Equals, 201)
//
//	return
//}
//
//func createNewAppendBlob(c *chk.C, container ContainerClient) (blob AppendBlobURL, name string) {
//	blob, name = getAppendBlobURL(c, container)
//
//	resp, err := blob.Create(ctx, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 201)
//	return
//}
//
//func createNewPageBlob(c *chk.C, container ContainerClient) (blob PageBlobURL, name string) {
//	blob, name = getPageBlobURL(c, container)
//
//	resp, err := blob.Create(ctx, PageBlobPageBytes*10, 0, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 201)
//	return
//}
//
//func createNewPageBlobWithSize(c *chk.C, container ContainerClient, sizeInBytes int64) (blob PageBlobURL, name string) {
//	blob, name = getPageBlobURL(c, container)
//
//	resp, err := blob.Create(ctx, sizeInBytes, 0, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 201)
//	return
//}
//
//func createBlockBlobWithPrefix(c *chk.C, container ContainerClient, prefix string) (blob BlockBlobClient, name string) {
//	name = prefix + generateName(blobPrefix)
//	blob = container.NewBlockBlobClient(name)
//
//	cResp, err := blob.Upload(ctx, strings.NewReader(blockBlobDefaultData), BlobHTTPHeaders{},
//		nil, BlobAccessConditions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(cResp.StatusCode(), chk.Equals, 201)
//	return
//}
//
//func deleteContainer(c *chk.C, container ContainerClient) {
//	resp, err := container.Delete(ctx, ContainerAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 202)
//}

func getGenericCredential(accountType string) (*shared.SharedKeyCredential, error) {
	accountNameEnvVar := accountType + "AZURE_STORAGE_ACCOUNT_NAME"
	accountKeyEnvVar := accountType + "AZURE_STORAGE_ACCOUNT_KEY"
	accountName, accountKey := os.Getenv(accountNameEnvVar), os.Getenv(accountKeyEnvVar)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(accountNameEnvVar + " and/or " + accountKeyEnvVar + " environment variables not specified.")
	}
	return shared.NewSharedKeyCredential(accountName, accountKey)
}

//
////getOAuthCredential can intake a OAuth credential from environment variables in one of the following ways:
////Direct: Supply a ADAL OAuth token in OAUTH_TOKEN and application ID in APPLICATION_ID to refresh the supplied token.
////Client secret: Supply a client secret in CLIENT_SECRET and application ID in APPLICATION_ID for SPN auth.
////TENANT_ID is optional and will be inferred as common if it is not explicitly defined.
//func getOAuthCredential(accountType string) (*TokenCredential, error) {
//	oauthTokenEnvVar := accountType + "OAUTH_TOKEN"
//	clientSecretEnvVar := accountType + "CLIENT_SECRET"
//	applicationIdEnvVar := accountType + "APPLICATION_ID"
//	tenantIdEnvVar := accountType + "TENANT_ID"
//	oauthToken, appId, tenantId, clientSecret := []byte(os.Getenv(oauthTokenEnvVar)), os.Getenv(applicationIdEnvVar), os.Getenv(tenantIdEnvVar), os.Getenv(clientSecretEnvVar)
//	if (len(oauthToken) == 0 && clientSecret == "") || appId == "" {
//		return nil, errors.New("(" + oauthTokenEnvVar + " OR " + clientSecretEnvVar + ") and/or " + applicationIdEnvVar + " environment variables not specified.")
//	}
//	if tenantId == "" {
//		tenantId = "common"
//	}
//
//	var Token adal.Token
//	if len(oauthToken) != 0 {
//		if err := json.Unmarshal(oauthToken, &Token); err != nil {
//			return nil, err
//		}
//	}
//
//	var spt *adal.ServicePrincipalToken
//
//	oauthConfig, err := adal.NewOAuthConfig("https://login.microsoftonline.com", tenantId)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(oauthToken) == 0 {
//		spt, err = adal.NewServicePrincipalToken(
//			*oauthConfig,
//			appId,
//			clientSecret,
//			"https://storage.azure.com")
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		spt, err = adal.NewServicePrincipalTokenFromManualToken(*oauthConfig,
//			appId,
//			"https://storage.azure.com",
//			Token,
//		)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	err = spt.Refresh()
//	if err != nil {
//		return nil, err
//	}
//
//	tc := NewTokenCredential(spt.Token().AccessToken, func(tc TokenCredential) time.Duration {
//		_ = spt.Refresh()
//		return time.Until(spt.Token().Expires())
//	})
//
//	return &tc, nil
//}

func getGenericBSU(accountType string) (ServiceClient, error) {
	credential, err := getGenericCredential(accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	blobPrimaryURL, _ := url.Parse("https://" + credential.AccountName() + ".blob.core.windows.net/")
	return NewServiceClient(blobPrimaryURL.String(), credential, nil)
}

func getBSU() ServiceClient {
	bsu, _ := getGenericBSU("")
	return bsu
}

func getAlternateBSU() (ServiceClient, error) {
	return getGenericBSU("SECONDARY_")
}

func getPremiumBSU() (ServiceClient, error) {
	return getGenericBSU("PREMIUM_")
}

func getBlobStorageBSU() (ServiceClient, error) {
	return getGenericBSU("BLOB_STORAGE_")
}

//func validateStorageError(c *chk.C, err error, code ServiceCodeType) {
//	serr, _ := err.(StorageError)
//	c.Assert(serr.ServiceCode(), chk.Equals, code)
//}

func getRelativeTimeGMT(amount time.Duration) time.Time {
	currentTime := time.Now().In(time.FixedZone("GMT", 0))
	currentTime = currentTime.Add(amount * time.Second)
	return currentTime
}

func generateCurrentTimeWithModerateResolution() time.Time {
	highResolutionTime := time.Now().UTC()
	return time.Date(highResolutionTime.Year(), highResolutionTime.Month(), highResolutionTime.Day(), highResolutionTime.Hour(), highResolutionTime.Minute(),
		highResolutionTime.Second(), 0, highResolutionTime.Location())
}

// Some tests require setting service properties. It can take up to 30 seconds for the new properties to be reflected across all FEs.
// We will enable the necessary property and try to run the test implementation. If it fails with an error that should be due to
// those changes not being reflected yet, we will wait 30 seconds and try the test again. If it fails this time for any reason,
// we fail the test. It is the responsibility of the the testImplFunc to determine which error string indicates the test should be retried.
// There can only be one such string. All errors that cannot be due to this detail should be asserted and not returned as an error string.
func runTestRequiringServiceProperties(c *chk.C, bsu ServiceClient, code string,
	enableServicePropertyFunc func(*chk.C, ServiceClient),
	testImplFunc func(*chk.C, ServiceClient) error,
	disableServicePropertyFunc func(*chk.C, ServiceClient)) {
	enableServicePropertyFunc(c, bsu)
	defer disableServicePropertyFunc(c, bsu)
	err := testImplFunc(c, bsu)
	// We cannot assume that the error indicative of slow update will necessarily be a StorageError. As in ListBlobs.
	if err != nil && err.Error() == code {
		time.Sleep(time.Second * 30)
		err = testImplFunc(c, bsu)
		c.Assert(err, chk.IsNil)
	}
}

//
//func enableSoftDelete(c *chk.C, bsu ServiceClient) {
//	days := int32(1)
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
//	c.Assert(err, chk.IsNil)
//}
//
//func disableSoftDelete(c *chk.C, bsu ServiceClient) {
//	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: false}})
//	c.Assert(err, chk.IsNil)
//}
//
//func validateUpload(c *chk.C, blobURL BlockBlobClient) {
//	resp, err := blobURL.Download(ctx, 0, 0, BlobAccessConditions{}, false)
//	c.Assert(err, chk.IsNil)
//	data, _ := ioutil.ReadAll(resp.Response().Body)
//	c.Assert(data, chk.HasLen, 0)
//}
