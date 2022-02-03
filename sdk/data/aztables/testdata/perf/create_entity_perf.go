// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type aztablesPerfTest struct {
	perf.PerfTestOptions
	client    aztables.Client
	entity    []byte
	tableName string
}

func (a *aztablesPerfTest) createClient() (*aztables.Client, error) {
	options := &aztables.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: a.ProxyInstance,
		},
	}

	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic(errors.New("could not find an environment variable for 'TABLES_STORAGE_ACCOUNT_NAME'"))
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic(errors.New("could not find an environment variable for 'TABLES_PRIMARY_STORAGE_ACCOUNT_KEY'"))
	}
	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	client, err := aztables.NewClientWithSharedKey(fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, a.tableName), cred, options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (a *aztablesPerfTest) GlobalSetup(ctx context.Context) error {
	client, err := a.createClient()
	if err != nil {
		return err
	}

	_, err = client.Create(context.Background(), nil)
	return err
}

func (a *aztablesPerfTest) Setup(ctx context.Context) error {
	client, err := a.createClient()
	if err != nil {
		return nil
	}
	a.client = *client
	return nil
}

func (a *aztablesPerfTest) Run(ctx context.Context) error {
	_, err := a.client.InsertEntity(ctx, a.entity, nil)
	return err
}

func (a *aztablesPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (a *aztablesPerfTest) GlobalCleanup(ctx context.Context) error {
	_, err := a.client.Delete(context.Background(), nil)
	return err
}

func (a *aztablesPerfTest) GetMetadata() perf.PerfTestOptions {
	return a.PerfTestOptions
}

func NewCreateEntityTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "CreateEntityTest"

	e := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Properties: map[string]interface{}{
			"StringTypeProperty": "StringTypeProperty",
			"IntTypeProperty":    5,
			"BoolTypeProperty":   true,
			"DateTimeProperty":   aztables.EDMDateTime(time.Date(1970, time.October, 4, 0, 0, 0, 0, time.UTC)),
			"GuidProperty":       aztables.EDMGUID("673f2ae2-de5d-4e18-8946-37699c36eb73"),
			"BinaryTypeProperty": aztables.EDMBinary([]byte("BinaryTypeProperty")),
			"DoubleTypeProperty": 200.23,
			"Int64TypeProperty":  aztables.EDMInt64(int(math.Pow(2.0, 33))),
		},
	}
	marshalled, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return &aztablesPerfTest{
		PerfTestOptions: *options,
		tableName:       "createEntityTable",
		entity:          marshalled,
	}
}
