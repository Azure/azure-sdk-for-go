//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsystemevents

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalInternalACSRouterCommunicationError(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		var unmarshalledErrs []*Error
		err := unmarshalInternalACSRouterCommunicationError(json.RawMessage("null"), "FieldName", &unmarshalledErrs)
		require.NoError(t, err)
		require.Empty(t, unmarshalledErrs)

		err = unmarshalInternalACSRouterCommunicationError(json.RawMessage("[{ \"Code\": null }]"), "FieldName", &unmarshalledErrs)
		require.NoError(t, err)
		require.Empty(t, unmarshalledErrs[0].Code)
	})

	t.Run("single", func(t *testing.T) {
		text := `[
			{
			"code": "Failure",
			"message": "Classification failed due to <reason>",
			"target": null,
			"innererror": {
							"code": "InnerFailure",
							"message": "Classification failed due to <reason>",
							"target": null},
			"details": null
			}
		]`

		var unmarshalledErrs []*Error
		err := unmarshalInternalACSRouterCommunicationError(json.RawMessage(text), "FieldName", &unmarshalledErrs)
		require.NoError(t, err)

		require.Equal(t, "Failure", unmarshalledErrs[0].Code)
		require.Equal(t, "Code: Failure\n"+
			"Message: Classification failed due to <reason>\n"+
			"InnerError:\n"+
			"  Code: InnerFailure\n"+
			"  Message: Classification failed due to <reason>\n", unmarshalledErrs[0].Error())
	})

	t.Run("Multiple", func(t *testing.T) {
		text := `[
			{
				"code": "Failure",
				"message": "Classification failed due to <reason>"
			},
			{
				"code": "FailureAsWell",
				"message": "Classification failed due to <reason>"
			}
		]`

		var unmarshalledErrs []*Error
		err := unmarshalInternalACSRouterCommunicationError(json.RawMessage(text), "FieldName", &unmarshalledErrs)
		require.NoError(t, err)

		require.Equal(t, "Failure", unmarshalledErrs[0].Code)
		require.Equal(t, "Code: Failure\n"+
			"Message: Classification failed due to <reason>\n", unmarshalledErrs[0].Error())
		require.Equal(t, "FailureAsWell", unmarshalledErrs[1].Code)
		require.Equal(t, "Code: FailureAsWell\n"+
			"Message: Classification failed due to <reason>\n", unmarshalledErrs[1].Error())
	})
}

func TestUnmarshalInternalAcsRouterCommunicationErrorRecursive(t *testing.T) {
	text := `[
		{
			"code": "Root.Failure",
			"message": "Root.Message",
			"innererror": {
				"code": "Root->Inner.Failure",
				"message": "Root->Inner.Message",
				"innererror": {
					"code": "Root->Inner->Inner.Failure",
					"message": "Root->Inner->Inner.Message"
				},
				"details": [
					{
						"code": "Root->Inner->Details[0].Failure",
						"message": "Root->Inner->Details[0].Message",
						"innererror": {
							"code": "Root->Inner->Details[0].Inner.Failure",
							"message": "Root->Inner->Details[0].Inner.Message"
						}
					}
				]
			},
			"details": [
				{
					"code": "Root->Details[0].Failure",
					"message": "Root->Details[0].Message",
					"innererror": {
						"code": "Root->Details[0]->Inner.Failure",
						"message": "Root->Details[0]->Inner.Message"
					}
				}
			]
		}
	]`

	var unmarshalledErrs []*Error
	err := unmarshalInternalACSRouterCommunicationError(json.RawMessage(text), "FieldName", &unmarshalledErrs)
	require.NoError(t, err)

	require.NotEmpty(t, unmarshalledErrs)
	require.Equal(t, "Root.Failure", unmarshalledErrs[0].Code)

	expectedMsg := "Code: Root.Failure\n" +
		"Message: Root.Message\n" +
		"Details[0]:\n" +
		"  Code: Root->Details[0].Failure\n" +
		"  Message: Root->Details[0].Message\n" +
		"  InnerError:\n" +
		"    Code: Root->Details[0]->Inner.Failure\n" +
		"    Message: Root->Details[0]->Inner.Message\n" +
		"InnerError:\n" +
		"  Code: Root->Inner.Failure\n" +
		"  Message: Root->Inner.Message\n" +
		"  Details[0]:\n" +
		"    Code: Root->Inner->Details[0].Failure\n" +
		"    Message: Root->Inner->Details[0].Message\n" +
		"    InnerError:\n" +
		"      Code: Root->Inner->Details[0].Inner.Failure\n" +
		"      Message: Root->Inner->Details[0].Inner.Message\n" +
		"  InnerError:\n" +
		"    Code: Root->Inner->Inner.Failure\n" +
		"    Message: Root->Inner->Inner.Message\n"

	require.Equal(t, expectedMsg, unmarshalledErrs[0].Error())
}

