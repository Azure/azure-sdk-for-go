//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

func (client *Client) uploadFileCreateRequest(ctx context.Context, file io.ReadSeeker, purpose FilePurpose, options *UploadFileOptions) (*policy.Request, error) {
	urlPath := client.formatURL("/files", nil)
	req, err := runtime.NewRequest(ctx, http.MethodPost, urlPath)
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
