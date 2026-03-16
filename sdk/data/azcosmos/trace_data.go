// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	traceDatumKeyClientSideRequestStats        = "Client Side Request Stats"
	traceDatumKeyTransportRequest              = "Microsoft.Azure.Documents.ServerStoreModel Transport Request"
	traceDatumKeyPointOperationStatistics      = "PointOperationStatisticsTraceDatum"
	traceDatumKeyPointOperationStatisticsError = "Point Operation Statistics"
	traceDatumKeyQueryMetrics                  = "Query Metrics"
)

type requestDiagnosticsState struct {
	requestTrace        *trace
	clientSideStats     *clientSideRequestStatisticsTraceDatum
	resourceType        resourceType
	requestSessionToken string
}

type requestDiagnosticsStateKey struct{}

func requestDiagnosticsStateFromContext(ctx context.Context) *requestDiagnosticsState {
	state, _ := ctx.Value(requestDiagnosticsStateKey{}).(*requestDiagnosticsState)
	return state
}

func withRequestDiagnosticsState(ctx context.Context, state *requestDiagnosticsState) context.Context {
	return context.WithValue(ctx, requestDiagnosticsStateKey{}, state)
}

func attachRequestDiagnostics(req *policy.Request, operationContext pipelineRequestOptions) *policy.Request {
	requestContext := req.Raw().Context()
	parentTrace := traceFromContext(requestContext)

	var requestTrace *trace
	if parentTrace != nil {
		requestTrace = parentTrace.StartChild(traceDatumKeyTransportRequest)
	} else {
		requestTrace = newRootTrace(traceDatumKeyTransportRequest)
	}

	clientSideStats := newClientSideRequestStatisticsTraceDatum(time.Now().UTC(), requestTrace)
	requestTrace.AddDatum(traceDatumKeyClientSideRequestStats, clientSideStats)

	state := &requestDiagnosticsState{
		requestTrace:        requestTrace,
		clientSideStats:     clientSideStats,
		resourceType:        operationContext.resourceType,
		requestSessionToken: req.Raw().Header.Get(cosmosHeaderSessionToken),
	}

	requestContext = withTrace(requestContext, requestTrace)
	requestContext = withRequestDiagnosticsState(requestContext, state)

	return req.WithContext(requestContext)
}

func addPointOperationStatisticsFromResponse(resp *http.Response, errorMessage string, key string) {
	if resp == nil || resp.Request == nil {
		return
	}

	state := requestDiagnosticsStateFromContext(resp.Request.Context())
	if state == nil || state.requestTrace == nil {
		return
	}

	requestURI := ""
	if resp.Request.URL != nil {
		requestURI = resp.Request.URL.String()
	}

	pointStats := pointOperationStatisticsTraceDatum{
		ActivityID:           resp.Header.Get(cosmosHeaderActivityId),
		ResponseTimeUTC:      time.Now().UTC(),
		StatusCode:           resp.StatusCode,
		SubStatusCode:        parseSubStatusCode(resp.Header.Get(cosmosHeaderSubstatus)),
		RequestCharge:        newResponse(resp).RequestCharge,
		RequestURI:           requestURI,
		ErrorMessage:         errorMessage,
		RequestSessionToken:  state.requestSessionToken,
		ResponseSessionToken: resp.Header.Get(cosmosHeaderSessionToken),
		BELatencyInMs:        resp.Header.Get(headerXmsRequestDurationMs),
	}

	state.requestTrace.AddOrUpdateDatum(key, pointStats)
}

func recordQueryMetricsFromResponse(resp *http.Response) {
	if resp == nil || resp.Request == nil {
		return
	}

	queryMetrics := resp.Header.Get(cosmosHeaderQueryMetrics)
	if queryMetrics == "" {
		return
	}

	state := requestDiagnosticsStateFromContext(resp.Request.Context())
	if state == nil || state.requestTrace == nil {
		return
	}

	state.requestTrace.AddOrUpdateDatum(traceDatumKeyQueryMetrics, queryMetrics)
}

type clientSideRequestStatisticsTraceDatum struct {
	mu                          sync.Mutex
	trace                       *trace
	requestStartTimeUTC         time.Time
	requestEndTimeUTC           *time.Time
	regionsContacted            []contactedRegion
	regionsContactedByURI       map[string]struct{}
	httpResponseStatistics      []httpResponseStatistics
	addressResolutionStatistics []addressResolutionStatistics
	forceAddressRefreshes       []forceAddressRefresh
}

type clientSideRequestStatisticsSnapshot struct {
	regionsContacted            []contactedRegion
	httpResponseStatistics      []httpResponseStatistics
	addressResolutionStatistics []addressResolutionStatistics
	forceAddressRefreshes       []forceAddressRefresh
}

type contactedRegion struct {
	name string
	uri  string
}

type httpResponseStatistics struct {
	startTimeUTC     time.Time
	duration         time.Duration
	requestURI       string
	resourceType     resourceType
	httpMethod       string
	activityID       string
	exceptionType    string
	exceptionMessage string
	statusCode       int
	statusCodeText   string
	reasonPhrase     string
	subStatusCode    int
}

type addressResolutionStatistics struct {
	startTimeUTC   time.Time
	endTimeUTC     *time.Time
	targetEndpoint string
}

type forceAddressRefresh struct {
	noChangeToCache []string
	original        []string
	newValues       []string
}

type pointOperationStatisticsTraceDatum struct {
	ActivityID           string
	ResponseTimeUTC      time.Time
	StatusCode           int
	SubStatusCode        int
	RequestCharge        float32
	RequestURI           string
	ErrorMessage         string
	RequestSessionToken  string
	ResponseSessionToken string
	BELatencyInMs        string
}

