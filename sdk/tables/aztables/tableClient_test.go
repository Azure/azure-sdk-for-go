// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	// "bytes"
	// "os"
	// "strconv"
	// "strings"
	// "time"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

func (s *aztestsSuite) TestContainerCreateAccessContainer(c *chk.C) {
	tableClient, _ := createTableClient(StorageEndpoint)

	_, err := tableClient.Create(ctx, generateName())
	// defer deleteContainer(c, containerClient)
	c.Assert(err, chk.IsNil)
}