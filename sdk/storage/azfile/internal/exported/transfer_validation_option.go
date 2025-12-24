// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"io"
)

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType interface {
	Apply(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)
	notPubliclyImplementable()
}

// TransferValidationTypeMD5 is a TransferValidationType used to provide a precomputed MD5.
type TransferValidationTypeMD5 []byte

func (c TransferValidationTypeMD5) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	cfg.SetMD5(c)
	return rsc, nil
}

func (TransferValidationTypeMD5) notPubliclyImplementable() {}
