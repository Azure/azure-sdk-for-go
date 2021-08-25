// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

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
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

type TransactionType string

const (
	Add           TransactionType = "add"
	UpdateMerge   TransactionType = "updatemerge"
	UpdateReplace TransactionType = "updatereplace"
	Delete        TransactionType = "delete"
	InsertMerge   TransactionType = "insertmerge"
	InsertReplace TransactionType = "insertreplace"
)

const (
	headerContentType             = "Content-Type"
	headerContentTransferEncoding = "Content-Transfer-Encoding"
	error_empty_transaction       = "transaction cannot be empty"
)

// Use azcore.ResponseError type, pass RawResponse, might have to create manually depending on constructor
// Int statuscode, Return the inner status code
// ErrorCode string set equal to OdataErrorMessage ("DuplicateRowKey")
// Lang/Value are useless at runtime, failedentity index as well


// type ODataErrorMessage struct {
// 	Lang  string `json:"lang"`
// 	Value string `json:"value"`
// }

// type ODataError struct {
// 	Code    string            `json:"code"`
// 	Message ODataErrorMessage `json:"message"`
// }

// type TableTransactionError struct {
// 	ODataError        ODataError `json:"odata.error"`
// 	FailedEntityIndex int
// }

// func (e *TableTransactionError) Error() string {
// 	return fmt.Sprintf("Code: %s, Message: %s", e.ODataError.Code, e.ODataError.Message.Value)
// }

type TransactionAction struct {
	ActionType TransactionType
	Entity     []byte
	ETag       *azcore.ETag
}

type TransactionResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
	// The response for a single table.
	TransactionResponses *[]http.Response
	// ContentType contains the information returned from the Content-Type header response.
	ContentType string
}

type SubmitTransactionOptions struct {
	RequestID *string
}

// SubmitTransaction submits the table transactional batch according to the slice of TableTransactionActions provided.
func (t *Client) SubmitTransaction(ctx context.Context, transactionActions []TransactionAction, tableSubmitTransactionOptions *SubmitTransactionOptions) (TransactionResponse, error) {
	u1, err := uuid.New()
	if err != nil {
		return TransactionResponse{}, err
	}
	u2, err := uuid.New()
	if err != nil {
		return TransactionResponse{}, err
	}
	return t.submitTransactionInternal(ctx, &transactionActions, u1, u2, tableSubmitTransactionOptions)
}

// submitTransactionInternal is the internal implementation for SubmitTransaction. It allows for explicit configuration of the batch and changeset UUID values for testing.
func (t *Client) submitTransactionInternal(ctx context.Context, transactionActions *[]TransactionAction, batchUuid uuid.UUID, changesetUuid uuid.UUID, tableSubmitTransactionOptions *SubmitTransactionOptions) (TransactionResponse, error) {
	if len(*transactionActions) == 0 {
		return TransactionResponse{}, errors.New(error_empty_transaction)
	}
	changesetBoundary := fmt.Sprintf("changeset_%s", changesetUuid.String())
	changeSetBody, err := t.generateChangesetBody(changesetBoundary, transactionActions)
	if err != nil {
		return TransactionResponse{}, err
	}
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(t.client.Con.Endpoint(), "$batch"))
	if err != nil {
		return TransactionResponse{}, err
	}
	req.Header.Set("x-ms-version", "2019-02-02")
	if tableSubmitTransactionOptions != nil && tableSubmitTransactionOptions.RequestID != nil {
		req.Header.Set("x-ms-client-request-id", *tableSubmitTransactionOptions.RequestID)
	}
	req.Header.Set("DataServiceVersion", "3.0")
	req.Header.Set("Accept", string(generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata))

	boundary := fmt.Sprintf("batch_%s", batchUuid.String())
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err = writer.SetBoundary(boundary)
	if err != nil {
		return TransactionResponse{}, err
	}
	h := make(textproto.MIMEHeader)
	h.Set(headerContentType, fmt.Sprintf("multipart/mixed; boundary=%s", changesetBoundary))
	batchWriter, err := writer.CreatePart(h)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = batchWriter.Write(changeSetBody.Bytes())
	if err != nil {
		return TransactionResponse{}, err
	}
	writer.Close()

	err = req.SetBody(azcore.NopCloser(bytes.NewReader(body.Bytes())), fmt.Sprintf("multipart/mixed; boundary=%s", boundary))
	if err != nil {
		return TransactionResponse{}, err
	}

	resp, err := t.client.Con.Pipeline().Do(req)
	if err != nil {
		return TransactionResponse{}, err
	}

	transactionResponse, err := buildTransactionResponse(req, resp, len(*transactionActions))
	if err != nil {
		return *transactionResponse, err
	}

	if !azcore.HasStatusCode(resp, http.StatusAccepted, http.StatusNoContent) {
		return TransactionResponse{}, azcore.NewResponseError(err, resp)
	}
	return *transactionResponse, nil
}

