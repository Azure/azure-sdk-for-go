// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const EventHubNamespace = "eventhHubsTest.servicebus.windows.net"
const EventHubName = "eventHub"

// Some useful queries:

/*

Give me a bar chart that shows me how many claims we grabbed per worker.
(this'll give you some confidence that you're seeing contention).

// pick a round number to check out
let round="1";
// timeframe.
let startAt=ago(24h);
// this is just your pod name.
let appRoleInstance="go18-checkpoints-goeh-5-4sjnn";
AppMetrics
| where TimeGenerated > startAt
| where AppRoleInstance == appRoleInstance
| where Name in ('stress.claimed')
| extend Round=tostring(Properties["Round"]), Worker=tostring(Properties["Worker"])
| where Round == round
| summarize Value=sum(Max) by Worker
| order by Value desc
| render columnchart
*/

type checkpointFlags struct {
	Rounds     int
	Partitions int
	Workers    int
}

func parseCheckpointsFlags() (checkpointFlags, error) {
	fs := flag.NewFlagSet("checkpoints", flag.ContinueOnError)

	cpFlags := checkpointFlags{}

	fs.IntVar(&cpFlags.Rounds, "rounds", 100, "Number of rounds.")
	fs.IntVar(&cpFlags.Partitions, "partitions", 100, "Number of partitions. This is just blob data so it can be an arbitrary number, not limited by Event Hubs.")
	fs.IntVar(&cpFlags.Workers, "workers", 100, "Number of workers that will attempt to claim partitions.")

	if err := fs.Parse(os.Args[2:]); err != nil {
		return checkpointFlags{}, err
	}

	return cpFlags, nil
}

// CheckpointStressTester ensures our test code to migrate from our (buggy) improperly-cased-blobs works in concurrent scenarios.
func CheckpointStressTester(ctx context.Context) error {
	cpFlags, err := parseCheckpointsFlags()

	if err != nil {
		return err
	}

	testData, err := newStressTestData("checkpoints", false, map[string]string{
		"Rounds":     fmt.Sprintf("%d", cpFlags.Rounds),
		"Partitions": fmt.Sprintf("%d", cpFlags.Partitions),
	})

	if err != nil {
		return err
	}

	log.Printf("Starting checkpoints test %s, rounds: %d, partitions: %d, workers: %d", testData.runID, cpFlags.Rounds, cpFlags.Partitions, cpFlags.Workers)
	defer log.Printf("Done with checkpoints test %s", testData.runID)

	if err != nil {
		return err
	}

	defer testData.Close()

	if _, err := testData.CC.Create(ctx, nil); err != nil {
		return err
	}

	defer func() {
		_, _ = testData.CC.Delete(context.Background(), nil)
	}()

	for round := 0; round < cpFlags.Rounds; round++ {
		fn := func(round int) error {
			if err := generateLegacyState(ctx, testData.CC, cpFlags); err != nil {
				return err
			}

			getShuffledOwnerships, err := newOwnershipShuffler(ctx, testData)

			if err != nil {
				return err
			}

			remaining := int64(cpFlags.Partitions)

			wg := &sync.WaitGroup{}

			for workerID := 0; workerID < cpFlags.Workers; workerID++ {
				wg.Add(1)

				work := func(worker int) error {
					defer wg.Done()

					for {
						checkpointStore, err := checkpoints.NewBlobStore(testData.MustNewCC(), nil)

						if err != nil {
							return fmt.Errorf("[r:%d,w:%d] failed to create blobstore: %s", round, worker, err)
						}

						// claiming is basically 1 call at a time with each blob. Since we're randomized (above)
						// we should get good distribution amongst all of the claimants.
						claimed, err := checkpointStore.ClaimOwnership(ctx, getShuffledOwnerships(), nil)

						if err != nil {
							return fmt.Errorf("[r:%d,w:%d] failed to claim ownership: %s", round, worker, err)
						}

						testData.TC.TrackMetric(string(MetricStressClaimed), float64(len(claimed)), map[string]string{
							"Round":  fmt.Sprintf("%d", round),
							"Worker": fmt.Sprintf("%d", worker),
						})

						n := atomic.AddInt64(&remaining, -int64(len(claimed)))

						if n == 0 {
							return nil
						} else if n < 0 {
							// we've somehow over-claimed.
							return fmt.Errorf("partitions were overclaimed")
						}
					}
				}

				go func(workerID int) {
					if err := work(workerID); err != nil {
						testData.TC.Track(&appinsights.ExceptionTelemetry{
							BaseTelemetry: appinsights.BaseTelemetry{
								Properties: map[string]string{
									"Round":     fmt.Sprintf("%d", round),
									"Remaining": fmt.Sprintf("%d", remaining),
								},
							},
							Error: err,
						})
					}
				}(workerID)
			}

			wg.Wait()

			log.Printf("Round %d ended, remaining: %d", round, remaining)

			// everything should be balanced.
			if remaining == 0 {
				// validate that blob storage _only_ has lowercase checkpoints.
				var badBlobPaths []string

				blobsPager := testData.CC.NewListBlobsFlatPager(nil)

				for blobsPager.More() {
					page, err := blobsPager.NextPage(context.Background())

					if err != nil {
						return err
					}

					for _, blob := range page.Segment.BlobItems {
						if strings.ToLower(*blob.Name) != *blob.Name {
							badBlobPaths = append(badBlobPaths, *blob.Name)
						}
					}
				}

				if len(badBlobPaths) > 0 {
					return fmt.Errorf("found blob paths that were NOT lowercase: %s", strings.Join(badBlobPaths, "\n"))
				}

				testData.TC.TrackEvent(string(MetricStressSuccess), map[string]string{
					"Round": fmt.Sprintf("%d", round),
				})
			} else if remaining < 0 {
				testData.TC.Track(&appinsights.ExceptionTelemetry{
					BaseTelemetry: appinsights.BaseTelemetry{
						Properties: map[string]string{
							"Round":     fmt.Sprintf("%d", round),
							"Remaining": fmt.Sprintf("%d", remaining),
						},
					},
					Error: fmt.Errorf("Partitions were overclaimed"),
				})
			}
			return nil
		}

		if err := fn(round); err != nil {
			return err
		}
	}

	return nil
}

