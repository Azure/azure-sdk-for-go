//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

/*
Package azfile provides access to Azure File Storage.
For more information please see https://learn.microsoft.com/rest/api/storageservices/file-service-rest-api

The azfile package is capable of :-
    - Creating, deleting, and querying shares in an account
    - Creating, deleting, and querying directories in a share
    - Creating, deleting, and querying files in a share or directory
    - Creating Shared Access Signature for authentication

Types of Resources

The azfile package allows you to interact with four types of resources :-

* Azure storage accounts.
* Shares within those storage accounts.
* Directories within those shares.
* Files within those shares or directories.

The Azure File Storage (azfile) client library for Go allows you to interact with each of these components through the use of a dedicated client object.
To create a client object, you will need the account's file service endpoint URL and a credential that allows you to access the account.

Types of Credentials

The clients support different forms of authentication.
The azfile library supports authorization via a shared key, Connection String,
or with a Shared Access Signature token.

Using a Shared Key

To use an account shared key (aka account key or access key), provide the key as a string.
This can be found in your storage account in the Azure Portal under the "Access Keys" section.

Use the key as the credential parameter to authenticate the client:

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	serviceURL := fmt.Sprintf("https://%s.file.core.windows.net/", accountName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handle(err)

	serviceClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handle(err)

	fmt.Println(serviceClient.URL())

Using a Connection String

Depending on your use case and authorization method, you may prefer to initialize a client instance with a connection string instead of providing the account URL and credential separately.
To do this, pass the connection string to the service client's `NewClientFromConnectionString` method.
The connection string can be found in your storage account in the Azure Portal under the "Access Keys" section.

	connStr := "DefaultEndpointsProtocol=https;AccountName=<my_account_name>;AccountKey=<my_account_key>;EndpointSuffix=core.windows.net"
	serviceClient, err := azfile.NewServiceClientFromConnectionString(connStr, nil)
	handle(err)

Using a Shared Access Signature (SAS) Token

To use a shared access signature (SAS) token, provide the token at the end of your service URL.
You can generate a SAS token from the Azure Portal under Shared Access Signature or use the service.Client.GetSASURL() functions.

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.file.core.windows.net/", accountName)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handle(err)
	serviceClient, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	handle(err)
	fmt.Println(serviceClient.URL())

	// Alternatively, you can create SAS on the fly

	resources := sas.AccountResourceTypes{Service: true}
	permission := sas.AccountPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(0, 0, 1)
	serviceURLWithSAS, err := serviceClient.GetSASURL(resources, permission, expiry, &service.GetSASURLOptions{StartTime: &start})
	handle(err)

	serviceClientWithSAS, err := service.NewClientWithNoCredential(serviceURLWithSAS, nil)
	handle(err)

	fmt.Println(serviceClientWithSAS.URL())

Types of Clients

There are four different clients provided to interact with the various components of the File Service:

1. **`ServiceClient`**
    * Get and set account settings.
    * Query, create, delete and restore shares within the account.

2. **`ShareClient`**
    * Get and set share access settings, properties, and metadata.
    * Create, delete, and query directories and files within the share.
    * `lease.ShareClient` to support share lease management.

3. **`DirectoryClient`**
	* Create or delete operations on a given directory.
    * Get and set directory properties.
    * List sub-directories and files within the given directory.

3. **`FileClient`**
    * Get and set file properties.
    * Perform CRUD operations on a given file.
    * `FileLeaseClient` to support file lease management.

Examples

	// Your account name and key can be obtained from the Azure Portal.
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	handle(err)

	// The service URL for file endpoints is usually in the form: http(s)://<account>.file.core.windows.net/
	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/", accountName), cred, nil)
	handle(err)

	// ===== 1. Create a share =====

	// First, create a share client, and use the Create method to create a new share in your account
	shareClient := serviceClient.NewShareClient("testshare")
	handle(err)

	// All APIs have an options' bag struct as a parameter.
	// The options' bag struct allows you to specify optional parameters such as metadata, quota, etc.
	// If you want to use the default options, pass in nil.
	_, err = shareClient.Create(context.TODO(), nil)
	handle(err)

	// ===== 2. Create a directory =====

	// First, create a directory client, and use the Create method to create a new directory in the share
	dirClient := shareClient.NewDirectoryClient("testdir")
	_, err = dirClient.Create(context.TODO(), nil)

	// ===== 3. Upload and Download a file =====
	uploadData := "Hello world!"

	// First, create a file client, and use the Create method to create a new file in the directory
	fileClient := dirClient.NewFileClient("HelloWorld.txt")
	_, err = fileClient.Create(context.TODO(), int64(len(uploadData)), nil)
	handle(err)

	// Upload data to the file
	_, err = fileClient.UploadRange(context.TODO(), 0, streaming.NopCloser(strings.NewReader(uploadData)), nil)
	handle(err)

	// Download the file's contents and ensure that the download worked properly
	fileDownloadResponse, err := fileClient.DownloadStream(context.TODO(), nil)
	handle(err)

	// Use io.readAll to read the downloaded data.
	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	reader := fileDownloadResponse.Body
	downloadData, err := io.ReadAll(reader)
	handle(err)
	if string(downloadData) != uploadData {
		handle(errors.New("uploaded data should be same as downloaded data"))
	}

	if err = reader.Close(); err != nil {
		handle(err)
		return
	}

	// ===== 3. List directories and files in a share =====
	// List methods returns a pager object which can be used to iterate over the results of a paging operation.
	// To iterate over a page use the NextPage(context.Context) to fetch the next page of results.
	// PageResponse() can be used to iterate over the results of the specific page.
	// Always check the Err() method after paging to see if an error was returned by the pager. A pager will return either an error or the page of results.
	// The below code lists the contents only for a single level of the directory hierarchy.
	rootDirClient := shareClient.NewRootDirectoryClient()
	pager := rootDirClient.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handle(err)
		for _, d := range resp.Segment.Directories {
			fmt.Println(*d.Name)
		}
		for _, f := range resp.Segment.Files {
			fmt.Println(*f.Name)
		}
	}

	// Delete the file.
	_, err = fileClient.Delete(context.TODO(), nil)
	handle(err)

	// Delete the directory.
	_, err = dirClient.Delete(context.TODO(), nil)
	handle(err)

	// Delete the share.
	_, err = shareClient.Delete(context.TODO(), nil)
	handle(err)
*/

package azfile
