package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// QueueServiceClient contains operations for Microsoft Azure Table Storage
// Service.
type TableServiceClient struct {
	client Client
}

const (
	TablesURIPath = "/Tables"
)

type createTableRequest struct {
	TableName string `json:"TableName"`
}

type queryTablesResponse struct {
	TableName []struct {
		TableName string `json:"TableName"`
	} `json:"value"`
}

func pathForTable(table string) string { return fmt.Sprintf("%s/%s", TablesURIPath, table) }

func (c *TableServiceClient) getStandardHeaders() map[string]string {
	return map[string]string{
		"x-ms-version":   "2015-02-21",
		"x-ms-date":      currentTimeRfc1123Formatted(),
		"accept":         "application/json;odata=nometadata",
		"Accept-Charset": "UTF-8",
		"Content-Type":   "application/json",
		//				"DataServiceVersion":    "1.0;NetFx",
		//				"MaxDataServiceVersion": "2.0;NetFx",
	}
}

func (c *TableServiceClient) QueryTables() ([]string, error) {
	uri := c.client.getEndpoint(tableServiceName, TablesURIPath, url.Values{})

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
	uri := c.client.getEndpoint(tableServiceName, TablesURIPath, url.Values{})

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

	//	log.Printf("resp == %s", resp)
	//	log.Printf("resp.statusCode == %d", resp.statusCode)

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

func (c *TableServiceClient) InsertEntity(tableName string, partitionKey string, rowKey string, entity interface{}) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})

	headers := c.getStandardHeaders()

	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(entity); err != nil {
		return err
	}

	dec := make(map[string]interface{})
	if err := json.NewDecoder(buf).Decode(&dec); err != nil {
		return err
	}

	// Inject PartitionKey and RowKey
	dec["PartitionKey"] = partitionKey
	dec["RowKey"] = rowKey

	buf.Reset()

	if err := json.NewEncoder(buf).Encode(&dec); err != nil {
		return err
	}

	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execLite("POST", uri, headers, buf)

	log.Printf("err == %s", err)

	//	log.Printf("resp == %s", resp)
	//	log.Printf("resp.statusCode == %d", resp.statusCode)

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
