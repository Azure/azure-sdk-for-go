// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore gocritic

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"

static const char *cosmos_go_header_version(void) {
	return AZURECOSMOSDRIVER_H_VERSION;
}
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// This is the build of the package that binds to azure_data_cosmos_driver_native. It is selected
// automatically when cgo is enabled on a target with a supported driver distribution module.

// driverAvailable reports whether this build can reach the Cosmos driver.
const driverAvailable = true

// nativeDriver owns the driver handles a client needs. close releases them in reverse acquisition
// order.
//
// Runtime ownership belongs to Runtime. Each client retains independent driver/account handles,
// credentials, options, metadata caches, and a completion queue.
//
// The runtime, account reference and completion queue are built locally. Initialize or the first
// operation creates the driver, which fetches account properties, seeds routing state and creates
// the account transport.
//
// The queue binds to the runtime, not to the driver, so it exists before driver creation and makes
// callers can stop waiting without cancelling native work.
type nativeDriver struct {
	// cfg is kept across the local setup in openDriver and the asynchronous driver creation that
	// follows it.
	cfg          driverConfig
	runtime      *C.cosmos_runtime_t
	ownedRuntime *nativeRuntime
	account      *C.cosmos_account_ref_t

	// Client defaults are immutable and borrowed read-only while snapshots are captured.
	clientOptions        *C.cosmos_operation_options_t
	releaseClientOptions func()

	tokenProvider *tokenProviderState
	// pending retains resources for submitted work after its Go caller stops waiting.
	pending sync.WaitGroup

	// mu guards everything below it. It is deliberately not held across driver creation because
	// that waits on the network, and holding it there would make a second caller block on the mutex
	// where it cannot honor its own context.
	//
	// Client.mu protects active callers; pending protects native submissions whose callers
	// have stopped waiting. Close drains both before freeing handles.
	mu sync.Mutex
	// created records that initialization succeeded. Failures are not cached because transport,
	// service and token acquisition can recover while this long-lived client remains in use.
	created bool
	closed  bool
	driver  *C.cosmos_driver_t

	// creating is non-nil while an attempt is in flight. Existing waiters receive that attempt's
	// result, while callers arriving after it finishes may start a retry.
	creating *driverCreation

	// beforeCreationWait is a test hook for observing that a caller joined an in-flight attempt.
	beforeCreationWait func()

	// beforePendingWait observes teardown reaching abandoned-operation cleanup in tests.
	beforePendingWait func()

	// reactor drains the completion queue operations are answered through. The queue binds to
	// the runtime rather than the driver, so it is created with the client and is what makes
	// driver creation itself awaitable.
	reactor *reactor

	cursors       map[*nativeQueryCursor]struct{}
	cursorReactor *cursorReactor

	// faultRules is a native-only test seam. Production clients never set it.
	faultRules []nativeFaultRule
}

type nativeFaultRule struct {
	id               string
	kind             int32
	region           string
	containerID      string
	transportKind    int32
	errorType        int32
	hitLimit         *int64
	delay            *time.Duration
	customStatusCode *int32
	customSubStatus  *int32
	retryAfter       *time.Duration
	startDelay       *time.Duration
	expireAfter      *time.Duration
}

const (
	nativeFaultOperationReadItem            = int32(C.COSMOS_FAULT_INJECTION_OPERATION_TYPE_READ_ITEM)
	nativeFaultErrorTooManyRequests         = int32(C.COSMOS_FAULT_INJECTION_ERROR_TYPE_TOO_MANY_REQUESTS)
	nativeFaultErrorReadSessionNotAvailable = int32(C.COSMOS_FAULT_INJECTION_ERROR_TYPE_READ_SESSION_NOT_AVAILABLE)
	nativeFaultTransportGatewayV2           = int32(C.COSMOS_FAULT_INJECTION_TRANSPORT_KIND_GATEWAY_V2)
	nativeFaultThrottleSubStatus            = int32(3200)
)

type driverCreation struct {
	done   chan struct{}
	driver *C.cosmos_driver_t
	err    error
}

// initialize eagerly creates the driver, which fills its account-properties and routing caches.
func (d *nativeDriver) initialize(ctx context.Context) error {
	ctx, _, release, err := d.snapshot(ctx, OperationOptions{})
	if err != nil {
		return err
	}
	defer release()
	_, err = d.ensureDriver(ctx)
	return err
}

