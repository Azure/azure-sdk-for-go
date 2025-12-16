// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"
	"fmt"
)

// MongoDBChatExtensionParametersEmbeddingDependency contains the embedding dependency for the [MongoDBChatExtensionParameters].
// NOTE: This should be created using [azopenai.NewMongoDBChatExtensionParametersEmbeddingDependency]
type MongoDBChatExtensionParametersEmbeddingDependency struct {
	value any
}

// NewMongoDBChatExtensionParametersEmbeddingDependency creates a [azopenai.MongoDBChatExtensionParametersEmbeddingDependency].
func NewMongoDBChatExtensionParametersEmbeddingDependency[T OnYourDataDeploymentNameVectorizationSource | OnYourDataEndpointVectorizationSource](value T) *MongoDBChatExtensionParametersEmbeddingDependency {
	switch any(value).(type) {
	case OnYourDataDeploymentNameVectorizationSource:
		return &MongoDBChatExtensionParametersEmbeddingDependency{value: value}
	case OnYourDataEndpointVectorizationSource:
		return &MongoDBChatExtensionParametersEmbeddingDependency{value: value}
	default:
		panic(fmt.Sprintf("Invalid type %T for MongoDBChatExtensionParametersEmbeddingDependency", value))
	}
}

// MarshalJSON implements the json.Marshaller interface for type MongoDBChatExtensionParametersEmbeddingDependency.
func (c MongoDBChatExtensionParametersEmbeddingDependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}
