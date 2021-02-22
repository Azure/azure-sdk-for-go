Generated from https://github.com/Azure/azure-rest-api-specs/tree/ea20b4f61fd31aeb6a72d0b0f76fdd6e68688351/specification/cognitiveservices/data-plane/Face/readme.md tag: `release_1_0`

Code generator 


### Breaking Changes

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

#### New Constants

1. DetectionModel.Detection03
