// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestServiceErrorsServiceClient(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service := initServiceTest(t, service, NewSpanValidator(t, SpanMatcher{
				Name: "ServiceClient.DeleteTable",
			}))

			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			_, err = service.CreateTable(ctx, tableName, nil)
			require.NoError(t, err)

			// Create a duplicate table to produce an error
			_, err = service.CreateTable(ctx, tableName, nil)
			require.Error(t, err)
			var httpErr *azcore.ResponseError
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, string(TableAlreadyExists), httpErr.ErrorCode)
			require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))

			_, err = service.DeleteTable(ctx, tableName, nil)
			require.NoError(t, err)
		})
	}
}

func TestCreateTableFromService(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service := initServiceTest(t, service, NewSpanValidator(t, SpanMatcher{
				Name: "ServiceClient.CreateTable",
			}))

			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			_, err = service.CreateTable(ctx, tableName, nil)
			deleteTable := func() {
				_, err := service.DeleteTable(ctx, tableName, nil)
				if err != nil {
					fmt.Printf("Error cleaning up test. %v\n", err.Error())
				}
			}
			t.Cleanup(deleteTable)

			require.NoError(t, err)
			// require.Equal(t, *resp.TableResponse.TableName, tableName)
		})
	}
}

func TestQueryTable(t *testing.T) {
	for _, svc := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), svc), func(t *testing.T) {
			service := initServiceTest(t, svc, tracing.Provider{})

			tableCount := 5
			tableNames := make([]string, tableCount)
			prefix1 := "zzza"
			prefix2 := "zzzb"

			// clean up the tables on the last test for that resource (storage, cosmos)
			if svc == cosmosTokenCredentialEndpoint || svc == storageTokenCredentialEndpoint {
				t.Cleanup(func() {
					require.NoError(t, clearAllTables(service))
				})
			}

			// create 10 tables with our exected prefix and 1 with a different prefix
			for i := 0; i < tableCount; i++ {
				if i < (tableCount - 1) {
					name := fmt.Sprintf("%v%v", prefix1, i)
					tableNames[i] = name
				} else {
					name := fmt.Sprintf("%v%v", prefix2, i)
					tableNames[i] = name
				}

				// only create the tables in the first test type for that resource (ie, storage, cosmos)
				if svc == cosmosEndpoint || svc == storageEndpoint {
					_, err := service.CreateTable(ctx, tableNames[i], nil)
					require.NoError(t, err)
				}
			}

			// Query for tables with no pagination. The filter should exclude one table from the results
			filter := fmt.Sprintf("TableName ge '%s' and TableName lt '%s'", prefix1, prefix2)
			pager := service.NewListTablesPager(&ListTablesOptions{Filter: &filter})

			resultCount := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				resultCount += len(resp.Tables)
			}

			require.Equal(t, resultCount, tableCount-1)

			// Query for tables with pagination
			top := int32(2)
			pager = service.NewListTablesPager(&ListTablesOptions{Filter: &filter, Top: &top})

			resultCount = 0
			pageCount := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				require.LessOrEqual(t, len(resp.Tables), 2)
				resultCount += len(resp.Tables)
				pageCount++
				for _, table := range resp.Tables {
					require.Nil(t, table.Value)
				}
			}

			require.Equal(t, resultCount, tableCount-1)
			if svc == "storage" {
				require.Equal(t, pageCount, int(top))
			}
		})
	}
}

type mdForListTables struct {
	EditLink string `json:"odata.editLink"`
	ID       string `json:"odata.id"`
	Type     string `json:"odata.type"`
}

func TestListTables(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			client := initServiceTest(t, service, NewSpanValidator(t, SpanMatcher{
				Name: "Pager[ListTablesResponse].NextPage",
			}))

			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			err = clearAllTables(client)
			require.NoError(t, err)

			for i := 0; i < 5; i++ {
				_, err := client.CreateTable(ctx, fmt.Sprintf("%v%v", tableName, i), nil)
				require.NoError(t, err)
			}

			count := 0
			pager := client.NewListTablesPager(&ListTablesOptions{
				Format: to.Ptr(MetadataFormatFull),
			})
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.Tables)

				for _, table := range resp.Tables {
					if service == storageEndpoint {
						// cosmos doesn't send full metadata
						require.NotEmpty(t, table.Value)
						var md mdForListTables
						require.NoError(t, json.Unmarshal(table.Value, &md))
						require.NotEmpty(t, md.EditLink)
						require.NotEmpty(t, md.ID)
						require.NotEmpty(t, md.Type)
					}
				}
			}

			require.Equal(t, 5, count)

			deleteTable := func() {
				for i := 0; i < 5; i++ {
					_, err := client.DeleteTable(ctx, fmt.Sprintf("%v%v", tableName, i), nil)
					if err != nil {
						fmt.Printf("Error cleaning up test. %v\n", err.Error())
					}
				}
			}
			t.Cleanup(deleteTable)
		})
	}
}

