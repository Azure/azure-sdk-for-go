// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

type ProcessorStrategy string

const (
	// ProcessorStrategyBalanced will attempt to claim a single partition at a time, until each active
	// owner has an equal share of partitions.
	// This is the default strategy.
	ProcessorStrategyBalanced ProcessorStrategy = "balanced"

	// ProcessorStrategyGreedy will attempt to claim as many partitions at a time as it can, ignoring
	// balance.
	ProcessorStrategyGreedy ProcessorStrategy = "greedy"
)

// NewProcessorOptions are the options for the NewProcessor
// function.
type NewProcessorOptions struct {
	// LoadBalancingStrategy dictates how concurrent Processor instances distribute
	// ownership of partitions between them.
	LoadBalancingStrategy ProcessorStrategy

	// UpdateInterval controls how often attempt to claim partitions.
	UpdateInterval time.Duration

	// PartitionExpirationDuration is the amount of time before a partition is considered
	// unowned.
	PartitionExpirationDuration time.Duration

	// StartPositions are the default start positions (configurable per partition, or with an overall
	// default value) if a checkpoint is not found in the CheckpointStore.
	StartPositions StartPositions
}

// StartPositions are used if there is no checkpoint for a partition in
// the checkpoint store.
type StartPositions struct {
	// PerPartition controls the start position for a specific partition,
	// by partition ID. If a partition is not configured here it will default
	// to Default start position.
	PerPartition map[string]StartPosition

	// Default is used if the partition is not found in the PerPartition map.
	Default StartPosition
}

// Processor uses a CheckpointStore, combined with a ConsumerClient, to provide
// automatic load balancing betweeen multiple consumers.
type Processor struct {
	ownershipUpdateInterval time.Duration
	defaultStartPositions   StartPositions
	checkpointStore         CheckpointStore

	// consumerClient is actually a *azeventhubs.ConsumerClient
	// it's an interface here to make testing easier.
	consumerClient consumerClientForProcessor

	nextClients           chan *ProcessorPartitionClient
	consumerClientDetails consumerClientDetails

	runCalled chan struct{}
	lb        *processorLoadBalancer
}

type consumerClientForProcessor interface {
	GetEventHubProperties(ctx context.Context, options *GetEventHubPropertiesOptions) (EventHubProperties, error)
	NewPartitionClient(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error)
	getDetails() consumerClientDetails
}

// NewProcessor creates a Processor.
func NewProcessor(consumerClient *ConsumerClient, checkpointStore CheckpointStore, options *NewProcessorOptions) (*Processor, error) {
	return newProcessorImpl(consumerClient, checkpointStore, options)
}

func newProcessorImpl(consumerClient consumerClientForProcessor, checkpointStore CheckpointStore, options *NewProcessorOptions) (*Processor, error) {
	if options == nil {
		options = &NewProcessorOptions{}
	}

	updateInterval := 10 * time.Second

	if options.UpdateInterval != 0 {
		updateInterval = options.UpdateInterval
	}

	partitionDurationExpiration := time.Minute

	if options.PartitionExpirationDuration != 0 {
		partitionDurationExpiration = options.PartitionExpirationDuration
	}

	startPosPerPartition := map[string]StartPosition{}

	if options.StartPositions.PerPartition != nil {
		for k, v := range options.StartPositions.PerPartition {
			startPosPerPartition[k] = v
		}
	}

	strategy := options.LoadBalancingStrategy

	switch strategy {
	case ProcessorStrategyBalanced:
	case ProcessorStrategyGreedy:
	case "":
		strategy = ProcessorStrategyBalanced
	default:
		return nil, fmt.Errorf("invalid load balancing strategy '%s'", strategy)
	}

	return &Processor{
		ownershipUpdateInterval: updateInterval,
		consumerClient:          consumerClient,
		checkpointStore:         checkpointStore,

		defaultStartPositions: StartPositions{
			PerPartition: startPosPerPartition,
			Default:      options.StartPositions.Default,
		},
		consumerClientDetails: consumerClient.getDetails(),
		runCalled:             make(chan struct{}),
		lb:                    newProcessorLoadBalancer(checkpointStore, consumerClient.getDetails(), strategy, partitionDurationExpiration),

		// `nextClients` will be initialized when the user calls Run() since it needs to query the #
		// of partitions on the Event Hub.
	}, nil
}

// NextPartitionClient will get the next available azeventhubs.PartitionProcessorClient if
// a partition is available or will block until a new one arrives or processor.Run() is
// cancelled.
//
// NOTE: this function will not return any values until processor.Run() is executing. If the
// Run() function is cancelled (or if this function is cancelled) the returned
// ProcessorPartitionClient will be nil.
func (p *Processor) NextPartitionClient(ctx context.Context) *ProcessorPartitionClient {
	<-p.runCalled

	select {
	case nextClient := <-p.nextClients:
		return nextClient
	case <-ctx.Done():
		return nil
	}
}

// Run runs the load balancing loop. Partitions that are claimed can be read using the
// DistributedPartitionClient returned from processor.Next().
//
// NOTE: If this function is cancelled the load balancing loop is cancelled and a nil error
// is returned. processor.NextPartitionClient() will return nil.
func (p *Processor) Run(ctx context.Context) error {
	err := p.runImpl(ctx)

	// the context is the proper way to close down the Run() loop, so it's not
	// an error and doesn't need to be returned.
	if ctx.Err() != nil {
		return nil
	}

	return err
}

