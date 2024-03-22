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
	Orig     gopls.Symbol
	New      string
	FileName string
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

	phase1Repls = append(phase1Repls, rename{
		FileName: "models.go",
		Orig:     modelsSyms["AcsRouterCommunicationError.Innererror"],
		New:      "InnerError",
	})

	modelRepls := getModelsReplacements(modelsSyms)

	total := len(phase1Repls) + len(modelRepls)

	now := time.Now()

	for i, repl := range phase1Repls {
		fmt.Printf("%s[%d/%d]: %s -> %s\n", repl.FileName, i+1, total, repl.Orig.Name, repl.New)
		if err := gopls.Rename(repl.FileName, repl.Orig, repl.New); err != nil {
			return err
		}
	}

	for i, repl := range modelRepls {
		idx := 1 + i + len(phase1Repls)
		fmt.Printf("%s[%d/%d]: %s -> %s\n", repl.FileName, idx, total, repl.Orig.Name, repl.New)
		if err := gopls.Rename(repl.FileName, repl.Orig, repl.New); err != nil {
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

func getConstantsReplacements(syms map[string]gopls.Symbol) []rename {
	typeRE := regexp.MustCompile(`^(Acs|Avs|Iot)([A-Za-z0-9]+)$`)
	possibleRE := regexp.MustCompile(`^(Possible)(Acs|Avs|Iot)([A-Za-z0-9]+)$`)

	var renames []rename

	for k, sym := range syms {
		matches := typeRE.FindStringSubmatch(k)

		if matches != nil {
			renames = append(renames, rename{
				FileName: constantsGoFile,
				Orig:     sym,
				New:      strings.ToUpper(matches[1]) + matches[2],
			})
		}

		matches = possibleRE.FindStringSubmatch(k)

		if matches != nil {
			renames = append(renames, rename{
				FileName: constantsGoFile,
				Orig:     sym,
				New:      matches[1] + strings.ToUpper(matches[2]) + matches[3],
			})
		}
	}

	sort.Slice(renames, func(i, j int) bool {
		return renames[i].New < renames[j].New
	})
	return renames
}

func getModelsReplacements(syms map[string]gopls.Symbol) []rename {
	typeRE := regexp.MustCompile(`^(Acs|Avs|Iot)([A-Za-z0-9]+)$`)

	var renames []rename

	for k, sym := range syms {
		matches := typeRE.FindStringSubmatch(k)

		if matches != nil {
			renames = append(renames, rename{
				FileName: modelsGoFile,
				Orig:     sym,
				New:      strings.ToUpper(matches[1]) + matches[2],
			})
		}
	}

	sort.Slice(renames, func(i, j int) bool {
		return renames[i].Orig.Name < renames[j].Orig.Name
	})

	return renames
}
