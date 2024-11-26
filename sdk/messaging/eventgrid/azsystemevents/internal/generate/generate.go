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

		if err := generateSystemEventEnum(); err != nil {
			return err
		}

		deleteUnneededFiles()
		return nil
	}

	if err := fn(); err != nil {
		fmt.Printf("./internal/generate: Failed with error: %s\n", err)
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

	// NOTE: the renaming of the error type is done in the propertyNameOverrideGo.tsp (AcsMessageChannelEventError)

	// NOTE: there appears to be a bug where my type name is automatically being capitalized, even though I marked it as internal.
	// Filed as https://github.com/Azure/autorest.go/issues/1467.
	if err := gopls.Rename(syms.Get("InternalACSMessageChannelEventError"), "internalACSMessageChannelEventError"); err != nil {
		return err
	}

	if err := gopls.Rename(syms.Get("InternalACSRouterCommunicationError"), "internalACSRouterCommunicationError"); err != nil {
		return err
	}

	// TODO: do I really need to handle these myself? Can I not use TypeSpec to do it?
	{
		if err := SwapType(syms.Get("AcsMessageReceivedEventData.Error"), "*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "ACSMessageReceivedEventData.Error", "unmarshalInternalACSMessageChannelEventError"); err != nil {
			return err
		}
	}

	{
		if err := SwapType(syms.Get("ACSMessageDeliveryStatusUpdatedEventData.Error"), "*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "ACSMessageDeliveryStatusUpdatedEventData.Error", "unmarshalInternalACSMessageChannelEventError"); err != nil {
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
			return fmt.Errorf("found error type which should have been deleted/renamed %q", sym.Name)
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
