//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"errors"
	"fmt"
	"regexp"
)

type transformFileOptions struct {
	AllowNoop bool // causes transformFile to fail if the text is not changed after running.
}

type replacer func(text string) (string, error)

func transformFiles(fileCache *FileCache, fileNames []string, replacer replacer, options *transformFileOptions) error {
	for _, fileName := range fileNames {
		origText, err := fileCache.LoadFile(fileName)

		if err != nil {
			return err
		}

		newText, err := replacer(origText)

		if err != nil {
			return err
		}

		if options != nil && options.AllowNoop && newText == origText {
			return errors.New("no replacements were made")
		}

		fileCache.UpdateFile(fileName, newText)
	}

	return nil
}

func removeType(fileCache *FileCache, typeName string) error {
	re := regexp.MustCompile(fmt.Sprintf(`(?s)type %s struct \{.+?\n\}`, typeName))

	err := transformFiles(fileCache, []string{"models.go"}, func(text string) (string, error) {
		return re.ReplaceAllString(text, ""), nil
	}, nil)

	if err != nil {
		return err
	}

	snipMarshallerRE := regexp.MustCompile(fmt.Sprintf(`(?s)// MarshalJSON implements the json.Marshaller interface for type %s.+?\n\}`, typeName))
	snipUnmarshallerRE := regexp.MustCompile(fmt.Sprintf(`(?s)// UnmarshalJSON implements the json.Unmarshaller interface for type %s.+?\n}`, typeName))

	return transformFiles(fileCache, []string{"models_serde.go"}, func(text string) (string, error) {
		text = snipMarshallerRE.ReplaceAllString(text, "")
		text = snipUnmarshallerRE.ReplaceAllString(text, "")
		return text, nil
	}, nil)
}

type updateFunctionOptions struct {
	IgnoreComment bool
}

func updateFunction(text string, objectName string, funcName string, replacer replacer, options *updateFunctionOptions) (string, error) {
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

func removeFunction(text string, objectName string, funcName string) (string, error) {
	return updateFunction(text, objectName, funcName, func(text string) (string, error) {
		return "", nil
	}, nil)
}
