// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azopenai

/**
 * This file is required for the generated code to be valid. The difference between the files
 * with a custom_ prefix and ones like this with an _extra suffix is that the _extra files are
 * not modified by the customization scripts, so they can be run safely. Files with the custom_
 * would change if they weren't ignored, so we have to keep them separate.
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

type clientData struct {
	endpoint string
	azure    bool
}

func (client *Client) newError(resp *http.Response) error {
	return newContentFilterResponseError(resp)
}

func newContentFilterResponseError(resp *http.Response) error {
	//nolint this error is an azcore.ResponseError by definition.
	respErr := runtime.NewResponseError(resp).(*azcore.ResponseError)

	if respErr.ErrorCode != "content_filter" {
		return respErr
	}

	body, err := runtime.Payload(resp)

	if err != nil {
		return err
	}

	var envelope *struct {
		Error struct {
			InnerError struct {
				ContentFilterResults *ContentFilterResults `json:"content_filter_result"`
			} `json:"innererror"`
		}
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}

	return &ContentFilterResponseError{
		ResponseError:        *respErr,
		ContentFilterResults: envelope.Error.InnerError.ContentFilterResults,
	}
}

func (client *Client) formatURL(path string, deployment *string) string {
	switch path {
	// https://learn.microsoft.com/en-us/azure/cognitive-services/openai/reference#image-generation
	case "/images/generations:submit":
		return runtime.JoinPaths(client.endpoint, path)
	default:
		if client.azure {
			if deployment != nil {
				escapedDeplID := url.PathEscape(*deployment)
				return runtime.JoinPaths(client.endpoint, "openai", "deployments", escapedDeplID, path)
			} else {
				return runtime.JoinPaths(client.endpoint, "openai", path)
			}
		}

		return runtime.JoinPaths(client.endpoint, path)
	}
}

// deserializeAudioTranscription handles deserializing the content if it's text/plain
// or a JSON object.
func deserializeAudioTranscription(resp *http.Response) (AudioTranscription, error) {
	defer func() {
		_ = resp.Request.Body.Close()
	}()

	contentType := resp.Header.Get("Content-type")

	if strings.Contains(contentType, "text/plain") {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return AudioTranscription{}, err
		}

		return AudioTranscription{
			Text: to.Ptr(string(body)),
		}, nil
	}

	var result *AudioTranscription
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return AudioTranscription{}, err
	}

	return *result, nil
}

func getAudioTranscriptionInternalHandleResponse(resp *http.Response) (getAudioTranscriptionInternalResponse, error) {
	at, err := deserializeAudioTranscription(resp)

	if err != nil {
		return getAudioTranscriptionInternalResponse{}, err
	}

	return getAudioTranscriptionInternalResponse{AudioTranscription: at}, nil
}

// deserializeAudioTranslation handles deserializing the content if it's text/plain
// or a JSON object.
func deserializeAudioTranslation(resp *http.Response) (AudioTranslation, error) {
	defer func() {
		_ = resp.Request.Body.Close()
	}()

	contentType := resp.Header.Get("Content-type")

	if strings.Contains(contentType, "text/plain") {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return AudioTranslation{}, err
		}

		return AudioTranslation{
			Text: to.Ptr(string(body)),
		}, nil
	}

	var result *AudioTranslation
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return AudioTranslation{}, err
	}

	return *result, nil
}

func getAudioTranslationInternalHandleResponse(resp *http.Response) (getAudioTranslationInternalResponse, error) {
	at, err := deserializeAudioTranslation(resp)

	if err != nil {
		return getAudioTranslationInternalResponse{}, err
	}

	return getAudioTranslationInternalResponse{AudioTranslation: at}, nil
}

func writeField[T interface {
	string | float32 | AudioTranscriptionFormat | AudioTranslationFormat
}](writer *multipart.Writer, fieldName string, v *T) error {
	if v == nil {
		return nil
	}

	switch v2 := any(v).(type) {
	case *string:
		return writer.WriteField(fieldName, *v2)
	case *float32:
		return writer.WriteField(fieldName, fmt.Sprintf("%f", *v2))
	case *AudioTranscriptionFormat:
		return writer.WriteField(fieldName, string(*v2))
	case *AudioTranslationFormat:
		return writer.WriteField(fieldName, string(*v2))
	default:
		return fmt.Errorf("no handler for type %T", v)
	}
}

func setMultipartFormData[T getAudioTranscriptionInternalOptions | getAudioTranslationInternalOptions | UploadFileOptions](req *policy.Request, file io.ReadSeekCloser, options T) error {
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)

	writeContent := func(fieldname, filename string, file io.ReadSeekCloser) error {
		fd, err := writer.CreateFormFile(fieldname, filename)

		if err != nil {
			return err
		}

		if _, err := io.Copy(fd, file); err != nil {
			return err
		}

		return err
	}

	var filename = "audio.mp3"

	switch opt := any(options).(type) {
	case getAudioTranscriptionInternalOptions:
		if opt.Filename != nil {
			filename = *opt.Filename
		}
	case getAudioTranslationInternalOptions:
		if opt.Filename != nil {
			filename = *opt.Filename
		}
	}

	if err := writeContent("file", filename, file); err != nil {
		return err
	}

	switch v := any(options).(type) {
	case getAudioTranslationInternalOptions:
		if err := writeField(writer, "model", v.DeploymentName); err != nil {
			return err
		}
		if err := writeField(writer, "prompt", v.Prompt); err != nil {
			return err
		}
		if err := writeField(writer, "response_format", v.ResponseFormat); err != nil {
			return err
		}
		if err := writeField(writer, "temperature", v.Temperature); err != nil {
			return err
		}
	case getAudioTranscriptionInternalOptions:
		if err := writeField(writer, "language", v.Language); err != nil {
			return err
		}
		if err := writeField(writer, "model", v.DeploymentName); err != nil {
			return err
		}
		if err := writeField(writer, "prompt", v.Prompt); err != nil {
			return err
		}
		if err := writeField(writer, "response_format", v.ResponseFormat); err != nil {
			return err
		}
		if err := writeField(writer, "temperature", v.Temperature); err != nil {
			return err
		}
	default:
		return fmt.Errorf("failed to serialize multipart for unhandled type %T", body)
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return req.SetBody(streaming.NopCloser(bytes.NewReader(body.Bytes())), writer.FormDataContentType())
}
