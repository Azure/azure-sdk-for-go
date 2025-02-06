# Go

These settings apply only when `--go` is specified on the command line.

``` yaml
input-file:
# PR: https://github.com/Azure/azure-rest-api-specs/pull/27076/files
#- https://raw.githubusercontent.com/Azure/azure-rest-api-specs/18c24352ad4a2e0959c0b4ec1404c3a250912f8b/specification/ai/data-plane/OpenAI.Assistants/OpenApiV2/preview/2024-02-15-preview/assistants_generated.json
- ./testdata/generated/openapi.json
output-folder: ../azopenaiassistants
clear-output-folder: false
module: github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: data-plane
go: true
title: "OpenAIAssistants"
use: "@autorest/go@4.0.0-preview.63"
slice-elements-byval: true
rawjson-as-bytes: true
# can't use this since it removes an innererror type that we want ()
# remove-non-reference-schema: true
```

## Transformations

Fix deployment and endpoint parameters so they show up in the right spots

``` yaml
directive:
  # Add x-ms-parameter-location to parameters in x-ms-parameterized-host
  - from: swagger-document
    where: $["x-ms-parameterized-host"].parameters.0
    transform: $["x-ms-parameter-location"] = "client";

  # fix a generation issue where "| null" in TypeSpec generates an allOf that
  # doesn't work with our polymorphic types.
  - from: swagger-document
    where: $.definitions.ThreadRun.properties.required_action
    transform: |
      $["$ref"] = "#/definitions/RequiredAction"
      delete $["allOf"];
      delete $["type"];
      return $;
```

Renaming the `Options` models to be `Body` so they don't clash with our methods names.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $.VectorStoreUpdateOptions["x-ms-client-name"] = "VectorStoreUpdateBody";
      $.VectorStoreOptions["x-ms-client-name"] = "VectorStoreBody";
      $.UpdateAssistantThreadOptions["x-ms-client-name"] = "UpdateThreadBody";
      $.UpdateAssistantOptions["x-ms-client-name"] = "UpdateAssistantBody";
      $.ThreadMessageOptions["x-ms-client-name"] = "CreateMessageBody";
      $.CreateRunOptions["x-ms-client-name"] = "CreateRunBody";
      $.CreateAndRunThreadOptions["x-ms-client-name"] = "CreateAndRunThreadBody";
      $.AssistantThreadCreationOptions["x-ms-client-name"] = "CreateThreadBody";
      $.AssistantCreationOptions["x-ms-client-name"] = "CreateAssistantBody";

      // These have 'Options' in the name, but they're not the main arguments for a function
      // and don't conflict with the normal functions options bag. So we don't need to rename these.
      // $.CreateFileSearchToolResourceVectorStoreOptions["x-ms-client-name"] = "CreateFileSearchToolResourceVectorStoreBody";
      // $.CreateToolResourcesOptions["x-ms-client-name"] = "CreateToolResourcesBody";
      // $.UpdateFileSearchToolResourceOptions["x-ms-client-name"] = "UpdateFileSearchToolResourceBody";
      // $.UpdateCodeInterpreterToolResourceOptions["x-ms-client-name"] = "UpdateCodeInterpreterToolResourceBody";
      // $.UpdateToolResourcesOptions["x-ms-client-name"] = "UpdateToolResourcesBody";
      // $.CreateCodeInterpreterToolResourceOptions["x-ms-client-name"] = "CreateCodeInterpreterToolResourceBody";
      return $;
