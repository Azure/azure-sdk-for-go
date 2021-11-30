package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

type perfTest struct {
	client  *azkeys.Client
	keyName string
}

func (p *perfTest) Setup() {
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZKEYS_TENANT_ID"),
		os.Getenv("AZKEYS_CLIENT_ID"),
		os.Getenv("AZKEYS_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(os.Getenv("AZURE_KEYVAULT_URL"), cred, nil)
	if err != nil {
		panic(err)
	}

	p.client = client
	s := rand.NewSource(time.Now().Unix())
	p.keyName = fmt.Sprintf("%d", s.Int63())

	_, err = p.client.CreateRSAKey(context.Background(), p.keyName, nil)
	if err != nil {
		panic(err)
	}
}

func (p *perfTest) Run() {
	_, err := p.client.GetKey(context.Background(), p.keyName, nil)
	if err != nil {
		panic(err)
	}
}

func (p *perfTest) TearDown() {
	poller, err := p.client.BeginDeleteKey(context.Background(), p.keyName, nil)
	if err != nil {
		panic(err)
	}
	_, err = poller.PollUntilDone(context.Background(), time.Second)
	if err != nil {
		panic(err)
	}
}

func main() {
	perf.RunPerfTest(&perfTest{})
}
