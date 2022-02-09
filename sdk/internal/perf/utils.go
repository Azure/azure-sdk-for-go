// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

// Helper function to sum a slice of integers
func sumInts(ints []int) int {
	var ret int
	for _, i := range ints {
		ret += i
	}
	return ret
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