```

## Unions

APIToolChoice

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      // we manually generate this type
      $.AssistantsApiToolChoiceOption = {
        "x-ms-external": true,
        "type": "object",
        "properties": { "ignored": { "type": "string" } },
        "x-ms-client-name": "AssistantsAPIToolChoiceOption"
      };

      // combine the two AssistantsApiToolChoiceOptionMode and AssistantsNamedToolChoiceType
      // into a single enum (similar to what they're doing in TypeSpec).
      $.AssistantsApiToolChoiceOptionMode["enum"] = [
        ...$.AssistantsApiToolChoiceOptionMode["enum"],
        ...$.AssistantsNamedToolChoiceType["enum"]
      ];
      $.AssistantsApiToolChoiceOptionMode["x-ms-enum"].values = [
        ...$.AssistantsApiToolChoiceOptionMode["x-ms-enum"].values,
        ...$.AssistantsNamedToolChoiceType["x-ms-enum"].values
      ];
      // fix the API casing of the mode
      $.AssistantsApiToolChoiceOptionMode["x-ms-client-name"] = "AssistantsAPIToolChoiceOptionMode";
      delete $.AssistantsNamedToolChoiceType;
      delete $.AssistantsNamedToolChoice;
```

ResponseFormat

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      // combine the two AssistantsApiResponseFormatMode and ApiResponseFormat
      // into a single enum (similar to what they're doing in TypeSpec).
      $.AssistantsApiResponseFormatOption = {
        "x-ms-external": true,
        "type": "object",
        "properties": { "ignored": { "type": "string" } },
        "x-ms-client-name": "AssistantResponseFormat"
      };

      const dest = $.AssistantsApiResponseFormatMode;

      dest["enum"] = [
        ...dest["enum"],
        ...$.ApiResponseFormat["enum"]
      ];
      dest["x-ms-enum"].values = [
        ...dest["x-ms-enum"].values,
        ...$.ApiResponseFormat["x-ms-enum"].values
      ];
      dest["x-ms-enum"].name = "AssistantResponseFormatType";

      // The 'none' option should be deleted - it seems to only be specified so you can get a 404? Let's just remove it.
      dest.enum = dest.enum.filter(name => name !== "none");
      dest["x-ms-enum"].values = dest["x-ms-enum"].values.filter(value => value.name !== "none");

      delete $.ApiResponseFormat;
      return $;
```

CreateFileSearchToolResourceOptions

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      // combine the two AssistantsApiResponseFormatMode and ApiResponseFormat
      // into a single enum (similar to what they're doing in TypeSpec).
      $.CreateFileSearchToolResourceOptions = {
        "x-ms-external": true,
        "type": "object",
        "properties": { "ignored": { "type": "string" } },
      };
      return $;
```

MessageAttachmentToolDefinition

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      // create the dummy type that I can use to point to my manually created union
      $.MessageAttachmentToolDefinition = {
        "x-ms-external": true,
        "type": "object",
        "properties": { "ignored": { "type": "string" } },
        "x-ms-client-name": "MessageAttachmentToolDefinition"
      };
```

## Model -> DeploymentName

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $.UpdateAssistantOptions.properties["model"]["x-ms-client-name"] = "DeploymentName";
      $.AssistantCreationOptions.properties["model"]["x-ms-client-name"] = "DeploymentName";
      $.CreateAndRunThreadOptions.properties["model"]["x-ms-client-name"] = "DeploymentName";
      return $;
```

## Fix anonymous named objects

Give names to anonymous objects defined within the operation definitions (paging operations). These are easily identified because the comment is always 'The response data for a requested list of items'

