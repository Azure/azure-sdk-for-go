//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestResponseFormatTypeUnion(t *testing.T) {
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
			require.NoError(t, err)
			defer mustDeleteAssistant(t, client, *resp.ID)
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
		require.NoError(t, err)

		defer mustDeleteAssistant(t, client, *createResp.ID)

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
			require.NoError(t, err)
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
