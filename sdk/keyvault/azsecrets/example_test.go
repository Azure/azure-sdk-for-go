//go:build go1.16
// +build go1.16

package azsecrets_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func ExampleNewClient() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}
	_ = client
}

func ExampleClient_SetSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	secret := "mySecret"
	value := "mySecretValue"

	resp, err := client.SetSecret(context.Background(), secret, value, nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		fmt.Println("Operation was not successful")
		fmt.Println(httpErr.RawResponse())
		panic(httpErr)
	}

	fmt.Printf("Set secret %s", *resp.ID)
}

func ExampleClient_GetSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.GetSecret(context.Background(), "mySecretName", nil)
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		fmt.Println("Operation was not successful")
		fmt.Println(httpErr.RawResponse())
		panic(httpErr)
	}

	fmt.Printf("Secret Name: %s\tSecret Value: %s", *resp.ID, *resp.Value)
}

func ExampleClient_BeginDeleteSecret() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.BeginDeleteSecret(context.Background(), "secretToDelete", nil)
	resp.PollUntilDone(context.Background(), 250*time.Millisecond)

	poller := resp.Poller
	finalResp, err := poller.FinalResponse(context.Background())
	var httpErr azcore.HTTPResponse
	if errors.As(err, &httpErr) {
		fmt.Println("Operation was not successful")
		fmt.Println(httpErr.RawResponse())
		panic(httpErr)
	}

	fmt.Printf("Deleted secret with name %s\n", *finalResp.ID)
}

func ExampleClient_ListSecrets() {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		panic(err)
	}

	pager := client.ListSecrets(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		fmt.Printf("Found %d secrets in this page.\n", len(resp.Secrets))
		for _, v := range resp.Secrets {
			fmt.Printf("Secret Name: %s\tSecret Tags: %v", *v.ID, v.Tags)
		}
	}

	if pager.Err() != nil {
		var httpErr azcore.HTTPResponse
		if errors.As(pager.Err(), &httpErr) {
			// handle error
		} else {
			// handle non HTTP Error
		}
	}
}
