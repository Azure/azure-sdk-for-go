package shared

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

// FormatContainerAccessConditions formats FilesystemAccessConditions into container's LeaseAccessConditions and ModifiedAccessConditions.
func FormatContainerAccessConditions(b *filesystem.AccessConditions) *container.AccessConditions {
	if b == nil {
		return nil
	}
	return &container.AccessConditions{
		LeaseAccessConditions: &container.LeaseAccessConditions{
			LeaseID: b.LeaseAccessConditions.LeaseID,
		},
		ModifiedAccessConditions: &container.ModifiedAccessConditions{
			IfMatch:           b.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       b.ModifiedAccessConditions.IfNoneMatch,
			IfModifiedSince:   b.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: b.ModifiedAccessConditions.IfUnmodifiedSince,
		},
	}
}

// FormatPathAccessConditions formats PathAccessConditions into path's LeaseAccessConditions and ModifiedAccessConditions.
func FormatPathAccessConditions(p *path.AccessConditions) (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if p == nil {
		return nil, nil
	}
	return &generated.LeaseAccessConditions{
			LeaseID: p.LeaseAccessConditions.LeaseID,
		}, &generated.ModifiedAccessConditions{
			IfMatch:           p.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       p.ModifiedAccessConditions.IfNoneMatch,
			IfModifiedSince:   p.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: p.ModifiedAccessConditions.IfUnmodifiedSince,
		}
}

// FormatBlobAccessConditions formats PathAccessConditions into blob's LeaseAccessConditions and ModifiedAccessConditions.
func FormatBlobAccessConditions(p *path.AccessConditions) *blob.AccessConditions {
	if p == nil {
		return nil
	}
	return &blob.AccessConditions{LeaseAccessConditions: &blob.LeaseAccessConditions{
		LeaseID: p.LeaseAccessConditions.LeaseID,
	}, ModifiedAccessConditions: &blob.ModifiedAccessConditions{
		IfMatch:           p.ModifiedAccessConditions.IfMatch,
		IfNoneMatch:       p.ModifiedAccessConditions.IfNoneMatch,
		IfModifiedSince:   p.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: p.ModifiedAccessConditions.IfUnmodifiedSince,
	}}
}
