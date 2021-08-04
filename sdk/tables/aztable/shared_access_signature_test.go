// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSasCreateTable(t *testing.T) {

	tempSAS := "blahblahblah"
	sasCredential, err := NewAzureSasCredential(tempSAS)
	require.NoError(t, err)

	client, err := NewTableClient("sastable", "https://seankaneprim.table.core.windows.net", sasCredential, nil)
	require.NoError(t, err)

	_, err = client.Create(context.Background())
	require.NoError(t, err)
}

func TestCreateSAS(t *testing.T) {
	accountName := os.Getenv("TABLES_PRIMARY_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	start := time.Date(2021, time.August, 3, 18, 0, 6, 0, time.UTC)
	end := time.Date(2021, time.August, 4, 2, 0, 6, 0, time.UTC)

	properties := AccountSignatureProperties{
		ResourceTypes: ResourceType{
			Service:   true,
			Container: true,
			Object:    true,
		},
		Permissions: AccountSasPermissions{
			Read:   true,
			Write:  true,
			Delete: true,
			List:   true,
			Add:    true,
			Create: true,
			Update: true,
		},
		Start:    &start,
		Expiry:   &end,
		Protocol: SasProtocolHttps,
	}

	key, err := GenerateAccountSignature(*cred, properties)
	require.NoError(t, err)
	fmt.Println(key)

	sascred, err := NewAzureSasCredential(key)
	require.NoError(t, err)

	client, err := NewTableClient("sastable", fmt.Sprintf("https://%v.table.core.windows.net", accountName), sascred, nil)
	require.NoError(t, err)
	_, err = client.Create(context.Background())
	require.NoError(t, err)

	_, err = client.Delete(context.Background())
	require.NoError(t, err)
}

func TestCreateTableSAS(t *testing.T) {
	accountName := os.Getenv("TABLES_PRIMARY_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	start := time.Now().UTC()       //time.Date(2021, time.August, 3, 18, 0, 6, 0, time.UTC)
	end := start.Add(time.Hour * 4) //time.Date(2021, time.August, 4, 2, 0, 6, 0, time.UTC)

	properties := TableSignatureProperties{
		TableName: "tablesastable",
		Permissions: AccountSasPermissions{
			Read:   true,
			Write:  true,
			Delete: true,
			List:   true,
			Add:    true,
			Create: true,
			Update: true,
		},
		Start:    &start,
		Expiry:   &end,
	}

	key, err := GenerateTableSignature(*cred, properties)
	require.NoError(t, err)
	fmt.Println(key)

	sascred, err := NewAzureSasCredential(key)
	require.NoError(t, err)

	client, err := NewTableClient(properties.TableName, fmt.Sprintf("https://%v.table.core.windows.net", accountName), sascred, nil)
	require.NoError(t, err)
	_, err = client.Create(context.Background())
	require.NoError(t, err)

	simpleEntity := map[string]string{
		"PartitionKey": "pk001",
		"RowKey":       "rk001",
		"Value":        "4",
	}
	marshalled, err := json.Marshal(simpleEntity)
	require.NoError(t, err)

	_, err = client.AddEntity(context.Background(), marshalled)
	require.NoError(t, err)

	_, err = client.Delete(context.Background())
	require.NoError(t, err)
}