// verifyDriverVersion checks header and library versions before any struct-sensitive ABI call.
func verifyDriverVersion() error {
	header := C.GoString(C.cosmos_go_header_version())
	linked := C.GoString(C.cosmos_version())
	return validateDriverVersions(header, linked)
}

func validateDriverVersions(header, linked string) error {
	if header == nativeDriverVersion && linked == nativeDriverVersion {
		return nil
	}
	return &Error{
		Code: CodeClientError,
		Message: fmt.Sprintf(
			"azcosmos: Cosmos driver version mismatch: vendored header %q, linked library %q, Go pin %q; "+
				"all three versions must match",
			header, linked, nativeDriverVersion),
	}
}

// openDriver acquires the resources a client can acquire locally: the runtime, the account
// reference and the completion queue, none of which touch the network. Initialize and operations
// create the driver separately so that network work receives their context.
func openDriver(cfg driverConfig) (*nativeDriver, error) {
	if err := verifyDriverVersion(); err != nil {
		return nil, err
	}

	d := &nativeDriver{cfg: cfg}
	if cfg.runtime != nil {
		d.runtime = cfg.runtime.native.handle
	} else {
		if err := d.buildRuntime(); err != nil {
			return nil, err
		}
	}
	d.clientOptions, d.releaseClientOptions = cfg.options.Operation.toNative()
	if err := d.buildAccount(cfg); err != nil {
		_ = d.close()
		return nil, err
	}
	// The account has copied the key, and the token-provider handle owns the credential reference.
	d.cfg.accountKey = ""
	d.cfg.tokenCredential = nil
	// Before the driver rather than after it, because creating the driver is itself answered
	// through this queue.
	var err error
	if d.reactor, err = newReactor(d.runtime, cfg.options.DiagnosticsVerbosity); err != nil {
		_ = d.close()
		return nil, err
	}
	return d, nil
}

