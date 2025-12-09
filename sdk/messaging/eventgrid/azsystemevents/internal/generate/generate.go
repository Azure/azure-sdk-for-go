//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var filesToDelete = []string{
	"options.go",
	"responses.go",
	"clientdeleteme_client.go",
}

func deleteType(name string, modelsGo string, modelsSerdeGo string) (newModelsGo string, newModelsSerdeGo string) {
	// get rid of the type in models.go

	// ex:
	// '// ACSCallParticipantEventData - Schema of common properties of all participant events'
	deleteModelTypeRE := regexp.MustCompile(fmt.Sprintf(`(?si)// %s .+?\n}\n`, name))
	modelsGo = deleteModelTypeRE.ReplaceAllString(modelsGo, "")

	// now get rid of the marshallers

	// NOTE: sometimes these _wrap_, so make sure you account for that.
	// '// MarshalJSON implements the json.Marshaller interface for type ACSCallParticipantEventData.'
	// '// UnmarshalJSON implements the json.Unmarshaller interface for type ACSCallParticipantEventData.'

	deleteSerdeRE := regexp.MustCompile(fmt.Sprintf(`(?si)// (MarshalJSON|UnmarshalJSON) implements the json\.(Unmarshaller|Marshaller) interface for type %s.+?\n}\n`, name))
	modelsSerdeGo = deleteSerdeRE.ReplaceAllString(modelsSerdeGo, "")

	return modelsGo, modelsSerdeGo
}

