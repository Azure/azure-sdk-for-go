// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package storage

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/aad"
	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-storage-blob-go/azblob"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/eph"
)

const (
	defaultTimeout = 1 * time.Minute
)

func (ts *testSuite) TestSingle() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	hub, delHub := ts.RandomHub()
	delContainer := ts.newTestContainerByName(*hub.Name)
	defer delContainer()

	processor, err := ts.newStorageBackedEPH(*hub.Name, *hub.Name)
	ts.Require().NoError(err)
	defer func() {
		closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := processor.Close(closeContext); err != nil {
			ts.Error(err)
		}
		cancel()
		delHub()
	}()

	messages, err := ts.sendMessages(*hub.Name, 10)
	ts.Require().NoError(err)

	var wg sync.WaitGroup
	wg.Add(len(messages))

	_, err = processor.RegisterHandler(ctx, func(c context.Context, event *eventhub.Event) error {
		wg.Done()
		return nil
	})
	ts.Require().NoError(err)

	ts.Require().NoError(processor.StartNonBlocking(ctx))
	end, _ := ctx.Deadline()
	waitUntil(ts.T(), &wg, time.Until(end))
}

func (ts *testSuite) TestMultiple() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout*2)
	defer cancel()

	hub, delHub := ts.RandomHub()
	defer delHub()
	delContainer := ts.newTestContainerByName(*hub.Name)
	defer delContainer()

	cred, err := NewAADSASCredential(ts.SubscriptionID, ts.ResourceGroupName, ts.AccountName, *hub.Name, AADSASCredentialWithEnvironmentVars())
	ts.Require().NoError(err)
	numPartitions := len(*hub.PartitionIds)
	processors := make(map[string]*eph.EventProcessorHost, numPartitions)
	processorNames := make([]string, numPartitions)
	for i := 0; i < numPartitions; i++ {
		leaserCheckpointer, err := NewStorageLeaserCheckpointer(cred, ts.AccountName, *hub.Name, ts.Env)
		ts.Require().NoError(err)

		processor, err := ts.newStorageBackedEPHOptions(*hub.Name, leaserCheckpointer, leaserCheckpointer)
		ts.Require().NoError(err)

		processors[processor.GetName()] = processor
		ts.Require().NoError(processor.StartNonBlocking(ctx))
		processorNames[i] = processor.GetName()
	}

	defer func() {
		for _, processor := range processors {
			closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := processor.Close(closeContext); err != nil {
				ts.Error(err)
			}
			cancel()
		}
		delHub()
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
			ts.Require().NoError(err)
			partitionsByProcessor[processor.GetName()] = partitionInts
		}

		if allHaveOnePartition(partitionsByProcessor, numPartitions) {
			balanced = true
			break
		}
	}
	if !balanced {
		ts.Fail("never balanced work within allotted time")
		return
	}

	closeContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ts.Require().NoError(processors[processorNames[numPartitions-1]].Close(closeContext)) // close the last partition
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
			ts.Require().NoError(err)
			partitionsByProcessor[processor.GetName()] = partitionInts
		}

		if allHandled(partitionsByProcessor, len(*hub.PartitionIds)) {
			balanced = true
			break
		}
	}
	if !balanced {
		ts.Fail("didn't balance after closing a processor")
	}
}

func (ts *testSuite) newTestContainerByName(containerName string) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cred, err := NewAADSASCredential(ts.SubscriptionID, ts.ResourceGroupName, ts.AccountName, containerName, AADSASCredentialWithEnvironmentVars())
	ts.Require().NoError(err)

	pipeline := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	fooURL, err := url.Parse("https://" + ts.AccountName + ".blob." + ts.Env.StorageEndpointSuffix + "/" + containerName)
	ts.NoError(err)

	containerURL := azblob.NewContainerURL(*fooURL, pipeline)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	ts.NoError(err)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if res, err := containerURL.Delete(ctx, azblob.ContainerAccessConditions{}); err != nil {
			msg := "error deleting container url"
			if res != nil {
				msg = fmt.Sprintf("status code: %q; error code: %q", res.StatusCode(), res.ErrorCode())
			}
			ts.NoError(err, msg)
		}
	}
}

func (ts *testSuite) sendMessages(hubName string, length int) ([]string, error) {
	client := ts.newClient(ts.T(), hubName)
	defer func() {
		closeContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		ts.NoError(client.Close(closeContext))
		cancel()
	}()

	messages := make([]string, length)
	for i := 0; i < length; i++ {
		messages[i] = ts.RandomName("message", 5)
	}

	events := make([]*eventhub.Event, length)
	for idx, msg := range messages {
		events[idx] = eventhub.NewEventFromString(msg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ebi := eventhub.NewEventBatchIterator(events...)
	ts.NoError(client.SendBatch(ctx, ebi))

	return messages, ctx.Err()
}

func (ts *testSuite) newStorageBackedEPH(hubName, containerName string) (*eph.EventProcessorHost, error) {
	cred, err := NewAADSASCredential(ts.SubscriptionID, ts.ResourceGroupName, ts.AccountName, containerName, AADSASCredentialWithEnvironmentVars())
	ts.Require().NoError(err)
	leaserCheckpointer, err := NewStorageLeaserCheckpointer(cred, ts.AccountName, containerName, ts.Env)
	ts.Require().NoError(err)
	return ts.newStorageBackedEPHOptions(hubName, leaserCheckpointer, leaserCheckpointer)
}

func (ts *testSuite) newStorageBackedEPHOptions(hubName string, leaser eph.Leaser, checkpointer eph.Checkpointer) (*eph.EventProcessorHost, error) {
	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	processor, err := eph.New(ctx, ts.Namespace, hubName, provider, leaser, checkpointer, eph.WithNoBanner())
	if err != nil {
		return nil, err
	}

	return processor, nil
}

func (ts *testSuite) newClient(t *testing.T, hubName string, opts ...eventhub.HubOption) *eventhub.Hub {
	provider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars(), aad.JWTProviderWithAzureEnvironment(&ts.Env))
	ts.Require().NoError(err)
	return ts.newClientWithProvider(t, hubName, provider, opts...)
}

func (ts *testSuite) newClientWithProvider(t *testing.T, hubName string, provider auth.TokenProvider, opts ...eventhub.HubOption) *eventhub.Hub {
	opts = append(opts, eventhub.HubWithEnvironment(ts.Env))
	client, err := eventhub.NewHub(ts.Namespace, hubName, provider, opts...)
	ts.Require().NoError(err)
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

	//var keys []string
	//for key := range partitionsByProcessor {
	//	keys = append(keys, key)
	//}
	//sort.Strings(keys)
	//for _, key := range keys {
	//	ints := partitionsByProcessor[key]
	//	sort.Ints(ints)
	//	fmt.Printf("Processor: %q, Partitions %+v\n", key, ints)
	//}
	//fmt.Println("========================================")
	//fmt.Println("========================================")

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
