//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

func (client *Client) uploadFileCreateRequest(ctx context.Context, file io.ReadSeeker, purpose FilePurpose, options *UploadFileOptions) (*policy.Request, error) {
	urlPath := client.formatURL("/files")
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	fileName := ""

	if options != nil && options.Filename != nil {
		fileName = *options.Filename
	}

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	if err := writeMultipart(req, fileBytes, fileName, purpose); err != nil {
		return nil, err
	}

	return req, nil
}

// GetFileContentResponse contains the response from the [Client.GetFileContent] function.
type GetFileContentResponse struct {
	// Content is the content of the file that's been downloaded.
	// NOTE: this must be Close()'d to avoid leaking resources.
	Content io.ReadCloser
}

// GetFileContentOptions contains the options for the [Client.GetFileContent] function.
type GetFileContentOptions struct {
	// For future expansion
}

// GetFileContent - Returns content for a specific file.
// If the operation fails it returns an *azcore.ResponseError type.
//
//   - fileID - The ID of the file to retrieve.
//   - options - GetFileContentOptions contains the optional parameters for the Client.GetFileContent method.
func (client *Client) GetFileContent(ctx context.Context, fileID string, _ *GetFileContentOptions) (GetFileContentResponse, error) {
	var err error

	req, err := func() (*policy.Request, error) {
		urlPath := client.formatURL("/files/{fileId}/content")
		if fileID == "" {
			return nil, errors.New("parameter fileID cannot be empty")
		}
		urlPath = strings.ReplaceAll(urlPath, "{fileId}", url.PathEscape(fileID))
		req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
		if err != nil {
			return nil, err
		}
		req.Raw().Header["Accept"] = []string{"application/octet-stream"}
		return req, nil
	}()

	if err != nil {
		return GetFileContentResponse{}, err
	}

	runtime.SkipBodyDownload(req)
	//nolint:bodyclose	// caller is responsible for closing this request after they stream the contents of the body.
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetFileContentResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GetFileContentResponse{}, err
	}

	return GetFileContentResponse{Content: httpResp.Body}, nil
}

func writeMultipart(req *policy.Request, fileContents []byte, filename string, purpose FilePurpose) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := createFormFile(writer, "file", filename)

	if err != nil {
		return err
	}

	if _, err := fileWriter.Write(fileContents); err != nil {
		return err
	}

	if err := writer.WriteField("purpose", string(purpose)); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return req.SetBody(streaming.NopCloser(bytes.NewReader(body.Bytes())), writer.FormDataContentType())
}

func createFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteReplacer.Replace(fieldname), quoteReplacer.Replace(filename)))

	contentType := openAIMimeTypes[filepath.Ext(filename)]

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

var openAIMimeTypes = map[string]string{
	".c":    "text/x-c",
	".cpp":  "text/x-c++",
	".csv":  "application/csv",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".html": "text/html",
	".java": "text/x-java",
	".json": "application/json",
	".md":   "text/markdown",
	".pdf":  "application/pdf",
	".php":  "text/x-php",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".py":   "text/x-python",
	".rb":   "text/x-ruby",
	".tex":  "text/x-tex",
	".txt":  "text/plain",
	".css":  "text/css",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript",
	".gif":  "image/gif",
	".png":  "image/png",
	".tar":  "application/x-tar",
	".ts":   "application/typescript",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xml":  "application/xml",
	".zip":  "application/z",
}

var quoteReplacer = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
