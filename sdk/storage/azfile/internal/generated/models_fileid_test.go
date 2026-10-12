// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHardLinkListUnmarshalXML(t *testing.T) {
	const body = `<?xml version="1.0" encoding="utf-8"?>
<HardLinks>
  <HardLink><ParentId>13835128424026341376</ParentId><FileName>plain.txt</FileName></HardLink>
  <HardLink><ParentId>0</ParentId><FileName Encoded="false">not%20decoded</FileName></HardLink>
  <HardLink><ParentId>13835128424026341376</ParentId><FileName Encoded="true">file%EF%BF%BE</FileName></HardLink>
</HardLinks>`

	var list HardLinkList
	require.NoError(t, xml.Unmarshal([]byte(body), &list))
	require.Len(t, list.HardLinks, 3)

	require.Equal(t, "13835128424026341376", *list.HardLinks[0].ParentID)
	require.Equal(t, "plain.txt", *list.HardLinks[0].FileName)

	// a name is only decoded when the service marks it as encoded
	require.Equal(t, "0", *list.HardLinks[1].ParentID)
	require.Equal(t, "not%20decoded", *list.HardLinks[1].FileName)

	// ￾ cannot be represented in XML, so the service percent-encodes it
	require.Equal(t, "file￾", *list.HardLinks[2].FileName)
}

func TestHardLinkListUnmarshalXMLEmpty(t *testing.T) {
	var list HardLinkList
	require.NoError(t, xml.Unmarshal([]byte(`<HardLinks />`), &list))
	require.Empty(t, list.HardLinks)

	var link HardLink
	require.NoError(t, xml.Unmarshal([]byte(`<HardLink><ParentId>1</ParentId><FileName Encoded="true"></FileName></HardLink>`), &link))
	require.Equal(t, "1", *link.ParentID)
}

func TestHardLinkUnmarshalXMLInvalidEncoding(t *testing.T) {
	var link HardLink
	require.Error(t, xml.Unmarshal([]byte(`<HardLink><ParentId>1</ParentId><FileName Encoded="true">%ZZ</FileName></HardLink>`), &link))
}
