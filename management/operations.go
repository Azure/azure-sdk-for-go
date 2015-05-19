package management

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrOperationCancelled from WaitForOperation when the polling loop is
	// cancelled through signaling the channel.
	ErrOperationCancelled = errors.New("Polling for operation status cancelled")
)

// GetOperationStatusResponse represents an in-flight operation. Use
// client.GetOperationStatus() to get the operation given the operation ID, or
// use WaitForOperation() to poll and wait until the operation has completed.
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

// GetOperationStatus gets the status of operation with given Operation ID.
// WaitForOperation utility method can be used for polling for operation status.
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

// WaitForOperation polls the Azure API for given operation ID indefinitely
// until the operation is completed with either success or failure.
// It is meant to be used for waiting for the result of the methods that
// return an OperationID value (meaning a long running operation has started).
//
// Cancellation of the polling loop (for instance, timing out) is done through
// cancel channel. If the user does not want to cancel, a nil chan can be provided.
// To cancel the method, it is recommended to close the channel provided to this
// method.
//
// If the operation was not successful or cancelling is signaled, an error
// is returned.
func (client *Client) WaitForOperation(operationID OperationID, cancel chan struct{}) error {
	for {
		done, err := client.checkOperationStatus(operationID)
		if err != nil || done {
			return err
		}
		select {
		case <-time.After(client.config.OperationPollInterval):
		case <-cancel:
			return ErrOperationCancelled
		}
	}
}

func (client *Client) checkOperationStatus(id OperationID) (done bool, err error) {
	op, err := client.GetOperationStatus(id)
	if err != nil {
		return false, fmt.Errorf("Failed to get operation status '%s': %v", id, err)
	}

	switch op.Status {
	case OperationStatusSucceeded:
		return true, nil
	case OperationStatusFailed:
		if op.Error != nil {
			return true, op.Error
		}
		return true, fmt.Errorf("Azure Operation (x-ms-request-id=%s) has failed", id)
	case OperationStatusInProgress:
		return false, nil
	default:
		return false, fmt.Errorf("Unknown operation status returned from API: %s (x-ms-request-id=%s)", op.Status, id)
	}
}
