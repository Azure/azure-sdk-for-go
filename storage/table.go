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

func pathForTable(queue string) string { return fmt.Sprintf("/%s", queue) }

//func pathForQueueMessages(queue string) string { return fmt.Sprintf("/%s/messages", queue) }
//func pathForMessage(queue, name string) string { return fmt.Sprintf("/%s/messages/%s", queue, name) }

func (c *TableServiceClient) getStandardHeaders() map[string]string {
	return map[string]string{
		"x-ms-version":   "2015-02-21",
		"x-ms-date":      currentTimeRfc1123Formatted(),
		"accept":         "application/json;odata=nometadata",
		"Accept-Charset": "UTF-8",
		//		"DataServiceVersion":    "1.0;NetFx",
		//		"MaxDataServiceVersion": "2.0;NetFx",
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
	
	log.Printf("%s", string(buf.Bytes()))

	// Reformat error
	return nil, nil

}

func (c *TableServiceClient) CreateTable(tableName string) error {
	uri := c.client.getEndpoint(tableServiceName, TablesURIPath, url.Values{})

	headers := c.getStandardHeaders()
	headers["Content-Type"] = "application/json"

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
