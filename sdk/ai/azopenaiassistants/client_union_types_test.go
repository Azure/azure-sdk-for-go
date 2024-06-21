//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestResponseFormatTypeUnion(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		for _, format := range azopenaiassistants.PossibleAssistantResponseFormatTypeValues() {
			resp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
				DeploymentName: &assistantsModel,
				ResponseFormat: &azopenaiassistants.AssistantResponseFormat{
					Type: format,
				},
			}, nil)
			requireNoErr(t, azure, err)
			defer mustDeleteAssistant(t, client, *resp.ID, azure)
		}
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestAPIToolChoiceUnion(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		createResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
			DeploymentName: &assistantsModel,
			Tools: []azopenaiassistants.ToolDefinitionClassification{
				&azopenaiassistants.CodeInterpreterToolDefinition{},
				&azopenaiassistants.FileSearchToolDefinition{},
				&azopenaiassistants.FunctionToolDefinition{
					Function: weatherFunctionDefn,
				},
			},
		}, nil)
		requireNoErr(t, azure, err)

		defer mustDeleteAssistant(t, client, *createResp.ID, azure)

		toolChoices := []*azopenaiassistants.AssistantsAPIToolChoiceOption{
			{Mode: azopenaiassistants.AssistantsAPIToolChoiceOptionModeAuto},
			{Mode: azopenaiassistants.AssistantsAPIToolChoiceOptionModeNone},
			{Mode: azopenaiassistants.AssistantsAPIToolChoiceOptionModeCodeInterpreter},
			{Mode: azopenaiassistants.AssistantsAPIToolChoiceOptionModeFileSearch},
			{Mode: azopenaiassistants.AssistantsAPIToolChoiceOptionModeFunction, Function: &azopenaiassistants.FunctionName{
				Name: to.Ptr("get_current_weather"),
			}},
		}

		for _, toolChoice := range toolChoices {
			resp, err := client.CreateThreadAndRun(context.Background(), azopenaiassistants.CreateAndRunThreadBody{
				AssistantID: createResp.ID,
				ToolChoice:  toolChoice,
				Thread: &azopenaiassistants.CreateThreadBody{
					Messages: []azopenaiassistants.CreateMessageBody{
						{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: to.Ptr("give me any answer")},
					},
				},
			}, nil)
			requireNoErr(t, azure, err)
			require.Equal(t, toolChoice, resp.ToolChoice)
		}
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}

func TestCodeInterpreterAndFileSearchMatchUnmarshaller(t *testing.T) {
	// NOTE: we deserialize these as part of the union type `MessageAttachmentToolDefinition`. If
	// fields are added to either one of these types they need to be accounted for in there so they
	// can be unmarshalled properly.
	getFieldNames := func(fields []reflect.StructField) []string {
		var names []string

		for _, field := range fields {
			names = append(names, field.Name)
		}

		return names
	}

	fields := reflect.VisibleFields(reflect.TypeOf(azopenaiassistants.CodeInterpreterToolDefinition{}))
	require.Equal(t, []string{"Type"}, getFieldNames(fields), "Fields match what we unmarshal")

	fields = reflect.VisibleFields(reflect.TypeOf(azopenaiassistants.FileSearchToolDefinition{}))
	require.Equal(t, []string{"Type"}, getFieldNames(fields), "Fields match what we unmarshal")

	fields = reflect.VisibleFields(reflect.TypeOf(azopenaiassistants.MessageAttachmentToolDefinition{}))
	require.Equal(t, []string{"CodeInterpreterToolDefinition", "FileSearchToolDefinition"}, getFieldNames(fields), "Fields match what we marshal")
}
