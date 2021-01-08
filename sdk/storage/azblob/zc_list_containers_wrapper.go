package azblob

import (
	"context"
	"time"
)

type listContainersSegmentAutoPager struct {
	pager ListContainersSegmentResponsePager
	channel chan ContainerItem
	errChan chan error
	ctx context.Context

	timeout time.Duration
	timer *time.Timer
}

func (p listContainersSegmentAutoPager) Go() {
	p.timer = time.NewTimer(p.timeout)

	for {
		resp := p.pager.PageResponse()

		for _, v := range *resp.EnumerationResults.ContainerItems {
			p.timer.Reset(p.timeout)

			select {
			case p.channel <- v:
			case <-p.timer.C:
				close(p.errChan)
				close(p.channel)
				return // break the queue
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
