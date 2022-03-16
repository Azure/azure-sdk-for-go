// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

type batchTestOptions struct {
	fullEDM       bool
	clientSharing bool
	count         int
}

var batchTestOpts batchTestOptions = batchTestOptions{
	fullEDM:       false,
	clientSharing: false,
	count:         100,
}

// batchTestRegister is called once per process
func batchTestRegister() {
	flag.IntVar(&listTestOpts.count, "count", 100, "Number of entities to batch create")
	flag.IntVar(&listTestOpts.count, "c", 100, "Number of entities to batch create")
	flag.BoolVar(&batchTestOpts.fullEDM, "full-edm", false, "whether to use entities that utiliza all EDM types for serialization/deserialization, or only strings. Default is only strings")
	flag.BoolVar(&batchTestOpts.clientSharing, "no-client-share", false, "create one ServiceClient per test instance. Default is to share a single ServiceClient")
}

type batchTestGlobal struct {
	perf.PerfTestOptions
	tableName string
}

// NewBatchTest is called once per process
func NewBatchTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	guid, err := uuid.New()
	if err != nil {
		return nil, err
	}
	tableName := fmt.Sprintf("table%s", strings.ReplaceAll(guid.String(), "-", ""))
	d := &batchTestGlobal{
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

func (d *batchTestGlobal) GlobalCleanup(ctx context.Context) error {
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

type batchEntityPerfTest struct {
	*batchTestGlobal
	perf.PerfTestOptions
	baseEDMEntity    aztables.EDMEntity
	baseStringEntity map[string]string
	tableClient      *aztables.Client
}

// NewPerfTest is called once per goroutine
func (g *batchTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	d := &batchEntityPerfTest{
		batchTestGlobal: g,
		PerfTestOptions: *options,
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

	pk, err := uuid.New()
	if err != nil {
		return nil, err
	}

	stringEntity["PartitionKey"] = pk.String()

	d.baseStringEntity = stringEntity

	edmEntity := fullEdm
	edmEntity.PartitionKey = pk.String()
	d.baseEDMEntity = edmEntity

	return d, nil
}

func (d *batchEntityPerfTest) Run(ctx context.Context) error {
	batch := make([]aztables.TransactionAction, batchTestOpts.count)

	for i := 0; i < batchTestOpts.count; i++ {

		if batchTestOpts.fullEDM {
			d.baseEDMEntity.RowKey = fmt.Sprint(i)
			marshalled, err := json.Marshal(d.baseEDMEntity)
			if err != nil {
				return err
			}

			batch[i] = aztables.TransactionAction{
				Entity:     marshalled,
				ActionType: aztables.TransactionTypeUpdateMerge,
			}
		} else {
			d.baseStringEntity["RowKey"] = fmt.Sprint(i)
			marshalled, err := json.Marshal(d.baseStringEntity)
			if err != nil {
				return err
			}

			batch[i] = aztables.TransactionAction{
				Entity:     marshalled,
				ActionType: aztables.TransactionTypeUpdateMerge,
			}
		}

	}

	_, err := d.tableClient.SubmitTransaction(ctx, batch, nil)
	return err
}

func (*batchEntityPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
