// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

type TableTransactionActionType string

const (
	Add           TableTransactionActionType = "add"
	UpdateMerge   TableTransactionActionType = "updatemerge"
	UpdateReplace TableTransactionActionType = "updatereplace"
	Delete        TableTransactionActionType = "delete"
	UpsertMerge   TableTransactionActionType = "upsertmerge"
	UpsertReplace TableTransactionActionType = "upsertreplace"
)

const (
	headerContentType             = "Content-Type"
	headerContentTransferEncoding = "Content-Transfer-Encoding"
)

type OdataErrorMessage struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type OdataError struct {
	Code    string            `json:"code"`
	Message OdataErrorMessage `json:"message"`
}

type TableError struct {
	OdataError OdataError `json:"odata.error"`
}

func (e *TableError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.OdataError.Code, e.OdataError.Message.Value)
}

type TableTransactionAction struct {
	ActionType TableTransactionActionType
	Entity     map[string]interface{}
	ETag       string
}

type TableTransactionResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// PreferenceApplied contains the information returned from the Preference-Applied header response.
	PreferenceApplied *string

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// The response for a single table.
	TransactionResponses *[]azcore.Response

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// ContentType contains the information returned from the Content-Type header response.
	ContentType *string
}

type TableSubmitTransactionOptions struct {
	RequestID *string
}

var defaultChangesetHeaders = map[string]string{
	"Accept":       "application/json;odata=minimalmetadata",
	"Content-Type": "application/json",
	"Prefer":       "return-no-content",
}

// SubmitTransaction submits the table transactional batch according to the slice of TableTransactionActions provided.
func (t *TableClient) SubmitTransaction(transactionActions []TableTransactionAction, tableSubmitTransactionOptions *TableSubmitTransactionOptions, ctx context.Context) (TableTransactionResponse, error) {
	return t.submitTransactionInternal(&transactionActions, uuid.New(), uuid.New(), tableSubmitTransactionOptions, ctx)
}

// submitTransactionInternal is the internal implementation for SubmitTransaction. It allows for explicit configuration of the batch and changeset UUID values for testing.
func (t *TableClient) submitTransactionInternal(transactionActions *[]TableTransactionAction, batchUuid uuid.UUID, changesetUuid uuid.UUID, tableSubmitTransactionOptions *TableSubmitTransactionOptions, ctx context.Context) (TableTransactionResponse, error) {

	changesetBoundary := fmt.Sprintf("changeset_%s", changesetUuid.String())
	changeSetBody, err := t.generateChangesetBody(changesetBoundary, transactionActions)
	if err != nil {
		return TableTransactionResponse{}, err
	}
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(t.client.con.Endpoint(), "$batch"))
	if err != nil {
		return TableTransactionResponse{}, err
	}
	req.Header.Set("x-ms-version", "2019-02-02")
	if tableSubmitTransactionOptions != nil && tableSubmitTransactionOptions.RequestID != nil {
		req.Header.Set("x-ms-client-request-id", *tableSubmitTransactionOptions.RequestID)
	}
	req.Header.Set("DataServiceVersion", "3.0")
	req.Header.Set("Accept", string(OdataMetadataFormatApplicationJSONOdataMinimalmetadata))

	boundary := fmt.Sprintf("batch_%s", batchUuid.String())
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)
	h := make(textproto.MIMEHeader)
	h.Set(headerContentType, fmt.Sprintf("multipart/mixed; boundary=%s", changesetBoundary))
	batchWriter, err := writer.CreatePart(h)
	if err != nil {
		return TableTransactionResponse{}, err
	}
	batchWriter.Write(changeSetBody.Bytes())
	writer.Close()

	req.SetBody(azcore.NopCloser(bytes.NewReader(body.Bytes())), fmt.Sprintf("multipart/mixed; boundary=%s", boundary))

	resp, err := t.client.con.Pipeline().Do(req)
	if err != nil {
		return TableTransactionResponse{}, err
	}

	transactionResponse, err := buildTransactionResponse(req, resp, len(*transactionActions))
	if err != nil {
		return TableTransactionResponse{}, err
	}

	if !resp.HasStatusCode(http.StatusAccepted, http.StatusNoContent) {
		return TableTransactionResponse{}, azcore.NewResponseError(err, resp.Response)
	}
	return transactionResponse, nil
}

