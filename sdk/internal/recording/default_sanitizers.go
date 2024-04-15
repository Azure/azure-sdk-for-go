// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// this file contains a set of default sanitizers applied to all recordings

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// sanitizer represents a single sanitizer configured via the test proxy's /Admin/AddSanitizers endpoint
type sanitizer struct {
	// Name is the name of a sanitizer type e.g. "BodyKeySanitizer"
	Name string        `json:"Name,omitempty"`
	Body sanitizerBody `json:"Body,omitempty"`
}

type sanitizerBody struct {
	// GroupForReplace is the name of the regex group to replace
	GroupForReplace string `json:"groupForReplace,omitempty"`
	// JSONPath is the JSON path to the value to replace
	JSONPath string `json:"jsonPath,omitempty"`
	// Key is the name of a header to sanitize
	Key string `json:"key,omitempty"`
	// Regex is the regular expression to match a value to sanitize
	Regex string `json:"regex,omitempty"`
	// Value is the string that replaces the matched value. The sanitizers in
	// this file accept the test proxy's default Value, "Sanitized".
	Value string `json:"value,omitempty"`
}

func newBodyKeySanitizer(jsonPath string) sanitizer {
	return sanitizer{
		Name: "BodyKeySanitizer",
		Body: sanitizerBody{
			JSONPath: jsonPath,
		},
	}
}

func newBodyRegexSanitizer(regex, groupForReplace string) sanitizer {
	return sanitizer{
		Name: "BodyRegexSanitizer",
		Body: sanitizerBody{
			GroupForReplace: groupForReplace,
			Regex:           regex,
		},
	}
}

func newGeneralRegexSanitizer(regex, groupForReplace string) sanitizer {
	return sanitizer{
		Name: "GeneralRegexSanitizer",
		Body: sanitizerBody{
			GroupForReplace: groupForReplace,
			Regex:           regex,
		},
	}
}

func newHeaderRegexSanitizer(key, regex, groupForReplace string) sanitizer {
	return sanitizer{
		Name: "HeaderRegexSanitizer",
		Body: sanitizerBody{
			GroupForReplace: groupForReplace,
			Key:             key,
			Regex:           regex,
		},
	}
}

// addSanitizers adds an arbitrary number of sanitizers with a single request. It
// isn't exported because SDK modules don't add enough sanitizers to benefit from it.
func addSanitizers(s []sanitizer, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizers", options.baseURL())
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	handleTestLevelSanitizer(req, options)
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewReader(b))
	req.ContentLength = int64(len(b))
	req.Header.Set("Content-Type", "application/json")
	return handleProxyResponse(client.Do(req))
}

