package azidentity

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func Test_GetToken_CloudShell(t *testing.T) {
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

// func Test_GetToken_NilScopes(t *testing.T) {
// 	managedClient, err := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
// 	if err != nil {
// 		t.Fatalf("Received an unexpected error when creating a New ManagedIdentityCredential")
// 	}

// 	_, err = managedClient.GetToken(context.Background(), azcore.TokenRequestOptions{})
// 	if err != nil {
// 		var authFailed *AuthenticationFailedError
// 		if !errors.As(err, &authFailed) {
// 			t.Fatalf("Expected an AuthenticationFailedError, but instead got: %T", err)
// 		}
// 	} else {
// 		t.Fatalf("Expected an error but did not receive one.")
// 	}
// }