func main() {
	fn := func() error {
		// NOTE: the renaming of the error type is done in the propertyNameOverrideGo.tsp (AcsMessageChannelEventError)
		modelsGoBytes, err := os.ReadFile(modelsGoFile)

		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", modelsGoFile, err)
		}

		modelsSerdeGoBytes, err := os.ReadFile(modelsSerdeGoFile)

		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", modelsSerdeGoFile, err)
		}

		modelsGo := string(modelsGoBytes)
		modelsSerdeGo := string(modelsSerdeGoBytes)

		modelsGo, modelsSerdeGo, err = swapErrorTypes(modelsGo, modelsSerdeGo)

		if err != nil {
			return fmt.Errorf("failed to swap error types: %w", err)
		}

		modelsGo, modelsSerdeGo = deleteUnneededTypes(modelsGo, modelsSerdeGo)

		if err := generateSystemEventEnum(modelsGo); err != nil {
			return fmt.Errorf("generateSystemEventEnum: %w", err)
		}

		deleteUnneededFiles()

		if err := os.WriteFile(modelsGoFile, []byte(modelsGo), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", modelsGoFile, err)
		}

		if err := os.WriteFile(modelsSerdeGoFile, []byte(modelsSerdeGo), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", modelsSerdeGoFile, err)
		}

		return nil
	}

	fmt.Printf("Running customizations in azsystemevents/internal/generate...\n")

	if err := fn(); err != nil {
		fmt.Printf("./internal/generate: Failed with error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("DONE\n")
}

// deleteUnneededTypes deletes types that are vestigial. Event Grid system events, in TypeSpec,
// use inheritance. We model this by pushing the parent properties into the child objects, so the
// parent object is not needed.
// NOTE: this is a workaround, I'm running into some conflicting behavior as I move between tsp-client
// and my "homebrew" version calling the powershell scripts directly.
func deleteUnneededTypes(modelsGo string, modelsSerdeGo string) (newModelsGo string, newModelsSerdeGo string) {
	typesToDelete := []string{
		"ACSCallParticipantEventData",
		"ACSChatEventBaseProperties",
		"ACSChatEventInThreadBaseProperties",
		"ACSChatMessageEventBaseProperties",
		"ACSChatMessageEventInThreadBaseProperties",
		"ACSChatThreadEventBaseProperties",
		"ACSChatThreadEventInThreadBaseProperties",
		"ACSMessageEventData",
		"ACSRouterEventData",
		"ACSRouterJobEventData",
		"ACSRouterWorkerEventData",
		"ACSSmsEventBaseProperties",
		"ACSSMSEventBaseProperties",
		"AppConfigurationSnapshotEventData",
		"AVSClusterEventData",
		"AVSPrivateCloudEventData",
		"AVSScriptExecutionEventData",
		"ContainerRegistryArtifactEventData",
		"ContainerRegistryEventData",
		"ContainerServiceClusterSupportEventData",
		"ContainerServiceNodePoolRollingEventData",
		"DeviceConnectionStateEventProperties",
		"DeviceLifeCycleEventProperties",
		"DeviceTelemetryEventProperties",
		"EventGridMQTTClientEventData",
		"MapsGeofenceEventProperties",
		"ResourceNotificationsResourceDeletedEventData",
		"ResourceNotificationsResourceUpdatedEventData",
	}

	newModelsGo = modelsGo
	newModelsSerdeGo = modelsSerdeGo

	for _, typeToDelete := range typesToDelete {
		// ex: // ACSChatEventBaseProperties - Schema of common properties of all chat events
		modelsRE := regexp.MustCompile(fmt.Sprintf("(?is)// %s.+?\n}\n", typeToDelete))
		newModelsGo = modelsRE.ReplaceAllString(newModelsGo, "")

		// ex: // MarshalJSON implements the json.Marshaller interface for type ACSChatEventBaseProperties.
		// ex: // UnmarshalJSON implements the json.Unmarshaller interface for type ACSChatEventBaseProperties.
		modelsSerdeRE := regexp.MustCompile(fmt.Sprintf("(?is)// (?:Marshal|Unmarshal)JSON\\s*implements\\s*the\\s*json\\.(?:Marshaller|Unmarshaller)\\s*interface\\s*for\\s*type\\s*%s.+?\n}\n", typeToDelete))
		newModelsSerdeGo = modelsSerdeRE.ReplaceAllString(newModelsSerdeGo, "")
	}

	return
}

// swapErrorTypes handles turning most of the auto-generated errors into a single consistent error type.
// The key is that the Error type doesn't export human readable strings as fields - it's all contained in
// the Error() field.
func swapErrorTypes(origModelsGo string, origModelsSerdeGo string) (newModelsGo string, newModelsSerdeGo string, err error) {
	newModelsGo = origModelsGo
	newModelsSerdeGo = origModelsSerdeGo

	// unexport the errors.
	newModelsGo = strings.ReplaceAll(newModelsGo, "InternalACSMessageChannelEventError", "internalACSMessageChannelEventError")
	newModelsGo = strings.ReplaceAll(newModelsGo, "InternalACSRouterCommunicationError", "internalACSRouterCommunicationError")
	newModelsSerdeGo = strings.ReplaceAll(newModelsSerdeGo, "InternalACSMessageChannelEventError", "internalACSMessageChannelEventError")
	newModelsSerdeGo = strings.ReplaceAll(newModelsSerdeGo, "InternalACSRouterCommunicationError", "internalACSRouterCommunicationError")

	// replace the types with the package level Error type that we use.
	newModelsGo = strings.ReplaceAll(newModelsGo, "Error *internalACSMessageChannelEventError", "Error *Error")
	newModelsGo = strings.ReplaceAll(newModelsGo, "Errors []internalACSRouterCommunicationError", "Errors []*Error")

	newModelsSerdeGo = UseCustomUnpopulate(newModelsSerdeGo, "ACSMessageReceivedEventData.Error", "unmarshalInternalACSMessageChannelEventError")
	newModelsSerdeGo = UseCustomUnpopulate(newModelsSerdeGo, "ACSMessageDeliveryStatusUpdatedEventData.Error", "unmarshalInternalACSMessageChannelEventError")

	newModelsSerdeGo = UseCustomUnpopulate(newModelsSerdeGo, "ACSRouterJobClassificationFailedEventData.Errors.Errors", "unmarshalInternalACSRouterCommunicationError")

	return newModelsGo, newModelsSerdeGo, nil
}

func generateSystemEventEnum(modelsGo string) error {
	constants, err := getConstantValues(modelsGo)

	if err != nil {
		return fmt.Errorf("failed to get constant values from file: %w", err)
	}

	if err := writeConstantsFile(systemEventsGoFile, constants); err != nil {
		return fmt.Errorf("failed to write constants file %s: %w", systemEventsGoFile, err)
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