func generateLegacyState(ctx context.Context, cc *container.Client, cpFlags checkpointFlags) error {
	// baseCheckpoint := azeventhubs.Checkpoint{
	// 	EventHubName:            "eventHub",
	// 	ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
	// 	FullyQualifiedNamespace: "eventhHubsTest.servicebus.windows.net",
	// 	PartitionID:             "0",
	// }

	// clear out the container first - we've got dirty state, etc..
	pager := cc.NewListBlobsFlatPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)

		if err != nil {
			return err
		}

		for _, item := range page.Segment.BlobItems {
			if _, err := cc.NewBlobClient(*item.Name).Delete(ctx, nil); err != nil {
				return err
			}
		}
	}

	baseOwnership := azeventhubs.Ownership{
		EventHubName:            EventHubName,
		ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
		FullyQualifiedNamespace: EventHubNamespace,
		OwnerID:                 "owner id",
	}

	// we're going to simulate an Event Hub but since we're not actually _using_ Event Hubs for
	// this test we can generate any number of partitions, even if it's impossible.
	for i := 0; i < cpFlags.Partitions; i++ {
		o := baseOwnership
		o.PartitionID = fmt.Sprintf("%d", i)

		blobPath := fmt.Sprintf("%s/%s/%s/ownership/%s", o.FullyQualifiedNamespace, o.EventHubName, o.ConsumerGroup, o.PartitionID)

		bc := cc.NewBlockBlobClient(blobPath)

		_, err := bc.UploadBuffer(context.Background(), nil, &blockblob.UploadBufferOptions{
			Metadata: map[string]*string{"ownerid": &o.OwnerID},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func newOwnershipShuffler(ctx context.Context, testData *stressTestData) (func() []azeventhubs.Ownership, error) {
	checkpointStore, err := checkpoints.NewBlobStore(testData.CC, nil)

	if err != nil {
		return nil, err
	}

	baseOwnerships, err := checkpointStore.ListOwnership(ctx, EventHubNamespace, EventHubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return nil, err
	}

	return func() []azeventhubs.Ownership {
		// shuffle the ownerships so our claiming is a bit random.
		shuffledOwnerships := make([]azeventhubs.Ownership, len(baseOwnerships))
		copy(shuffledOwnerships, baseOwnerships)

		rand.Shuffle(len(shuffledOwnerships), func(i, j int) { shuffledOwnerships[i] = shuffledOwnerships[j] })

		return shuffledOwnerships
	}, nil
}
