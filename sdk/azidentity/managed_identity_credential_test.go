package azidentity

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestManagedIdentityCredential_GetTokenInCloudShell(t *testing.T) {
	managedClient, err := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
	if err != nil {
		fmt.Println("Managed ID error: ", err)
	} else {
		managedAT, err := managedClient.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println(managedAT)
		}
	}
}
