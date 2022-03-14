// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

var stringEntity = map[string]string{
	"PartitionKey":        "",
	"RowKey":              "",
	"StringTypeProperty1": "StringTypeProperty",
	"StringTypeProperty2": "1970-10-04T00:00:00+00:00",
	"StringTypeProperty3": "c9da6455-213d-42c9-9a79-3e9149a57833",
	"StringTypeProperty4": "BinaryTypeProperty",
	"StringTypeProperty5": fmt.Sprint(2 ^ 32 + 1),
	"StringTypeProperty6": "200.23",
	"StringTypeProperty7": "5",
}

type downloadTestOptions struct{}

//nolint
var downloadTestOpts downloadTestOptions = downloadTestOptions{}

// downloadTestRegister is called once per process
func downloadTestRegister() {

}

type downloadTestGlobal struct {
	perf.PerfTestOptions
	tableName string
}

// NewInsertEntityTest is called once per process
func NewInsertEntityTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	guid, err := uuid.New()
	if err != nil {
		return nil, err
	}
	tableName := fmt.Sprintf("table%s",strings.ReplaceAll(guid.String(), "-", ""))
	d := &downloadTestGlobal{
		PerfTestOptions: options,
		tableName:       tableName,
	}

	connStr, ok := os.LookupEnv("AZURE_TABLES_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_TABLES_CONNECTION_STRING' could not be found")
	}

	svcClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	_, err = svcClient.CreateTable(context.Background(), d.tableName, nil)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *downloadTestGlobal) GlobalCleanup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_TABLES_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_TABLES_CONNECTION_STRING' could not be found")
	}

	svcClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		return err
	}

	_, err = svcClient.DeleteTable(context.Background(), d.tableName, nil)
	return err
}

type downloadPerfTest struct {
	*downloadTestGlobal
	perf.PerfTestOptions
	entity      []byte
	tableClient *aztables.Client
}

// NewPerfTest is called once per goroutine
func (g *downloadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	d := &downloadPerfTest{
		downloadTestGlobal: g,
		PerfTestOptions:    *options,
	}

	connStr, ok := os.LookupEnv("AZURE_TABLES_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_TABLES_CONNECTION_STRING' could not be found")
	}

	svcClient, err := aztables.NewServiceClientFromConnectionString(connStr, &aztables.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: d.PerfTestOptions.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}

	d.tableClient = svcClient.NewClient(g.tableName)

	rk, err := uuid.New()
	if err != nil {
		return nil, err
	}
	pk, err := uuid.New()
	if err != nil {
		return nil, err
	}

	stringEntity["PartitionKey"] = pk.String()
	stringEntity["RowKey"] = rk.String()

	bytes, err := json.Marshal(stringEntity)
	if err != nil {
		return nil, err
	}

	d.entity = bytes

	return d, nil
}

func (d *downloadPerfTest) Run(ctx context.Context) error {
	_, err := d.tableClient.InsertEntity(ctx, d.entity, &aztables.InsertEntityOptions{
		UpdateMode: aztables.EntityUpdateModeMerge,
	})
	return err
}

func (*downloadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
