# Troubleshooting Azure Monitor Ingestion module

This troubleshooting guide contains instructions to diagnose frequently encountered issues while using the Azure Monitor Ingestion module for Go.

## Table of contents

* [General troubleshooting](#general-troubleshooting)
    * [Error handling](#error-handling)
    * [Enable client logging](#enable-client-logging)
* [Troubleshooting logs ingestion](#troubleshooting-logs-ingestion)
    * [Troubleshooting authorization errors](#troubleshooting-authorization-errors)
    * [Troubleshooting missing logs](#troubleshooting-missing-logs)
    * [Troubleshooting slow logs upload](#troubleshooting-slow-logs-upload)

## General troubleshooting

### Error handling

All methods which send HTTP requests return `*azcore.ResponseError` when these requests fail. `ResponseError` has error details and the raw response from Monitor Ingestion.

### Enable client logging

To troubleshoot issues with the module, first enable logging to monitor the behavior of the application. The errors and warnings in the logs generally provide useful insights into what went wrong and sometimes include corrective actions to fix issues.

This module uses the logging implementation in `azcore`. To turn on logging for all Azure SDK modules, set `AZURE_SDK_GO_LOGGING` to `all`. By default, the logger writes to stderr. Use the `azcore/log` package to control log output. For example, logging only HTTP request and response events, and printing them to stdout:

```go
import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

// print log output to stdout
azlog.SetListener(func(event azlog.Event, s string) {
    fmt.Println(s)
})

// Includes only requests and responses in logs
azlog.SetEvents(azlog.EventRequest, azlog.EventResponse)
```

## Troubleshooting logs ingestion

### Troubleshooting authorization errors

If you get an error with HTTP status code 403 (Forbidden), it means the provided credentials lack sufficient permissions to upload logs to the specified Data Collection Endpoint (DCE) and Data Collection Rule (DCR) ID.

```text
"error": {
"code": "OperationFailed",
"message": "The authentication token provided does not have access to ingest data for the data collection rule with immutable Id '***' PipelineAccessResult: AccessGranted: False, IsRbacPresent: False, IsDcrDceBindingValid: , DcrArmId: /subscriptions/***/resourceGroups/***/providers/Microsoft.Insights/dataCollectionRules/az-dcr, Message: Required authorization action was not found for tenantId *** objectId *** on resourceId /subscriptions/***/resourceGroups/***/providers/Microsoft.Insights/dataCollectionRules/az-dcr ConfigurationId: ***.."
}
```

1. Check that the application or user making the request has sufficient permissions:
   * See this document to [manage access to DCR](https://learn.microsoft.com/azure/azure-monitor/logs/tutorial-logs-ingestion-portal#assign-permissions-to-the-dcr).
   * To ingest logs, ensure the service principal is assigned the **Monitoring Metrics Publisher** role for the DCR.
1. If the user or application is granted sufficient privileges to upload logs, ensure you're authenticating as that user/application. If you're authenticating using the [DefaultAzureCredential](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential), check the logs to verify that the credential used is the one you expected. To enable logging, see the [Enable client logging](#enable-client-logging) section.
1. The permissions may take up to 30 minutes to propagate. If the permissions were granted recently, retry after some time.

### Troubleshooting missing logs

When you send logs to Azure Monitor for ingestion, the request may succeed, but you may not see the data appear in the designated Log Analytics workspace table as configured in the DCR. To investigate and resolve this issue, ensure the following:

* Double-check that you're using the correct DCE when configuring the `Client`. Using the wrong endpoint can result in data not being properly sent to the Log Analytics workspace.

* Make sure you provide the correct DCR ID to the `Upload` method. The DCR ID is an immutable identifier that determines the transformation rules applied to the uploaded logs and directs them to the appropriate Log Analytics workspace table.

* Verify that the custom table specified in the DCR exists in the Log Analytics workspace. Ensure that you provide the accurate name of the custom table to the upload method. Mismatched table names can lead to logs not being stored correctly.

* Confirm that the logs you're sending adhere to the format expected by the DCR. The data should be in the form of a JSON object or array, structured according to the requirements specified in the DCR. Additionally, it's essential to encode the request body in UTF-8 to avoid any data transmission issues.

* Keep in mind that data ingestion may take some time, especially if you're sending data to a specific table for the first time. In such cases, allow up to 15 minutes for the data to be fully ingested and available for querying and analysis.

### Troubleshooting slow logs upload

If you experience delays when uploading logs, it could be due to reaching service limits, which may trigger the rate limiter to throttle your client. To determine if your client has reached service limits, you can enable logging and check if the service is returning errors with an HTTP status code 429. For more information on service limits, see the [Azure Monitor service limits documentation](https://learn.microsoft.com/azure/azure-monitor/service-limits#logs-ingestion-api).

To enable client logging and to troubleshoot this issue further, see the instructions provided in the section titled [Enable client logging](#enable-client-logging).