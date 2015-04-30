// Package testutils contains some test utilities for the Azure SDK
package testutils

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/management"
)

// Returns a management Client for testing. Expects AZSUBSCRIPTIONID
// and AZCERTDATA to be present in the environment. AZCERTDATA is the
// base64encoded binary representation of the PEM certificate data.
func GetTestClient(t *testing.T) management.Client {
	subid := os.Getenv("AZSUBSCRIPTIONID")
	certdata := os.Getenv("AZCERTDATA")
	if subid == "" || certdata == "" {
		t.Skip("AZSUBSCRIPTIONID or AZCERTDATA not set, skipping test")
	}
	cert, err := base64.StdEncoding.DecodeString(certdata)
	if err != nil {
		t.Fatal(err)
	}

	client, err := management.NewClient(subid, cert)
	if err != nil {
		t.Fatal(err)
	}
	return client
}
