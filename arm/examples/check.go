package examples

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/azure/azure-sdk-for-go/arm/storage"
	"github.com/azure/go-autorest/autorest"
	"github.com/azure/go-autorest/autorest/azure"
)

type inspectors struct{}

func (i inspectors) WithInspection() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			fmt.Printf("Inspecting Request: %s %s\n", r.Method, r.URL)
			return p.Prepare(r)
		})
	}
}

func (i inspectors) ByInspecting() autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			fmt.Printf("Inspecting Response: %s for %s %s\n", resp.Status, resp.Request.Method, resp.Request.URL)
			return r.Respond(resp)
		})
	}
}

func checkName(name string) {
	c, err := helpers.LoadCredentials()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	sac := storage.NewStorageAccountsClient(c["subscriptionID"])

	spt, err := helpers.NewServicePrincipalTokenFromCredentials(c, azure.AzureResourceManagerScope)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	sac.Authorizer = spt

	sac.Sender = autorest.CreateSender(
		autorest.WithLogging(log.New(os.Stdout, "sdk-example: ", log.LstdFlags)))

	i := inspectors{}
	sac.RequestInspector = i
	sac.ResponseInspector = i

	cna, err := sac.CheckNameAvailability(
		storage.StorageAccountCheckNameAvailabilityParameters{
			Name: name,
			Type: "Microsoft.Storage/storageAccounts"})

	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		if cna.NameAvailable {
			fmt.Printf("The name '%s' is available\n", name)
		} else {
			fmt.Printf("The name '%s' is unavailable because %s\n", name, cna.Message)
		}
	}
}
