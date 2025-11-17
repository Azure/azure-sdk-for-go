//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteTypeProperly(t *testing.T) {
	modelsGo := `before
// ACSCallParticipantEventData - Schema of common properties of all participant events
type ACSCallParticipantEventData struct {
	// REQUIRED; The correlationId of calling event
	CorrelationID *string

	// REQUIRED; The call id of the server
	ServerCallID *string

	// REQUIRED; The call participant who initiated the call.
	StartedBy *ACSCallParticipantProperties

	// The display name of the participant.
	DisplayName *string

	// The group metadata
	Group *ACSCallGroupProperties

	// Is the calling event a room call.
	IsRoomsCall *bool

	// Is two-party in calling event.
	IsTwoParty *bool

	// The id of the participant.
	ParticipantID *string

	// The room metadata
	Room *ACSCallRoomProperties

	// The user of the call participant
	User *ACSCallParticipantProperties

	// The user agent of the participant.
	UserAgent *string
}
after
`

	modelsSerdeGo := `before

// MarshalJSON implements the json.Marshaller interface for type ACSCallParticipantEventData.
func (a ACSCallParticipantEventData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "correlationId", a.CorrelationID)
	populate(objectMap, "displayName", a.DisplayName)
	populate(objectMap, "group", a.Group)
	populate(objectMap, "isRoomsCall", a.IsRoomsCall)
	populate(objectMap, "isTwoParty", a.IsTwoParty)
	populate(objectMap, "participantId", a.ParticipantID)
	populate(objectMap, "room", a.Room)
	populate(objectMap, "serverCallId", a.ServerCallID)
	populate(objectMap, "startedBy", a.StartedBy)
	populate(objectMap, "user", a.User)
	populate(objectMap, "userAgent", a.UserAgent)
	return json.Marshal(objectMap)
}

between

// UnmarshalJSON implements the json.Unmarshaller interface for type ACSCallParticipantEventData.
func (a *ACSCallParticipantEventData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "correlationId":
			err = unpopulate(val, "CorrelationID", &a.CorrelationID)
			delete(rawMsg, key)
		case "displayName":
			err = unpopulate(val, "DisplayName", &a.DisplayName)
			delete(rawMsg, key)
		case "group":
			err = unpopulate(val, "Group", &a.Group)
			delete(rawMsg, key)
		case "isRoomsCall":
			err = unpopulate(val, "IsRoomsCall", &a.IsRoomsCall)
			delete(rawMsg, key)
		case "isTwoParty":
			err = unpopulate(val, "IsTwoParty", &a.IsTwoParty)
			delete(rawMsg, key)
		case "participantId":
			err = unpopulate(val, "ParticipantID", &a.ParticipantID)
			delete(rawMsg, key)
		case "room":
			err = unpopulate(val, "Room", &a.Room)
			delete(rawMsg, key)
		case "serverCallId":
			err = unpopulate(val, "ServerCallID", &a.ServerCallID)
			delete(rawMsg, key)
		case "startedBy":
			err = unpopulate(val, "StartedBy", &a.StartedBy)
			delete(rawMsg, key)
		case "user":
			err = unpopulate(val, "User", &a.User)
			delete(rawMsg, key)
		case "userAgent":
			err = unpopulate(val, "UserAgent", &a.UserAgent)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}

after
`

	newModelsGo, newModelSerdeGo := deleteType("ACSCallParticipantEventData", modelsGo, modelsSerdeGo)
	require.Equal(t, `before
after
`, newModelsGo)

	require.Equal(t, `before


between


after
`, newModelSerdeGo)
}

func TestUseCustomUnpopulate(t *testing.T) {
	goCode := `// UnmarshalJSON implements the json.Unmarshaller interface for type MyTypeName.
func (a *MyTypeName) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "context":
			err = unpopulate(val, "Context", &a.Context)
			delete(rawMsg, key)
		case "error":
			err = unpopulate(val, "Error", &a.Error)
			delete(rawMsg, key)
		case "from":
			err = unpopulate(val, "From", &a.From)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}
`

	expectedGoCode := `// UnmarshalJSON implements the json.Unmarshaller interface for type MyTypeName.
func (a *MyTypeName) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", a, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "context":
			err = unpopulate(val, "Context", &a.Context)
			delete(rawMsg, key)
		case "error":
			err =  customUnmarshaller(val, "Error", &a.Error)
			delete(rawMsg, key)
		case "from":
			err = unpopulate(val, "From", &a.From)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", a, err)
		}
	}
	return nil
}
`

	actualGoCode := UseCustomUnpopulate(goCode, "MyTypeName.Error", "customUnmarshaller")
	require.Equal(t, expectedGoCode, actualGoCode)
}
