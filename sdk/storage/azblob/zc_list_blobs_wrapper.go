package azblob

import (
	"context"
	"time"
)

type listBlobsFlatSegmentAutoPager struct {
	pager   ListBlobsFlatSegmentResponsePager
	channel chan BlobItemInternal
	errChan chan error
	ctx context.Context

	// Set to 0 for no time-out
	timeout time.Duration
	timer   *time.Timer
}

type listBlobsHierarchySegmentAutoPager struct {
	pager   ListBlobsHierarchySegmentResponsePager
	channel chan BlobItemInternal
	errChan chan error
	ctx context.Context

	// Set to 0 for no time-out
	timeout time.Duration
	timer   *time.Timer
}

func (p listBlobsFlatSegmentAutoPager) Go() {
	p.timer = time.NewTimer(p.timeout)

	for {
		resp := p.pager.PageResponse()

		if resp.RawResponse != nil {
			for _, v := range *resp.EnumerationResults.Segment.BlobItems {
				p.timer.Reset(p.timeout)

				select {
				case p.channel <- v:
				case <-p.timer.C:
					close(p.errChan)
					close(p.channel)
					return // break the queue
				}
			}
		}

		if !p.pager.NextPage(p.ctx) {
			err := p.pager.Err()
			if err != nil {
				p.errChan <- err
			}

			close(p.errChan)
			close(p.channel)
			return
		}
	}
}

func (p listBlobsHierarchySegmentAutoPager) Go() {
	p.timer = time.NewTimer(p.timeout)

	// Stop it immediately
	// This way, as the user requested, we just don't time out.
	if p.timeout == 0 {
		p.timer.Stop()
	}

	for {
		resp := p.pager.PageResponse()

		if resp.RawResponse != nil {
			for _, v := range *resp.EnumerationResults.Segment.BlobItems {
				if p.timeout != 0 {
					p.timer.Reset(p.timeout)
				}

				select {
				case p.channel <- v:
				case <-p.timer.C:
					close(p.errChan)
					close(p.channel)
					return // break the queue
				}
			}
		}

		if !p.pager.NextPage(p.ctx) {
			err := p.pager.Err()
			if err != nil {
				p.errChan <- err
			}

			close(p.errChan)
			close(p.channel)
			return
		}
	}
}
