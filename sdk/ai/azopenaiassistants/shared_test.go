//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func stringize(v azopenaiassistants.MessageContentClassification) string {
	switch m := v.(type) {
	case *azopenaiassistants.MessageTextContent:
		return fmt.Sprintf("Text = %s\n", *m.Text.Value)
	case *azopenaiassistants.MessageImageFileContent:
		return fmt.Sprintf("Image = %s\n", *m.ImageFile.FileID)
	}

	panic("Unhandled type for stringizing")
}

type mustGetClientWithAssistantArgs struct {
	newClientArgs
	Instructions string
}

type newClientArgs struct {
	Azure       bool
	UseIdentity bool
}

func mustGetClient(t *testing.T, args newClientArgs) *azopenaiassistants.Client {
	var httpClient policy.Transporter
	// var recordingPolicy
	// PerRetryPolicies: []{&mimeTypeRecordingPolicy{}}
	var perRetryPolicy policy.Policy

	if recording.GetRecordMode() != recording.LiveMode {
		err := recording.Start(t, RecordingDirectory, nil)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := recording.Stop(t, nil)
			require.NoError(t, err)
		})

		tmpHttpClient, err := recording.NewRecordingHTTPClient(t, nil)
		require.NoError(t, err)

		if recording.GetRecordMode() == recording.RecordingMode {
			err = recording.AddURISanitizer("https://openai.azure.com", strings.TrimRight(tv.AOAIEndpoint, "/"), nil)
			require.NoError(t, err)

			err = recording.AddURISanitizer("https://openai.azure.com", strings.TrimRight(tv.OpenAIEndpoint, "/"), nil)
			require.NoError(t, err)

			err = recording.AddHeaderRegexSanitizer("Api-Key", "key", "", nil)
			require.NoError(t, err)

			// add regexes for all of our resources
			for _, res := range []string{"thread", "run", "asst", "step", "msg"} {
				err = recording.AddURISanitizer("$1Sanitized$2", fmt.Sprintf("^(.+?/)%s_[^?/]+(.*)$", res), nil)
				require.NoError(t, err)

				err = recording.AddURISanitizer("$1Sanitized$2", fmt.Sprintf("^(.+?/)%s_[^?/]+(/.*)$", res), nil)
				require.NoError(t, err)
			}

			err := recording.AddURISanitizer("after=SANITIZED", `after=msg_[^&+]+`, nil)
			require.NoError(t, err)

			for _, res := range []string{"thread", "run", "assistant", "step", "message", "first", "last"} {
				// ie: "run_id": "run_SSUooivFMcO9KVYY99xUoFUG",
				err = recording.AddBodyRegexSanitizer(`$1:"Sanitized"`, fmt.Sprintf(`("%s_id")\s*:\s*"[^"]+"`, res), nil)
				require.NoError(t, err)
			}
		}

		httpClient = tmpHttpClient
		perRetryPolicy = &mimeTypeRecordingPolicy{}
	} else if os.Getenv("SSLKEYLOGFILE") != "" {
		file, err := os.OpenFile(os.Getenv("SSLKEYLOGFILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := file.Close()
			require.NoError(t, err)
		})

		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{
			KeyLogWriter: file,
		}
		httpClient = &http.Client{Transport: transport}
	}

	opts := &azopenaiassistants.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Logging: policy.LogOptions{
				IncludeBody: true,
			},
			Transport: httpClient,
		},
	}

	if perRetryPolicy != nil {
		opts.PerRetryPolicies = append(opts.PerRetryPolicies, perRetryPolicy)
	}

	if args.Azure {
		if args.UseIdentity {
			dac, err := azidentity.NewDefaultAzureCredential(nil)
			require.NoError(t, err)

			tmpClient, err := azopenaiassistants.NewClient(tv.AOAIEndpoint, dac, opts)
			require.NoError(t, err)
			return tmpClient
		} else {
			tmpClient, err := azopenaiassistants.NewClientWithKeyCredential(tv.AOAIEndpoint, azcore.NewKeyCredential(tv.AOAIKey), opts)
			require.NoError(t, err)
			return tmpClient
		}
	} else {
		tmpClient, err := azopenaiassistants.NewClientForOpenAI(tv.OpenAIEndpoint, azcore.NewKeyCredential(tv.OpenAIKey), opts)
		require.NoError(t, err)
		return tmpClient
	}
}

var weatherFunctionDefn = &azopenaiassistants.FunctionDefinition{
	Name: to.Ptr("get_current_weather"),
	Parameters: map[string]any{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]any{
			"location": map[string]any{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]any{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	},
}

func mustGetClientWithAssistant(t *testing.T, args mustGetClientWithAssistantArgs) (*azopenaiassistants.Client, azopenaiassistants.CreateAssistantResponse) {
	client := mustGetClient(t, args.newClientArgs)

	// give the assistant a random-ish name.
	id, err := recording.GenerateAlphaNumericID(t, "your-assistant-name", 6+len("your-assistant-name"), true)
	require.NoError(t, err)

	assistantName := id

	createResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
		Name:           &assistantName,
		DeploymentName: &assistantsModel,
		Instructions:   to.Ptr("You are a personal math tutor. Write and run code to answer math questions."),
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			&azopenaiassistants.CodeInterpreterToolDefinition{},

			// others...
			&azopenaiassistants.FunctionToolDefinition{
				Function: weatherFunctionDefn,
			},
			// &azopenaiassistants.RetrievalToolDefinition{}
		},
	}, nil)
	requireNoErr(t, args.Azure, err)

	t.Cleanup(func() {
		_, err := client.DeleteAssistant(context.Background(), *createResp.ID, nil)
		require.NoError(t, err)
	})

	return client, createResp
}

