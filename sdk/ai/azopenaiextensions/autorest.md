# Go

These settings apply only when `--go` is specified on the command line.

``` yaml
input-file:
# this file is generated using the ./testdata/genopenapi.ps1 file.
- ./testdata/generated/openapi.json
output-folder: ../azopenaiextensions
clear-output-folder: false
module: github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions
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

``` yaml
directive:
  - from: swagger-document
    where: $["x-ms-paths"]
    transform: |
      return {};
    # NOTE: this is where we decide what models to keep. Anything not included in here just gets
    # removed from the swagger definition.
  - from: swagger-document
    where: $
    transform: |
      const newDefs = {};
      const newPaths = {};

      // add types here if they're Azure related, and we want to keep them and
      // they're not covered by the oydModelRegex below.
      const keep = {};

      // this'll catch the Azure "on your data" models.
      const oydModelRegex = /^(OnYour|Azure|Pinecone|ContentFilter).+$/;

      for (const key in $.definitions) {
        if (!(key in keep) && !key.match(oydModelRegex)) {
          continue
        }

        $lib.log(`Including ${key}`);
        newDefs[key] = $.definitions[key];
      }

      $.definitions = newDefs;

      // clear out any operations, we aren't going to use them.
      $.paths = {};
      $.parameters = {};

      return $;
  - from: swagger-document
    debug: true
    where: $.definitions
    transform: |
      $["Azure.Core.Foundations.Error"]["x-ms-client-name"] = "Error"; 
      delete $["Azure.Core.Foundations.Error"].properties["innererror"];
      delete $["Azure.Core.Foundations.Error"].properties["details"];
      delete $["Azure.Core.Foundations.Error"].properties["target"];

      $["Azure.Core.Foundations.InnerError"]["x-ms-external"] = true;
      $["Azure.Core.Foundations.ErrorResponse"]["x-ms-external"] = true; 
      return $;
```
