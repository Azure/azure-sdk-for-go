# Go

These settings apply only when `--go` is specified on the command line.

```yaml
input-file:
  # this file is generated using the ./testdata/genopenapi.ps1 file.
  - ./testdata/generated/openapi.json
output-folder: ../azopenai
clear-output-folder: false
module: github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai
license-header: MICROSOFT_MIT_NO_VERSION
openapi-type: data-plane
go: true
use: "@autorest/go@4.0.0-preview.63"
title: "OpenAI"
slice-elements-byval: true
rawjson-as-bytes: true
# can't use this since it removes an innererror type that we want ()
# remove-non-reference-schema: true
```

## Transformations

Fix deployment and endpoint parameters so they show up in the right spots

```yaml
directive:
  # Add x-ms-parameter-location to parameters in x-ms-parameterized-host
  - from: swagger-document
    where: $["x-ms-parameterized-host"].parameters.0
    transform: $["x-ms-parameter-location"] = "client";
  # Make deploymentId a client parameter
  # This must be done in each operation as the parameter is not defined in the components section
  - from: swagger-document
    where: $.paths..parameters..[?(@.name=='deploymentId')]
    transform: $["x-ms-parameter-location"] = "client";
```

## Model -> DeploymentName

```yaml
directive:
  - from: swagger-document
    where: $.definitions.GenerateSpeechFromTextOptions.properties
    transform: |
      $["model"] = {
          "type": "string",
          "description": "The model to use for this request."
      }
      return $;
  - from:
      - models.go
      - models_serde.go
      - options.go
      - client.go
    where: $
    transform: |
      return $.replace(/Model \*string/g, 'DeploymentName *string')
        .replace(/populate\(objectMap, "model", (.)\.Model\)/g, 'populate(objectMap, "model", $1.DeploymentName)')
        .replace(/err = unpopulate\(val, "Model", &(.)\.Model\)/g, 'err = unpopulate(val, "Model", &$1.DeploymentName)')
        .replace(/Model:/g, "DeploymentName: ");
  # hack - we have _one_ spot where we want to keep it as Model (ChatCompletions.Model) since
  # it is the actual model name, not the deployment, in Azure OpenAI or OpenAI.
  - from: models.go
    where: $
    transform: return $.replace(/(ChatCompletions.+?)DeploymentName/s, "$1Model");
  - from: models_serde.go
    where: $
    transform: |
      return $
        .replace(/(func \(c ChatCompletions\) MarshalJSON.+?populate\(objectMap, "model", c\.)DeploymentName/s, "$1Model")
        .replace(/(func \(c \*ChatCompletions\) UnmarshalJSON.+?unpopulate\(val, "Model", &c\.)DeploymentName/s, "$1Model");
```

## Polymorphic adjustments

The polymorphic _input_ models all expose the discriminator but it's ignored when serializing
(ie, each type already knows the value and fills it in). So we'll just hide it.

`ChatRequestMessageClassification.Role`

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatRequestMessage
    transform: $.properties.role["x-ms-client-name"] = "InternalRoleRename"
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/InternalRoleRename/g, "role")
```

`AzureChatExtensionConfigurationClassification.Type`

```yaml
directive:
  - from: swagger-document
    where: $.definitions.AzureChatExtensionConfiguration
    transform: $.properties.type["x-ms-client-name"] = "InternalChatExtensionTypeRename"
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/InternalChatExtensionTypeRename/g, "configType")
```

`OnYourDataAuthenticationOptionsClassification.Type`

```yaml
directive:
  - from: swagger-document
    where: $.definitions.OnYourDataAuthenticationOptions
    transform: $.properties.type["x-ms-client-name"] = "InternalOYDAuthTypeRename"
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/InternalOYDAuthTypeRename/g, "configType")
```

`ChatCompletionsResponseFormat.Type`

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionsResponseFormat
    transform: $.properties.type["x-ms-client-name"] = "InternalChatCompletionsResponseFormat"
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/InternalChatCompletionsResponseFormat/g, "respType")
```

Fix casing of some fields