type runThreadArgs struct {
	newClientArgs
	Assistant azopenaiassistants.CreateAssistantBody
	Thread    azopenaiassistants.CreateAndRunThreadBody
}

func mustRunThread(ctx context.Context, t *testing.T, args runThreadArgs) (*azopenaiassistants.Client, []azopenaiassistants.ThreadMessage) {
	client := mustGetClient(t, args.newClientArgs)

	// give the assistant a random-ish name.
	assistantName, err := recording.GenerateAlphaNumericID(t, "your-assistant-name", 6+len("your-assistant-name"), true)
	require.NoError(t, err)

	if args.Assistant.Name == nil {
		args.Assistant.Name = &assistantName
	}

	args.Assistant.DeploymentName = &assistantsModel

	createResp, err := client.CreateAssistant(ctx, args.Assistant, nil)
	requireNoErr(t, args.Azure, err)

	t.Cleanup(func() {
		_, err := client.DeleteAssistant(ctx, *createResp.ID, nil)
		requireNoErr(t, args.Azure, err)
	})

	// create a thread and run it
	args.Thread.AssistantID = createResp.ID
	threadRunResp, err := client.CreateThreadAndRun(ctx, args.Thread, nil)
	requireNoErr(t, args.Azure, err)

	// poll for the thread end
	runStatus := pollForTests(t, ctx, client, *threadRunResp.ThreadID, *threadRunResp.ID, args.Azure)
	require.Equal(t, *runStatus.Status, azopenaiassistants.RunStatusCompleted)

	var allMessages []azopenaiassistants.ThreadMessage

	messagePager := client.NewListMessagesPager(*threadRunResp.ThreadID, &azopenaiassistants.ListMessagesOptions{
		Order: to.Ptr(azopenaiassistants.ListSortOrderAscending),
	})

	for messagePager.More() {
		page, err := messagePager.NextPage(ctx)
		requireNoErr(t, args.Azure, err)

		allMessages = append(allMessages, page.Data...)
	}

	return client, allMessages
}

type mimeTypeRecordingPolicy struct{}

// Do changes out the boundary for a multipart message. This makes it simpler to write
// recordings.
func (mrp *mimeTypeRecordingPolicy) Do(req *policy.Request) (*http.Response, error) {
	if recording.GetRecordMode() == recording.LiveMode {
		// this is strictly to make the IDs in the multipart body stable for test recordings.
		return req.Next()
	}

	// we'll fix up the multipart to make it more predictable for test recordings.
	//    Content-Type: multipart/form-data; boundary=787c880ce3dd11f9b6384d625c399c8490fc8989ceb6b7d208ec7426c12e

	contentType := req.Raw().Header[http.CanonicalHeaderKey("Content-type")]

	if len(contentType) == 0 {
		return req.Next()
	}

	mediaType, params, err := mime.ParseMediaType(contentType[0])

	if err != nil || mediaType != "multipart/form-data" {
		// we'll just assume our policy doesn't apply here.
		return req.Next()
	}

	origBoundary := params["boundary"]

	if origBoundary == "" {
		return nil, errors.New("Invalid use of this policy - no boundary was passed as part of the multipart mime type")
	}

	params["boundary"] = "boundary-for-recordings"

	// now let's update the body itself - we'll just do a simple string replacement. The entire purpose of the boundary string is to provide a
	// separator, which is distinct from the content.
	body := req.Body()
	defer body.Close()

	origBody, err := io.ReadAll(body)

	if err != nil {
		return nil, err
	}

	newBody := bytes.ReplaceAll(origBody, []byte(origBoundary), []byte("boundary-for-recordings"))

	if err := req.SetBody(streaming.NopCloser(bytes.NewReader(newBody)), mime.FormatMediaType(mediaType, params)); err != nil {
		return nil, err
	}

	return req.Next()
}

func getFileName(t *testing.T, ext string) *string {
	fileName := fmt.Sprintf("go-%s-%d.%s",
		t.Name(), rand.Int63(), ext)

	return &fileName
}

func requireNoErr(t *testing.T, azure bool, err error) {
	if responseErr := (*azcore.ResponseError)(nil); azure && errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusTooManyRequests {
		t.Skipf("Assistants API is being throttled: %s", responseErr)
	}

	require.NoError(t, err)
}

func requireSuccessfulPolling(t *testing.T, azure bool, resp azopenaiassistants.GetRunResponse, err error) {
	// it's possible we're oversubscribed, so we need to just skip this test
	if azure && resp.LastError != nil && resp.LastError.Code != nil && *resp.LastError.Code == "rate_limit_exceeded" {
		t.Skipf("Test being skipped, we're rate limited")
	}

	requireNoErr(t, azure, err)
}

func requireType[T azopenaiassistants.StreamEventDataClassification](t *testing.T, event azopenaiassistants.StreamEvent) T {
	v, ok := event.Event.(T)

	var zero T
	require.Truef(t, ok, "Expecting %T (%T)", zero, event.Event)
	return v
}

func getValue[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}
