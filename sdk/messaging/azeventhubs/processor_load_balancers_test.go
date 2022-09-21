// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProcessorLoadBalancers_Greedy_EnoughUnownedPartitions(t *testing.T) {
	cps := newCheckpointStoreForTest()

	_, err := cps.ClaimOwnership(context.Background(), []Ownership{
		newTestOwnership("0", "some-other-client"),
		newTestOwnership("3", "some-other-client"),
	}, nil)
	require.NoError(t, err)

	lb := newProcessorLoadBalancer(cps, newTestConsumerDetails("new-client"), ProcessorStrategyGreedy, time.Hour)

	// "0" and "3" are already claimed, so we'll pick up the 2 free partitions.
	ownerships, err := lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3"})
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2"}, mapToStrings(ownerships, extractPartitionID))

	finalOwnerships, err := cps.ListOwnership(context.Background(), "fqdn", "event-hub", "consumer-group", nil)
	require.NoError(t, err)

	require.Equal(t, 4, len(finalOwnerships), "all partitions claimed")
}

func TestProcessorLoadBalancers_Balanced_UnownedPartitions(t *testing.T) {
	cps := newCheckpointStoreForTest()

	_, err := cps.ClaimOwnership(context.Background(), []Ownership{
		newTestOwnership("0", "some-other-client"),
		newTestOwnership("3", "some-other-client"),
	}, nil)
	require.NoError(t, err)

	lb := newProcessorLoadBalancer(cps, newTestConsumerDetails("new-client"), ProcessorStrategyBalanced, time.Hour)

	// "0" and "3" are already claimed, so we'll pick up one partition each time we load balance
	ownerships, err := lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3"})
	require.NoError(t, err)
	require.Equal(t, 1, len(ownerships))

	ownerships, err = lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3"})
	require.NoError(t, err)
	require.Equal(t, 2, len(ownerships))

	requireBalanced(t, cps, 4, 2)
}

func TestProcessorLoadBalancers_Greedy_ForcedToSteal(t *testing.T) {
	cps := newCheckpointStoreForTest()

	const someOtherClientID = "some-other-client-id"
	const stealingClientID = "stealing-client-id"

	_, err := cps.ClaimOwnership(context.Background(), []Ownership{
		// All the partitions are owned.
		//
		// You can picture this happening if there was only one consumer for
		// awhile and then we brought up a second instance to try to distribute
		// the load.
		//
		// Since _all_ partitions are owned we need to steal some, which can result
		// in some flux.
		newTestOwnership("0", someOtherClientID),
		newTestOwnership("1", someOtherClientID),
		newTestOwnership("2", someOtherClientID),
		newTestOwnership("3", someOtherClientID),
		newTestOwnership("4", someOtherClientID),
	}, nil)
	require.NoError(t, err)

	lb := newProcessorLoadBalancer(cps, newTestConsumerDetails(stealingClientID), ProcessorStrategyGreedy, time.Hour)

	ownerships, err := lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3", "4"})
	require.NoError(t, err)
	require.NotEmpty(t, mapToStrings(ownerships, extractPartitionID))

	finalOwnerships, err := cps.ListOwnership(context.Background(), "fqdn", "event-hub", "consumer-group", nil)
	require.NoError(t, err)

	ownersMap := groupBy(finalOwnerships, extractOwnerID)

	// there should be no partitions in common
	require.Empty(t, findCommon(
		mapToStrings(ownersMap[someOtherClientID], extractPartitionID),
		mapToStrings(ownersMap[stealingClientID], extractPartitionID),
	))

	require.Equal(t, 5, len(finalOwnerships), "all partitions claimed")
}

func TestProcessorLoadBalancers_AnyStrategy_GrabExpiredPartition(t *testing.T) {
	for _, strategy := range []ProcessorStrategy{ProcessorStrategyBalanced, ProcessorStrategyGreedy} {
		t.Run(string(strategy), func(t *testing.T) {
			cps := newCheckpointStoreForTest()

			const clientA = "clientA"
			const clientB = "clientB"
			const clientCWithExpiredPartition = "clientC"

			middleOwnership := newTestOwnership("2", clientCWithExpiredPartition)

			_, err := cps.ClaimOwnership(context.Background(), []Ownership{
				newTestOwnership("0", clientA),
				newTestOwnership("1", clientA),
				middleOwnership,
				newTestOwnership("3", clientB),
				newTestOwnership("4", clientB),
			}, nil)
			require.NoError(t, err)

			// expire the middle partition (simulating that ClientC died, so nobody's updated it's ownership in awhile)
			cps.ExpireOwnership(middleOwnership)

			lb := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientB), strategy, time.Hour)

			ownerships, err := lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3", "4"})
			require.NoError(t, err)
			require.NotEmpty(t, mapToStrings(ownerships, extractPartitionID))

			requireBalanced(t, cps, 5, 2)
		})
	}
}

