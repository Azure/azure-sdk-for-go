Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82//specification/cognitiveservices/data-plane/TranslatorText/readme.md tag: `release_1_0_preview.1`

Code generator @microsoft.azure/autorest.go@2.1.175


## Breaking Changes

## Struct Changes

### Removed Structs

1. SetObject

### Removed Struct Fields

1. ErrorResponseV2.autorest.Response

## Signature Changes

### Funcs

1. TranslationClient.CancelOperation
	- Returns
		- From: SetObject, error
		- To: BatchStatusDetail, error
1. TranslationClient.CancelOperationResponder
	- Returns
		- From: SetObject, error
		- To: BatchStatusDetail, error
1. TranslationClient.GetDocumentFormats
	- Returns
		- From: SetObject, error
		- To: FileFormatListResult, error
1. TranslationClient.GetDocumentFormatsResponder
	- Returns
		- From: SetObject, error
		- To: FileFormatListResult, error
1. TranslationClient.GetDocumentStatus
	- Returns
		- From: SetObject, error
		- To: DocumentStatusDetail, error
1. TranslationClient.GetDocumentStatusResponder
	- Returns
		- From: SetObject, error
		- To: DocumentStatusDetail, error
1. TranslationClient.GetDocumentStorageSource
	- Returns
		- From: SetObject, error
		- To: StorageSourceListResult, error
1. TranslationClient.GetDocumentStorageSourceResponder
	- Returns
		- From: SetObject, error
		- To: StorageSourceListResult, error
1. TranslationClient.GetGlossaryFormats
	- Returns
		- From: SetObject, error
		- To: FileFormatListResult, error
1. TranslationClient.GetGlossaryFormatsResponder
	- Returns
		- From: SetObject, error
		- To: FileFormatListResult, error
1. TranslationClient.GetOperationDocumentsStatus
	- Returns
		- From: SetObject, error
		- To: DocumentStatusResponse, error
1. TranslationClient.GetOperationDocumentsStatusResponder
	- Returns
		- From: SetObject, error
		- To: DocumentStatusResponse, error
1. TranslationClient.GetOperationStatus
	- Returns
		- From: SetObject, error
		- To: BatchStatusDetail, error
1. TranslationClient.GetOperationStatusResponder
	- Returns
		- From: SetObject, error
		- To: BatchStatusDetail, error
1. TranslationClient.GetOperations
	- Returns
		- From: SetObject, error
		- To: BatchStatusResponse, error
1. TranslationClient.GetOperationsResponder
	- Returns
		- From: SetObject, error
		- To: BatchStatusResponse, error
1. TranslationClient.SubmitBatchRequest
	- Returns
		- From: ErrorResponseV2, error
		- To: autorest.Response, error
1. TranslationClient.SubmitBatchRequestResponder
	- Returns
		- From: ErrorResponseV2, error
		- To: autorest.Response, error

## Struct Changes

### New Struct Fields

1. BatchStatusDetail.autorest.Response
1. BatchStatusResponse.autorest.Response
1. DocumentStatusDetail.autorest.Response
1. DocumentStatusResponse.autorest.Response
1. FileFormatListResult.autorest.Response
1. StorageSourceListResult.autorest.Response
