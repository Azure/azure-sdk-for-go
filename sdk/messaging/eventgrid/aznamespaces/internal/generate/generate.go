// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces/internal/generate/gopls"
)

func main() {
	if err := checkRequestMethods(); err != nil {
		panic(err)
	}

	if err := deleteModels(); err != nil {
		panic(err)
	}

	if err := injectClientData(); err != nil {
		panic(err)
	}

	if err := exportPublicSymbols(); err != nil {
		panic(err)
	}

	if err := useAZCoreCloudEvent(); err != nil {
		panic(err)
	}

	if err := transformError(); err != nil {
		panic(err)
	}

	// Do a double check, we shouldn't have any other unexported symbols
	// in generated code. If we do it's probably fallout from our unexporting
	// in TypeSpec.
	if err := checkUnexportedSymbols(); err != nil {
		panic(err)
	}
}

func useAZCoreCloudEvent() error {
	if err := replaceAllInFile("models.go", regexp.MustCompile(`(?m)^\s+Event \*cloudEvent$`), []byte("\tEvent messaging.CloudEvent")); err != nil {
		return err
	}

	if err := replaceAllInFile("sender_client.go", regexp.MustCompile(`(?m)event cloudEvent`), []byte("\tevent *messaging.CloudEvent")); err != nil {
		return err
	}

	if err := replaceAllInFile("sender_client.go", regexp.MustCompile(`(?m)events \[\]cloudEvent`), []byte("\tevents []*messaging.CloudEvent")); err != nil {
		return err
	}

	return nil
}

func deleteModels() error {
	models := "cloudEvent|innerError"

	if err := replaceAllInFile(
		"models_serde.go",
		regexp.MustCompile(
			fmt.Sprintf(`(?s)// (Unmarshal|Marshal)JSON implements the json\.(Unmarshaller|Marshaller) interface for type (%s).+?\n}`, models)), nil); err != nil {
		return err
	}

	if err := replaceAllInFile(
		"models.go",
		regexp.MustCompile(fmt.Sprintf(`(?s)// (%s).+?type (%s) struct.+?\n}`, models, models)), nil); err != nil {
		return err
	}

	if err := replaceAllInFile(
		"options.go",
		regexp.MustCompile(`(?s)// (senderClientsendOptions).+?type (senderClientsendOptions) struct.+?\n}`), nil); err != nil {
		return err
	}

	if err := replaceAllInFile("sender_client.go", regexp.MustCompile(`senderClientsendOptions`), []byte("SendOptions")); err != nil {
		return err
	}

	return nil
}

// injectClientData adds our senderData and receiverData structs as
// fields for their respective clients. We stash some data in there
// for our client methods to use later.
func injectClientData() error {
	data, err := os.ReadFile("sender_client.go")

	if err != nil {
		return err
	}

	newData := bytes.Replace(data, []byte("type SenderClient struct {"), []byte("type SenderClient struct {\ndata senderData\n"), 1)

	if bytes.Equal(data, newData) {
		return fmt.Errorf("failed applying senderData patch")
	}

	if err := os.WriteFile("sender_client.go", newData, 0600); err != nil {
		return err
	}

	data, err = os.ReadFile("receiver_client.go")

	if err != nil {
		return err
	}

	newData = bytes.Replace(data, []byte("type ReceiverClient struct {"), []byte("type ReceiverClient struct {\ndata receiverData\n"), 1)

	if bytes.Equal(data, newData) {
		return fmt.Errorf("failed applying receiverData patch")
	}

	if err := os.WriteFile("receiver_client.go", newData, 0600); err != nil {
		return err
	}

	return nil
}

func checkRequestMethods() error {
	typesToFlatten := map[string]bool{
		// these were customized to already be unexported within the TypeSpec
		// by unexporting the individual client methods.
		"acknowledgeOptions": true,
		"rejectOptions":      true,
		"releaseOptions":     true,
	}

	// these are actually models, not options
	modelSymbols, err := gopls.SymbolsSlice("models.go")

	if err != nil {
		return err
	}

	for _, sym := range modelSymbols {
		if sym.Type != gopls.SymbolTypeStruct || !typesToFlatten[sym.Name] {
			continue
		}

		if len(sym.Children) != 1 || sym.Children[0].Name != "LockTokens" {
			return fmt.Errorf("Type %s contains more than just LockTokens (%#v), check to see if the Receiver or Sender signatures need to be adjusted", sym.Name, sym.Children)
		}
	}

	return nil
}

