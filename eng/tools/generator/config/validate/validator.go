// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
)

type Validator interface {
	Validate(cfg config.Config) error
}

func NewLocalValidator(specRoot string) Validator {
	return &localValidator{
		specRoot: specRoot,
	}
}

func NewRemoteValidator(ctx context.Context, client *query.Client) Validator {
	return &remoteValidator{
		ctx:    ctx,
		client: client,
	}
}
