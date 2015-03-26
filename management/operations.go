package management

import (
	"encoding/xml"
	"fmt"
	"time"
)

//  GetOperationStatusResponse  represents an in-flight operation. Use client.GetOperationStatus()
// to get the operation given the operation ID, or use
// client.WaitAsyncOperation() to block until the operation has completed.
// See https://msdn.microsoft.com/en-us/library/azure/ee460783.aspx
type GetOperationStatusResponse struct {
	XMLName        xml.Name `xml:"http://schemas.microsoft.com/windowsazure Operation"`
	ID             string
	Status         OperationStatus
	HttpStatusCode string
	Error          *AzureError
}

type OperationStatus string

const (
	OperationStatusInProgress OperationStatus = "InProgress"
	OperationStatusSucceeded  OperationStatus = "Succeeded"
	OperationStatusFailed     OperationStatus = "Failed"
)

// OperationId is assigned by Azure API and can be used to look up the status of an operation
type OperationId string

// getOperationStatus gets an operation given the operation ID.
func (client *Client) GetOperationStatus(operationId OperationId) (GetOperationStatusResponse, error) {
	operation := GetOperationStatusResponse{}
	if operationId == "" {
		return operation, fmt.Errorf(errParamNotSpecified, "operationId")
	}

	url := fmt.Sprintf("operations/%s", operationId)
	response, azureErr := client.SendAzureGetRequest(url)
	if azureErr != nil {
		return operation, azureErr
	}

	err := xml.Unmarshal(response, &operation)
	if err != nil {
		return operation, err
	}

	return operation, nil
}

const pollInterval = 30 * time.Second

// WaitAsyncOperation blocks until the operation with the given operationId is
// no longer in the InProgress state. If the operation was successful, nothing is
// returned, otherwise an error is returned.
func (client *Client) WaitAsyncOperation(operationId OperationId) error {
	for {
		operation, err := client.GetOperationStatus(operationId)
		if err != nil {
			return fmt.Errorf("Failed to get operation status '%s': %v", operationId, err)
		}

		switch operation.Status {
		case OperationStatusSucceeded:
			return nil
		case OperationStatusFailed:
			if operation.Error != nil {
				return operation.Error
			}
			return fmt.Errorf("Azure operation %s failed", operationId)
		case OperationStatusInProgress:
			time.Sleep(pollInterval)
		default:
			return fmt.Errorf("Unknown operation status: %s", operation.Status)
		}
	}
}
