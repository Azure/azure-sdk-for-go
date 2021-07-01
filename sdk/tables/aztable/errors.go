// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import "errors"

func checkEntityForPkRk(entity *map[string]interface{}, err error) error {

	if _, ok := (*entity)["PartitionKey"]; !ok {
		return errors.New("Entity must have a PartitionKey and RowKey")
	}

	if _, ok := (*entity)["RowKey"]; !ok {
		return errors.New("Entity must have a PartitionKey and RowKey")
	}

	return err
}