```yaml
directive:
  # Filepath -> FilePath
  - from: swagger-document
    where: $.definitions.AzureChatExtensionRetrievedDocument
    transform: $.properties.filepath["x-ms-client-name"] = "FilePath"
  - from: swagger-document
    where: $.definitions.AzureChatExtensionDataSourceResponseCitation
    transform: $.properties.filepath["x-ms-client-name"] = "FilePath"

  # FilepathField -> FilePathField
  - from: swagger-document
    where: $.definitions.AzureCosmosDBFieldMappingOptions
    transform: $.properties.filepath_field["x-ms-client-name"] = "FilePathField"
  - from: swagger-document
    where: $.definitions.AzureSearchIndexFieldMappingOptions
    transform: $.properties.filepath_field["x-ms-client-name"] = "FilePathField"
  - from: swagger-document
    where: $.definitions.ElasticsearchIndexFieldMappingOptions
    transform: $.properties.filepath_field["x-ms-client-name"] = "FilePathField"
  - from: swagger-document
    where: $.definitions.PineconeFieldMappingOptions
    transform: $.properties.filepath_field["x-ms-client-name"] = "FilePathField"

  # CustomBlocklist -> CustomBlockList
  - from: swagger-document
    where: $.definitions.ImageGenerationPromptFilterResults
    transform: $.properties.custom_blocklists["x-ms-client-name"] = "CustomBlockLists"
  - from: swagger-document
    where: $.definitions.ContentFilterResultDetailsForPrompt
    transform: $.properties.custom_blocklists["x-ms-client-name"] = "CustomBlockLists"
  - from: swagger-document
    where: $.definitions.ContentFilterResultsForChoice
    transform: $.properties.custom_blocklists["x-ms-client-name"] = "CustomBlockLists"
```

## Cleanup the audio transcription APIs

We're wrapping the audio translation and transcription APIs, so we can eliminate some of
these autogenerated models and functions.

```yaml
directive:
  # kill models
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: |
      const typesToRemove = [
        'XMSPathsHksgfdDeploymentsDeploymentidAudioTranscriptionsOverloadGetaudiotranscriptionasresponseobjectPostRequestbodyContentMultipartFormDataSchema',
        'XMSPaths1Ak7Ov3DeploymentsDeploymentidAudioTranslationsOverloadGetaudiotranslationasresponseobjectPostRequestbodyContentMultipartFormDataSchema',
        'Paths1G1Yr9HDeploymentsDeploymentidAudioTranslationsPostRequestbodyContentMultipartFormDataSchema',
        'Paths1MlipaDeploymentsDeploymentidAudioTranscriptionsPostRequestbodyContentMultipartFormDataSchema',
        'Paths1Filz8PFilesPostRequestbodyContentMultipartFormDataSchema',
        'Paths46Ul4XUploadsUploadIDPartsPostRequestbodyContentMultipartFormDataSchema',
        'BatchCreateResponse',
      ];

      for (let name of typesToRemove) {
        $ = $.replace(new RegExp(`(// ${name} - [\\w\\s\\.,]+\n)?type ${name} struct.+?\s*\}`, "s"), "")
        .replace(new RegExp(`// MarshalJSON implements the json.Marshaller interface for type ${name}.+?\n}`, "s"), "")
        .replace(new RegExp(`// UnmarshalJSON implements the json.Unmarshaller interface for type ${name}.+?\n}`, "s"), "");
      }
      return $;
  # kill API functions
  - from:
      - client.go
    where: $
    transform: |
      return $.replace(/\/\/ GetAudioTranscriptionAsPlainText -.+?\n\}\n/s, "")
        .replace(/\/\/ GetAudioTranslationAsPlainText -.+?\n\}\n/s, "")
        .replace(/\/\/ getAudioTranscriptionAsPlainTextCreateRequest.+?\n}\n/s, "")
        .replace(/\/\/ getAudioTranscriptionAsPlainTextHandleResponse.+?\n}\n/s, "")
        .replace(/\/\/ getAudioTranslationAsPlainTextCreateRequest.+?\n}\n/s, "")
        .replace(/\/\/ getAudioTranslationAsPlainTextHandleResponse.+?\n}\n/s, "");
  # remove other plain text models/options
  - from:
      - options.go
    where: $
    transform: |
      return $.replace(/\/\/ ClientGetAudioTranslationAsPlainTextOptions .+?\n\}\n/s, "")
      .replace(/\/\/ ClientGetAudioTranscriptionAsPlainTextOptions .+?\n\}\n/s, "");
  # remove other plain text models/options
  - from:
      - responses.go
    where: $
    transform: |
      return $.replace(/\/\/ ClientGetAudioTranscriptionAsPlainTextResponse .+?\n\}\n/s, "")
      .replace(/\/\/ ClientGetAudioTranslationAsPlainTextResponse .+?\n\}\n/s, "");

  # fix any calls that don't use 'deploymentID'
  - from: client.go
    where: $
    transform: |
      return $
        .replace(/, deploymentID string,/g, ",")
        .replace(/ctx, deploymentID, /g, "ctx, ")
