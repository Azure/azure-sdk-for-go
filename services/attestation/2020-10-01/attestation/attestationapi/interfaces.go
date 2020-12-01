package attestationapi

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/attestation/2020-10-01/attestation"
)

// PolicyClientAPI contains the set of methods on the PolicyClient type.
type PolicyClientAPI interface {
	Get(ctx context.Context, instanceURL string, attestationType attestation.Type) (result attestation.PolicyResponse, err error)
	Reset(ctx context.Context, instanceURL string, attestationType attestation.Type, policyJws string) (result attestation.PolicyResponse, err error)
	Set(ctx context.Context, instanceURL string, attestationType attestation.Type, newAttestationPolicy string) (result attestation.PolicyResponse, err error)
}

var _ PolicyClientAPI = (*attestation.PolicyClient)(nil)

// PolicyCertificatesClientAPI contains the set of methods on the PolicyCertificatesClient type.
type PolicyCertificatesClientAPI interface {
	Add(ctx context.Context, instanceURL string, policyCertificateToAdd string) (result attestation.PolicyCertificatesModifyResponse, err error)
	Get(ctx context.Context, instanceURL string) (result attestation.PolicyCertificatesResponse, err error)
	Remove(ctx context.Context, instanceURL string, policyCertificateToRemove string) (result attestation.PolicyCertificatesModifyResponse, err error)
}

var _ PolicyCertificatesClientAPI = (*attestation.PolicyCertificatesClient)(nil)

// ClientAPI contains the set of methods on the Client type.
type ClientAPI interface {
	AttestOpenEnclave(ctx context.Context, instanceURL string, request attestation.AttestOpenEnclaveRequest) (result attestation.Response, err error)
	AttestSgxEnclave(ctx context.Context, instanceURL string, request attestation.AttestSgxEnclaveRequest) (result attestation.Response, err error)
	AttestTpm(ctx context.Context, instanceURL string, request attestation.TpmAttestationRequest) (result attestation.TpmAttestationResponse, err error)
}

var _ ClientAPI = (*attestation.Client)(nil)

// SigningCertificatesClientAPI contains the set of methods on the SigningCertificatesClient type.
type SigningCertificatesClientAPI interface {
	Get(ctx context.Context, instanceURL string) (result attestation.JSONWebKeySet, err error)
}

var _ SigningCertificatesClientAPI = (*attestation.SigningCertificatesClient)(nil)

// MetadataConfigurationClientAPI contains the set of methods on the MetadataConfigurationClient type.
type MetadataConfigurationClientAPI interface {
	Get(ctx context.Context, instanceURL string) (result attestation.SetObject, err error)
}

var _ MetadataConfigurationClientAPI = (*attestation.MetadataConfigurationClient)(nil)
