// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/cmd/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/spf13/cobra"
)

var aztablesCmd = &cobra.Command{
	Use:   "CreateEntityTest",
	Short: "aztables perf test for creating an entity",
	Long:  "aztables perf test for creating an entity. This test uses the `Client.InsertEntity` method to merge an entity on insertion.",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return RunPerfTest(&aztablesPerfTest{})
	},
}

func init() {
	rootCmd.AddCommand(aztablesCmd)
}

type aztablesPerfTest struct {
	client       aztables.Client
	entity       []byte
	partitionKey string
	rowKey       string
	tableName    string
}

func (a *aztablesPerfTest) createClient() error {

	options := &aztables.ClientOptions{}
	if TestProxy == "http" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: true, TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &aztables.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
	} else if TestProxy == "https" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: true, TestName: a.GetMetadata()})
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
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	client, err := aztables.NewClient(fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, a.tableName), cred, options)
	if err != nil {
		return err
	}
	a.client = *client

	return nil
}

func (a *aztablesPerfTest) GlobalSetup(ctx context.Context) error {
	a.tableName = "randomTableName"
	if TestProxy != "" {
		err := recording.Start(a.GetMetadata(), nil)
		if err != nil {
			return err
		}
	}

	err := a.createClient()
	if err != nil {
		if TestProxy != "" {
			recording.Stop(a.GetMetadata(), nil)
		}
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
		if TestProxy != "" {
			recording.Stop(a.GetMetadata(), nil)
		}
		return err
	}
	a.entity = marshalled

	_, err = a.client.Create(context.Background(), nil)
	return err
}

func (a *aztablesPerfTest) GlobalTearDown(ctx context.Context) error {
	_, err := a.client.Delete(context.Background(), nil)
	if err != nil {
		return err
	}
	return recording.Stop(a.GetMetadata(), nil)
}

func (a *aztablesPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (a *aztablesPerfTest) Run(ctx context.Context) error {
	_, err := a.client.InsertEntity(ctx, a.entity, nil)
	return err
}

func (a *aztablesPerfTest) TearDown(ctx context.Context) error {
	_, err := a.client.DeleteEntity(ctx, a.partitionKey, a.rowKey, nil)
	var azErr azcore.HTTPResponse
	if errors.As(err, &azErr) {
		if azErr.RawResponse().StatusCode == http.StatusNotFound {
			return nil
		}
	}
	return nil
}

func (a *aztablesPerfTest) GetMetadata() string {
	return "aztables-update"
}
