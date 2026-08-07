// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type TransactionalContentSetter interface {
	SetCRC64([]byte)
}

func (a *PathClientAppendDataOptions) SetCRC64(v []byte) {
	a.TransactionalContentCRC64 = v
}

// MarshalJSON implements the json.Marshaller interface for type Path.
func (p Path) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]any)
	populate(objectMap, "contentLength", p.ContentLength)
	populate(objectMap, "creationTime", p.CreationTime)
	populate(objectMap, "etag", p.ETag)
	populate(objectMap, "EncryptionContext", p.EncryptionContext)
	populate(objectMap, "EncryptionScope", p.EncryptionScope)
	populate(objectMap, "expiryTime", p.ExpiryTime)
	populate(objectMap, "group", p.Group)
	populate(objectMap, "isDirectory", p.IsDirectory)
	populate(objectMap, "lastModified", p.LastModified)
	populate(objectMap, "name", p.Name)
	populate(objectMap, "owner", p.Owner)
	populate(objectMap, "permissions", p.Permissions)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type Path.
func (p *Path) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", p, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "contentLength":
			var rawVal string
			err = unpopulate(val, "ContentLength", &rawVal)
			intVal, _ := strconv.ParseInt(rawVal, 10, 64)
			p.ContentLength = &intVal
			delete(rawMsg, key)
		case "creationTime":
			err = unpopulate(val, "CreationTime", &p.CreationTime)
			delete(rawMsg, key)
		case "etag":
			err = unpopulate(val, "ETag", &p.ETag)
			delete(rawMsg, key)
		case "EncryptionContext":
			err = unpopulate(val, "EncryptionContext", &p.EncryptionContext)
			delete(rawMsg, key)
		case "EncryptionScope":
			err = unpopulate(val, "EncryptionScope", &p.EncryptionScope)
			delete(rawMsg, key)
		case "expiryTime":
			err = unpopulate(val, "ExpiryTime", &p.ExpiryTime)
			delete(rawMsg, key)
		case "group":
			err = unpopulate(val, "Group", &p.Group)
			delete(rawMsg, key)
		case "isDirectory":
			var rawVal string
			err = unpopulate(val, "IsDirectory", &rawVal)
			boolVal, _ := strconv.ParseBool(rawVal)
			p.IsDirectory = &boolVal
			delete(rawMsg, key)
		case "lastModified":
			err = unpopulate(val, "LastModified", &p.LastModified)
			delete(rawMsg, key)
		case "name":
			err = unpopulate(val, "Name", &p.Name)
			delete(rawMsg, key)
		case "owner":
			err = unpopulate(val, "Owner", &p.Owner)
			delete(rawMsg, key)
		case "permissions":
			err = unpopulate(val, "Permissions", &p.Permissions)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", p, err)
		}
	}
	return nil
}
