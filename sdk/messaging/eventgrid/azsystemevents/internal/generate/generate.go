//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

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

		if err := doRenames(); err != nil {
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
		if err := gopls.Rename(syms.Get("AcsAdvancedMessageChannelEventError"), "internalACSAdvancedMessageChannelEventError"); err != nil {
			return err
		}

		if err := SwapType(syms.Get("AcsAdvancedMessageReceivedEventData.Error"), "*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "AcsAdvancedMessageReceivedEventData.Error", "unmarshalInternalACSAdvancedMessageChannelEventError"); err != nil {
			return err
		}
	}

	{
		if err := SwapType(syms.Get("AcsRouterJobClassificationFailedEventData.Errors"), "[]*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "AcsRouterJobClassificationFailedEventData.Errors", "unmarshalInternalACSRouterCommunicationError"); err != nil {
			return err
		}
	}

	{
		if err := gopls.Rename(syms.Get("AcsRouterCommunicationError"), "internalAcsRouterCommunicationError"); err != nil {
			return err
		}

		if err := SwapType(syms.Get("AcsRouterJobClassificationFailedEventData.Errors"), "[]*Error"); err != nil {
			return err
		}

		if err := UseCustomUnpopulate(modelsSerdeGoFile, "AcsRouterJobClassificationFailedEventData.Errors", "unmarshalInternalACSRouterCommunicationError"); err != nil {
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
		log.Fatalf("Failed to write constants file %s: %s", systemEventsGoFile, err)
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

type rename struct {
	Orig *gopls.Symbol
	New  string
}

func doRenames() error {
	constantSyms, err := gopls.Symbols(constantsGoFile)

	if err != nil {
		return err
	}

	modelsSyms, err := gopls.Symbols(modelsGoFile)

	if err != nil {
		return err
	}

	phase1Repls := getConstantsReplacements(constantSyms)

	// We don't need to do this since we're not holding onto this error type anymore.
	// phase1Repls = append(phase1Repls, rename{
	// 	Orig: modelsSyms.Get("AcsRouterCommunicationError.Innererror"),
	// 	New:  "InnerError",
	// })

	modelRepls := getModelsReplacements(modelsSyms)

	total := len(phase1Repls) + len(modelRepls)

	now := time.Now()

	for i, repl := range phase1Repls {
		fmt.Printf("%s[%d/%d]: %s -> %s\n", repl.Orig.File, i+1, total, repl.Orig.Name, repl.New)
		if err := gopls.Rename(repl.Orig, repl.New); err != nil {
			return err
		}
	}

	for i, repl := range modelRepls {
		idx := 1 + i + len(phase1Repls)
		fmt.Printf("%s[%d/%d]: %s -> %s\n", repl.Orig.File, idx, total, repl.Orig.Name, repl.New)
		if err := gopls.Rename(repl.Orig, repl.New); err != nil {
			return err
		}
	}

	fmt.Printf("Took %s to do gopls based renames\n", time.Since(now))

	return fixComments(modelRepls)
}

func fixComments(modelRepls []rename) error {
	// one annoyance here is that if the type name is not a doc reference (like [typename]) then
	// it doesn't get renamed, so we'll have to fix up the comments on our own.
	modelsSerdeBytes, err := os.ReadFile("models_serde.go")

	if err != nil {
		return err
	}

	constantsBytes, err := os.ReadFile(constantsGoFile)

	if err != nil {
		return err
	}

	for _, repl := range modelRepls {
		modelsSerdeBytes = bytes.Replace(modelsSerdeBytes, []byte(repl.Orig.Name), []byte(repl.New), -1)
		constantsBytes = bytes.Replace(constantsBytes, []byte(repl.Orig.Name), []byte(repl.New), -1)
	}

	if err := os.WriteFile("models_serde.go", modelsSerdeBytes, 0700); err != nil {
		return err
	}

	return os.WriteFile(constantsGoFile, constantsBytes, 0700)
}

func getConstantsReplacements(syms *gopls.SymbolMap) []rename {
	typeRE := regexp.MustCompile(`^(Acs|Avs|Iot)([A-Za-z0-9]+)$`)
	possibleRE := regexp.MustCompile(`^(Possible)(Acs|Avs|Iot)([A-Za-z0-9]+)$`)

	var renames []rename

	for _, sym := range syms.All() {
		matches := typeRE.FindStringSubmatch(sym.Name)

		if matches != nil {
			renames = append(renames, rename{
				Orig: sym,
				New:  strings.ToUpper(matches[1]) + matches[2],
			})
		}

		matches = possibleRE.FindStringSubmatch(sym.Name)

		if matches != nil {
			renames = append(renames, rename{
				Orig: sym,
				New:  matches[1] + strings.ToUpper(matches[2]) + matches[3],
			})
		}
	}

	sort.Slice(renames, func(i, j int) bool {
		return renames[i].New < renames[j].New
	})
	return renames
}

func getModelsReplacements(syms *gopls.SymbolMap) []rename {
	typeRE := regexp.MustCompile(`^(Acs|Avs|Iot)([A-Za-z0-9]+)$`)

	var renames []rename

	for _, sym := range syms.All() {
		matches := typeRE.FindStringSubmatch(sym.Name)

		if matches != nil {
			renames = append(renames, rename{
				Orig: sym,
				New:  strings.ToUpper(matches[1]) + matches[2],
			})
		}
	}

	sort.Slice(renames, func(i, j int) bool {
		return renames[i].Orig.Name < renames[j].Orig.Name
	})

	return renames
}
