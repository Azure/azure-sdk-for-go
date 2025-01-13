// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		err := godotenv.Load()

		if err != nil {
			log.Printf("Failed to load env file: %s", err.Error())
		}

		_ = test.MustGetEnvVars([]test.EnvKey{
			test.EnvKeyEndpoint,
			test.EnvKeyEndpointPremium,
			test.EnvKeyConnectionString,
			test.EnvKeyConnectionStringPremium,
			test.EnvKeyConnectionStringNoManage,
			test.EnvKeyConnectionStringSendOnly,
			test.EnvKeyConnectionStringListenOnly,
		})
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
