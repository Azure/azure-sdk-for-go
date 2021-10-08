# Change History

## Additive Changes

### New Constants

1. DiskCreateOption.DiskCreateOptionCopyStart
1. DiskState.DiskStateActiveSASFrozen
1. DiskState.DiskStateFrozen
1. PublicNetworkAccess.PublicNetworkAccessDisabled
1. PublicNetworkAccess.PublicNetworkAccessEnabled

### New Funcs

1. *CommunityGallery.UnmarshalJSON([]byte) error
1. *CommunityGalleryImage.UnmarshalJSON([]byte) error
1. *CommunityGalleryImageVersion.UnmarshalJSON([]byte) error
1. *PirCommunityGalleryResource.UnmarshalJSON([]byte) error
1. CommunityGalleriesClient.Get(context.Context, string, string) (CommunityGallery, error)
1. CommunityGalleriesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. CommunityGalleriesClient.GetResponder(*http.Response) (CommunityGallery, error)
1. CommunityGalleriesClient.GetSender(*http.Request) (*http.Response, error)
1. CommunityGallery.MarshalJSON() ([]byte, error)
1. CommunityGalleryImage.MarshalJSON() ([]byte, error)
1. CommunityGalleryImageVersion.MarshalJSON() ([]byte, error)
1. CommunityGalleryImageVersionsClient.Get(context.Context, string, string, string, string) (CommunityGalleryImageVersion, error)
1. CommunityGalleryImageVersionsClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. CommunityGalleryImageVersionsClient.GetResponder(*http.Response) (CommunityGalleryImageVersion, error)
1. CommunityGalleryImageVersionsClient.GetSender(*http.Request) (*http.Response, error)
1. CommunityGalleryImagesClient.Get(context.Context, string, string, string) (CommunityGalleryImage, error)
1. CommunityGalleryImagesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. CommunityGalleryImagesClient.GetResponder(*http.Response) (CommunityGalleryImage, error)
1. CommunityGalleryImagesClient.GetSender(*http.Request) (*http.Response, error)
1. NewCommunityGalleriesClient(string) CommunityGalleriesClient
1. NewCommunityGalleriesClientWithBaseURI(string, string) CommunityGalleriesClient
1. NewCommunityGalleryImageVersionsClient(string) CommunityGalleryImageVersionsClient
1. NewCommunityGalleryImageVersionsClientWithBaseURI(string, string) CommunityGalleryImageVersionsClient
1. NewCommunityGalleryImagesClient(string) CommunityGalleryImagesClient
1. NewCommunityGalleryImagesClientWithBaseURI(string, string) CommunityGalleryImagesClient
1. PirCommunityGalleryResource.MarshalJSON() ([]byte, error)
1. PossiblePublicNetworkAccessValues() []PublicNetworkAccess

### Struct Changes

#### New Structs

1. CommunityGalleriesClient
1. CommunityGallery
1. CommunityGalleryIdentifier
1. CommunityGalleryImage
1. CommunityGalleryImageProperties
1. CommunityGalleryImageVersion
1. CommunityGalleryImageVersionProperties
1. CommunityGalleryImageVersionsClient
1. CommunityGalleryImagesClient
1. PirCommunityGalleryResource
1. SupportedCapabilities

#### New Struct Fields

1. DiskAccess.ExtendedLocation
1. DiskProperties.CompletionPercent
1. DiskProperties.PublicNetworkAccess
1. DiskProperties.SupportedCapabilities
1. DiskRestorePointProperties.CompletionPercent
1. DiskRestorePointProperties.DiskAccessID
1. DiskRestorePointProperties.NetworkAccessPolicy
1. DiskRestorePointProperties.PublicNetworkAccess
1. DiskRestorePointProperties.SupportedCapabilities
1. DiskUpdateProperties.PublicNetworkAccess
1. DiskUpdateProperties.SupportedCapabilities
1. EncryptionSetProperties.AutoKeyRotationError
1. SnapshotProperties.CompletionPercent
1. SnapshotProperties.PublicNetworkAccess
1. SnapshotProperties.SupportedCapabilities
1. SnapshotUpdateProperties.PublicNetworkAccess
