// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

	// jwtRegex is used to redact JWTs (e.g. access tokens) in log output sent to a test client, although
	// that output should never contain tokens because it's sent only when a test fails i.e., the request
	// handler couldn't obtain an access token
	jwtRegex   = regexp.MustCompile(`ey\S+\.\S+\.\S+`)
	logOptions = policy.LogOptions{
		AllowedQueryParams: []string{"client_id", "msi_res_id", "object_id", "resource"},
		IncludeBody:        true,
	}
	// logs collects log output from a test run to help debug failures. Note that its usage isn't
	// concurrency-safe and that's okay because live managed identity tests targeting this server
	// don't send concurrent requests.
	logs          strings.Builder
	missingConfig string
)

func credential(id azidentity.ManagedIDKind) (azcore.TokenCredential, error) {
	co := azcore.ClientOptions{Logging: logOptions}
	if config.workloadID {
		// the identity is determined by service account configuration
		return azidentity.NewWorkloadIdentityCredential(&azidentity.WorkloadIdentityCredentialOptions{ClientOptions: co})
	}
	return azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: co,
		ID:            id,
	})
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
	logs.Reset()
	log.Print("received a request")
	if missingConfig != "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "need a value for "+missingConfig)
		return
	}

	cred, err := credential(nil)
	if err == nil {
		name := "ManagedIdentityCredential"
		if config.workloadID {
			name = "WorkloadIdentityCredential"
		}
		logs.WriteString("\n*** testing " + name + "\n\n")
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
		// discard logs from the successful tests above
		logs.Reset()
		logs.WriteString("*** testing DefaultAzureCredential\n\n")
		cred, err = azidentity.NewDefaultAzureCredential(
			&azidentity.DefaultAzureCredentialOptions{
				ClientOptions: azcore.ClientOptions{Logging: logOptions},
			},
		)
		if err == nil {
			err = listContainers(config.storageName, cred)
		}
	}

	msg := "test passed"
	if err != nil {
		logs.WriteString("\n*** test failed with error: " + err.Error() + "\n")
		msg = logs.String()
	}
	fmt.Fprint(w, msg)
	log.Print(msg)
}

func main() {
	azlog.SetListener(func(_ azlog.Event, msg string) {
		msg = jwtRegex.ReplaceAllString(msg, "***")
		logs.WriteString(msg + "\n\n")
	})
	azlog.SetEvents(azidentity.EventAuthentication, azlog.EventRequest, azlog.EventResponse)

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
