//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

func (a *AppendBlobClientAppendBlockOptions) SetCRC64(v []byte) {
	a.TransactionalContentCRC64 = v
}

func (a *AppendBlobClientAppendBlockOptions) SetMD5(v []byte) {
	a.TransactionalContentMD5 = v
}

func (a *BlockBlobClientStageBlockOptions) SetCRC64(v []byte) {
	a.TransactionalContentCRC64 = v
}

func (a *BlockBlobClientStageBlockOptions) SetMD5(v []byte) {
	a.TransactionalContentMD5 = v
}

func (a *PageBlobClientUploadPagesOptions) SetCRC64(v []byte) {
	a.TransactionalContentCRC64 = v
}

func (a *PageBlobClientUploadPagesOptions) SetMD5(v []byte) {
	a.TransactionalContentMD5 = v
}
