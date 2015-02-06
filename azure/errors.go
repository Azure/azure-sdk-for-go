package azure

const (
	publishSettingsConfigurationError = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."

	managementCertificateConfigurationError = "Both ManagementCertificatePath and SubscriptionId should be set, and PublishSettingsFilePath must not be set."

	paramNotSpecifiedError = "Parameter %s is not specified."
)
