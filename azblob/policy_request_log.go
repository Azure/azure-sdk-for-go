package azblob

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// RequestLogOptions configures the retry policy's behavior.
type RequestLogOptions struct {
	// LogWarningIfTryOverThreshold logs a warning if a tried operation takes longer than the specified duration (0=no logging).
	LogWarningIfTryOverThreshold time.Duration
}

// NewRequestLogPolicyFactory creates a RequestLogPolicyFactory object configured using the specified options.
func NewRequestLogPolicyFactory(o RequestLogOptions) pipeline.Factory {
	if o.LogWarningIfTryOverThreshold == 0 {
		// It would be good to relate this to https://azure.microsoft.com/en-us/support/legal/sla/storage/v1_2/
		// But this monitors the time to get the HTTP response; NOT the time to download the response body.
		o.LogWarningIfTryOverThreshold = 2 * time.Second // Default to 2 seconds
	}
	return &requestLogPolicyFactory{o: o}
}

type requestLogPolicyFactory struct {
	o RequestLogOptions
}

func (f *requestLogPolicyFactory) New(node pipeline.Node) pipeline.Policy {
	return &requestLogPolicy{node: node, o: f.o}
}

type requestLogPolicy struct {
	node pipeline.Node
	o    RequestLogOptions
	try  int
}

func (p *requestLogPolicy) Do(ctx context.Context, request pipeline.Request) (response pipeline.Response, err error) {
	p.try++ // The first try is #1 (not #0)
	operationStart := time.Now()

	// Log the outgoing request as informational
	if p.node.WouldLog(pipeline.LogInfo) {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "===== OUTGOING REQUEST (Try=%d):\n", p.try)
		pipeline.WriteRequest(b, request.Request)
		p.node.Log(pipeline.LogInfo, b.String())
	}

	// Set the time for this particular retry operation and then Do the operation.
	tryStart := time.Now()
	response, err = p.node.Do(ctx, request) // Make the request
	tryEnd := time.Now()
	tryDuration := tryEnd.Sub(tryStart)
	opDuration := tryEnd.Sub(operationStart)

	severity := pipeline.LogInfo // Assume success and default to informational logging
	logMsg := func(b *bytes.Buffer) {
		b.WriteString("SUCCESS:\n")
		pipeline.WriteResponseWithRequest(b, response.Response())
	}

	// If the response took too long, we'll upgrade to warning.
	if p.o.LogWarningIfTryOverThreshold != 0 && tryDuration > p.o.LogWarningIfTryOverThreshold {
		// Log a warning if the try duration exceeded the specified threshold
		severity = pipeline.LogWarning
		logMsg = func(b *bytes.Buffer) {
			fmt.Fprintf(b, "SLOW [tryDuration > %v]:\n", p.o.LogWarningIfTryOverThreshold)
			pipeline.WriteResponseWithRequest(b, response.Response())
		}
	}

	if err != nil {
		if serr, ok := err.(StorageError); ok && serr.Response() != nil {
			// This Storage Error did get a HTTP response from the service, so we won't change the severity
			logMsg = func(b *bytes.Buffer) {
				fmt.Fprintf(b, "OPERATION ERROR:\n%v\n", serr)
			}
		} else {
			// This error did not get an HTTP response from the service; upgrade the severity to Error
			severity = pipeline.LogError
			logMsg = func(b *bytes.Buffer) {
				// Write the error, the originating requestm and the stack
				fmt.Fprintf(b, "NETWORK ERROR:\n%v\n", err)
				pipeline.WriteRequest(b, request.Request)
				b.Write(stack()) // For errors, we append the stack trace (an expensive operation)
			}
		}
	}

	if p.node.WouldLog(severity) || true { // true is for testing
		// We're going to log this; build the string to log
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "RESPONSE (Try=%d, TryDuration=%v, OpDuration=%v) -- ", p.try, tryDuration, opDuration)
		logMsg(b)
		p.node.Log(severity, b.String())
	}
	return response, err
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
