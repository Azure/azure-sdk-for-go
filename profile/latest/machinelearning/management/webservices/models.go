package webservices

import (
	 original "github.com/Azure/azure-sdk-for-go/service/machinelearning/management/2017-01-01/webservices"
)

type (
	 ManagementClient = original.ManagementClient
	 AssetType = original.AssetType
	 ColumnFormat = original.ColumnFormat
	 ColumnType = original.ColumnType
	 DiagnosticsLevel = original.DiagnosticsLevel
	 InputPortType = original.InputPortType
	 OutputPortType = original.OutputPortType
	 ParameterType = original.ParameterType
	 ProvisioningState = original.ProvisioningState
	 AssetItem = original.AssetItem
	 AsyncOperationErrorInfo = original.AsyncOperationErrorInfo
	 AsyncOperationStatus = original.AsyncOperationStatus
	 BlobLocation = original.BlobLocation
	 ColumnSpecification = original.ColumnSpecification
	 CommitmentPlan = original.CommitmentPlan
	 DiagnosticsConfiguration = original.DiagnosticsConfiguration
	 ExampleRequest = original.ExampleRequest
	 GraphEdge = original.GraphEdge
	 GraphNode = original.GraphNode
	 GraphPackage = original.GraphPackage
	 GraphParameter = original.GraphParameter
	 GraphParameterLink = original.GraphParameterLink
	 InputPort = original.InputPort
	 Keys = original.Keys
	 MachineLearningWorkspace = original.MachineLearningWorkspace
	 ModeValueInfo = original.ModeValueInfo
	 ModuleAssetParameter = original.ModuleAssetParameter
	 OutputPort = original.OutputPort
	 PaginatedWebServicesList = original.PaginatedWebServicesList
	 Parameter = original.Parameter
	 Properties = original.Properties
	 PropertiesForGraph = original.PropertiesForGraph
	 RealtimeConfiguration = original.RealtimeConfiguration
	 Resource = original.Resource
	 ServiceInputOutputSpecification = original.ServiceInputOutputSpecification
	 StorageAccount = original.StorageAccount
	 TableSpecification = original.TableSpecification
	 WebService = original.WebService
	 GroupClient = original.GroupClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 AssetTypeModule = original.AssetTypeModule
	 AssetTypeResource = original.AssetTypeResource
	 Byte = original.Byte
	 Char = original.Char
	 Complex128 = original.Complex128
	 Complex64 = original.Complex64
	 DateTime = original.DateTime
	 DateTimeOffset = original.DateTimeOffset
	 Double = original.Double
	 Duration = original.Duration
	 Float = original.Float
	 Int16 = original.Int16
	 Int32 = original.Int32
	 Int64 = original.Int64
	 Int8 = original.Int8
	 Uint16 = original.Uint16
	 Uint32 = original.Uint32
	 Uint64 = original.Uint64
	 Uint8 = original.Uint8
	 Boolean = original.Boolean
	 Integer = original.Integer
	 Number = original.Number
	 String = original.String
	 All = original.All
	 Error = original.Error
	 None = original.None
	 Dataset = original.Dataset
	 OutputPortTypeDataset = original.OutputPortTypeDataset
	 ParameterTypeBoolean = original.ParameterTypeBoolean
	 ParameterTypeColumnPicker = original.ParameterTypeColumnPicker
	 ParameterTypeCredential = original.ParameterTypeCredential
	 ParameterTypeDataGatewayName = original.ParameterTypeDataGatewayName
	 ParameterTypeDouble = original.ParameterTypeDouble
	 ParameterTypeEnumerated = original.ParameterTypeEnumerated
	 ParameterTypeFloat = original.ParameterTypeFloat
	 ParameterTypeInt = original.ParameterTypeInt
	 ParameterTypeMode = original.ParameterTypeMode
	 ParameterTypeParameterRange = original.ParameterTypeParameterRange
	 ParameterTypeScript = original.ParameterTypeScript
	 ParameterTypeString = original.ParameterTypeString
	 Failed = original.Failed
	 Provisioning = original.Provisioning
	 Succeeded = original.Succeeded
	 Unknown = original.Unknown
)
