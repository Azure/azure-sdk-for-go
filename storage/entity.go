package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"io/ioutil"

	"github.com/satori/uuid"
)

// Annotating as secure for gas scanning
/* #nosec */
const (
	partitionKeyNode  = "PartitionKey"
	rowKeyNode        = "RowKey"
	etagErrorTemplate = "Etag didn't match: %v"
)

// Entity represents an entity inside an Azure table.
type Entity struct {
	tsc           *TableServiceClient
	Table         *Table
	PartitionKey  string
	RowKey        string
	TimeStamp     time.Time
	OdataMetadata string
	OdataType     string
	OdataID       string
	OdataEtag     string
	OdataEditLink string
	Properties    map[string]interface{}
}

// GetEntityReference returns an Entity object with the specified
// partition key and row key.
func (t *Table) GetEntityReference(partitionKey, rowKey string) Entity {
	return Entity{
		PartitionKey: partitionKey,
		RowKey:       rowKey,
		Table:        t,
		tsc:          t.tsc,
	}
}

// Insert inserts the referenced entity in its table.
// The function fails if there is an entity with the same
// PartitionKey and RowKey in the table.
// See: https://docs.microsoft.com/rest/api/storageservices/fileservices/insert-entity
func (e *Entity) Insert(getResponse bool) error {
	uri := e.tsc.client.getEndpoint(tableServiceName, e.Table.buildPath(), nil)

	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	headers := e.tsc.client.getStandardHeaders()
	addBodyRelatedHeaders(&headers, len(body))
	if getResponse {
		headers[headerPrefer] = "return-content"
		headers[headerAccept] = "application/json;odata=fullmetadata"
	} else {
		headers[headerPrefer] = "return-no-content"
	}

	resp, err := e.tsc.client.execInternalJSON(http.MethodPost, uri, headers, bytes.NewReader(body), e.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	data, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return err
	}

	if getResponse {
		if err = checkRespCode(resp.statusCode, []int{http.StatusCreated}); err != nil {
			return err
		}
		if err = e.UnmarshalJSON(data); err != nil {
			return err
		}
	} else {
		if err = checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
			return err
		}
	}

	return nil
}

// Update updates the contents of an entity. The function fails if there is no entity
// with the same PartitionKey and RowKey in the table or if the ETag is different
// than the one in Azure.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/update-entity2
func (e *Entity) Update(force bool) error {
	return e.updateMerge(force, http.MethodPut)
}

// Merge merges the contents of entity specified with PartitionKey and RowKey
// with the content specified in Properties.
// The function fails if there is no entity with the same PartitionKey and
// RowKey in the table or if the ETag is different than the one in Azure.
// Read more: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/merge-entity
func (e *Entity) Merge(force bool) error {
	return e.updateMerge(force, "MERGE")
}

// Delete deletes the entity.
// The function fails if there is no entity with the same PartitionKey and
// RowKey in the table or if the ETag is different than the one in Azure.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/delete-entity1
func (e *Entity) Delete(force bool) error {
	uri := e.tsc.client.getEndpoint(tableServiceName, e.buildPath(), nil)

	headers := e.tsc.client.getStandardHeaders()
	addIfMatchHeader(&headers, force, e.OdataEtag)

	resp, err := e.tsc.client.execInternalJSON(http.MethodDelete, uri, headers, nil, e.tsc.auth)
	if err != nil {
		if resp.statusCode == http.StatusPreconditionFailed {
			return fmt.Errorf(etagErrorTemplate, err)
		}
		return err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err
	}

	return e.updateTimestamp(resp.headers)
}

// InsertOrReplace inserts an entity or replaces the existing one.
// Read more: https://docs.microsoft.com/rest/api/storageservices/fileservices/insert-or-replace-entity
func (e *Entity) InsertOrReplace() error {
	return e.insertOr(http.MethodPut)
}

// InsertOrMerge inserts an entity or merges the existing one.
// Read more: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/insert-or-merge-entity
func (e *Entity) InsertOrMerge() error {
	return e.insertOr("MERGE")
}

func (e *Entity) buildPath() string {
	return fmt.Sprintf("%s(PartitionKey='%s', RowKey='%s')", e.Table.buildPath(), e.PartitionKey, e.RowKey)
}

// MarshalJSON is a custom marshaller for entity
func (e *Entity) MarshalJSON() ([]byte, error) {
	completeMap := map[string]interface{}{}
	completeMap[partitionKeyNode] = e.PartitionKey
	completeMap[rowKeyNode] = e.RowKey
	for k, v := range e.Properties {
		typeKey := strings.Join([]string{k, OdataTypeSuffix}, "")
		switch t := v.(type) {
		case []byte:
			completeMap[typeKey] = OdataBinary
			completeMap[k] = string(t)
		case time.Time:
			completeMap[typeKey] = OdataDateTime
			completeMap[k] = t.Format(time.RFC3339Nano)
		case uuid.UUID:
			completeMap[typeKey] = OdataGUID
			completeMap[k] = t.String()
		case int64:
			completeMap[typeKey] = OdataInt64
			completeMap[k] = fmt.Sprintf("%v", v)
		default:
			completeMap[k] = v
		}
		if strings.HasSuffix(k, OdataTypeSuffix) {
			if !(completeMap[k] == OdataBinary ||
				completeMap[k] == OdataDateTime ||
				completeMap[k] == OdataGUID ||
				completeMap[k] == OdataInt64) {
				return nil, fmt.Errorf("Odata.type annotation %v value is not valid", k)
			}
			valueKey := strings.TrimSuffix(k, OdataTypeSuffix)
			if _, ok := completeMap[valueKey]; !ok {
				return nil, fmt.Errorf("Odata.type annotation %v defined without value defined", k)
			}
		}
	}
	return json.Marshal(completeMap)
}