var defaultSanitizers = []sanitizer{
	newGeneralRegexSanitizer(`("|;)Secret=(?<secret>[^;]+)`, "secret"),
	newBodyKeySanitizer("$..refresh_token"),
	newHeaderRegexSanitizer("api-key", "", ""),
	newBodyKeySanitizer("$..access_token"),
	newBodyKeySanitizer("$..connectionString"),
	newBodyKeySanitizer("$..applicationSecret"),
	newBodyKeySanitizer("$..apiKey"),
	newBodyRegexSanitizer(`client_secret=(?<secret>[^&"]+)`, "secret"),
	newBodyRegexSanitizer(`client_assertion=(?<secret>[^&"]+)`, "secret"),
	newHeaderRegexSanitizer("x-ms-rename-source", "", ""),
	newHeaderRegexSanitizer("x-ms-file-rename-source-authorization", "", ""),
	newHeaderRegexSanitizer("x-ms-file-rename-source", "", ""),
	newHeaderRegexSanitizer("x-ms-encryption-key-sha256", "", ""),
	newHeaderRegexSanitizer("x-ms-encryption-key", "", ""),
	newHeaderRegexSanitizer("x-ms-copy-source-authorization", "", ""),
	newHeaderRegexSanitizer("x-ms-copy-source", "", ""),
	newBodyRegexSanitizer("token=(?<token>[^&]+)($|&)", "token"),
	newHeaderRegexSanitizer("subscription-key", "", ""),
	newBodyKeySanitizer("$..sshPassword"),
	newBodyKeySanitizer("$..secondaryKey"),
	newBodyKeySanitizer("$..runAsPassword"),
	newBodyKeySanitizer("$..primaryKey"),
	newHeaderRegexSanitizer("Location", "", ""),
	newGeneralRegexSanitizer(`("|;)[Aa]ccess[Kk]ey=(?<secret>[^;]+)`, "secret"),
	newGeneralRegexSanitizer(`("|;)[Aa]ccount[Kk]ey=(?<secret>[^;]+)`, "secret"),
	newBodyKeySanitizer("$..aliasSecondaryConnectionString"),
	newGeneralRegexSanitizer(`("|;)[Ss]hared[Aa]ccess[Kk]ey=(?<secret>[^;\"]+)`, "secret"),
	newHeaderRegexSanitizer("aeg-sas-token", "", ""),
	newHeaderRegexSanitizer("aeg-sas-key", "", ""),
	newHeaderRegexSanitizer("aeg-channel-name", "", ""),
	newBodyKeySanitizer("$..adminPassword"),
	newBodyKeySanitizer("$..administratorLoginPassword"),
	newBodyKeySanitizer("$..accessToken"),
	newBodyKeySanitizer("$..accessSAS"),
	newGeneralRegexSanitizer(`(?:(sv|sig|se|srt|ss|sp)=)(?<secret>[^&\"]+)`, "secret"), // SAS tokens
	newBodyKeySanitizer("$.value[*].key"),
	newBodyKeySanitizer("$.key"),
	newBodyKeySanitizer("$..userId"),
	newBodyKeySanitizer("$..urlSource"),
	newBodyKeySanitizer("$..uploadUrl"),
	newBodyKeySanitizer("$..token"),
	newBodyKeySanitizer("$..to"),
	newBodyKeySanitizer("$..tenantId"),
	newBodyKeySanitizer("$..targetResourceId"),
	newBodyKeySanitizer("$..targetModelLocation"),
	newBodyKeySanitizer("$..storageContainerWriteSas"),
	newBodyKeySanitizer("$..storageContainerUri"),
	newBodyKeySanitizer("$..storageContainerReadListSas"),
	newBodyKeySanitizer("$..storageAccountPrimaryKey"),
	newBodyKeySanitizer("$..storageAccount"),
	newBodyKeySanitizer("$..source"),
	newBodyKeySanitizer("$..secondaryReadonlyMasterKey"),
	newBodyKeySanitizer("$..secondaryMasterKey"),
	newBodyKeySanitizer("$..secondaryConnectionString"),
	newBodyKeySanitizer("$..scriptUrlSasToken"),
	newBodyKeySanitizer("$..scan"),
	newBodyKeySanitizer("$..sasUri"),
	newBodyKeySanitizer("$..resourceGroup"),
	newBodyKeySanitizer("$..privateKey"),
	newBodyKeySanitizer("$..principalId"),
	newBodyKeySanitizer("$..primaryReadonlyMasterKey"),
	newBodyKeySanitizer("$..primaryMasterKey"),
	newBodyKeySanitizer("$..primaryConnectionString"),
	newBodyKeySanitizer("$..password"),
	newBodyKeySanitizer("$..outputDataUri"),
	newBodyKeySanitizer("$..managedResourceGroupName"),
	newBodyKeySanitizer("$..logLink"),
	newBodyKeySanitizer("$..lastModifiedBy"),
	newBodyKeySanitizer("$..keyVaultClientSecret"),
	newBodyKeySanitizer("$..inputDataUri"),
	newBodyKeySanitizer("$..id"),
	newBodyKeySanitizer("$..httpHeader"),
	newBodyKeySanitizer("$..guardian"),
	newBodyKeySanitizer("$..functionKey"),
	newBodyKeySanitizer("$..from"),
	newBodyKeySanitizer("$..fencingClientPassword"),
	newBodyKeySanitizer("$..encryptedCredential"),
	newBodyKeySanitizer("$..credential"),
	newBodyKeySanitizer("$..createdBy"),
	newBodyKeySanitizer("$..containerUri"),
	newBodyKeySanitizer("$..clientSecret"),
	newBodyKeySanitizer("$..certificatePassword"),
	newBodyKeySanitizer("$..catalog"),
	newBodyKeySanitizer("$..azureBlobSource.containerUrl"),
	newBodyKeySanitizer("$..authHeader"),
	newBodyKeySanitizer("$..atlasKafkaSecondaryEndpoint"),
	newBodyKeySanitizer("$..atlasKafkaPrimaryEndpoint"),
	newBodyKeySanitizer("$..appkey"),
	newBodyKeySanitizer("$..appId"),
	newBodyKeySanitizer("$..acrToken"),
	newBodyKeySanitizer("$..accountKey"),
	newBodyKeySanitizer("$..AccessToken"),
	newBodyKeySanitizer("$..WEBSITE_AUTH_ENCRYPTION_KEY"),
	newBodyRegexSanitizer("-----BEGIN PRIVATE KEY-----\\\\n(?<key>.+\\\\n)*-----END PRIVATE KEY-----\\\\n", "key"),
	newBodyKeySanitizer("$..adminPassword.value"),
	newBodyKeySanitizer("$..decryptionKey"),
}
