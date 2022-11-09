# Release History

## 0.3.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

## 0.3.0 (2022-11-08)

### Features Added
* Added `ColumnIndexLookup` field to Table struct
* Added type `Row`
* Added sovereign cloud support

### Breaking Changes
* Added error return values to `NewLogsClient` and `NewMetricsClient`
* Rename `Batch` to `QueryBatch`
* Rename `NewListMetricDefinitionsPager` to `NewListDefinitionsPager`
* Rename `NewListMetricNamespacesPager` to `NewListNamespacesPager`
* Changed type of `Render` and `Statistics` from interface{} to []byte

### Other Changes
* Updated docs with more detailed examples

## 0.2.0 (2022-10-11)

### Breaking Changes
* Changed format of logs `ErrorInfo` struct to custom error type

## 0.1.0 (2022-09-08)
* This is the initial release of the `azquery` library
