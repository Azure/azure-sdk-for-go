// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type processorLoadBalancer struct {
	checkpointStore             CheckpointStore
	details                     consumerClientDetails
	strategy                    ProcessorStrategy
	partitionExpirationDuration time.Duration

	// NOTE: when you create your own *rand.Rand it is not thread safe.
	rnd *rand.Rand
}

func newProcessorLoadBalancer(checkpointStore CheckpointStore, details consumerClientDetails, strategy ProcessorStrategy, partitionExpiration time.Duration) *processorLoadBalancer {
	return &processorLoadBalancer{
		checkpointStore:             checkpointStore,
		details:                     details,
		strategy:                    strategy,
		partitionExpirationDuration: partitionExpiration,
		rnd:                         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type loadBalancerInfo struct {
	// current are the partitions that _we_ own
	current []Ownership

	// unownedOrExpired partitions either had no claim _ever_ or were once
	// owned but the ownership claim has expired.
	unownedOrExpired []Ownership

	// aboveMax are ownerships where the specific owner has too many partitions
	// it contains _all_ the partitions for that particular consumer.
	aboveMax []Ownership

	// maxAllowed is the maximum number of partitions a consumer should have
	// If partitions do not divide evenly this will be the "theoretical" max
	// with the assumption that this particular consumer will get an extra
	// partition.
	maxAllowed int

	// extraPartitionPossible is true if the partitions cannot split up evenly
	// amongst all the known consumers.
	extraPartitionPossible bool

	raw []Ownership
}

// loadBalance calls through to the user's configured load balancing algorithm.
// NOTE: this function is NOT thread safe!
func (lb *processorLoadBalancer) LoadBalance(ctx context.Context, partitionIDs []string) ([]Ownership, error) {
	lbinfo, err := lb.getAvailablePartitions(ctx, partitionIDs)

	if err != nil {
		return nil, err
	}

	claimMorePartitions := true

	if len(lbinfo.current) >= lbinfo.maxAllowed {
		// - I have _exactly_ the right amount
		// or
		// - I have too many. We expect to have some stolen from us, but we'll maintain
		//    ownership for now.
		claimMorePartitions = false
	} else if lbinfo.extraPartitionPossible && len(lbinfo.current) == lbinfo.maxAllowed-1 {
		// In the 'extraPartitionPossible' scenario, some consumers will have an extra partition
		// since things don't divide up evenly. We're one under the max, which means we _might_
		// be able to claim another one.
		//
		// We will attempt to grab _one_ more but only if there are free partitions available
		// or if one of the consumers has more than the max allowed.
		claimMorePartitions = len(lbinfo.unownedOrExpired) > 0 || len(lbinfo.aboveMax) > 0
	}

	ownerships := lbinfo.current

	if claimMorePartitions {
		switch lb.strategy {
		case ProcessorStrategyGreedy:
			ownerships = lb.greedyLoadBalancer(ctx, lbinfo)
		case ProcessorStrategyBalanced:
			o := lb.balancedLoadBalancer(ctx, lbinfo)

			if o != nil {
				ownerships = append(lbinfo.current, *o)
			}
		default:
			return nil, fmt.Errorf("invalid load balancing strategy '%s'", lb.strategy)
		}
	}

	return lb.checkpointStore.ClaimOwnership(ctx, ownerships, nil)
}

// getAvailablePartitions finds all partitions that are either completely unowned _or_
// their ownership is stale.
func (lb *processorLoadBalancer) getAvailablePartitions(ctx context.Context, partitionIDs []string) (loadBalancerInfo, error) {
	ownerships, err := lb.checkpointStore.ListOwnership(ctx, lb.details.FullyQualifiedNamespace, lb.details.EventHubName, lb.details.ConsumerGroup, nil)

	if err != nil {
		return loadBalancerInfo{}, err
	}

	alreadyAdded := map[string]bool{}
	groupedByOwner := map[string][]Ownership{
		lb.details.ClientID: nil,
	}

	var unownedOrExpired []Ownership

	// split out partitions by whether they're currently owned
	// and if they're expired.
	for _, o := range ownerships {
		alreadyAdded[o.PartitionID] = true

		if time.Since(o.LastModifiedTime.UTC()) > lb.partitionExpirationDuration {
			unownedOrExpired = append(unownedOrExpired, o)
			continue
		}

		groupedByOwner[o.OwnerID] = append(groupedByOwner[o.OwnerID], o)
	}

	// add in all the unowned partitions
	for _, partID := range partitionIDs {
		if alreadyAdded[partID] {
			continue
		}

		unownedOrExpired = append(unownedOrExpired, Ownership{
			FullyQualifiedNamespace: lb.details.FullyQualifiedNamespace,
			ConsumerGroup:           lb.details.ConsumerGroup,
			EventHubName:            lb.details.EventHubName,
			PartitionID:             partID,
			OwnerID:                 lb.details.ClientID,
			// note that we don't have etag info here since nobody has
			// ever owned this partition.
		})
	}

	maxAllowed := len(partitionIDs) / len(groupedByOwner)
	hasRemainder := len(partitionIDs)%len(groupedByOwner) > 0

	if hasRemainder {
		maxAllowed += 1
	}

	var aboveMax []Ownership

	for id, ownerships := range groupedByOwner {
		if id == lb.details.ClientID {
			continue
		}

		if len(ownerships) > maxAllowed {
			aboveMax = append(aboveMax, ownerships...)
		}
	}

	return loadBalancerInfo{
		current:                groupedByOwner[lb.details.ClientID],
		unownedOrExpired:       unownedOrExpired,
		aboveMax:               aboveMax,
		maxAllowed:             maxAllowed,
		extraPartitionPossible: hasRemainder,
		raw:                    ownerships,
	}, nil
}

// greedyLoadBalancer will attempt to grab as many free partitions as it needs to balance
// in each round.
func (lb *processorLoadBalancer) greedyLoadBalancer(ctx context.Context, lbinfo loadBalancerInfo) []Ownership {
	ours := lbinfo.current

	// try claiming from the completely unowned or expires ownerships _first_
	randomOwnerships := getRandomOwnerships(lb.rnd, lbinfo.unownedOrExpired, lbinfo.maxAllowed-len(ours))
	ours = append(ours, randomOwnerships...)

	if len(ours) < lbinfo.maxAllowed {
		// if that's not enough then we'll randomly steal from any owners that had partitions
		// above the maximum.
		randomOwnerships := getRandomOwnerships(lb.rnd, lbinfo.aboveMax, lbinfo.maxAllowed-len(ours))
		ours = append(ours, randomOwnerships...)
	}

	for i := 0; i < len(ours); i++ {
		ours[i] = lb.resetOwnership(ours[i])
	}
	return ours
}

// balancedLoadBalancer attempts to split the partition load out between the available
// consumers so each one has an even amount (or even + 1, if the # of consumers and #
// of partitions doesn't divide evenly).
//
// NOTE: the checkpoint store itself does not have a concept of 'presence' that doesn't
// ALSO involve owning a partition. It's possible for a consumer to get boxed out for a
// bit until it manages to steal at least one partition since the other consumers don't
// know it exists until then.
func (lb *processorLoadBalancer) balancedLoadBalancer(ctx context.Context, lbinfo loadBalancerInfo) *Ownership {
	if len(lbinfo.unownedOrExpired) > 0 {
		idx := lb.rnd.Intn(len(lbinfo.unownedOrExpired))
		o := lb.resetOwnership(lbinfo.unownedOrExpired[idx])
		return &o
	}

	if len(lbinfo.aboveMax) > 0 {
		idx := lb.rnd.Intn(len(lbinfo.aboveMax))
		o := lb.resetOwnership(lbinfo.aboveMax[idx])
		return &o
	}

	return nil
}

func (lb *processorLoadBalancer) resetOwnership(o Ownership) Ownership {
	o.ETag = nil
	o.OwnerID = lb.details.ClientID
	return o
}

func getRandomOwnerships(rnd *rand.Rand, ownerships []Ownership, count int) []Ownership {
	limit := int(math.Min(float64(count), float64(len(ownerships))))

	if limit == 0 {
		return nil
	}

	choices := rnd.Perm(limit)

	var newOwnerships []Ownership

	for i := 0; i < len(choices); i++ {
		newOwnerships = append(newOwnerships, ownerships[choices[i]])
	}

	return newOwnerships
}