func TestFormattingErrors(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		msg := unpackMsg(internalACSRouterCommunicationError{
			Code:    to.Ptr("code"),
			Message: to.Ptr("message"),
		}, "")

		require.Equal(t, "Code: code\n"+
			"Message: message\n", msg)
	})

	t.Run("InnerError", func(t *testing.T) {
		msg := unpackMsg(internalACSRouterCommunicationError{
			Code:    to.Ptr("code"),
			Message: to.Ptr("message"),
			Innererror: &internalACSRouterCommunicationError{
				Code:    to.Ptr("inner.code"),
				Message: to.Ptr("inner.message"),
			},
		}, "")

		require.Equal(t, "Code: code\n"+
			"Message: message\n"+
			"InnerError:\n"+
			"  Code: inner.code\n"+
			"  Message: inner.message\n", msg)
	})

	t.Run("InnerInnerError", func(t *testing.T) {
		msg := unpackMsg(internalACSRouterCommunicationError{
			Code:    to.Ptr("code"),
			Message: to.Ptr("message"),
			Innererror: &internalACSRouterCommunicationError{
				Code:    to.Ptr("inner.code"),
				Message: to.Ptr("inner.message"),
				Innererror: &internalACSRouterCommunicationError{
					Code:    to.Ptr("inner->inner.code"),
					Message: to.Ptr("inner->inner.message"),
				},
			},
		}, "")

		require.Equal(t, "Code: code\n"+
			"Message: message\n"+
			"InnerError:\n"+
			"  Code: inner.code\n"+
			"  Message: inner.message\n"+
			"  InnerError:\n"+
			"    Code: inner->inner.code\n"+
			"    Message: inner->inner.message\n", msg)
	})

	t.Run("Details", func(t *testing.T) {
		msg := unpackMsg(internalACSRouterCommunicationError{
			Code:    to.Ptr("code"),
			Message: to.Ptr("message"),
			Details: []internalACSRouterCommunicationError{
				{
					Code:    to.Ptr("details[0].code"),
					Message: to.Ptr("details[0].message"),
					Details: []internalACSRouterCommunicationError{
						{
							Code:    to.Ptr("details[0]->details[0].code"),
							Message: to.Ptr("details[0]->details[0].message"),
						},
						{
							Code:    to.Ptr("details[0]->details[1].code"),
							Message: to.Ptr("details[0]->details[1].message"),
						},
					}},
			},
		}, "")

		require.Equal(t, "Code: code\n"+
			"Message: message\n"+
			"Details[0]:\n"+
			"  Code: details[0].code\n"+
			"  Message: details[0].message\n"+
			"  Details[0]:\n"+
			"    Code: details[0]->details[0].code\n"+
			"    Message: details[0]->details[0].message\n"+
			"  Details[1]:\n"+
			"    Code: details[0]->details[1].code\n"+
			"    Message: details[0]->details[1].message\n", msg)
	})
}
