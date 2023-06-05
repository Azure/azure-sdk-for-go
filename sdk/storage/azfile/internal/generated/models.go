//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

type TransactionalContentSetter interface {
	SetMD5([]byte)
	// add SetCRC64() when Azure File service starts supporting it.
}

func (f *FileClientUploadRangeOptions) SetMD5(v []byte) {
	f.ContentMD5 = v
}

type SourceContentSetter interface {
	SetSourceContentCRC64(v []byte)
	// add SetSourceContentMD5() when Azure File service starts supporting it.
}

func (f *FileClientUploadRangeFromURLOptions) SetSourceContentCRC64(v []byte) {
	f.SourceContentCRC64 = v
}
