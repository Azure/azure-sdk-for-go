// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// This is the build of the package that binds to azure_data_cosmos_driver_native. It is selected
// automatically when cgo is enabled on a target for which this module carries a native archive.

// driverAvailable reports whether this build can reach the Cosmos driver.
const driverAvailable = true

// nativeDriver owns the driver handles a client needs. close releases them in reverse acquisition
// order.
//
// The runtime is per client rather than shared. It caches drivers by endpoint and only evicts that
// cache when it is freed, so a process-wide runtime would keep a closed client's driver alive and
// hand it to the next client for the same endpoint, defeating [Client.Close].
//
// The runtime, account reference and completion queue are built locally. Initialize or the first
// operation creates the driver, which fetches account properties, seeds routing state and creates
// the account transport.
//
// The queue binds to the runtime, not to the driver, so it exists before driver creation and makes
// that network work cancellable through cosmos_driver_get_or_create_submit.
type nativeDriver struct {
	// cfg is kept across the local setup in openDriver and the asynchronous driver creation that
	// follows it.
	cfg     driverConfig
	runtime *C.cosmos_runtime_t
	account *C.cosmos_account_ref_t

	tokenProvider *tokenProviderState

	// mu guards everything below it. It is deliberately not held across driver creation or
	// container resolution: both wait on the network, and holding it there would make a second
	// caller block on the mutex where it cannot honor its own context.
	//
	// That means mu does not keep these handles alive for an operation's duration, and close
	// does not wait for one to finish. Client.mu is what does: Close takes it for write, which
	// blocks until every operation holding it for read has returned, so close only ever runs
	// with nothing in flight. Calling into a nativeDriver outside Client.acquire breaks that,
	// which is why nothing but a test does.
	mu sync.Mutex
	// created records that creation reached a verdict, so a failure is reported to every later
	// caller rather than retried on every operation. A cancelled attempt is not a verdict; see
	// ensureDriver.
	created   bool
	closed    bool
	driver    *C.cosmos_driver_t
	driverErr error

	// creating is non-nil while an attempt is in flight and closed when it finishes. Waiters
	// block on it rather than on mu, so a second caller can still honor its own context while
	// the first caller's attempt runs.
	creating chan struct{}

	// reactor drains the completion queue operations are answered through. The queue binds to
	// the runtime rather than the driver, so it is created with the client and is what makes
	// driver creation itself awaitable.
	reactor *reactor

	// containers caches resolved container handles by "database/container". Resolving reads
	// container metadata from the gateway on a miss, so an item operation would otherwise pay
	// for that lookup on every call.
	containers map[string]*C.cosmos_container_ref_t
}

// initialize eagerly creates the driver, which fills its account-properties and routing caches.
func (d *nativeDriver) initialize(ctx context.Context) error {
	_, err := d.ensureDriver(ctx)
	return err
}

// cancel stops host token acquisition before Close waits for in-flight operations.
func (d *nativeDriver) cancel() {
	if d != nil && d.tokenProvider != nil {
		d.tokenProvider.cancel()
	}
}

// verifyDriverVersion reports whether the library that was linked is the one this package was
// built against.
//
// The header is vendored, so the binding compiles against one version and links against whatever
// archive it was given. A mismatch is not a compile error: it is a struct layout or a calling
// convention that has quietly moved, which surfaces as corrupted values or a crash somewhere far
// from the cause. Checking on first initialization turns that into a message that names the problem.
//
// The ABI carries no major-version concept yet, so this compares the whole version. Once it does,
// this should relax to a compatible range rather than an exact match, which is what the
// distribution design calls for.
func verifyDriverVersion() error {
	linked := C.GoString(C.cosmos_version())
	if linked == nativeDriverVersion {
		return nil
	}
	return &Error{
		Code: CodeClientError,
		Message: fmt.Sprintf(
			"azcosmos: the linked Cosmos driver is version %q, but this package was built against %q; "+
				"the vendored header and the driver library have to come from the same version",
			linked, nativeDriverVersion),
	}
}

