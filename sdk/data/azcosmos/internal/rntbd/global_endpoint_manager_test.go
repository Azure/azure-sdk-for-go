// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"
)

var (
	dbAccountJson1 = `{"_self":"","id":"testaccount","_rid":"testaccount.documents.azure.com","media":"//media/","addresses":"//addresses/","_dbs":"//dbs/","writableLocations":[{"name":"East US","databaseAccountEndpoint":"https://testaccount-eastus.documents.azure.com:443/"}],"readableLocations":[{"name":"East US","databaseAccountEndpoint":"https://testaccount-eastus.documents.azure.com:443/"},{"name":"East Asia","databaseAccountEndpoint":"https://testaccount-eastasia.documents.azure.com:443/"}],"enableMultipleWriteLocations":false}`

	dbAccountJson2 = `{"_self":"","id":"testaccount","_rid":"testaccount.documents.azure.com","media":"//media/","addresses":"//addresses/","_dbs":"//dbs/","writableLocations":[{"name":"East Asia","databaseAccountEndpoint":"https://testaccount-eastasia.documents.azure.com:443/"}],"readableLocations":[{"name":"East Asia","databaseAccountEndpoint":"https://testaccount-eastasia.documents.azure.com:443/"}],"enableMultipleWriteLocations":false}`

	dbAccountJson3 = `{"_self":"","id":"testaccount","_rid":"testaccount.documents.azure.com","media":"//media/","addresses":"//addresses/","_dbs":"//dbs/","writableLocations":[{"name":"West US","databaseAccountEndpoint":"https://testaccount-westus.documents.azure.com:443/"}],"readableLocations":[{"name":"West US","databaseAccountEndpoint":"https://testaccount-westus.documents.azure.com:443/"}],"enableMultipleWriteLocations":false}`

	dbAccountJson4 = `{"_self":"","id":"testaccount","_rid":"testaccount.documents.azure.com","media":"//media/","addresses":"//addresses/","_dbs":"//dbs/","writableLocations":[{"name":"East US","databaseAccountEndpoint":"https://testaccount-eastus.documents.azure.com:443/"},{"name":"East Asia","databaseAccountEndpoint":"https://testaccount-eastasia.documents.azure.com:443/"}],"readableLocations":[{"name":"East US","databaseAccountEndpoint":"https://testaccount-eastus.documents.azure.com:443/"},{"name":"East Asia","databaseAccountEndpoint":"https://testaccount-eastasia.documents.azure.com:443/"}],"enableMultipleWriteLocations":true}`
)

type mockDatabaseAccountManager struct {
	mu              sync.Mutex
	accounts        []*DatabaseAccount
	currentIndex    int
	serviceEndpoint *url.URL
}

func newMockDatabaseAccountManager(serviceEndpoint string) *mockDatabaseAccountManager {
	endpoint, _ := url.Parse(serviceEndpoint)
	return &mockDatabaseAccountManager{
		serviceEndpoint: endpoint,
		accounts:        []*DatabaseAccount{},
	}
}

func (m *mockDatabaseAccountManager) SetDatabaseAccount(account *DatabaseAccount) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts = []*DatabaseAccount{account}
	m.currentIndex = 0
}

func (m *mockDatabaseAccountManager) SetDatabaseAccounts(accounts ...*DatabaseAccount) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts = accounts
	m.currentIndex = 0
}

func (m *mockDatabaseAccountManager) GetDatabaseAccountFromEndpoint(ctx context.Context, endpoint *url.URL) (*DatabaseAccount, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.accounts) == 0 {
		return nil, nil
	}
	account := m.accounts[m.currentIndex]
	if m.currentIndex < len(m.accounts)-1 {
		m.currentIndex++
	}
	return account, nil
}

func (m *mockDatabaseAccountManager) GetServiceEndpoint() *url.URL {
	return m.serviceEndpoint
}