// we specifically override the visiblity of the normal service operations because we
// have to do mods like removing the topic and subscription parameters, or making some
// parameters positional instead of being in a struct argument.
func exportPublicSymbols() error {
	err := multiRename("responses.go", map[string]string{
		// receiver
		"receiverClientacknowledgeResponse": "AcknowledgeResponse",
		"receiverClientreceiveResponse":     "ReceiveResponse",
		"receiverClientrejectResponse":      "RejectResponse",
		"receiverClientreleaseResponse":     "ReleaseResponse",
		"receiverClientrenewLockResponse":   "RenewLocksResponse",

		// sender
		"senderClientsendEventsResponse": "SendEventsResponse",
		"senderClientsendResponse":       "SendResponse",
	})

	if err != nil {
		return err
	}

	err = multiRename("options.go", map[string]string{
		// receiver
		"receiverClientacknowledgeOptions": "AcknowledgeOptions",
		"receiverClientreceiveOptions":     "ReceiveOptions",
		"receiverClientrejectOptions":      "RejectOptions",
		"receiverClientreleaseOptions":     "ReleaseOptions",
		"receiverClientrenewLockOptions":   "RenewLocksOptions",

		// sender
		"senderClientsendEventsOptions": "SendEventsOptions",
		"senderClientsendOptions":       "SendOptions",
	})

	if err != nil {
		return err
	}

	err = multiRename("models.go", map[string]string{
		// receiver
		"acknowledgeResult": "AcknowledgeResult",
		"receiveResult":     "ReceiveResult",
		"rejectResult":      "RejectResult",
		"releaseResult":     "ReleaseResult",
		"renewLocksResult":  "RenewLocksResult",

		"failedLockToken": "FailedLockToken",
		"receiveDetails":  "ReceiveDetails",

		// sender
		"publishResult": "PublishResult",

		// other models that are embedded in results
		"brokerProperties": "BrokerProperties",
	})

	if err != nil {
		return err
	}

	err = multiRename("constants.go", map[string]string{
		"PossiblereleaseDelayValues": "PossibleReleaseDelayValues",
		"releaseDelay":               "ReleaseDelay",
		"releaseDelayNoDelay":        "ReleaseDelayNoDelay",
		"releaseDelayOneHour":        "ReleaseDelayOneHour",
		"releaseDelayOneMinute":      "ReleaseDelayOneMinute",
		"releaseDelayTenMinutes":     "ReleaseDelayTenMinutes",
		"releaseDelayTenSeconds":     "ReleaseDelayTenSeconds",
	})

	if err != nil {
		return err
	}

	return replaceAllInFile("constants.go", regexp.MustCompile("releaseDelay const type"), []byte("ReleaseDelay const type"))
}

func checkUnexportedSymbols() error {
	var unexported []string

	whiteListed := map[string]bool{
		"rejectOptions":      true,
		"releaseOptions":     true,
		"publishResult":      true,
		"renewLockOptions":   true,
		"acknowledgeOptions": true,
	}

	for _, file := range []string{"models.go", "options.go", "responses.go"} {
		symbols, err := gopls.SymbolsSlice(file)

		if err != nil {
			return err
		}

		for _, sym := range symbols {
			if sym.Type != gopls.SymbolTypeStruct || whiteListed[sym.Name] {
				continue
			}

			rn, _ := utf8.DecodeRuneInString(sym.Name)

			if unicode.IsLower(rn) {
				unexported = append(unexported, fmt.Sprintf("[%s] %s", file, sym.Name))
			}
		}
	}

	sort.Strings(unexported)

	if len(unexported) > 0 {
		return fmt.Errorf("symbols are unexported:\n%s", strings.Join(unexported, "\n"))
	}

	return nil
}

func multiRename(file string, renames map[string]string) error {
	symbols, err := gopls.SymbolsSlice(file)

	if err != nil {
		return err
	}

	// we do it in reverse order so we can do all the renames without reloading the file
	// since our positions should not shift.
	for i := len(symbols) - 1; i >= 0; i-- {
		if renames[symbols[i].Name] == "" {
			continue
		}

		if err := gopls.Rename(symbols[i], renames[symbols[i].Name]); err != nil {
			return err
		}
	}

	return nil
}

func replaceAllInFile(file string, re *regexp.Regexp, replacement []byte) error {
	oldBytes, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	newBytes := re.ReplaceAll(oldBytes, replacement)

	if bytes.Equal(oldBytes, newBytes) {
		return fmt.Errorf("no changes in %s", file)
	}

	if err := os.WriteFile(file, newBytes, 0600); err != nil {
		return err
	}

	return nil
}

func transformError() error {
	symbols, err := gopls.SymbolsSlice("models.go")

	if err != nil {
		return err
	}

	renamed := false

	for _, sym := range symbols {
		if sym.Name == "errorModel" {
			if err := gopls.Rename(sym, "Error"); err != nil {
				return err
			}
			renamed = true
		}
	}

	if !renamed {
		return errors.New("errorModel wasn't present in models.go")
	}

	// clip out some of the internal fields
	err = replaceAllInFile("models.go", regexp.MustCompile(`(?s)// [^/]+?\s+(Innererror \*innerError|Details \[\]Error|Target \*string)\n`), nil)

	if err != nil {
		return err
	}

	// delete these lines from the marshallers and ummarshallers
	/*
		populate(objectMap, "details", e.Details)
		populate(objectMap, "innererror", e.Innererror)
		populate(objectMap, "target", e.Target)
	*/

	err = replaceAllInFile("models_serde.go", regexp.MustCompile(`(?m)^\s+populate\(objectMap, "(details|innererror|target)", e\.(Details|Innererror|Target)\)$`), nil)

	if err != nil {
		return fmt.Errorf("removing Error populate lines: %w", err)
	}

	/*
		case "details":
			err = unpopulate(val, "Details", &e.Details)
			delete(rawMsg, key)
		case "innererror":
			err = unpopulate(val, "Innererror", &e.Innererror)
			delete(rawMsg, key)
		case "target":
			err = unpopulate(val, "Target", &e.Target)
			delete(rawMsg, key)
		}
	*/

	err = replaceAllInFile("models_serde.go", regexp.MustCompile(`(?s)case "(details|innererror|target)":.+?delete\(rawMsg, key\)\n`), nil)

	if err != nil {
		return fmt.Errorf("removing Error unpopulate lines: %w", err)
	}

	return nil
}
