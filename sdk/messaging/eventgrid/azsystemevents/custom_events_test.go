package azsystemevents

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalInternalACSRouterCommunicationError(t *testing.T) {
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
		require.Equal(t, []*Error{
			{
				Code:    "Failure",
				message: "Classification failed due to <reason>",
				InnerError: &Error{
					Code: "InnerFailure", message: "Classification failed due to <reason>",
				},
			},
		}, unmarshalledErrs)
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
		require.Equal(t, []*Error{
			{Code: "Failure", message: "Classification failed due to <reason>"},
			{Code: "FailureAsWell", message: "Classification failed due to <reason>"},
		}, unmarshalledErrs)
	})

	t.Run("Recursive", func(t *testing.T) {
		text := `[
			{
				"code": "Failure",
				"message": "FailureMessage",
				"innererror": {
					"code": "InnerFailure",
					"message": "InnerFailureMessage",
					"innererror": {
						"code": "Inner.InnerFailure",
						"message": "Inner.InnerFailureMessage"
					}
				},
				"details": [
					{
						"code": "Detail",
						"message": "DetailMessage",
						"innererror": {
							"code": "Detail.InnerFailure",
							"message": "Detail.InnerFailure.Message"
						}
					}							
				]
			}
		]`

		var unmarshalledErrs []*Error
		err := unmarshalInternalACSRouterCommunicationError(json.RawMessage(text), "FieldName", &unmarshalledErrs)
		require.NoError(t, err)
		require.Equal(t, []*Error{
			{
				Code:    "Failure",
				message: "FailureMessage",
				InnerError: &Error{
					Code:    "InnerFailure",
					message: "InnerFailureMessage",
					InnerError: &Error{
						Code:    "Inner.InnerFailure",
						message: "Inner.InnerFailureMessage",
					},
				},
				Details: []*Error{
					{
						Code:    "Detail",
						message: "DetailMessage",
						InnerError: &Error{
							Code:    "Detail.InnerFailure",
							message: "Detail.InnerFailure.Message",
						},
					},
				},
			},
		}, unmarshalledErrs)
	})
}
