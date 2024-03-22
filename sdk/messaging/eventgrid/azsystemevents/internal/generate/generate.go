//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/internal/gopls"
)

var filesToDelete = []string{
	"options.go",
	"responses.go",
	"clientdeleteme_client.go",
}

const systemEventsGoFile = "system_events.go"
const modelsGoFile = "models.go"
const constantsGoFile = "constants.go"

const header = `//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

// Type represents the value set in EventData.EventType or messaging.CloudEvent.Type
// for system events.
type Type string

const (
`

const footer = `)`

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

type constant struct {
	ConstantName  string
	GoType        string
	ConstantValue string
}

func getConstantValues(reader io.ReadCloser) (map[string]constant, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	var currentComment []string

	var comments []struct {
		Text    string
		ForType string
	}

	for scanner.Scan() {
		line := scanner.Text()

		if commentText, found := cutPrefix(line, "// "); found {
			currentComment = append(currentComment, commentText)
		} else if len(currentComment) > 0 && strings.HasPrefix(line, "type ") {
			comments = append(comments, struct {
				Text    string
				ForType string
			}{
				Text:    strings.Join(currentComment, " "),
				ForType: strings.Split(line, " ")[1],
			})
			currentComment = nil
		}
	}

	consts := map[string]constant{}
	// ex: 'StorageBlobCreatedEventData - Schema of the Data property of an Event for a Microsoft.Storage.BlobCreated event.'
	//     'IotHubDeviceConnectedEventData - Event data for Microsoft.Devices.DeviceConnected event.'
	typeRE := regexp.MustCompile(`([^ ]+) -.+?([A-Za-z0-9]+(?:\.[A-Za-z0-9]+){2,3})`)

	ignoredTypes := []string{
		"DeviceConnectionStateEventProperties",
		"DeviceLifeCycleEventProperties",
		"DeviceTelemetryEventProperties",
		"EventGridMQTTClientEventData",
		"MapsGeofenceEventProperties",
	}

	for _, comment := range comments {
		matches := typeRE.FindStringSubmatch(comment.Text)

		if strings.Contains(comment.Text, "Schema of the Data") && len(matches) == 0 {
			ignorable := false

			// we have a few types that we can ignore - they're not top-level SystemEvents (they're
			// embedded or used by other events).
			for _, prefix := range ignoredTypes {
				if strings.HasPrefix(comment.Text, prefix) {
					ignorable = true
					break
				}
			}

			if !ignorable {
				log.Printf("===========> DIDN'T MATCH REGEX: %q ", comment)
			}
		}

		if len(matches) == 0 {
			continue
		}

		goTypeName, stringValue := matches[1], matches[2]
		constantName := strings.Replace(goTypeName, "EventData", "", 1)

		consts[stringValue] = constant{
			ConstantName:  "Type" + constantName,
			GoType:        goTypeName,
			ConstantValue: stringValue,
		}
	}

	return consts, nil
}

func writeConstantsFile(path string, constants map[string]constant) error {
	writer, err := os.Create(path)

	if err != nil {
		return err
	}

	defer writer.Close()

	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	for _, key := range sortedKeys(constants) {
		c := constants[key]

		// ex:
		// TypeAPIManagementAPICreated Type = "Microsoft.ApiManagement.APICreated" // maps to APIManagementAPICreatedEventData
		buff := fmt.Sprintf("%s Type = \"%s\" // maps to %s\n", c.ConstantName, c.ConstantValue, c.GoType)
		_, err := writer.WriteString(buff)

		if err != nil {
			return err
		}
	}

	if _, err := writer.Write([]byte(footer)); err != nil {
		return err
	}

	return nil
}

func sortedKeys[T any](m map[string]T) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// copy of strings.CutPrefix, which doesn't exist in our oldest support compiler (1.18)
func cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
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