// openDriver acquires the resources a client can acquire locally: the runtime, the account
// reference and the completion queue, none of which touch the network. Initialize and operations
// create the driver separately so that network work receives their context.
func openDriver(cfg driverConfig) (*nativeDriver, error) {
	d := &nativeDriver{cfg: cfg}
	if err := d.buildRuntime(); err != nil {
		return nil, err
	}
	if err := d.buildAccount(cfg); err != nil {
		_ = d.close()
		return nil, err
	}
	// The token-provider handle now owns the credential reference.
	d.cfg.tokenCredential = nil
	// Before the driver rather than after it, because creating the driver is itself answered
	// through this queue.
	var err error
	if d.reactor, err = newReactor(d.runtime); err != nil {
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
// A verdict is cached, so a driver that cannot be created is reported to every later caller rather
// than retried on every operation. A cancelled attempt is not a verdict: one caller's deadline
// says nothing about the account, so the next caller starts fresh instead of inheriting a
// cancellation it never asked for.
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
			driver, err := d.driver, d.driverErr
			d.mu.Unlock()
			return driver, err
		}
		if inFlight := d.creating; inFlight != nil {
			d.mu.Unlock()
			select {
			case <-inFlight:
				// The attempt finished; loop to read whatever verdict it reached, or to start
				// a fresh attempt if it was cancelled.
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		done := make(chan struct{})
		d.creating = done
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
		case isCancellation(err):
			// Deliberately not recorded; see the doc comment.
		default:
			d.created, d.driver, d.driverErr = true, driver, err
		}
		d.mu.Unlock()
		close(done)

		return driver, err
	}
}

// createDriver runs one driver-creation attempt.
func (d *nativeDriver) createDriver(ctx context.Context) (*C.cosmos_driver_t, error) {
	// Checked here rather than in openDriver so it runs once per client and on the path that
	// actually calls into the driver, rather than on every construction.
	if err := verifyDriverVersion(); err != nil {
		return nil, err
	}

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

// isCancellation reports whether an attempt was abandoned rather than answered.
//
// It covers the driver's own cancelled completion as well as the caller's context, because the
// distinction that matters is whether anything was learned about the account, and a cancellation
// of either kind means nothing was.
func isCancellation(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var cosmosErr *Error
	return errors.As(err, &cosmosErr) && cosmosErr.Code == CodeOperationCancelled
}

// errClientClosed is the failure every entry point reports once the client is closed.
func errClientClosed() error {
	return &Error{Code: CodeClientClosed, Message: "the client has been closed"}
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
	status := C.cosmos_driver_options_build(d.account, config, &options) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if err := statusError(status, nil, "building the driver options"); err != nil {
		return nil, err
	}
	return options, nil
}

// buildRuntime creates the Tokio runtime the driver executes on.
//
// ApplicationID is applied here rather than with the other client options because the C ABI carries
// the user agent on the runtime, not on the driver.
func (d *nativeDriver) buildRuntime() error {
	// Seeded from the defaults so that fields this binding does not set keep the driver's values
	// rather than a Go zero.
	options := C.cosmos_runtime_options_default()

	if d.cfg.options.ApplicationID != "" {
		// Copied into the runtime before the call returns, so freeing it here is safe.
		suffix := C.CString(d.cfg.options.ApplicationID)
		defer C.free(unsafe.Pointer(suffix))
		options.user_agent_suffix = suffix
	}

	var richErr *C.cosmos_error_t
	status := C.cosmos_runtime_build(&options, &d.runtime, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	err := statusError(status, richErr, "building the driver runtime")
	if err == nil || d.cfg.options.ApplicationID == "" {
		return err
	}

	// ApplicationID is the only runtime option this binding changes from the driver's defaults.
	// The C ABI reports an invalid value as a bare status with no field name, so identify the field
	// without duplicating the driver's validation rule.
	var cosmosErr *Error
	if errors.As(err, &cosmosErr) &&
		cosmosErr.SubStatus == int(C.COSMOS_SUB_STATUS_CLIENT_FFI_INVALID_OPTION_VALUE) {
		cosmosErr.Message = "azcosmos: the Cosmos driver rejected ClientOptions.ApplicationID"
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

	endpoint := C.CString(cfg.endpoint)
	defer C.free(unsafe.Pointer(endpoint))
	// The key is copied into a Rust Secret before the call returns, so freeing it here is safe.
	key := C.CString(cfg.accountKey)
	defer C.free(unsafe.Pointer(key))

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
	// Freeing here is only safe because the caller guarantees no operation is in flight; see
	// nativeDriver.mu. The lock below orders this against a concurrent state read, not against
	// an operation, which it can no longer wait for.
	d.mu.Lock()
	defer d.mu.Unlock()
	d.closed = true

	// The reactor goes first: it holds the completion queue, and stopping it is what guarantees
	// no thread is still reading completions that reference the driver.
	d.reactor.close()
	d.reactor = nil

	for key, container := range d.containers {
		C.cosmos_container_ref_free(container)
		delete(d.containers, key)
	}

	C.cosmos_driver_free(d.driver)
	d.driver = nil
	C.cosmos_account_ref_free(d.account)
	d.account = nil
	C.cosmos_runtime_free(d.runtime)
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
