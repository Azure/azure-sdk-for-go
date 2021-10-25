# Release History

## 0.1.1 (2021-10-25)
### Breaking Changes

### New Content

- New function `*CommunityGalleryImagesClient.Get(context.Context, string, string, string, *CommunityGalleryImagesGetOptions) (CommunityGalleryImagesGetResponse, error)`
- New function `NewCommunityGalleryImagesClient(*arm.Connection, string) *CommunityGalleryImagesClient`
- New function `*CommunityGalleriesClient.Get(context.Context, string, string, *CommunityGalleriesGetOptions) (CommunityGalleriesGetResponse, error)`
- New function `NewCommunityGalleryImageVersionsClient(*arm.Connection, string) *CommunityGalleryImageVersionsClient`
- New function `NewCommunityGalleriesClient(*arm.Connection, string) *CommunityGalleriesClient`
- New function `CommunityGalleryImageProperties.MarshalJSON() ([]byte, error)`
- New function `*CommunityGalleryImageProperties.UnmarshalJSON([]byte) error`
- New function `CommunityGalleryImageVersionProperties.MarshalJSON() ([]byte, error)`
- New function `*CommunityGalleryImageVersionProperties.UnmarshalJSON([]byte) error`
- New function `*CommunityGalleryImageVersionsClient.Get(context.Context, string, string, string, string, *CommunityGalleryImageVersionsGetOptions) (CommunityGalleryImageVersionsGetResponse, error)`
- New struct `CommunityGalleriesClient`
- New struct `CommunityGalleriesGetOptions`
- New struct `CommunityGalleriesGetResponse`
- New struct `CommunityGalleriesGetResult`
- New struct `CommunityGallery`
- New struct `CommunityGalleryIdentifier`
- New struct `CommunityGalleryImage`
- New struct `CommunityGalleryImageProperties`
- New struct `CommunityGalleryImageVersion`
- New struct `CommunityGalleryImageVersionProperties`
- New struct `CommunityGalleryImageVersionsClient`
- New struct `CommunityGalleryImageVersionsGetOptions`
- New struct `CommunityGalleryImageVersionsGetResponse`
- New struct `CommunityGalleryImageVersionsGetResult`
- New struct `CommunityGalleryImagesClient`
- New struct `CommunityGalleryImagesGetOptions`
- New struct `CommunityGalleryImagesGetResponse`
- New struct `CommunityGalleryImagesGetResult`
- New struct `PirCommunityGalleryResource`

Total 0 breaking change(s), 48 additive change(s).


## 0.1.0 (2021-09-29)
- To better align with the Azure SDK guidelines (https://azure.github.io/azure-sdk/general_introduction.html), we have decided to change the module path to "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute". Therefore, we are deprecating the old module path (which is "github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute") to avoid confusion. 