// This functionality is only available on storage accounts
func TestGetStatistics(t *testing.T) {
	var cred *SharedKeyCredential
	var err error

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	stop := func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
	defer stop()

	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeAccountKey")

	if recording.GetRecordMode() == recording.PlaybackMode {
		cred, err = NewSharedKeyCredential("fakeaccount", "fakeAccountKey==")
	} else {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
	}

	serviceURL := storageURI(accountName + "-secondary")
	service, err := createServiceClientForRecordingForSharedKey(t, serviceURL, *cred, NewSpanValidator(t, SpanMatcher{
		Name: "ServiceClient.GetStatistics",
	}))
	require.NoError(t, err)

	resp, err := service.GetStatistics(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GeoReplication.LastSyncTime)
	require.NotNil(t, resp.GeoReplication.Status)
}

// Functionality is only available on storage accounts
func TestGetProperties(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, NewSpanValidator(t, SpanMatcher{
		Name: "ServiceClient.GetProperties",
	}))

	resp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// Logging is only available on storage accounts
func TestSetLogging(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, NewSpanValidator(t, SpanMatcher{
		Name: "ServiceClient.SetProperties",
	}))

	getResp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	getResp.Logging = &Logging{
		Read:    to.Ptr(true),
		Write:   to.Ptr(true),
		Delete:  to.Ptr(true),
		Version: to.Ptr("1.0"),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.Ptr(true),
			Days:    to.Ptr(int32(5)),
		},
	}

	resp, err := service.SetProperties(ctx, getResp.ServiceProperties, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(time.Second * 45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *getResp.Logging.Read, *received.Logging.Read)
	require.Equal(t, *getResp.Logging.Write, *received.Logging.Write)
	require.Equal(t, *getResp.Logging.Delete, *received.Logging.Delete)
	require.Equal(t, *getResp.Logging.RetentionPolicy.Enabled, *received.Logging.RetentionPolicy.Enabled)
	require.Equal(t, *getResp.Logging.RetentionPolicy.Days, *received.Logging.RetentionPolicy.Days)
}

func TestSetHoursMetrics(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, tracing.Provider{})

	getResp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	getResp.HourMetrics = &Metrics{
		Enabled:     to.Ptr(true),
		IncludeAPIs: to.Ptr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.Ptr(true),
			Days:    to.Ptr(int32(5)),
		},
		Version: to.Ptr("1.0"),
	}

	resp, err := service.SetProperties(ctx, getResp.ServiceProperties, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(time.Second * 45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *getResp.HourMetrics.Enabled, *received.HourMetrics.Enabled)
	require.Equal(t, *getResp.HourMetrics.IncludeAPIs, *received.HourMetrics.IncludeAPIs)
	require.Equal(t, *getResp.HourMetrics.RetentionPolicy.Days, *received.HourMetrics.RetentionPolicy.Days)
	require.Equal(t, *getResp.HourMetrics.RetentionPolicy.Enabled, *received.HourMetrics.RetentionPolicy.Enabled)
}

func TestSetMinuteMetrics(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, tracing.Provider{})

	getResp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	getResp.MinuteMetrics = &Metrics{
		Enabled:     to.Ptr(true),
		IncludeAPIs: to.Ptr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.Ptr(true),
			Days:    to.Ptr(int32(5)),
		},
		Version: to.Ptr("1.0"),
	}

	resp, err := service.SetProperties(ctx, getResp.ServiceProperties, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(time.Second * 45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *getResp.MinuteMetrics.Enabled, *received.MinuteMetrics.Enabled)
	require.Equal(t, *getResp.MinuteMetrics.IncludeAPIs, *received.MinuteMetrics.IncludeAPIs)
	require.Equal(t, *getResp.MinuteMetrics.RetentionPolicy.Days, *received.MinuteMetrics.RetentionPolicy.Days)
	require.Equal(t, *getResp.MinuteMetrics.RetentionPolicy.Enabled, *received.MinuteMetrics.RetentionPolicy.Enabled)
}

func TestSetCors(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, tracing.Provider{})

	getResp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	getResp.Cors = []*CorsRule{
		{
			AllowedHeaders:  to.Ptr("x-ms-meta-data*"),
			AllowedMethods:  to.Ptr("PUT"),
			AllowedOrigins:  to.Ptr("www.xyz.com"),
			ExposedHeaders:  to.Ptr("x-ms-meta-source*"),
			MaxAgeInSeconds: to.Ptr(int32(500)),
		},
	}

	resp, err := service.SetProperties(ctx, getResp.ServiceProperties, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(time.Second * 45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *getResp.Cors[0].AllowedHeaders, *received.Cors[0].AllowedHeaders)
	require.Equal(t, *getResp.Cors[0].AllowedMethods, *received.Cors[0].AllowedMethods)
	require.Equal(t, *getResp.Cors[0].AllowedOrigins, *received.Cors[0].AllowedOrigins)
	require.Equal(t, *getResp.Cors[0].ExposedHeaders, *received.Cors[0].ExposedHeaders)
	require.Equal(t, *getResp.Cors[0].MaxAgeInSeconds, *received.Cors[0].MaxAgeInSeconds)
}

