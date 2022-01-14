# Change History

## Additive Changes

### New Constants

1. TypeBasicActivity.TypeBasicActivityTypeFail

### New Funcs

1. *FailActivity.UnmarshalJSON([]byte) error
1. Activity.AsFailActivity() (*FailActivity, bool)
1. AppendVariableActivity.AsFailActivity() (*FailActivity, bool)
1. AzureDataExplorerCommandActivity.AsFailActivity() (*FailActivity, bool)
1. AzureFunctionActivity.AsFailActivity() (*FailActivity, bool)
1. AzureMLBatchExecutionActivity.AsFailActivity() (*FailActivity, bool)
1. AzureMLExecutePipelineActivity.AsFailActivity() (*FailActivity, bool)
1. AzureMLUpdateResourceActivity.AsFailActivity() (*FailActivity, bool)
1. ControlActivity.AsFailActivity() (*FailActivity, bool)
1. CopyActivity.AsFailActivity() (*FailActivity, bool)
1. CustomActivity.AsFailActivity() (*FailActivity, bool)
1. DataLakeAnalyticsUSQLActivity.AsFailActivity() (*FailActivity, bool)
1. DatabricksNotebookActivity.AsFailActivity() (*FailActivity, bool)
1. DatabricksSparkJarActivity.AsFailActivity() (*FailActivity, bool)
1. DatabricksSparkPythonActivity.AsFailActivity() (*FailActivity, bool)
1. DeleteActivity.AsFailActivity() (*FailActivity, bool)
1. ExecuteDataFlowActivity.AsFailActivity() (*FailActivity, bool)
1. ExecutePipelineActivity.AsFailActivity() (*FailActivity, bool)
1. ExecuteSSISPackageActivity.AsFailActivity() (*FailActivity, bool)
1. ExecuteWranglingDataflowActivity.AsFailActivity() (*FailActivity, bool)
1. ExecutionActivity.AsFailActivity() (*FailActivity, bool)
1. FailActivity.AsActivity() (*Activity, bool)
1. FailActivity.AsAppendVariableActivity() (*AppendVariableActivity, bool)
1. FailActivity.AsAzureDataExplorerCommandActivity() (*AzureDataExplorerCommandActivity, bool)
1. FailActivity.AsAzureFunctionActivity() (*AzureFunctionActivity, bool)
1. FailActivity.AsAzureMLBatchExecutionActivity() (*AzureMLBatchExecutionActivity, bool)
1. FailActivity.AsAzureMLExecutePipelineActivity() (*AzureMLExecutePipelineActivity, bool)
1. FailActivity.AsAzureMLUpdateResourceActivity() (*AzureMLUpdateResourceActivity, bool)
1. FailActivity.AsBasicActivity() (BasicActivity, bool)
1. FailActivity.AsBasicControlActivity() (BasicControlActivity, bool)
1. FailActivity.AsBasicExecutionActivity() (BasicExecutionActivity, bool)
1. FailActivity.AsControlActivity() (*ControlActivity, bool)
1. FailActivity.AsCopyActivity() (*CopyActivity, bool)
1. FailActivity.AsCustomActivity() (*CustomActivity, bool)
1. FailActivity.AsDataLakeAnalyticsUSQLActivity() (*DataLakeAnalyticsUSQLActivity, bool)
1. FailActivity.AsDatabricksNotebookActivity() (*DatabricksNotebookActivity, bool)
1. FailActivity.AsDatabricksSparkJarActivity() (*DatabricksSparkJarActivity, bool)
1. FailActivity.AsDatabricksSparkPythonActivity() (*DatabricksSparkPythonActivity, bool)
1. FailActivity.AsDeleteActivity() (*DeleteActivity, bool)
1. FailActivity.AsExecuteDataFlowActivity() (*ExecuteDataFlowActivity, bool)
1. FailActivity.AsExecutePipelineActivity() (*ExecutePipelineActivity, bool)
1. FailActivity.AsExecuteSSISPackageActivity() (*ExecuteSSISPackageActivity, bool)
1. FailActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. FailActivity.AsExecutionActivity() (*ExecutionActivity, bool)
1. FailActivity.AsFailActivity() (*FailActivity, bool)
1. FailActivity.AsFilterActivity() (*FilterActivity, bool)
1. FailActivity.AsForEachActivity() (*ForEachActivity, bool)
1. FailActivity.AsGetMetadataActivity() (*GetMetadataActivity, bool)
1. FailActivity.AsHDInsightHiveActivity() (*HDInsightHiveActivity, bool)
1. FailActivity.AsHDInsightMapReduceActivity() (*HDInsightMapReduceActivity, bool)
1. FailActivity.AsHDInsightPigActivity() (*HDInsightPigActivity, bool)
1. FailActivity.AsHDInsightSparkActivity() (*HDInsightSparkActivity, bool)
1. FailActivity.AsHDInsightStreamingActivity() (*HDInsightStreamingActivity, bool)
1. FailActivity.AsIfConditionActivity() (*IfConditionActivity, bool)
1. FailActivity.AsLookupActivity() (*LookupActivity, bool)
1. FailActivity.AsSQLServerStoredProcedureActivity() (*SQLServerStoredProcedureActivity, bool)
1. FailActivity.AsSetVariableActivity() (*SetVariableActivity, bool)
1. FailActivity.AsSwitchActivity() (*SwitchActivity, bool)
1. FailActivity.AsUntilActivity() (*UntilActivity, bool)
1. FailActivity.AsValidationActivity() (*ValidationActivity, bool)
1. FailActivity.AsWaitActivity() (*WaitActivity, bool)
1. FailActivity.AsWebActivity() (*WebActivity, bool)
1. FailActivity.AsWebHookActivity() (*WebHookActivity, bool)
1. FailActivity.MarshalJSON() ([]byte, error)
1. FilterActivity.AsFailActivity() (*FailActivity, bool)
1. ForEachActivity.AsFailActivity() (*FailActivity, bool)
1. GetMetadataActivity.AsFailActivity() (*FailActivity, bool)
1. HDInsightHiveActivity.AsFailActivity() (*FailActivity, bool)
1. HDInsightMapReduceActivity.AsFailActivity() (*FailActivity, bool)
1. HDInsightPigActivity.AsFailActivity() (*FailActivity, bool)
1. HDInsightSparkActivity.AsFailActivity() (*FailActivity, bool)
1. HDInsightStreamingActivity.AsFailActivity() (*FailActivity, bool)
1. IfConditionActivity.AsFailActivity() (*FailActivity, bool)
1. LookupActivity.AsFailActivity() (*FailActivity, bool)
1. SQLServerStoredProcedureActivity.AsFailActivity() (*FailActivity, bool)
1. SetVariableActivity.AsFailActivity() (*FailActivity, bool)
1. SwitchActivity.AsFailActivity() (*FailActivity, bool)
1. UntilActivity.AsFailActivity() (*FailActivity, bool)
1. ValidationActivity.AsFailActivity() (*FailActivity, bool)
1. WaitActivity.AsFailActivity() (*FailActivity, bool)
1. WebActivity.AsFailActivity() (*FailActivity, bool)
1. WebHookActivity.AsFailActivity() (*FailActivity, bool)

### Struct Changes

#### New Structs

1. FailActivity
1. FailActivityTypeProperties

#### New Struct Fields

1. AzureBlobFSLinkedServiceTypeProperties.ServicePrincipalCredential
1. AzureBlobFSLinkedServiceTypeProperties.ServicePrincipalCredentialType
1. AzureDatabricksDetltaLakeLinkedServiceTypeProperties.Credential
1. AzureDatabricksDetltaLakeLinkedServiceTypeProperties.WorkspaceResourceID
1. CosmosDbLinkedServiceTypeProperties.Credential
1. DynamicsLinkedServiceTypeProperties.Credential
1. GoogleAdWordsLinkedServiceTypeProperties.ConnectionProperties
1. LinkedIntegrationRuntimeRbacAuthorization.Credential
