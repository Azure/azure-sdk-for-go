//go:build ignore

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// tracePolicy injects httptrace hooks so we can see the 100-continue handshake.
type tracePolicy struct{}

func (t *tracePolicy) Do(req *policy.Request) (*http.Response, error) {
	raw := req.Raw()
	expect := raw.Header.Get("Expect")
	fmt.Printf("\n--- REQUEST ---\n")
	fmt.Printf("> %s %s HTTP/1.1\n", raw.Method, raw.URL.Path)
	fmt.Printf("> Host: %s\n", raw.URL.Host)
	fmt.Printf("> Content-Length: %d\n", raw.ContentLength)
	if expect != "" {
		fmt.Printf("> Expect: %s\n", expect)
	}
	fmt.Printf(">\n")

	start := time.Now()

	trace := &httptrace.ClientTrace{
		WroteHeaders: func() {
			fmt.Printf("  [%s] Headers sent, waiting for server...\n", time.Since(start).Round(time.Millisecond))
		},
		Got100Continue: func() {
			fmt.Printf("< HTTP/1.1 100 Continue\n")
			fmt.Printf("  [%s] Got 100 Continue — now sending body\n", time.Since(start).Round(time.Millisecond))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			fmt.Printf("  [%s] Full request body sent\n", time.Since(start).Round(time.Millisecond))
		},
		TLSHandshakeStart: func() {},
		TLSHandshakeDone:  func(state tls.ConnectionState, err error) {},
	}
	ctx := httptrace.WithClientTrace(raw.Context(), trace)
	raw = raw.WithContext(ctx)
	req.Raw().Header = raw.Header

	// Replace the request context
	*req.Raw() = *raw

	resp, err := req.Next()
	if err == nil {
		fmt.Printf("< HTTP/1.1 %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		fmt.Printf("  [%s] Response received\n", time.Since(start).Round(time.Millisecond))
	} else {
		fmt.Printf("  [%s] Error: %v\n", time.Since(start).Round(time.Millisecond), err)
	}
	fmt.Printf("--- END ---\n")
	return resp, err
}

func main() {
	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if connStr == "" {
		fmt.Fprintln(os.Stderr, "ERROR: Set AZURE_STORAGE_CONNECTION_STRING")
		os.Exit(1)
	}

	containerName := "expect100test"
	cntClient, _ := container.NewClientFromConnectionString(connStr, containerName, nil)
	cntClient.Create(context.Background(), nil)

	opts := &blockblob.ClientOptions{}
	opts.PerRetryPolicies = []policy.Policy{&tracePolicy{}}

	client, err := blockblob.NewClientFromConnectionString(connStr, containerName, "test-blob", opts)
	if err != nil {
		panic(err)
	}

	// 9 MiB — above threshold
	fmt.Println("========== Upload 9 MiB (SHOULD trigger 100-continue) ==========")
	body := make([]byte, 9*1024*1024)
	rand.Read(body)
	_, err = client.Upload(context.Background(), &readSeekCloser{bytes.NewReader(body)}, nil)
	if err != nil {
		fmt.Printf("Upload error: %v\n", err)
	}

	// 4 MiB — below threshold
	fmt.Println("\n========== Upload 4 MiB (should NOT trigger 100-continue) ==========")
	small := make([]byte, 4*1024*1024)
	rand.Read(small)
	_, err = client.Upload(context.Background(), &readSeekCloser{bytes.NewReader(small)}, nil)
	if err != nil {
		fmt.Printf("Upload error: %v\n", err)
	}
}

type readSeekCloser struct{ *bytes.Reader }

func (r *readSeekCloser) Close() error { return nil }
