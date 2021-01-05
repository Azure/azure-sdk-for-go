package azblob

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A ContainerClient represents a URL to the Azure Storage container allowing you to manipulate its blobs.
type ContainerClient struct {
	client *containerClient
	u      url.URL
}

// NewContainerClient creates a ContainerClient object using the specified URL and request policy pipeline.
func NewContainerClient(containerURL string, cred azcore.Credential, options *connectionOptions) (ContainerClient, error) {
	u, err := url.Parse(containerURL)
	if err != nil {
		return ContainerClient{}, err
	}
	return ContainerClient{client: &containerClient{
		con: newConnection(containerURL, cred, options),
	}, u: *u}, nil
}

// URL returns the URL endpoint used by the ContainerClient object.
func (c ContainerClient) URL() url.URL {
	return c.u
}

// String returns the URL as a string.
func (c ContainerClient) String() string {
	u := c.URL()
	return u.String()
}

// WithPipeline creates a new ContainerClient object identical to the source but with the specified request policy pipeline.
func (c ContainerClient) WithPipeline(pipeline azcore.Pipeline) ContainerClient {
	con := newConnectionWithPipeline(c.u.String(), pipeline)
	return ContainerClient{client: &containerClient{con: con}, u: c.u}
}

// NewBlobClient creates a new BlobClient object by concatenating blobName to the end of
// ContainerClient's URL. The new BlobClient uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the BlobClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewBlobClient instead of calling this object's
// NewBlobClient method.
func (c ContainerClient) NewBlobURL(blobName string, mode *PathRenameMode) BlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)
	newCon := newConnectionWithPipeline(blobURL.String(), c.client.con.p)

	return BlobClient{
		client: &blobClient{newCon, mode},
		u:      blobURL,
	}
}

// NewAppendBlobURL creates a new AppendBlobURL object by concatenating blobName to the end of
// ContainerClient's URL. The new AppendBlobURL uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the AppendBlobURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewAppendBlobURL instead of calling this object's
// NewAppendBlobURL method.
func (c ContainerClient) NewAppendBlobURL(blobName string) AppendBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)
	newCon := newConnectionWithPipeline(blobURL.String(), c.client.con.p)

	return AppendBlobClient{
		client: &appendBlobClient{newCon},
		u:      blobURL,
	}
}

// NewBlockBlobClient creates a new BlockBlobClient object by concatenating blobName to the end of
// ContainerClient's URL. The new BlockBlobClient uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the BlockBlobClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewBlockBlobClient instead of calling this object's
// NewBlockBlobClient method.
func (c ContainerClient) NewBlockBlobClient(blobName string) BlockBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)
	newCon := newConnectionWithPipeline(blobURL.String(), c.client.con.p)

	return BlockBlobClient{
		client:     &blockBlobClient{newCon},
		u:          blobURL,
		BlobClient: BlobClient{client: &blobClient{con: newCon}},
	}
}

// NewPageBlobURL creates a new PageBlobURL object by concatenating blobName to the end of
// ContainerClient's URL. The new PageBlobURL uses the same request policy pipeline as the ContainerClient.
// To change the pipeline, create the PageBlobURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewPageBlobURL instead of calling this object's
// NewPageBlobURL method.
func (c ContainerClient) NewPageBlobClient(blobName string) PageBlobClient {
	blobURL := appendToURLPath(c.URL(), blobName)
	newCon := newConnectionWithPipeline(blobURL.String(), c.client.con.p)

	return PageBlobClient{
		client: &pageBlobClient{newCon},
		u:      blobURL,
	}
}

func (c ContainerClient) GetAccountInfo(ctx context.Context) (ContainerGetAccountInfoResponse, error) {
	return c.client.GetAccountInfo(ctx, nil)
}

// Create creates a new container within a storage account. If a container with the same name already exists, the operation fails.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/create-container.
func (c ContainerClient) Create(ctx context.Context, options *CreateContainerOptions) (ContainerCreateResponse, error) {
	basics, cpkInfo := options.pointers()
	return c.client.Create(ctx, basics, cpkInfo)
}

