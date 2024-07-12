//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	t := &transformer{fileCache: NewFileCache()}

	if err := t.Do(); err != nil {
		log.Fatalf("Errors with transformer.Do(): %s", err)
	}

	if err := validate(); err != nil {
		log.Fatalf("Errors with validation(): %s", err)
	}
}

func validate() error {
	// check that nothing has an 'any' field
	// make sure we don't have any Paths<> blah types

	modelsBytes, err := os.ReadFile("models.go")

	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?m)^\s*([A-Za-z0-9]+)\s+any$`)
	matches := re.FindAllStringSubmatch(string(modelsBytes), -1)

	if matches == nil {
		fmt.Printf("PASS: No fields with type 'any' found!\n")
		return nil
	}

	// a type being 'any' is really just an indication that it was in the
	// Definitions in swagger but didn't contain any fields.
	var fields []string

	for _, match := range matches {
		// we can ignore this one - it's for function definitions
		// and it's intentionally an 'any' field
		if match[1] == "Parameters" {
			continue
		}
		fields = append(fields, match[1])
	}

	if len(fields) > 0 {
		sort.Strings(fields)
		return fmt.Errorf("%d field(s) were found with 'any' as the type\n  %s", len(fields), strings.Join(fields, "\n  "))
	}

	return nil
}

type transformer struct {
	fileCache *FileCache
}

func (t *transformer) Do() error {
	transforms := []func() error{
		t.removeClientPrefix,
		t.injectClientData,
		t.injectFormatURLHelper,
		t.hideListFunctions,
		t.fixFiles,
		t.fixStreaming,
		t.applyHacks,
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

	return transformFiles(t.fileCache, "injectFormatURLHelper", []string{"client.go"}, func(text string) (string, error) {
		return re.ReplaceAllString(text, "urlPath := client.formatURL($1)"), nil
	}, nil)
}

// injectClientData adds in our own user-defined struct so we don't have to keep
// editing client.go just to add in a new field we need.
func (t *transformer) injectClientData() error {
	return transformFiles(t.fileCache, "injectClientData", []string{"client.go"}, func(text string) (string, error) {
		newText := strings.Replace(text, "type Client struct {\n", "type Client struct {\ncd clientData\n", 1)

		return newText, nil
	}, &transformFileOptions{AllowNoop: true})
}

// hideListFunctions hides all the lists since we're supposed to expose pagers
// for these. (they don't fit the standard Azure pager pattern so aren't auto-generated)
func (t *transformer) hideListFunctions() error {
	return transformFiles(t.fileCache, "hideListFunctions", []string{"client.go"}, func(text string) (string, error) {
		funcsToHide := []string{
			"ListAssistantFiles",
			"ListAssistants",
			"ListMessageFiles",
			"ListMessages",
			"ListRunSteps",
			"ListRuns",
			"ListVectorStoreFileBatchFiles",
			"ListVectorStoreFiles",
			"ListVectorStores",
		}

		for _, funcToHide := range funcsToHide {
			text = strings.ReplaceAll(text, "func (client *Client) "+funcToHide, "func (client *Client) internal"+funcToHide)
		}

		return text, nil
	}, nil)
}

// removeClientPrefix removes the leading `Client` that gets prefixed onto every model.
func (t *transformer) removeClientPrefix() error {
	re := regexp.MustCompile(`Client([A-Z][A-Za-z]+)`)

	return transformFiles(t.fileCache, "removeClientPrefix", []string{"client.go", "models.go", "models_serde.go", "options.go", "responses.go"}, func(text string) (string, error) {
		return re.ReplaceAllString(text, "$1"), nil
	}, nil)
}

func (t *transformer) fixFiles() error {
	err := transformFiles(t.fileCache, "fixFiles", []string{"client.go"}, func(text string) (string, error) {
		text, err := removeFunctions(text, "Client", "uploadFileCreateRequest")

		if err != nil {
			return "", err
		}

		// removing these for now - the TypeSpec -> OpenAPI2 -> go generation is causing us to treat the `bytes`
		// field in TypeSpec as a deserialized []byte, instead of returning the Body's stream.
		getFileContentFunctions := []string{
			"getFileContentCreateRequest",
			"getFileContentHandleResponse",
			"GetFileContent",
		}

		return removeFunctions(text, "Client", getFileContentFunctions...)
	}, nil)

	if err != nil {
		return err
	}

	// this type is generated because UploadFile takes a multipart input argument - we have a manually generated
	// argument type for this instead.
	if err := removeTypes(t.fileCache, []string{"Paths1Filz8PFilesPostRequestbodyContentMultipartFormDataSchema"}, &removeTypesOptions{
		// this auto-gen'd type doesn't have a comment
		IgnoreComment: true,
	}); err != nil {
		return err
	}

	if err := removeTypes(t.fileCache, []string{"GetFileContentOptions", "GetFileContentResponse"}, nil); err != nil {
		return err
	}

	return transformFiles(t.fileCache, "fixFiles", []string{"client.go"}, func(text string) (string, error) {
		return strings.Replace(
			text,
			"func (client *Client) UploadFile(ctx context.Context, file io.ReadSeekCloser, purpose FilePurpose, options *UploadFileOptions) (UploadFileResponse, error) {",
			"func (client *Client) UploadFile(ctx context.Context, file io.ReadSeeker, purpose FilePurpose, options *UploadFileOptions) (UploadFileResponse, error) {", 1), nil
	}, nil)
}

func (t *transformer) applyHacks() error {
	// TODO: the 'fileID' parameter for "createVectorStoreFileCreateRequest" isn't being encoded as a body field, it's being encoded as the entire body.
	re := regexp.MustCompile(`(?s)(func \(client \*Client\) createVectorStoreFileCreateRequest\(.+?if err := runtime.MarshalAsJSON\(req, )fileID`)

	err := transformFiles(t.fileCache, "createVectorStoreFileCreateRequest fix", []string{"client.go"}, func(text string) (string, error) {
		text = re.ReplaceAllString(text, "${1}fileIDStruct{fileID}")
		return text, nil
	}, nil)

	if err != nil {
		return err
	}

	// TODO: for some reason the doc comments aren't making it over.
	docs := map[string]string{
		"CreateVectorStoreFileBatchBody": "// CreateVectorStoreFileBatchBody contains arguments for the [CreateVectorStoreFileBatch] method.",
		"SubmitToolOutputsToRunBody":     "// SubmitToolOutputsToRunBody contains arguments for the [SubmitToolOutputsToRun] method.",
		"UpdateMessageBody":              "// UpdateMessageBody contains arguments for the [UpdateMessage] method.",
		"UpdateRunBody":                  "// UpdateRunBody contains arguments for the [UpdateRun] method.",
	}

	for typeName, comment := range docs {
		err := transformFiles(t.fileCache, "update doc comments("+typeName+")", []string{"models.go"}, func(text string) (string, error) {
			return strings.Replace(text,
				fmt.Sprintf("type %s struct {", typeName),
				fmt.Sprintf("%s\ntype %s struct {", comment, typeName), 1), nil
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *transformer) fixStreaming() error {
	log.Printf("fixStreaming()...")

	// internalize all the 'stream' variables. We expose custom methods so we can give them back our streaming type.
	re := regexp.MustCompile(`(?s)// If true, returns a stream of events that.+?Stream \*bool`)

	err := transformFiles(t.fileCache, "make Stream an internal member with a doc that points to our stream functions", []string{"models.go"}, func(text string) (string, error) {
		text = re.ReplaceAllString(text, "// NOTE: Use the Stream version of this function (ex: [CreateThreadAndRunStream], [CreateRunStream] or [SubmitToolOutputsToRunStream]) to stream results.\nstream *bool")
		return text, nil
	}, nil)

	if err != nil {
		return err
	}

	err = transformFiles(t.fileCache, "fix serde for Stream", []string{"models_serde.go"}, func(text string) (string, error) {

		text = strings.ReplaceAll(text, `populate(objectMap, "stream", s.Stream)`, `populate(objectMap, "stream", s.stream)`)
		text = strings.ReplaceAll(text, `err = unpopulate(val, "Stream", &s.Stream)`, `err = unpopulate(val, "Stream", &s.stream)`)
		text = strings.ReplaceAll(text, `populate(objectMap, "stream", c.Stream)`, `populate(objectMap, "stream", c.stream)`)
		text = strings.ReplaceAll(text, `err = unpopulate(val, "Stream", &c.Stream)`, `err = unpopulate(val, "Stream", &c.stream)`)
		return text, nil
	}, nil)

	return err
}
