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
	client *client
}

// NewServiceClient creates a ServiceClient object using the specified URL, credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClient(serviceURL string, cred azcore.Credential, options *clientOptions) (ServiceClient, error) {
	client, err := newClient(serviceURL, cred, options)

	if err != nil {
		return ServiceClient{}, err
	}

	return ServiceClient{client: client}, err
}

// URL returns the URL endpoint used by the ServiceClient object.
func (s ServiceClient) URL() url.URL {
	return *s.client.u
}

// String returns the URL as a string.
func (s ServiceClient) String() string {
	u := s.URL()
	return u.String()
}

// WithPipeline creates a new ServiceClient object identical to the source but with the specified request policy pipeline.
func (s ServiceClient) WithPipeline(pipeline azcore.Pipeline) (ServiceClient, error) {
	client, err := newClientWithPipeline(s.client.u.String(), pipeline)

	if err != nil {
		return ServiceClient{}, err
	}

	return ServiceClient{client: client}, err
}

// NewContainerClient creates a new ContainerClient object by concatenating containerName to the end of
// ServiceClient's URL. The new ContainerClient uses the same request policy pipeline as the ServiceClient.
// To change the pipeline, create the ContainerClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewContainerClient instead of calling this object's
// NewContainerClient method.
func (s ServiceClient) NewContainerClient(containerName string) ContainerClient {
	containerURL := appendToURLPath(*s.client.u, containerName)
	containerClient, _ := newClientWithPipeline(containerURL.String(), s.client.p)
	return ContainerClient{
		client: containerClient,
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

func (s ServiceClient) GetAccountInfo(ctx context.Context) (*ServiceGetAccountInfoResponse, error) {
	return s.client.ServiceOperations().GetAccountInfo(ctx)
}

//GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceClient object.
//OAuth is required for this call, as well as any role that can delegate access to the storage account.
func (s ServiceClient) GetUserDelegationCredential(ctx context.Context, info KeyInfo) (UserDelegationCredential, error) {
	udk, err := s.client.ServiceOperations().GetUserDelegationKey(ctx, info, nil)
	if err != nil {
		return UserDelegationCredential{}, err
	}
	return NewUserDelegationCredential(strings.Split(s.String(), ".")[0], *udk.UserDelegationKey), nil
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

// The List Containers Segment operation returns a list of the containers under the specified account.
// A page is returned to allow the user to walk through segments of the list. Please refer to the examples on its ussage.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (s ServiceClient) ListContainersSegment(o *ListContainersSegmentOptions) (ListContainersSegmentResponsePager, error) {
	pager, err := s.client.ServiceOperations().ListContainersSegment(o.pointers())

	// override the generated advancer, which is incorrect
	if err == nil {
		p := pager.(*listContainersSegmentResponsePager) // cast to the internal type first

		p.advancer = func(response *ListContainersSegmentResponseResponse) (*azcore.Request, error) {
			if response.EnumerationResults.NextMarker == nil {
				return nil, errors.New("unexpected missing NextMarker")
			}

			queryValues, _ := url.ParseQuery(p.request.URL.RawQuery)
			queryValues.Set("marker", *response.EnumerationResults.NextMarker)

			p.request.URL.RawQuery = queryValues.Encode()
			return p.request, nil
		}
	}

	return pager, err
}

func (s ServiceClient) GetProperties(ctx context.Context) (*StorageServicePropertiesResponse, error) {
	return s.client.ServiceOperations().GetProperties(ctx, nil)
}

func (s ServiceClient) SetProperties(ctx context.Context, properties StorageServiceProperties) (*ServiceSetPropertiesResponse, error) {
	return s.client.ServiceOperations().SetProperties(ctx, properties, nil)
}

func (s ServiceClient) GetStatistics(ctx context.Context) (*StorageServiceStatsResponse, error) {
	return s.client.ServiceOperations().GetStatistics(ctx, nil)
}
