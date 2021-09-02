// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

const tableNamePrefix = "tableName"

func TestServiceErrorsServiceClient(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service, delete := initServiceTest(t, service)
			defer delete()

			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			_, err = service.CreateTable(context.Background(), tableName, nil)
			require.NoError(t, err)

			// Create a duplicate table to produce an error
			_, err = service.CreateTable(context.Background(), tableName, nil)
			require.Error(t, err)

			_, err = service.DeleteTable(context.Background(), tableName, nil)
			require.NoError(t, err)
		})
	}
}

func TestCreateTableFromService(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service, delete := initServiceTest(t, service)
			defer delete()

			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			_, err = service.CreateTable(ctx, tableName, nil)
			deleteTable := func() {
				_, err := service.DeleteTable(ctx, tableName, nil)
				if err != nil {
					fmt.Printf("Error cleaning up test. %v\n", err.Error())
				}
			}
			defer deleteTable()

			require.NoError(t, err)
			// require.Equal(t, *resp.TableResponse.TableName, tableName)
		})
	}
}

func TestQueryTable(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service, delete := initServiceTest(t, service)
			defer delete()

			tableCount := 5
			tableNames := make([]string, tableCount)
			prefix1 := "zzza"
			prefix2 := "zzzb"

			defer clearAllTables(service) //nolint
			//create 10 tables with our exected prefix and 1 with a different prefix
			for i := 0; i < tableCount; i++ {
				if i < (tableCount - 1) {
					name := fmt.Sprintf("%v%v", prefix1, i)
					tableNames[i] = name
				} else {
					name := fmt.Sprintf("%v%v", prefix2, i)
					tableNames[i] = name
				}
				_, err := service.CreateTable(ctx, tableNames[i], nil)
				require.NoError(t, err)
			}

			// Query for tables with no pagination. The filter should exclude one table from the results
			filter := fmt.Sprintf("TableName ge '%s' and TableName lt '%s'", prefix1, prefix2)
			pager := service.ListTables(&ListTablesOptions{Filter: &filter})

			resultCount := 0
			for pager.NextPage(ctx) {
				resp := pager.PageResponse()
				resultCount += len(resp.Tables)
			}

			require.NoError(t, pager.Err())
			require.Equal(t, resultCount, tableCount-1)

			// Query for tables with pagination
			top := int32(2)
			pager = service.ListTables(&ListTablesOptions{Filter: &filter, Top: &top})

			resultCount = 0
			pageCount := 0
			for pager.NextPage(ctx) {
				resp := pager.PageResponse()
				resultCount += len(resp.Tables)
				pageCount++
			}

			require.NoError(t, pager.Err())
			require.Equal(t, resultCount, tableCount-1)
			require.Equal(t, pageCount, int(top))

		})
	}
}

func TestListTables(t *testing.T) {
	for _, service := range services {
		t.Run(fmt.Sprintf("%v_%v", t.Name(), service), func(t *testing.T) {
			service, delete := initServiceTest(t, service)
			defer delete()
			tableName, err := createRandomName(t, tableNamePrefix)
			require.NoError(t, err)

			err = clearAllTables(service)
			require.NoError(t, err)

			for i := 0; i < 5; i++ {
				_, err := service.CreateTable(ctx, fmt.Sprintf("%v%v", tableName, i), nil)
				require.NoError(t, err)
			}

			count := 0
			pager := service.ListTables(nil)
			for pager.NextPage(ctx) {
				resp := pager.PageResponse()
				count += len(resp.Tables)
			}

			require.NoError(t, pager.Err())
			require.Equal(t, 5, count)

			deleteTable := func() {
				for i := 0; i < 5; i++ {
					_, err := service.DeleteTable(ctx, fmt.Sprintf("%v%v", tableName, i), nil)
					if err != nil {
						fmt.Printf("Error cleaning up test. %v\n", err.Error())
					}
				}
			}
			defer deleteTable()

		})
	}
}

