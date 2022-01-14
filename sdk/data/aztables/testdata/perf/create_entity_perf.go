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
	client    aztables.Client
	entity    []byte
	tableName string
}

func (a *aztablesPerfTest) createClient() error {
	options := &aztables.ClientOptions{}
	if perf.TestProxy == "http" {
		t, err := perf.NewProxyTransport(&perf.TransportOptions{TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &aztables.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
	} else if perf.TestProxy == "https" {
		t, err := perf.NewProxyTransport(&perf.TransportOptions{TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &aztables.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
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
		return err
	}
	client, err := aztables.NewClientWithSharedKey(fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, a.tableName), cred, options)
	if err != nil {
		return err
	}
	a.client = *client

	return nil
}

func (a *aztablesPerfTest) GlobalSetup(ctx context.Context) error {
	a.tableName = "randomTableName"

	err := a.createClient()
	if err != nil {
		return err
	}
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
		return err
	}
	a.entity = marshalled

	_, err = a.client.Create(context.Background(), nil)
	return err
}

func (a *aztablesPerfTest) GlobalTearDown(ctx context.Context) error {
	_, err := a.client.Delete(context.Background(), nil)
	return err
}

func (a *aztablesPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (a *aztablesPerfTest) Run(ctx context.Context) error {
	_, err := a.client.InsertEntity(ctx, a.entity, nil)
	return err
}

func (a *aztablesPerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (a *aztablesPerfTest) GetMetadata() string {
	return "CreateEntity"
}
