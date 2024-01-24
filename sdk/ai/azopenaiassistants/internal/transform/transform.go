//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {
	t := &transformer{fileCache: NewFileCache()}

	if err := t.Do(); err != nil {
		log.Fatal(err)
	}
}

type transformer struct {
	fileCache *FileCache
}

func (t *transformer) Do() error {
	transforms := []func() error{
		t.injectClientData,
		t.injectFormatURLHelper,
		t.hideListFunctions,
		t.fixBodyArgs,
		t.removeUnusedMultipartModel,
		t.renameInnerPageObjects,
		t.renameModelToDeploymentName,
		t.fixNilCheck,
		t.hackFixTimestamps,
		// /files changes
		t.fixMultipart,
		t.fixFilenameType,
		//t.fixFileName,

		t.removeClientPrefix,
	}

	for _, tr := range transforms {
		if err := tr(); err != nil {
			return err
		}
	}

	// write all modified files
	if err := t.fileCache.WriteAll(); err != nil {
		return err
	}

	return nil
}

func (t *transformer) injectFormatURLHelper() error {
	// urlPath := "/threads/{threadId}/runs/{runId}/cancel"
	re := regexp.MustCompile(`(?m)^\s+urlPath := (.+)$`)

	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		return re.ReplaceAllString(text, "urlPath := client.formatURL($1)"), nil
	}, nil)
}

// injectClientData adds in our own user-defined struct so we don't have to keep
// editing client.go just to add in a new field we need.
func (t *transformer) injectClientData() error {
	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		newText := strings.Replace(text, "type Client struct {\n", "type Client struct {\nclientData\n", 1)

		return newText, nil
	}, &transformFileOptions{AllowNoop: true})
}

func (t *transformer) renameModelToDeploymentName() error {
	// we've standardized on 'DeploymentName' when you're specifying a model.

	// Fix the names of the structs
	// Model *string

	err := transformFiles(t.fileCache, []string{"models.go"}, func(text string) (string, error) {
		return strings.Replace(text, "Model *string", "DeploymentName *string", -1), nil
	}, nil)

	if err != nil {
		return err
	}

	// Fix the marshalling of the struct
	// err = unpopulate(val, "Model", &a.Model)
	// populate(objectMap, "model", a.Model)
	popRE := regexp.MustCompile(`(?m)^\s+populate\(objectMap, "model", ([a-zA-Z]).Model\)`)
	unpopRE := regexp.MustCompile(`(?m)^\s+err = unpopulate\(val, "Model", &([a-zA-Z]).Model\)`)

	return transformFiles(t.fileCache, []string{"models_serde.go"}, func(text string) (string, error) {
		text = popRE.ReplaceAllString(text, `populate(objectMap, "model", $1.DeploymentName)`)
		text = unpopRE.ReplaceAllString(text, `err = unpopulate(val, "Model", &$1.DeploymentName)`)
		return text, nil
	}, nil)
}

// hideListFunctions hides all the lists since we're supposed to expose pagers
// for these. (they don't fit the standard Azure pager pattern so aren't auto-generated)
func (t *transformer) hideListFunctions() error {
	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		funcsToHide := []string{
			"ListAssistantFiles",
			"ListAssistants",
			"ListMessageFiles",
			"ListMessages",
			"ListRunSteps",
			"ListRuns",
		}

		for _, funcToHide := range funcsToHide {
			text = strings.Replace(text, "func (client *Client) "+funcToHide, "func (client *Client) internal"+funcToHide, -1)
		}

		return text, nil
	}, nil)
}

