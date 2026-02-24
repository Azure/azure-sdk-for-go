// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// MultiBalanceTester runs the BalanceTest multiple times against different
// combinations of partition acquisition strategy and number of processors.
//
// NOTE: this test assumes that the Event Hub you're using has 32 partitions.
func MultiBalanceTester(ctx context.Context) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	rounds := fs.Int("rounds", 1, "Number of rounds to run")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	ch := make(chan string, 10000)

	log.SetEvents(EventBalanceTest, azeventhubs.EventConsumer)
	log.SetListener(func(e log.Event, s string) {
		if e == azeventhubs.EventConsumer &&
			!strings.Contains(s, "Asked for") {
			return
		}

		ch <- fmt.Sprintf("[%s] %s", e, s)
	})

	go func() {
		for {
			select {
			case s := <-ch:
				golog.Println(s)
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < *rounds; i++ {
		testData := []struct {
			Processors int
			Strategy   azeventhubs.ProcessorStrategy
		}{
			{32, azeventhubs.ProcessorStrategyGreedy},
			{31, azeventhubs.ProcessorStrategyGreedy},
			{16, azeventhubs.ProcessorStrategyGreedy},
			{5, azeventhubs.ProcessorStrategyGreedy},
			{1, azeventhubs.ProcessorStrategyGreedy},

			{32, azeventhubs.ProcessorStrategyBalanced},
			{31, azeventhubs.ProcessorStrategyBalanced},
			{16, azeventhubs.ProcessorStrategyBalanced},
			{5, azeventhubs.ProcessorStrategyBalanced},
			{1, azeventhubs.ProcessorStrategyBalanced},
		}

		for _, td := range testData {
			log.Writef(EventBalanceTest, "----- BEGIN[%d]: %s, %d processors -----", i, td.Strategy, td.Processors)

			if err := balanceTesterImpl(ctx, td.Processors, td.Strategy); err != nil {
				log.Writef(EventBalanceTest, "----- END[%d]: FAIL: %s, %d processors, %s -----", i, td.Strategy, td.Processors, err)
				return err
			}

			log.Writef(EventBalanceTest, "----- END[%d]: %s, %d processors -----", i, td.Strategy, td.Processors)
		}
	}

	return nil
}
