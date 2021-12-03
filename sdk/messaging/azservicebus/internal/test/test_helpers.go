// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RandomString generates a random string with prefix
func RandomString(prefix string, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}

func GetConnectionString(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func GetConnectionStringWithoutManagePerms(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING_NO_MANAGE")

	if cs == "" {
		t.Skip()
	}

	return cs
}
