// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/rntbd"
	"github.com/google/uuid"
)

type directModeTransport struct {
	endpointProvider *rntbd.EndpointProvider
	addressResolver  AddressResolver
	fallback         http.RoundTripper

	mu     sync.RWMutex
	closed bool
}

// AddressResolver resolves backend addresses for Direct mode requests.
// It determines the physical replica endpoints for a given request based on
// partition key range information from request headers.
type AddressResolver interface {
	Resolve(ctx context.Context, req *http.Request) ([]*url.URL, error)
}

// DirectModeTransportOptions configures the Direct mode HTTP transport.
// It allows customization of connection pooling, address resolution, and
// fallback behavior when Direct mode cannot be used.
type DirectModeTransportOptions struct {
	PoolOptions     *rntbd.PoolOptions
	AddressResolver AddressResolver
	Fallback        http.RoundTripper
}

func newDirectModeTransport(opts *DirectModeTransportOptions) *directModeTransport {
	if opts == nil {
		opts = &DirectModeTransportOptions{}
	}

	poolOpts := opts.PoolOptions
	if poolOpts == nil {
		poolOpts = rntbd.DefaultPoolOptions()
	}

	fallback := opts.Fallback
	if fallback == nil {
		fallback = http.DefaultTransport
	}

	return &directModeTransport{
		endpointProvider: rntbd.NewEndpointProvider(poolOpts),
		addressResolver:  opts.AddressResolver,
		fallback:         fallback,
	}
}

func (t *directModeTransport) Do(req *http.Request) (*http.Response, error) {
	if !t.shouldUseDirect(req) {
		return t.fallbackWithCleanHeaders(req)
	}

	svcReq, err := t.httpToServiceRequest(req)
	if err != nil {
		return nil, fmt.Errorf("direct mode: failed to convert request: %w", err)
	}

	addresses, err := t.resolveAddresses(req.Context(), req)
	if err != nil {
		return t.fallbackWithCleanHeaders(req)
	}

	if len(addresses) == 0 {
		return t.fallbackWithCleanHeaders(req)
	}

	var lastErr error
	for _, addr := range addresses {
		resp, err := t.sendToEndpoint(req.Context(), addr, svcReq)
		if err != nil {
			lastErr = err
			continue
		}
		return resp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("direct mode: all replicas failed: %w", lastErr)
	}
	return nil, errors.New("direct mode: no replicas available")
}

func (t *directModeTransport) fallbackWithCleanHeaders(req *http.Request) (*http.Response, error) {
	// Gateway rejects requests with both partitionKey AND partitionKeyRangeId/EPK
	req.Header.Del(cosmosHeaderPartitionKeyRangeId)
	req.Header.Del(cosmosHeaderCollectionRid)
	req.Header.Del(cosmosHeaderEffectivePartitionKey)
	return t.fallback.RoundTrip(req)
}

func (t *directModeTransport) shouldUseDirect(req *http.Request) bool {
	if t.addressResolver == nil {
		return false
	}

	resType := req.Header.Get("x-ms-cosmos-resource-type")
	if resType == "" {
		return false
	}

	switch strings.ToLower(resType) {
	case "document", "storedprocedure", "trigger", "userdefinedfunctions", "conflict", "attachment":
		return true
	default:
		return false
	}
}

func (t *directModeTransport) httpToServiceRequest(req *http.Request) (*rntbd.ServiceRequest, error) {
	activityIDStr := req.Header.Get("x-ms-activity-id")
	var activityID uuid.UUID
	if activityIDStr != "" {
		var err error
		activityID, err = uuid.Parse(activityIDStr)
		if err != nil {
			activityID = uuid.New()
		}
	} else {
		activityID = uuid.New()
	}

	opType := t.httpMethodToOperationType(req.Method, req.Header)
	resType := t.parseResourceType(req.Header.Get("x-ms-cosmos-resource-type"))

	resourceAddress := strings.TrimPrefix(req.URL.Path, "/")

	headers := make(map[string]string)
	for key, values := range req.Header {
		if len(values) > 0 {
			// Normalize header keys to lowercase to match RNTBD header mapping
			headers[strings.ToLower(key)] = values[0]
		}
	}

	var content []byte
	if req.Body != nil {
		var err error
		content, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(content))
	}

	// For Direct Mode document operations, we need the collection's ResourceID (RID)
	// even when using name-based addressing. This is required by the RNTBD protocol.
	collectionRid := headers["x-ms-cosmos-collection-rid"]

	// Java SDK uses PartitionKeyRangeId internally for address resolution but does NOT send it as an RNTBD token.
	// The emulator rejects requests that have BOTH PartitionKey AND PartitionKeyRangeId tokens.
	if headers["x-ms-documentdb-partitionkeyrangeid"] != "" {
		delete(headers, "x-ms-documentdb-partitionkeyrangeid")
		delete(headers, "x-ms-effective-partition-key")
	}

	svcReq := &rntbd.ServiceRequest{
		OperationType:   opType,
		ResourceType:    resType,
		ResourceID:      collectionRid, // Collection RID required for document operations
		ResourceAddress: resourceAddress,
		IsNameBased:     !t.isRidBased(resourceAddress),
		ActivityID:      activityID,
		Headers:         headers,
		Content:         content,
	}

	return svcReq, nil
}

