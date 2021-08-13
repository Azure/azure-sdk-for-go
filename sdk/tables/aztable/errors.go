// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

func checkEntityForPkRk(entity *map[string]interface{}, err error) error {
	if _, ok := (*entity)[partitionKey]; !ok {
		return errPartitionKeyRowKeyError
	}

	if _, ok := (*entity)[rowKey]; !ok {
		return errPartitionKeyRowKeyError
	}

	return err
}
