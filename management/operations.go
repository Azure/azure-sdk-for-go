package management

import (
	"encoding/xml"
	"fmt"
	"time"
)

// GetOperationStatusResponse represents an in-flight operation. Use
// client.GetOperationStatus() to get the operation given the operation ID, or
// use client.WaitAsyncOperation() to block until the operation has completed.
// See https://msdn.microsoft.com/en-us/library/azure/ee460783.aspx
type GetOperationStatusResponse struct {
	XMLName        xml.Name `xml:"http://schemas.microsoft.com/windowsazure Operation"`
	ID             string
	Status         OperationStatus
	HTTPStatusCode string
	Error          *AzureError
}

// OperationStatus describes the states an Microsoft Azure Service Management
// operation an be in.
type OperationStatus string

// List of states an operation can be reported as
const (
	OperationStatusInProgress OperationStatus = "InProgress"
	OperationStatusSucceeded  OperationStatus = "Succeeded"
	OperationStatusFailed     OperationStatus = "Failed"
)

// OperationID is assigned by Azure API and can be used to look up the status of
// an operation
type OperationID string

// GetOperationStatus gets an operation given the operation ID.
func (client *Client) GetOperationStatus(operationID OperationID) (GetOperationStatusResponse, error) {
	operation := GetOperationStatusResponse{}
	if operationID == "" {
		return operation, fmt.Errorf(errParamNotSpecified, "operationID")
	}

	url := fmt.Sprintf("operations/%s", operationID)
	response, azureErr := client.SendAzureGetRequest(url)
	if azureErr != nil {
		return operation, azureErr
	}

	err := xml.Unmarshal(response, &operation)
	return operation, err
}

const pollInterval = 30 * time.Second

// WaitAsyncOperation blocks until the operation with the given operationId is
// no longer in the InProgress state. If the operation was successful, nothing
// is returned, otherwise an error is returned.
func (client *Client) WaitAsyncOperation(operationID OperationID) error {
	for {
		operation, err := client.GetOperationStatus(operationID)
		if err != nil {
			return fmt.Errorf("Failed to get operation status '%s': %v", operationID, err)
		}

		switch operation.Status {
		case OperationStatusSucceeded:
			return nil
		case OperationStatusFailed:
			if operation.Error != nil {
				return operation.Error
			}
			return fmt.Errorf("Azure Operation ID=%s has failed", operationID)
		case OperationStatusInProgress:
			time.Sleep(pollInterval)
		default:
			return fmt.Errorf("Unknown operation status:%s (ID=%s)", operation.Status, operationID)
		}
	}
}
