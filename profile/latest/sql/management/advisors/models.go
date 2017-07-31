package advisors

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2015-05-01-preview/advisors"
)

type (
	 AutoExecuteStatus = original.AutoExecuteStatus
	 AutoExecuteStatusInheritedFrom = original.AutoExecuteStatusInheritedFrom
	 ImplementationMethod = original.ImplementationMethod
	 IsRetryable = original.IsRetryable
	 RecommendedActionCurrentState = original.RecommendedActionCurrentState
	 RecommendedActionInitiatedBy = original.RecommendedActionInitiatedBy
	 Status = original.Status
	 Advisor = original.Advisor
	 ListAdvisor = original.ListAdvisor
	 ListRecommendedAction = original.ListRecommendedAction
	 Properties = original.Properties
	 ProxyResource = original.ProxyResource
	 RecommendedAction = original.RecommendedAction
	 RecommendedActionErrorInfo = original.RecommendedActionErrorInfo
	 RecommendedActionImpactRecord = original.RecommendedActionImpactRecord
	 RecommendedActionImplementationInfo = original.RecommendedActionImplementationInfo
	 RecommendedActionMetricInfo = original.RecommendedActionMetricInfo
	 RecommendedActionProperties = original.RecommendedActionProperties
	 RecommendedActionStateInfo = original.RecommendedActionStateInfo
	 Resource = original.Resource
	 ServerAdvisorsClient = original.ServerAdvisorsClient
	 ManagementClient = original.ManagementClient
	 DatabaseAdvisorsClient = original.DatabaseAdvisorsClient
	 DatabaseRecommendedActionsClient = original.DatabaseRecommendedActionsClient
)

const (
	 Default = original.Default
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 AutoExecuteStatusInheritedFromDatabase = original.AutoExecuteStatusInheritedFromDatabase
	 AutoExecuteStatusInheritedFromDefault = original.AutoExecuteStatusInheritedFromDefault
	 AutoExecuteStatusInheritedFromElasticPool = original.AutoExecuteStatusInheritedFromElasticPool
	 AutoExecuteStatusInheritedFromServer = original.AutoExecuteStatusInheritedFromServer
	 AutoExecuteStatusInheritedFromSubscription = original.AutoExecuteStatusInheritedFromSubscription
	 AzurePowerShell = original.AzurePowerShell
	 TSQL = original.TSQL
	 No = original.No
	 Yes = original.Yes
	 Active = original.Active
	 Error = original.Error
	 Executing = original.Executing
	 Expired = original.Expired
	 Ignored = original.Ignored
	 Monitoring = original.Monitoring
	 Pending = original.Pending
	 PendingRevert = original.PendingRevert
	 Resolved = original.Resolved
	 RevertCancelled = original.RevertCancelled
	 Reverted = original.Reverted
	 Reverting = original.Reverting
	 Success = original.Success
	 Verifying = original.Verifying
	 System = original.System
	 User = original.User
	 GA = original.GA
	 LimitedPublicPreview = original.LimitedPublicPreview
	 PrivatePreview = original.PrivatePreview
	 PublicPreview = original.PublicPreview
	 DefaultBaseURI = original.DefaultBaseURI
)
