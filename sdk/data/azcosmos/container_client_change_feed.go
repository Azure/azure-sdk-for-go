// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
)

// NewChangeFeedPager reads changes for an explicit scope and initial position.
// Construction performs no I/O; NextPage reports argument errors. Nil options use LatestVersion.
// Complete keys and full-container scopes follow native routing. HPK-prefix routing follows the
// Rust driver's logical-key path; multi-range prefix fan-out is not established by this ABI.
// Defer Close. Empty and HTTP 304 pages remain usable; callers choose when to poll again.
func (c *ContainerClient) NewChangeFeedPager(scope FeedScope, startFrom ChangeFeedStartFrom, options *ChangeFeedOptions) *ChangeFeedPager {
	req, err := newChangeFeedRequest(scope, startFrom, options)
	req.databaseID, req.containerID = c.database.id, c.id
	return &ChangeFeedPager{client: c.database.client, req: req, validationErr: err}
}

// ChangeFeedPager owns a retained native change-feed plan.
// Iteration and checkpoint calls must not be concurrent. Close waits for an active Go call.
// Context abandonment terminates the pager without cancelling admitted native execution.
type ChangeFeedPager struct {
	mu            sync.Mutex
	client        *Client
	req           changeFeedRequest
	validationErr error
	cursor        changeFeedCursor
	done          bool
}

type changeFeedCursor interface {
	next(context.Context) (ChangeFeedResponse, error)
	checkpoint(context.Context) (string, error)
	close()
}

// More reports whether the pager is open, not whether changes currently exist.
func (p *ChangeFeedPager) More() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return !p.done
}

// NextPage advances one native page. Empty and HTTP 304 pages never mark permanent exhaustion.
// Failures after execution starts terminate the pager; retries belong to the Rust driver.
// The context bounds the Go wait. Client.Close can wait for native work after cancellation.
func (p *ChangeFeedPager) NextPage(ctx context.Context) (ChangeFeedResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return ChangeFeedResponse{}, p.validationErr
	}
	release, err := p.client.acquire()
	if err != nil {
		return ChangeFeedResponse{}, err
	}
	defer release()
	if p.done {
		return ChangeFeedResponse{}, &Error{Code: CodeClientError, Message: "azcosmos: change feed pager is closed"}
	}
	if err := ctx.Err(); err != nil {
		return ChangeFeedResponse{}, err
	}
	ctx, cancel := contextWithEndToEndTimeout(ctx, p.req.options.Operation.EndToEndTimeout)
	defer cancel()
	var setup Response
	if p.cursor == nil {
		p.cursor, setup, err = p.client.openChangeFeed(ctx, &p.req)
		if err != nil {
			p.finish()
			return ChangeFeedResponse{}, addFeedSetupCharge(err, setup, "fetching change feed page")
		}
	}
	page, err := p.cursor.next(ctx)
	if err != nil {
		p.finish()
		return ChangeFeedResponse{}, addFeedSetupCharge(err, setup, "fetching change feed page")
	}
	page.RequestCharge += setup.RequestCharge
	if page.ActivityID == "" && setup.RequestCharge != 0 {
		page.ActivityID = setup.ActivityID
	}
	return page, nil
}

// ContinuationToken snapshots progress delivered to Go without advancing the feed.
// Call after a successful NextPage, including an idle page, and before Close.
// Resume using the same container, scope and mode. StartFrom is validated but the token's saved
// positions take precedence. Go does not parse tokens, restore their scope or enforce scope equality.
// The checkpoint does not acknowledge application processing or guarantee exactly-once delivery.
func (p *ChangeFeedPager) ContinuationToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.validationErr != nil {
		return "", p.validationErr
	}
	release, err := p.client.acquire()
	if err != nil {
		return "", err
	}
	defer release()
	if p.done || p.cursor == nil {
		return "", &Error{Code: CodeClientError, Message: "azcosmos: checkpoint requires an open change feed pager"}
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	ctx, cancel := contextWithEndToEndTimeout(ctx, p.req.options.Operation.EndToEndTimeout)
	defer cancel()
	token, err := p.cursor.checkpoint(ctx)
	if err != nil {
		p.finish()
	}
	return token, err
}

// Close releases the retained plan. It is idempotent and safe alongside Client.Close.
func (p *ChangeFeedPager) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.finish()
	return nil
}

func (p *ChangeFeedPager) finish() {
	p.done = true
	if p.cursor != nil {
		p.cursor.close()
		p.cursor = nil
	}
}
