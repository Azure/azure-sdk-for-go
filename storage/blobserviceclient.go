package storage

// BlobStorageClient contains operations for Microsoft Azure Blob Storage
// Service.
type BlobStorageClient struct {
	client Client
	auth   authentication
}

// GetServiceProperties gets the properties of your storage account's blob service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/get-blob-service-properties
func (b *BlobStorageClient) GetServiceProperties() (*ServiceProperties, error) {
	return b.client.getServiceProperties(blobServiceName, b.auth)
}

// SetServiceProperties sets the properties of your storage account's blob service.
// See: https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/set-blob-service-properties
func (b *BlobStorageClient) SetServiceProperties(props ServiceProperties) error {
	return b.client.setServiceProperties(props, blobServiceName, b.auth)
}
