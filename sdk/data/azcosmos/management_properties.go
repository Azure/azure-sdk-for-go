// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// DatabaseProperties represents the properties of a database.
type DatabaseProperties struct {
	// ID is the unique id of the database. Required on [Client.CreateDatabase]; server-assigned
	// fields below are populated by the service and ignored on create.
	ID string `json:"id"`
	// ETag is the entity tag of the database (`_etag`).
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink is the self-link of the database (`_self`).
	SelfLink string `json:"_self,omitempty"`
	// ResourceID is the resource id of the database (`_rid`).
	ResourceID string `json:"_rid,omitempty"`
	// LastModified is when the database was last modified (`_ts`).
	LastModified time.Time `json:"-"`
}

type databasePropertiesWire struct {
	ID           string       `json:"id"`
	ETag         *azcore.ETag `json:"_etag,omitempty"`
	SelfLink     string       `json:"_self,omitempty"`
	ResourceID   string       `json:"_rid,omitempty"`
	LastModified int64        `json:"_ts"`
}

func (p databasePropertiesWire) properties() DatabaseProperties {
	properties := DatabaseProperties{
		ID:         p.ID,
		ETag:       p.ETag,
		SelfLink:   p.SelfLink,
		ResourceID: p.ResourceID,
	}
	if p.LastModified != 0 {
		properties.LastModified = time.Unix(p.LastModified, 0).UTC()
	}
	return properties
}

// PartitionKeyKind represents the type of partition key used by a container.
type PartitionKeyKind string

const (
	// PartitionKeyKindHash selects a single-path hash partition key.
	PartitionKeyKindHash PartitionKeyKind = "Hash"
	// PartitionKeyKindMultiHash selects a hierarchical (multi-path) partition key.
	PartitionKeyKindMultiHash PartitionKeyKind = "MultiHash"
)

// PartitionKeyDefinition defines the partition key path(s) of a container.
type PartitionKeyDefinition struct {
	// Kind is the partition key kind. Empty infers [PartitionKeyKindHash] for a single path and
	// [PartitionKeyKindMultiHash] for more than one.
	Kind PartitionKeyKind `json:"kind,omitempty"`
	// Paths is the list of partition key paths, for example "/tenantId".
	Paths []string `json:"paths"`
	// Version is the partition key hash version. Zero creates the container with version 2.
	// Service definitions that omit the version decode as version 1.
	Version int `json:"version,omitempty"`
	extra   map[string]json.RawMessage
}

// MarshalJSON applies the [PartitionKeyDefinition.Kind] inference documented on the field: an
// empty Kind is resolved from the path count before encoding, since the service expects an
// explicit kind rather than inferring one itself.
func (p PartitionKeyDefinition) MarshalJSON() ([]byte, error) {
	kind := p.Kind
	if kind == "" {
		switch len(p.Paths) {
		case 1:
			kind = PartitionKeyKindHash
		case 0:
		default:
			kind = PartitionKeyKindMultiHash
		}
	}
	version := p.Version
	if version == 0 {
		version = 2
	}
	type alias PartitionKeyDefinition
	return encodeAdditionalManagementProperties(alias{Kind: kind, Paths: p.Paths, Version: version}, p.extra)
}

// UnmarshalJSON preserves version 1 for legacy service definitions that omit the version.
func (p *PartitionKeyDefinition) UnmarshalJSON(data []byte) error {
	type alias PartitionKeyDefinition
	decoded := alias{Version: 1}
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "kind", "paths", "version")
	if err != nil {
		return err
	}
	*p = PartitionKeyDefinition(decoded)
	p.extra = extra
	return nil
}

// IndexingMode selects how a container's indexing policy is applied.
type IndexingMode string

const (
	// IndexingModeConsistent updates the index synchronously with every write.
	IndexingModeConsistent IndexingMode = "consistent"
	// IndexingModeNone disables indexing.
	IndexingModeNone IndexingMode = "none"
)

// IncludedPath is a JSON path included in a container's indexing policy.
type IncludedPath struct {
	Path  string `json:"path"`
	extra map[string]json.RawMessage
}

