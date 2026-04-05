// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"fmt"
	"strings"
)

type IAddressResolver interface {
	ResolveAsync(ctx context.Context, request *DocumentServiceRequest, forceRefresh bool) ([]AddressInformation, error)
}

type DocumentServiceRequest struct {
	DefaultReplicaIndex *int
}

func (r *DocumentServiceRequest) GetDefaultReplicaIndex() *int {
	return r.DefaultReplicaIndex
}

type AddressSelector struct {
	addressResolver IAddressResolver
	protocol        Protocol
}

func NewAddressSelector(addressResolver IAddressResolver, protocol Protocol) *AddressSelector {
	return &AddressSelector{
		addressResolver: addressResolver,
		protocol:        protocol,
	}
}

func (s *AddressSelector) ResolveAllUriAsync(ctx context.Context, request *DocumentServiceRequest, includePrimary bool, forceRefresh bool) ([]Uri, error) {
	allReplicaAddresses, err := s.ResolveAddressesAsync(ctx, request, forceRefresh)
	if err != nil {
		return nil, err
	}

	var result []Uri
	for _, addr := range allReplicaAddresses {
		if includePrimary || !addr.IsPrimary() {
			result = append(result, addr.GetPhysicalUri())
		}
	}
	return result, nil
}

func (s *AddressSelector) ResolvePrimaryUriAsync(ctx context.Context, request *DocumentServiceRequest, forceAddressRefresh bool) (Uri, error) {
	replicaAddresses, err := s.ResolveAddressesAsync(ctx, request, forceAddressRefresh)
	if err != nil {
		return Uri{}, err
	}

	return GetPrimaryUri(request, replicaAddresses)
}

func GetPrimaryUri(request *DocumentServiceRequest, replicaAddresses []AddressInformation) (Uri, error) {
	var primaryAddress *AddressInformation

	if request.GetDefaultReplicaIndex() != nil {
		defaultReplicaIndex := *request.GetDefaultReplicaIndex()
		if defaultReplicaIndex >= 0 && defaultReplicaIndex < len(replicaAddresses) {
			addr := replicaAddresses[defaultReplicaIndex]
			primaryAddress = &addr
		}
	} else {
		for i := range replicaAddresses {
			addr := &replicaAddresses[i]
			if addr.IsPrimary() && !strings.Contains(addr.GetPhysicalUri().GetURIAsString(), "[") {
				primaryAddress = addr
				break
			}
		}
	}

	if primaryAddress == nil {
		var addresses []string
		for _, addr := range replicaAddresses {
			addresses = append(addresses, addr.GetPhysicalUri().GetURIAsString())
		}
		return Uri{}, NewGoneException(fmt.Sprintf(
			"The requested resource is no longer available at the server. Returned addresses are {%s}",
			strings.Join(addresses, ",")))
	}

	return primaryAddress.GetPhysicalUri(), nil
}

func (s *AddressSelector) ResolveAddressesAsync(ctx context.Context, request *DocumentServiceRequest, forceAddressRefresh bool) ([]AddressInformation, error) {
	addresses, err := s.addressResolver.ResolveAsync(ctx, request, forceAddressRefresh)
	if err != nil {
		return nil, err
	}

	var filtered []AddressInformation
	for _, addr := range addresses {
		uriStr := addr.GetPhysicalUri().GetURIAsString()
		if uriStr != "" && strings.EqualFold(addr.GetProtocolScheme(), s.protocol.Scheme()) {
			filtered = append(filtered, addr)
		}
	}

	var privateAddrs []AddressInformation
	var publicAddrs []AddressInformation
	for _, addr := range filtered {
		if !addr.IsPublic() {
			privateAddrs = append(privateAddrs, addr)
		} else {
			publicAddrs = append(publicAddrs, addr)
		}
	}

	if len(privateAddrs) > 0 {
		return privateAddrs, nil
	}
	return publicAddrs, nil
}

type GoneException struct {
	message string
}

func NewGoneException(message string) *GoneException {
	return &GoneException{message: message}
}

func (e *GoneException) Error() string {
	return e.message
}
