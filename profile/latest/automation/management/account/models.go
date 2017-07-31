package account

import (
	 original "github.com/Azure/azure-sdk-for-go/service/automation/management/2015-10-31/account"
)

type (
	 OperationsClient = original.OperationsClient
	 StatisticsClient = original.StatisticsClient
	 UsagesClient = original.UsagesClient
	 AutomationAccountClient = original.AutomationAccountClient
	 ManagementClient = original.ManagementClient
	 AutomationAccountState = original.AutomationAccountState
	 ContentSourceType = original.ContentSourceType
	 DscConfigurationProvisioningState = original.DscConfigurationProvisioningState
	 DscConfigurationState = original.DscConfigurationState
	 ModuleProvisioningState = original.ModuleProvisioningState
	 RunbookProvisioningState = original.RunbookProvisioningState
	 RunbookState = original.RunbookState
	 RunbookTypeEnum = original.RunbookTypeEnum
	 SkuNameEnum = original.SkuNameEnum
	 AutomationAccount = original.AutomationAccount
	 AutomationAccountCreateOrUpdateParameters = original.AutomationAccountCreateOrUpdateParameters
	 AutomationAccountCreateOrUpdateProperties = original.AutomationAccountCreateOrUpdateProperties
	 AutomationAccountListResult = original.AutomationAccountListResult
	 AutomationAccountProperties = original.AutomationAccountProperties
	 AutomationAccountUpdateParameters = original.AutomationAccountUpdateParameters
	 AutomationAccountUpdateProperties = original.AutomationAccountUpdateProperties
	 ContentHash = original.ContentHash
	 ContentLink = original.ContentLink
	 ContentSource = original.ContentSource
	 DscConfiguration = original.DscConfiguration
	 DscConfigurationParameter = original.DscConfigurationParameter
	 DscConfigurationProperties = original.DscConfigurationProperties
	 DscNode = original.DscNode
	 DscNodeConfigurationAssociationProperty = original.DscNodeConfigurationAssociationProperty
	 ErrorResponse = original.ErrorResponse
	 Module = original.Module
	 ModuleErrorInfo = original.ModuleErrorInfo
	 ModuleProperties = original.ModuleProperties
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 Resource = original.Resource
	 Runbook = original.Runbook
	 RunbookDraft = original.RunbookDraft
	 RunbookParameter = original.RunbookParameter
	 RunbookProperties = original.RunbookProperties
	 Sku = original.Sku
	 Statistics = original.Statistics
	 StatisticsListResult = original.StatisticsListResult
	 Usage = original.Usage
	 UsageCounterName = original.UsageCounterName
	 UsageListResult = original.UsageListResult
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Ok = original.Ok
	 Suspended = original.Suspended
	 Unavailable = original.Unavailable
	 EmbeddedContent = original.EmbeddedContent
	 URI = original.URI
	 Succeeded = original.Succeeded
	 DscConfigurationStateEdit = original.DscConfigurationStateEdit
	 DscConfigurationStateNew = original.DscConfigurationStateNew
	 DscConfigurationStatePublished = original.DscConfigurationStatePublished
	 ModuleProvisioningStateActivitiesStored = original.ModuleProvisioningStateActivitiesStored
	 ModuleProvisioningStateCancelled = original.ModuleProvisioningStateCancelled
	 ModuleProvisioningStateConnectionTypeImported = original.ModuleProvisioningStateConnectionTypeImported
	 ModuleProvisioningStateContentDownloaded = original.ModuleProvisioningStateContentDownloaded
	 ModuleProvisioningStateContentRetrieved = original.ModuleProvisioningStateContentRetrieved
	 ModuleProvisioningStateContentStored = original.ModuleProvisioningStateContentStored
	 ModuleProvisioningStateContentValidated = original.ModuleProvisioningStateContentValidated
	 ModuleProvisioningStateCreated = original.ModuleProvisioningStateCreated
	 ModuleProvisioningStateCreating = original.ModuleProvisioningStateCreating
	 ModuleProvisioningStateFailed = original.ModuleProvisioningStateFailed
	 ModuleProvisioningStateModuleDataStored = original.ModuleProvisioningStateModuleDataStored
	 ModuleProvisioningStateModuleImportRunbookComplete = original.ModuleProvisioningStateModuleImportRunbookComplete
	 ModuleProvisioningStateRunningImportModuleRunbook = original.ModuleProvisioningStateRunningImportModuleRunbook
	 ModuleProvisioningStateStartingImportModuleRunbook = original.ModuleProvisioningStateStartingImportModuleRunbook
	 ModuleProvisioningStateSucceeded = original.ModuleProvisioningStateSucceeded
	 ModuleProvisioningStateUpdating = original.ModuleProvisioningStateUpdating
	 RunbookProvisioningStateSucceeded = original.RunbookProvisioningStateSucceeded
	 RunbookStateEdit = original.RunbookStateEdit
	 RunbookStateNew = original.RunbookStateNew
	 RunbookStatePublished = original.RunbookStatePublished
	 Graph = original.Graph
	 GraphPowerShell = original.GraphPowerShell
	 GraphPowerShellWorkflow = original.GraphPowerShellWorkflow
	 PowerShell = original.PowerShell
	 PowerShellWorkflow = original.PowerShellWorkflow
	 Script = original.Script
	 Basic = original.Basic
	 Free = original.Free
)
