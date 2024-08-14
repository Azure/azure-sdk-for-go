//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/internal/gopls"
)

var filesToDelete = []string{
	"options.go",
	"responses.go",
	"clientdeleteme_client.go",
}

func main() {
	fn := func() error {
		if err := swapErrorTypes(); err != nil {
			return err
		}

		// remove extraneous types
		if err := doRemove(); err != nil {
			return err
		}

		if err := generateSystemEventEnum(); err != nil {
			return err
		}

		deleteUnneededFiles()
		return nil
	}

	if err := fn(); err != nil {
		fmt.Printf("Failed with error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("DONE\n")
}

// swapErrorTypes handles turning most of the auto-generated errors into a single consistent error type.
// The key is that the Error type doesn't export human readable strings as fields - it's all contained in
// the Error() field.
func swapErrorTypes() error {
	syms, err := gopls.Symbols(modelsGoFile)

	if err != nil {
		return err
	}

	{
		if err := SwapType(syms.Get("AcsMessageReceivedEventData.Error"), "*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "ACSMessageReceivedEventData.Error", "unmarshalInternalACSMessageChannelEventError"); err != nil {
			return err
		}
	}

	{
		if err := SwapType(syms.Get("AcsRouterJobClassificationFailedEventData.Errors"), "[]*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "ACSRouterJobClassificationFailedEventData.Errors", "unmarshalInternalACSRouterCommunicationError"); err != nil {
			return err
		}
	}

	syms, err = gopls.Symbols(modelsGoFile)

	if err != nil {
		return err
	}

	allowedErrs := map[string]bool{
		"MediaJobError": true,
	}

	for _, sym := range syms.All() {
		if allowedErrs[sym.Name] {
			continue
		}

		if strings.HasSuffix(sym.Name, "Error") && !strings.HasPrefix(sym.Name, "internal") && sym.Type == "Struct" {
			return fmt.Errorf("found redundant unhandled error type %s\n", sym.Name)
		}
	}

	return nil
}

func generateSystemEventEnum() error {
	reader, err := os.Open(modelsGoFile)

	if err != nil {
		return fmt.Errorf("Failed to open %s: %w", modelsGoFile, err)
	}

	defer reader.Close()

	constants, err := getConstantValues(reader)

	if err != nil {
		return fmt.Errorf("Failed to get constant values from file: %w", err)
	}

	if err := writeConstantsFile(systemEventsGoFile, constants); err != nil {
		return fmt.Errorf("Failed to write constants file %s: %w", systemEventsGoFile, err)
	}

	return nil
}

func deleteUnneededFiles() {
	// we don't need these files since we're (intentionally) not exporting a Client from this
	// package.
	fmt.Printf("Deleting unneeded files\n")

	for _, file := range filesToDelete {
		fmt.Printf("Deleting %s since it only has client types\n", file)
		if _, err := os.Stat(file); err == nil {
			_ = os.Remove(file)
		}
	}
}

func doRemove() error {
	modelsToRemove := []string{
		// These types are base objects of some of our system events in the TypeSpec, giving them a simple way to share fields.
		// Our generator handles this parent/child relationship by just inlining those properties into the children, so the base struct is just vestigial.
		// Note that these have been annotated with @output, which is why they're not just clipped out using our normal "unused/unreferenced" type logic
		// in the Go emitter.
		"ACSChatEventBaseProperties",
		"ACSChatEventInThreadBaseProperties",
		"ACSChatMessageEventBaseProperties",
		"ACSChatMessageEventInThreadBaseProperties",
		"ACSChatThreadEventBaseProperties",
		"ACSChatThreadEventInThreadBaseProperties",
		"ACSRouterEventData",
		"ACSRouterJobEventData",
		"ACSRouterWorkerEventData",
		"ACSSmsEventBaseProperties",
		"AppConfigurationSnapshotEventData",
		"AVSClusterEventData",
		"AVSPrivateCloudEventData",
		"AVSScriptExecutionEventData",
		"ContainerServiceClusterSupportEventData",
		"ContainerServiceNodePoolRollingEventData",
		"ResourceNotificationsResourceDeletedEventData",
		"ResourceNotificationsResourceUpdatedEventData",
	}

	for _, m := range modelsToRemove {
		if err := DeleteType(m); err != nil {
			return err
		}
	}

	return nil
}
