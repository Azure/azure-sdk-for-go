package storage

import (
	"fmt"
	"net/url"
	"strings"
)

// This consts are meant to help with opdata supported operations
const (
	OdataTypeSuffix = "@odata.type"

	// Types

	OdataBinary   = "Edm.Binary"
	OdataDateTime = "Edm.DateTime"
	OdataGUID     = "Edm.Guid"
	OdataInt64    = "Edm.Int64"

	// Query options

	OdataFilter  = "$filter"
	OdataOrderBy = "$orderby"
	OdataTop     = "$top"
	OdataSkip    = "$skip"
	OdataCount   = "$count"
	OdataExpand  = "$expand"
	OdataSelect  = "$select"
	OdataSearch  = "$search"
)

func fixOdataQuery(odataQuery url.Values) url.Values {
	for k, v := range odataQuery {
		if !strings.HasPrefix(k, "$") {
			newkey := fmt.Sprintf("$%v", k)
			odataQuery[newkey] = v
			odataQuery.Del(k)
		}
	}
	return odataQuery
}
