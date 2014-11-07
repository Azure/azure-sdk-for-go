package storage

type BlobStorageClient struct {
	client StorageClient
}

func (c BlobStorageClient) getBaseUrl() string {
	return c.client.getBaseUrl(blobServiceName)
}
