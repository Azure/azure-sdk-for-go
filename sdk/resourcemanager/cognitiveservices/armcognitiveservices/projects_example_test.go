// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcognitiveservices_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices/v4"
	"github.com/stretchr/testify/require"
)

func ExampleProjectsClient_BeginUpdate_clearAFoundryProjectSToolControlsWithExplicitNull() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	ctx := context.TODO()
	clientFactory, err := armcognitiveservices.NewClientFactory("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	poller, err := clientFactory.NewProjectsClient().BeginUpdate(ctx, "myResourceGroup", "myAccount", "myProject", armcognitiveservices.Project{
		Properties: &armcognitiveservices.ProjectProperties{
			ToolControls: azcore.NullValue[*armcognitiveservices.ToolControls](),
		},
	}, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func TestProjectToolControlsNullJSON(t *testing.T) {
	for _, test := range []struct {
		name         string
		toolControls *armcognitiveservices.ToolControls
		expected     string
	}{
		{name: "explicit null clears controls", toolControls: azcore.NullValue[*armcognitiveservices.ToolControls](), expected: `{"properties":{"toolControls":null}}`},
		{name: "nil omits controls", expected: `{"properties":{}}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(armcognitiveservices.Project{
				Properties: &armcognitiveservices.ProjectProperties{ToolControls: test.toolControls},
			})
			require.NoError(t, err)
			require.JSONEq(t, test.expected, string(body))
		})
	}
}