// ExcludedPath is a JSON path excluded from a container's indexing policy.
type ExcludedPath struct {
	Path string `json:"path"`
}

// CompositeIndexOrder is the sort order of a [CompositeIndex] path.
type CompositeIndexOrder string

const (
	CompositeIndexAscending  CompositeIndexOrder = "ascending"
	CompositeIndexDescending CompositeIndexOrder = "descending"
)

// CompositeIndex is one path within a composite index, used to serve an ORDER BY query over
// multiple properties.
type CompositeIndex struct {
	Path  string              `json:"path"`
	Order CompositeIndexOrder `json:"order"`
}

// IndexingPolicy configures how a container indexes items.
//
// Includes spatial, composite, vector, and full-text index definitions.
type IndexingPolicy struct {
	// Automatic selects whether the policy applies automatically on write.
	Automatic bool `json:"automatic"`
	// IndexingMode selects the indexing mode. Empty uses the service default
	// ([IndexingModeConsistent]).
	IndexingMode IndexingMode `json:"indexingMode,omitempty"`
	// IncludedPaths lists the paths to index. Empty indexes every path.
	IncludedPaths []IncludedPath `json:"includedPaths,omitempty"`
	// ExcludedPaths lists the paths to exclude from indexing.
	ExcludedPaths []ExcludedPath `json:"excludedPaths,omitempty"`
	// CompositeIndexes lists the composite indexes, each an ordered slice of paths.
	CompositeIndexes [][]CompositeIndex `json:"compositeIndexes,omitempty"`
	// SpatialIndexes selects paths and spatial geometry types.
	SpatialIndexes []SpatialIndex `json:"spatialIndexes,omitempty"`
	// VectorIndexes selects vector index algorithms and tuning.
	VectorIndexes []VectorIndex `json:"vectorIndexes,omitempty"`
	// FullTextIndexes selects paths declared in the container's full-text policy.
	FullTextIndexes []FullTextIndex `json:"fullTextIndexes,omitempty"`
	extra           map[string]json.RawMessage
}

// UniqueKey is a set of paths whose combined values must be unique across the items of a
// container.
type UniqueKey struct {
	Paths []string `json:"paths"`
}

// UniqueKeyPolicy configures a container's unique key constraints. Unique keys can only be set
// when the container is created.
type UniqueKeyPolicy struct {
	UniqueKeys []UniqueKey `json:"uniqueKeys"`
}

// ContainerProperties represents the properties of a container.
type ContainerProperties struct {
	// ID is the unique id of the container. Required on [DatabaseClient.CreateContainer].
	ID string `json:"id"`
	// ETag is the entity tag of the container (`_etag`).
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink is the self-link of the container (`_self`).
	SelfLink string `json:"_self,omitempty"`
	// ResourceID is the resource id of the container (`_rid`).
	ResourceID string `json:"_rid,omitempty"`
	// LastModified is when the container was last modified (`_ts`).
	LastModified time.Time `json:"-"`
	// DefaultTimeToLive is the default time to live, in seconds, for items in the container. Nil
	// disables TTL expiration.
	// https://learn.microsoft.com/azure/cosmos-db/time-to-live
	DefaultTimeToLive *int32 `json:"defaultTtl,omitempty"`
	// PartitionKeyDefinition is the partition key path(s) of the container. Required on create and
	// immutable afterward.
	PartitionKeyDefinition PartitionKeyDefinition `json:"partitionKey"`
	// IndexingPolicy configures indexing. Nil uses the service default policy.
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	// UniqueKeyPolicy configures unique key constraints. Can only be set on create.
	UniqueKeyPolicy *UniqueKeyPolicy `json:"uniqueKeyPolicy,omitempty"`
	// ConflictResolutionPolicy selects multi-writer conflict resolution.
	ConflictResolutionPolicy *ConflictResolutionPolicy `json:"conflictResolutionPolicy,omitempty"`
	// VectorEmbeddingPolicy describes vector fields and similarity functions.
	VectorEmbeddingPolicy *VectorEmbeddingPolicy `json:"vectorEmbeddingPolicy,omitempty"`
	// FullTextPolicy describes text fields and languages.
	FullTextPolicy *FullTextPolicy `json:"fullTextPolicy,omitempty"`
	// ChangeFeedPolicy configures all-versions-and-deletes retention.
	ChangeFeedPolicy *ChangeFeedPolicy `json:"changeFeedPolicy,omitempty"`
	// AnalyticalStorageTimeToLive is analytical-store TTL in seconds. Nil leaves it unset.
	AnalyticalStorageTimeToLive *int32 `json:"analyticalStorageTtl,omitempty"`
	extra                       map[string]json.RawMessage
}

