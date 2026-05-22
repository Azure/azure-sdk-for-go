//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

const systemEventsGoFile = "system_events.go"
const modelsGoFile = "models.go"
const modelsSerdeGoFile = "models_serde.go"

const header = `//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

const (
`

const footer = `)`

type constant struct {
	ConstantName  string
	GoType        string
	ConstantValue string
}

func getConstantValues(modelsGo string) (map[string]constant, error) {
	scanner := bufio.NewScanner(strings.NewReader(modelsGo))

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
				log.Fatalf("Non-system event type not classified: %q ", comment)
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

	defer func() {
		_ = writer.Close()
	}()

	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	for _, key := range sortedKeys(constants) {
		c := constants[key]

		// ex:
		// TypeAPIManagementAPICreated Type = "Microsoft.ApiManagement.APICreated" // maps to APIManagementAPICreatedEventData
		buff := fmt.Sprintf("%s = \"%s\" // maps to %s\n", c.ConstantName, c.ConstantValue, c.GoType)
		_, err := writer.WriteString(buff)

		if err != nil {
			return err
		}
	}

	if _, err := writer.Write([]byte(footer)); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
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
