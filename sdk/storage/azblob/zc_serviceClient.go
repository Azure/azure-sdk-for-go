package azblob

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	// ContainerNameRoot is the special Azure Storage name used to identify a storage account's root container.
	ContainerNameRoot = "$root"

	// ContainerNameLogs is the special Azure Storage name used to identify a storage account's logs container.
	ContainerNameLogs = "$logs"
)

// A ServiceClient represents a URL to the Azure Storage Blob service allowing you to manipulate blob containers.
type ServiceClient struct {
	client *serviceClient
	u      url.URL
}

// NewServiceClient creates a ServiceClient object using the specified URL, credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClient(serviceURL string, cred azcore.Credential, options *connectionOptions) (ServiceClient, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return ServiceClient{}, err
	}
	return ServiceClient{client: &serviceClient{
		con: newConnection(serviceURL, cred, options),
	}, u: *u}, nil
}

// URL returns the URL endpoint used by the ServiceClient object.
func (s ServiceClient) URL() url.URL {
	return s.u
}

// String returns the URL as a string.
func (s ServiceClient) String() string {
	return s.u.String()
}

// WithPipeline creates a new ServiceClient object identical to the source but with the specified request policy pipeline.
func (s ServiceClient) WithPipeline(pipeline azcore.Pipeline) ServiceClient {
	connection := newConnectionWithPipeline(s.u.String(), pipeline)
	return ServiceClient{client: &serviceClient{con: connection},
		u: s.u,
	}
}

// NewContainerClient creates a new ContainerClient object by concatenating containerName to the end of
// ServiceClient's URL. The new ContainerClient uses the same request policy pipeline as the ServiceClient.
// To change the pipeline, create the ContainerClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewContainerClient instead of calling this object's
// NewContainerClient method.
func (s ServiceClient) NewContainerClient(containerName string) ContainerClient {
	containerURL := appendToURLPath(s.u, containerName)
	containerConnection := newConnectionWithPipeline(containerURL.String(), s.client.con.p)
	return ContainerClient{
		client: &containerClient{
			con: containerConnection,
		},
		u: containerURL,
	}
}

// appendToURLPath appends a string to the end of a URL's path (prefixing the string with a '/' if required)
func appendToURLPath(u url.URL, name string) url.URL {
	// e.g. "https://ms.com/a/b/?k1=v1&k2=v2#f"
	// When you call url.Parse() this is what you'll get:
	//     Scheme: "https"
	//     Opaque: ""
	//       User: nil
	//       Host: "ms.com"
	//       Path: "/a/b/"	This should start with a / and it might or might not have a trailing slash
	//    RawPath: ""
	// ForceQuery: false
	//   RawQuery: "k1=v1&k2=v2"
	//   Fragment: "f"
	if len(u.Path) == 0 || u.Path[len(u.Path)-1] != '/' {
		u.Path += "/" // Append "/" to end before appending name
	}
	u.Path += name
	return u
}

func (s ServiceClient) GetAccountInfo(ctx context.Context) (ServiceGetAccountInfoResponse, error) {
	return s.client.GetAccountInfo(ctx, nil)
}

//GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceClient object.
//OAuth is required for this call, as well as any role that can delegate access to the storage account.
func (s ServiceClient) GetUserDelegationCredential(ctx context.Context, info KeyInfo) (UserDelegationCredential, error) {
	udk, err := s.client.GetUserDelegationKey(ctx, info, nil)
	if err != nil {
		return UserDelegationCredential{}, err
	}
	urlParts := NewBlobURLParts(s.URL())
	return NewUserDelegationCredential(strings.Split(urlParts.Host, ".")[0], *udk.UserDelegationKey), nil
}

//NewKeyInfo creates a new KeyInfo struct with the correct time formatting & conversion
func NewKeyInfo(Start, Expiry time.Time) KeyInfo {
	start := Start.UTC().Format(SASTimeFormat)
	expiry := Expiry.UTC().Format(SASTimeFormat)
	return KeyInfo{
		Start:  &start,
		Expiry: &expiry,
	}
}

// The List Containers Segment operation returns a channel of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
// The returned channel contains individual container items.
// AutoPagerTimeout specifies the amount of time with no read operations before the channel times out and closes. Specify no time and it will be ignored.
// AutoPagerBufferSize specifies the channel's buffer size.
// Both the blob item channel and error channel should be watched. Only one error will be released via this channel (or a nil error, to register a clean exit.)
func (s ServiceClient) ListContainersSegment(ctx context.Context, AutoPagerBufferSize uint, AutoPagerTimeout time.Duration, o *ListContainersSegmentOptions) (chan ContainerItem, chan error) {
	output := make(chan ContainerItem, AutoPagerBufferSize)
	errChan := make(chan error, 1)

	listOptions := o.pointers()
	pager := s.client.ListContainersSegment(listOptions)
	// override the generated advancer, which is incorrect
	if pager.Err() != nil {
		errChan <- pager.Err()
		close(output)
		close(errChan)
		return output, errChan
	}

	p := pager.(*listContainersSegmentResponsePager) // cast to the internal type first
	p.advancer = func(cxt context.Context, response ListContainersSegmentResponseResponse) (*azcore.Request, error) {
		if response.EnumerationResults.NextMarker == nil {
			return nil, errors.New("unexpected missing NextMarker")
		}
		req, err := s.client.listContainersSegmentCreateRequest(ctx, listOptions)
		if err != nil {
			return nil, err
		}
		queryValues, _ := url.ParseQuery(req.URL.RawQuery)
		queryValues.Set("marker", *response.EnumerationResults.NextMarker)

		req.URL.RawQuery = queryValues.Encode()
		return req, nil
	}

	go listContainersSegmentAutoPager{
		pager,
		output,
		errChan,
		ctx,
		AutoPagerTimeout,
		nil,
	}.Go()

	return output, errChan
}

func (s ServiceClient) GetProperties(ctx context.Context) (StorageServicePropertiesResponse, error) {
	return s.client.GetProperties(ctx, nil)
}

func (s ServiceClient) SetProperties(ctx context.Context, properties StorageServiceProperties) (ServiceSetPropertiesResponse, error) {
	return s.client.SetProperties(ctx, properties, nil)
}

func (s ServiceClient) GetStatistics(ctx context.Context) (StorageServiceStatsResponse, error) {
	return s.client.GetStatistics(ctx, nil)
}
