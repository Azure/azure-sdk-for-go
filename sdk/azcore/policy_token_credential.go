// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// TODO: move to azidentity

var _ TokenCredential = (*tokenCredentialWithRefresh)(nil)
var _ TokenCredential = (*tokenCredential)(nil)

// TokenRefresher represents a callback method that you write; this method is called periodically
// so you can refresh the token credential's value.
type TokenRefresher func() (string, time.Duration)

// TokenCredential represents a token credential (which is also a pipeline.Factory).
type TokenCredential interface {
	Credential
	Token() string
	SetToken(newToken string)
}

// NewTokenCredential creates a token credential for use with role-based access control (RBAC) access to Azure Storage
// resources. You initialize the TokenCredential with an initial token value. If you pass a non-nil value for
// tokenRefresher, then the function you pass will be called immediately so it can refresh and change the
// TokenCredential's token value by calling SetToken. Your tokenRefresher function must return a time.Duration
// indicating how long the TokenCredential object should wait before calling your tokenRefresher function again.
// If your tokenRefresher callback fails to refresh the token, you can return a duration of 0 to stop your
// TokenCredential object from ever invoking tokenRefresher again. Also, oen way to deal with failing to refresh a
// token is to cancel a context.Context object used by requests that have the TokenCredential object in their pipeline.
func NewTokenCredential(initialToken string, tokenRefresher TokenRefresher) TokenCredential {
	tc := &tokenCredential{}
	tc.SetToken(initialToken) // We don't set it above to guarantee atomicity
	if tokenRefresher == nil {
		return tc // If no callback specified, return the simple tokenCredential
	}

	tcwr := &tokenCredentialWithRefresh{token: tc}
	tcwr.token.startRefresh(tokenRefresher)
	runtime.SetFinalizer(tcwr, func(deadTC *tokenCredentialWithRefresh) {
		deadTC.token.stopRefresh()
		deadTC.token = nil //  Sanity (not really required)
	})
	return tcwr
}

/*/////////////////////////////////////////////////////////////////////////////

Why are there 2 structures (tokenCredentialWithRefresh and tokenCredential)?

Having a timer invoke an object's method keeps the object from getting GC'd and
we'd like the pipeline to control the lifetime of the token credential. That is,
if no pipeline is keeping the TokenCredential object alive, then the times should
not keep the TC object alive any longer.

To get this behavior, the pipeline holds on to a TC wrapper object (tokenCredentialWithRefesh);
which refers to the real tokenCredential object. The timer keeps the real tokenCredential
object from getting GC'd. However, if the tokenCredentialWithRefresh wrapper gets GC'd, it's
Finalize method stops the timer which allows the real tokenCredential object to also be GC'd.

NOTE: This design requires that there is a 1:1 realtionship between wrapper and real object.
If multiple wrappers existed for a single real TC, then when any werapper got GC'd it would
stop the timer and this is not desired behavior as other client objects may still be using the TC.

/////////////////////////////////////////////////////////////////////////////*/

// tokenCredentialWithRefresh is a wrapper over a token credential.
// When this wrapper object gets GC'd, it stops the tokenCredential's timer
// which allows the tokenCredential object to also be GC'd.
type tokenCredentialWithRefresh struct {
	token *tokenCredential
}

// marker satisfies the Credential interface making Credential policies "special"
func (p *tokenCredentialWithRefresh) marker() {}

// Token returns the current token value
func (c *tokenCredentialWithRefresh) Token() string { return c.token.Token() }

// SetToken changes the current token value
func (c *tokenCredentialWithRefresh) SetToken(token string) { c.token.SetToken(token) }

// Do ...
func (c *tokenCredentialWithRefresh) Do(ctx context.Context, req *Request) (*Response, error) {
	return c.token.Do(ctx, req)
}

///////////////////////////////////////////////////////////////////////////////

// tokenCredential is a pipeline.Factory is the credential's policy factory.
type tokenCredential struct {
	token atomic.Value

	// The members below are only used if the user specified a tokenRefresher callback function.
	timer          *time.Timer
	tokenRefresher TokenRefresher
	lock           sync.Mutex
	stopped        bool
}

// marker satisfies the Credential interface making Credential policies "special"
func (c *tokenCredential) marker() {}

// Token returns the current token value
func (c *tokenCredential) Token() string { return c.token.Load().(string) }

// SetToken changes the current token value
func (c *tokenCredential) SetToken(token string) { c.token.Store(token) }

// startRefresh calls refresh which immediately calls tokenRefresher
// and then starts a timer to call tokenRefresher in the future.
func (c *tokenCredential) startRefresh(tokenRefresher TokenRefresher) {
	c.tokenRefresher = tokenRefresher
	c.stopped = false // In case user calls StartRefresh, StopRefresh, & then StartRefresh again
	c.refresh()
}

// refresh calls the user's tokenRefresher so they can refresh the token (returning it) and then
// starts another time (based on the returned duration) in order to refresh the token again in the future.
func (c *tokenCredential) refresh() {
	newToken, delay := c.tokenRefresher() // Invoke the user's refresh callback outside of the lock
	if delay > 0 {                        // If duration is 0 or negative, refresher wants to not be called again
		c.lock.Lock()
		if !c.stopped {
			c.SetToken(newToken)
			c.timer = time.AfterFunc(delay, c.refresh)
		}
		c.lock.Unlock()
	}
}

// stopRefresh stops any pending timer and sets stopped field to true to prevent
// any new timer from starting.
// NOTE: Stopping the timer allows the GC to destroy the tokenCredential object.
func (c *tokenCredential) stopRefresh() {
	c.lock.Lock()
	c.stopped = true
	if c.timer != nil {
		c.timer.Stop()
	}
	c.lock.Unlock()
}

// Do ...
func (c *tokenCredential) Do(ctx context.Context, req *Request) (*Response, error) {
	if req.Request.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, errors.New("token credentials require a URL using the https protocol scheme")
	}
	req.Request.Header[headerAuthorization] = []string{"Bearer " + c.Token()}
	return req.Do(ctx)
}