func buildTransactionResponse(req *azcore.Request, resp *azcore.Response, itemCount int) (TableTransactionResponse, error) {
	innerResponses := make([]azcore.Response, itemCount)
	result := TableTransactionResponse{RawResponse: resp.Response, TransactionResponses: &innerResponses}

	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return TableTransactionResponse{}, err
		}
		result.Date = &date
	}

	if val := resp.Header.Get("Preference-Applied"); val != "" {
		result.PreferenceApplied = &val
	}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TableTransactionResponse{}, err
	}
	reader := bytes.NewReader(bytesBody)
	if bytes.IndexByte(bytesBody, '{') == 0 {
		// This is a failure and the body is json
		return TableTransactionResponse{}, errors.New(string(bytesBody))
	}
	outerBoundary := getBoundaryName(bytesBody)
	mpReader := multipart.NewReader(reader, outerBoundary)
	outerPart, err := mpReader.NextPart()
	innerBytes, err := ioutil.ReadAll(outerPart)
	innerBoundary := getBoundaryName(innerBytes)
	reader = bytes.NewReader(innerBytes)
	mpReader = multipart.NewReader(reader, innerBoundary)
	i := 0
	innerPart, err := mpReader.NextPart()
	for ; err == nil; innerPart, err = mpReader.NextPart() {
		part, err := ioutil.ReadAll(innerPart)
		if err != nil {
			break
		}
		r, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(part)), req.Request)
		if err != nil {
			return TableTransactionResponse{}, err
		}
		if r.StatusCode >= 400 {
			errorBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return TableTransactionResponse{}, err
			} else {
				return TableTransactionResponse{}, errors.New(string(errorBody))
			}
		}
		innerResponses[i] = azcore.Response{Response: r}
		i++
	}

	return result, nil
}

func getBoundaryName(bytesBody []byte) string {
	end := bytes.Index(bytesBody, []byte("\n"))
	if end > 0 && bytesBody[end-1] == '\r' {
		end -= 1
	}
	return string(bytesBody[2:end])
}

// transactionHandleError handles the SubmitTransaction error response.
func (client *tableClient) transactionHandleError(errorBody error) error {
	oe := TableError{}
	b := []byte(errorBody.Error())
	if err := json.Unmarshal(b, &oe); err == nil {
		return &oe
	}
	return errors.New("Unknown error.")
}

// generateChangesetBody generates the individual changesets for the various operations within the batch request.
// There is a changeset for Insert, Delete, Merge etc.
func (t *TableClient) generateChangesetBody(changesetBoundary string, transactionActions *[]TableTransactionAction) (*bytes.Buffer, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.SetBoundary(changesetBoundary)

	for _, be := range *transactionActions {
		err := t.generateEntitySubset(&be, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	return body, nil
}

// generateEntitySubset generates body payload for particular batch entity
func (t *TableClient) generateEntitySubset(transactionAction *TableTransactionAction, writer *multipart.Writer) error {

	h := make(textproto.MIMEHeader)
	h.Set(headerContentTransferEncoding, "binary")
	h.Set(headerContentType, "application/http")
	qo := &QueryOptions{Format: OdataMetadataFormatApplicationJSONOdataMinimalmetadata.ToPtr()}

	operationWriter, err := writer.CreatePart(h)
	if err != nil {
		return err
	}
	var req *azcore.Request
	var entity map[string]interface{} = transactionAction.Entity

	if _, ok := entity[PartitionKey]; !ok {
		return fmt.Errorf("entity properties must contain a %s property", PartitionKey)
	}
	if _, ok := entity[RowKey]; !ok {
		return fmt.Errorf("entity properties must contain a %s property", RowKey)
	}
	// Consider empty ETags as '*'
	if len(transactionAction.ETag) == 0 {
		transactionAction.ETag = "*"
	}

	switch transactionAction.ActionType {
	case Delete:
		req, err = t.client.deleteEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), transactionAction.ETag, &TableDeleteEntityOptions{}, qo)
	case Add:
		toOdataAnnotatedDictionary(&entity)
		req, err = t.client.insertEntityCreateRequest(ctx, t.name, &TableInsertEntityOptions{TableEntityProperties: &entity, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, qo)
	case UpdateMerge:
		fallthrough
	case UpsertMerge:
		toOdataAnnotatedDictionary(&entity)
		opts := &TableMergeEntityOptions{TableEntityProperties: &entity}
		if len(transactionAction.ETag) > 0 {
			opts.IfMatch = &transactionAction.ETag
		}
		req, err = t.client.mergeEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), opts, qo)
		if isCosmosEndpoint(t.client.con.Endpoint()) {
			transformPatchToCosmosPost(req)
		}
	case UpdateReplace:
		fallthrough
	case UpsertReplace:
		toOdataAnnotatedDictionary(&entity)
		req, err = t.client.updateEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), &TableUpdateEntityOptions{TableEntityProperties: &entity, IfMatch: &transactionAction.ETag}, qo)
	}

	urlAndVerb := fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL)
	operationWriter.Write([]byte(urlAndVerb))
	writeHeaders(req.Header, &operationWriter)
	operationWriter.Write([]byte("\r\n")) // additional \r\n is needed per changeset separating the "headers" and the body.
	if req.Body != nil {
		io.Copy(operationWriter, req.Body)
	}

	return nil
}

func writeHeaders(h http.Header, writer *io.Writer) {
	// This way it is guaranteed the headers will be written in a sorted order
	var keys []string
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		(*writer).Write([]byte(fmt.Sprintf("%s: %s\r\n", k, h.Get(k))))
	}
}
