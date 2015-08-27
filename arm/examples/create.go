package examples

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func withWatcher() autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
			fmt.Printf("Sending %s %s\n", r.Method, r.URL)
			resp, err := s.Do(r)
			fmt.Printf("...received status %s\n", resp.Status)
			if autorest.ResponseRequiresPolling(resp) {
				fmt.Printf("...will poll after %d seconds\n",
					int(autorest.GetPollingDelay(resp, time.Duration(0))/time.Second))
				fmt.Printf("...will poll at %s\n", autorest.GetPollingLocation(resp))
			}
			fmt.Println("")
			return resp, err
		})
	}
}

func createAccount(resourceGroup, name string) {
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

	cna, err := sac.CheckNameAvailability(
		storage.StorageAccountCheckNameAvailabilityParameters{
			Name: name,
			Type: "Microsoft.Storage/storageAccounts"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if !cna.NameAvailable {
		fmt.Printf("%s is unavailable -- try again\n", name)
		return
	}
	fmt.Printf("%s is available\n\n", name)

	sac.Sender = autorest.CreateSender(withWatcher())
	sac.PollingMode = autorest.PollUntilAttempts
	sac.PollingAttempts = 5

	cp := storage.StorageAccountCreateParameters{}
	cp.Location = "westus"
	cp.Properties.AccountType = storage.StandardLRS

	sa, err := sac.Create(resourceGroup, name, cp)
	if err != nil {
		if sa.Response.StatusCode != http.StatusAccepted {
			fmt.Printf("Creation of %s.%s failed with err -- %v\n", resourceGroup, name, err)
			return
		} else {
			fmt.Printf("Create initiated for %s.%s -- poll %s to check status\n",
				resourceGroup,
				name,
				sa.GetPollingLocation())
			return
		}
	}

	fmt.Printf("Successfully created %s.%s\n\n", resourceGroup, name)

	sac.Sender = nil
	r, err := sac.Delete(resourceGroup, name)
	if err != nil {
		fmt.Printf("Delete of %s.%s failed with status %s\n...%v\n", resourceGroup, name, r.Status, err)
		return
	}
	fmt.Printf("Deletion of %s.%s succeeded -- %s\n", resourceGroup, name, r.Status)
}
