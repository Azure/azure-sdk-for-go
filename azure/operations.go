package azure

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
	XMLName        xml.Name `xml:"operation"`
	ID             string
	Status         string
	HttpStatusCode string
	Error          AzureError
}

//getOperationStatus gets an operation given the operation ID.
func (client *Client) getOperationStatus(operationId string) (*operation, error) {
	if len(operationId) == 0 {
		return nil, fmt.Errorf(paramNotSpecifiedError, "operationId")
	}

	operation := new(operation)
	url := "operations/" + operationId
	response, azureErr := client.SendAzureGetRequest(url)
	if azureErr != nil {
		return nil, azureErr
	}

	err := xml.Unmarshal(response, operation)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

//waitAsyncOperation blocks until the operation with the given operationId is
//no longer in the InProgress state. If the operation was successful, nothing is
//returned, otherwise an error is returned.
func (client *Client) WaitAsyncOperation(operationId string) error {
	if len(operationId) == 0 {
		return fmt.Errorf(paramNotSpecifiedError, "operationId")
	}

	status := "InProgress"
	operation := new(operation)
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
		return errors.New(operation.Error.Message)
	}

	return nil
}
