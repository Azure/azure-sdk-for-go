package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// QueueServiceClient contains operations for Microsoft Azure Table Storage
// Service.
type TableServiceClient struct {
	client Client
}

type TableName string

const (
	tablesURIPath                       = "/Tables"
	partitionKeyNode                    = "PartitionKey"
	rowKeyNode                          = "RowKey"
	tag                                 = "table"
	tagIgnore                           = "-"
	continuationTokenPartitionKeyHeader = "X-Ms-Continuation-Nextpartitionkey"
	continuationTokenRowHeader          = "X-Ms-Continuation-Nextrowkey"
)

type createTableRequest struct {
	TableName string `json:"TableName"`
}

type queryTablesResponse struct {
	TableName []struct {
		TableName string `json:"TableName"`
	} `json:"value"`
}

type TableEntry interface {
	PartitionKey() string
	RowKey() string
	SetPartitionKey(string) error
	SetRowKey(string) error
}

type ContinuationToken struct {
	NextPartitionKey string
	NextRowKey       string
}

type getTableEntriesResponse struct {
	Elements []map[string]interface{} `json:"value"`
}

func pathForTable(table TableName) string { return fmt.Sprintf("%s", table) }

func (c *TableServiceClient) getStandardHeaders() map[string]string {
	return map[string]string{
		"x-ms-version":   "2015-02-21",
		"x-ms-date":      currentTimeRfc1123Formatted(),
		"Accept":         "application/json;odata=nometadata",
		"Accept-Charset": "UTF-8",
		"Content-Type":   "application/json",
	}
}

func (c *TableServiceClient) QueryTables() ([]TableName, error) {
	uri := c.client.getEndpoint(tableServiceName, tablesURIPath, url.Values{})

	headers := c.getStandardHeaders()
	headers["Content-Length"] = "0"

	resp, err := c.client.execTable("GET", uri, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.body)

	var respArray queryTablesResponse
	if err := json.Unmarshal(buf.Bytes(), &respArray); err != nil {
		return nil, err
	}

	s := make([]TableName, len(respArray.TableName))
	for i, elem := range respArray.TableName {
		s[i] = TableName(elem.TableName)
	}

	return s, nil
}

func (c *TableServiceClient) CreateTable(tableName TableName) error {
	uri := c.client.getEndpoint(tableServiceName, tablesURIPath, url.Values{})

	headers := c.getStandardHeaders()

	req := createTableRequest{TableName: string(tableName)}
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return err
	}

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execTable("POST", uri, headers, buf)

	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusCreated}); err != nil {
		return err
	} else {
		return nil
	}
}

func injectPartitionAndRowKeys(entry TableEntry, buf *bytes.Buffer) error {
	if err := json.NewEncoder(buf).Encode(entry); err != nil {
		return err
	}

	dec := make(map[string]interface{})
	if err := json.NewDecoder(buf).Decode(&dec); err != nil {
		return err
	}

	// Inject PartitionKey and RowKey
	dec[partitionKeyNode] = entry.PartitionKey()
	dec[rowKeyNode] = entry.RowKey()

	// Remove tagged fields
	// The tag is defined in the const section
	// This is useful to avoid storing the PartitionKey and RowKey twice.
	numFields := reflect.ValueOf(entry).Elem().NumField()
	for i := 0; i < numFields; i++ {
		f := reflect.ValueOf(entry).Elem().Type().Field(i)

		if f.Tag.Get(tag) == tagIgnore {
			// we must look for its JSON name in the dictionary
			// as the user can rename it using a tag
			jsonName := f.Name
			if f.Tag.Get("json") != "" {
				jsonName = f.Tag.Get("json")
			}
			delete(dec, jsonName)
		}
	}

	buf.Reset()

	if err := json.NewEncoder(buf).Encode(&dec); err != nil {
		return err
	}

	return nil
}

