//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azopenaiassistants

import "encoding/json"

func unmarshalMessageContentClassification(rawMsg json.RawMessage) (MessageContentClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b MessageContentClassification
	switch m["type"] {
	case "image_file":
		b = &MessageImageFileContent{}
	case "text":
		b = &MessageTextContent{}
	default:
		b = &MessageContent{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalMessageContentClassificationArray(rawMsg json.RawMessage) ([]MessageContentClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]MessageContentClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalMessageContentClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalMessageTextAnnotationClassification(rawMsg json.RawMessage) (MessageTextAnnotationClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b MessageTextAnnotationClassification
	switch m["type"] {
	case "file_citation":
		b = &MessageTextFileCitationAnnotation{}
	case "file_path":
		b = &MessageTextFilePathAnnotation{}
	default:
		b = &MessageTextAnnotation{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalMessageTextAnnotationClassificationArray(rawMsg json.RawMessage) ([]MessageTextAnnotationClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]MessageTextAnnotationClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalMessageTextAnnotationClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalRequiredToolCallClassification(rawMsg json.RawMessage) (RequiredToolCallClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b RequiredToolCallClassification
	switch m["type"] {
	case "function":
		b = &RequiredFunctionToolCall{}
	default:
		b = &RequiredToolCall{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalRequiredToolCallClassificationArray(rawMsg json.RawMessage) ([]RequiredToolCallClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]RequiredToolCallClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalRequiredToolCallClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalRunStepCodeInterpreterToolCallOutputClassification(rawMsg json.RawMessage) (RunStepCodeInterpreterToolCallOutputClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b RunStepCodeInterpreterToolCallOutputClassification
	switch m["type"] {
	case "image":
		b = &RunStepCodeInterpreterImageOutput{}
	case "logs":
		b = &RunStepCodeInterpreterLogOutput{}
	default:
		b = &RunStepCodeInterpreterToolCallOutput{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalRunStepCodeInterpreterToolCallOutputClassificationArray(rawMsg json.RawMessage) ([]RunStepCodeInterpreterToolCallOutputClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]RunStepCodeInterpreterToolCallOutputClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalRunStepCodeInterpreterToolCallOutputClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalRunStepDetailsClassification(rawMsg json.RawMessage) (RunStepDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b RunStepDetailsClassification
	switch m["type"] {
	case string(RunStepTypeMessageCreation):
		b = &RunStepMessageCreationDetails{}
	case string(RunStepTypeToolCalls):
		b = &RunStepToolCallDetails{}
	default:
		b = &RunStepDetails{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalRunStepToolCallClassification(rawMsg json.RawMessage) (RunStepToolCallClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b RunStepToolCallClassification
	switch m["type"] {
	case "code_interpreter":
		b = &RunStepCodeInterpreterToolCall{}
	case "function":
		b = &RunStepFunctionToolCall{}
	case "retrieval":
		b = &RunStepRetrievalToolCall{}
	default:
		b = &RunStepToolCall{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalRunStepToolCallClassificationArray(rawMsg json.RawMessage) ([]RunStepToolCallClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]RunStepToolCallClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalRunStepToolCallClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalToolDefinitionClassification(rawMsg json.RawMessage) (ToolDefinitionClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ToolDefinitionClassification
	switch m["type"] {
	case "code_interpreter":
		b = &CodeInterpreterToolDefinition{}
	case "function":
		b = &FunctionToolDefinition{}
	case "retrieval":
		b = &RetrievalToolDefinition{}
	default:
		b = &ToolDefinition{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalToolDefinitionClassificationArray(rawMsg json.RawMessage) ([]ToolDefinitionClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ToolDefinitionClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalToolDefinitionClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}
