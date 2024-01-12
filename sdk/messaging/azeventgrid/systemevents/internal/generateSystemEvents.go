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
)

func main() {
	if err := impl(); err != nil {
		log.Fatalf("Failed to generate system events: %s", err)
	}
}

const clientFile = "clientdeleteme_client.go"

const header = `//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package systemevents

// Type represents the value set in EventData.EventType or messaging.CloudEvent.Type
// for system events.
type Type string

const (
`

const footer = `)`

func impl() error {
	if _, err := os.Stat(clientFile); err == nil {
		if err := os.Remove(clientFile); err != nil {
			return err
		}
	}

	reader, err := os.Open("models.go")

	if err != nil {
		return err
	}

	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	writer, err := os.Create("system_events.go")

	if err != nil {
		return err
	}

	defer writer.Close()

	if _, err := writer.Write([]byte(header)); err != nil {
		return err
	}

	// ex: '// APIManagementAPICreatedEventData - Schema of the Data property of an Event for a Microsoft.ApiManagement.APICreated'
	typeRE := regexp.MustCompile(`^// ([^ ]+) - Schema of the Data property of an Event for a (.+?)$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := typeRE.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		goTypeName, stringType := matches[1], matches[2]

		buff := []byte(fmt.Sprintf("Type%s Type = \"%s\" // maps to %s\n", goTypeName, stringType, goTypeName))

		if _, err := writer.Write(buff); err != nil {
			return err
		}
	}

	if _, err := writer.Write([]byte(footer)); err != nil {
		return err
	}

	return nil
}
