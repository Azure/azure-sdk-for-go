package advisor

import (
	 original "github.com/Azure/azure-sdk-for-go/service/advisor/management/2017-04-19/advisor"
)

type (
	 Category = original.Category
	 Impact = original.Impact
	 Risk = original.Risk
	 OperationDisplayInfo = original.OperationDisplayInfo
	 OperationEntity = original.OperationEntity
	 OperationEntityListResult = original.OperationEntityListResult
	 RecommendationProperties = original.RecommendationProperties
	 Resource = original.Resource
	 ResourceRecommendationBase = original.ResourceRecommendationBase
	 ResourceRecommendationBaseListResult = original.ResourceRecommendationBaseListResult
	 ShortDescription = original.ShortDescription
	 SuppressionContract = original.SuppressionContract
	 SuppressionContractListResult = original.SuppressionContractListResult
	 SuppressionProperties = original.SuppressionProperties
	 OperationsClient = original.OperationsClient
	 RecommendationsClient = original.RecommendationsClient
	 SuppressionsClient = original.SuppressionsClient
	 ManagementClient = original.ManagementClient
)

const (
	 Cost = original.Cost
	 HighAvailability = original.HighAvailability
	 Performance = original.Performance
	 Security = original.Security
	 High = original.High
	 Low = original.Low
	 Medium = original.Medium
	 Error = original.Error
	 None = original.None
	 Warning = original.Warning
	 DefaultBaseURI = original.DefaultBaseURI
)
