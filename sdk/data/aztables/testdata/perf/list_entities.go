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

type listEntitiesTestOptions struct {
	fullEDM       bool
	clientSharing bool
	count         int
}

var listTestOpts listEntitiesTestOptions = listEntitiesTestOptions{
	fullEDM:       false,
	clientSharing: false,
	count:         100,
}

// listTestRegister is called once per process
func listTestRegister() {
	flag.IntVar(&listTestOpts.count, "count", 100, "Number of entities to list")
	flag.IntVar(&listTestOpts.count, "c", 100, "Number of entities to list")
	flag.BoolVar(&listTestOpts.fullEDM, "full-edm", false, "whether to use entities that utiliza all EDM types for serialization/deserialization, or only strings. Default is only strings")
	flag.BoolVar(&listTestOpts.clientSharing, "no-client-share", false, "create one ServiceClient per test instance. Default is to share a single ServiceClient")
}

type listEntityTestGlobal struct {
	perf.PerfTestOptions
	tableName string
	svcClient *aztables.ServiceClient
}

// NewListEntitiesTest is called once per process
func NewListEntitiesTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	guid, err := uuid.New()
	if err != nil {
		return nil, err
	}
	tableName := fmt.Sprintf("table%s", strings.ReplaceAll(guid.String(), "-", ""))
	d := &listEntityTestGlobal{
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
	d.svcClient = svcClient

	client := d.svcClient.NewClient(d.tableName)

	baseEntityEDM := fullEdm
	baseEntityString := stringEntity

	u, err := uuid.New()
	if err != nil {
		return nil, err
	}

	baseEntityEDM.PartitionKey = u.String()
	baseEntityString["PartitionKey"] = u.String()

	for i := 0; i < listTestOpts.count; i++ {
		if listTestOpts.fullEDM {
			u, err := uuid.New()
			if err != nil {
				return nil, err
			}
			baseEntityEDM.RowKey = u.String()

			marshalled, err := json.Marshal(baseEntityEDM)
			if err != nil {
				return nil, err
			}

			_, err = client.UpsertEntity(ctx, marshalled, nil)
			if err != nil {
				return nil, err
			}
		} else {
			u, err := uuid.New()
			if err != nil {
				return nil, err
			}
			baseEntityString["RowKey"] = u.String()

			marshalled, err := json.Marshal(baseEntityString)
			if err != nil {
				return nil, err
			}

			_, err = client.UpsertEntity(ctx, marshalled, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return d, nil
}

func (d *listEntityTestGlobal) GlobalCleanup(ctx context.Context) error {
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

type listEntitiesPerfTest struct {
	*listEntityTestGlobal
	perf.PerfTestOptions
	client *aztables.Client
}

// NewPerfTest is called once per goroutine
func (g *listEntityTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	d := &listEntitiesPerfTest{
		listEntityTestGlobal: g,
		PerfTestOptions:      *options,
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

	d.client = svcClient.NewClient(g.tableName)

	return d, nil
}

func (d *listEntitiesPerfTest) Run(ctx context.Context) error {
	pager := d.client.NewListEntitiesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		_ = resp
	}
	return nil
}

func (*listEntitiesPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