```

## Move the Azure extensions into their own section of the options

We've moved these 'extension' data types into their own field.

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionsOptions.properties.data_sources
    transform: $["x-ms-client-name"] = "AzureExtensionsOptions"
```

## Trim the Error object to match our Go conventions

```yaml
directive:
  - from: swagger-document
    where: $.definitions["Azure.Core.Foundations.Error"]
    transform: |
      $.properties = {
        code: $.properties["code"],
        message: {
          ...$.properties["message"],
          "x-ms-client-name": "InternalErrorMessageRename"
        },
      };
      $["x-ms-client-name"] = "Error";

  - from: swagger-document
    where: $.definitions
    transform: delete $["Azure.Core.Foundations.InnerError"];

  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/InternalErrorMessageRename/g, "message");
```

## Splice in some hooks for custom code

```yaml
directive:
  # Allow interception of formatting the URL path
  - from: client.go
    where: $
    transform: |
      const urlPaths = [
          ".+?/audio/speech",
          ".+?/audio/transcriptions",
          ".+?/audio/translations",
          ".+?/chat/completions",
          ".+?/completions",
          ".+?/embeddings",
          ".+?/images/generations"
      ].join("|");

      const re = new RegExp(
        '(urlPath := "(?:' + urlPaths + ')"\\s+' +
        '.+?)runtime\.JoinPaths\\(client\\.endpoint, urlPath\\)',
        'gs');

      return $.replace(re, "$1client.formatURL(urlPath, body.DeploymentName)");
  - from: client.go
    where: $
    transform: return $.replace(/runtime\.JoinPaths\(client\.endpoint, urlPath\)/g, "client.formatURL(urlPath, nil)");
  # Allow custom parsing of the returned error, mostly for handling the content filtering errors.
  - from: client.go
    where: $
    transform: return $.replace(/runtime\.NewResponseError/sg, "client.newError");
```

Other misc fixes

```yaml
directive:
  - from: swagger-document
    where: $..paths["/deployments/{deploymentId}/completions"].post.requestBody
    transform: $["required"] = true;
  - from: swagger-document
    where: $.paths["/deployments/{deploymentId}/embeddings"].post.requestBody
    transform: $["required"] = true;

  # get rid of these auto-generated LRO status methods that aren't exposed.
  - from: swagger-document
    where: $.paths
    transform: delete $["/operations/images/{operationId}"]
```

Changes for audio/whisper APIs.

