package azidentity

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestManagedIdentityCredential_GetTokenInCloudShell(t *testing.T) {
	msiEndpoint := os.Getenv("MSI_ENDPOINT")
	if len(msiEndpoint) == 0 {
		t.Skip()
	}
	msiCred := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
	_, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
		if err != nil {
			t.Fatalf("Received an error when attempting to retrieve a token")
		}
	}
