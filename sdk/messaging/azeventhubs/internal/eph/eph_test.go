// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package eph

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/aad"
	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/stretchr/testify/suite"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
)

const (
	defaultTimeout = 1 * time.Minute
)

type (
	// eventHubSuite encapsulates a end to end test of Event Hubs with build up and tear down of all EH resources
	testSuite struct {
		test.BaseSuite
	}
)

func TestEPH(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) TestRegisterUnRegisterHandler() {
	hub, del := s.RandomHub()
	defer del()

	p, err := s.newInMemoryEPH(*hub.Name)
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	s.Len(p.RegisteredHandlerIDs(), 0, "should have no registered handlers")
	handlerID1, _ := p.RegisterHandler(ctx, func(ctx context.Context, evt *eventhub.Event) error {
		return nil
	})

	handlerID2, _ := p.RegisterHandler(ctx, func(ctx context.Context, evt *eventhub.Event) error {
		return nil
	})

	s.Len(p.RegisteredHandlerIDs(), 2, "should have 2 registered handlers")
	p.UnregisterHandler(ctx, handlerID2)
	s.Require().Len(p.RegisteredHandlerIDs(), 1, "should have 1 registered handlers")
	s.Equal(handlerID1, p.RegisteredHandlerIDs()[0], "should only contain handlerID1")
}

func (s *testSuite) TestNewWithConnectionString() {
	hub, del := s.RandomHub()
	defer del()

	leaserCheckpointer := newMemoryLeaserCheckpointer(DefaultLeaseDuration, new(sharedStore))
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	host, err := NewFromConnectionString(ctx, os.Getenv("EVENTHUB_CONNECTION_STRING")+";EntityPath="+*hub.Name, leaserCheckpointer, leaserCheckpointer)
	s.NoError(err)
	s.NotNil(host)
}

func (s *testSuite) TestSingle() {
	hub, del := s.RandomHub()
	processor, err := s.newInMemoryEPH(*hub.Name)
	s.Require().NoError(err)

	messages, err := s.sendMessages(*hub.Name, 10)
	s.Require().NoError(err)

	var wg sync.WaitGroup
	wg.Add(len(messages))

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	_, _ = processor.RegisterHandler(ctx, func(c context.Context, event *eventhub.Event) error {
		wg.Done()
		return nil
	})

	s.NoError(processor.StartNonBlocking(context.Background()))
	defer func() {
		closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_ = processor.Close(closeContext)
		cancel()
		del()
	}()

	end, _ := ctx.Deadline()
	waitUntil(s.T(), &wg, time.Until(end))
}

func (s *testSuite) TestMultiple() {
	hub, del := s.RandomHub()
	defer del()

	numPartitions := len(*hub.PartitionIds)
	sharedStore := new(sharedStore)
	processors := make(map[string]*EventProcessorHost, numPartitions)
	processorNames := make([]string, numPartitions)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout*2)
	defer cancel()
	for i := 0; i < numPartitions; i++ {
		processor, err := s.newInMemoryEPHWithOptions(*hub.Name, sharedStore)
		s.Require().NoError(err)

		processors[processor.GetName()] = processor
		s.Require().NoError(processor.StartNonBlocking(ctx))
		processorNames[i] = processor.GetName()
	}

	defer func() {
		for _, processor := range processors {
			closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_ = processor.Close(closeContext)
			cancel()
		}
		del()
	}()

	count := 0
	var partitionsByProcessor map[string][]int
	balanced := false
	for {
		<-time.After(3 * time.Second)
		count++
		if count > 50 {
			break
		}

		partitionsByProcessor = make(map[string][]int, len(*hub.PartitionIds))
		for _, processor := range processors {
			partitions := processor.PartitionIDsBeingProcessed()
			partitionInts, err := stringsToInts(partitions)
			s.Require().NoError(err)
			partitionsByProcessor[processor.GetName()] = partitionInts
		}

		if allHaveOnePartition(partitionsByProcessor, numPartitions) {
			balanced = true
			break
		}
	}
	s.Require().True(balanced, "never balanced work within allotted time")

	closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	s.Require().NoError(processors[processorNames[numPartitions-1]].Close(closeContext)) // close the last partition
	delete(processors, processorNames[numPartitions-1])
	cancel()

	count = 0
	balanced = false
	for {
		<-time.After(3 * time.Second)
		count++
		if count > 50 {
			break
		}

		partitionsByProcessor = make(map[string][]int, len(*hub.PartitionIds))
		for _, processor := range processors {
			partitions := processor.PartitionIDsBeingProcessed()
			partitionInts, err := stringsToInts(partitions)
			s.Require().NoError(err)

			partitionsByProcessor[processor.GetName()] = partitionInts
		}

		if allHandled(partitionsByProcessor, len(*hub.PartitionIds)) {
			balanced = true
			break
		}
	}
	if !balanced {
		s.T().Error("didn't balance after closing a processor")
	}
}

