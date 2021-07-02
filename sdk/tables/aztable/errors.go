// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

func checkEntityForPkRk(entity *map[string]interface{}, err error) error {
	if _, ok := (*entity)[partitionKey]; !ok {
		return partitionKeyRowKeyError
	}

	if _, ok := (*entity)[rowKey]; !ok {
		return partitionKeyRowKeyError
	}

	return err
}
