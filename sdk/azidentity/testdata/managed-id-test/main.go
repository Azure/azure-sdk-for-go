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
		// resourceID of a managed identity permitted to list blobs in the account specified by storageNameUserAssigned.
		resourceID string
		// storageName is the name of a storage account accessible by the default or system-assigned identity
		storageName string
		// storageNameUserAssigned is the name of a storage account accessible by the identity specified by
		// resourceID. The default or system-assigned identity shouldn't have any permission for this account.
		storageNameUserAssigned string
		// workloadID determines whether this app tests ManagedIdentityCredential or WorkloadIdentityCredential.
		// When true, the app ignores resourceID and storageNameUserAssigned.
		workloadID bool
	}{
		resourceID:              os.Getenv("AZIDENTITY_USER_ASSIGNED_IDENTITY"),
		storageName:             os.Getenv("AZIDENTITY_STORAGE_NAME"),
		storageNameUserAssigned: os.Getenv("AZIDENTITY_STORAGE_NAME_USER_ASSIGNED"),
		workloadID:              os.Getenv("AZIDENTITY_USE_WORKLOAD_IDENTITY") != "",
	}

	missingConfig string
)

func credential(resourceID string) (azcore.TokenCredential, error) {
	if config.workloadID {
		// the identity is determined by service account configuration
		return azidentity.NewWorkloadIdentityCredential(nil)
	}
	opts := azidentity.ManagedIdentityCredentialOptions{}
	if resourceID != "" {
		opts.ID = azidentity.ResourceID(resourceID)
	}
	return azidentity.NewManagedIdentityCredential(&opts)
}

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

	cred, err := credential("")
	if err == nil {
		err = listContainers(config.storageName, cred)
		if !config.workloadID && err == nil {
			cred, err = credential(config.resourceID)
			if err == nil {
				err = listContainers(config.storageNameUserAssigned, cred)
			}
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
	if config.storageName == "" {
		v = append(v, "AZIDENTITY_STORAGE_NAME")
	}
	if config.workloadID {
		log.Print("Testing WorkloadIdentityCredential")
	} else {
		log.Print("Testing ManagedIdentityCredential")
		if config.resourceID == "" {
			v = append(v, "AZIDENTITY_USER_ASSIGNED_IDENTITY")
		}
		if config.storageNameUserAssigned == "" {
			v = append(v, "AZIDENTITY_STORAGE_NAME_USER_ASSIGNED")
		}
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
