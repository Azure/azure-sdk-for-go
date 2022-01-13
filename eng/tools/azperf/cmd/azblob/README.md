# Blob Storage Performance Tests

### Setup for test resources
These tests will run against a pre-configured Storage account. The following environment variable will need to be set for the tests to access the live resources:

```
AZURE_STORAGE_CONNECTION_STRING=<live storage account connection string>
```

### Common perf command line options
These options are available for all perf tests:
- `-d --duration=10` Number of seconds to run as many operations (the "run" function) as possible. Default is 10.
- `-w --warm-up=5` Number of seconds to spend warming up the connection before measuring begins. Default is 5.
- `-x --test-proxies` Whether to run the tests against the test proxy server. Specfiy the URL(s) for the proxy endpoint(s) (e.g. "https://localhost:5001"). WARNING: When using with Legacy tests - only HTTPS is supported.

### Common Blob command line options
The options are available for all Blob perf tests:
- `--size=10240` Size in bytes of data to be transferred in upload or download tests. Default is 10240.
- `--count=100` Number of blobs to list. Defaults to 100.
<!-- - `--max-put-size` Maximum size of data uploading in single HTTP PUT. Default is 64\*1024\*1024.
- `--max-block-size` Maximum size of data in a block within a blob. Defaults to 4\*1024\*1024.
- `--buffer-threshold` Minimum block size to prevent full block buffering. Defaults to 4\*1024\*1024+1. -->

### Available Performance Tests
* `BlobUploadTest`: Uploads a stream of `size` bytes to a new Blob
* `BlobDownloadTest`: Download a stream of `size` bytes.
* `BlobListTest`: List a specified number of blobs.