func deserializeEntry(retType reflect.Type, reader io.Reader) ([](*TableEntry), error) {
	buf := new(bytes.Buffer)

	var ret getTableEntriesResponse
	if err := json.NewDecoder(reader).Decode(&ret); err != nil {
		return nil, err
	}

	tEntries := make([]*TableEntry, len(ret.Elements))

	for i, entry := range ret.Elements {

		buf.Reset()
		if err := json.NewEncoder(buf).Encode(entry); err != nil {
			return nil, err
		}

		dec := make(map[string]interface{})
		if err := json.NewDecoder(buf).Decode(&dec); err != nil {
			return nil, err
		}

		var pKey, rKey string
		// strip pk and rk
		for key, val := range dec {
			switch {
			case key == partitionKeyNode:
				pKey = val.(string)
			case key == rowKeyNode:
				rKey = val.(string)
			}
		}

		delete(dec, partitionKeyNode)
		delete(dec, rowKeyNode)

		buf.Reset()
		if err := json.NewEncoder(buf).Encode(dec); err != nil {
			return nil, err
		}

		// Create a empty retType instance
		e := reflect.New(retType.Elem()).Interface().(TableEntry)

		// Popolate it with the values
		if err := json.NewDecoder(buf).Decode(&e); err != nil {
			return nil, err
		}

		// Reset PartitionKey and RowKey
		e.SetPartitionKey(pKey)
		e.SetRowKey(rKey)

		// store the pointer
		tEntries[i] = &e

	}

	return tEntries, nil
}

func (c *TableServiceClient) QueryTableEntries(tableName TableName, previousContToken *ContinuationToken, retType reflect.Type, top int, query string) ([](*TableEntry), *ContinuationToken, error) {
	buf := new(bytes.Buffer)

	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("()?$top=%d", top)
	if query != "" {
		uri += fmt.Sprintf("&$filter=%s", url.QueryEscape(query))
	}

	if previousContToken != nil {
		uri += fmt.Sprintf("&NextPartitionKey=%s&NextRowKey=%s", previousContToken.NextPartitionKey, previousContToken.NextRowKey)
	}

	headers := c.getStandardHeaders()
	headers["DataServiceVersion"] = "1.0;NetFx"
	headers["MaxDataServiceVersion"] = "3.0;NetFx"

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execTable("GET", uri, headers, buf)

	contToken := ExtractContinuationTokenFromHeaders(resp.headers)

	//	for key, val := range resp.headers {
	//		log.Printf("headers[%s] = %s", key, val)
	//	}

	if err != nil {
		return nil, contToken, err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return nil, contToken, err
	}

	retEntries, err := deserializeEntry(retType, resp.body)
	if err != nil {
		return nil, contToken, err
	}

	return retEntries, contToken, nil
}

func (c *TableServiceClient) GetTableEntries(tableName TableName, previousContToken *ContinuationToken, retType reflect.Type, top int) ([](*TableEntry), *ContinuationToken, error) {
	buf := new(bytes.Buffer)

	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("()?$top=%d", top)

	if previousContToken != nil {
		uri += fmt.Sprintf("&NextPartitionKey=%s&NextRowKey=%s", previousContToken.NextPartitionKey, previousContToken.NextRowKey)
	}

	headers := c.getStandardHeaders()
	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execTable("GET", uri, headers, buf)

	contToken := ExtractContinuationTokenFromHeaders(resp.headers)

	//	for key, val := range resp.headers {
	//		log.Printf("headers[%s] = %s", key, val)
	//	}

	if err != nil {
		return nil, contToken, err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return nil, contToken, err
	}

	retEntries, err := deserializeEntry(retType, resp.body)
	if err != nil {
		return nil, contToken, err
	}

	return retEntries, contToken, nil
}

func ExtractContinuationTokenFromHeaders(headers map[string][]string) *ContinuationToken {
	if len(headers[continuationTokenPartitionKeyHeader]) > 0 {
		return &ContinuationToken{headers[continuationTokenPartitionKeyHeader][0], headers[continuationTokenRowHeader][0]}
	}
	return nil
}

func (c *TableServiceClient) InsertEntry(tableName TableName, entry TableEntry) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	headers := c.getStandardHeaders()

	buf := new(bytes.Buffer)

	if err := injectPartitionAndRowKeys(entry, buf); err != nil {
		return err
	}

	//	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execTable("POST", uri, headers, buf)

	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusCreated}); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *TableServiceClient) InsertOrReplaceEntry(tableName TableName, entry TableEntry) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("(PartitionKey='%s',RowKey='%s')", url.QueryEscape(entry.PartitionKey()), url.QueryEscape(entry.RowKey()))

	headers := c.getStandardHeaders()

	buf := new(bytes.Buffer)

	if err := injectPartitionAndRowKeys(entry, buf); err != nil {
		return err
	}

	//	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execTable("PUT", uri, headers, buf)

	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err
	} else {
		return nil
	}
}