func (p *Processor) runImpl(ctx context.Context) error {
	consumers := &sync.Map{}
	defer closeConsumers(ctx, consumers)

	// size the channel to the # of partitions. We can never exceed this size since
	// we'll never reclaim a partition that we already have ownership of.
	eventHubProperties, err := p.initNextClientsCh(ctx)

	if err != nil {
		return err
	}

	// do one dispatch immediately
	if err := p.dispatch(ctx, eventHubProperties, consumers); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(calculateUpdateInterval(p.ownershipUpdateInterval)):
			if err := p.dispatch(ctx, eventHubProperties, consumers); err != nil {
				return err
			}
		}
	}
}

func calculateUpdateInterval(updateInterval time.Duration) time.Duration {
	// Introduce some jitter:  [0.0, 1.0) / 2 = [0.0, 0.5) + 0.8 = [0.8, 1.3)
	// (copied from the retry code for calculating jitter)
	return time.Duration(updateInterval.Seconds() * (rand.Float64()/2 + 0.8) * float64(time.Second))
}

func (p *Processor) initNextClientsCh(ctx context.Context) (EventHubProperties, error) {
	eventHubProperties, err := p.consumerClient.GetEventHubProperties(ctx, nil)

	if err != nil {
		return EventHubProperties{}, err
	}

	p.nextClients = make(chan *ProcessorPartitionClient, len(eventHubProperties.PartitionIDs))
	close(p.runCalled)

	return eventHubProperties, nil
}

func (p *Processor) dispatch(ctx context.Context, eventHubProperties EventHubProperties, consumers *sync.Map) error {
	ownerships, err := p.lb.LoadBalance(ctx, eventHubProperties.PartitionIDs)

	if err != nil {
		return err
	}

	checkpoints, err := p.getCheckpointsMap(ctx)

	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, ownership := range ownerships {
		wg.Add(1)

		go func(o Ownership) {
			defer wg.Done()

			err := p.addPartitionClient(ctx, o, checkpoints, consumers)

			if err != nil {
				azlog.Writef(EventConsumer, "failed to create partition client for partition '%s': %s", o.PartitionID, err.Error())
			}
		}(ownership)
	}

	wg.Wait()

	return nil
}

// addPartitionClient creates a ProcessorPartitionClient
func (p *Processor) addPartitionClient(ctx context.Context, ownership Ownership, checkpoints map[string]Checkpoint, consumers *sync.Map) error {
	processorPartClient := &ProcessorPartitionClient{
		consumerClientDetails: p.consumerClientDetails,
		checkpointStore:       p.checkpointStore,
		innerClient:           nil,
		partitionID:           ownership.PartitionID,
		cleanupFn: func() {
			consumers.Delete(ownership.PartitionID)
		},
	}

	if _, alreadyExists := consumers.LoadOrStore(ownership.PartitionID, processorPartClient); alreadyExists {
		return nil
	}

	partClient, err := p.consumerClient.NewPartitionClient(ownership.PartitionID, &NewPartitionClientOptions{
		StartPosition: p.getStartPosition(checkpoints, ownership),
	})

	if err != nil {
		consumers.Delete(ownership.PartitionID)
		return err
	}

	// make sure we create the link _now_ - if we're stealing we want to stake a claim _now_, rather than
	// later when the user actually calls ReceiveEvents(), since the acquisition of the link is lazy.
	if err := partClient.init(ctx); err != nil {
		consumers.Delete(ownership.PartitionID)
		_ = partClient.Close(ctx)
		return err
	}

	processorPartClient.innerClient = partClient

	select {
	case p.nextClients <- processorPartClient:
		return nil
	default:
		processorPartClient.Close(ctx)
		return fmt.Errorf("partitions channel full, consumer for partition %s could not be returned", ownership.PartitionID)
	}
}

func (p *Processor) getStartPosition(checkpoints map[string]Checkpoint, ownership Ownership) StartPosition {
	startPosition := p.defaultStartPositions.Default
	cp, hasCheckpoint := checkpoints[ownership.PartitionID]

	if hasCheckpoint {
		startPosition = StartPosition{
			SequenceNumber: &cp.SequenceNumber,
		}
	} else if p.defaultStartPositions.PerPartition != nil {
		defaultStartPosition, exists := p.defaultStartPositions.PerPartition[ownership.PartitionID]

		if exists {
			startPosition = defaultStartPosition
		}
	}

	return startPosition
}

func (p *Processor) getCheckpointsMap(ctx context.Context) (map[string]Checkpoint, error) {
	details := p.consumerClient.getDetails()
	checkpoints, err := p.checkpointStore.ListCheckpoints(ctx, details.FullyQualifiedNamespace, details.EventHubName, details.ConsumerGroup, nil)

	if err != nil {
		return nil, err
	}

	m := map[string]Checkpoint{}

	for _, cp := range checkpoints {
		m[cp.PartitionID] = cp
	}

	return m, nil
}

func closeConsumers(ctx context.Context, consumersMap *sync.Map) {
	consumersMap.Range(func(key, value any) bool {
		client := value.(*ProcessorPartitionClient)

		if client != nil {
			client.Close(ctx)
		}

		return true
	})
}
