package queue

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/shared"
	"net/http"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure Queue Storage service allowing you to manipulate queues.
type Client base.Client[generated.ServiceClient]

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.queue.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewServiceClient(serviceURL, pl, nil)), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.queue.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewServiceClient(serviceURL, pl, cred)), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (s *Client) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return s.generated().Endpoint()
}

// TODO: CreateQueue
// TODO: DeleteQueue

// GetProperties - gets the properties of a storage account's Queue service, including properties for Storage Analytics
// and CORS (Cross-Origin Resource Sharing) rules.
func (s *Client) GetProperties(ctx context.Context, o *GetPropertiesOptions) (GetPropertiesResponse, error) {
	getPropertiesOptions := o.format()
	resp, err := s.generated().GetProperties(ctx, getPropertiesOptions)
	return resp, err
}

// SetProperties Sets the properties of a storage account's Queue service, including Azure Storage Analytics.
// If an element (e.g. analytics_logging) is left as None, the existing settings on the service for that functionality are preserved.
func (s *Client) SetProperties(ctx context.Context, o *SetPropertiesOptions) (SetPropertiesResponse, error) {
	properties, setPropertiesOptions := o.format()
	resp, err := s.generated().SetProperties(ctx, properties, setPropertiesOptions)
	return resp, err
}

// GetStatistics Retrieves statistics related to replication for the Queue service.
func (s *Client) GetStatistics(ctx context.Context, o *GetStatisticsOptions) (GetStatisticsResponse, error) {
	getStatisticsOptions := o.format()
	resp, err := s.generated().GetStatistics(ctx, getStatisticsOptions)

	return resp, err
}

// NewListQueuesPager operation returns a pager of the queues under the specified account.
// Use an empty Marker to start enumeration from the beginning. Queue names are returned in lexicographic order.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-queues1.
func (s *Client) NewListQueuesPager(o *ListQueuesOptions) *runtime.Pager[ListQueuesResponse] {
	listOptions := generated.ServiceClientListQueuesSegmentOptions{}
	if o != nil {
		if o.Include.Metadata {
			listOptions.Include = append(listOptions.Include, "metadata")
		}
		listOptions.Marker = o.Marker
		listOptions.Maxresults = o.MaxResults
		listOptions.Prefix = o.Prefix
	}
	return runtime.NewPager(runtime.PagingHandler[ListQueuesResponse]{
		More: func(page ListQueuesResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListQueuesResponse) (ListQueuesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = s.generated().ListQueuesSegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.NextMarker
				req, err = s.generated().ListQueuesSegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListQueuesResponse{}, err
			}
			resp, err := s.generated().Pipeline().Do(req)
			if err != nil {
				return ListQueuesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListQueuesResponse{}, runtime.NewResponseError(resp)
			}
			return s.generated().ListQueuesSegmentHandleResponse(resp)
		},
	})
}

// TODO: GetSASURL()
