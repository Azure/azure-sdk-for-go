// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"math"
	"math/rand"
	"sort"
)

func perWorkerSampleLimit(limit, workers, workerIndex int) int {
	if limit <= 0 || workers <= 1 {
		return limit
	}
	perWorker := limit / workers
	if workerIndex < limit%workers {
		perWorker++
	}
	if perWorker == 0 {
		return -1
	}
	return perWorker
}

// mergeReservoirs combines two uniform reservoir samples using weighted
// priority sampling. Each retained value represents seen/retained source
// observations, preventing worker merge order from biasing the result.
func mergeReservoirs[T any](dst []T, dstSeen int64, src []T, srcSeen int64, limit int, rng *rand.Rand) ([]T, int64) {
	totalSeen := dstSeen + srcSeen
	if limit <= 0 {
		return append(dst, src...), totalSeen
	}
	target := limit
	if available := len(dst) + len(src); available < target {
		target = available
	}
	type weightedValue struct {
		value T
		key   float64
	}
	candidates := make([]weightedValue, 0, len(dst)+len(src))
	appendCandidates := func(values []T, seen int64) {
		if len(values) == 0 {
			return
		}
		weight := float64(seen) / float64(len(values))
		for _, value := range values {
			// Exponential ranks implement weighted sampling without replacement;
			// lower keys have higher priority.
			key := -math.Log(1-rng.Float64()) / weight
			candidates = append(candidates, weightedValue{value: value, key: key})
		}
	}
	appendCandidates(dst, dstSeen)
	appendCandidates(src, srcSeen)
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].key < candidates[j].key })
	merged := make([]T, target)
	for i := range target {
		merged[i] = candidates[i].value
	}
	return merged, totalSeen
}