func TestProcessorLoadBalancers_AnyStrategy_FullyBalancedOdd(t *testing.T) {
	for _, strategy := range []ProcessorStrategy{ProcessorStrategyBalanced, ProcessorStrategyGreedy} {
		t.Run(string(strategy), func(t *testing.T) {
			cps := newCheckpointStoreForTest()

			const clientA = "clientA"
			const clientB = "clientB"

			_, err := cps.ClaimOwnership(context.Background(), []Ownership{
				newTestOwnership("0", clientA),
				newTestOwnership("1", clientA),
				newTestOwnership("2", clientA),
				newTestOwnership("3", clientB),
				newTestOwnership("4", clientB),
			}, nil)
			require.NoError(t, err)

			{
				lbB := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientB), strategy, time.Hour)

				ownerships, err := lbB.LoadBalance(context.Background(), []string{"0", "1", "2", "3", "4"})
				require.NoError(t, err)
				require.Equal(t, []string{"3", "4"}, mapToStrings(ownerships, extractPartitionID))
				requireBalanced(t, cps, 5, 2)
			}

			{
				lbA := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientA), strategy, time.Hour)

				ownerships, err := lbA.LoadBalance(context.Background(), []string{"0", "1", "2", "3", "4"})
				require.NoError(t, err)
				require.Equal(t, []string{"0", "1", "2"}, mapToStrings(ownerships, extractPartitionID))
				requireBalanced(t, cps, 5, 2)
			}
		})
	}
}

func TestProcessorLoadBalancers_AnyStrategy_FullyBalancedEven(t *testing.T) {
	for _, strategy := range []ProcessorStrategy{ProcessorStrategyBalanced, ProcessorStrategyGreedy} {
		t.Run(string(strategy), func(t *testing.T) {
			cps := newCheckpointStoreForTest()

			const clientA = "clientA"
			const clientB = "clientB"

			_, err := cps.ClaimOwnership(context.Background(), []Ownership{
				newTestOwnership("0", clientA),
				newTestOwnership("1", clientA),
				newTestOwnership("2", clientB),
				newTestOwnership("3", clientB),
			}, nil)
			require.NoError(t, err)

			{
				lbB := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientB), strategy, time.Hour)

				ownerships, err := lbB.LoadBalance(context.Background(), []string{"0", "1", "2", "3"})
				require.NoError(t, err)
				require.Equal(t, []string{"2", "3"}, mapToStrings(ownerships, extractPartitionID))
				requireBalanced(t, cps, 4, 2)
			}

			{
				lbA := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientA), strategy, time.Hour)

				ownerships, err := lbA.LoadBalance(context.Background(), []string{"0", "1", "2", "3"})
				require.NoError(t, err)
				require.Equal(t, []string{"0", "1"}, mapToStrings(ownerships, extractPartitionID))
				requireBalanced(t, cps, 4, 2)
			}
		})
	}
}

func TestProcessorLoadBalancers_Any_GrabExtraPartitionBecauseAboveMax(t *testing.T) {
	for _, strategy := range []ProcessorStrategy{ProcessorStrategyBalanced, ProcessorStrategyGreedy} {
		t.Run(string(strategy), func(t *testing.T) {
			cps := newCheckpointStoreForTest()

			const clientA = "clientA"
			const clientB = "clientB"

			_, err := cps.ClaimOwnership(context.Background(), []Ownership{
				newTestOwnership("0", clientA),
				newTestOwnership("1", clientA),
				// nobody owns "2"
				newTestOwnership("3", clientB),
				newTestOwnership("4", clientB),
			}, nil)
			require.NoError(t, err)

			lb := newProcessorLoadBalancer(cps, newTestConsumerDetails(clientB), strategy, time.Hour)

			ownerships, err := lb.LoadBalance(context.Background(), []string{"0", "1", "2", "3", "4"})
			require.NoError(t, err)
			require.NotEmpty(t, mapToStrings(ownerships, extractPartitionID))

			requireBalanced(t, cps, 5, 2)
		})
	}
}

