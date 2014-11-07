package storage

import (
	"fmt"
	"net/url"
)

const (
	StorageServiceBaseUrl = "core.windows.net"
	defaultUseHttps       = true

	blobServiceName  = "blob"
	tableServiceName = "table"
	queueServiceName = "queue"
)

type StorageClient struct {
	accountName string
	accountKey  string
	useHttps    bool
	baseUrl     string
}

func NewBasicClient(accountName, accountKey string) (*StorageClient, error) {
	return NewClient(accountName, accountKey, StorageServiceBaseUrl, defaultUseHttps)
}

func NewClient(accountName, accountKey, blobServiceBaseUrl string, useHttps bool) (*StorageClient, error) {
	if accountName == "" {
		return nil, fmt.Errorf("azure: account name required")
	} else if accountKey == "" {
		return nil, fmt.Errorf("azure: account key required")
	} else if blobServiceBaseUrl == "" {
		return nil, fmt.Errorf("azure: base storage service url required")
	}

	return &StorageClient{
		accountName: accountName,
		accountKey:  accountKey,
		useHttps:    useHttps,
		baseUrl:     blobServiceBaseUrl}, nil
}

func (c StorageClient) getBaseUrl(service string) string {
	scheme := "http"
	if c.useHttps {
		scheme = "https"
	}

	host := fmt.Sprintf("%s.%s.%s", c.accountName, service, c.baseUrl)

	u := &url.URL{
		Scheme: scheme,
		Host:   host}
	return u.String()
}

func (c StorageClient) GetBlobService() *BlobStorageClient {
	return &BlobStorageClient{c}
}
