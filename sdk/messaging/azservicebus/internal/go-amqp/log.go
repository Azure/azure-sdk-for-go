// Copyright (C) 2017 Kale Blankenship
// Portions Copyright (c) Microsoft Corporation

//go:build !debug
// +build !debug

package amqp

// dummy functions used when debugging is not enabled
// nolint:deadcode,unused
func debug(_ int, _ string, _ ...interface{}) {}