type containerPropertiesWire struct {
	ID                          string                    `json:"id"`
	ETag                        *azcore.ETag              `json:"_etag,omitempty"`
	SelfLink                    string                    `json:"_self,omitempty"`
	ResourceID                  string                    `json:"_rid,omitempty"`
	LastModified                int64                     `json:"_ts"`
	DefaultTimeToLive           *int32                    `json:"defaultTtl,omitempty"`
	PartitionKeyDefinition      PartitionKeyDefinition    `json:"partitionKey"`
	IndexingPolicy              *IndexingPolicy           `json:"indexingPolicy,omitempty"`
	UniqueKeyPolicy             *UniqueKeyPolicy          `json:"uniqueKeyPolicy,omitempty"`
	ConflictResolutionPolicy    *ConflictResolutionPolicy `json:"conflictResolutionPolicy,omitempty"`
	VectorEmbeddingPolicy       *VectorEmbeddingPolicy    `json:"vectorEmbeddingPolicy,omitempty"`
	FullTextPolicy              *FullTextPolicy           `json:"fullTextPolicy,omitempty"`
	ChangeFeedPolicy            *ChangeFeedPolicy         `json:"changeFeedPolicy,omitempty"`
	AnalyticalStorageTimeToLive *int32                    `json:"analyticalStorageTtl,omitempty"`
	extra                       map[string]json.RawMessage
}

func (p containerPropertiesWire) properties() ContainerProperties {
	properties := ContainerProperties{
		ID:                          p.ID,
		ETag:                        p.ETag,
		SelfLink:                    p.SelfLink,
		ResourceID:                  p.ResourceID,
		DefaultTimeToLive:           p.DefaultTimeToLive,
		PartitionKeyDefinition:      p.PartitionKeyDefinition,
		IndexingPolicy:              p.IndexingPolicy,
		UniqueKeyPolicy:             p.UniqueKeyPolicy,
		ConflictResolutionPolicy:    p.ConflictResolutionPolicy,
		VectorEmbeddingPolicy:       p.VectorEmbeddingPolicy,
		FullTextPolicy:              p.FullTextPolicy,
		ChangeFeedPolicy:            p.ChangeFeedPolicy,
		AnalyticalStorageTimeToLive: p.AnalyticalStorageTimeToLive,
		extra:                       p.extra,
	}
	if p.LastModified != 0 {
		properties.LastModified = time.Unix(p.LastModified, 0).UTC()
	}
	return properties
}

// UnmarshalJSON decodes service system properties and retains unmodeled properties for replacement.
func (p *ContainerProperties) UnmarshalJSON(data []byte) error {
	var wire containerPropertiesWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*p = wire.properties()
	return nil
}

func (p *containerPropertiesWire) UnmarshalJSON(data []byte) error {
	type alias containerPropertiesWire
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for _, name := range []string{"id", "_etag", "_self", "_rid", "_ts", "defaultTtl", "partitionKey", "indexingPolicy", "uniqueKeyPolicy", "conflictResolutionPolicy", "vectorEmbeddingPolicy", "fullTextPolicy", "changeFeedPolicy", "analyticalStorageTtl"} {
		delete(fields, name)
	}
	*p = containerPropertiesWire(decoded)
	p.extra = fields
	return nil
}

