//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armvideoanalyzer

// AudioEncoderBaseClassification provides polymorphic access to related types.
// Call the interface's GetAudioEncoderBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AudioEncoderAac, *AudioEncoderBase
type AudioEncoderBaseClassification interface {
	// GetAudioEncoderBase returns the AudioEncoderBase content of the underlying type.
	GetAudioEncoderBase() *AudioEncoderBase
}

// AuthenticationBaseClassification provides polymorphic access to related types.
// Call the interface's GetAuthenticationBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AuthenticationBase, *JwtAuthentication
type AuthenticationBaseClassification interface {
	// GetAuthenticationBase returns the AuthenticationBase content of the underlying type.
	GetAuthenticationBase() *AuthenticationBase
}

// CertificateSourceClassification provides polymorphic access to related types.
// Call the interface's GetCertificateSource() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *CertificateSource, *PemCertificateList
type CertificateSourceClassification interface {
	// GetCertificateSource returns the CertificateSource content of the underlying type.
	GetCertificateSource() *CertificateSource
}

// CredentialsBaseClassification provides polymorphic access to related types.
// Call the interface's GetCredentialsBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *CredentialsBase, *UsernamePasswordCredentials
type CredentialsBaseClassification interface {
	// GetCredentialsBase returns the CredentialsBase content of the underlying type.
	GetCredentialsBase() *CredentialsBase
}

// EncoderPresetBaseClassification provides polymorphic access to related types.
// Call the interface's GetEncoderPresetBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EncoderCustomPreset, *EncoderPresetBase, *EncoderSystemPreset
type EncoderPresetBaseClassification interface {
	// GetEncoderPresetBase returns the EncoderPresetBase content of the underlying type.
	GetEncoderPresetBase() *EncoderPresetBase
}

// EndpointBaseClassification provides polymorphic access to related types.
// Call the interface's GetEndpointBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EndpointBase, *TLSEndpoint, *UnsecuredEndpoint
type EndpointBaseClassification interface {
	// GetEndpointBase returns the EndpointBase content of the underlying type.
	GetEndpointBase() *EndpointBase
}

// NodeBaseClassification provides polymorphic access to related types.
// Call the interface's GetNodeBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EncoderProcessor, *NodeBase, *ProcessorNodeBase, *RtspSource, *SinkNodeBase, *SourceNodeBase, *VideoSink, *VideoSource
type NodeBaseClassification interface {
	// GetNodeBase returns the NodeBase content of the underlying type.
	GetNodeBase() *NodeBase
}

// ProcessorNodeBaseClassification provides polymorphic access to related types.
// Call the interface's GetProcessorNodeBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EncoderProcessor, *ProcessorNodeBase
type ProcessorNodeBaseClassification interface {
	NodeBaseClassification
	// GetProcessorNodeBase returns the ProcessorNodeBase content of the underlying type.
	GetProcessorNodeBase() *ProcessorNodeBase
}

// SinkNodeBaseClassification provides polymorphic access to related types.
// Call the interface's GetSinkNodeBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *SinkNodeBase, *VideoSink
type SinkNodeBaseClassification interface {
	NodeBaseClassification
	// GetSinkNodeBase returns the SinkNodeBase content of the underlying type.
	GetSinkNodeBase() *SinkNodeBase
}

// SourceNodeBaseClassification provides polymorphic access to related types.
// Call the interface's GetSourceNodeBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *RtspSource, *SourceNodeBase, *VideoSource
type SourceNodeBaseClassification interface {
	NodeBaseClassification
	// GetSourceNodeBase returns the SourceNodeBase content of the underlying type.
	GetSourceNodeBase() *SourceNodeBase
}

// TimeSequenceBaseClassification provides polymorphic access to related types.
// Call the interface's GetTimeSequenceBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *TimeSequenceBase, *VideoSequenceAbsoluteTimeMarkers
type TimeSequenceBaseClassification interface {
	// GetTimeSequenceBase returns the TimeSequenceBase content of the underlying type.
	GetTimeSequenceBase() *TimeSequenceBase
}

// TokenKeyClassification provides polymorphic access to related types.
// Call the interface's GetTokenKey() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *EccTokenKey, *RsaTokenKey, *TokenKey
type TokenKeyClassification interface {
	// GetTokenKey returns the TokenKey content of the underlying type.
	GetTokenKey() *TokenKey
}

// TunnelBaseClassification provides polymorphic access to related types.
// Call the interface's GetTunnelBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *SecureIotDeviceRemoteTunnel, *TunnelBase
type TunnelBaseClassification interface {
	// GetTunnelBase returns the TunnelBase content of the underlying type.
	GetTunnelBase() *TunnelBase
}

// VideoEncoderBaseClassification provides polymorphic access to related types.
// Call the interface's GetVideoEncoderBase() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *VideoEncoderBase, *VideoEncoderH264
type VideoEncoderBaseClassification interface {
	// GetVideoEncoderBase returns the VideoEncoderBase content of the underlying type.
	GetVideoEncoderBase() *VideoEncoderBase
}

