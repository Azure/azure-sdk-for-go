// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

// QueryOptions contains a group of parameters for the Table.Query method.
type ListOptions struct {
	// OData filter expression.
	Filter *string
	// Specifies the media type for the response.
	Format *OdataMetadataFormat
	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string
	// Maximum number of records to return.
	Top *int32
}

func (l *ListOptions) toQueryOptions() *QueryOptions {
	if l == nil {
		return &QueryOptions{}
	}

	return &QueryOptions{
		Filter: l.Filter,
		Format: l.Format,
		Select: l.Select,
		Top:    l.Top,
	}
}
