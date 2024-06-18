// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	golog "log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

const (
	EventBalanceTest log.Event = "balance.test"
)

// BalanceTester checks that we can properly distribute partitions and
// maintain it over time.
func BalanceTester(ctx context.Context) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	numProcessors := fs.Int("processors", 32, "The # of processor instances to run")
	strategy := fs.String("strategy", string(azeventhubs.ProcessorStrategyBalanced), "The partition acquisition strategy to use (balanced, greedy)")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	log.SetEvents(EventBalanceTest, azeventhubs.EventConsumer)
	log.SetListener(func(e log.Event, s string) {
		// we don't have structured logging in our SDK so this is the most reasonable way to
		// see what partitions each processor
		if e == azeventhubs.EventConsumer &&
			!strings.Contains(s, "Asked for") {
			return
		}

		golog.Printf("[%s] %s", e, s)
	})

	return balanceTesterImpl(ctx, *numProcessors, azeventhubs.ProcessorStrategy(*strategy))
}

func balanceTesterImpl(ctx context.Context, numProcessors int, strategy azeventhubs.ProcessorStrategy) error {
	testData, err := newStressTestData("balancetester", map[string]string{
		"processors": fmt.Sprintf("%d", numProcessors),
		"strategy":   string(strategy),
	})

	if err != nil {
		return err
	}

	args := balanceTester{
		stressTestData: testData,
		numProcessors:  numProcessors,
		strategy:       strategy,
	}

	args.numPartitions, err = func(ctx context.Context) (int, error) {
		client, err := azeventhubs.NewProducerClient(args.Namespace, args.HubName, args.Cred, nil)

		if err != nil {
			return 0, err
		}

		defer func() {
			_ = client.Close(ctx)
		}()

		props, err := client.GetEventHubProperties(ctx, nil)

		if err != nil {
			return 0, err
		}

		return len(props.PartitionIDs), nil
	}(ctx)

	if err != nil {
		return err
	}

	return args.Run(ctx)
}

type balanceTester struct {
	*stressTestData

	strategy      azeventhubs.ProcessorStrategy
	numProcessors int
	numPartitions int
}

func (bt *balanceTester) Run(ctx context.Context) error {
	defer bt.cleanupContainer()

	wg := sync.WaitGroup{}
	failuresChan := make(chan error, bt.numProcessors)

	testCtx, cancelTest := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelTest()

	mu := sync.Mutex{}
	var lastBalanceError error
	startTime := time.Now()

	go func() {
		balancedCount := 0
		var firstBalance time.Duration

	Loop:
		// poll every 5 seconds to see if the checkpoint store is "balanced" (all owners
		// own a fair-share of the partitions).
		for {
			select {
			case <-ctx.Done():
				break Loop
			case <-time.After(5 * time.Second):
				err := bt.checkBalance(ctx)

				if ibErr := (unbalancedError)(nil); errors.As(err, &ibErr) {
					mu.Lock()
					lastBalanceError = err
					mu.Unlock()

					log.Writef(EventBalanceTest, "Balance not achieved, resetting balancedCount: %s", ibErr)
					balancedCount = 0

					bt.TC.TrackEvent("Unbalanced", map[string]string{
						"Message": ibErr.Error(),
					})
					continue
				} else if err != nil {
					mu.Lock()
					lastBalanceError = err
					mu.Unlock()

					bt.TC.TrackException(err)
					break Loop
				}

				if balancedCount == 0 {
					firstBalance = time.Since(startTime)
				}

				balancedCount++
				log.Writef(EventBalanceTest, "Balanced, with %d consecutive checks", balancedCount)

				bt.TC.TrackEvent("Balanced", map[string]string{
					"Count":           fmt.Sprintf("%d", balancedCount),
					"DurationSeconds": fmt.Sprintf("%d", firstBalance/time.Second),
				})

				if balancedCount == 3 {
					log.Writef(EventBalanceTest, "Balanced at %d seconds (approx)", firstBalance/time.Second)

					mu.Lock()
					lastBalanceError = nil
					mu.Unlock()

					cancelTest()
					break Loop
				}
			}
		}
	}()

	for i := 0; i < bt.numProcessors; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if err := bt.process(testCtx, fmt.Sprintf("proc%02d", i)); err != nil {
				failuresChan <- err
				cancelTest()
				return
			}
		}(i)
	}

	wg.Wait()
	close(failuresChan)
	cancelTest()

	// any errors?
	for err := range failuresChan {
		bt.TC.TrackException(err)
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	mu.Lock()
	err := lastBalanceError
	mu.Unlock()

	if err != nil {
		bt.TC.TrackException(err)
		return err
	}

	log.Writef(EventBalanceTest, "BALANCED")
	return nil
}

