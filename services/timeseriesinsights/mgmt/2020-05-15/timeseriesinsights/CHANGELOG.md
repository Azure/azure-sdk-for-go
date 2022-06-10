# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Struct Fields

1. EventHubEventSourceMutableProperties.LocalTimestamp
1. EventSourceMutableProperties.LocalTimestamp
1. Gen1EnvironmentResourceProperties.CreationTime
1. Gen1EnvironmentResourceProperties.ProvisioningState
1. IoTHubEventSourceMutableProperties.LocalTimestamp

### Signature Changes

#### Funcs

1. EnvironmentsClient.Update
	- Params
		- From: context.Context, string, string, EnvironmentUpdateParameters
		- To: context.Context, string, string, BasicEnvironmentUpdateParameters
1. EnvironmentsClient.UpdatePreparer
	- Params
		- From: context.Context, string, string, EnvironmentUpdateParameters
		- To: context.Context, string, string, BasicEnvironmentUpdateParameters
1. EventSourcesClient.Update
	- Params
		- From: context.Context, string, string, string, EventSourceUpdateParameters
		- To: context.Context, string, string, string, BasicEventSourceUpdateParameters
1. EventSourcesClient.UpdatePreparer
	- Params
		- From: context.Context, string, string, string, EventSourceUpdateParameters
		- To: context.Context, string, string, string, BasicEventSourceUpdateParameters

## Additive Changes

### New Constants

1. IngressStartAtType.CustomEnqueuedTime
1. IngressStartAtType.EarliestAvailable
1. IngressStartAtType.EventSourceCreationTime
1. KindBasicEnvironmentUpdateParameters.KindBasicEnvironmentUpdateParametersKindEnvironmentUpdateParameters
1. KindBasicEnvironmentUpdateParameters.KindBasicEnvironmentUpdateParametersKindGen1
1. KindBasicEnvironmentUpdateParameters.KindBasicEnvironmentUpdateParametersKindGen2
1. KindBasicEventSourceUpdateParameters.KindBasicEventSourceUpdateParametersKindEventSourceUpdateParameters
1. KindBasicEventSourceUpdateParameters.KindBasicEventSourceUpdateParametersKindMicrosoftEventHub
1. KindBasicEventSourceUpdateParameters.KindBasicEventSourceUpdateParametersKindMicrosoftIoTHub

### New Funcs

