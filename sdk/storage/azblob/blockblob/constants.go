package blockblob

//nolint
const (
	_1MiB = 1024 * 1024

	// MaxUploadBlobBytes indicates the maximum number of bytes that can be sent in a call to Upload.
	MaxUploadBlobBytes = 256 * 1024 * 1024 // 256MB

	// MaxStageBlockBytes indicates the maximum number of bytes that can be sent in a call to StageBlock.
	MaxStageBlockBytes = 4000 * 1024 * 1024 // 4GB

	// MaxBlocks indicates the maximum number of blocks allowed in a block blob.
	MaxBlocks = 50000
)