// Delete marks the specified container for deletion. The container and any blobs contained within it are later deleted during garbage collection.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-container.
func (c ContainerClient) Delete(ctx context.Context, options *DeleteContainerOptions) (ContainerDeleteResponse, error) {
	basics, leaseInfo, accessConditions := options.pointers()
	return c.client.Delete(ctx, basics, leaseInfo, accessConditions)
}

func (c ContainerClient) GetMetadata(ctx context.Context, gpo *GetPropertiesOptionsContainer) (map[string]string, error) {
	resp, err := c.GetProperties(ctx, gpo)

	if err != nil {
		return nil, err
	}

	return *resp.Metadata, nil
}

//
//// GetProperties returns the container's properties.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-container-metadata.
func (c ContainerClient) GetProperties(ctx context.Context, gpo *GetPropertiesOptionsContainer) (ContainerGetPropertiesResponse, error) {
	// NOTE: GetMetadata actually calls GetProperties internally because GetProperties returns the metadata AND the properties.
	// This allows us to not expose a GetProperties method at all simplifying the API.
	// The optionals are nil, like they were in track 1.5
	options, leaseAccess := gpo.pointers()

	return c.client.GetProperties(ctx, options, leaseAccess)
}

//// SetMetadata sets the container's metadata.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-container-metadata.
func (c ContainerClient) SetMetadata(ctx context.Context, metadata map[string]string, ac *ContainerAccessConditions) (ContainerSetMetadataResponse, error) {
	if !ac.ModifiedAccessConditions.IfUnmodifiedSince.IsZero() || *ac.ModifiedAccessConditions.IfMatch != ETagNone || *ac.ModifiedAccessConditions.IfNoneMatch != ETagNone {
		return ContainerSetMetadataResponse{}, errors.New("the IfUnmodifiedSince, IfMatch, and IfNoneMatch must have their default values because they are ignored by the blob service")
	}
	metadataOptions := ContainerSetMetadataOptions{
		Metadata: &metadata,
	}

	mac, lac := ac.pointers()

	return c.client.SetMetadata(ctx, &metadataOptions, lac, mac)
}

// GetAccessPolicy returns the container's access policy. The access policy indicates whether container's blobs may be accessed publicly.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-container-acl.
func (c ContainerClient) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (SignedIDentifierArrayResponse, error) {
	o, ac := options.pointers()

	return c.client.GetAccessPolicy(ctx, o, ac)
}

// SetAccessPolicy sets the container's permissions. The access policy indicates whether blobs in a container may be accessed publicly.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-container-acl.
func (c ContainerClient) SetAccessPolicy(ctx context.Context, options SetAccessPolicyOptions) (ContainerSetAccessPolicyResponse, error) {

	accessPolicy := options.ContainerSetAccessPolicyOptions
	if accessPolicy.Access == nil || accessPolicy.ContainerAcl == nil {
		return ContainerSetAccessPolicyResponse{}, errors.New("ContainerSetAccess must be specified with AT LEAST Access and ContainerAcl")
	}

	ac := options.ContainerAccessConditions
	if *ac.ModifiedAccessConditions.IfMatch != ETagNone || *ac.ModifiedAccessConditions.IfNoneMatch != ETagNone {
		return ContainerSetAccessPolicyResponse{}, errors.New("the IfMatch and IfNoneMatch access conditions must have their default values because they are ignored by the service")
	}

	mac, lac := ac.pointers()

	return c.client.SetAccessPolicy(ctx, &accessPolicy, lac, mac)
}

