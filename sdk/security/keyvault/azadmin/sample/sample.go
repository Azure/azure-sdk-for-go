package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
)

func main() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	client, err := ekm.NewClient("https://contosohsm.managedhsm.azure.net/", credential, nil)
	if err != nil {
		panic(err)
	}

	ekmConnection, err := client.GetEkmConnection(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("EKM Connection Host: %s\n", *ekmConnection.Host)

	ekmCertificate, err := client.GetEkmCertificate(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("EKM Certificate Subject: %s\n", *ekmCertificate.SubjectCommonName)
}