func (t *directModeTransport) httpMethodToOperationType(method string, headers http.Header) rntbd.OperationType {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return rntbd.OperationRead
	case http.MethodPost:
		// Check if this is a query request (POST with x-ms-documentdb-query header)
		if strings.EqualFold(headers.Get("x-ms-documentdb-query"), "true") {
			return rntbd.OperationQuery
		}
		return rntbd.OperationCreate
	case http.MethodPut:
		return rntbd.OperationReplace
	case http.MethodDelete:
		return rntbd.OperationDelete
	case http.MethodPatch:
		return rntbd.OperationPatch
	default:
		return rntbd.OperationInvalid
	}
}

func (t *directModeTransport) parseResourceType(resTypeStr string) rntbd.ResourceType {
	switch strings.ToLower(resTypeStr) {
	case "document":
		return rntbd.ResourceDocument
	case "collection":
		return rntbd.ResourceCollection
	case "database":
		return rntbd.ResourceDatabase
	case "storedprocedure":
		return rntbd.ResourceStoredProcedure
	case "trigger":
		return rntbd.ResourceTrigger
	case "userdefinedfunctions":
		return rntbd.ResourceUserDefinedFunction
	case "conflict":
		return rntbd.ResourceConflict
	case "attachment":
		return rntbd.ResourceAttachment
	case "user":
		return rntbd.ResourceUser
	case "permission":
		return rntbd.ResourcePermission
	case "offer":
		return rntbd.ResourceOffer
	case "partitionkeyrange":
		return rntbd.ResourcePartitionKeyRange
	default:
		return rntbd.ResourceDocument
	}
}

func (t *directModeTransport) isRidBased(resourceAddress string) bool {
	if strings.HasPrefix(resourceAddress, "dbs/") {
		return false
	}
	return len(resourceAddress) > 0 && !strings.Contains(resourceAddress, "/")
}

func (t *directModeTransport) resolveAddresses(ctx context.Context, req *http.Request) ([]*url.URL, error) {
	if t.addressResolver == nil {
		return nil, errors.New("no address resolver configured")
	}
	return t.addressResolver.Resolve(ctx, req)
}

func (t *directModeTransport) sendToEndpoint(ctx context.Context, addr *url.URL, svcReq *rntbd.ServiceRequest) (*http.Response, error) {
	endpoint, err := t.endpointProvider.GetOrCreate(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	conn, err := endpoint.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Java SDK strips trailing slash from replica path (see RntbdRequestArgs line 59)
	svcReq.ReplicaPath = strings.TrimSuffix(addr.Path, "/")

	reqMsg, err := rntbd.BuildRequestMessage(svcReq)
	if err != nil {
		endpoint.RecordError()
		return nil, fmt.Errorf("failed to build RNTBD request: %w", err)
	}

	svcReq.MarkSendingRequestStarted(time.Now().UnixNano())

	respMsg, err := conn.Send(ctx, reqMsg)
	if err != nil {
		endpoint.RecordError()
		return nil, fmt.Errorf("failed to send RNTBD request: %w", err)
	}

	endpoint.RecordSuccess()

	return t.rntbdToHttpResponse(respMsg, svcReq)
}

func (t *directModeTransport) rntbdToHttpResponse(respMsg *rntbd.ResponseMessage, svcReq *rntbd.ServiceRequest) (*http.Response, error) {
	storeResp, err := rntbd.ParseResponseMessage(respMsg, svcReq.ResourceAddress)
	if err != nil {
		return nil, fmt.Errorf("direct mode: failed to parse response: %w", err)
	}

	statusCode := int(storeResp.StatusCode)

	header := make(http.Header)
	for key, value := range storeResp.Headers {
		header.Set(key, value)
	}

	var body io.ReadCloser
	if len(storeResp.Content) > 0 {
		body = io.NopCloser(bytes.NewReader(storeResp.Content))
	} else {
		body = io.NopCloser(bytes.NewReader(nil))
	}

	resp := &http.Response{
		Status:        fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode:    statusCode,
		Proto:         "RNTBD/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        header,
		Body:          body,
		ContentLength: int64(len(storeResp.Content)),
	}

	if lsnStr := header.Get(rntbd.RespHeaderLSN); lsnStr != "" {
		resp.Header.Set("x-ms-lsn", lsnStr)
	}
	if ruStr := header.Get(rntbd.RespHeaderRequestCharge); ruStr != "" {
		resp.Header.Set("x-ms-request-charge", ruStr)
	}
	if sessionToken := header.Get(rntbd.RespHeaderSessionToken); sessionToken != "" {
		resp.Header.Set("x-ms-session-token", sessionToken)
	}
	if etag := header.Get(rntbd.RespHeaderETag); etag != "" {
		resp.Header.Set("etag", etag)
	}
	if continuation := header.Get(rntbd.RespHeaderContinuation); continuation != "" {
		resp.Header.Set("x-ms-continuation", continuation)
	}

	return resp, nil
}

func (t *directModeTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}
	t.closed = true

	return t.endpointProvider.Close()
}