func (s *testSuite) sendMessages(hubName string, length int) ([]string, error) {
	client := s.newClient(s.T(), hubName)
	defer func() {
		closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_ = client.Close(closeContext)
		cancel()
	}()

	messages := make([]string, length)
	for i := 0; i < length; i++ {
		messages[i] = s.RandomName("message", 5)
	}

	events := make([]*eventhub.Event, length)
	for idx, msg := range messages {
		events[idx] = eventhub.NewEventFromString(msg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ebi := eventhub.NewEventBatchIterator(events...)
	err := client.SendBatch(ctx, ebi)
	if err != nil {
		return nil, err
	}

	return messages, ctx.Err()
}

func (s *testSuite) newInMemoryEPH(hubName string) (*EventProcessorHost, error) {
	return s.newInMemoryEPHWithOptions(hubName, new(sharedStore))
}

func (s *testSuite) newInMemoryEPHWithOptions(hubName string, store *sharedStore) (*EventProcessorHost, error) {
	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	leaserCheckpointer := newMemoryLeaserCheckpointer(DefaultLeaseDuration, store)
	processor, err := New(ctx, s.Namespace, hubName, provider, leaserCheckpointer, leaserCheckpointer, WithNoBanner())
	if err != nil {
		return nil, err
	}

	return processor, nil
}

func (s *testSuite) newClient(t *testing.T, hubName string, opts ...eventhub.HubOption) *eventhub.Hub {
	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars(), aad.JWTProviderWithAzureEnvironment(&s.Env))
	if err != nil {
		t.Fatal(err)
	}
	return s.newClientWithProvider(t, hubName, provider, opts...)
}

func (s *testSuite) newClientWithProvider(t *testing.T, hubName string, provider auth.TokenProvider, opts ...eventhub.HubOption) *eventhub.Hub {
	opts = append(opts, eventhub.HubWithEnvironment(s.Env))
	client, err := eventhub.NewHub(s.Namespace, hubName, provider, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func waitUntil(t *testing.T, wg *sync.WaitGroup, d time.Duration) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(d):
		t.Error("took longer than " + fmtDuration(d))
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second) / time.Second
	return fmt.Sprintf("%d seconds", d)
}

func allHaveOnePartition(partitionsByProcessor map[string][]int, numberOfPartitions int) bool {
	for _, partitions := range partitionsByProcessor {
		if len(partitions) != 1 {
			return false
		}
	}

	countByPartition := make(map[int]int, numberOfPartitions)
	for i := 0; i < numberOfPartitions; i++ {
		countByPartition[i] = 0
	}
	for _, partitions := range partitionsByProcessor {
		for _, partition := range partitions {
			if count, ok := countByPartition[partition]; ok {
				countByPartition[partition] = count + 1
			}
		}
	}
	for i := 0; i < numberOfPartitions; i++ {
		if countByPartition[i] != 1 {
			return false
		}
	}
	return true
}

func allHandled(partitionsByProcessor map[string][]int, numberOfPartitions int) bool {
	countByPartition := make(map[int]int, numberOfPartitions)
	for i := 0; i < numberOfPartitions; i++ {
		countByPartition[i] = 0
	}
	for _, partitions := range partitionsByProcessor {
		for _, partition := range partitions {
			if count, ok := countByPartition[partition]; ok {
				countByPartition[partition] = count + 1
			}
		}
	}

	for _, count := range countByPartition {
		if count != 1 {
			return false
		}
	}
	return true
}

func stringsToInts(strs []string) ([]int, error) {
	ints := make([]int, len(strs))
	for idx, str := range strs {
		strInt, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		ints[idx] = strInt
	}
	return ints, nil
}
