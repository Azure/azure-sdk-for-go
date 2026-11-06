// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !go_v1 && !go_v2

package main

import (
	"context"
	"errors"
)

func newSDKClient(context.Context, options) (sdkClient, error) {
	return nil, errors.New("build with exactly one of -tags go_v1 or -tags go_v2")
}

func isNotFound(_ error) bool {
	return false
}
