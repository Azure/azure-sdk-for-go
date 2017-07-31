package searchindex

import (
	 original "github.com/Azure/azure-sdk-for-go/service/search/2016-09-01/searchindex"
)

type (
	 ManagementClient = original.ManagementClient
	 DocumentsProxyClient = original.DocumentsProxyClient
	 IndexActionType = original.IndexActionType
	 QueryType = original.QueryType
	 SearchMode = original.SearchMode
	 DocumentIndexResult = original.DocumentIndexResult
	 IndexingResult = original.IndexingResult
	 Int64 = original.Int64
	 SearchParametersPayload = original.SearchParametersPayload
	 SuggestParametersPayload = original.SuggestParametersPayload
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Delete = original.Delete
	 Merge = original.Merge
	 MergeOrUpload = original.MergeOrUpload
	 Upload = original.Upload
	 Full = original.Full
	 Simple = original.Simple
	 All = original.All
	 Any = original.Any
)
