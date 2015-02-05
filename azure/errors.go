package azure

const (
	publishSettingsConfigurationError = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."

	managementCertificateConfigurationError = "ManagementCertificatePath is set. Consequently SubscriptionId must also be set, and PublishSettingsFilePath must not be set."

	paramNotSpecifiedError = "Parameter %s is not specified."
)