func TestRefreshLocationAsyncForConnectivityIssue(t *testing.T) {
	manager := newMockDatabaseAccountManager("https://testaccount.documents.azure.com:443")

	account1, _ := NewDatabaseAccountFromJSON(dbAccountJson1)
	manager.SetDatabaseAccount(account1)

	policy := &ConnectionPolicy{
		EnableEndpointDiscovery:     true,
		UsingMultipleWriteLocations: true,
	}

	gem := NewGlobalEndpointManager(manager, policy)
	if err := gem.Init(context.Background()); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	locationCache := gem.GetLocationCache()
	availableWriteByLocation := locationCache.GetAvailableWriteEndpointByLocation()
	availableReadByLocation := locationCache.GetAvailableReadEndpointByLocation()

	if len(availableWriteByLocation) != 1 {
		t.Errorf("expected 1 available write location, got %d", len(availableWriteByLocation))
	}
	if len(availableReadByLocation) != 2 {
		t.Errorf("expected 2 available read locations, got %d", len(availableReadByLocation))
	}
	if _, ok := availableWriteByLocation["East US"]; !ok {
		t.Error("expected East US in available write locations")
	}
	if _, ok := availableReadByLocation["East US"]; !ok {
		t.Error("expected East US in available read locations")
	}
	if _, ok := availableReadByLocation["East Asia"]; !ok {
		t.Error("expected East Asia in available read locations")
	}

	account2, _ := NewDatabaseAccountFromJSON(dbAccountJson2)
	manager.SetDatabaseAccount(account2)

	eastUSEndpoint, _ := url.Parse("https://testaccount-eastus.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForRead(eastUSEndpoint)
	gem.RefreshLocationAsync(context.Background(), false)

	locationCache = gem.GetLocationCache()
	availableReadByLocation = locationCache.GetAvailableReadEndpointByLocation()
	if len(availableReadByLocation) != 1 {
		t.Errorf("expected 1 available read location after refresh, got %d", len(availableReadByLocation))
	}
	if _, ok := availableReadByLocation["East Asia"]; !ok {
		t.Error("expected East Asia in available read locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}

	account3, _ := NewDatabaseAccountFromJSON(dbAccountJson3)
	manager.SetDatabaseAccount(account3)

	eastAsiaEndpoint, _ := url.Parse("https://testaccount-eastasia.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForRead(eastAsiaEndpoint)
	gem.RefreshLocationAsync(context.Background(), false)

	availableReadByLocation = gem.GetLocationCache().GetAvailableReadEndpointByLocation()
	if _, ok := availableReadByLocation["West US"]; !ok {
		t.Error("expected West US in available read locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}
}

func TestRefreshLocationAsyncForConnectivityIssueWithPreferredLocation(t *testing.T) {
	manager := newMockDatabaseAccountManager("https://testaccount.documents.azure.com:443")

	policy := &ConnectionPolicy{
		EnableEndpointDiscovery:     true,
		PreferredLocations:          []string{"East US", "East Asia"},
		UsingMultipleWriteLocations: true,
	}

	account1, _ := NewDatabaseAccountFromJSON(dbAccountJson1)
	manager.SetDatabaseAccount(account1)

	gem := NewGlobalEndpointManager(manager, policy)
	if err := gem.Init(context.Background()); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	account2, _ := NewDatabaseAccountFromJSON(dbAccountJson2)
	manager.SetDatabaseAccount(account2)

	eastUSEndpoint, _ := url.Parse("https://testaccount-eastus.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForRead(eastUSEndpoint)
	gem.RefreshLocationAsync(context.Background(), false)

	locationCache := gem.GetLocationCache()
	if len(locationCache.GetReadEndpoints()) != 2 {
		t.Errorf("expected 2 read endpoints, got %d", len(locationCache.GetReadEndpoints()))
	}

	availableReadByLocation := locationCache.GetAvailableReadEndpointByLocation()
	if len(availableReadByLocation) != 2 {
		t.Errorf("expected 2 available read locations, got %d", len(availableReadByLocation))
	}

	if _, ok := availableReadByLocation["East Asia"]; !ok {
		t.Error("expected East Asia in available read locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}

	account3, _ := NewDatabaseAccountFromJSON(dbAccountJson3)
	manager.SetDatabaseAccount(account3)

	eastAsiaEndpoint, _ := url.Parse("https://testaccount-eastasia.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForRead(eastAsiaEndpoint)
	gem.RefreshLocationAsync(context.Background(), false)

	availableReadByLocation = gem.GetLocationCache().GetAvailableReadEndpointByLocation()
	if _, ok := availableReadByLocation["West US"]; !ok {
		t.Error("expected West US in available read locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}
}

func TestRefreshLocationAsyncForWriteForbidden(t *testing.T) {
	manager := newMockDatabaseAccountManager("https://testaccount.documents.azure.com:443")

	account1, _ := NewDatabaseAccountFromJSON(dbAccountJson1)
	manager.SetDatabaseAccount(account1)

	policy := &ConnectionPolicy{
		EnableEndpointDiscovery:     true,
		UsingMultipleWriteLocations: true,
	}

	gem := NewGlobalEndpointManager(manager, policy)
	if err := gem.Init(context.Background()); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	account2, _ := NewDatabaseAccountFromJSON(dbAccountJson2)
	manager.SetDatabaseAccount(account2)

	eastUSEndpoint, _ := url.Parse("https://testaccount-eastus.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForWrite(eastUSEndpoint)
	gem.RefreshLocationAsync(context.Background(), true)

	locationCache := gem.GetLocationCache()
	if len(locationCache.GetReadEndpoints()) != 1 {
		t.Errorf("expected 1 read endpoint, got %d", len(locationCache.GetReadEndpoints()))
	}

	availableWriteByLocation := locationCache.GetAvailableWriteEndpointByLocation()
	if _, ok := availableWriteByLocation["East Asia"]; !ok {
		t.Error("expected East Asia in available write locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}

	account3, _ := NewDatabaseAccountFromJSON(dbAccountJson3)
	manager.SetDatabaseAccount(account3)

	eastAsiaEndpoint, _ := url.Parse("https://testaccount-eastasia.documents.azure.com:443/")
	gem.MarkEndpointUnavailableForWrite(eastAsiaEndpoint)
	gem.RefreshLocationAsync(context.Background(), true)

	locationCache = gem.GetLocationCache()
	if len(locationCache.GetReadEndpoints()) != 1 {
		t.Errorf("expected 1 read endpoint, got %d", len(locationCache.GetReadEndpoints()))
	}

	availableWriteByLocation = locationCache.GetAvailableWriteEndpointByLocation()
	if _, ok := availableWriteByLocation["West US"]; !ok {
		t.Error("expected West US in available write locations")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}
	if !gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be true")
	}
}

func TestBackgroundRefreshForMultiMaster(t *testing.T) {
	manager := newMockDatabaseAccountManager("https://testaccount.documents.azure.com:443")

	account4, _ := NewDatabaseAccountFromJSON(dbAccountJson4)
	manager.SetDatabaseAccount(account4)

	policy := &ConnectionPolicy{
		EnableEndpointDiscovery:     true,
		UsingMultipleWriteLocations: true,
	}

	gem := NewGlobalEndpointManager(manager, policy)
	if err := gem.Init(context.Background()); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	gem.RefreshLocationAsync(context.Background(), false)

	if gem.RefreshInBackground() {
		t.Error("expected refreshInBackground to be false for multi-master")
	}
}

func TestStartRefreshLocationTimerAsync(t *testing.T) {
	manager := newMockDatabaseAccountManager("https://testaccount.documents.azure.com:443")

	account1, _ := NewDatabaseAccountFromJSON(dbAccountJson1)
	account2, _ := NewDatabaseAccountFromJSON(dbAccountJson2)
	account3, _ := NewDatabaseAccountFromJSON(dbAccountJson3)

	manager.SetDatabaseAccounts(account1, account2, account3)

	policy := &ConnectionPolicy{
		EnableEndpointDiscovery:     true,
		UsingMultipleWriteLocations: true,
	}

	gem := NewGlobalEndpointManager(manager, policy)
	gem.SetBackgroundRefreshLocationTimeIntervalInMS(500)

	if err := gem.Init(context.Background()); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gem.StartRefreshLocationTimerAsync(ctx)

	time.Sleep(700 * time.Millisecond)

	locationCache := gem.GetLocationCache()
	availableReadByLocation := locationCache.GetAvailableReadEndpointByLocation()

	if len(availableReadByLocation) != 1 {
		t.Errorf("expected 1 available read location after first refresh, got %d", len(availableReadByLocation))
	}

	var firstKey string
	for k := range availableReadByLocation {
		firstKey = k
		break
	}
	if firstKey != "East Asia" {
		t.Errorf("expected East Asia, got %s", firstKey)
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}

	time.Sleep(700 * time.Millisecond)

	availableReadByLocation = gem.GetLocationCache().GetAvailableReadEndpointByLocation()
	if _, ok := availableReadByLocation["West US"]; !ok {
		t.Error("expected West US in available read locations after second refresh")
	}

	if gem.IsRefreshing() {
		t.Error("expected isRefreshing to be false")
	}

	gem.Close()
}
