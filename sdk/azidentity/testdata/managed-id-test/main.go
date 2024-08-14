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
		// clientID, objectID, resourceID of a managed identity permitted to list blobs in the account specified by storageNameUserAssigned
		clientID, objectID, resourceID azidentity.ManagedIDKind
		// storageName is the name of a storage account accessible by the default or system-assigned identity
		storageName string
		// storageNameUserAssigned is the name of a storage account accessible by the identity specified by
		// resourceID. The default or system-assigned identity shouldn't have any permission for this account.
		storageNameUserAssigned string
		// workloadID determines whether this app tests ManagedIdentityCredential or WorkloadIdentityCredential.
		// When true, the app ignores clientID, objectID, resourceID and storageNameUserAssigned.
		workloadID bool
	}{
		clientID:                azidentity.ClientID(os.Getenv("AZIDENTITY_USER_ASSIGNED_IDENTITY_CLIENT_ID")),
		objectID:                azidentity.ObjectID(os.Getenv("AZIDENTITY_USER_ASSIGNED_IDENTITY_OBJECT_ID")),
		resourceID:              azidentity.ResourceID(os.Getenv("AZIDENTITY_USER_ASSIGNED_IDENTITY")),
		storageName:             os.Getenv("AZIDENTITY_STORAGE_NAME"),
		storageNameUserAssigned: os.Getenv("AZIDENTITY_STORAGE_NAME_USER_ASSIGNED"),
		workloadID:              os.Getenv("AZIDENTITY_USE_WORKLOAD_IDENTITY") != "",
	}

	missingConfig string
)

func credential(id azidentity.ManagedIDKind) (azcore.TokenCredential, error) {
	if config.workloadID {
		// the identity is determined by service account configuration
		return azidentity.NewWorkloadIdentityCredential(nil)
	}
	return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{ID: id})
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

	cred, err := credential(nil)
	if err == nil {
		err = listContainers(config.storageName, cred)
	}
	if err == nil && !config.workloadID {
		for _, id := range []azidentity.ManagedIDKind{config.clientID, config.objectID, config.resourceID} {
			cred, err = credential(id)
			if err == nil {
				err = listContainers(config.storageNameUserAssigned, cred)
			}
			if err != nil {
				break
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
		if config.clientID.String() == "" {
			v = append(v, "AZIDENTITY_USER_ASSIGNED_IDENTITY_CLIENT_ID")
		}
		if config.objectID.String() == "" {
			v = append(v, "AZIDENTITY_USER_ASSIGNED_IDENTITY_OBJECT_ID")
		}
		if config.resourceID.String() == "" {
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
