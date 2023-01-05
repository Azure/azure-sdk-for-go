package azqueue

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/sas"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// URLParts object represents the components that make up an Azure Storage Queue URL.
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts = sas.URLParts

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the URLParts object.
func ParseURL(u string) (URLParts, error) {
	return sas.ParseURL(u)
}
