//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
)

type testVars struct {
	OpenAIKey      string
	OpenAIEndpoint string

	AOAIKey      string
	AOAIEndpoint string
}

var tv testVars

const RecordingDirectory = "sdk/ai/azopenaiassistants/testdata"

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf(".env file couldn't load: %s", err)
	}

	tv.OpenAIKey = recording.GetEnvVariable("OPENAI_API_KEY", "key")
	tv.OpenAIEndpoint = recording.GetEnvVariable("OPENAI_ENDPOINT", "endpoint")

	tv.AOAIKey = recording.GetEnvVariable("AOAI_ASSISTANTS_KEY", "key")
	tv.AOAIEndpoint = recording.GetEnvVariable("AOAI_ASSISTANTS_ENDPOINT", "endpoint")

	os.Exit(run(m))
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(RecordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
		}
	}

	return m.Run()
}