1. *AzureEventSourceProperties.UnmarshalJSON([]byte) error
1. *EventHubEventSourceCommonProperties.UnmarshalJSON([]byte) error
1. *EventHubEventSourceCreationProperties.UnmarshalJSON([]byte) error
1. *EventHubEventSourceResourceProperties.UnmarshalJSON([]byte) error
1. *EventSourceCommonProperties.UnmarshalJSON([]byte) error
1. *IoTHubEventSourceCommonProperties.UnmarshalJSON([]byte) error
1. *IoTHubEventSourceCreationProperties.UnmarshalJSON([]byte) error
1. *IoTHubEventSourceResourceProperties.UnmarshalJSON([]byte) error
1. *Operation.UnmarshalJSON([]byte) error
1. EnvironmentUpdateParameters.AsBasicEnvironmentUpdateParameters() (BasicEnvironmentUpdateParameters, bool)
1. EnvironmentUpdateParameters.AsEnvironmentUpdateParameters() (*EnvironmentUpdateParameters, bool)
1. EnvironmentUpdateParameters.AsGen1EnvironmentUpdateParameters() (*Gen1EnvironmentUpdateParameters, bool)
1. EnvironmentUpdateParameters.AsGen2EnvironmentUpdateParameters() (*Gen2EnvironmentUpdateParameters, bool)
1. EventHubEventSourceUpdateParameters.AsBasicEventSourceUpdateParameters() (BasicEventSourceUpdateParameters, bool)
1. EventHubEventSourceUpdateParameters.AsEventHubEventSourceUpdateParameters() (*EventHubEventSourceUpdateParameters, bool)
1. EventHubEventSourceUpdateParameters.AsEventSourceUpdateParameters() (*EventSourceUpdateParameters, bool)
1. EventHubEventSourceUpdateParameters.AsIoTHubEventSourceUpdateParameters() (*IoTHubEventSourceUpdateParameters, bool)
1. EventSourceUpdateParameters.AsBasicEventSourceUpdateParameters() (BasicEventSourceUpdateParameters, bool)
1. EventSourceUpdateParameters.AsEventHubEventSourceUpdateParameters() (*EventHubEventSourceUpdateParameters, bool)
1. EventSourceUpdateParameters.AsEventSourceUpdateParameters() (*EventSourceUpdateParameters, bool)
1. EventSourceUpdateParameters.AsIoTHubEventSourceUpdateParameters() (*IoTHubEventSourceUpdateParameters, bool)
1. Gen1EnvironmentUpdateParameters.AsBasicEnvironmentUpdateParameters() (BasicEnvironmentUpdateParameters, bool)
1. Gen1EnvironmentUpdateParameters.AsEnvironmentUpdateParameters() (*EnvironmentUpdateParameters, bool)
1. Gen1EnvironmentUpdateParameters.AsGen1EnvironmentUpdateParameters() (*Gen1EnvironmentUpdateParameters, bool)
1. Gen1EnvironmentUpdateParameters.AsGen2EnvironmentUpdateParameters() (*Gen2EnvironmentUpdateParameters, bool)
1. Gen2EnvironmentUpdateParameters.AsBasicEnvironmentUpdateParameters() (BasicEnvironmentUpdateParameters, bool)
1. Gen2EnvironmentUpdateParameters.AsEnvironmentUpdateParameters() (*EnvironmentUpdateParameters, bool)
1. Gen2EnvironmentUpdateParameters.AsGen1EnvironmentUpdateParameters() (*Gen1EnvironmentUpdateParameters, bool)
1. Gen2EnvironmentUpdateParameters.AsGen2EnvironmentUpdateParameters() (*Gen2EnvironmentUpdateParameters, bool)
1. IoTHubEventSourceUpdateParameters.AsBasicEventSourceUpdateParameters() (BasicEventSourceUpdateParameters, bool)
1. IoTHubEventSourceUpdateParameters.AsEventHubEventSourceUpdateParameters() (*EventHubEventSourceUpdateParameters, bool)
1. IoTHubEventSourceUpdateParameters.AsEventSourceUpdateParameters() (*EventSourceUpdateParameters, bool)
1. IoTHubEventSourceUpdateParameters.AsIoTHubEventSourceUpdateParameters() (*IoTHubEventSourceUpdateParameters, bool)
1. PossibleIngressStartAtTypeValues() []IngressStartAtType
1. PossibleKindBasicEnvironmentUpdateParametersValues() []KindBasicEnvironmentUpdateParameters
1. PossibleKindBasicEventSourceUpdateParametersValues() []KindBasicEventSourceUpdateParameters

### Struct Changes

#### New Structs

1. Dimension
1. IngressStartAtProperties
1. LogSpecification
1. MetricAvailability
1. MetricSpecification
1. OperationProperties
1. ServiceSpecification

#### New Struct Fields

1. AzureEventSourceProperties.*IngressStartAtProperties
1. AzureEventSourceProperties.LocalTimestamp
1. EnvironmentUpdateParameters.Kind
1. EventHubEventSourceCommonProperties.*IngressStartAtProperties
1. EventHubEventSourceCommonProperties.LocalTimestamp
1. EventHubEventSourceCreationProperties.*IngressStartAtProperties
1. EventHubEventSourceCreationProperties.LocalTimestamp
1. EventHubEventSourceResourceProperties.*IngressStartAtProperties
1. EventHubEventSourceResourceProperties.LocalTimestamp
1. EventHubEventSourceUpdateParameters.Kind
1. EventSourceCommonProperties.*IngressStartAtProperties
1. EventSourceCommonProperties.LocalTimestamp
1. EventSourceUpdateParameters.Kind
1. Gen1EnvironmentUpdateParameters.Kind
1. Gen2EnvironmentUpdateParameters.Kind
1. IoTHubEventSourceCommonProperties.*IngressStartAtProperties
1. IoTHubEventSourceCommonProperties.LocalTimestamp
1. IoTHubEventSourceCreationProperties.*IngressStartAtProperties
1. IoTHubEventSourceCreationProperties.LocalTimestamp
1. IoTHubEventSourceResourceProperties.*IngressStartAtProperties
1. IoTHubEventSourceResourceProperties.LocalTimestamp
1. IoTHubEventSourceUpdateParameters.Kind
1. Operation.*OperationProperties
1. Operation.Origin
