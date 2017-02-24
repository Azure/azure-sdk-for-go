package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	headerAccept          = "Accept"
	headerEtag            = "Etag"
	headerPrefer          = "Prefer"
	headerXmsContinuation = "x-ms-Continuation-NextTableName"
)

// TableServiceClient contains operations for Microsoft Azure Table Storage
// Service.
type TableServiceClient struct {
	client Client
	auth   authentication
}

// TableQueryResult contains the response from
// QueryTables and QueryTablesNextResults functions.
type TableQueryResult struct {
	OdataMetadata string  `json:"odata.metadata"`
	Tables        []Table `json:"value"`
	NextLink      *string
}

// GetTableReference returns a Table object for the specified table name.
func (t *TableServiceClient) GetTableReference(name string) Table {
	return Table{
		tsc:  t,
		Name: name,
	}
}

// QueryTables returns the tables in the storage account.
// You can use query options defined by the OData Protocol specification.
//
// See https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/query-tables
func (t *TableServiceClient) QueryTables(odataQuery url.Values) (*TableQueryResult, error) {
	uri := t.client.getEndpoint(tableServiceName, tablesURIPath, fixOdataQuery(odataQuery))
	return t.queryTables(uri)

}

// QueryTablesNextResults returns the next page of results
// from a QueryTables or a QueryTablesNextResults operation.
//
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/query-tables
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/query-timeout-and-pagination
func (t *TableServiceClient) QueryTablesNextResults(results *TableQueryResult) (*TableQueryResult, error) {
	if results.NextLink == nil {
		return nil, errors.New("There are no more pages in this query results")
	}
	return t.queryTables(*results.NextLink)
}

func (t *TableServiceClient) queryTables(uri string) (*TableQueryResult, error) {
	headers := t.client.getStandardHeaders()
	headers[headerAccept] = "application/json;odata=fullmetadata"

	resp, err := t.client.execInternalJSON(http.MethodGet, uri, headers, nil, t.auth)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return nil, err
	}
	var out TableQueryResult
	err = json.Unmarshal(respBody, &out)
	if err != nil {
		return nil, err
	}

	for i := range out.Tables {
		out.Tables[i].tsc = t
	}

	nextLink := resp.headers.Get(http.CanonicalHeaderKey(headerXmsContinuation))
	if nextLink == "" {
		out.NextLink = nil
	} else {
		originalURI, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		v := originalURI.Query()
		v.Set(nextTableQueryParameter, nextLink)
		newURI := t.client.getEndpoint(tableServiceName, tablesURIPath, v)
		out.NextLink = &newURI
	}

	return &out, nil
}

func addBodyRelatedHeaders(h map[string]string, length int) map[string]string {
	h[headerContentType] = "application/json"
	h[headerContentLength] = fmt.Sprintf("%v", length)
	return h
}

func addReturnContentHeaders(h map[string]string, ml MetadataLevel) map[string]string {
	if ml != EmptyPayload {
		h[headerPrefer] = "return-content"
		h[headerAccept] = string(ml)
	} else {
		h[headerPrefer] = "return-no-content"
	}
	return h
}

// GetServiceProperties gets the properties of your storage account's table service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/get-table-service-properties
func (t *TableServiceClient) GetServiceProperties() (*ServiceProperties, error) {
	return t.client.getServiceProperties(tableServiceName, t.auth)
}

// SetServiceProperties sets the properties of your storage account's table service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/set-table-service-properties
func (t *TableServiceClient) SetServiceProperties(props ServiceProperties) error {
	return t.client.setServiceProperties(props, tableServiceName, t.auth)
}
