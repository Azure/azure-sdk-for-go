// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

/*
#cgo CFLAGS: -I${SRCDIR}/internal/native
#cgo LDFLAGS: -lazurecosmosdriver

#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// This is the build of the package that binds to azure_data_cosmos_driver_native. It is selected
// only with `-tags azcosmos_driver` and CGO_ENABLED=1, because it needs a C toolchain at build
// time and libazurecosmosdriver on the linker path. The default build uses driver_stub.go.

// driverAvailable reports whether this build can reach the Cosmos driver.
const driverAvailable = true

// nativeDriver owns the driver handles a client needs. close releases them in reverse acquisition
// order.
//
// The runtime is per client rather than shared. It caches drivers by endpoint and only evicts that
// cache when it is freed, so a process-wide runtime would keep a closed client's driver alive and
// hand it to the next client for the same endpoint, defeating [Client.Close].
//
// The runtime and the account reference are built when the client is constructed, because both are
// local and validate their inputs. The driver itself is not: creating it fetches the account's
// properties, so it does network I/O. Doing that in the constructor would make an unreachable
// account a construction failure rather than an operation failure, so it is created on first use.
//
// First use does not make it cancellable. cosmos_driver_get_or_create_blocking blocks a thread
// until the driver's own transport timeout, so an operation with a shorter deadline still waits
// for it. Cancellation needs cosmos_driver_get_or_create_submit, which delivers its result through
// a completion queue and so lands with the queue.
type nativeDriver struct {
	runtime *C.cosmos_runtime_t
	account *C.cosmos_account_ref_t

	// mu guards everything below it. It is held across driver creation so that a second caller
	// waits for the first rather than starting its own, and so that close cannot free a driver
	// that is still being created.
	mu sync.Mutex
	// created records that creation was attempted, so a failure is reported to every later
	// caller rather than retried on every operation.
	created   bool
	closed    bool
	driver    *C.cosmos_driver_t
	driverErr error

	// reactor drains the completion queue operations are answered through. It is created with
	// the driver, because the queue is bound to the runtime and only operations need it.
	reactor *reactor

	// containers caches resolved container handles by "database/container". Resolving reads
	// container metadata from the gateway on a miss, so an item operation would otherwise pay
	// for that lookup on every call.
	containers map[string]*C.cosmos_container_ref_t
}

// errTokenCredentialUnsupported is returned when a client is created with a token credential.
//
// The driver supports them (AccountReference::with_credential), but the C ABI does not expose a
// constructor for one: bridging an async Rust TokenCredential through C is deferred upstream, so
// cosmos_account_ref_with_master_key is the only account reference this binding can build.
var errTokenCredentialUnsupported = &Error{
	Code:    CodeClientError,
	Message: "token credentials are not supported by the Cosmos driver yet; use NewClientWithKey",
}

// openDriver acquires the resources a client can acquire locally: the runtime and the account
// reference, which together validate the endpoint and the key. The driver itself is created on
// first use; see [nativeDriver].
func openDriver(cfg driverConfig) (*nativeDriver, error) {
	if cfg.usesTokenCredential() {
		return nil, errTokenCredentialUnsupported
	}

	d := &nativeDriver{}
	if err := d.buildRuntime(); err != nil {
		return nil, err
	}
	if err := d.buildAccount(cfg); err != nil {
		_ = d.close()
		return nil, err
	}
	return d, nil
}

// ensureDriver creates the driver on first use and returns it, or the failure that first attempt
// produced. It is safe for concurrent use.
//
// It reports [CodeClientClosed] once the client is closed rather than a NULL handle, so that a
// caller which reaches it without going through [Client.acquire] cannot pass NULL into the C ABI.
//
//nolint:unused // consumed once operations reach the driver.
func (d *nativeDriver) ensureDriver() (*C.cosmos_driver_t, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return nil, &Error{Code: CodeClientClosed, Message: "the client has been closed"}
	}
	if d.created {
		return d.driver, d.driverErr
	}
	d.created = true

	var richErr *C.cosmos_error_t
	// NULL driver options selects the defaults; per-client options are applied per operation.
	status := C.cosmos_driver_get_or_create_blocking(d.runtime, d.account, nil, &d.driver, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if d.driverErr = statusError(status, richErr, "creating the driver"); d.driverErr != nil {
		return nil, d.driverErr
	}

	// The queue is only needed once operations can run, so it is created with the driver rather
	// than at construction.
	if d.reactor, d.driverErr = newReactor(d.runtime); d.driverErr != nil {
		return nil, d.driverErr
	}
	return d.driver, nil
}

// buildRuntime creates the Tokio runtime the driver executes on.
func (d *nativeDriver) buildRuntime() error {
	var richErr *C.cosmos_error_t
	// A NULL options pointer selects the driver's defaults for every field.
	status := C.cosmos_runtime_build(nil, &d.runtime, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	return statusError(status, richErr, "building the driver runtime")
}

// buildAccount creates the account reference that names the endpoint and carries the key.
func (d *nativeDriver) buildAccount(cfg driverConfig) error {
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
	// Taking the lock waits for a driver creation that is still running, so it is never freed
	// while the C call that produces it is in flight.
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
