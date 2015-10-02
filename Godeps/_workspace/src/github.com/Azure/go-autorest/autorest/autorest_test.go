package autorest

import (
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/mocks"
)

func TestResponseHasStatusCode(t *testing.T) {
	codes := []int{http.StatusOK, http.StatusAccepted}
	resp := &http.Response{StatusCode: http.StatusAccepted}
	if !ResponseHasStatusCode(resp, codes...) {
		t.Errorf("autorest: ResponseHasStatusCode failed to find %v in %v", resp.StatusCode, codes)
	}
}

func TestResponseHasStatusCodeNotPresent(t *testing.T) {
	codes := []int{http.StatusOK, http.StatusAccepted}
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	if ResponseHasStatusCode(resp, codes...) {
		t.Errorf("autorest: ResponseHasStatusCode unexpectedly found %v in %v", resp.StatusCode, codes)
	}
}

func TestResponseRequiresPollingIgnoresSuccess(t *testing.T) {
	resp := mocks.NewResponse()

	if ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling did not ignore a successful response")
	}
}

func TestResponseRequiresPollingLeavesBodyOpen(t *testing.T) {
	resp := mocks.NewResponse()

	ResponseRequiresPolling(resp)
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: ResponseRequiresPolling closed the responise body while ignoring a successful response")
	}
}

func TestResponseRequiresPollingDefaultsToAcceptedStatusCode(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	if !ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling failed to create a request for default 202 Accepted status code")
	}
}

func TestResponseRequiresPollingReturnsFalseForUnexpectedStatusCodes(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 InternalServerError", http.StatusInternalServerError)
	mocks.SetAcceptedHeaders(resp)

	if ResponseRequiresPolling(resp) {
		t.Error("autorest: ResponseRequiresPolling did not return false when ignoring a status code")
	}
}

func TestNewPollingRequestLeavesBodyOpenWhenLocationHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 InternalServerError", http.StatusInternalServerError)

	NewPollingRequest(resp, NullAuthorizer{})
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: NewPollingRequest closed the http.Request Body when the Location header was missing")
	}
}

func TestNewPollingRequestDoesNotReturnARequestWhenLocationHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("500 InternalServerError", http.StatusInternalServerError)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: NewPollingRequest returned an http.Request when the Location header was missing")
	}
}

func TestNewPollingRequestReturnsAnErrorWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	_, err := NewPollingRequest(resp, mockFailingAuthorizer{})
	if err == nil {
		t.Error("autorest: NewPollingRequest failed to return an error when Prepare fails")
	}
}

func TestNewPollingRequestLeavesBodyOpenWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerLocation), testBadURL)

	_, err := NewPollingRequest(resp, NullAuthorizer{})
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: NewPollingRequest closed the http.Request Body when Prepare returned an error (%v)", err)
	}
}

func TestNewPollingRequestDoesNotReturnARequestWhenPrepareFails(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerLocation), testBadURL)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})
	if req != nil {
		t.Error("autorest: NewPollingRequest returned an http.Request when Prepare failed")
	}
}

func TestNewPollingRequestClosesTheResponseBody(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	NewPollingRequest(resp, NullAuthorizer{})
	if resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: NewPollingRequest failed to close the response body when creating a new request")
	}
}

func TestNewPollingRequestReturnsAGetRequest(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})
	if req.Method != "GET" {
		t.Errorf("autorest: NewPollingRequest did not create an HTTP GET request -- actual method %v", req.Method)
	}
}

func TestNewPollingRequestProvidesTheURL(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})
	if req.URL.String() != mocks.TestURL {
		t.Errorf("autorest: NewPollingRequest did not create an HTTP with the expected URL -- received %v, expected %v", req.URL, mocks.TestURL)
	}
}

func TestNewPollingRequestAppliesAuthorization(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, mockAuthorizer{})
	if req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)) != testAuthorizationHeader {
		t.Errorf("autorest: NewPollingRequest did not apply authorization -- received %v, expected %v",
			req.Header.Get(http.CanonicalHeaderKey(headerAuthorization)), testAuthorizationHeader)
	}
}

func TestGetPollingLocation(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	l := GetPollingLocation(resp)
	if len(l) == 0 {
		t.Errorf("autorest: GetPollingLocation failed to return Location header -- expected %v, received %v", mocks.TestURL, l)
	}
}

func TestGetPollingLocationReturnsEmptyStringForMissingLocation(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)

	l := GetPollingLocation(resp)
	if len(l) != 0 {
		t.Errorf("autorest: GetPollingLocation return a value without a Location header -- received %v", l)
	}
}