```yaml
directive:
  # the whisper operations are really long since they are a conglomeration of _all_ the
  # possible return types.
  - rename-operation:
      from: GetAudioTranscriptionAsResponseObject
      to: GetAudioTranscriptionInternal
  - rename-operation:
      from: GetAudioTranslationAsResponseObject
      to: GetAudioTranslationInternal

  - from: swagger-document
    where: $["x-ms-paths"]["/deployments/{deploymentId}/audio/translations?_overload=getAudioTranslationAsResponseObject"].post
    transform: $.operationId = "GetAudioTranslationInternal"
  - from: swagger-document
    where: $["x-ms-paths"]["/deployments/{deploymentId}/audio/transcriptions?_overload=getAudioTranscriptionAsResponseObject"].post
    transform: $.operationId = "GetAudioTranscriptionInternal"

  # hide the generated functions, in favor of our public wrappers.
  - from:
      - client.go
      - models.go
      - models_serde.go
      - responses.go
      - options.go
    where: $
    transform: |
      return $
        .replace(/GetAudioTranscriptionInternal([^){ ]*)/g, "getAudioTranscriptionInternal$1")
        .replace(/GetAudioTranslationInternal([^){ ]*)/g, "getAudioTranslationInternal$1");

  # some multipart fixing
  - from: client.go
    where: $
    transform: |
      return $
        .replace(/(func.* getAudio(?:Translation|Transcription)InternalCreateRequest\(.+?)options/g, "$1body")
        .replace(/(func.* uploadFileCreateRequest\(.+?)options/g, "$1body");
  - from: client.go
    where: $
    transform: |
      return $
        .replace(/((?:getAudioTranscriptionInternalCreateRequest|getAudioTranslationInternalCreateRequest).+?)runtime\.SetMultipartFormData\(.+?\)/sg, "$1setMultipartFormData(req, file, *body)");

  # response type parsing (can be text/plain _or_ JSON)
  - from: client.go
    where: $
    transform: |
      return $
      .replace(/client\.getAudioTranscriptionInternalHandleResponse/g, "getAudioTranscriptionInternalHandleResponse")
        .replace(/client\.getAudioTranslationInternalHandleResponse/g, "getAudioTranslationInternalHandleResponse")

  # fix the file parameter to be a []byte.
  - from: client.go
    where: $
    transform: return $.replace(/^(func \(client \*Client\) getAudioTrans.+?)file string,(.+)$/mg, "$1file []byte,$2")
```

## Logprob casing fixes

```yaml
directive:
  - from:
      - options.go
      - models_serde.go
      - models.go
    where: $
    transform: |
      return $
        .replace(/AvgLogprob \*float32/g, "AvgLogProb *float32")
        .replace(/(a|c)\.AvgLogprob/g, "$1.AvgLogProb")
  - from:
      - client.go
      - models.go
      - models_serde.go
      - options.go
      - responses.go
    where: $
    transform: return $.replace(/Logprobs/g, "LogProbs")
```

```yaml
directive:
  #
  # strip out the deploymentID validation code - we absorbed this into the endpoint.
  #
	# urlPath := "/deployments/{deploymentId}/embeddings"
	# if client.deploymentID == "" {
	# 	return nil, errors.New("parameter client.deploymentID cannot be empty")
	# }
	# urlPath = strings.ReplaceAll(urlPath, "{deploymentId}", url.PathEscape(client.deploymentID))
  - from: client.go
    where: $
    transform: >-
      return $.replace(
        /(\s+)urlPath\s*:=\s*"\/deployments\/\{deploymentId\}\/([^"]+)".+?url\.PathEscape.+?\n/gs,
        "$1urlPath := \"$2\"\n")

  # splice out the auto-generated `deploymentID` field from the client
  - from: client.go
    where: $
    transform: >-
      return $.replace(
        /(type Client struct[^}]+})/s,
        "type Client struct {\ninternal *azcore.Client; clientData;\n}")

  - from:
    - models_serde.go
    - models.go
    where: $
    transform: |
      return $
        // remove some types that were generated to support the recursive error.
        .replace(/\/\/ AzureCoreFoundationsInnerErrorInnererror.+?\n}/s, "")
        // also, remove its marshalling functions
        .replace(/\/\/ (Unmarshal|Marshal)JSON implements[^\n]+?for type AzureCoreFoundationsInnerErrorInnererror.+?\n}/sg, "")
        .replace(/\/\/ AzureCoreFoundationsErrorInnererror.+?\n}/s, "")
        .replace(/\/\/ (Unmarshal|Marshal)JSON implements[^\n]+?for type AzureCoreFoundationsErrorInnererror.+?\n}/sg, "")
        .replace(/\/\/ AzureCoreFoundationsErrorResponseError.+?\n}/s, "")
        .replace(/\/\/ (Unmarshal|Marshal)JSON implements[^\n]+?for type AzureCoreFoundationsErrorResponseError.+?\n}/sg, "")
        .replace(/\/\/ AzureCoreFoundationsErrorResponse.+?\n}/s, "")
        .replace(/\/\/ (Unmarshal|Marshal)JSON implements[^\n]+?for type AzureCoreFoundationsErrorResponse.+?\n}/sg, "")

        // Remove any references to the type and replace them with InnerError.
        .replace(/Innererror \*(AzureCoreFoundationsInnerErrorInnererror|AzureCoreFoundationsErrorInnererror)/g, "InnerError *InnerError")

        // Fix the marshallers/unmarshallers to use the right case.
        .replace(/(a|c).Innererror/g, '$1.InnerError')

        // We have two "inner error" types that are identical (ErrorInnerError and InnerError). Let's eliminate the one that's not actually directly referenced.
        .replace(/\/\/azureCoreFoundationsInnerError.+?\n}/s, "")

        //
        // Fix the AzureCoreFoundation naming to match our style.
        //
        .replace(/AzureCoreFoundations/g, "")
  - from: constants.go
    where: $
    transform: >-
      return $.replace(
        /type ServiceAPIVersions string.+PossibleServiceAPIVersionsValues.+?\n}/gs,
        "")

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - models.go
      - options.go
      - responses.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");

  # Make the Azure extensions internal - we expose these through the GetChatCompletions*() functions
  # and just treat which endpoint we use as an implementation detail.
  # - from: client.go
  #   where: $
  #   transform: |
  #     return $
  #       .replace(/GetChatCompletionsWithAzureExtensions([ (])/g, "getChatCompletionsWithAzureExtensions$1")
  #       .replace(/GetChatCompletions([ (])/g, "getChatCompletions$1");
```

