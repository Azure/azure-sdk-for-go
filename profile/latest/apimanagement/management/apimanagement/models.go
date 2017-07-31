package apimanagement

import (
	 original "github.com/Azure/azure-sdk-for-go/service/apimanagement/management/2017-03-01/apimanagement"
)

type (
	 ManagementClient = original.ManagementClient
	 PolicyScopeContract = original.PolicyScopeContract
	 ErrorFieldContract = original.ErrorFieldContract
	 ErrorResponse = original.ErrorResponse
	 PolicyCollection = original.PolicyCollection
	 PolicyContract = original.PolicyContract
	 PolicyContractProperties = original.PolicyContractProperties
	 PolicySnippetContract = original.PolicySnippetContract
	 PolicySnippetsCollection = original.PolicySnippetsCollection
	 RegionContract = original.RegionContract
	 RegionListResult = original.RegionListResult
	 Resource = original.Resource
	 PolicyClient = original.PolicyClient
	 PolicySnippetsClient = original.PolicySnippetsClient
	 RegionsClient = original.RegionsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 All = original.All
	 API = original.API
	 Operation = original.Operation
	 Product = original.Product
	 Tenant = original.Tenant
)