// MarshalJSON preserves service properties unknown to this SDK during read-modify-replace.
func (p ContainerProperties) MarshalJSON() ([]byte, error) {
	type alias ContainerProperties
	encoded, err := json.Marshal(alias(p))
	if err != nil || len(p.extra) == 0 {
		return encoded, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return nil, err
	}
	for key, value := range p.extra {
		if _, exists := fields[key]; !exists {
			fields[key] = value
		}
	}
	return json.Marshal(fields)
}

// UnmarshalJSON decodes database system properties, including the Unix timestamp.
func (p *DatabaseProperties) UnmarshalJSON(data []byte) error {
	var wire databasePropertiesWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	*p = wire.properties()
	return nil
}

// SpatialIndex indexes geometry on a path.
type SpatialIndex struct {
	Path string `json:"path"`
	// Types contains Point, LineString, Polygon, or MultiPolygon.
	Types []string `json:"types"`
	extra map[string]json.RawMessage
}

// VectorIndex describes a vector index.
type VectorIndex struct {
	Path string `json:"path"`
	// Type is flat, quantizedFlat, or diskANN.
	Type string `json:"type"`
	// QuantizerType is the quantizer algorithm, if configured.
	QuantizerType          *string  `json:"quantizerType,omitempty"`
	QuantizationByteSize   *uint32  `json:"quantizationByteSize,omitempty"`
	IndexingSearchListSize *uint32  `json:"indexingSearchListSize,omitempty"`
	VectorIndexShardKey    []string `json:"vectorIndexShardKey,omitempty"`
	extra                  map[string]json.RawMessage
}

// FullTextIndex indexes a path declared in FullTextPolicy.
type FullTextIndex struct {
	Path  string `json:"path"`
	extra map[string]json.RawMessage
}

// ConflictResolutionPolicy selects LastWriterWins or Custom conflict resolution.
type ConflictResolutionPolicy struct {
	Mode                        string `json:"mode"`
	ConflictResolutionPath      string `json:"conflictResolutionPath"`
	ConflictResolutionProcedure string `json:"conflictResolutionProcedure"`
}

// VectorEmbeddingPolicy describes vector-valued item properties.
type VectorEmbeddingPolicy struct {
	VectorEmbeddings []VectorEmbedding `json:"vectorEmbeddings"`
	extra            map[string]json.RawMessage
}

// VectorEmbedding describes a vector field.
type VectorEmbedding struct {
	Path string `json:"path"`
	// DataType is float16, float32, int8, or uint8.
	DataType   string `json:"dataType"`
	Dimensions uint32 `json:"dimensions"`
	// DistanceFunction is euclidean, cosine, or dotproduct.
	DistanceFunction string `json:"distanceFunction"`
	extra            map[string]json.RawMessage
}

// FullTextPolicy describes language analysis for text fields.
type FullTextPolicy struct {
	DefaultLanguage *string        `json:"defaultLanguage,omitempty"`
	FullTextPaths   []FullTextPath `json:"fullTextPaths"`
	extra           map[string]json.RawMessage
}

// FullTextPath describes one text field.
type FullTextPath struct {
	Path     string  `json:"path"`
	Language *string `json:"language,omitempty"`
	extra    map[string]json.RawMessage
}

func decodeAdditionalManagementProperties(data []byte, target any, names ...string) (map[string]json.RawMessage, error) {
	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	for _, name := range names {
		delete(fields, name)
	}
	return fields, nil
}

func encodeAdditionalManagementProperties(value any, extra map[string]json.RawMessage) ([]byte, error) {
	body, err := json.Marshal(value)
	if err != nil || len(extra) == 0 {
		return body, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, err
	}
	for name, value := range extra {
		fields[name] = value
	}
	return json.Marshal(fields)
}

