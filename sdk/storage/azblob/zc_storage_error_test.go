//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestErrorResponseUnmarshal() {
	t := s.T()

	cases := []struct {
		name  string
		input string
	}{
		{"singleline", "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><Error><Code>ContainerAlreadyExists</Code><Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message></Error>"},
		{"multiline", "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n<Error>\n  <Code>ContainerAlreadyExists</Code>\n  <Message>The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z</Message>\n</Error>"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_require := assert.New(t)
			se := StorageError{}
			_require.Nil(xml.Unmarshal([]byte(c.input), &se))

			_require.Contains(se.details, "Code")
			_require.Equal("ContainerAlreadyExists", se.details["Code"])

			_require.Equal("The specified container already exists.\nRequestId:73b2473b-c1c8-4162-97bb-dc171bff61c9\nTime:2021-12-13T19:45:40.679Z", se.description)
		})
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestInternalErrorWrapping() {
	t := s.T()

	cases := []struct {
		name   string
		err    *InternalError
		target error
		result bool
	}{
		{
			"errors.Is reports true if checked against an internal error",
			&InternalError{errors.New("foo")},
			&InternalError{errors.New("bar")},
			true,
		},
		{
			"errors.Is reports true if checked against wrapped error",
			&InternalError{context.Canceled},
			context.Canceled,
			true,
		},
		{
			"errors.Is reports false if no match",
			&InternalError{context.Canceled},
			errors.New("foo"),
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.result, errors.Is(c.err, c.target))
		})
	}
}
