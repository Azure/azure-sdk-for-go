// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type emulatorTests struct {
	host string
	key  string
}

func newEmulatorTests() *emulatorTests {
	return &emulatorTests{
		host: "https://localhost:8081/",
		key:  "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==",
	}
}
