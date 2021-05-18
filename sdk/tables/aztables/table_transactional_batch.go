// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	TransactionResponses *[]TableResponse

	// Version contains the information returned from the x-ms-version header response.
	Version *string
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
func (t *TableClient) SubmitTransaction(transactionActions []TableTransactionAction, tableSubmitTransactionOptions *TableSubmitTransactionOptions, ctx context.Context) error {
	return t.submitTransactionInternal(&transactionActions, uuid.New(), uuid.New(), tableSubmitTransactionOptions, ctx)
}

// submitTransactionInternal is the internal implementation for SubmitTransaction. It allows for explicit configuration of the batch and changeset UUID values for testing.
func (t *TableClient) submitTransactionInternal(transactionActions *[]TableTransactionAction, batchUuid uuid.UUID, changesetUuid uuid.UUID, tableSubmitTransactionOptions *TableSubmitTransactionOptions, ctx context.Context) error {

	changesetBoundary := fmt.Sprintf("changeset_%s", changesetUuid.String())
	changeSetBody, err := t.generateChangesetBody(changesetBoundary, transactionActions)
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(t.client.con.Endpoint(), "$batch"))
	if err != nil {
		return err
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
		return err
	}
	batchWriter.Write(changeSetBody.Bytes())
	writer.Close()

	req.SetBody(azcore.NopCloser(bytes.NewReader(body.Bytes())), fmt.Sprintf("multipart/mixed; boundary=%s", boundary))

	resp, err := t.client.con.Pipeline().Do(req)
	if err != nil {
		return err
	}
	if !resp.HasStatusCode(http.StatusAccepted, http.StatusNoContent) {
		return t.client.transactionHandleError(resp)
	}
	return nil
}

// transactionHandleError handles the InsertEntity error response.
func (client *tableClient) transactionHandleError(resp *azcore.Response) error {
	var err TableServiceError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// generateChangesetBody generates the individual changesets for the various operations within the batch request.
// There is a changeset for Insert, Delete, Merge etc.
func (t *TableClient) generateChangesetBody(changesetBoundary string, transactionActions *[]TableTransactionAction) (*bytes.Buffer, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.SetBoundary(changesetBoundary)

	for _, be := range *transactionActions {
		t.generateEntitySubset(&be, writer)
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

	switch transactionAction.ActionType {
	case Delete:
		req, err = t.client.deleteEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), transactionAction.ETag, &TableDeleteEntityOptions{}, qo)
	case Add:
		toOdataAnnotatedDictionary(&entity)
		req, err = t.client.insertEntityCreateRequest(ctx, t.name, &TableInsertEntityOptions{TableEntityProperties: &entity, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, qo)
	case UpdateMerge:
	case UpsertReplace:
		toOdataAnnotatedDictionary(&entity)
		opts := &TableMergeEntityOptions{TableEntityProperties: &entity}
		if len(transactionAction.ETag) > 0 {
			opts.IfMatch = &transactionAction.ETag
		}
		req, err = t.client.mergeEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), opts, qo)
	case UpdateReplace:
	case UpsertMerge:
		toOdataAnnotatedDictionary(&entity)
		req, err = t.client.updateEntityCreateRequest(ctx, t.name, entity[PartitionKey].(string), entity[RowKey].(string), &TableUpdateEntityOptions{TableEntityProperties: &entity, IfMatch: &transactionAction.ETag}, qo)
	}

	urlAndVerb := fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL)
	operationWriter.Write([]byte(urlAndVerb))
	writeHeaders(req.Header, &operationWriter)
	operationWriter.Write([]byte("\r\n")) // additional \r\n is needed per changeset separating the "headers" and the body.
	io.Copy(operationWriter, req.Body)

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