// ensureDriver returns the initialized driver, creating it if necessary. It is safe for concurrent
// use.
//
// Creation is awaited rather than blocked on, so the caller's context bounds it. Only one attempt
// runs at a time: a second caller waits on the first rather than starting its own, and waits on a
// channel rather than on the mutex so its own context still applies.
//
// Only successful initialization is cached. A failed transport, service or token request can
// recover, so the next caller starts a fresh attempt rather than inheriting a stale failure.
//
// It reports [CodeClientClosed] once the client is closed rather than a NULL handle, so that a
// caller which reaches it without going through [Client.acquire] cannot pass NULL into the C ABI.
func (d *nativeDriver) ensureDriver(ctx context.Context) (*C.cosmos_driver_t, error) {
	for {
		d.mu.Lock()
		if d.closed {
			d.mu.Unlock()
			return nil, errClientClosed()
		}
		if d.created {
			driver := d.driver
			d.mu.Unlock()
			return driver, nil
		}
		if inFlight := d.creating; inFlight != nil {
			d.mu.Unlock()
			if d.beforeCreationWait != nil {
				d.beforeCreationWait()
			}
			select {
			case <-inFlight.done:
				if err := ctx.Err(); err != nil {
					return nil, err
				}
				if isCancellation(inFlight.err) {
					continue
				}
				return inFlight.driver, cloneError(inFlight.err)
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		attempt := &driverCreation{done: make(chan struct{})}
		d.creating = attempt
		d.mu.Unlock()

		driver, err := d.createDriver(ctx)

		d.mu.Lock()
		d.creating = nil
		switch {
		case d.closed:
			// Closed while this was in flight, which only a caller outside Client.acquire can
			// arrange. Nothing else owns the handle: it is published to d.driver by the default
			// branch below and this branch is taken instead, so freeing it here is the only way
			// it gets freed at all.
			C.cosmos_driver_free(driver)
			driver, err = nil, errClientClosed()
		case err == nil:
			d.created, d.driver = true, driver
		}
		attempt.driver, attempt.err = driver, err
		close(attempt.done)
		d.mu.Unlock()

		return driver, cloneError(err)
	}
}

// createDriver runs one driver-creation attempt.
func (d *nativeDriver) createDriver(ctx context.Context) (*C.cosmos_driver_t, error) {
	options, err := d.buildDriverOptions()
	if err != nil {
		return nil, err
	}
	// The submit call clones the options, so freeing them once it returns is safe.
	defer C.cosmos_driver_options_free(options)

	result, err := d.awaitCompletion(ctx, "creating the driver",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			return C.cosmos_driver_get_or_create_submit(d.runtime, d.account, options, queue, cookie, preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
		})
	if err != nil {
		return nil, err
	}
	if result.err != nil {
		result.release()
		return nil, result.err
	}
	if result.driver == nil {
		// A successful completion that carries no driver would otherwise be returned as a NULL
		// handle and fail later, somewhere with no connection to this call.
		result.release()
		return nil, &Error{
			Code:    CodeClientError,
			Message: "azcosmos: the driver was created but the completion carried no handle",
		}
	}
	return result.driver, nil
}

// errClientClosed is the failure every entry point reports once the client is closed.
func errClientClosed() error {
	return &Error{Code: CodeClientClosed, Message: "the client has been closed"}
}

func isCancellation(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// buildDriverOptions builds the driver's options from the client's. The account is cloned into
// them, so the returned handle outlives the reference passed here.
func (d *nativeDriver) buildDriverOptions() (*C.cosmos_driver_options_t, error) {
	config, release, err := d.cfg.options.toNative()
	if err != nil {
		return nil, err
	}
	defer release()

	var options *C.cosmos_driver_options_t
	var status C.cosmos_status_code_t
	if len(d.faultRules) == 0 {
		status = C.cosmos_driver_options_build(d.account, config, &options) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	} else {
		status = d.buildFaultInjectionOptions(config, &options)
	}
	if err := statusError(status, nil, "building the driver options"); err != nil {
		return nil, err
	}
	return options, nil
}

func freeNativeDriverOptions(options *C.cosmos_driver_options_t) {
	C.cosmos_driver_options_free(options)
}

func nativeOptionalFaultMilliseconds(duration *time.Duration) C.int64_t {
	if duration == nil {
		return -1
	}
	return C.int64_t(duration.Milliseconds())
}

func nativeOptionalFaultInt64(value *int64) C.int64_t {
	if value == nil {
		return -1
	}
	return C.int64_t(*value)
}

func nativeFaultInjectionRules(rules []nativeFaultRule) (*C.cosmos_fault_injection_rule_t, func()) {
	if len(rules) == 0 {
		return nil, func() {}
	}

	array := (*C.cosmos_fault_injection_rule_t)(C.malloc(
		C.size_t(len(rules)) * C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_rule_t{}))))
	releases := []func(){func() { C.free(unsafe.Pointer(array)) }}
	stringView := func(value string) C.cosmos_string_view_t {
		if value == "" {
			return C.cosmos_string_view_t{}
		}
		view, allocation := toNativeString(value)
		releases = append(releases, func() { C.free(allocation) })
		return view
	}

	for i, rule := range rules {
		nativeRule := &unsafe.Slice(array, len(rules))[i]
		*nativeRule = C.cosmos_fault_injection_rule_default()
		nativeRule.id = stringView(rule.id)
		nativeRule.hit_limit = nativeOptionalFaultInt64(rule.hitLimit)
		nativeRule.start_delay_ms = nativeOptionalFaultMilliseconds(rule.startDelay)
		nativeRule.expire_after_ms = nativeOptionalFaultMilliseconds(rule.expireAfter)

		condition := (*C.cosmos_fault_injection_condition_t)(C.malloc(
			C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_condition_t{}))))
		*condition = C.cosmos_fault_injection_condition_default()
		releases = append(releases, func() { C.free(unsafe.Pointer(condition)) })
		condition.operation_type = C.int32_t(rule.kind)
		condition.region = stringView(rule.region)
		condition.container_id = stringView(rule.containerID)
		condition.transport_kind = C.int32_t(rule.transportKind)
		nativeRule.condition = condition

		result := (*C.cosmos_fault_injection_result_t)(C.malloc(
			C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_result_t{}))))
		*result = C.cosmos_fault_injection_result_default()
		releases = append(releases, func() { C.free(unsafe.Pointer(result)) })
		result.error_type = C.int32_t(rule.errorType)
		result.probability = 1
		result.delay_ms = nativeOptionalFaultMilliseconds(rule.delay)
		result.retry_after_ms = nativeOptionalFaultMilliseconds(rule.retryAfter)
		if rule.retryAfter != nil {
			result.custom_status_code = http.StatusTooManyRequests
			result.custom_sub_status = C.int32_t(nativeFaultThrottleSubStatus)
		}
		if rule.customStatusCode != nil {
			result.custom_status_code = C.int32_t(*rule.customStatusCode)
		}
		if rule.customSubStatus != nil {
			result.custom_sub_status = C.int32_t(*rule.customSubStatus)
		}
		nativeRule.result = result
	}

	return array, func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
}

