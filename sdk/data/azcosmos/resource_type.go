// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// resourceType defines supported values for resources.
type resourceType int

const (
	resourceTypeDatabase            resourceType = 0
	resourceTypeCollection          resourceType = 1
	resourceTypeDocument            resourceType = 2
	resourceTypeUser                resourceType = 4
	resourceTypePermission          resourceType = 5
	resourceTypeConflict            resourceType = 107
	resourceTypeStoredProcedure     resourceType = 109
	resourceTypeTrigger             resourceType = 110
	resourceTypeUserDefinedFunction resourceType = 111
	resourceTypeOffer               resourceType = 113
	resourceTypeDatabaseAccount     resourceType = 118
	resourceTypePartitionKeyRange   resourceType = 125
	resourceTypeClientEncryptionKey resourceType = 141
)

func (r resourceType) String() string {
	switch r {
	case resourceTypeDatabase:
		return "database"
	case resourceTypeCollection:
		return "collection"
	case resourceTypeDocument:
		return "document"
	case resourceTypeUser:
		return "user"
	case resourceTypePermission:
		return "permission"
	case resourceTypeConflict:
		return "conflict"
	case resourceTypeStoredProcedure:
		return "storedprocedure"
	case resourceTypeTrigger:
		return "trigger"
	case resourceTypeUserDefinedFunction:
		return "userdefinedfunctions"
	case resourceTypeOffer:
		return "offer"
	case resourceTypeDatabaseAccount:
		return "databaseaccount"
	case resourceTypePartitionKeyRange:
		return "partitionkeyrange"
	case resourceTypeClientEncryptionKey:
		return "clientencryptionkey"
	default:
		return ""
	}
}
