package storage

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
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
		"Accept":         "application/json;odata=nometadata",
		"Accept-Charset": "UTF-8",
		"Content-Type":   "application/json",
		"DataServiceVersion": "1.0;NetFx",
		"MaxDataServiceVersion":"2.0;NetFx",
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

func (c *TableServiceClient) InsertEntity(tableName string, partitionKey string, rowKey string, entries map[string]interface{}) error {
	uri := c.client.getEndpoint(tableServiceName, pathForTable(tableName), url.Values{})
	uri += fmt.Sprintf("(PartitionKey='%s', RowKey='%s')", partitionKey, rowKey)

	headers := make(map[string]string)

	headers["Accept"] = "application/atom+xml,application/xml"
	headers["Accept-Charset"] = "UTF-8"
	headers["User-Agent"] = "Microsoft ADO.NET Data Services"
	headers["Content-Type"] = "application/atom+xml"
	
	headers["DataServiceVersion"] = "1.0;NetFx"
	headers["MaxDataServiceVersion"] = "2.0;NetFx"
	
	headers["x-ms-date"] = currentTimeRfc1123Formatted()
	headers["x-ms-version"] = "2015-02-21"

	buf := new(bytes.Buffer)
	if err := xml.NewEncoder(buf).Encode(time.Now()); err != nil {
		log.Printf("XML Encoding of time.Now() failed: %s", err)
		return err
	}
	
	xmlTimestamp := string(buf.Bytes())
	// strip xmlTimeStamp of outer nodes
	xmlTimestamp = xmlTimestamp[6 : len(xmlTimestamp)-7]
	
	sXML := "<?xml version=\"1.0\" encoding=\"utf-8\"?><entry xmlns=\"http://www.w3.org/2005/Atom\" xmlns:d=\"http://schemas.microsoft.com/ado/2007/08/dataservices\" xmlns:m=\"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata\">"
	sXML += fmt.Sprintf("<id /><title /><updated>%s</updated><author><name /></author><content type=\"application/xml\"><m:properties>", xmlTimestamp)

	sXML += fmt.Sprintf("<d:PartitionKey>%s</d:PartitionKey><d:RowKey>%s</d:RowKey><d:Timestamp m:type=\"Edm.DateTime\" m:null=\"true\" />", partitionKey, rowKey)

	for key, val := range entries {
		val = val
		sXML += fmt.Sprintf("<d:%s m:type=\"Edm.String\">%s</d:%s>", key, "oooo", key)
	}
	
	sXML += "</m:properties></content></entry>"


	buf.Reset()
	buf.WriteString(sXML)

	log.Printf("request.body == %s", string(buf.Bytes()))

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())
	 
	resp, err := c.client.execLite("PUT", uri, headers, buf)

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
