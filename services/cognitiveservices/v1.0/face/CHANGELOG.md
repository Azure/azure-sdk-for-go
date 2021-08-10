# Unreleased

## Breaking Changes

### Removed Constants

1. AccessoryType.Glasses
1. AccessoryType.HeadWear
1. AccessoryType.Mask

### Signature Changes

#### Funcs

1. Client.DetectWithStream
	- Params
		- From: context.Context, io.ReadCloser, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel
		- To: context.Context, io.ReadCloser, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel, *int32
1. Client.DetectWithStreamPreparer
	- Params
		- From: context.Context, io.ReadCloser, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel
		- To: context.Context, io.ReadCloser, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel, *int32
1. Client.DetectWithURL
	- Params
		- From: context.Context, ImageURL, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel
		- To: context.Context, ImageURL, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel, *int32
1. Client.DetectWithURLPreparer
	- Params
		- From: context.Context, ImageURL, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel
		- To: context.Context, ImageURL, *bool, *bool, []AttributeType, RecognitionModel, *bool, DetectionModel, *int32
1. LargeFaceListClient.List
	- Params
		- From: context.Context, *bool
		- To: context.Context, *bool, string, *int32
1. LargeFaceListClient.ListPreparer
	- Params
		- From: context.Context, *bool
		- To: context.Context, *bool, string, *int32

## Additive Changes

### New Constants

1. AccessoryType.AccessoryTypeGlasses
1. AccessoryType.AccessoryTypeHeadWear
1. AccessoryType.AccessoryTypeMask
1. AttributeType.AttributeTypeMask
1. DetectionModel.Detection03
1. MaskType.FaceMask
1. MaskType.NoMask
1. MaskType.OtherMaskOrOcclusion
1. MaskType.Uncertain
1. RecognitionModel.Recognition04

### New Funcs

1. PossibleMaskTypeValues() []MaskType

### Struct Changes

#### New Structs

1. Mask
1. NonNullableNameAndNullableUserDataContract

#### New Struct Fields

1. Attributes.Mask