## Workarounds

This handles a case where (depending on mixture of older and newer resources) we can potentially see
_either_ of these fields that represents the same data (prompt filter results).

```yaml
directive:
  - from: models_serde.go
    where: $
    transform: return $.replace(/case "prompt_filter_results":/g, 'case "prompt_annotations":\nfallthrough\ncase "prompt_filter_results":')
```

Add in some types that are incorrectly not being exported in the generation

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatCompletionRequestMessageContentPartType"] = {
        "type": "string",
        "enum": [
          "image_url",
          "text"
        ],
        "description": "The type of the content part.",
        "x-ms-enum": {
          "name": "ChatCompletionRequestMessageContentPartType",
          "modelAsString": true,
          "values": [
            {
              "name": "image_url",
              "value": "image_url",
              "description": "Chat content contains an image URL"
            },
            {
              "name": "text",
              "value": "text",
              "description": "Chat content contains text"
            },
          ]
        }
      };
      $["ChatCompletionRequestMessageContentPart"] = {
        "title": "represents either an image URL or text content for a prompt",
        "type": "object",
        "discriminator": "type",
        "properties": {
          "type": {
            "$ref": "#/definitions/ChatCompletionRequestMessageContentPartType"
          }
        },
        "required": [
          "type"
        ],
      };
      $["ChatCompletionRequestMessageContentPartImage"] = {
        "type": "object",
        "title": "represents an image URL, to be used as part of a prompt",
        "properties": {
          "image_url": {
            "type": "object",
            "title": "contains the URL and level of detail for an image prompt",
            "properties": {
              "url": {
                "type": "string",
                "description": "Either a URL of the image or the base64 encoded image data.",
                "format": "uri"
              },
              "detail": {
                "type": "string",
                "description": "Specifies the detail level of the image. Learn more in the [Vision guide](/docs/guides/vision/low-or-high-fidelity-image-understanding).",
                "enum": [
                  "auto",
                  "low",
                  "high"
                ],
                "default": "auto"
              }
            },
            "required": [
              "url"
            ]
          }
        },
        "allOf": [
          {
            "$ref": "#/definitions/ChatCompletionRequestMessageContentPart"
          }
        ],
        "required": [
          "image_url"
        ],
        "x-ms-discriminator-value": "image_url"
      };
      $["ChatCompletionRequestMessageContentPartText"] = {
        "type": "object",
        "title": "represents text content, to be used as part of a prompt",
        "properties": {
          "text": {
            "type": "string",
            "description": "The text content."
          }
        },
        "allOf": [
          {
            "$ref": "#/definitions/ChatCompletionRequestMessageContentPart"
          }
        ],
        "required": [
          "text"
        ],
        "x-ms-discriminator-value": "text"
      };
