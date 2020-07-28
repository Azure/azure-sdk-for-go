package azblob

import (
	"context"
	"net/url"

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

////GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceClient object.
////OAuth is required for this call, as well as any role that can delegate access to the storage account.
//func (s ServiceClient) GetUserDelegationCredential(ctx context.Context, info KeyInfo, timeout *int32, requestID *string) (UserDelegationCredential, error) {
//	sc := newServiceClient(s.client.url, s.client.p)
//	udk, err := sc.GetUserDelegationKey(ctx, info, timeout, requestID)
//	if err != nil {
//		return UserDelegationCredential{}, err
//	}
//	return NewUserDelegationCredential(strings.Split(s.client.url.Host, ".")[0], *udk), nil
//}
//
////TODO this was supposed to be generated
////NewKeyInfo creates a new KeyInfo struct with the correct time formatting & conversion
//func NewKeyInfo(Start, Expiry time.Time) KeyInfo {
//	return KeyInfo{
//		Start:  Start.UTC().Format(SASTimeFormat),
//		Expiry: Expiry.UTC().Format(SASTimeFormat),
//	}
//}
//
//// ListContainersFlatSegment returns a single segment of containers starting from the specified Marker. Use an empty
//// Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
//// After getting a segment, process it, and then call ListContainersFlatSegment again (passing the the
//// previously-returned Marker) to get the next segment. For more information, see
//// https://docs.microsoft.com/rest/api/storageservices/list-containers2.
//func (s ServiceClient) ListContainersSegment(ctx context.Context, marker Marker, o ListContainersSegmentOptions) (*ListContainersSegmentResponse, error) {
//	prefix, include, maxResults := o.pointers()
//	return s.client.ListContainersSegment(ctx, prefix, marker.Val, maxResults, include, nil, nil)
//}
//
//// ListContainersOptions defines options available when calling ListContainers.
//type ListContainersSegmentOptions struct {
//	Detail     ListContainersDetail // No IncludeType header is produced if ""
//	Prefix     string               // No Prefix header is produced if ""
//	MaxResults int32                // 0 means unspecified
//	// TODO: update swagger to generate this type?
//}
//
//func (o *ListContainersSegmentOptions) pointers() (prefix *string, include ListContainersIncludeType, maxResults *int32) {
//	if o.Prefix != "" {
//		prefix = &o.Prefix
//	}
//	if o.MaxResults != 0 {
//		maxResults = &o.MaxResults
//	}
//	include = ListContainersIncludeType(o.Detail.string())
//	return
//}
//
//// ListContainersFlatDetail indicates what additional information the service should return with each container.
//type ListContainersDetail struct {
//	// Tells the service whether to return metadata for each container.
//	Metadata bool
//}
//
//// string produces the Include query parameter's value.
//func (d *ListContainersDetail) string() string {
//	items := make([]string, 0, 1)
//	// NOTE: Multiple strings MUST be appended in alphabetic order or signing the string for authentication fails!
//	if d.Metadata {
//		items = append(items, string(ListContainersIncludeMetadata))
//	}
//	if len(items) > 0 {
//		return strings.Join(items, ",")
//	}
//	return string(ListContainersIncludeNone)
//}
//
//func (s ServiceClient) GetProperties(ctx context.Context) (*StorageServiceProperties, error) {
//	return s.client.GetProperties(ctx, nil, nil)
//}
//
//func (s ServiceClient) SetProperties(ctx context.Context, properties StorageServiceProperties) (*ServiceSetPropertiesResponse, error) {
//	return s.client.SetProperties(ctx, properties, nil, nil)
//}
//
//func (s ServiceClient) GetStatistics(ctx context.Context) (*StorageServiceStats, error) {
//	return s.client.GetStatistics(ctx, nil, nil)
//}
