// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"testing"
)

type mockAddressResolver struct {
	addresses []AddressInformation
	err       error
}

func (m *mockAddressResolver) ResolveAsync(ctx context.Context, request *DocumentServiceRequest, forceRefresh bool) ([]AddressInformation, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.addresses, nil
}

func TestGetPrimaryUri_NoAddress(t *testing.T) {
	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	replicaAddresses := []AddressInformation{}

	_, err := GetPrimaryUri(request, replicaAddresses)
	if err == nil {
		t.Fatal("expected GoneException, got nil")
	}

	if _, ok := err.(*GoneException); !ok {
		t.Fatalf("expected GoneException, got %T", err)
	}
}

func TestGetPrimaryUri_NoPrimaryAddress(t *testing.T) {
	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos2", ProtocolHTTPS),
	}

	_, err := GetPrimaryUri(request, replicaAddresses)
	if err == nil {
		t.Fatal("expected GoneException, got nil")
	}

	goneErr, ok := err.(*GoneException)
	if !ok {
		t.Fatalf("expected GoneException, got %T", err)
	}

	expectedMsg := "The requested resource is no longer available at the server. Returned addresses are {https://cosmos1/,https://cosmos2/}"
	if goneErr.Error() != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, goneErr.Error())
	}
}

func TestGetPrimaryUri(t *testing.T) {
	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	result, err := GetPrimaryUri(request, replicaAddresses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := UriCreate("https://cosmos2/")
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetPrimaryUri_WithRequestReplicaIndex(t *testing.T) {
	replicaIndex := 1
	request := &DocumentServiceRequest{DefaultReplicaIndex: &replicaIndex}
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	result, err := GetPrimaryUri(request, replicaAddresses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := UriCreate("https://cosmos2/")
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestResolvePrimaryUriAsync(t *testing.T) {
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos4", ProtocolTCP),
		NewAddressInformation(true, true, "https://cosmos5", ProtocolTCP),
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	addressResolver := &mockAddressResolver{addresses: replicaAddresses}
	selector := NewAddressSelector(addressResolver, ProtocolHTTPS)

	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	result, err := selector.ResolvePrimaryUriAsync(context.Background(), request, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := UriCreate("https://cosmos2/")
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestResolveAllUriAsync(t *testing.T) {
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos4", ProtocolTCP),
		NewAddressInformation(true, true, "https://cosmos5", ProtocolTCP),
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	addressResolver := &mockAddressResolver{addresses: replicaAddresses}
	selector := NewAddressSelector(addressResolver, ProtocolHTTPS)

	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	result, err := selector.ResolveAllUriAsync(context.Background(), request, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []Uri{
		UriCreate("https://cosmos1/"),
		UriCreate("https://cosmos2/"),
		UriCreate("https://cosmos3/"),
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d URIs, got %d", len(expected), len(result))
	}

	for i, uri := range result {
		if !uri.Equal(expected[i]) {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], uri)
		}
	}
}

func TestResolveAddressesAsync(t *testing.T) {
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos4", ProtocolTCP),
		NewAddressInformation(true, true, "https://cosmos5", ProtocolTCP),
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	addressResolver := &mockAddressResolver{addresses: replicaAddresses}
	selector := NewAddressSelector(addressResolver, ProtocolHTTPS)

	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	result, err := selector.ResolveAddressesAsync(context.Background(), request, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedAddresses := []AddressInformation{
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	if len(result) != len(expectedAddresses) {
		t.Fatalf("expected %d addresses, got %d", len(expectedAddresses), len(result))
	}

	for i, addr := range result {
		if !addr.Equal(expectedAddresses[i]) {
			t.Errorf("at index %d: expected %v, got %v", i, expectedAddresses[i], addr)
		}
	}
}

func TestResolveAllUriAsync_RNTBD(t *testing.T) {
	replicaAddresses := []AddressInformation{
		NewAddressInformation(true, false, "rntbd://cosmos1", ProtocolTCP),
		NewAddressInformation(true, true, "rntbd://cosmos2", ProtocolTCP),
		NewAddressInformation(true, false, "https://cosmos1", ProtocolHTTPS),
		NewAddressInformation(true, true, "https://cosmos2", ProtocolHTTPS),
		NewAddressInformation(true, false, "https://cosmos3", ProtocolHTTPS),
	}

	addressResolver := &mockAddressResolver{addresses: replicaAddresses}
	selector := NewAddressSelector(addressResolver, ProtocolTCP)

	request := &DocumentServiceRequest{DefaultReplicaIndex: nil}
	result, err := selector.ResolveAllUriAsync(context.Background(), request, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []Uri{
		UriCreate("rntbd://cosmos1/"),
		UriCreate("rntbd://cosmos2/"),
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d URIs, got %d", len(expected), len(result))
	}

	for i, uri := range result {
		if !uri.Equal(expected[i]) {
			t.Errorf("at index %d: expected %v, got %v", i, expected[i], uri)
		}
	}
}
