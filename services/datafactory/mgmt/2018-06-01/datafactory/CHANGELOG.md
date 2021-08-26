# Change History

## Additive Changes

### New Constants

1. TypeBasicActivity.TypeBasicActivityTypeExecuteWranglingDataflow
1. TypeBasicDataFlow.TypeBasicDataFlowTypeWranglingDataFlow

### New Funcs

1. *ExecuteWranglingDataflowActivity.UnmarshalJSON([]byte) error
1. *WranglingDataFlow.UnmarshalJSON([]byte) error
1. Activity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AppendVariableActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AzureDataExplorerCommandActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AzureFunctionActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AzureMLBatchExecutionActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AzureMLExecutePipelineActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. AzureMLUpdateResourceActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ControlActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. CopyActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. CustomActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. DataFlow.AsWranglingDataFlow() (*WranglingDataFlow, bool)
1. DataLakeAnalyticsUSQLActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. DatabricksNotebookActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. DatabricksSparkJarActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. DatabricksSparkPythonActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. DeleteActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ExecuteDataFlowActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ExecutePipelineActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ExecutePowerQueryActivityTypeProperties.MarshalJSON() ([]byte, error)
1. ExecuteSSISPackageActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ExecuteWranglingDataflowActivity.AsActivity() (*Activity, bool)
1. ExecuteWranglingDataflowActivity.AsAppendVariableActivity() (*AppendVariableActivity, bool)
1. ExecuteWranglingDataflowActivity.AsAzureDataExplorerCommandActivity() (*AzureDataExplorerCommandActivity, bool)
1. ExecuteWranglingDataflowActivity.AsAzureFunctionActivity() (*AzureFunctionActivity, bool)
1. ExecuteWranglingDataflowActivity.AsAzureMLBatchExecutionActivity() (*AzureMLBatchExecutionActivity, bool)
1. ExecuteWranglingDataflowActivity.AsAzureMLExecutePipelineActivity() (*AzureMLExecutePipelineActivity, bool)
1. ExecuteWranglingDataflowActivity.AsAzureMLUpdateResourceActivity() (*AzureMLUpdateResourceActivity, bool)
1. ExecuteWranglingDataflowActivity.AsBasicActivity() (BasicActivity, bool)
1. ExecuteWranglingDataflowActivity.AsBasicControlActivity() (BasicControlActivity, bool)
1. ExecuteWranglingDataflowActivity.AsBasicExecutionActivity() (BasicExecutionActivity, bool)
1. ExecuteWranglingDataflowActivity.AsControlActivity() (*ControlActivity, bool)
1. ExecuteWranglingDataflowActivity.AsCopyActivity() (*CopyActivity, bool)
1. ExecuteWranglingDataflowActivity.AsCustomActivity() (*CustomActivity, bool)
1. ExecuteWranglingDataflowActivity.AsDataLakeAnalyticsUSQLActivity() (*DataLakeAnalyticsUSQLActivity, bool)
1. ExecuteWranglingDataflowActivity.AsDatabricksNotebookActivity() (*DatabricksNotebookActivity, bool)
1. ExecuteWranglingDataflowActivity.AsDatabricksSparkJarActivity() (*DatabricksSparkJarActivity, bool)
1. ExecuteWranglingDataflowActivity.AsDatabricksSparkPythonActivity() (*DatabricksSparkPythonActivity, bool)
1. ExecuteWranglingDataflowActivity.AsDeleteActivity() (*DeleteActivity, bool)
1. ExecuteWranglingDataflowActivity.AsExecuteDataFlowActivity() (*ExecuteDataFlowActivity, bool)
1. ExecuteWranglingDataflowActivity.AsExecutePipelineActivity() (*ExecutePipelineActivity, bool)
1. ExecuteWranglingDataflowActivity.AsExecuteSSISPackageActivity() (*ExecuteSSISPackageActivity, bool)
1. ExecuteWranglingDataflowActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ExecuteWranglingDataflowActivity.AsExecutionActivity() (*ExecutionActivity, bool)
1. ExecuteWranglingDataflowActivity.AsFilterActivity() (*FilterActivity, bool)
1. ExecuteWranglingDataflowActivity.AsForEachActivity() (*ForEachActivity, bool)
1. ExecuteWranglingDataflowActivity.AsGetMetadataActivity() (*GetMetadataActivity, bool)
1. ExecuteWranglingDataflowActivity.AsHDInsightHiveActivity() (*HDInsightHiveActivity, bool)
1. ExecuteWranglingDataflowActivity.AsHDInsightMapReduceActivity() (*HDInsightMapReduceActivity, bool)
1. ExecuteWranglingDataflowActivity.AsHDInsightPigActivity() (*HDInsightPigActivity, bool)
1. ExecuteWranglingDataflowActivity.AsHDInsightSparkActivity() (*HDInsightSparkActivity, bool)
1. ExecuteWranglingDataflowActivity.AsHDInsightStreamingActivity() (*HDInsightStreamingActivity, bool)
1. ExecuteWranglingDataflowActivity.AsIfConditionActivity() (*IfConditionActivity, bool)
1. ExecuteWranglingDataflowActivity.AsLookupActivity() (*LookupActivity, bool)
1. ExecuteWranglingDataflowActivity.AsSQLServerStoredProcedureActivity() (*SQLServerStoredProcedureActivity, bool)
1. ExecuteWranglingDataflowActivity.AsSetVariableActivity() (*SetVariableActivity, bool)
1. ExecuteWranglingDataflowActivity.AsSwitchActivity() (*SwitchActivity, bool)
1. ExecuteWranglingDataflowActivity.AsUntilActivity() (*UntilActivity, bool)
1. ExecuteWranglingDataflowActivity.AsValidationActivity() (*ValidationActivity, bool)
1. ExecuteWranglingDataflowActivity.AsWaitActivity() (*WaitActivity, bool)
1. ExecuteWranglingDataflowActivity.AsWebActivity() (*WebActivity, bool)
1. ExecuteWranglingDataflowActivity.AsWebHookActivity() (*WebHookActivity, bool)
1. ExecuteWranglingDataflowActivity.MarshalJSON() ([]byte, error)
1. ExecutionActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. FilterActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ForEachActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. GetMetadataActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. HDInsightHiveActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. HDInsightMapReduceActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. HDInsightPigActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. HDInsightSparkActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. HDInsightStreamingActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. IfConditionActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. LookupActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. MappingDataFlow.AsWranglingDataFlow() (*WranglingDataFlow, bool)
1. SQLServerStoredProcedureActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. SetVariableActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. SwitchActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. UntilActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. ValidationActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. WaitActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. WebActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. WebHookActivity.AsExecuteWranglingDataflowActivity() (*ExecuteWranglingDataflowActivity, bool)
1. WranglingDataFlow.AsBasicDataFlow() (BasicDataFlow, bool)
1. WranglingDataFlow.AsDataFlow() (*DataFlow, bool)
1. WranglingDataFlow.AsMappingDataFlow() (*MappingDataFlow, bool)
1. WranglingDataFlow.AsWranglingDataFlow() (*WranglingDataFlow, bool)
1. WranglingDataFlow.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ExecutePowerQueryActivityTypeProperties
1. ExecuteWranglingDataflowActivity
1. PowerQuerySink
1. PowerQuerySource
1. PowerQueryTypeProperties
1. WranglingDataFlow