func buildTransactionResponse(req *azcore.Request, resp *http.Response, itemCount int) (*TransactionResponse, error) {
	innerResponses := make([]http.Response, itemCount)
	result := TransactionResponse{RawResponse: resp, TransactionResponses: &innerResponses}

	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = val
	}

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &TransactionResponse{}, err
	}
	reader := bytes.NewReader(bytesBody)
	if bytes.IndexByte(bytesBody, '{') == 0 {
		// This is a failure and the body is json
		return &TransactionResponse{}, newTableTransactionError(bytesBody)
	}

	outerBoundary := getBoundaryName(bytesBody)
	mpReader := multipart.NewReader(reader, outerBoundary)
	outerPart, err := mpReader.NextPart()
	if err != nil {
		return &TransactionResponse{}, err
	}

	innerBytes, err := ioutil.ReadAll(outerPart)
	if err != nil {
		return &TransactionResponse{}, err
	}
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
			return &TransactionResponse{}, err
		}
		if r.StatusCode >= 400 {
			errorBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return &TransactionResponse{}, err
			} else {
				innerResponses = []http.Response{*r}
				return &result, newTableTransactionError(errorBody)
			}
		}
		innerResponses[i] = *r
		i++
	}

	return &result, nil
}

func getBoundaryName(bytesBody []byte) string {
	end := bytes.Index(bytesBody, []byte("\n"))
	if end > 0 && bytesBody[end-1] == '\r' {
		end -= 1
	}
	return string(bytesBody[2:end])
}

// newTableTransactionError handles the SubmitTransaction error response.
func newTableTransactionError(errorBody []byte) error {
	oe := TableTransactionError{}
	if err := json.Unmarshal(errorBody, &oe); err == nil {
		if i := strings.Index(oe.ODataError.Message.Value, ":"); i > 0 {
			if val, err := strconv.Atoi(oe.ODataError.Message.Value[0:i]); err == nil {
				oe.FailedEntityIndex = val
			}
		}
		return &oe
	}
	return fmt.Errorf("unknown error: %s", string(errorBody))
}

// generateChangesetBody generates the individual changesets for the various operations within the batch request.
// There is a changeset for Insert, Delete, Merge etc.
func (t *Client) generateChangesetBody(changesetBoundary string, transactionActions *[]TransactionAction) (*bytes.Buffer, error) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.SetBoundary(changesetBoundary)
	if err != nil {
		return nil, err
	}

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
func (t *Client) generateEntitySubset(transactionAction *TransactionAction, writer *multipart.Writer) error {
	h := make(textproto.MIMEHeader)
	h.Set(headerContentTransferEncoding, "binary")
	h.Set(headerContentType, "application/http")
	qo := &generated.QueryOptions{Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr()}

	operationWriter, err := writer.CreatePart(h)
	if err != nil {
		return err
	}
	var req *azcore.Request
	var entity map[string]interface{}
	err = json.Unmarshal(transactionAction.Entity, &entity)
	if err != nil {
		return err
	}

	if _, ok := entity[partitionKey]; !ok {
		return fmt.Errorf("entity properties must contain a %s property", partitionKey)
	}
	if _, ok := entity[rowKey]; !ok {
		return fmt.Errorf("entity properties must contain a %s property", rowKey)
	}
	// Consider empty ETags as '*'
	if transactionAction.ETag == nil {
		star := azcore.ETagAny
		transactionAction.ETag = &star
	}

	switch transactionAction.ActionType {
	case Delete:
		req, err = t.client.DeleteEntityCreateRequest(ctx, t.name, entity[partitionKey].(string), entity[rowKey].(string), string(*transactionAction.ETag), &generated.TableDeleteEntityOptions{}, qo)
		if err != nil {
			return err
		}
	case Add:
		req, err = t.client.InsertEntityCreateRequest(ctx, t.name, &generated.TableInsertEntityOptions{TableEntityProperties: entity, ResponsePreference: generated.ResponseFormatReturnNoContent.ToPtr()}, qo)
		if err != nil {
			return err
		}
	case UpdateMerge:
		fallthrough
	case InsertMerge:
		opts := &generated.TableMergeEntityOptions{TableEntityProperties: entity}
		if transactionAction.ETag != nil {
			opts.IfMatch = to.StringPtr(string(*transactionAction.ETag))
		}
		req, err = t.client.MergeEntityCreateRequest(ctx, t.name, entity[partitionKey].(string), entity[rowKey].(string), opts, qo)
		if err != nil {
			return err
		}
		if isCosmosEndpoint(t.client.Con.Endpoint()) {
			transformPatchToCosmosPost(req)
		}
	case UpdateReplace:
		fallthrough
	case InsertReplace:
		req, err = t.client.UpdateEntityCreateRequest(ctx, t.name, entity[partitionKey].(string), entity[rowKey].(string), &generated.TableUpdateEntityOptions{TableEntityProperties: entity, IfMatch: to.StringPtr(string(*transactionAction.ETag))}, qo)
		if err != nil {
			return err
		}
	}

	urlAndVerb := fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL)
	_, err = operationWriter.Write([]byte(urlAndVerb))
	if err != nil {
		return err
	}
	err = writeHeaders(req.Header, &operationWriter)
	if err != nil {
		return err
	}
	_, err = operationWriter.Write([]byte("\r\n")) // additional \r\n is needed per changeset separating the "headers" and the body.
	if err != nil {
		return err
	}
	if req.Body != nil {
		_, err = io.Copy(operationWriter, req.Body)

	}

	return err
}

func writeHeaders(h http.Header, writer *io.Writer) error {
	// This way it is guaranteed the headers will be written in a sorted order
	var keys []string
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var err error
	for _, k := range keys {
		_, err = (*writer).Write([]byte(fmt.Sprintf("%s: %s\r\n", k, h.Get(k))))

	}
	return err
}