func TestProcessorLoadBalancers_AnyStrategy_StealsToBalance(t *testing.T) {
	for _, strategy := range []ProcessorStrategy{ProcessorStrategyBalanced, ProcessorStrategyGreedy} {
		t.Run(string(strategy), func(t *testing.T) {
			cps := newCheckpointStoreForTest()

			const lotsClientID = "has-too-many-client-id"
			const littleClientID = "has-too-few-id"
			allPartitionIDs := []string{"0", "1", "2", "3"}

			_, err := cps.ClaimOwnership(context.Background(), []Ownership{
				newTestOwnership(allPartitionIDs[0], lotsClientID),
				newTestOwnership(allPartitionIDs[1], lotsClientID),
				newTestOwnership(allPartitionIDs[2], lotsClientID),
				newTestOwnership(allPartitionIDs[3], littleClientID),
			}, nil)
			require.NoError(t, err)

			{
				tooManyLB := newProcessorLoadBalancer(cps, newTestConsumerDetails(lotsClientID), strategy, time.Hour)
				require.NoError(t, err)

				ownerships, err := tooManyLB.LoadBalance(context.Background(), allPartitionIDs)
				require.NoError(t, err)

				// it'll just keep reclaiming the partitions it owns (ie, nobody gives up partitions, they are
				// only taken)
				require.Equal(t, []string{"0", "1", "2"}, mapToStrings(ownerships, extractPartitionID))
			}

			{
				tooFewLB := newProcessorLoadBalancer(cps, newTestConsumerDetails(littleClientID), strategy, time.Hour)
				require.NoError(t, err)

				// either strategy will balance here by stealing.
				ownerships, err := tooFewLB.LoadBalance(context.Background(), allPartitionIDs)
				require.NoError(t, err)
				require.Equal(t, 2, len(ownerships))
			}

			requireBalanced(t, cps, len(allPartitionIDs), 2)
		})
	}
}

func TestProcessorLoadBalancers_InvalidStrategy(t *testing.T) {
	cps := newCheckpointStoreForTest()

	lb := newProcessorLoadBalancer(cps, newTestConsumerDetails("does not matter"), "", time.Hour)
	ownerships, err := lb.LoadBalance(context.Background(), []string{"0"})
	require.Nil(t, ownerships)
	require.EqualError(t, err, "invalid load balancing strategy ''")

	lb = newProcessorLoadBalancer(cps, newTestConsumerDetails("does not matter"), "super-greedy", time.Hour)
	ownerships, err = lb.LoadBalance(context.Background(), []string{"0"})
	require.Nil(t, ownerships)
	require.EqualError(t, err, "invalid load balancing strategy 'super-greedy'")
}

func mapToStrings[T any](src []T, fn func(t T) string) []string {
	var dest []string

	for _, t := range src {
		dest = append(dest, fn(t))
	}

	sort.Strings(dest)
	return dest
}

const testEventHubFQDN = "fqdn"
const testConsumerGroup = "consumer-group"
const testEventHubName = "event-hub"

func newTestOwnership(partitionID string, ownerID string) Ownership {
	return Ownership{
		OwnershipData: OwnershipData{
			OwnerID: ownerID,
		},
		CheckpointStoreAddress: CheckpointStoreAddress{
			PartitionID:             partitionID,
			ConsumerGroup:           testConsumerGroup,
			EventHubName:            testEventHubName,
			FullyQualifiedNamespace: testEventHubFQDN,
		},
	}
}

func newTestConsumerDetails(clientID string) consumerClientDetails {
	return consumerClientDetails{
		ConsumerGroup:           testConsumerGroup,
		EventHubName:            testEventHubName,
		FullyQualifiedNamespace: testEventHubFQDN,
		ClientID:                clientID,
	}
}

func requireBalanced(t *testing.T, cps CheckpointStore, totalPartitions int, numConsumers int) {
	ownerships, err := cps.ListOwnership(context.Background(), testEventHubFQDN, testEventHubName, testConsumerGroup, nil)
	require.NoError(t, err)

	min := totalPartitions / numConsumers
	max := min

	if totalPartitions%numConsumers > 0 {
		max++
	}

	require.Equal(t, len(ownerships), totalPartitions)

	ownershipMap := groupBy(ownerships, extractOwnerID)
	require.Equal(t, numConsumers, len(ownershipMap))

	for owner, partitions := range ownershipMap {
		require.Truef(t, len(partitions) == min || len(partitions) == max, "partitions for %s was %d, needed to be %d or %d", owner, len(partitions), min, max)
	}
}

func findCommon(left []string, right []string) []string {
	if len(left) == 0 || len(right) == 0 {
		panic("probable test error - left or right is empty")
	}

	leftMap := map[string]bool{}

	for _, l := range left {
		leftMap[l] = true
	}

	var common []string

	for _, r := range right {
		if leftMap[r] {
			common = append(common, r)
		}
	}

	return common
}

func extractPartitionID(t Ownership) string { return t.PartitionID }
func extractOwnerID(t Ownership) string     { return t.OwnerID }

func groupBy[T any](src []T, fn func(t T) string) map[string][]T {
	dest := map[string][]T{}

	for _, s := range src {
		key := fn(s)
		dest[key] = append(dest[key], s)
	}

	return dest
}
