package azblob

import (
	"context"
	"io"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	AppendBlobMaxAppendBlockBytes = 4 * 1024 * 1024
	AppendBlobMaxBlocks           = 5000
)

type AppendBlobClient struct {
	BlobClient
	client *appendBlobClient
	u      url.URL
}

func NewAppendBlobClient(blobURL string, cred azcore.Credential, options *connectionOptions) (AppendBlobClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return AppendBlobClient{}, err
	}
	con := newConnection(blobURL, cred, options)
	return AppendBlobClient{client: &appendBlobClient{con: con}, u: *u}, nil
}

func (ab AppendBlobClient) WithPipeline(pipeline azcore.Pipeline) AppendBlobClient {
	con := newConnectionWithPipeline(ab.u.String(), pipeline)
	return AppendBlobClient{client: &appendBlobClient{con}, u: ab.u}
}

func (ab AppendBlobClient) GetAccountInfo(ctx context.Context) (BlobGetAccountInfoResponse, error) {
	blobClient := BlobClient{client: &blobClient{ab.client.con, nil}}

	return blobClient.GetAccountInfo(ctx)
}

func (ab AppendBlobClient) URL() url.URL {
	return ab.u
}

func (ab AppendBlobClient) WithSnapshot(snapshot string) AppendBlobClient {
	p := NewBlobURLParts(ab.URL())
	p.Snapshot = snapshot
	snapshotURL := p.URL()
	return AppendBlobClient{
		client: &appendBlobClient{
			newConnectionWithPipeline(snapshotURL.String(), ab.client.con.p),
		},
		u: ab.u,
	}
}

func (ab AppendBlobClient) AppendBlock(ctx context.Context, body io.ReadSeeker, options *AppendBlockOptions) (AppendBlobAppendBlockResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return AppendBlobAppendBlockResponse{}, nil
	}

	appendOptions, aac, cpkinfo, cpkscope, mac, lac := options.pointers()

	return ab.client.AppendBlock(ctx, count, azcore.NopCloser(body), appendOptions, lac, aac, cpkinfo, cpkscope, mac)
}

func (ab AppendBlobClient) AppendBlockFromURL(ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions) (AppendBlobAppendBlockFromURLResponse, error) {
	appendOptions, aac, cpkinfo, cpkscope, mac, lac, smac := options.pointers()

	return ab.client.AppendBlockFromURL(ctx, source, contentLength, appendOptions, cpkinfo, cpkscope, lac, aac, mac, smac)
}
