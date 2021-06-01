package validate

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/tools/generator/config"
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