// UnmarshalJSON is a custom unmarshaller for entities
func (e *Entity) UnmarshalJSON(data []byte) error {
	errorTemplate := "Deserializing error: %v"

	props := map[string]interface{}{}
	err := json.Unmarshal(data, &props)
	if err != nil {
		return err
	}

	// deselialize metadata
	e.OdataMetadata = getAndDelete(&props, "odata.metadata")
	e.OdataType = getAndDelete(&props, "odata.type")
	e.OdataID = getAndDelete(&props, "odata.id")
	e.OdataEtag = getAndDelete(&props, "odata.etag")
	e.OdataEditLink = getAndDelete(&props, "odata.editLink")
	e.PartitionKey = getAndDelete(&props, partitionKeyNode)
	e.RowKey = getAndDelete(&props, rowKeyNode)

	// deserialize timestamp
	str, ok := props["Timestamp"].(string)
	if !ok {
		return fmt.Errorf(errorTemplate, "Timestamp casting error")
	}
	t, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	e.TimeStamp = t
	delete(props, "Timestamp")
	delete(props, "Timestamp@odata.type")

	// deserialize entity (user defined fields)
	for k, v := range props {
		if strings.HasSuffix(k, OdataTypeSuffix) {
			valueKey := strings.TrimSuffix(k, OdataTypeSuffix)
			str, ok := props[valueKey].(string)
			if !ok {
				return fmt.Errorf(errorTemplate, fmt.Sprintf("%v casting error", v))
			}
			switch v {
			case OdataBinary:
				props[valueKey] = []byte(str)
			case OdataDateTime:
				t, err := time.Parse("2006-01-02T15:04:05Z", str)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				props[valueKey] = t
			case OdataGUID:
				props[valueKey] = uuid.FromStringOrNil(str)
			case OdataInt64:
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return fmt.Errorf(errorTemplate, err)
				}
				props[valueKey] = i
			default:
				return fmt.Errorf(errorTemplate, fmt.Sprintf("%v is not supported", v))
			}
			delete(props, k)
		}
	}

	e.Properties = props
	return nil
}

func getAndDelete(props *map[string]interface{}, key string) string {
	str, _ := (*props)[key].(string)
	delete(*props, key)
	return str
}

func addBodyRelatedHeaders(h *map[string]string, length int) {
	(*h)[headerContentType] = "application/json"
	(*h)[headerContentLength] = fmt.Sprintf("%v", length)
}

func addIfMatchHeader(h *map[string]string, force bool, etag string) {
	if force {
		(*h)[headerIfMatch] = "*"
	} else {
		(*h)[headerIfMatch] = etag
	}
}

// updates Etag and timestamp
func (e *Entity) updateEtagAndTimestamp(headers http.Header) error {
	e.OdataEtag = headers.Get(headerEtag)
	return e.updateTimestamp(headers)
}

func (e *Entity) updateTimestamp(headers http.Header) error {
	str := headers.Get(headerDate)
	t, err := time.Parse(time.RFC1123, str)
	if err != nil {
		return fmt.Errorf("Update timestamp error: %v", err)
	}
	e.TimeStamp = t
	return nil
}

func (e *Entity) insertOr(verb string) error {
	uri := e.tsc.client.getEndpoint(tableServiceName, e.buildPath(), nil)

	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	headers := e.tsc.client.getStandardHeaders()
	addBodyRelatedHeaders(&headers, len(body))

	resp, err := e.tsc.client.execInternalJSON(verb, uri, headers, bytes.NewReader(body), e.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err
	}

	return e.updateEtagAndTimestamp(resp.headers)
}

func (e *Entity) updateMerge(force bool, verb string) error {
	uri := e.tsc.client.getEndpoint(tableServiceName, e.buildPath(), nil)

	body, err := json.Marshal(e)
	if err != nil {
		return err
	}

	headers := e.tsc.client.getStandardHeaders()
	addBodyRelatedHeaders(&headers, len(body))
	addIfMatchHeader(&headers, force, e.OdataEtag)

	resp, err := e.tsc.client.execInternalJSON(verb, uri, headers, bytes.NewReader(body), e.tsc.auth)
	if err != nil {
		if resp.statusCode == http.StatusPreconditionFailed {
			return fmt.Errorf(etagErrorTemplate, err)
		}
		return err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err
	}

	return e.updateEtagAndTimestamp(resp.headers)
}
