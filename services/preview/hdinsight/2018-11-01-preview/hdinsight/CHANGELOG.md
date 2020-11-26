
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Type of `SparkStatement.State` has been changed from `*string` to `StatementState`
- Type of `SparkStatementOutput.Status` has been changed from `*string` to `StatementExecutionStatus`
- Type of `SparkStatementRequest.Kind` has been changed from `*string` to `SessionJobKind`
- Type of `SparkBatchJob.State` has been changed from `*string` to `JobState`
- Type of `SparkJobState.State` has been changed from `*string` to `JobState`
- Type of `SparkSessionJob.State` has been changed from `*string` to `JobState`
- Type of `SparkSessionJob.Kind` has been changed from `*string` to `SessionJobKind`

## New Content

- Const `StatementStateCancelled` is added
- Const `StatementExecutionStatusOk` is added
- Const `NotStarted` is added
- Const `StatementStateCancelling` is added
- Const `StatementStateError` is added
- Const `Starting` is added
- Const `Success` is added
- Const `StatementStateRunning` is added
- Const `Idle` is added
- Const `Running` is added
- Const `Busy` is added
- Const `Error` is added
- Const `ShuttingDown` is added
- Const `Dead` is added
- Const `StatementStateWaiting` is added
- Const `StatementExecutionStatusError` is added
- Const `Recovering` is added
- Const `Killed` is added
- Const `StatementStateAvailable` is added
- Const `StatementExecutionStatusAbort` is added
- Function `PossibleStatementStateValues() []StatementState` is added
- Function `PossibleJobStateValues() []JobState` is added
- Function `PossibleStatementExecutionStatusValues() []StatementExecutionStatus` is added

