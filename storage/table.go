package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	tablesURIPath                  = "/Tables"
	nextTableQueryParameter        = "NextTableName"
	headerNextPartitionKey         = "x-ms-continuation-NextPartitionKey"
	headerNextRowKey               = "x-ms-continuation-NextRowKey"
	nextPartitionKeyQueryParameter = "NextPartitionKey"
	nextRowKeyQueryParameter       = "NextRowKey"
)

// TableAccessPolicy are used for SETTING table policies
type TableAccessPolicy struct {
	ID         string
	StartTime  time.Time
	ExpiryTime time.Time
	CanRead    bool
	CanAppend  bool
	CanUpdate  bool
	CanDelete  bool
}

// Table represents an Azure table.
type Table struct {
	tsc            *TableServiceClient
	Name           string `json:"TableName"`
	OdataEditLink  string `json:"odata.editLink"`
	OdataID        string `json:"odata.id"`
	OdataMetadata  string `json:"odata.metadata"`
	OdataType      string `json:"odata.type"`
	AccessPolicies []TableAccessPolicy
}

// EntityQueryResult contains the response from
// QueryEntities and QueryEntitiesNextResults functions.
type EntityQueryResult struct {
	OdataMetadata string   `json:"odata.metadata"`
	Entities      []Entity `json:"value"`
	NextLink      *string
}

type continuationToken struct {
	NextPartitionKey string
	NextRowKey       string
}

func (t *Table) buildPath() string {
	return fmt.Sprintf("/%s", t.Name)
}

// Create creates the referenced table.
// This function fails if the name is not compliant
// with the specification or the tables already exists.
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/create-table
func (t *Table) Create(getResponse bool) error {
	uri := t.tsc.client.getEndpoint(tableServiceName, tablesURIPath, nil)

	headers := t.tsc.client.getStandardHeaders()
	headers[headerContentType] = "application/json"
	if getResponse {
		headers[headerPrefer] = "return-content"
		headers[headerAccept] = "application/json;odata=fullmetadata"
	} else {
		headers[headerPrefer] = "return-no-content"
	}

	type createTableRequest struct {
		TableName string `json:"TableName"`
	}
	req := createTableRequest{TableName: t.Name}

	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return err
	}

	headers["Content-Length"] = fmt.Sprintf("%d", buf.Len())

	resp, err := t.tsc.client.execInternalJSON(http.MethodPost, uri, headers, buf, t.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusCreated, http.StatusNoContent}); err != nil {
		return err
	}

	if getResponse {
		data, err := ioutil.ReadAll(resp.body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, t)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes the referenced table.
// This function fails if the table is not present.
// Be advised: Delete deletes all the entries that may be present.
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/delete-table
func (t *Table) Delete() error {
	uri := t.tsc.client.getEndpoint(tableServiceName, tablesURIPath, url.Values{})
	uri += fmt.Sprintf("('%s')", t.Name)

	headers := t.tsc.client.getStandardHeaders()

	resp, err := t.tsc.client.execInternalJSON(http.MethodDelete, uri, headers, nil, t.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err

	}
	return nil
}

// QueryEntities returns the entities in the table.
// You can use query options defined by the OData Protocol specification.
//
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/query-entities
func (t *Table) QueryEntities(odataQuery url.Values) (*EntityQueryResult, error) {
	uri := t.tsc.client.getEndpoint(tableServiceName, t.buildPath(), fixOdataQuery(odataQuery))
	return t.queryEntities(uri)
}

// QueryEntitiesNextResults returns the next page of results
// from a QueryEntities or QueryEntitiesNextResults operation.
//
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/query-entities
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/query-timeout-and-pagination
func (t *Table) QueryEntitiesNextResults(results *EntityQueryResult) (*EntityQueryResult, error) {
	if results.NextLink == nil {
		return nil, errors.New("There are no more pages in this query results")
	}
	return t.queryEntities(*results.NextLink)
}

// SetPermissions sets up table ACL permissions
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/Set-Table-ACL
func (t *Table) SetPermissions(timeout uint) error {
	params := url.Values{"comp": {"acl"}}

	if timeout > 0 {
		params.Add("timeout", fmt.Sprint(timeout))
	}

	uri := t.tsc.client.getEndpoint(tableServiceName, t.Name, params)
	headers := t.tsc.client.getStandardHeaders()

	body, length, err := generateTableACLPayload(t.AccessPolicies)
	if err != nil {
		return err
	}
	headers["Content-Length"] = fmt.Sprintf("%v", length)

	resp, err := t.tsc.client.execInternalJSON(http.MethodPut, uri, headers, body, t.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err := checkRespCode(resp.statusCode, []int{http.StatusNoContent}); err != nil {
		return err
	}
	return nil
}

func generateTableACLPayload(policies []TableAccessPolicy) (io.Reader, int, error) {
	sil := SignedIdentifiers{
		SignedIdentifiers: []SignedIdentifier{},
	}
	for _, tap := range policies {
		permission := generateTablePermissions(&tap)
		signedIdentifier := convertAccessPolicyToXMLStructs(tap.ID, tap.StartTime, tap.ExpiryTime, permission)
		sil.SignedIdentifiers = append(sil.SignedIdentifiers, signedIdentifier)
	}
	return xmlMarshal(sil)
}

// GetPermissions gets the table ACL permissions
// See https://docs.microsoft.com/rest/api/storageservices/fileservices/get-table-acl
func (t *Table) GetPermissions(timeout int) error {
	params := url.Values{"comp": {"acl"}}

	if timeout > 0 {
		params.Add("timeout", strconv.Itoa(timeout))
	}

	uri := t.tsc.client.getEndpoint(tableServiceName, t.Name, params)
	headers := t.tsc.client.getStandardHeaders()

	resp, err := t.tsc.client.execInternalJSON(http.MethodGet, uri, headers, nil, t.tsc.auth)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return err
	}

	var ap AccessPolicy
	err = xmlUnmarshal(resp.body, &ap.SignedIdentifiersList)
	if err != nil {
		return err
	}
	t.updateAccessPolicy(ap)

	return nil
}

func (t *Table) queryEntities(uri string) (*EntityQueryResult, error) {
	headers := t.tsc.client.getStandardHeaders()
	headers[headerAccept] = "application/json;odata=fullmetadata"

	resp, err := t.tsc.client.execInternalJSON(http.MethodGet, uri, headers, nil, t.tsc.auth)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return nil, err
	}
	var entities EntityQueryResult
	err = json.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	for i := range entities.Entities {
		entities.Entities[i].tsc = t.tsc
		entities.Entities[i].table = t
	}

	contToken := extractContinuationTokenFromHeaders(resp.headers)
	if contToken == nil {
		entities.NextLink = nil
	} else {
		originalURI, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		v := originalURI.Query()
		v.Set(nextPartitionKeyQueryParameter, contToken.NextPartitionKey)
		v.Set(nextRowKeyQueryParameter, contToken.NextRowKey)
		newURI := t.tsc.client.getEndpoint(tableServiceName, t.buildPath(), v)
		entities.NextLink = &newURI
	}

	return &entities, nil
}

