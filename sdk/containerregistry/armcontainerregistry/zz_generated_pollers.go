// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcontainerregistry

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"time"
)

// AgentPoolPoller provides polling facilities until the operation reaches a terminal state.
type AgentPoolPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final AgentPoolResponse will be returned.
	FinalResponse(ctx context.Context) (AgentPoolResponse, error)
}

type agentPoolPoller struct {
	pt *armcore.LROPoller
}

func (p *agentPoolPoller) Done() bool {
	return p.pt.Done()
}

func (p *agentPoolPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *agentPoolPoller) FinalResponse(ctx context.Context) (AgentPoolResponse, error) {
	respType := AgentPoolResponse{AgentPool: &AgentPool{}}
	resp, err := p.pt.FinalResponse(ctx, respType.AgentPool)
	if err != nil {
		return AgentPoolResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *agentPoolPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *agentPoolPoller) pollUntilDone(ctx context.Context, freq time.Duration) (AgentPoolResponse, error) {
	respType := AgentPoolResponse{AgentPool: &AgentPool{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.AgentPool)
	if err != nil {
		return AgentPoolResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ConnectedRegistryPoller provides polling facilities until the operation reaches a terminal state.
type ConnectedRegistryPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final ConnectedRegistryResponse will be returned.
	FinalResponse(ctx context.Context) (ConnectedRegistryResponse, error)
}

type connectedRegistryPoller struct {
	pt *armcore.LROPoller
}

func (p *connectedRegistryPoller) Done() bool {
	return p.pt.Done()
}

func (p *connectedRegistryPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *connectedRegistryPoller) FinalResponse(ctx context.Context) (ConnectedRegistryResponse, error) {
	respType := ConnectedRegistryResponse{ConnectedRegistry: &ConnectedRegistry{}}
	resp, err := p.pt.FinalResponse(ctx, respType.ConnectedRegistry)
	if err != nil {
		return ConnectedRegistryResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *connectedRegistryPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *connectedRegistryPoller) pollUntilDone(ctx context.Context, freq time.Duration) (ConnectedRegistryResponse, error) {
	respType := ConnectedRegistryResponse{ConnectedRegistry: &ConnectedRegistry{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.ConnectedRegistry)
	if err != nil {
		return ConnectedRegistryResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ExportPipelinePoller provides polling facilities until the operation reaches a terminal state.
type ExportPipelinePoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final ExportPipelineResponse will be returned.
	FinalResponse(ctx context.Context) (ExportPipelineResponse, error)
}

type exportPipelinePoller struct {
	pt *armcore.LROPoller
}

func (p *exportPipelinePoller) Done() bool {
	return p.pt.Done()
}

func (p *exportPipelinePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *exportPipelinePoller) FinalResponse(ctx context.Context) (ExportPipelineResponse, error) {
	respType := ExportPipelineResponse{ExportPipeline: &ExportPipeline{}}
	resp, err := p.pt.FinalResponse(ctx, respType.ExportPipeline)
	if err != nil {
		return ExportPipelineResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *exportPipelinePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *exportPipelinePoller) pollUntilDone(ctx context.Context, freq time.Duration) (ExportPipelineResponse, error) {
	respType := ExportPipelineResponse{ExportPipeline: &ExportPipeline{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.ExportPipeline)
	if err != nil {
		return ExportPipelineResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// GenerateCredentialsResultPoller provides polling facilities until the operation reaches a terminal state.
type GenerateCredentialsResultPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final GenerateCredentialsResultResponse will be returned.
	FinalResponse(ctx context.Context) (GenerateCredentialsResultResponse, error)
}

type generateCredentialsResultPoller struct {
	pt *armcore.LROPoller
}

func (p *generateCredentialsResultPoller) Done() bool {
	return p.pt.Done()
}

func (p *generateCredentialsResultPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *generateCredentialsResultPoller) FinalResponse(ctx context.Context) (GenerateCredentialsResultResponse, error) {
	respType := GenerateCredentialsResultResponse{GenerateCredentialsResult: &GenerateCredentialsResult{}}
	resp, err := p.pt.FinalResponse(ctx, respType.GenerateCredentialsResult)
	if err != nil {
		return GenerateCredentialsResultResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *generateCredentialsResultPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *generateCredentialsResultPoller) pollUntilDone(ctx context.Context, freq time.Duration) (GenerateCredentialsResultResponse, error) {
	respType := GenerateCredentialsResultResponse{GenerateCredentialsResult: &GenerateCredentialsResult{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.GenerateCredentialsResult)
	if err != nil {
		return GenerateCredentialsResultResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// HTTPPoller provides polling facilities until the operation reaches a terminal state.
type HTTPPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final *http.Response will be returned.
	FinalResponse(ctx context.Context) (*http.Response, error)
}

type httpPoller struct {
	pt *armcore.LROPoller
}

func (p *httpPoller) Done() bool {
	return p.pt.Done()
}

func (p *httpPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *httpPoller) FinalResponse(ctx context.Context) (*http.Response, error) {
	return p.pt.FinalResponse(ctx, nil)
}

func (p *httpPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *httpPoller) pollUntilDone(ctx context.Context, freq time.Duration) (*http.Response, error) {
	return p.pt.PollUntilDone(ctx, freq, nil)
}

// ImportPipelinePoller provides polling facilities until the operation reaches a terminal state.
type ImportPipelinePoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final ImportPipelineResponse will be returned.
	FinalResponse(ctx context.Context) (ImportPipelineResponse, error)
}

type importPipelinePoller struct {
	pt *armcore.LROPoller
}

func (p *importPipelinePoller) Done() bool {
	return p.pt.Done()
}

func (p *importPipelinePoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *importPipelinePoller) FinalResponse(ctx context.Context) (ImportPipelineResponse, error) {
	respType := ImportPipelineResponse{ImportPipeline: &ImportPipeline{}}
	resp, err := p.pt.FinalResponse(ctx, respType.ImportPipeline)
	if err != nil {
		return ImportPipelineResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *importPipelinePoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *importPipelinePoller) pollUntilDone(ctx context.Context, freq time.Duration) (ImportPipelineResponse, error) {
	respType := ImportPipelineResponse{ImportPipeline: &ImportPipeline{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.ImportPipeline)
	if err != nil {
		return ImportPipelineResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// PipelineRunPoller provides polling facilities until the operation reaches a terminal state.
type PipelineRunPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final PipelineRunResponseType will be returned.
	FinalResponse(ctx context.Context) (PipelineRunResponseType, error)
}

type pipelineRunPoller struct {
	pt *armcore.LROPoller
}

func (p *pipelineRunPoller) Done() bool {
	return p.pt.Done()
}

func (p *pipelineRunPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *pipelineRunPoller) FinalResponse(ctx context.Context) (PipelineRunResponseType, error) {
	respType := PipelineRunResponseType{PipelineRun: &PipelineRun{}}
	resp, err := p.pt.FinalResponse(ctx, respType.PipelineRun)
	if err != nil {
		return PipelineRunResponseType{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *pipelineRunPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *pipelineRunPoller) pollUntilDone(ctx context.Context, freq time.Duration) (PipelineRunResponseType, error) {
	respType := PipelineRunResponseType{PipelineRun: &PipelineRun{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.PipelineRun)
	if err != nil {
		return PipelineRunResponseType{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// PrivateEndpointConnectionPoller provides polling facilities until the operation reaches a terminal state.
type PrivateEndpointConnectionPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final PrivateEndpointConnectionResponse will be returned.
	FinalResponse(ctx context.Context) (PrivateEndpointConnectionResponse, error)
}

type privateEndpointConnectionPoller struct {
	pt *armcore.LROPoller
}

func (p *privateEndpointConnectionPoller) Done() bool {
	return p.pt.Done()
}

func (p *privateEndpointConnectionPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *privateEndpointConnectionPoller) FinalResponse(ctx context.Context) (PrivateEndpointConnectionResponse, error) {
	respType := PrivateEndpointConnectionResponse{PrivateEndpointConnection: &PrivateEndpointConnection{}}
	resp, err := p.pt.FinalResponse(ctx, respType.PrivateEndpointConnection)
	if err != nil {
		return PrivateEndpointConnectionResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *privateEndpointConnectionPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *privateEndpointConnectionPoller) pollUntilDone(ctx context.Context, freq time.Duration) (PrivateEndpointConnectionResponse, error) {
	respType := PrivateEndpointConnectionResponse{PrivateEndpointConnection: &PrivateEndpointConnection{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.PrivateEndpointConnection)
	if err != nil {
		return PrivateEndpointConnectionResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// RegistryPoller provides polling facilities until the operation reaches a terminal state.
type RegistryPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final RegistryResponse will be returned.
	FinalResponse(ctx context.Context) (RegistryResponse, error)
}

type registryPoller struct {
	pt *armcore.LROPoller
}

func (p *registryPoller) Done() bool {
	return p.pt.Done()
}

func (p *registryPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *registryPoller) FinalResponse(ctx context.Context) (RegistryResponse, error) {
	respType := RegistryResponse{Registry: &Registry{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Registry)
	if err != nil {
		return RegistryResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *registryPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *registryPoller) pollUntilDone(ctx context.Context, freq time.Duration) (RegistryResponse, error) {
	respType := RegistryResponse{Registry: &Registry{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Registry)
	if err != nil {
		return RegistryResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ReplicationPoller provides polling facilities until the operation reaches a terminal state.
type ReplicationPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final ReplicationResponse will be returned.
	FinalResponse(ctx context.Context) (ReplicationResponse, error)
}

type replicationPoller struct {
	pt *armcore.LROPoller
}

func (p *replicationPoller) Done() bool {
	return p.pt.Done()
}

func (p *replicationPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *replicationPoller) FinalResponse(ctx context.Context) (ReplicationResponse, error) {
	respType := ReplicationResponse{Replication: &Replication{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Replication)
	if err != nil {
		return ReplicationResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *replicationPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *replicationPoller) pollUntilDone(ctx context.Context, freq time.Duration) (ReplicationResponse, error) {
	respType := ReplicationResponse{Replication: &Replication{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Replication)
	if err != nil {
		return ReplicationResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// RunPoller provides polling facilities until the operation reaches a terminal state.
type RunPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final RunResponse will be returned.
	FinalResponse(ctx context.Context) (RunResponse, error)
}

type runPoller struct {
	pt *armcore.LROPoller
}

func (p *runPoller) Done() bool {
	return p.pt.Done()
}

func (p *runPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *runPoller) FinalResponse(ctx context.Context) (RunResponse, error) {
	respType := RunResponse{Run: &Run{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Run)
	if err != nil {
		return RunResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *runPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *runPoller) pollUntilDone(ctx context.Context, freq time.Duration) (RunResponse, error) {
	respType := RunResponse{Run: &Run{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Run)
	if err != nil {
		return RunResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// ScopeMapPoller provides polling facilities until the operation reaches a terminal state.
type ScopeMapPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final ScopeMapResponse will be returned.
	FinalResponse(ctx context.Context) (ScopeMapResponse, error)
}

type scopeMapPoller struct {
	pt *armcore.LROPoller
}

func (p *scopeMapPoller) Done() bool {
	return p.pt.Done()
}

func (p *scopeMapPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *scopeMapPoller) FinalResponse(ctx context.Context) (ScopeMapResponse, error) {
	respType := ScopeMapResponse{ScopeMap: &ScopeMap{}}
	resp, err := p.pt.FinalResponse(ctx, respType.ScopeMap)
	if err != nil {
		return ScopeMapResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *scopeMapPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *scopeMapPoller) pollUntilDone(ctx context.Context, freq time.Duration) (ScopeMapResponse, error) {
	respType := ScopeMapResponse{ScopeMap: &ScopeMap{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.ScopeMap)
	if err != nil {
		return ScopeMapResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// TaskPoller provides polling facilities until the operation reaches a terminal state.
type TaskPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final TaskResponse will be returned.
	FinalResponse(ctx context.Context) (TaskResponse, error)
}

type taskPoller struct {
	pt *armcore.LROPoller
}

func (p *taskPoller) Done() bool {
	return p.pt.Done()
}

func (p *taskPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *taskPoller) FinalResponse(ctx context.Context) (TaskResponse, error) {
	respType := TaskResponse{Task: &Task{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Task)
	if err != nil {
		return TaskResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *taskPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *taskPoller) pollUntilDone(ctx context.Context, freq time.Duration) (TaskResponse, error) {
	respType := TaskResponse{Task: &Task{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Task)
	if err != nil {
		return TaskResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// TaskRunPoller provides polling facilities until the operation reaches a terminal state.
type TaskRunPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final TaskRunResponse will be returned.
	FinalResponse(ctx context.Context) (TaskRunResponse, error)
}

type taskRunPoller struct {
	pt *armcore.LROPoller
}

func (p *taskRunPoller) Done() bool {
	return p.pt.Done()
}

func (p *taskRunPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *taskRunPoller) FinalResponse(ctx context.Context) (TaskRunResponse, error) {
	respType := TaskRunResponse{TaskRun: &TaskRun{}}
	resp, err := p.pt.FinalResponse(ctx, respType.TaskRun)
	if err != nil {
		return TaskRunResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *taskRunPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *taskRunPoller) pollUntilDone(ctx context.Context, freq time.Duration) (TaskRunResponse, error) {
	respType := TaskRunResponse{TaskRun: &TaskRun{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.TaskRun)
	if err != nil {
		return TaskRunResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// TokenPoller provides polling facilities until the operation reaches a terminal state.
type TokenPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final TokenResponse will be returned.
	FinalResponse(ctx context.Context) (TokenResponse, error)
}

type tokenPoller struct {
	pt *armcore.LROPoller
}

func (p *tokenPoller) Done() bool {
	return p.pt.Done()
}

func (p *tokenPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *tokenPoller) FinalResponse(ctx context.Context) (TokenResponse, error) {
	respType := TokenResponse{Token: &Token{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Token)
	if err != nil {
		return TokenResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *tokenPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *tokenPoller) pollUntilDone(ctx context.Context, freq time.Duration) (TokenResponse, error) {
	respType := TokenResponse{Token: &Token{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Token)
	if err != nil {
		return TokenResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

// WebhookPoller provides polling facilities until the operation reaches a terminal state.
type WebhookPoller interface {
	azcore.Poller
	// FinalResponse performs a final GET to the service and returns the final response
	// for the polling operation. If there is an error performing the final GET then an error is returned.
	// If the final GET succeeded then the final WebhookResponse will be returned.
	FinalResponse(ctx context.Context) (WebhookResponse, error)
}

type webhookPoller struct {
	pt *armcore.LROPoller
}

func (p *webhookPoller) Done() bool {
	return p.pt.Done()
}

func (p *webhookPoller) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

func (p *webhookPoller) FinalResponse(ctx context.Context) (WebhookResponse, error) {
	respType := WebhookResponse{Webhook: &Webhook{}}
	resp, err := p.pt.FinalResponse(ctx, respType.Webhook)
	if err != nil {
		return WebhookResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}

func (p *webhookPoller) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

func (p *webhookPoller) pollUntilDone(ctx context.Context, freq time.Duration) (WebhookResponse, error) {
	respType := WebhookResponse{Webhook: &Webhook{}}
	resp, err := p.pt.PollUntilDone(ctx, freq, respType.Webhook)
	if err != nil {
		return WebhookResponse{}, err
	}
	respType.RawResponse = resp
	return respType, nil
}