func TestSetTooManyCors(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, tracing.Provider{})

	corsRules1 := CorsRule{
		AllowedHeaders:  to.Ptr("x-ms-meta-data*"),
		AllowedMethods:  to.Ptr("PUT"),
		AllowedOrigins:  to.Ptr("www.xyz.com"),
		ExposedHeaders:  to.Ptr("x-ms-meta-source*"),
		MaxAgeInSeconds: to.Ptr(int32(500)),
	}
	props := ServiceProperties{Cors: make([]*CorsRule, 0)}
	for i := 0; i < 6; i++ {
		props.Cors = append(props.Cors, &corsRules1)
	}

	_, err := service.SetProperties(ctx, props, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))
}

func TestRetentionTooLong(t *testing.T) {
	service := initServiceTest(t, storageEndpoint, tracing.Provider{})

	metrics := Metrics{
		Enabled:     to.Ptr(true),
		IncludeAPIs: to.Ptr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.Ptr(true),
			Days:    to.Ptr(int32(366)),
		},
		Version: to.Ptr("1.0"),
	}
	props := ServiceProperties{MinuteMetrics: &metrics}

	_, err := service.SetProperties(ctx, props, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Contains(t, PossibleTableErrorCodeValues(), TableErrorCode(httpErr.ErrorCode))
}

func TestGetAccountSASToken(t *testing.T) {
	cred, err := NewSharedKeyCredential("myAccountName", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	require.NoError(t, err)
	service, err := NewServiceClientWithSharedKey("https://myAccountName.table.core.windows.net", cred, nil)
	require.NoError(t, err)

	resources := AccountSASResourceTypes{Service: true}
	perms := AccountSASPermissions{Read: true}
	start := time.Date(2021, time.September, 8, 14, 30, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	sas, err := service.GetAccountSASURL(resources, perms, start, end)
	require.NoError(t, err)
	require.Equal(t, "https://myAccountName.table.core.windows.net/?se=2021-09-09T14%3A30%3A00Z&sig=m%2F%2FxhMvxidHaswzZRpyuiHykqnTppPi%2BQ9S5xHMksIQ%3D&sp=r&spr=https&srt=s&ss=t&st=2021-09-08T14%3A30%3A00Z&sv=2019-02-02", sas)
}

func TestGetAccountSASTokenError(t *testing.T) {
	cred := NewFakeCredential("fakeaccount", "fakekey")
	service, err := NewServiceClient("https://myAccountName.table.core.windows.net", cred, nil)
	require.NoError(t, err)

	resources := AccountSASResourceTypes{Service: true}
	perms := AccountSASPermissions{Read: true}

	_, err = service.GetAccountSASURL(resources, perms, time.Now(), time.Now().Add(time.Hour))
	require.Error(t, err)
}

type tokenCredFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (t tokenCredFunc) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if l := len(tro.Scopes); l != 1 {
		return azcore.AccessToken{}, fmt.Errorf("unexpected scopes len %d", l)
	}
	return t(ctx, tro)
}

type fakeTransport struct{}

func (fakeTransport) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Request:    req,
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
		Header:     http.Header{},
	}, nil
}

func TestNewServiceClient_sovereignClouds(t *testing.T) {
	tests := []struct {
		label    string
		endpoint string
		scope    string
		cfg      cloud.Configuration
	}{
		{
			label:    "storage China",
			endpoint: "https://myAccountName.table.core.windows.net",
			scope:    "https://storage.azure.cn/.default",
			cfg:      cloud.AzureChina,
		},
		{
			label:    "cosmos China",
			endpoint: "https://myAccountName.table.cosmos.windows.net",
			scope:    "https://cosmos.azure.cn/.default",
			cfg:      cloud.AzureChina,
		},
		{
			label:    "storage USGov",
			endpoint: "https://myAccountName.table.core.windows.net",
			scope:    "https://storage.azure.us/.default",
			cfg:      cloud.AzureGovernment,
		},
		{
			label:    "cosmos USGov",
			endpoint: "https://myAccountName.table.cosmos.windows.net",
			scope:    "https://cosmos.azure.us/.default",
			cfg:      cloud.AzureGovernment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			client, err := NewServiceClient(tt.endpoint, tokenCredFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				if s := tro.Scopes[0]; s != tt.scope {
					return azcore.AccessToken{}, fmt.Errorf("incorrect scope %s", s)
				}
				return azcore.AccessToken{Token: "fake_token", ExpiresOn: time.Now().Add(time.Hour)}, nil
			}), &ClientOptions{
				ClientOptions: policy.ClientOptions{
					Cloud:     tt.cfg,
					Transport: &fakeTransport{},
				},
			})
			require.NoError(t, err)

			// we just call some API so that the pipeline is triggered which will call GetToken on our fake cred
			_, err = client.DeleteTable(context.Background(), "fake-table", nil)
			require.NoError(t, err)
		})
	}
}