func TestGetPollingDelay(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	d := GetPollingDelay(resp, DefaultPollingDelay)
	if d != mocks.TestDelay {
		t.Errorf("autorest: GetPollingDelay failed to returned the expected delay -- expected %v, received %v", mocks.TestDelay, d)
	}
}

func TestGetPollingDelayReturnsDefaultDelayIfRetryHeaderIsMissing(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)

	d := GetPollingDelay(resp, DefaultPollingDelay)
	if d != DefaultPollingDelay {
		t.Errorf("autorest: GetPollingDelay failed to returned the default delay for a missing Retry-After header -- expected %v, received %v",
			DefaultPollingDelay, d)
	}
}

func TestGetPollingDelayReturnsDefaultDelayIfRetryHeaderIsMalformed(t *testing.T) {
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	resp.Header.Set(http.CanonicalHeaderKey(headerRetryAfter), "a very bad non-integer value")

	d := GetPollingDelay(resp, DefaultPollingDelay)
	if d != DefaultPollingDelay {
		t.Errorf("autorest: GetPollingDelay failed to returned the default delay for a malformed Retry-After header -- expected %v, received %v",
			DefaultPollingDelay, d)
	}
}

func TestPollForAttemptsStops(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	PollForAttempts(client, req, time.Duration(0), 5)
	if client.Attempts() < 5 || client.Attempts() > 5 {
		t.Errorf("autorest: PollForAttempts stopped incorrectly -- expected %v attempts, actual attempts were %v", 5, client.Attempts())
	}
}

func TestPollForDurationsStops(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	d := 10 * time.Millisecond
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) < d {
		t.Error("autorest: PollForDuration stopped too soon")
	}
}

func TestPollForDurationsStopsWithinReason(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	d := 10 * time.Millisecond
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	mocks.SetRetryHeader(resp, d)
	client.SetResponse(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	start := time.Now()
	PollForDuration(client, req, time.Duration(0), d)
	if time.Now().Sub(start) > (time.Duration(5.0) * d) {
		t.Error("autorest: PollForDuration took too long to stop -- exceeded 5 times expected duration")
	}
}

func TestPollingHonorsDelay(t *testing.T) {
	client := mocks.NewSender()
	client.EmitErrors(-1)

	d1 := 10 * time.Millisecond
	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	mocks.SetRetryHeader(resp, d1)
	client.SetResponse(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	start := time.Now()
	PollForAttempts(client, req, time.Duration(0), 2)
	d2 := time.Now().Sub(start)
	if d2 < d1 {
		t.Errorf("autorest: Polling failed to honor delay -- expected %v, actual %v", d1.Seconds(), d2.Seconds())
	}
}

func TestPollingReturnsErrorForExpectedStatusCode(t *testing.T) {
	client := mocks.NewSender()

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	client.SetResponse(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, http.StatusAccepted)
	if err == nil {
		t.Error("autorest: Polling failed to emit error for known status code")
	}
}

func TestPollingReturnsNoErrorForUnexpectedStatusCode(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("500 InternalServerError", http.StatusInternalServerError)

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1, http.StatusAccepted)
	if err != nil {
		t.Error("autorest: Polling emitted error for unknown status code")
	}
}

func TestPollingReturnsDefaultsToAcceptedStatusCode(t *testing.T) {
	client := mocks.NewSender()

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	client.SetResponse(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	resp, err := PollForAttempts(client, req, time.Duration(0), 1)
	if err == nil {
		t.Error("autorest: Polling failed to default to HTTP 202")
	}
}

func TestPollingLeavesFinalBodyOpen(t *testing.T) {
	client := mocks.NewSender()
	client.EmitStatus("500 InternalServerError", http.StatusInternalServerError)

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})

	resp, _ = PollForAttempts(client, req, time.Duration(0), 1)
	if !resp.Body.(*mocks.Body).IsOpen() {
		t.Error("autorest: Polling unexpectedly closed the response body")
	}
}

func TestDecorateForPollingCloseBodyOnEachAttempt(t *testing.T) {
	client := mocks.NewSender()

	resp := mocks.NewResponseWithStatus("202 Accepted", http.StatusAccepted)
	mocks.SetAcceptedHeaders(resp)
	client.SetResponse(resp)

	req, _ := NewPollingRequest(resp, NullAuthorizer{})
	resp, _ = PollForAttempts(client, req, time.Duration(0), 5)
	if resp.Body.(*mocks.Body).CloseAttempts() < 5 {
		t.Errorf("autorest: decorateForPolling failed to close the response Body between requests -- expected %v, received %v",
			5, resp.Body.(*mocks.Body).CloseAttempts())
	}
}