type gatewayAddressResolver struct {
	client      *http.Client
	accountHost string
	cred        *KeyCredential

	mu    sync.RWMutex
	cache map[string][]*url.URL
}

// GatewayAddressResolverOptions configures the gateway-based address resolver.
// This resolver queries the Cosmos DB gateway to discover backend replica addresses
// for a given partition key range.
type GatewayAddressResolverOptions struct {
	Client      *http.Client
	AccountHost string
	Credential  *KeyCredential
}

func newGatewayAddressResolver(opts *GatewayAddressResolverOptions) *gatewayAddressResolver {
	if opts == nil {
		opts = &GatewayAddressResolverOptions{}
	}

	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}

	return &gatewayAddressResolver{
		client:      client,
		accountHost: opts.AccountHost,
		cred:        opts.Credential,
		cache:       make(map[string][]*url.URL),
	}
}

func (r *gatewayAddressResolver) Resolve(ctx context.Context, req *http.Request) ([]*url.URL, error) {
	pkRangeID := req.Header.Get("x-ms-documentdb-partitionkeyrangeid")
	collectionRid := req.Header.Get("x-ms-cosmos-collection-rid")

	if collectionRid == "" || pkRangeID == "" {
		return nil, nil
	}

	// Extract the resource path from the request URL (e.g., "dbs/dbName/colls/collName/docs")
	resourcePath := strings.TrimPrefix(req.URL.Path, "/")

	cacheKey := collectionRid + ":" + pkRangeID

	r.mu.RLock()
	addresses, ok := r.cache[cacheKey]
	r.mu.RUnlock()
	if ok {
		return addresses, nil
	}

	addresses, err := r.fetchAddresses(ctx, collectionRid, pkRangeID, resourcePath)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.cache[cacheKey] = addresses
	r.mu.Unlock()

	return addresses, nil
}

func (r *gatewayAddressResolver) fetchAddresses(ctx context.Context, collectionRid, pkRangeID, resourcePath string) ([]*url.URL, error) {
	resolveFor := url.QueryEscape(resourcePath)
	filter := url.QueryEscape("protocol eq rntbd")
	path := fmt.Sprintf("/addresses/?$resolveFor=%s&$filter=%s&$partitionKeyRangeIds=%s", resolveFor, filter, pkRangeID)
	addrURL := r.accountHost + path

	addrReq, err := http.NewRequestWithContext(ctx, http.MethodGet, addrURL, nil)
	if err != nil {
		return nil, err
	}

	xmsDate := time.Now().UTC().Format(http.TimeFormat)
	addrReq.Header.Set(headerXmsDate, xmsDate)
	addrReq.Header.Set(headerXmsVersion, apiVersion)

	if r.cred != nil {
		resourceIdForAuth := r.extractCollectionPath(resourcePath, collectionRid)
		authHeader := r.cred.buildCanonicalizedAuthHeader(false, http.MethodGet, "docs", resourceIdForAuth, xmsDate, "master", "1.0")
		addrReq.Header.Set(headerAuthorization, authHeader)
	}

	resp, err := r.client.Do(addrReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("address resolution failed with status %d: %s", resp.StatusCode, string(body))
	}

	return r.parseAddressResponse(body)
}

func (r *gatewayAddressResolver) parseAddressResponse(body []byte) ([]*url.URL, error) {
	var response struct {
		Addresss []struct {
			PhysicalUri string `json:"physcialUri"`
			IsPrimary   bool   `json:"isPrimary"`
			Protocol    string `json:"protocol"`
		} `json:"Addresss"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse address response: %w", err)
	}

	var addresses []*url.URL
	for _, addr := range response.Addresss {
		if addr.Protocol == "rntbd" && addr.PhysicalUri != "" {
			parsed, err := url.Parse(addr.PhysicalUri)
			if err != nil {
				continue
			}
			addresses = append(addresses, parsed)
		}
	}

	return addresses, nil
}

func (r *gatewayAddressResolver) InvalidateCache(collectionRid, pkRangeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cache, collectionRid+":"+pkRangeID)
}

func (r *gatewayAddressResolver) ClearCache() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache = make(map[string][]*url.URL)
}

func (r *gatewayAddressResolver) extractCollectionPath(resourcePath string, collectionRid string) string {
	parts := strings.Split(resourcePath, "/")
	if len(parts) >= 4 && parts[0] == "dbs" && parts[2] == "colls" {
		return strings.Join(parts[:4], "/")
	}
	return collectionRid
}

var _ = strconv.Itoa