// This functionality is only available on storage accounts
func TestGetStatistics(t *testing.T) {
	t.Skip()
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	accountName := recording.GetEnvVariable(t, "TABLES_STORAGE_ACCOUNT_NAME", "fakestorageaccount")
	serviceURL := storageURI(accountName+"-secondary", "core.windows.net")
	service, err := createServiceClientForRecording(t, serviceURL, cred)
	require.NoError(t, err)

	// s.T().Skip() // TODO: need to change URL to -secondary https://docs.microsoft.com/en-us/rest/api/storageservices/get-table-service-stats
	resp, err := service.GetStatistics(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// Functionality is only available on storage accounts
func TestGetProperties(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	resp, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// Logging is only available on storage accounts
func TestSetLogging(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	logging := Logging{
		Read:    to.BoolPtr(true),
		Write:   to.BoolPtr(true),
		Delete:  to.BoolPtr(true),
		Version: to.StringPtr("1.0"),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
	}
	props := ServiceProperties{Logging: &logging}

	resp, err := service.SetProperties(ctx, props, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *props.Logging.Read, *received.Logging.Read)
	require.Equal(t, *props.Logging.Write, *received.Logging.Write)
	require.Equal(t, *props.Logging.Delete, *received.Logging.Delete)
	require.Equal(t, *props.Logging.RetentionPolicy.Enabled, *received.Logging.RetentionPolicy.Enabled)
	require.Equal(t, *props.Logging.RetentionPolicy.Days, *received.Logging.RetentionPolicy.Days)
}

func TestSetHoursMetrics(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
		Version: to.StringPtr("1.0"),
	}
	props := ServiceProperties{HourMetrics: &metrics}

	resp, err := service.SetProperties(ctx, props, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *props.HourMetrics.Enabled, *received.HourMetrics.Enabled)
	require.Equal(t, *props.HourMetrics.IncludeAPIs, *received.HourMetrics.IncludeAPIs)
	require.Equal(t, *props.HourMetrics.RetentionPolicy.Days, *received.HourMetrics.RetentionPolicy.Days)
	require.Equal(t, *props.HourMetrics.RetentionPolicy.Enabled, *received.HourMetrics.RetentionPolicy.Enabled)
}

func TestSetMinuteMetrics(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(5),
		},
		Version: to.StringPtr("1.0"),
	}
	props := ServiceProperties{MinuteMetrics: &metrics}

	resp, err := service.SetProperties(ctx, props, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *props.MinuteMetrics.Enabled, *received.MinuteMetrics.Enabled)
	require.Equal(t, *props.MinuteMetrics.IncludeAPIs, *received.MinuteMetrics.IncludeAPIs)
	require.Equal(t, *props.MinuteMetrics.RetentionPolicy.Days, *received.MinuteMetrics.RetentionPolicy.Days)
	require.Equal(t, *props.MinuteMetrics.RetentionPolicy.Enabled, *received.MinuteMetrics.RetentionPolicy.Enabled)
}

func TestSetCors(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	corsRules1 := CorsRule{
		AllowedHeaders:  to.StringPtr("x-ms-meta-data*"),
		AllowedMethods:  to.StringPtr("PUT"),
		AllowedOrigins:  to.StringPtr("www.xyz.com"),
		ExposedHeaders:  to.StringPtr("x-ms-meta-source*"),
		MaxAgeInSeconds: to.Int32Ptr(500),
	}
	props := ServiceProperties{Cors: []*CorsRule{&corsRules1}}

	resp, err := service.SetProperties(ctx, props, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	recording.Sleep(45)

	received, err := service.GetProperties(ctx, nil)
	require.NoError(t, err)

	require.Equal(t, *props.Cors[0].AllowedHeaders, *received.Cors[0].AllowedHeaders)
	require.Equal(t, *props.Cors[0].AllowedMethods, *received.Cors[0].AllowedMethods)
	require.Equal(t, *props.Cors[0].AllowedOrigins, *received.Cors[0].AllowedOrigins)
	require.Equal(t, *props.Cors[0].ExposedHeaders, *received.Cors[0].ExposedHeaders)
	require.Equal(t, *props.Cors[0].MaxAgeInSeconds, *received.Cors[0].MaxAgeInSeconds)
}

func TestSetTooManyCors(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	corsRules1 := CorsRule{
		AllowedHeaders:  to.StringPtr("x-ms-meta-data*"),
		AllowedMethods:  to.StringPtr("PUT"),
		AllowedOrigins:  to.StringPtr("www.xyz.com"),
		ExposedHeaders:  to.StringPtr("x-ms-meta-source*"),
		MaxAgeInSeconds: to.Int32Ptr(500),
	}
	props := ServiceProperties{Cors: make([]*CorsRule, 0)}
	for i := 0; i < 6; i++ {
		props.Cors = append(props.Cors, &corsRules1)
	}

	_, err := service.SetProperties(ctx, props, nil)
	require.Error(t, err)
}

func TestRetentionTooLong(t *testing.T) {
	service, delete := initServiceTest(t, "storage")
	defer delete()

	metrics := Metrics{
		Enabled:     to.BoolPtr(true),
		IncludeAPIs: to.BoolPtr(true),
		RetentionPolicy: &RetentionPolicy{
			Enabled: to.BoolPtr(true),
			Days:    to.Int32Ptr(366),
		},
		Version: to.StringPtr("1.0"),
	}
	props := ServiceProperties{MinuteMetrics: &metrics}

	_, err := service.SetProperties(ctx, props, nil)
	require.Error(t, err)
}
