package management

import (
	"encoding/xml"
	"fmt"
)

// AzureError represents an error returned by the management API. It has an error
// code (for example, ResourceNotFound) and a descriptive message.
type AzureError struct {
	XMLName xml.Name `xml:"Error"`
	Code    string
	Message string
}

//Error implements the error interface for the AzureError type.
func (e AzureError) Error() string {
	return fmt.Sprintf("Error response from Azure. Code: %s, Message: %s", e.Code, e.Message)
}

// IsResourceNotFoundError returns true if the provided error is an AzureError
// reporting that a given resource has not been found.
func IsResourceNotFoundError(err error) bool {
	azureErr, ok := err.(AzureError)
	return ok && azureErr.Code == "ResourceNotFound"
}