// fixBodyArgs fixes the generated models from TypeSpec that are (apparently) supposed
// to have an implicit name generated for them. I think this should overall just go away and be replaced by
// what Joel's adding to our direct-from-TypeSpec generator.
func (t *transformer) fixBodyArgs() error {
	// find them
	// ex: func (client *Client) UpdateMessage(ctx context.Context, threadID string, messageID string, body Paths12Hz0B8ThreadsThreadidMessagesMessageidPostRequestbodyContentApplicationJSONSchema, options *ClientUpdateMessageOptions) (ClientUpdateMessageResponse, error) {

	replacements := map[string]string{}

	// match functions that have a 'body' parameter that's got the long
	// PathsWsxzpAssistantsAssistantidFilesGetResponses200ContentApplicationJSONSchema style name.
	anonModelRE := regexp.MustCompile(`(?m)^func \(client \*Client\) ([A-Z].+?)\(ctx.+body (Paths[^,]+),`)

	err := transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		matches := anonModelRE.FindAllStringSubmatch(text, -1)

		for _, match := range matches {
			operation, anonModelName := match[1], match[2]

			newModelName := fmt.Sprintf("%sBody", operation)
			replacements[anonModelName] = newModelName

			text = strings.Replace(text, anonModelName, newModelName, -1)
		}

		return text, nil
	}, nil)

	if err != nil {
		return err
	}

	err = transformFiles(t.fileCache, []string{"models.go", "models_serde.go"}, func(text string) (string, error) {
		for oldName, newName := range replacements {
			text = strings.Replace(text, oldName, newName, -1)
		}
		return text, nil
	}, nil)

	if err != nil {
		return err
	}

	// We have a few that have to be replaced manually.
	err = transformFiles(t.fileCache, []string{"models.go", "models_serde.go"}, func(text string) (string, error) {
		// rename the types so they're 'Body' instead of 'Options' (these weren't the Options bag types for the function)
		text = strings.Replace(text, "CreateRunOptions", "CreateRunBody", -1)
		text = strings.Replace(text, "UpdateAssistantOptions", "UpdateAssistantBody", -1)
		text = strings.Replace(text, "AssistantCreationOptions", "AssistantCreationBody", -1)
		return text, nil
	}, nil)

	if err != nil {
		return err
	}

	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		text = strings.Replace(text, "createRunOptions CreateRunOptions", "body CreateRunBody", -1)
		text = strings.ReplaceAll(text,
			`req, err := client.createRunCreateRequest(ctx, threadID, createRunOptions, options)`,
			`req, err := client.createRunCreateRequest(ctx, threadID, body, options)`)
		text = strings.ReplaceAll(text,
			`if err := runtime.MarshalAsJSON(req, createRunOptions); err != nil {`,
			`if err := runtime.MarshalAsJSON(req, body); err != nil {`)

		text = strings.Replace(text, "AssistantCreationOptions", "AssistantCreationBody", -1)
		return text, nil
	}, nil)
}

// renameInnerPageObjects gives names to the anonymous inner objects the Swagger has for unnamed data contained
// within a single page of results.
// For now, I'm just renaming the inner ones manually.
func (t *transformer) renameInnerPageObjects() error {
	regexp.MustCompile(`^`)

	renames := map[string]string{
		"PathsWsxzpAssistantsAssistantidFilesGetResponses200ContentApplicationJSONSchema":              "AssistantFilesPage",
		"Paths1Ih5M1JAssistantsGetResponses200ContentApplicationJSONSchema":                            "AssistantsPage",
		"Paths17M2HqjThreadsThreadidMessagesMessageidFilesGetResponses200ContentApplicationJSONSchema": "MessageFilesPage",
		"Paths783Jj4ThreadsThreadidMessagesGetResponses200ContentApplicationJSONSchema":                "MessagesPage",
		"PathsPia9TjThreadsThreadidRunsRunidStepsGetResponses200ContentApplicationJSONSchema":          "RunStepsPage",
		"PathsMc8ByoThreadsThreadidRunsGetResponses200ContentApplicationJSONSchema":                    "ThreadRunsPage",
	}

	return transformFiles(t.fileCache, []string{"client.go", "models.go", "models_serde.go", "response_types.go"}, func(text string) (string, error) {
		for search, replace := range renames {
			text = strings.ReplaceAll(text, search, replace)
		}
		return text, nil
	}, nil)
}

func (t *transformer) fixNilCheck() error {
	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		text = strings.ReplaceAll(text,
			`func (client *Client) uploadFileCreateRequest(ctx context.Context, file string, purpose FilePurpose, options *UploadFileOptions) (*policy.Request, error) {`,
			"func (client *Client) uploadFileCreateRequest(ctx context.Context, file string, purpose FilePurpose, options *UploadFileOptions) (*policy.Request, error) {\nvar fileName *string\nif options == nil { fileName = options.Filename }\n")

		text = strings.ReplaceAll(text,
			`"Filename": options.Filename,`,
			`"Filename": fileName,`)

		return text, nil
	}, nil)
}

// removeClientPrefix removes the leading `Client` that gets prefixed onto every model.
func (t *transformer) removeClientPrefix() error {
	re := regexp.MustCompile(`Client([A-Z][A-Za-z]+)`)

	return transformFiles(t.fileCache, []string{"client.go", "models.go", "models_serde.go", "options.go", "response_types.go"}, func(text string) (string, error) {
		return re.ReplaceAllString(text, "$1"), nil
	}, nil)
}

func (t *transformer) removeUnusedMultipartModel() error {
	return removeType(t.fileCache, "Paths1Filz8PFilesPostRequestbodyContentMultipartFormDataSchema")
}

func (t *transformer) fixMultipart() error {
	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		return removeFunction(text, "Client", "uploadFileCreateRequest")
	}, nil)
}

func (t *transformer) fixFilenameType() error {
	return transformFiles(t.fileCache, []string{"client.go"}, func(text string) (string, error) {
		return strings.Replace(
			text,
			"func (client *Client) UploadFile(ctx context.Context, file string",
			"func (client *Client) UploadFile(ctx context.Context, file []byte", 1), nil
	}, nil)
}
