// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

func deserializeEmbeddingsArray(msg json.RawMessage) ([]float32, error) {
	if len(msg) == 0 {
		return nil, nil
	}

	if msg[0] == '"' && len(msg) > 2 && msg[len(msg)-1] == '"' {
		// this is a base64 string, not an array of numbers.
		msg = msg[1 : len(msg)-1] // splice out the "'s from the beginning and end of the base64 string
		destBytes, err := base64.StdEncoding.DecodeString(string(msg))

		if err != nil {
			return nil, err
		}

		floats := make([]float32, len(destBytes)/4) // it's a binary serialization of float32s.
		var reader = bytes.NewReader(destBytes)

		if err := binary.Read(reader, binary.LittleEndian, floats); err != nil {
			return nil, err
		}

		return floats, nil
	}

	var v []float32

	if err := json.Unmarshal(msg, &v); err != nil {
		return nil, err
	}

	return v, nil
}