```

Polymorphic removal of the Type field: `ChatCompletionRequestMessageContentPartClassification.Type`

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionRequestMessageContentPart
    transform: $.properties.type["x-ms-client-name"] = "ChatCompletionRequestMessageContentPartTypeRename"
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $.replace(/ChatCompletionRequestMessageContentPartTypeRename/g, "partType")
```

Another workaround - streaming results don't contain the discriminator field so we'll
inject it when we can infer it properly ('function' property exists).

```yaml
directive:
  - from: polymorphic_helpers.go
    where: $
    transform: |
      return $.replace(/(func unmarshalChatCompletionsToolCallClassification.+?var b ChatCompletionsToolCallClassification\n)/s,
        `$1\n` +
        `if m["type"] == nil && m["function"] != nil {\n` +
        `  // WORKAROUND: the streaming results don't contain the proper role for functions, so we need to add these in.\n` +
        `  m["type"] = string(ChatRoleFunction)\n` +
        `}\n`);
```

Embedding has two ways of coming back - base64 or already decoded into floats. This has to be handled manually.

```yaml
directive:
  - from: models_serde.go
    where: $
    transform: return $.replace(/err = unpopulate\(val, "Embedding", &e.Embedding\)/g, "err = deserializeEmbeddingsArray(val, e)");
  - from: models.go
    where: $
    transform: return $.replace(/\/\/ EmbeddingItem - .+?type EmbeddingItem struct \{.+?\n}\n/s, "");
```

Fix some doc comments

```yaml
directive:
  - from: models.go
    where: $
    transform: |
      const text = "// NOTE: This field is not available when using [Client.GetChatCompletionsStream].\n$1";
      return $.replace(/(Usage \*CompletionsUsage)/, text);
  - from: models.go
    where: $
    transform: |
      const text = "// - If using EmbeddingEncodingFormatFloat (the default), the value will be a []float32, in [EmbeddingItem.Embedding]\n" +
        "// - If using EmbeddingEncodingFormatBase64, the value will be a base-64 string in [EmbeddingItem.EmbeddingBase64]\n";

      return $.replace(/(EncodingFormat \*EmbeddingEncodingFormat)/, `${text}$1`);
```

Update docs for FunctionDefinition.Parameters to indicate it's intended to be serialized JSON bytes,
not an object or map[string]any.

```yaml
directive:
  - from: models.go
    where: $
    transform: |
      const comment = `	// REQUIRED; The function definition details for the function tool. \n` +
        `	// NOTE: this field is JSON text that describes a JSON schema. You can marshal a data\n` +
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
        `	//	funcDef := &azopenai.FunctionDefinition{\n` +
        `	// 		Name:        to.Ptr("get_current_weather"),\n` +
        `	// 		Description: to.Ptr("Get the current weather in a given location"),\n` +
        `	// 		Parameters:  jsonBytes,\n` +
        `	// 	}`;
      return $.replace(/Parameters \[\]byte/, comment + "\nParameters []byte");
```

## Unions

Update the ChatRequestUserMessage to allow for []ChatCompletionRequestMessageContentPartText _or_
a string.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatRequestUserMessageContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatRequestUserMessage.properties.content
    transform: $["$ref"] = "#/definitions/ChatRequestUserMessageContent"; return $;
```

Update ChatRequestAssistantMessage.Content to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatRequestAssistantMessageContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatRequestAssistantMessage.properties.content
    transform: $["$ref"] = "#/definitions/ChatRequestAssistantMessageContent"; return $;
```

Update ChatRequestSystemMessage.content to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatRequestSystemMessageContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatRequestSystemMessage.properties.content
    transform: $["$ref"] = "#/definitions/ChatRequestSystemMessageContent"; return $;
```

Update ChatRequestDeveloperMessage.content to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatRequestDeveloperMessageContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatRequestDeveloperMessage.properties.content
    transform: $["$ref"] = "#/definitions/ChatRequestDeveloperMessageContent"; return $;
```

Update ChatRequestToolMessage.content to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatRequestToolMessageContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatRequestToolMessage.properties.content
    transform: $["$ref"] = "#/definitions/ChatRequestToolMessageContent"; return $;
```

