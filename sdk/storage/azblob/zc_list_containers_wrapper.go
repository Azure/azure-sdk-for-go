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

// func (c ContainerItem) Duplicate() ContainerItem {
// 	out := ContainerItem{
// 		Deleted:    nil,
// 		Metadata:   nil,
// 		Name:       nil,
// 		Properties: nil,
// 		Version:    nil,
// 	}
//
// 	out.Name = &*c.Name
//
// 	if out.Metadata != nil {
// 		meta := make(map[string]string)
// 		for k, v := range *c.Metadata {
// 			meta[k] = v
// 		}
// 		out.Metadata = &meta
// 	}
//
// 	if out.Properties != nil {
// 		out.Properties.
// 	}
//
// }

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
				p.errChan <- err
			} else {
				p.errChan <- nil
			}

			close(p.errChan)
			close(p.channel)
			return
		}
	}
}