func extractContinuationTokenFromHeaders(h http.Header) *continuationToken {
	ct := continuationToken{
		NextPartitionKey: h.Get(headerNextPartitionKey),
		NextRowKey:       h.Get(headerNextRowKey),
	}

	if ct.NextPartitionKey != "" && ct.NextRowKey != "" {
		return &ct
	}
	return nil
}

func (t *Table) updateAccessPolicy(ap AccessPolicy) {
	t.AccessPolicies = []TableAccessPolicy{}
	for _, policy := range ap.SignedIdentifiersList.SignedIdentifiers {
		tap := TableAccessPolicy{
			ID:         policy.ID,
			StartTime:  policy.AccessPolicy.StartTime,
			ExpiryTime: policy.AccessPolicy.ExpiryTime,
		}
		tap.CanRead = updatePermissions(policy.AccessPolicy.Permission, "r")
		tap.CanAppend = updatePermissions(policy.AccessPolicy.Permission, "a")
		tap.CanUpdate = updatePermissions(policy.AccessPolicy.Permission, "u")
		tap.CanDelete = updatePermissions(policy.AccessPolicy.Permission, "d")

		t.AccessPolicies = append(t.AccessPolicies, tap)
	}
}

func generateTablePermissions(tap *TableAccessPolicy) (permissions string) {
	// generate the permissions string (raud).
	// still want the end user API to have bool flags.
	permissions = ""

	if tap.CanRead {
		permissions += "r"
	}

	if tap.CanAppend {
		permissions += "a"
	}

	if tap.CanUpdate {
		permissions += "u"
	}

	if tap.CanDelete {
		permissions += "d"
	}
	return permissions
}