Update MongoDBChatExtensionParameters.embedding_dependency to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["MongoDBChatExtensionParametersEmbeddingDependency"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.MongoDBChatExtensionParameters.properties.embedding_dependency
    transform: $["$ref"] = "#/definitions/MongoDBChatExtensionParametersEmbeddingDependency"; return $;
```

Update PredictionContent.content to use its custom type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["PredictionContentContent"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.PredictionContent.properties.content
    transform: $["$ref"] = "#/definitions/PredictionContentContent"; return $;
```

\*ChatCompletionsToolChoice

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatCompletionsToolChoice"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatCompletionsOptions.properties.tool_choice
    transform: $["$ref"] = "#/definitions/ChatCompletionsToolChoice"; return $;
```

\*ChatCompletionsOptionsFunctionCall

```yaml
directive:
  - from: swagger-document
    where: $.definitions
    transform: |
      $["ChatCompletionsOptionsFunctionCall"] = {
        "x-ms-external": true,
        "type": "object", "properties": { "stub": { "type": "string" }}
      };
      return $;
  - from: swagger-document
    where: $.definitions.ChatCompletionsOptions.properties.function_call
    transform: $["$ref"] = "#/definitions/ChatCompletionsOptionsFunctionCall"; return $;
```

<!--
Fix ToolChoice discriminated union

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionsOptions.properties
    transform: $["tool_choice"]["x-ms-client-name"] = "ToolChoiceRenameMe"
  - from:
    - models.go
    - models_serde.go
    where: $
    transform: |
      return $
        .replace(/^\s+ToolChoiceRenameMe.+$/m, "ToolChoice *ChatCompletionsToolChoice")   // update the name _and_ type for the field
        .replace(/ToolChoiceRenameMe/g, "ToolChoice")    // rename all other references
        .replace(/populateAny\(objectMap, "tool_choice", c\.ToolChoice\)/, 'populate(objectMap, "tool_choice", c.ToolChoice)');   // treat field as typed so nil means omit.
``` -->

<!-- ToolChoice and FunctionCall have types so they don't need the 'RawMessage' treatment.

```yaml
directive:
  - from: models.go
    where: $
    transform: return $.replace(/FunctionCall \[\]byte/, "FunctionCall *ChatCompletionsOptionsFunctionCall");
  - from: models_serde.go
    where: $
    transform: return $.replace(/json\.RawMessage\(([a-z]\.(?:FunctionCall|ToolChoice))\)/g, "$1")

``` -->

## Pagers

ListBatches is a pageable API.

```yaml
directive:
  - from: client.go
    where: $
    transform: return $
      .replace(/ListBatches([ (])/g, "listBatches$1")
      .replace(/result\.ListBatchesResponse/g, "result.ListBatchesPage");
  - from: responses.go
    where: $
    transform: return $.replace(/ListBatchesResponse\s+\}/g, "ListBatchesPage\n}");
  - from:
      - models.go
      - models_serde.go
    where: $
    transform: return $
      .replace(/ ListBatchesResponse/g, " ListBatchesPage")
      .replace(/ListBatchesResponse\)/g, " ListBatchesPage)");
```

## Files

```yaml
directive:
  - from: client.go
    where: $
    transform: return $.replace(/\/\/ uploadFileCreateRequest creates .+?return req, nil\s+}/s, "");
```

## Anonymous types

Give names to anonymous types

```yaml
directive:
  - from: swagger-document
    where: $.paths['/batches'].get.responses['200'].schema
    transform: $["x-ms-client-name"] = "ListBatchesPage"; return $;
```

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionsOptions.properties.stream_options
    transform: $["$ref"] = "#/definitions/ChatCompletionStreamOptions"
  - from: swagger-document
    where: $.definitions.CompletionsOptions.properties.stream_options
    transform: $["$ref"] = "#/definitions/ChatCompletionStreamOptions"
```

## Doc updates

Hoisting the description for an anonymous type.

```yaml
directive:
  - from: swagger-document
    where: $.definitions.ChatCompletionsJsonSchemaResponseFormat.properties.json_schema
    transform: $.description = $.properties.description.description; return $;
```
