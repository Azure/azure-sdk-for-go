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

		// apply temporary fixes
		if err := applyTempFixes(); err != nil {
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

// TODO: this is temporary while the PR is still being developed.
func applyTempFixes() error {
	// See: https://github.com/Azure/azure-rest-api-specs/pull/28706/files#r1596180693
	{
		err := Replace(polymorphicHelpersGoFile, "@odata.type polymorphic fix", func(data string) (string, error) {
			return strings.Replace(string(data), `switch m["@odataType"] {`, `	switch m["@odata.type"] {`, 1), nil
		})

		if err != nil {
			return err
		}

		err = Replace(modelsSerdeGoFile, "@odata.type populate", func(data string) (string, error) {
			return strings.Replace(string(data), `populate(objectMap, "@odataType", m.ODataType)`, `populate(objectMap, "@odata.type", m.ODataType)`, 1), nil
		})

		if err != nil {
			return err
		}

		err = Replace(modelsSerdeGoFile, "@odata.type deserialize", func(data string) (string, error) {
			return strings.Replace(string(data), `case "@odataType":`, `case "@odata.type":`, 1), nil
		})

		if err != nil {
			return err
		}
	}

	// See: https://github.com/Azure/azure-rest-api-specs/pull/28706/files#r1597270162
	fixups := []string{
		"Microsoft.EventGrid.SystemEvents.MQTTClientCreatedOrUpdated",
		"Microsoft.EventGrid.SystemEvents.MQTTClientDeleted",
		"Microsoft.EventGrid.SystemEvents.MQTTClientSessionConnected",
		"Microsoft.EventGrid.SystemEvents.MQTTClientSessionDisconnected",
		"Microsoft.EventGrid.SystemEvents.SubscriptionDeletedEvent",
		"Microsoft.EventGrid.SystemEvents.SubscriptionValidationEvent",
	}

	for _, fixup := range fixups {
		// there's some badly named constants for the system events.
		err := Replace(modelsGoFile, fmt.Sprintf("fixing ID for %s", fixup), func(data string) (string, error) {
			return strings.Replace(string(data), fixup, strings.Replace(fixup, ".SystemEvents.", ".", 1), 1), nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func doRemove() error {
	modelsToRemove := []string{
		// these are not events, themselves, but are just contained within
		// other events.
		"AVSClusterEventData",
		"AVSPrivateCloudEventData",
		"AVSScriptExecutionEventData",
	}

	for _, m := range modelsToRemove {
		if err := DeleteType(m); err != nil {
			return err
		}
	}

	return nil
}
