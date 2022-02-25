// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorResponseUnmarshal(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"singleline", "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>"},
		{"multiline", "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n<Error>\n  <Code>ContainerAlreadyExists</Code>\n  <Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message>\n</Error>"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			se := StorageError{}
			require.Nil(t, xml.Unmarshal([]byte(c.input), &se))

			require.Contains(t, se.details, "Code")
			require.Equal(t, "ContainerAlreadyExists", se.details["Code"])

			require.Equal(t, "The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z", se.description)
		})
	}
}
