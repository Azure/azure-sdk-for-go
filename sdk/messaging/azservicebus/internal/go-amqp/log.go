//go:build !debug
// +build !debug

// Copyright (C) 2017 Kale Blankenship
// Portions Copyright (c) Microsoft Corporation
package amqp

// dummy functions used when debugging is not enabled

func debug(_ int, _ string, _ ...interface{}) {}
