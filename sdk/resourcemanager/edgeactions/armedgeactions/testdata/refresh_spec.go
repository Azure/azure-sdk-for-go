//go:build ignore

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Refresh the response-only OpenAPI snapshot. Run from the module with go generate.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	if err := refresh(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func refresh() error {
	location, err := os.ReadFile("tsp-location.yaml")
	if err != nil {
		return err
	}
	source := map[string]string{}
	for _, key := range []string{"repo", "commit", "directory"} {
		match := regexp.MustCompile(`(?m)^` + key + `:\s*(\S+)\s*$`).FindSubmatch(location)
		if len(match) != 2 {
			return fmt.Errorf("missing %s in tsp-location.yaml", key)
		}
		source[key] = strings.Trim(string(match[1]), `"'`)
	}
	if !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(source["commit"]) {
		return fmt.Errorf("tsp-location.yaml must pin a full commit")
	}
	root := "https://raw.githubusercontent.com/" + source["repo"] + "/" + source["commit"] + "/"
	specURL := root + source["directory"] + "/stable/2026-10-01/openapi.json"
	source["url"] = specURL
	client := &http.Client{Timeout: 60 * time.Second}
	documents := map[string]map[string]any{}
	load := func(address string) (map[string]any, error) {
		if !strings.HasPrefix(address, root) {
			return nil, fmt.Errorf("schema reference leaves the pinned repository")
		}
		if doc, ok := documents[address]; ok {
			return doc, nil
		}
		response, err := client.Get(address)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return nil, fmt.Errorf("schema download returned HTTP %d", response.StatusCode)
		}
		var doc map[string]any
		if err = json.NewDecoder(response.Body).Decode(&doc); err != nil {
			return nil, err
		}
		documents[address] = doc
		return doc, nil
	}
	definitions := map[string]any{}
	var bundle func(any, string) (any, error)
	bundle = func(value any, documentURL string) (any, error) {
		switch value := value.(type) {
		case []any:
			for i := range value {
				item, err := bundle(value[i], documentURL)
				if err != nil {
					return nil, err
				}
				value[i] = item
			}
		case map[string]any:
			for key, item := range value {
				if key != "$ref" {
					converted, err := bundle(item, documentURL)
					if err != nil {
						return nil, err
					}
					value[key] = converted
					continue
				}
				reference, ok := item.(string)
				if !ok {
					return nil, fmt.Errorf("invalid schema reference")
				}
				base, err := url.Parse(documentURL)
				if err != nil {
					return nil, err
				}
				relative, err := url.Parse(reference)
				if err != nil {
					return nil, err
				}
				target := base.ResolveReference(relative)
				pointer := target.Fragment
				target.Fragment = ""
				if !strings.HasPrefix(pointer, "/definitions/") {
					return nil, fmt.Errorf("unsupported schema reference %s", pointer)
				}
				hash := sha256.Sum256([]byte(target.String() + "#" + pointer))
				name := fmt.Sprintf("%x_%s", hash[:8], pointer[strings.LastIndex(pointer, "/")+1:])
				value[key] = "#/definitions/" + name
				if _, ok := definitions[name]; ok {
					continue
				}
				doc, err := load(target.String())
				if err != nil {
					return nil, err
				}
				var definition any = doc
				for _, part := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
					object, ok := definition.(map[string]any)
					if !ok {
						return nil, fmt.Errorf("invalid reference pointer %s", pointer)
					}
					definition, ok = object[strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")]
					if !ok {
						return nil, fmt.Errorf("missing reference pointer %s", pointer)
					}
				}
				// Clone before rewriting references: other paths may use this cached document.
				data, err := json.Marshal(definition)
				if err != nil {
					return nil, err
				}
				if err = json.Unmarshal(data, &definition); err != nil {
					return nil, err
				}
				definitions[name] = nil
				definitions[name], err = bundle(definition, target.String())
				if err != nil {
					return nil, err
				}
			}
		}
		return value, nil
	}
	spec, err := load(specURL)
	if err != nil {
		return err
	}
	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		return fmt.Errorf("OpenAPI paths missing")
	}
	for _, methods := range paths {
		operations, ok := methods.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid OpenAPI path")
		}
		for method, value := range operations {
			op, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("invalid OpenAPI operation")
			}
			if _, ok := op["operationId"]; !ok {
				delete(operations, method)
				continue
			}
			responses, ok := op["responses"].(map[string]any)
			if !ok {
				return fmt.Errorf("operation has no responses")
			}
			responsesBundle, err := bundle(responses, specURL)
			if err != nil {
				return err
			}
			operations[method] = map[string]any{
				"operationId": op["operationId"], "responses": responsesBundle,
				"x-ms-long-running-operation": op["x-ms-long-running-operation"] == true,
			}
		}
	}
	data, err := json.MarshalIndent(map[string]any{"source": source, "paths": paths, "definitions": definitions}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("testdata/response_spec.json", append(data, '\n'), 0644)
}