type nativeFaultRuleValues struct {
	id               string
	kind             int32
	region           string
	containerID      string
	transportKind    int32
	errorType        int32
	hitLimit         int64
	delayMS          int64
	probability      float32
	customStatusCode int32
	customSubStatus  int32
	retryAfterMS     int64
	startDelayMS     int64
	expireAfterMS    int64
}

// inspectNativeFaultRule converts one rule and reads it back in Go types for ABI mapping tests.
func inspectNativeFaultRule(rule nativeFaultRule) nativeFaultRuleValues {
	rules, release := nativeFaultInjectionRules([]nativeFaultRule{rule})
	defer release()

	nativeRule := unsafe.Slice(rules, 1)[0]
	return nativeFaultRuleValues{
		id:               fromNativeString(nativeRule.id),
		kind:             int32(nativeRule.condition.operation_type),
		region:           fromNativeString(nativeRule.condition.region),
		containerID:      fromNativeString(nativeRule.condition.container_id),
		transportKind:    int32(nativeRule.condition.transport_kind),
		errorType:        int32(nativeRule.result.error_type),
		hitLimit:         int64(nativeRule.hit_limit),
		delayMS:          int64(nativeRule.result.delay_ms),
		probability:      float32(nativeRule.result.probability),
		customStatusCode: int32(nativeRule.result.custom_status_code),
		customSubStatus:  int32(nativeRule.result.custom_sub_status),
		retryAfterMS:     int64(nativeRule.result.retry_after_ms),
		startDelayMS:     int64(nativeRule.start_delay_ms),
		expireAfterMS:    int64(nativeRule.expire_after_ms),
	}
}

// buildFaultInjectionOptions uses the versioned ABI only for native integration tests.
// The driver copies all rule memory before this call returns.
func (d *nativeDriver) buildFaultInjectionOptions(
	config *C.cosmos_driver_options_config_t,
	out **C.cosmos_driver_options_t,
) C.cosmos_status_code_t {
	rules, releaseRules := nativeFaultInjectionRules(d.faultRules)
	defer releaseRules()
	v2 := C.cosmos_driver_options_config_v2_default()
	v2.preferred_regions = config.preferred_regions
	v2.preferred_regions_len = config.preferred_regions_len
	v2.operation_options = config.operation_options
	v2.fault_injection_rules = rules
	v2.fault_injection_rules_len = C.uintptr_t(len(d.faultRules))
	v2.fault_injection_rule_stride = C.uintptr_t(unsafe.Sizeof(C.cosmos_fault_injection_rule_t{}))
	return C.cosmos_driver_options_build_v2(d.account, &v2, out) //nolint:gocritic // dupSubExpr targets cgo-generated code.
}

// buildRuntime creates the Tokio runtime the driver executes on.
func (d *nativeDriver) buildRuntime() error {
	native, err := openRuntime(RuntimeOptions{})
	if err == nil {
		d.ownedRuntime = native
		d.runtime = native.handle
		return nil
	}
	return err
}

// nativeInvalidOptionSubStatus exposes the C ABI value to tests, which cannot import C directly.
func nativeInvalidOptionSubStatus() int {
	return int(C.COSMOS_SUB_STATUS_CLIENT_FFI_INVALID_OPTION_VALUE)
}