// UnmarshalJSON retains legacy index definitions and other unmodeled path properties.
func (p *IncludedPath) UnmarshalJSON(data []byte) error {
	type alias IncludedPath
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path")
	if err != nil {
		return err
	}
	*p = IncludedPath(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unmodeled included-path configuration.
func (p IncludedPath) MarshalJSON() ([]byte, error) {
	type alias IncludedPath
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains spatial configuration such as geometry bounding boxes.
func (p *SpatialIndex) UnmarshalJSON(data []byte) error {
	type alias SpatialIndex
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path", "types")
	if err != nil {
		return err
	}
	*p = SpatialIndex(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unmodeled spatial configuration.
func (p SpatialIndex) MarshalJSON() ([]byte, error) {
	type alias SpatialIndex
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown indexing properties for read-modify-replace.
func (p *IndexingPolicy) UnmarshalJSON(data []byte) error {
	type alias IndexingPolicy
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "automatic", "indexingMode", "includedPaths", "excludedPaths", "compositeIndexes", "spatialIndexes", "vectorIndexes", "fullTextIndexes")
	if err != nil {
		return err
	}
	*p = IndexingPolicy(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown indexing properties.
func (p IndexingPolicy) MarshalJSON() ([]byte, error) {
	type alias IndexingPolicy
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown vector-index properties for read-modify-replace.
func (p *VectorIndex) UnmarshalJSON(data []byte) error {
	type alias VectorIndex
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path", "type", "quantizerType", "quantizationByteSize", "indexingSearchListSize", "vectorIndexShardKey")
	if err != nil {
		return err
	}
	*p = VectorIndex(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown vector-index properties.
func (p VectorIndex) MarshalJSON() ([]byte, error) {
	type alias VectorIndex
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown text-index properties for read-modify-replace.
func (p *FullTextIndex) UnmarshalJSON(data []byte) error {
	type alias FullTextIndex
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path")
	if err != nil {
		return err
	}
	*p = FullTextIndex(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown text-index properties.
func (p FullTextIndex) MarshalJSON() ([]byte, error) {
	type alias FullTextIndex
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown vector-policy properties for read-modify-replace.
func (p *VectorEmbeddingPolicy) UnmarshalJSON(data []byte) error {
	type alias VectorEmbeddingPolicy
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "vectorEmbeddings")
	if err != nil {
		return err
	}
	*p = VectorEmbeddingPolicy(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown vector-policy properties.
func (p VectorEmbeddingPolicy) MarshalJSON() ([]byte, error) {
	type alias VectorEmbeddingPolicy
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown vector-field properties for read-modify-replace.
func (p *VectorEmbedding) UnmarshalJSON(data []byte) error {
	type alias VectorEmbedding
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path", "dataType", "dimensions", "distanceFunction")
	if err != nil {
		return err
	}
	*p = VectorEmbedding(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown vector-field properties.
func (p VectorEmbedding) MarshalJSON() ([]byte, error) {
	type alias VectorEmbedding
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown full-text-policy properties for read-modify-replace.
func (p *FullTextPolicy) UnmarshalJSON(data []byte) error {
	type alias FullTextPolicy
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "defaultLanguage", "fullTextPaths")
	if err != nil {
		return err
	}
	*p = FullTextPolicy(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown full-text-policy properties.
func (p FullTextPolicy) MarshalJSON() ([]byte, error) {
	type alias FullTextPolicy
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// UnmarshalJSON retains unknown text-path properties for read-modify-replace.
func (p *FullTextPath) UnmarshalJSON(data []byte) error {
	type alias FullTextPath
	var decoded alias
	extra, err := decodeAdditionalManagementProperties(data, &decoded, "path", "language")
	if err != nil {
		return err
	}
	*p = FullTextPath(decoded)
	p.extra = extra
	return nil
}

// MarshalJSON preserves unknown text-path properties.
func (p FullTextPath) MarshalJSON() ([]byte, error) {
	type alias FullTextPath
	return encodeAdditionalManagementProperties(alias(p), p.extra)
}

// ChangeFeedPolicy configures the retention window for all versions and deletes.
type ChangeFeedPolicy struct {
	// RetentionDuration is a whole number of minutes. Zero disables retention.
	RetentionDuration uint64 `json:"retentionDuration"`
}
