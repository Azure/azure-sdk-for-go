package azblob

import (
	"context"
	"io"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	AppendBlobMaxAppendBlockBytes = 4 * 1024 * 1024
	AppendBlobMaxBlocks           = 50_000
)

type AppendBlobClient struct {
	BlobClient
	client *appendBlobClient
	u      url.URL
}

func NewAppendBlobClient(blobURL string, cred azcore.Credential, options *ClientOptions) (AppendBlobClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return AppendBlobClient{}, err
	}
	con := newConnection(blobURL, cred, options.getConnectionOptions())
	return AppendBlobClient{
		client:     &appendBlobClient{con: con},
		u:          *u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}, nil
}

func (ab AppendBlobClient) URL() url.URL {
	return ab.u
}

// WithSnapshot creates a new AppendBlobURL object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (ab AppendBlobClient) WithSnapshot(snapshot string) AppendBlobClient {
	p := NewBlobURLParts(ab.URL())
	p.Snapshot = snapshot
	snapshotURL := p.URL()
	con := newConnectionWithPipeline(snapshotURL.String(), ab.client.con.p)
	return AppendBlobClient{
		client:     &appendBlobClient{con: con},
		u:          snapshotURL,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// WithVersionID creates a new AppendBlobURL object identical to the source but with the specified version id.
// Pass "" to remove the versionID returning a URL to the base blob.
func (ab AppendBlobClient) WithVersionID(versionID string) AppendBlobClient {
	p := NewBlobURLParts(ab.URL())
	p.VersionID = versionID
	versionIDURL := p.URL()
	con := newConnectionWithPipeline(versionIDURL.String(), ab.client.con.p)
	return AppendBlobClient{
		client:     &appendBlobClient{con: con},
		u:          versionIDURL,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// Create creates a 0-size append blob. Call AppendBlock to append data to an append blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (ab AppendBlobClient) Create(ctx context.Context, options *CreateAppendBlobOptions) (AppendBlobCreateResponse, error) {
	appendBlobAppendBlockOptions, blobHttpHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := options.pointers()
	resp, err := ab.client.Create(ctx, 0, appendBlobAppendBlockOptions, blobHttpHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return resp, handleError(err)
}

// AppendBlock writes a stream to a new block of data to the end of the existing append blob.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/append-block.
func (ab AppendBlobClient) AppendBlock(ctx context.Context, body io.ReadSeeker, options *AppendBlockOptions) (AppendBlobAppendBlockResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return AppendBlobAppendBlockResponse{}, nil
	}

	appendOptions, aac, cpkinfo, cpkscope, mac, lac := options.pointers()

	resp, err := ab.client.AppendBlock(ctx, count, azcore.NopCloser(body), appendOptions, lac, aac, cpkinfo, cpkscope, mac)

	return resp, handleError(err)
}

// AppendBlockFromURL copies a new block of data from source URL to the end of the existing append blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/append-block-from-url.
func (ab AppendBlobClient) AppendBlockFromURL(ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions) (AppendBlobAppendBlockFromURLResponse, error) {
	appendOptions, aac, cpkinfo, cpkscope, mac, lac, smac := options.pointers()

	resp, err := ab.client.AppendBlockFromURL(ctx, source, contentLength, appendOptions, cpkinfo, cpkscope, lac, aac, mac, smac)

	return resp, handleError(err)
}
