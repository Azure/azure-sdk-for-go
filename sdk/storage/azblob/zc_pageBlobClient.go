package azblob

import (
	"context"
	"io"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type PageBlobClient struct {
	BlobClient
	client *pageBlobClient
	u      url.URL
}

func NewPageBlobClient(blobURL string, cred azcore.Credential, options *connectionOptions) (PageBlobClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return PageBlobClient{}, err
	}
	con := newConnection(blobURL, cred, options)
	return PageBlobClient{client: &pageBlobClient{con: con}, u: *u}, nil
}

func (pb PageBlobClient) WithPipeline(pipeline azcore.Pipeline) PageBlobClient {
	con := newConnectionWithPipeline(pb.u.String(), pipeline)
	return PageBlobClient{client: &pageBlobClient{con}, u: pb.u}
}

func (pb PageBlobClient) URL() url.URL {
	return pb.u
}

func (pb PageBlobClient) WithSnapshot(snapshot string) PageBlobClient {
	p := NewBlobURLParts(pb.URL())
	p.Snapshot = snapshot
	snapshotURL := p.URL()
	return PageBlobClient{
		client: &pageBlobClient{
			newConnectionWithPipeline(snapshotURL.String(), pb.client.con.p),
		},
		u: pb.u,
	}
}

func (pb PageBlobClient) GetAccountInfo(ctx context.Context) (BlobGetAccountInfoResponse, error) {
	blobClient := BlobClient{client: &blobClient{pb.client.con, nil}}

	return blobClient.GetAccountInfo(ctx)
}

func (pb PageBlobClient) Create(ctx context.Context, size int64, options *CreatePageBlobOptions) (PageBlobCreateResponse, error) {
	creationOptions, httpHeaders, cpkInfo, cpkScope, lac, mac := options.pointers()

	return pb.client.Create(ctx, 0, size, creationOptions, httpHeaders, lac, cpkInfo, cpkScope, mac)
}

func (pb PageBlobClient) UploadPages(ctx context.Context, offset int, body io.ReadSeeker, options *UploadPagesOptions) (PageBlobUploadPagesResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)

	if err != nil {
		return PageBlobUploadPagesResponse{}, err
	}

	uploadOptions, cpkInfo, cpkScope, snac, lac, mac := options.pointers()

	return pb.client.UploadPages(ctx, count, azcore.NopCloser(body), uploadOptions, lac, cpkInfo, cpkScope, snac, mac)
}

func (pb PageBlobClient) UploadPagesFromURL(ctx context.Context, source url.URL, sourceOffset, blobOffset, count int64, options *UploadPagesFromURLOptions) (PageBlobUploadPagesFromURLResponse, error) {
	uploadOptions, cpkInfo, cpkScope, snac, smac, lac, mac := options.pointers()

	return pb.client.UploadPagesFromURL(ctx, source, rangeToString(sourceOffset, count), 0, rangeToString(blobOffset, count), uploadOptions, cpkInfo, cpkScope, lac, snac, mac, smac)
}

func (pb PageBlobClient) ClearPages(ctx context.Context, offset, count int64, options *ClearPagesOptions) (PageBlobClearPagesResponse, error) {
	clearOptions := &PageBlobClearPagesOptions{
		RangeParameter: rangeToStringPtr(offset, count),
	}

	cpkInfo, cpkScope, snac, lac, mac := options.pointers()

	return pb.client.ClearPages(ctx, 0, clearOptions, lac, cpkInfo, cpkScope, snac, mac)
}

func (pb PageBlobClient) GetPageRanges(ctx context.Context, offset, count int64, options *GetPageRangesOptions) (PageListResponse, error) {
	snapshot, lac, mac := options.pointers()

	getRangesOptions := &PageBlobGetPageRangesOptions{
		RangeParameter: rangeToStringPtr(offset, count),
		Snapshot:       snapshot,
	}

	return pb.client.GetPageRanges(ctx, getRangesOptions, lac, mac)
}

func (pb PageBlobClient) GetPageRangesDiff(ctx context.Context, offset, count int64, prevSnapshot string, options *GetPageRangesOptions) (PageListResponse, error) {
	snapshot, lac, mac := options.pointers()

	diffOptions := &PageBlobGetPageRangesDiffOptions{
		Prevsnapshot:   &prevSnapshot,
		RangeParameter: rangeToStringPtr(offset, count),
		Snapshot:       snapshot,
	}

	return pb.client.GetPageRangesDiff(ctx, diffOptions, lac, mac)
}

func (pb PageBlobClient) Resize(ctx context.Context, size int64, options *ResizePageBlobOptions) (PageBlobResizeResponse, error) {
	cpkInfo, cpkScope, lac, mac := options.pointers()

	return pb.client.Resize(ctx, size, nil, lac, cpkInfo, cpkScope, mac)
}

func (pb PageBlobClient) UpdateSequenceNumber(ctx context.Context, actionType SequenceNumberActionType, options *UpdateSequenceNumberPageBlob) (PageBlobUpdateSequenceNumberResponse, error) {
	updateOptions, lac, mac := options.pointers()

	return pb.client.UpdateSequenceNumber(ctx, actionType, updateOptions, lac, mac)
}

func (pb PageBlobClient) StartCopyIncremental(ctx context.Context, source url.URL, conditions *ModifiedAccessConditions) (PageBlobCopyIncrementalResponse, error) {
	return pb.client.CopyIncremental(ctx, source, nil, conditions)
}