// buildAccount creates the account reference that names the endpoint and carries the key.
func (d *nativeDriver) buildAccount(cfg driverConfig) error {
	if cfg.usesTokenCredential() {
		return d.buildTokenAccount(cfg)
	}

	endpoint, endpointAllocation := toNativeString(cfg.endpoint)
	defer C.free(endpointAllocation)
	// The key is copied into a Rust Secret before the call returns, so freeing it here is safe.
	key, keyAllocation := toNativeString(cfg.accountKey)
	defer C.free(keyAllocation)

	var richErr *C.cosmos_error_t
	status := C.cosmos_account_ref_with_master_key(endpoint, key, &d.account, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	return statusError(status, richErr, "building the account reference")
}

// close releases the handles in reverse acquisition order. Every handle is released even if an
// earlier one fails, so there is nothing left to retry.
//
// Freeing is a no-op on NULL, which is what makes it safe to call from openDriver's partial
// failure paths and when the driver was never created.
func (d *nativeDriver) close() error {
	if d == nil {
		return nil
	}
	// Client.mu prevents new submissions while pending drains; mu orders remaining state reads.
	d.mu.Lock()
	d.closed = true
	cursors := make([]*nativeQueryCursor, 0, len(d.cursors))
	for cursor := range d.cursors {
		cursors = append(cursors, cursor)
	}
	d.mu.Unlock()
	for _, cursor := range cursors {
		cursor.close()
	}
	if d.beforePendingWait != nil {
		d.beforePendingWait()
	}
	d.pending.Wait()
	d.cursorReactor.close()
	d.cursorReactor = nil
	d.mu.Lock()
	defer d.mu.Unlock()

	// The reactor goes first: it holds the completion queue, and stopping it is what guarantees
	// no thread is still reading completions that reference the driver.
	d.reactor.close()
	d.reactor = nil

	C.cosmos_driver_free(d.driver)
	d.driver = nil
	C.cosmos_account_ref_free(d.account)
	d.account = nil
	if d.releaseClientOptions != nil {
		d.releaseClientOptions()
		d.releaseClientOptions = nil
		d.clientOptions = nil
	}
	d.ownedRuntime.close()
	d.ownedRuntime = nil
	d.runtime = nil
	return nil
}

// statusError converts a packed cosmos_status_code_t and its optional rich error into an [Error],
// or nil on success. It always frees the rich error.
func statusError(status C.cosmos_status_code_t, richErr *C.cosmos_error_t, doing string) error {
	if richErr != nil {
		defer C.cosmos_error_free(richErr)
	}
	if status == C.COSMOS_STATUS_SUCCESS {
		return nil
	}

	httpStatus, subStatus := unpackStatus(status)
	err := &Error{
		Code:       codeForRichError(false, httpStatus, subStatus),
		StatusCode: httpStatus,
		SubStatus:  subStatus,
		Message:    fmt.Sprintf("azcosmos: %s", doing),
	}
	if richErr != nil {
		// The rich error's own fields are authoritative: the packed return value can be a
		// pre-flight rejection synthesized by the wrapper rather than the driver's own status.
		err.SubStatus = normalizeSubStatus(int32(richErr.sub_status))
		err.StatusCode = int(richErr.http_status_code)
		err.FromWire = richErr.is_from_wire == 1
		err.Code = codeForRichError(err.FromWire, err.StatusCode, err.SubStatus)
		if richErr.message != nil {
			err.Message = fmt.Sprintf("azcosmos: %s: %s", doing, C.GoString(richErr.message))
		}
		if richErr.activity_id != nil {
			err.ActivityID = C.GoString(richErr.activity_id)
		}
		if richErr.session_token != nil {
			err.SessionToken = SessionToken(C.GoString(richErr.session_token))
		}
		if richErr.etag != nil {
			err.ETag = azcore.ETag(C.GoString(richErr.etag))
		}
		if richErr.retry_after_ms >= 0 {
			err.RetryAfter = time.Duration(richErr.retry_after_ms) * time.Millisecond
		}
	}
	return err
}

// unpackStatus splits a packed cosmos_status_code_t into its HTTP status and sub-status, mirroring
// the COSMOS_STATUS_HTTP and COSMOS_STATUS_SUB macros. They are macros rather than functions, so
// cgo cannot call them and the arithmetic is reproduced here.
func unpackStatus(status C.cosmos_status_code_t) (httpStatus, subStatus int) {
	packed := uint32(status)
	return int(packed >> 16), int(packed & 0xFFFF)
}