```yaml
directive:
  - from: swagger-document
    where: $.paths
    transform: |
      //
      // assistants
      //
      $["/assistants"].get.responses["200"].schema["x-ms-client-name"] = "AssistantsPage";

      //
      // threads
      //

      const threadsBase = '/threads/{threadId}';

      // GETs
      $[threadsBase + "/messages"].get.responses["200"].schema["x-ms-client-name"] = "ThreadMessagesPage";
      $[threadsBase + "/runs"].get.responses["200"].schema["x-ms-client-name"] = "ThreadRunsPage";
      $[threadsBase + "/runs/{runId}/steps"].get.responses["200"].schema["x-ms-client-name"] = "ThreadRunStepsPage";

      // POSTs
      $[threadsBase + "/messages/{messageId}"].post.parameters[2].schema["x-ms-client-name"] = "UpdateMessageBody";
      $[threadsBase + "/runs/{runId}"].post.parameters[2].schema["x-ms-client-name"] = "UpdateRunBody";
      $[threadsBase + "/runs"].post.parameters[2]["x-ms-client-name"] = "CreateRunBody";
      $[threadsBase + "/runs/{runId}/submit_tool_outputs"].post.parameters[2].schema["x-ms-client-name"] = "SubmitToolOutputsToRunBody";

      //
      // vector stores
      //
      const vectorStoresBase = '/vector_stores';

      // GETs
      $[vectorStoresBase].get.responses["200"].schema["x-ms-client-name"] = "VectorStoresPage";
      $[vectorStoresBase + "/{vectorStoreId}/file_batches/{batchId}/files"].get.responses["200"].schema["x-ms-client-name"] = "VectorStoreFileBatchesPage";
      $[vectorStoresBase + "/{vectorStoreId}/files"].get.responses["200"].schema["x-ms-client-name"] = "VectorStoreFilesPage";

      // POSTs
      $[vectorStoresBase + "/{vectorStoreId}/file_batches"].post.parameters[1].schema["x-ms-client-name"] = "CreateVectorStoreFileBatchBody";

      return $;
```


## Docs

Fix docs for []byte fields that are intended as JSON bytes.


### AssistantsAPIResponseFormatJSONSchemaJSONSchema.Schema
```yaml
directive:
  - from: models.go
    where: $
    transform: |
      const comment = `	// NOTE: this field is JSON text that describes a JSON schema. You can marshal a data\n` +
        ` // structure using code similar to this:\n` +
        ` //\n` +
        `	//	jsonBytes, err := json.Marshal(map[string]any{\n` +
        `	//		"required": []string{"location"},\n` +
        `	// 		"type":     "object",\n` +
        `	// 		"properties": map[string]any{\n` +
        `	// 			"location": map[string]any{\n` +
        `	//	 			"type":        "string",\n` +
        `	// 				"description": "The city and state, e.g. San Francisco, CA",\n` +
        `	// 			},\n` +
        `	//		},\n` +
        `	//	})\n` +
        `	//\n` +
        `	//	if err != nil {\n` +
        `	// 		panic(err)\n` +
        `	//	}\n` +
        `	// \n` +
        `	//	funcDef := &azopenaiassistants.AssistantsAPIResponseFormatJSONSchemaJSONSchema{\n` +
        `	// 		Name:        to.Ptr("get_current_weather"),\n` +
        `	// 		Description: to.Ptr("Get the current weather in a given location"),\n` +
        `	// 		Schema:      jsonBytes,\n` +
        ` //      Strict:      to.Ptr(false),\n` +
        `	// 	}`;
      return $.replace(/Schema \[\]byte/, comment + "\nSchema []byte");
```

### FunctionDefinition.Schema
```yaml
directive:
  - from: models.go
    where: $
    transform: |
      const comment = `	// NOTE: this field is JSON text that describes a JSON schema. You can marshal a data\n` +
        ` // structure using code similar to this:\n` +
        ` //\n` +
        `	//	jsonBytes, err := json.Marshal(map[string]any{\n` +
        `	//		"required": []string{"location"},\n` +
        `	// 		"type":     "object",\n` +
        `	// 		"properties": map[string]any{\n` +
        `	// 			"location": map[string]any{\n` +
        `	//	 			"type":        "string",\n` +
        `	// 				"description": "The city and state, e.g. San Francisco, CA",\n` +
        `	// 			},\n` +
        `	//		},\n` +
        `	//	})\n` +
        `	//\n` +
        `	//	if err != nil {\n` +
        `	// 		panic(err)\n` +
        `	//	}\n` +
        `	// \n` +
        `	//	funcDef := &azopenaiassistants.FunctionDefinition{\n` +
        `	// 		Name:        to.Ptr("get_current_weather"),\n` +
        `	// 		Description: to.Ptr("Get the current weather in a given location"),\n` +
        `	// 		Parameters:  jsonBytes,\n` +
        `	// 	}`;
      return $.replace(/Parameters \[\]byte/, comment + "\nParameters []byte");
```

