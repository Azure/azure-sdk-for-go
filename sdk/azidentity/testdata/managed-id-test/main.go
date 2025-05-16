// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
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
	clientOptions = azcore.ClientOptions{
		Logging: policy.LogOptions{
			AllowedQueryParams: []string{"client_id", "mi_res_id", "msi_res_id", "object_id", "principal_id", "resource"},
			IncludeBody:        true,
		},
	}

	// jwtRegex is used to redact JWTs (e.g. access tokens) in log output sent to a client, although
	// that output should never contain tokens because it's sent only when authentication fails
	jwtRegex = regexp.MustCompile(`ey\S+\.\S+\.\S+`)

	// logs collects SDK log output from a test run to help debug failures. Note that its usage
	// isn't concurrency-safe and that's okay because live managed identity tests targeting this
	// server don't send concurrent requests.
	logs strings.Builder
)

type testOptions struct {
	id          azidentity.ManagedIDKind
	storageName string
}

func options(r *http.Request) (testOptions, error) {
	c := testOptions{}

	q := r.URL.Query()
	if c.storageName = q.Get("storage-name"); c.storageName == "" {
		return testOptions{}, errors.New("storage-name parameter is required")
	}

	if id := q.Get("client-id"); id != "" {
		c.id = azidentity.ClientID(id)
	}
	if id := q.Get("object-id"); id != "" {
		if c.id != nil {
			return testOptions{}, errors.New("multiple IDs specified")
		}
		c.id = azidentity.ObjectID(id)
	}
	if resourceID := q.Get("resource-id"); resourceID != "" {
		if c.id != nil {
			return testOptions{}, errors.New("multiple IDs specified")
		}
		c.id = azidentity.ResourceID(resourceID)
	}

	return c, nil
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

func handleDefaultAzureCredential(w http.ResponseWriter, r *http.Request) {
	logs.Reset()
	logs.WriteString("testing DefaultAzureCredential\n\n")

	o, err := options(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{
			ClientOptions: clientOptions,
		},
	)
	if err == nil {
		err = listContainers(o.storageName, cred)
	}

	msg := "test passed"
	if err != nil {
		logs.WriteString("\ntest failed with error: " + err.Error() + "\n")
		msg = logs.String()
	}
	fmt.Fprint(w, msg)
	log.Print(msg)
}

func handleManagedIdentityCredential(w http.ResponseWriter, r *http.Request) {
	logs.Reset()
	logs.WriteString("testing ManagedIdentityCredential\n\n")

	config, err := options(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientOptions,
		ID:            config.id,
	})

	if err == nil {
		err = listContainers(config.storageName, cred)
	}

	if err == nil {
		cred, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ClientOptions: clientOptions,
			ID:            config.id,
		})
		if err == nil {
			err = listContainers(config.storageName, cred)
		}
	}

	msg := "test passed"
	if err != nil {
		logs.WriteString("\ntest failed with error: " + err.Error() + "\n")
		msg = logs.String()
	}
	fmt.Fprint(w, msg)
	log.Print(msg)
}

func handleWorkloadIdentityCredential(w http.ResponseWriter, r *http.Request) {
	logs.Reset()
	logs.WriteString("testing WorkloadIdentityCredential\n\n")

	config, err := options(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	cred, err := azidentity.NewWorkloadIdentityCredential(&azidentity.WorkloadIdentityCredentialOptions{
		ClientOptions: clientOptions,
	})
	if err == nil {
		err = listContainers(config.storageName, cred)
	}

	msg := "test passed"
	if err != nil {
		logs.WriteString("\ntest failed with error: " + err.Error() + "\n")
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

	http.HandleFunc("/dac", handleDefaultAzureCredential)
	http.HandleFunc("/mic", handleManagedIdentityCredential)
	http.HandleFunc("/wic", handleWorkloadIdentityCredential)

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on http://127.0.0.1:%s", port)
	log.Print(http.ListenAndServe(":"+port, nil))
}
