package storage

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"log"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/stretchr/testify/suite"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
)

type (
	// eventHubSuite encapsulates a end to end test of Event Hubs with build up and tear down of all EH resources
	testSuite struct {
		test.BaseSuite
		AccountName string
		ServiceURL  *azblob.ServiceURL
	}
)

func TestStorage(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (ts *testSuite) SetupSuite() {
	ts.BaseSuite.SetupSuite()
	ts.AccountName = test.MustGetEnv("STORAGE_ACCOUNT_NAME")
}

func (ts *testSuite) TearDownSuite() {
	ts.BaseSuite.TearDownSuite()
}

func (ts *testSuite) TestCredential() {
	containerName := "foo"
	blobName := "bar"
	message := "Hello World!!"
	tokenProvider, err := NewAADSASCredential(ts.SubscriptionID, ts.ResourceGroupName, ts.AccountName, containerName, AADSASCredentialWithEnvironmentVars())
	if err != nil {
		ts.T().Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), shortTimeout)
	defer cancel()
	pipeline := azblob.NewPipeline(tokenProvider, azblob.PipelineOptions{})
	fooURL, err := url.Parse("https://" + ts.AccountName + ".blob." + ts.Env.StorageEndpointSuffix + "/" + containerName)
	if err != nil {
		ts.T().Error(err)
	}

	containerURL := azblob.NewContainerURL(*fooURL, pipeline)
	defer func() {
		if res, err := containerURL.Delete(ctx, azblob.ContainerAccessConditions{}); err != nil {
			log.Fatal(err, res)
		}
	}()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		ts.T().Error(err)
	}

	blobURL := containerURL.NewBlobURL(blobName).ToBlockBlobURL()
	_, err = blobURL.Upload(ctx, strings.NewReader(message), azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		ts.T().Error(err)
	}
}
