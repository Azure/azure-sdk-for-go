Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Type of `SparkStatementRequest.Kind` has been changed from `*string` to `SessionJobKind`
- Type of `SparkBatchJob.State` has been changed from `*string` to `JobState`
- Type of `SparkSessionJob.Kind` has been changed from `*string` to `SessionJobKind`
- Type of `SparkSessionJob.State` has been changed from `*string` to `JobState`
- Type of `SparkStatementOutput.Status` has been changed from `*string` to `StatementExecutionStatus`
- Type of `SparkJobState.State` has been changed from `*string` to `JobState`
- Type of `SparkStatement.State` has been changed from `*string` to `StatementState`

## New Content

- New const `Busy`
- New const `StatementStateAvailable`
- New const `Idle`
- New const `Success`
- New const `StatementStateError`
- New const `Recovering`
- New const `Starting`
- New const `Killed`
- New const `ShuttingDown`
- New const `Error`
- New const `StatementStateWaiting`
- New const `Dead`
- New const `NotStarted`
- New const `StatementExecutionStatusOk`
- New const `StatementExecutionStatusAbort`
- New const `Running`
- New const `StatementStateCancelled`
- New const `StatementStateCancelling`
- New const `StatementExecutionStatusError`
- New const `StatementStateRunning`
- New function `PossibleStatementStateValues() []StatementState`
- New function `PossibleJobStateValues() []JobState`
- New function `PossibleStatementExecutionStatusValues() []StatementExecutionStatus`
