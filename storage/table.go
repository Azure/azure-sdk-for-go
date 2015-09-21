package storage

import (
	"encoding/xml"
	"time"
	"bytes"
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
		"x-ms-version":          "2015-02-21",
		"x-ms-date":             currentTimeRfc1123Formatted(),
		"accept":                "application/atom+xml",
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
	headers["accept"] = "application/atom+xml"
	headers["Content-Type"] = "application/atom+xml"

	buf := new(bytes.Buffer)
	
	nowBytes, err := xml.Marshal(time.Now())
	if err != nil {
		log.Printf("Cannot convert time.Now() to XML datetime: %s", err.Error())
		return err
	}

	fmt.Fprintf(buf, "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<entry xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\" xmlns=\"http://www.w3.org/2005/Atom\"> <title /> <updated>%s</updated> <author><name/></author><id/><content type=\"application/xml\"><m:properties> <d:TableName>%s</d:TableName> </m:properties> </content> </entry>", string(nowBytes), tableName)

//	log.Printf(string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := c.client.execLite("POST", uri, headers, buf)
	
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