func (bt *balanceTester) process(ctx context.Context, name string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := azeventhubs.NewConsumerClient(bt.Namespace, bt.HubName, azeventhubs.DefaultConsumerGroup, bt.Cred, &azeventhubs.ConsumerClientOptions{
		InstanceID: name,
	})

	if err != nil {
		return err
	}

	defer func() { _ = client.Close(ctx) }()

	blobClient, err := azblob.NewClient(bt.StorageEndpoint, bt.Cred, nil)

	if err != nil {
		return err
	}

	containerClient := blobClient.ServiceClient().NewContainerClient(bt.runID)

	if _, err := containerClient.Create(ctx, nil); err != nil {
		if !bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
			return err
		}
	}

	blobStore, err := checkpoints.NewBlobStore(containerClient, nil)

	if err != nil {
		return err
	}

	processor, err := azeventhubs.NewProcessor(client, blobStore, &azeventhubs.ProcessorOptions{
		LoadBalancingStrategy: bt.strategy,
	})

	if err != nil {
		return err
	}

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for {
			pc := processor.NextPartitionClient(ctx)

			if pc == nil {
				break
			}

			go bt.keepAlive(ctx, pc)
		}
	}()

	err = processor.Run(ctx)
	cancel()
	<-ch

	return err
}

func (bt *balanceTester) keepAlive(ctx context.Context, pc *azeventhubs.ProcessorPartitionClient) {
	defer func() {
		_ = pc.Close(context.Background())
	}()

	for {
		if _, err := pc.ReceiveEvents(ctx, 1, nil); err != nil {
			break
		}
	}
}

type unbalancedError error

// checkBalance queries the checkpoint store.
// It returns `nil` if no error occurred and the checkpoint store was balanced.
// If the checkpoint store is NOT balanced it returns an unbalancedError
func (bt *balanceTester) checkBalance(ctx context.Context) error {
	blobClient, err := azblob.NewClient(bt.StorageEndpoint, bt.Cred, nil)

	if err != nil {
		return err
	}

	blobStore, err := checkpoints.NewBlobStore(
		blobClient.ServiceClient().NewContainerClient(bt.runID),
		nil)

	if err != nil {
		return err
	}

	ownerships, err := blobStore.ListOwnership(ctx, bt.Namespace, bt.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return err
	}

	stats := bt.summarizeBalance(ownerships)

	if !stats.Balanced {
		return unbalancedError(fmt.Errorf("unbalanced: %s", stats.String()))
	}

	return nil
}

func (bt *balanceTester) cleanupContainer() {
	blobClient, err := azblob.NewClient(bt.StorageEndpoint, bt.Cred, nil)

	if err != nil {
		return
	}

	containerClient := blobClient.ServiceClient().NewContainerClient(bt.runID)

	_, _ = containerClient.Delete(context.Background(), nil)
}

func (bt *balanceTester) summarizeBalance(ownerships []azeventhubs.Ownership) stats {
	counts := map[string]int{}

	for _, o := range ownerships {
		counts[o.OwnerID]++
	}

	// now let's make sure everyone only took a fair share
	min := bt.numPartitions / bt.numProcessors
	max := min

	if bt.numPartitions%bt.numProcessors != 0 {
		max += 1
	}

	tooFew := 0
	tooMany := 0

	for _, owned := range counts {
		if owned < min {
			tooFew++
		} else if owned > max {
			tooMany++
		}
	}

	sum := 0

	for _, v := range counts {
		sum += v
	}

	return stats{
		Processors: fmt.Sprintf("%d/%d", len(counts), bt.numProcessors),
		Partitions: fmt.Sprintf("%d/%d", sum, bt.numPartitions),
		OwnTooFew:  tooFew,
		OwnTooMany: tooMany,
		Balanced: len(counts) == bt.numProcessors &&
			sum == bt.numPartitions &&
			tooFew == 0 &&
			tooMany == 0,
		Raw: counts,
	}
}

type stats struct {
	Processors string
	Partitions string
	OwnTooFew  int
	OwnTooMany int
	Balanced   bool
	Raw        map[string]int
}

func (s *stats) String() string {
	jsonBytes, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
