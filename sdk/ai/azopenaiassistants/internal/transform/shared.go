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

type transformFileOptions struct {
	AllowNoop bool // causes transformFile to fail if the text is not changed after running.
}

type replacer func(text string) (string, error)

func transformFiles(fileCache *FileCache, purpose string, fileNames []string, replacer replacer, options *transformFileOptions) error {
	if options == nil {
		options = &transformFileOptions{}
	}

	replaced := false

	for _, fileName := range fileNames {
		origText, err := fileCache.LoadFile(fileName)

		if err != nil {
			return err
		}

		newText, err := replacer(origText)

		if err != nil {
			return err
		}

		if newText != origText {
			replaced = true
		}

		fileCache.UpdateFile(fileName, newText)
	}

	if !replaced && !options.AllowNoop {
		return fmt.Errorf("(%s) no replacements were made in files %#v", purpose, fileNames)
	}

	return nil
}

type removeTypesOptions struct {
	IgnoreComment bool
}

func removeTypes(fileCache *FileCache, typeNames []string, options *removeTypesOptions) error {
	if options == nil {
		options = &removeTypesOptions{}
	}

	for _, typeName := range typeNames {
		purpose := fmt.Sprintf("Removing type %s", typeName)
		log.Println(purpose)

		reText := fmt.Sprintf(`type %s struct \{.+?\n\}`, typeName)

		if !options.IgnoreComment {
			reText = fmt.Sprintf(`// %s.+?`, typeName) + reText
		}

		reText = "(?s)" + reText

		re := regexp.MustCompile(reText)

		err := transformFiles(fileCache, purpose, []string{"models.go", "responses.go", "options.go"}, func(text string) (string, error) {
			return re.ReplaceAllString(text, ""), nil
		}, nil)

		if err != nil {
			return err
		}

		if strings.HasSuffix(typeName, "Response") || strings.HasSuffix(typeName, "Options") {
			// only model types have actual serde functions to remove.
			continue
		}

		snipMarshallerRE := regexp.MustCompile(fmt.Sprintf(`(?s)// MarshalJSON implements the json.Marshaller interface for type %s.+?\n\}`, typeName))
		snipUnmarshallerRE := regexp.MustCompile(fmt.Sprintf(`(?s)// UnmarshalJSON implements the json.Unmarshaller interface for type %s.+?\n}`, typeName))

		err = transformFiles(fileCache, purpose, []string{"models_serde.go"}, func(text string) (string, error) {
			text = snipMarshallerRE.ReplaceAllString(text, "")
			text = snipUnmarshallerRE.ReplaceAllString(text, "")
			return text, nil
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

type updateFunctionOptions struct {
	IgnoreComment bool
}

func updateFunction(text string, objectName string, funcName string, replacer replacer, options *updateFunctionOptions) (string, error) {
	log.Printf("Updating function %s.%s", objectName, funcName)
	return updateFunctionImpl(text, objectName, funcName, replacer, options)
}

func removeFunctions(text string, objectName string, funcNames ...string) (string, error) {
	for _, funcName := range funcNames {
		log.Printf("Removing function %s.%s", objectName, funcName)

		var err error
		text, err = updateFunctionImpl(text, objectName, funcName, func(_ string) (string, error) {
			return "", nil
		}, nil)

		if err != nil {
			return "", err
		}
	}

	return text, nil
}

func updateFunctionImpl(text string, objectName string, funcName string, replacer replacer, options *updateFunctionOptions) (string, error) {
	// ex: func (client *Client) uploadFileCreateRequest(ctx context.Context, file string, purpose FilePurpose, options *ClientUploadFileOptions) (*policy.Request, error) {
	regexpText := fmt.Sprintf(`func \([^ ]+\s+\*%s\) %s\(.+?\n}`, objectName, funcName)

	if options == nil || !options.IgnoreComment {
		regexpText = fmt.Sprintf("// %s .+?", funcName) + regexpText
	}

	re := regexp.MustCompile("(?s)" + regexpText)
	funcText := re.FindString(text)

	if funcText == "" {
		return "", fmt.Errorf("no match for object %s, function name %s", objectName, funcName)
	}

	newFuncText, err := replacer(funcText)

	if err != nil {
		return "", err
	}

	newText := re.ReplaceAllString(text, newFuncText)
	return newText, nil
}
