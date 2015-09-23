package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

// QueueServiceClient contains operations for Microsoft Azure Table Storage
// Service.
type TableServiceClient struct {
	client Client
}

const (
	tablesURIPath           = "/Tables"
	partitionKeyNode        = "PartitionKey"
	rowKeyNode              = "RowKey"
	tag                     = "table"
	tagIgnore               = "-"
	continuationTokenHeader = "x-ms-continuation-NextTableName"
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

type getTableEntriesResponse struct {
	Elements []map[string]interface{} `json:"value"`
}

type ContinuationToken string

func pathForTable(table string) string { return fmt.Sprintf("%s", table) }

func (c *TableServiceClient) getStandardHeaders() map[string]string {
	return map[string]string{
		"x-ms-version":   "2015-02-21",
		"x-ms-date":      currentTimeRfc1123Formatted(),
		"Accept":         "application/json;odata=nometadata",
		"Accept-Charset": "UTF-8",
		"Content-Type":   "application/json",
	}
}

func (c *TableServiceClient) QueryTables() ([]string, error) {
	uri := c.client.getEndpoint(tableServiceName, tablesURIPath, url.Values{})

	headers := c.getStandardHeaders()
	headers["Content-Length"] = "0"

	resp, err := c.client.execLite("GET", uri, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		log.Printf("resp.body after error == %s \t%s", err.Error(), resp.body)
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.body)

	var respArray queryTablesResponse
	if err := json.Unmarshal(buf.Bytes(), &respArray); err != nil {
		return nil, err
	}

	s := make([]string, len(respArray.TableName))
	for i, elem := range respArray.TableName {
		s[i] = elem.TableName
	}

	return s, nil
}

func (c *TableServiceClient) CreateTable(tableName string) error {
	uri := c.client.getEndpoint(tableServiceName, tablesURIPath, url.Values{})

	headers := c.getStandardHeaders()

	req := createTableRequest{TableName: tableName}
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return err
	}

	log.Printf(string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execLite("POST", uri, headers, buf)

	log.Printf("err == %s", err)

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

		log.Printf("f.Name == %s, f.Tag == %s", f.Name, f.Tag)

		if f.Tag.Get(tag) == tagIgnore {
			log.Printf("\tIgnoring %s", f.Name)
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
		//log.Printf("entry == %s", entry)

		buf.Reset()
		if err := json.NewEncoder(buf).Encode(entry); err != nil {
			return nil, err
		}
		//log.Printf("buf == %s", buf.Bytes())

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

		//		log.Printf("dec == %s", dec)

		buf.Reset()
		if err := json.NewEncoder(buf).Encode(dec); err != nil {
			return nil, err
		}

		e := reflect.New(retType.Elem()).Interface().(TableEntry)

		//		log.Printf("e == %s", e)

		//		log.Printf("buf mangled == %s", buf.Bytes())

		if err := json.NewDecoder(buf).Decode(&e); err != nil {
			return nil, err
		}

		// Reset PartitionKey and RowKey
		e.SetPartitionKey(pKey)
		e.SetRowKey(rKey)

		// store the pointer
		tEntries[i] = &e

		//		log.Printf("e == %s", e)

		//		log.Printf("")
	}

	//	for _, elem := range ret.Elements {
	//		*entriesToPopolate = append(*entriesToPopolate, TableEntry{PartitionKey:elem[partitionKeyNode].(string), RowKey:elem[rowKeyNode].(string)})
	//	}

	//	log.Printf("*entriesToPopolate == %s", *entriesToPopolate)

	return tEntries, nil
}

func (c *TableServiceClient) GetTableEntries(tableName string, previousContToken ContinuationToken, retType reflect.Type) ([](*TableEntry), ContinuationToken, error) {
	buf := new(bytes.Buffer)

	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("()")

	log.Printf("uri == %s ", uri)

	headers := c.getStandardHeaders()
	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())
	if previousContToken != "" {
		log.Printf("Setting continuationNextTableName to %s", string(previousContToken))
		headers["x-ms-continuation-NextTableName"] = string(previousContToken)
	}

	resp, err := c.client.execLite("GET", uri, headers, buf)

	var contToken ContinuationToken
	tcontToken := resp.headers[continuationTokenHeader]

	log.Printf("tcontToken == %s", tcontToken)

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

func (c *TableServiceClient) InsertEntry(tableName string, entry TableEntry) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	headers := c.getStandardHeaders()

	buf := new(bytes.Buffer)

	if err := injectPartitionAndRowKeys(entry, buf); err != nil {
		return err
	}

	//	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execLite("POST", uri, headers, buf)

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

func (c *TableServiceClient) InsertOrReplaceEntry(tableName string, entry TableEntry) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("(PartitionKey='%s',RowKey='%s')", url.QueryEscape(entry.PartitionKey()), url.QueryEscape(entry.RowKey()))

	headers := c.getStandardHeaders()

	buf := new(bytes.Buffer)

	if err := injectPartitionAndRowKeys(entry, buf); err != nil {
		return err
	}

	//	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execLite("PUT", uri, headers, buf)

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
