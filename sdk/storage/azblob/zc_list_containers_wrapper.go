// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"time"
)

type listContainersSegmentAutoPager struct {
	pager ListContainersSegmentResponsePager
	channel chan ContainerItem
	errChan chan error
	ctx     context.Context

	timeout time.Duration
	timer *time.Timer
}

func (p listContainersSegmentAutoPager) Go() {
	p.timer = time.NewTimer(p.timeout)

	for {
		resp := p.pager.PageResponse()

		if resp.EnumerationResults != nil && resp.EnumerationResults.ContainerItems != nil {
			for _, v := range *resp.EnumerationResults.ContainerItems {
				if p.timeout != 0 {
					p.timer.Reset(p.timeout)
				} else {
					p.timer.C = nil
				}

				select {
				case p.channel <- v:
				case <-p.timer.C:
					p.errChan <- nil

					close(p.channel)
					close(p.errChan)
					return // break the queue
				}
			}
		}

		if !p.pager.NextPage(p.ctx) {
			err := p.pager.Err()
			if err != nil {
				p.errChan <- handleError(err)
			} else {
				p.errChan <- nil
			}

			close(p.errChan)
			close(p.channel)
			return
		}
	}
}