// AcquireLease acquires a lease on the container for delete operations. The lease duration must be between 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c ContainerClient) AcquireLease(ctx context.Context, leaseOptions *AcquireLeaseOptionsContainer) (ContainerAcquireLeaseResponse, error) {

	if leaseOptions == nil || leaseOptions.ContainerSetAccessPolicyOptions == nil || leaseOptions.ContainerSetAccessPolicyOptions.Duration == nil || leaseOptions.ContainerSetAccessPolicyOptions.ProposedLeaseId == nil {
		return ContainerAcquireLeaseResponse{}, errors.New("leaseOptions must be specified, with at least ProposedLeaseID and Duration specified under ContainerAcquireLeaseOptions")
	}

	return c.client.AcquireLease(ctx, leaseOptions.ContainerSetAccessPolicyOptions, leaseOptions.ModifiedAccessConditions)
}

// RenewLease renews the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c ContainerClient) RenewLease(ctx context.Context, leaseId string, leaseOptions *RenewLeaseOptionsContainer) (ContainerRenewLeaseResponse, error) {
	renewOptions, accessConditions := leaseOptions.pointers()

	return c.client.RenewLease(ctx, leaseId, renewOptions, accessConditions)
}

// ReleaseLease releases the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c ContainerClient) ReleaseLease(ctx context.Context, leaseID string, leaseOptions *ReleaseLeaseOptionsContainer) (ContainerReleaseLeaseResponse, error) {
	options, ac := leaseOptions.pointers()

	return c.client.ReleaseLease(ctx, leaseID, options, ac)
}

// BreakLease breaks the container's previously-acquired lease (if it exists).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c ContainerClient) BreakLease(ctx context.Context, container *BreakLeaseOptionsContainer) (ContainerBreakLeaseResponse, error) {
	options, ac := container.pointers()

	return c.client.BreakLease(ctx, options, ac)
}

// ChangeLease changes the container's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c ContainerClient) ChangeLease(ctx context.Context, leaseID string, proposedID string, options *ChangeLeaseOptionsContainer) (ContainerChangeLeaseResponse, error) {
	clo, ac := options.pointers()

	return c.client.ChangeLease(ctx, leaseID, proposedID, clo, ac)
}

// ListBlobsFlatSegment returns a channel of blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
// The returned channel contains individual blob items.
// AutoPagerTimeout specifies the amount of time with no read operations before the channel times out and closes. Specify no time and it will be ignored.
// AutoPagerBufferSize specifies the channel's buffer size.
func (c ContainerClient) ListBlobsFlatSegment(ctx context.Context, AutoPagerBufferSize uint, AutoPagerTimeout time.Duration, listOptions *ContainerListBlobFlatSegmentOptions) (chan BlobItemInternal, error) {
	pager := c.client.ListBlobFlatSegment(listOptions)

	output := make(chan BlobItemInternal, AutoPagerBufferSize)
	go listBlobsFlatSegmentAutoPager{
		pager,
		output,
		ctx,
		AutoPagerTimeout,
		nil,
	}.Go()

	return output, nil
}

// ListBlobsHierarchySegment returns a channel of blobs starting from the specified Marker. Use an empty
// Marker to start enumeration from the beginning. Blob names are returned in lexicographic order.
// After getting a segment, process it, and then call ListBlobsHierarchicalSegment again (passing the the
// previously-returned Marker) to get the next segment.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-blobs.
// AutoPagerTimeout specifies the amount of time with no read operations before the channel times out and closes. Specify no time and it will be ignored.
// AutoPagerBufferSize specifies the channel's buffer size.
func (c ContainerClient) ListBlobsHierarchySegment(ctx context.Context, delimiter string, AutoPagerBufferSize uint, AutoPagerTimeout time.Duration, listOptions *ContainerListBlobHierarchySegmentOptions) (chan BlobItemInternal, error) {
	pager := c.client.ListBlobHierarchySegment(delimiter, listOptions)

	output := make(chan BlobItemInternal, AutoPagerBufferSize)
	go listBlobsHierarchySegmentAutoPager{
		pager,
		output,
		ctx,
		AutoPagerTimeout,
		nil,
	}.Go()

	return output, nil
}
