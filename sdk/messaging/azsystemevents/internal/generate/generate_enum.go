//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

var filesToDelete = []string{
	"options.go",
	"response_types.go",
	"clientdeleteme_client.go",
}

const systemEventsPath = "system_events.go"
const goModelsFile = "models.go"

const tsOutputPath = "./internal/generate/testdata/tsevents.txt"
const javaOutputPath = "./internal/generate/testdata/javaevents.txt"
const pyOutputPath = "./internal/generate/testdata/pyevents.txt"
const goOutputPath = "./internal/generate/testdata/goevents.txt"

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
	if len(os.Args) < 2 {
		fmt.Printf("Usage: generate_system_events (generate|validate|download)")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "download":
		err = writeEventFiles()
	case "generate":
		err = generateConstants()
	case "validate":
		var m map[string]bool
		m, err = readLines(pyOutputPath)

		if err != nil {
			break
		}

		err = validateAgainstAnotherLang(m)
	default:
		err = errors.New(os.Args[1] + " is not a valid command")
	}

	if err != nil {
		log.Fatalf("Failed with error: %s", err)
	}
}

func downloadEventTypes(url string, lineRE *regexp.Regexp, outputFile string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return writeEvents(resp.Body, lineRE, outputFile)
}

func writeEvents(reader io.Reader, lineRE *regexp.Regexp, outputPath string) error {
	scanner := bufio.NewScanner(reader)

	allMatches := map[string]bool{}

	for scanner.Scan() {
		matches := lineRE.FindStringSubmatch(scanner.Text())

		if len(matches) == 0 {
			continue
		}
		allMatches[matches[1]] = true
	}

	writer, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	defer writer.Close()

	for _, key := range sortedKeys(allMatches) {
		if _, err := writer.WriteString(fmt.Sprintf("%s\n", key)); err != nil {
			return err
		}
	}

	return nil
}

func writeEventFiles() error {
	err := downloadEventTypes("https://raw.githubusercontent.com/Azure/azure-sdk-for-js/main/sdk/eventgrid/eventgrid/src/predicates.ts", regexp.MustCompile(`^\s+"([^"]+)"`), tsOutputPath)

	if err != nil {
		return err
	}

	err = downloadEventTypes("https://raw.githubusercontent.com/Azure/azure-sdk-for-java/main/sdk/eventgrid/azure-messaging-eventgrid/src/main/java/com/azure/messaging/eventgrid/SystemEventNames.java", regexp.MustCompile(`"([^"]+?)"`), javaOutputPath)

	if err != nil {
		return err
	}

	err = downloadEventTypes("https://raw.githubusercontent.com/Azure/azure-sdk-for-python/main/sdk/eventgrid/azure-eventgrid/azure/eventgrid/_event_mappings.py", regexp.MustCompile(`^.+ = '([^']+)'`), pyOutputPath)

	if err != nil {
		return err
	}

	reader, err := os.Open(goModelsFile)

	if err != nil {
		return err
	}

	defer reader.Close()

	err = writeEvents(reader, regexp.MustCompile(`([A-Za-z0-9]+(?:\.[A-Za-z0-9]+){2,3})`), goOutputPath)

	if err != nil {
		return err
	}

	return nil
}

func generateConstants() error {
	reader, err := os.Open(goModelsFile)

	if err != nil {
		return fmt.Errorf("Failed to open %s: %w", goModelsFile, err)
	}

	defer reader.Close()

	constants, err := getConstantValues(reader)

	if err != nil {
		return fmt.Errorf("Failed to get constant values from file: %w", err)
	}

	if err := writeConstantsFile(systemEventsPath, constants); err != nil {
		log.Fatalf("Failed to write constants file %s: %s", systemEventsPath, err)
	}

	// we don't need these files since we're (intentionally) not exporting a Client from this
	// package.
	for _, file := range filesToDelete {
		log.Printf("Deleting %s since it only has client types", file)
		if _, err := os.Stat(file); err == nil {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateAgainstAnotherLang(otherLang map[string]bool) error {
	reader, err := os.Open(goModelsFile)

	if err != nil {
		return err
	}

	defer reader.Close()

	consts, err := getConstantValues(reader)

	if err != nil {
		return err
	}

	constsMap := map[string]constant{}

	for _, c := range consts {
		constsMap[c.ConstantValue] = c
	}

	var missingConstants []string

	// now cross-check our generated models against the models that we _should_ have
	// if we are caught up to the latest production Event Grid package for Java.
	for value := range otherLang {
		if _, ok := constsMap[value]; !ok {

			// there are some events that were deprecated
			if value == "Microsoft.Communication.ChatMemberAddedToThreadWithUser" || value == "Microsoft.Communication.ChatMemberRemovedFromThreadWithUser" {
				continue
			}

			// we're missing a value!
			missingConstants = append(missingConstants, value)
		}
	}

	if len(missingConstants) > 0 {
		sort.Strings(missingConstants)
		return errors.New("Missing constants:\n" + strings.Join(missingConstants, "\n"))
	}

	return nil
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

		if commentText, found := strings.CutPrefix(line, "// "); found {
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
			// embeded or used by other events).
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

func readLines(path string) (map[string]bool, error) {
	reader, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(reader)

	m := map[string]bool{}

	for scanner.Scan() {
		m[scanner.Text()] = true
	}

	return m, nil
}

/**
"ServiceBusDeadletterMessagesAvailableWithNoListenersEventData" => "ServiceBusDeadletterMessagesAvailableWithNoListener",
"SubscriptionDeletedEventData" => "EventGridSubscriptionDeleted",
"SubscriptionValidationEventData" => "EventGridSubscriptionValidation",
*/
