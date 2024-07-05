// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"encoding/json"
	"fmt"
)

// OpenAI's reference doc for Assistants: https://platform.openai.com/docs/api-reference/assistants

// MessageAttachmentToolDefinition specifies the tools the files being attached will be
// associated with.
type MessageAttachmentToolDefinition struct {
	// CodeInterpreterToolDefinition should be set if you attach files for the CodeInterpreter tool.
	CodeInterpreterToolDefinition *CodeInterpreterToolDefinition

	// FileSearchToolDefinition should be set if you attach files for the FileSearchToolDefinition tool.
	FileSearchToolDefinition *FileSearchToolDefinition
}

// MarshalJSON implements the json.Marshaller interface for type MessageAttachmentToolDefinition.
func (m MessageAttachmentToolDefinition) MarshalJSON() ([]byte, error) {
	if m.CodeInterpreterToolDefinition != nil {
		return json.Marshal(m.CodeInterpreterToolDefinition)
	}

	// if by chance neither of them is filled out then we'll just serialize a JSON null, which is correct. The service
	// does its own validation.
	return json.Marshal(m.FileSearchToolDefinition)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type MessageAttachmentToolDefinition.
func (m *MessageAttachmentToolDefinition) UnmarshalJSON(data []byte) error {
	var v *struct {
		Type string
		// NOTE: This type needs to include, uniquely, all the fields of the individual types.
		// You can't just embed them because fields that are common between them will
		// just get ignored.
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v.Type {
	case "code_interpreter":
		m.CodeInterpreterToolDefinition = &CodeInterpreterToolDefinition{
			Type: &v.Type,
		}
		return nil
	case "file_search":
		m.FileSearchToolDefinition = &FileSearchToolDefinition{
			Type: &v.Type,
		}
		return nil
	default:
		return fmt.Errorf("unhandled Type (%q) for MessageAttachmentToolDefinition", v.Type)
	}
}

// AssistantsAPIToolChoiceOption controls which tools are called by the model.
type AssistantsAPIToolChoiceOption struct {
	// Mode controls the behavior for this particular tool.
	//   - [AssistantsAPIToolChoiceOptionModeAuto] (default) means the model can pick between
	//     generating a message or calling a tool.
	//   - [AssistantsAPIToolChoiceOptionModeNone] means the model will not call any tools
	//     and instead generates a message.
	//
	// Alternately, this value can also be set to:
	//   - [AssistantsAPIToolChoiceOptionModeCodeInterpreter] - use the code interpreter tool.
	//   - [AssistantsAPIToolChoiceOptionModeFileSearch] - use the file search tool.
	//   - [AssistantsAPIToolChoiceOptionModeFunction] - use a function. The function name
	//     is set in the [AssistantsAPIToolChoiceOption.Function] field.
	Mode AssistantsAPIToolChoiceOptionMode

	// Function sets the name of the function to call.
	Function *FunctionName
}

// toolChoiceJSON matches the underlying JSON format we use when serializing
// the union types for toolChoice.
type toolChoiceJSON struct {
	Type     string        `json:"type"`
	Function *FunctionName `json:"function"`
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AssistantsAPIToolChoiceOption.
func (a *AssistantsAPIToolChoiceOption) UnmarshalJSON(data []byte) error {
	strValue, modelValue, err := unmarshalStringOrObject[toolChoiceJSON](data)

	if err != nil {
		return err
	}

	if modelValue != nil {
		a.Mode = AssistantsAPIToolChoiceOptionMode(modelValue.Type)
		a.Function = modelValue.Function
	} else {
		a.Mode = AssistantsAPIToolChoiceOptionMode(strValue)
	}

	return nil
}

// MarshalJSON implements the json.Marshaller interface for type AssistantsAPIToolChoiceOption.
func (a AssistantsAPIToolChoiceOption) MarshalJSON() ([]byte, error) {
	switch a.Mode {
	case AssistantsAPIToolChoiceOptionModeAuto, AssistantsAPIToolChoiceOptionModeNone:
		return json.Marshal(a.Mode)
	case AssistantsAPIToolChoiceOptionModeCodeInterpreter, AssistantsAPIToolChoiceOptionModeFileSearch, AssistantsAPIToolChoiceOptionModeFunction:
		return json.Marshal(toolChoiceJSON{
			Type:     string(a.Mode),
			Function: a.Function,
		})
	}

	return nil, nil
}

// AssistantResponseFormat controls the response format of tool calls made by an assistant.
type AssistantResponseFormat struct {
	// Type indicates which format type should be used for tool calls.
	// The default is [AssistantResponseFormatTypeAuto].
	Type AssistantResponseFormatType
}

// responseFormatJSON matches the underlying JSON format we use when serializing
// the union types for toolChoice.
type responseFormatJSON struct {
	Type string `json:"type"`
}

// UnmarshalJSON implements the json.Unmarshaller interface for type AssistantResponseFormat.
func (a *AssistantResponseFormat) UnmarshalJSON(data []byte) error {
	strValue, modelValue, err := unmarshalStringOrObject[responseFormatJSON](data)

	if err != nil {
		return err
	}

	if modelValue != nil {
		a.Type = AssistantResponseFormatType(modelValue.Type)
	} else {
		a.Type = AssistantResponseFormatType(strValue)
	}

	return nil
}

// MarshalJSON implements the json.Marshaller interface for type AssistantResponseFormat.
func (a AssistantResponseFormat) MarshalJSON() ([]byte, error) {
	switch a.Type {
	case AssistantResponseFormatTypeAuto:
		return json.Marshal(a.Type)
	case AssistantResponseFormatTypeJSONObject, AssistantResponseFormatTypeText:
		return json.Marshal(responseFormatJSON{
			Type: string(a.Type),
		})
	default:
		return nil, fmt.Errorf("unknown type %s, failed to marshal value", a.Type)
	}
}

// CreateFileSearchToolResourceOptions is set of resources that are used by the file search
// tool.
type CreateFileSearchToolResourceOptions struct {
	// VectorStoreIDs are the vector stores that will be attached to this assistant.
	// NOTE: There can be a maximum of 1 vector store attached to the assistant.
	VectorStoreIDs []string `json:"vector_store_ids"`

	// VectorStores can be set to create a vector store with file_ids and attach it to
	// this assistant.
	// NOTE: There can be a maximum of 1 vector store attached to the assistant.
	VectorStores []CreateFileSearchToolResourceVectorStoreOptions `json:"vector_stores"`
}

// unmarshalStringOrObject checks to see if the jsonBytes are actually a JSON serialized string
// and, otherwise, assumes it's a JSON object.
func unmarshalStringOrObject[T any](jsonBytes []byte) (string, *T, error) {
	if len(jsonBytes) == 0 {
		return "", nil, fmt.Errorf("can't deserialize from an empty slice of bytes")
	}

	if jsonBytes[0] == '"' { // ie: it's a JSON string, not a JSON object
		var str *string

		if err := json.Unmarshal(jsonBytes, &str); err != nil {
			return "", nil, err
		}

		return *str, nil, nil
	}

	var model *T

	if err := json.Unmarshal(jsonBytes, &model); err != nil {
		return "", nil, err
	}

	return "", model, nil
}

// this is a workaround until I figure out where the generation is going awry for the vector store
type fileIDStruct struct {
	FileID string `json:"file_id"`
}
