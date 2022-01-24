// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import "fmt"

// Helper function to sum a slice of integers
func sumInts(ints []int) int {
	var ret int
	for _, i := range ints {
		ret += i
	}
	return ret
}

// convert an integer to a string, left padding with zeros
func leftPad(i int) string {
	if i >= 100 {
		return fmt.Sprintf("%d", i)
	} else if i >= 10 {
		return fmt.Sprintf("0%d", i)
	} else if i > 0 {
		return fmt.Sprintf("00%d", i)
	}
	return "000"
}

// Convert an integer to a string with comma separation in the format 1234 -> 1,234
func commaIze(i int) string {
	if i < 1000 {
		return fmt.Sprintf("%d", i)
	}

	copy := i
	ret := ""
	for copy >= 1000 {
		temp := copy % 1000
		tempS := leftPad(temp)
		ret = fmt.Sprintf("%s,%s", ret, tempS)

		copy /= 1000
	}

	if copy == 0 {
		return ret
	}

	return fmt.Sprintf("%d%s", copy, ret)
}
