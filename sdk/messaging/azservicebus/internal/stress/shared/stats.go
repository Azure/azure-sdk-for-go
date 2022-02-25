// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

type Stats struct {
	name string

	Sent     int32
	Received int32
	Errors   int32
}

func NewStats(name string) *Stats {
	return &Stats{
		name: name,
	}
}

func (s *Stats) AddSent(add int32) {
	atomic.AddInt32(&s.Sent, add)
}

func (s *Stats) AddReceived(add int32) {
	atomic.AddInt32(&s.Received, add)
}

func (s *Stats) AddError(reason string, err error) {
	if err != nil {
		return
	}

	atomic.AddInt32(&s.Errors, 1)
}

func (s *Stats) String() string {
	sent := atomic.LoadInt32(&s.Sent)
	received := atomic.LoadInt32(&s.Received)
	errors := atomic.LoadInt32(&s.Errors)

	if s.name == "" {
		return fmt.Sprintf("(s:%d, r:%d, err:%d)", sent, received, errors)
	} else {
		return fmt.Sprintf("(n:%s, s:%d, r:%d, err:%d)", s.name, sent, received, errors)
	}
}

type statsPrinter struct {
	tc  appinsights.TelemetryClient
	mu  sync.RWMutex
	all []*Stats
}

func newStatsPrinter(ctx context.Context, prefix string, interval time.Duration, telemetryClient appinsights.TelemetryClient) *statsPrinter {
	sp := &statsPrinter{
		tc: telemetryClient,
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(interval)

	TickerLoop:
		for range ticker.C {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break TickerLoop
			default:
			}

			sp.mu.RLock()
			sp.PrintStats()
			sp.mu.RUnlock()
		}
	}(ctx)

	return sp
}

func (sp *statsPrinter) PrintStats() {
	log.Printf("Stats:")

	for _, stats := range sp.all {
		log.Printf("  %s", stats.String())

		if stats.Sent > 0 {
			sp.tc.TrackMetric(fmt.Sprintf("%s.TotalSent", stats.name), float64(stats.Sent))
		}

		if stats.Received > 0 {
			sp.tc.TrackMetric(fmt.Sprintf("%s.TotalReceived", stats.name), float64(stats.Received))
		}

		sp.tc.TrackMetric(fmt.Sprintf("%s.TotalErrors", stats.name), float64(stats.Errors))
	}
}

// NewStat creates a new stat with `name` and adds it to the list of statistics that will
// be printed and reported regularly.
func (sp *statsPrinter) NewStat(name string) *Stats {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	stats := NewStats(name)
	sp.all = append(sp.all, stats)

	return stats
}
