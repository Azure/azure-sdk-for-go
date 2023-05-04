//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueuePermissions_String(t *testing.T) {
	testdata := []struct {
		input    QueuePermissions
		expected string
	}{
		{input: QueuePermissions{Read: true}, expected: "r"},
		{input: QueuePermissions{Add: true}, expected: "a"},
		{input: QueuePermissions{Update: true}, expected: "u"},
		{input: QueuePermissions{Process: true}, expected: "p"},
		{input: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, expected: "raup"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestQueuePermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected QueuePermissions
	}{
		{expected: QueuePermissions{Read: true}, input: "r"},
		{expected: QueuePermissions{Add: true}, input: "a"},
		{expected: QueuePermissions{Update: true}, input: "u"},
		{expected: QueuePermissions{Process: true}, input: "p"},
		{expected: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, input: "raup"},
		{expected: QueuePermissions{
			Read:    true,
			Add:     true,
			Update:  true,
			Process: true,
		}, input: "puar"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseQueuePermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestQueuePermissions_ParseNegative(t *testing.T) {
	_, err := parseQueuePermissions("puatr") // Here 't' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "116")
}

func TestGetCanonicalName(t *testing.T) {
	testdata := []struct {
		inputAccount string
		inputQueue   string
		expected     string
	}{
		{inputAccount: "fakestorageaccount", inputQueue: "fakestoragequeue", expected: "/queue/fakestorageaccount/fakestoragequeue"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getCanonicalName(c.inputAccount, c.inputQueue))
	}
}
