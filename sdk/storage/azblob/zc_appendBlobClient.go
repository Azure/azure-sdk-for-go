package azblob

import (
	"context"
	"io"
	"net/url"

	"github.com/azure/azure-sdk-for-go/sdk/azcore"
)

const (
	AppendBlobMaxAppendBlockBytes = 4 * 1024 * 1024
	AppendBlobMaxBlocks = 5000
)

type AppendBlobClient struct {
	client appendBlobClient
	cred azcore.Credential
	options *clientOptions
}

func NewAppendBlobClient(blobURL string, cred azcore.Credential, options *clientOptions) AppendBlobClient {
	client := newClient(blobURL, cred, options)

	return AppendBlobClient{
		client: appendBlobClient{client },
		cred: cred,
		options: options,
	}
}

func (ab AppendBlobClient) WithPipeline(pipeline azcore.Pipeline) AppendBlobClient {
	client := newClientWithPipeline(ab.client.u, pipeline)
	blobClient := appendBlobClient{ client }

	return AppendBlobClient{ blobClient, ab.cred, ab.options }
}

func (ab AppendBlobClient) GetAccountInfo(ctx context.Context) (*BlobGetAccountInfoResponse, error) {
	blobClient := BlobClient{client: &blobClient{ab.client.client, nil} }

	return blobClient.GetAccountInfo(ctx)
}

func (ab AppendBlobClient) URL() url.URL {
	bURL, _ := url.Parse(ab.client.u)

	return *bURL
}

func (ab AppendBlobClient) WithSnapshot(snapshot string) AppendBlobClient {
	p := NewBlobURLParts(ab.URL())
	p.Snapshot = snapshot

	uri := p.URL()

	return NewAppendBlobClient(uri.String(), ab.cred, ab.options)
}

func (ab AppendBlobClient) AppendBlock(ctx context.Context, body io.ReadSeeker, options *AppendBlockOptions) (*AppendBlobAppendBlockResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return nil, nil
	}

	appendOptions, aac, cpkinfo, cpkscope, mac, lac := options.pointers()

	return ab.client.AppendBlock(ctx, count, azcore.NopCloser(body), appendOptions, lac, aac, cpkinfo, cpkscope, mac)
}

func (ab AppendBlobClient) AppendBlockFromURL(ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions) (*AppendBlobAppendBlockFromURLResponse, error) {
	appendOptions, aac, cpkinfo, cpkscope, mac, lac, smac := options.pointers()

	return ab.client.AppendBlockFromURL(ctx, source, contentLength, appendOptions, cpkinfo, cpkscope, lac, aac, mac, smac)
}