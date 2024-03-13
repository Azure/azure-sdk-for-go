// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	config = struct {
		id                      string
		storageName             string
		storageNameUserAssigned string
	}{
		id:                      os.Getenv("AZIDENTITY_USER_ASSIGNED_IDENTITY"),
		storageName:             os.Getenv("AZIDENTITY_STORAGE_NAME"),
		storageNameUserAssigned: os.Getenv("AZIDENTITY_STORAGE_NAME_USER_ASSIGNED"),
	}

	missingConfig string
)

func listContainers(account string, cred azcore.TokenCredential) error {
	url := fmt.Sprintf("https://%s.blob.core.windows.net", account)
	log.Printf("listing containers in %s", url)
	client, err := azblob.NewClient(url, cred, nil)
	if err == nil {
		_, err = client.NewListContainersPager(nil).NextPage(context.Background())
	}
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("received a request")
	if missingConfig != "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "need a value for "+missingConfig)
		return
	}

	cred, err := azidentity.NewManagedIdentityCredential(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}
	err = listContainers(config.storageName, cred)
	if err == nil {
		cred, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ResourceID(config.id),
		})
		if err == nil {
			err = listContainers(config.storageNameUserAssigned, cred)
		}
	}

	if err == nil {
		fmt.Fprint(w, "test passed")
		log.Print("test passed")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Print(err)
	}
}

func main() {
	v := []string{}
	if config.id == "" {
		v = append(v, "AZIDENTITY_USER_ASSIGNED_IDENTITY")
	}
	if config.storageName == "" {
		v = append(v, "AZIDENTITY_STORAGE_NAME")
	}
	if config.storageNameUserAssigned == "" {
		v = append(v, "AZIDENTITY_STORAGE_NAME_USER_ASSIGNED")
	}
	if len(v) > 0 {
		missingConfig = strings.Join(v, ", ")
		log.Print("missing values for " + missingConfig)
	}

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler)
	log.Printf("listening on http://127.0.0.1:%s", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