func newClientSideRequestStatisticsTraceDatum(startTime time.Time, trace *trace) *clientSideRequestStatisticsTraceDatum {
	return &clientSideRequestStatisticsTraceDatum{
		trace:                       trace,
		requestStartTimeUTC:         startTime,
		regionsContacted:            []contactedRegion{},
		regionsContactedByURI:       map[string]struct{}{},
		httpResponseStatistics:      []httpResponseStatistics{},
		addressResolutionStatistics: []addressResolutionStatistics{},
		forceAddressRefreshes:       []forceAddressRefresh{},
	}
}

func (d *clientSideRequestStatisticsTraceDatum) recordHTTPResponse(startTime time.Time, response *http.Response, resourceType resourceType, regionName string) {
	if response == nil {
		return
	}

	requestURI := ""
	if response.Request != nil && response.Request.URL != nil {
		requestURI = response.Request.URL.String()
	}

	d.recordRegion(regionName, requestURI)

	endTime := time.Now().UTC()
	d.updateRequestEndTime(endTime)

	stat := httpResponseStatistics{
		startTimeUTC:   startTime.UTC(),
		duration:       endTime.Sub(startTime),
		requestURI:     requestURI,
		resourceType:   resourceType,
		httpMethod:     response.Request.Method,
		activityID:     response.Header.Get(cosmosHeaderActivityId),
		statusCode:     response.StatusCode,
		statusCodeText: formatHTTPStatusCodeString(response.StatusCode),
		reasonPhrase:   http.StatusText(response.StatusCode),
		subStatusCode:  parseSubStatusCode(response.Header.Get(cosmosHeaderSubstatus)),
	}

	d.mu.Lock()
	d.httpResponseStatistics = append(d.httpResponseStatistics, stat)
	d.mu.Unlock()

	if response.StatusCode >= http.StatusBadRequest && d.trace != nil && d.trace.summary != nil {
		d.trace.summary.incrementFailedCount()
	}
}

func (d *clientSideRequestStatisticsTraceDatum) recordHTTPError(startTime time.Time, req *http.Request, err error, resourceType resourceType, regionName string) {
	requestURI := ""
	method := http.MethodGet
	if req != nil {
		method = req.Method
		if req.URL != nil {
			requestURI = req.URL.String()
		}
	}

	d.recordRegion(regionName, requestURI)

	endTime := time.Now().UTC()
	d.updateRequestEndTime(endTime)

	stat := httpResponseStatistics{
		startTimeUTC:     startTime.UTC(),
		duration:         endTime.Sub(startTime),
		requestURI:       requestURI,
		resourceType:     resourceType,
		httpMethod:       method,
		exceptionType:    fmt.Sprintf("%T", err),
		exceptionMessage: err.Error(),
	}

	d.mu.Lock()
	d.httpResponseStatistics = append(d.httpResponseStatistics, stat)
	d.mu.Unlock()

	if d.trace != nil && d.trace.summary != nil {
		d.trace.summary.incrementFailedCount()
	}
}

func (d *clientSideRequestStatisticsTraceDatum) snapshot() clientSideRequestStatisticsSnapshot {
	d.mu.Lock()
	defer d.mu.Unlock()

	regionsContacted := append([]contactedRegion(nil), d.regionsContacted...)
	httpResponseStatistics := append([]httpResponseStatistics(nil), d.httpResponseStatistics...)
	addressResolutionStatistics := append([]addressResolutionStatistics(nil), d.addressResolutionStatistics...)
	forceAddressRefreshes := append([]forceAddressRefresh(nil), d.forceAddressRefreshes...)

	return clientSideRequestStatisticsSnapshot{
		regionsContacted:            regionsContacted,
		httpResponseStatistics:      httpResponseStatistics,
		addressResolutionStatistics: addressResolutionStatistics,
		forceAddressRefreshes:       forceAddressRefreshes,
	}
}

func (d *clientSideRequestStatisticsTraceDatum) updateRequestEndTime(endTime time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.requestEndTimeUTC == nil || endTime.After(*d.requestEndTimeUTC) {
		value := endTime.UTC()
		d.requestEndTimeUTC = &value
	}
}

func (d *clientSideRequestStatisticsTraceDatum) recordRegion(regionName string, requestURI string) {
	if requestURI == "" {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.regionsContactedByURI[requestURI]; ok {
		return
	}

	d.regionsContactedByURI[requestURI] = struct{}{}
	d.regionsContacted = append(d.regionsContacted, contactedRegion{
		name: regionName,
		uri:  requestURI,
	})
}

func parseSubStatusCode(value string) int {
	if value == "" {
		return 0
	}

	subStatusCode, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return subStatusCode
}

func formatHTTPStatusCodeString(statusCode int) string {
	if statusCode == 0 {
		return ""
	}

	if text := http.StatusText(statusCode); text != "" {
		replacer := strings.NewReplacer(" ", "", "-", "", "'", "", ".", "")
		return replacer.Replace(text)
	}

	return strconv.Itoa(statusCode)
}

func resourceTypeName(resourceType resourceType) string {
	switch resourceType {
	case resourceTypeDatabase:
		return "Database"
	case resourceTypeCollection:
		return "DocumentCollection"
	case resourceTypeDocument:
		return "Document"
	case resourceTypeOffer:
		return "Offer"
	case resourceTypeDatabaseAccount:
		return "DatabaseAccount"
	case resourceTypePartitionKeyRange:
		return "PartitionKeyRange"
	default:
		return fmt.Sprintf("%d", resourceType)
	}
}
