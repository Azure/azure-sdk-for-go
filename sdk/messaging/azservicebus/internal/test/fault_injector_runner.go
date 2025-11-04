// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"slices"

	"github.com/stretchr/testify/require"
)

func StartSlowTransferFaultInjector(t *testing.T, host string, after int, delay time.Duration) string {
	return StartFaultInjector(t, host, []string{"transfer_delay",
		"--after", strconv.Itoa(after),
		"--delay", delay.String()})
}

func StartFaultInjector(t *testing.T, host string, args []string) string {
	faultInjectorPath := os.Getenv("AMQP_FAULT_INJECTOR_PATH")

	if faultInjectorPath == "" {
		t.Skip("AMQP_FAULT_INJECTOR_PATH not set to a path, skipping fault injector test")
	}

	tempDir := t.TempDir()

	addressFile := filepath.Join(tempDir, "fault-injector-address-file.txt")

	t.Logf("Address file: %s", addressFile)

	args = slices.Clone(args)
	args = append(args,
		"--cert", tempDir,
		"--logs", tempDir,
		"--address-file", addressFile,
		"--host", host)

	cmd := exec.Command(faultInjectorPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	require.NoError(t, err)

	t.Cleanup(func() {
		err := cmd.Process.Signal(os.Kill)

		if err != nil {
			t.Logf("Fault injector exited with error: %s", err)
		} else {
			t.Logf("Fault injector exited")
		}
	})

	// wait for our port file to show up
	for i := 0; i < 200; i++ {
		time.Sleep(200 * time.Millisecond)
		data, err := os.ReadFile(addressFile)

		if err != nil || len(data) < 4 {
			continue
		}

		t.Logf("Got %s from fault injector address file", string(data))
		return string(data)
	}

	require.Fail(t, "Fault injector didn't write an address file in a timely manner")
	return ""
}
