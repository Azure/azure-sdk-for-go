// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces/internal/generate/gopls"
)

func main() {
	fmt.Printf("checkRequestMethods...\n")
	if err := checkRequestMethods(); err != nil {
		panic(err)
	}

	fmt.Printf("deleteModels...\n")
	if err := deleteModels(); err != nil {
		panic(err)
	}

	fmt.Printf("injectClientData...\n")
	if err := injectClientData(); err != nil {
		panic(err)
	}

	fmt.Printf("exportPublicSymbols...\n")
	if err := renameSymbols(); err != nil {
		panic(err)
	}

	fmt.Printf("useAZCoreCloudEvent...\n")
	if err := useAZCoreCloudEvent(); err != nil {
		panic(err)
	}

	fmt.Printf("transformError...\n")
	if err := transformError(); err != nil {
		panic(err)
	}

	// fix result types
	fmt.Printf("fix result types...\n")
	if err := fixResultTypes(); err != nil {
		panic(err)
	}
}

func useAZCoreCloudEvent() error {
	if err := replaceAllInFile("models.go", regexp.MustCompile(`(?m)^\s+Event \*CloudEvent$`), []byte("\tEvent messaging.CloudEvent")); err != nil {
		return err
	}

	if err := replaceAllInFile("sender_client.go", regexp.MustCompile(`(?m)event CloudEvent`), []byte("\tevent *messaging.CloudEvent")); err != nil {
		return err
	}

	if err := replaceAllInFile("sender_client.go", regexp.MustCompile(`(?m)events \[\]CloudEvent`), []byte("\tevents []*messaging.CloudEvent")); err != nil {
		return err
	}

	return nil
}

func deleteModels() error {
	models := "CloudEvent|InnerError|PublishResult"

	if err := replaceAllInFile(
		"models_serde.go",
		regexp.MustCompile(
			fmt.Sprintf(`(?s)// (Unmarshal|Marshal)JSON implements the json\.(Unmarshaller|Marshaller) interface for type (%s).+?\n}`, models)), nil); err != nil {
		return fmt.Errorf("failed to remove marshaller/unmarshaller for %s: %w", models, err)
	}

	if err := replaceAllInFile(
		"models.go",
		regexp.MustCompile(fmt.Sprintf(`(?s)// (%s).+?type (%s) struct.+?\n}`, models, models)), nil); err != nil {
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
		"AcknowledgeEventsOptions": true,
		"RejectEventsOptions":      true,
		"ReleaseEventsOptions":     true,
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
func renameSymbols() error {
	// find all the things we need to rename
	if err := renameStructs("responses.go", regexp.MustCompile(`^(?:ReceiverClient|SenderClient)(.+)$`), "$1"); err != nil {
		return fmt.Errorf("failed to rename ReceiverClient/SenderClient prefixed types in responses.go")
	}

	if err := renameStructs("options.go", regexp.MustCompile(`^(?:ReceiverClient|SenderClient)(.+)$`), "$1"); err != nil {
		return fmt.Errorf("failed to rename ReceiverClient/SenderClient prefixed types in options.go")
	}

	re := regexp.MustCompile(`(?m)^func \(client \*SenderClient\) (SendEvent|SendEvents)\(`)

	if err := replaceAllInFile("sender_client.go", re, []byte("func (client *SenderClient) internal$1(")); err != nil {
		return fmt.Errorf("failed to rename SendEvent/SendEvents to be internal: %w", err)
	}

	re = regexp.MustCompile(`(?m)^func \(client \*ReceiverClient\) (AcknowledgeEvents|ReceiveEvents|ReleaseEvents|RejectEvents|RenewEventLocks)\(`)
	if err := replaceAllInFile("receiver_client.go", re, []byte("func (client *ReceiverClient) internal$1(")); err != nil {
		return fmt.Errorf("failed to rename Receiver functions to be internal: %w", err)
	}

	return nil
}

func renameStructs(file string, re *regexp.Regexp, replace string) error {
	symbols, err := gopls.SymbolsSlice(file)

	if err != nil {
		return err
	}

	// we do it in reverse order so we can do all the renames without reloading the file
	// since our positions should not shift.
	for i := len(symbols) - 1; i >= 0; i-- {
		if symbols[i].Type != gopls.SymbolTypeStruct {
			continue
		}

		replaced := re.ReplaceAllString(symbols[i].Name, replace)

		if replaced == symbols[i].Name {
			continue
		}

		fmt.Printf("Renaming %s -> %s\n", symbols[i].Name, replaced)

		if err := gopls.Rename(symbols[i], replaced); err != nil {
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
	// clip out some of the internal fields
	err := replaceAllInFile("models.go", regexp.MustCompile(`(?s)// [^/]+?\s+(Innererror \*InnerError|Details \[\]Error|Target \*string)\n`), nil)

	if err != nil {
		return fmt.Errorf("failed to splice out the error fields in models.go: %w", err)
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

func fixResultTypes() error {
	symbols, err := gopls.SymbolsSlice("models.go")

	if err != nil {
		return err
	}

	renames := map[string]string{
		"AcknowledgeResult": "AcknowledgeEventsResult",
		"ReceiveResult":     "ReceiveEventsResult",
		"RejectResult":      "RejectEventsResult",
		"ReleaseResult":     "ReleaseEventsResult",
		"RenewLocksResult":  "RenewEventLocksResult",
	}

	for i := len(symbols) - 1; i >= 0; i-- {
		newName := renames[symbols[i].Name]

		if newName == "" {
			continue
		}

		if err := gopls.Rename(symbols[i], newName); err != nil {
			return err
		}
	}

	{
		re := regexp.MustCompile(`(?s)if err := runtime.UnmarshalAsJSON\(resp, &result.PublishResult\); err != nil \{.+?return result, nil`)
		data, err := os.ReadFile("sender_client.go")

		if err != nil {
			return err
		}

		newData := re.ReplaceAll(data, []byte("return result, nil"))

		newData = bytes.Replace(newData, []byte(`sendEventHandleResponse(resp *http.Response)`), []byte(`sendEventHandleResponse(_ *http.Response)`), 1)
		newData = bytes.Replace(newData, []byte(`sendEventsHandleResponse(resp *http.Response)`), []byte(`sendEventsHandleResponse(_ *http.Response)`), 1)

		if err := os.WriteFile("sender_client.go", newData, 0600); err != nil {
			return err
		}
	}

	{
		data, err := os.ReadFile("responses.go")

		if err != nil {
			return err
		}

		re := regexp.MustCompile(`(?s)// The result of the Publish operation.\s+PublishResult`)

		newData := re.ReplaceAll(data, nil)

		if err := os.WriteFile("responses.go", newData, 0600); err != nil {
			return err
		}
	}

	return nil
}
