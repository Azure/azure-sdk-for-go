package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func handler(w http.ResponseWriter, r *http.Request) {
	v := os.Getenv("AZURE_IDENTITY_TEST_VAULT_URL")
	if v == "" {
		fmt.Fprint(w, "test failed: no value for AZURE_IDENTITY_TEST_VAULT_URL")
		return
	}
	o := azidentity.ManagedIdentityCredentialOptions{}
	if id, ok := os.LookupEnv("AZURE_IDENTITY_TEST_MANAGED_IDENTITY_CLIENT_ID"); ok {
		o.ID = azidentity.ClientID(id)
	}
    var err error
    var client *azsecrets.Client
	cred, err := azidentity.NewManagedIdentityCredential(&o)
	if err == nil {
		client, err = azsecrets.NewClient(v, cred, nil)
		if err == nil {
			pager := client.ListSecrets(nil)
			for pager.NextPage(context.Background()) {
				pager.PageResponse()
			}
			err = pager.Err()
		}
	}
    if err == nil {
        fmt.Fprint(w, "test passed")
    } else {
		fmt.Fprintf(w, "test failed: %s", err.Error())
	}
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/HttpTrigger1", handler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
