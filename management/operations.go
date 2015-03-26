package management

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

//operation represents an in-flight operation. Use client.getOperationStatus()
//to get the operation given the operation ID, or use
//client.waitAsyncOperation() to block until the operation has completed.
type operation struct {
	XMLName        xml.Name `xml:"Operation"`
	ID             string
	Status         string
	HttpStatusCode string
	Error          AzureError
}

// OperationId is assigned by Azure API and can be used to look up the status of an operation
type OperationId string

//getOperationStatus gets an operation given the operation ID.
func (client *Client) getOperationStatus(operationId OperationId) (operation, error) {
	operation := operation{}
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

//waitAsyncOperation blocks until the operation with the given operationId is
//no longer in the InProgress state. If the operation was successful, nothing is
//returned, otherwise an error is returned.
func (client *Client) WaitAsyncOperation(operationId OperationId) error {
	if operationId == "" {
		return fmt.Errorf(errParamNotSpecified, "operationId")
	}

	status := "InProgress"
	operation := operation{}
	err := errors.New("")
	for status == "InProgress" {
		time.Sleep(2000 * time.Millisecond)
		operation, err = client.getOperationStatus(operationId)
		if err != nil {
			return err
		}

		status = operation.Status
	}

	if status == "Failed" {
		return fmt.Errorf("Azure operation %s failed. Code: %s, Message: %s", operationId, operation.Error.Code, operation.Error.Message)
	}

	return nil
}
