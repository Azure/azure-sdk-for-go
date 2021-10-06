// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type cosmosRequestOptions interface {
	toHeaders() *map[string]string
}
