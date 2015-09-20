package storage

import (
	//	"encoding/xml"
	"bytes"
	//	"encoding/json"
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
		"x-ms-version":          c.client.apiVersion,
		"x-ms-date":             currentTimeRfc1123Formatted(),
		"accept":                "application/json;odata=nometadata",
		"Accept-Charset":        "UTF-8",
		"DataServiceVersion":    "1.0;NetFx",
		"MaxDataServiceVersion": "2.0;NetFx",
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
		return nil, err
	} else {
		return nil, nil
	}
}

func (c *TableServiceClient) CreateTable(tableName string) error {
	uri := c.client.getEndpoint(tableServiceName, TablesURIPath, url.Values{})

	headers := c.getStandardHeaders()
	headers["accept"] = "application/atom+xml,application/xml"
	headers["Content-Type"] = "application/atom+xml"

	//	req := createTableRequest{TableName: tableName}

	buf := new(bytes.Buffer)
	//	if err := json.NewEncoder(buf).Encode(req); err != nil {
	//		return err
	//	}

	fmt.Fprintf(buf, "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?><entry xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" xmlns=\"http://www.w3.org/2005/Atom\"> <title /> <updated></updated> <author><name/></author><id/><content type=\"application/xml\"><m:properties> <d:TableName>%s</d:TableName> </m:properties> </content> </entry>", tableName)

	log.Printf(string(buf.Bytes()))

	headers["Content-Length"] = string(buf.Len())
	buf.Reset()

	resp, err := c.client.execLite("POST", uri, headers, buf)
	log.Printf("resp == %s", resp)
	log.Printf("resp.statusCode == %d", resp.statusCode)
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
