variable "location" {
  description = "Azure datacenter to deploy to."
  default = "westus2"
}

variable "servicebus_name_prefix" {
  description = "Input your unique Azure Service Bus Namespace name"
  default = "azuresbtests"
}

variable "resource_group_name" {
  description = "Resource group to provision test infrastructure in."
  default = "servicebus-go-tests"
}

resource "random_string" "name" {
  keepers = {
    # Generate a new id each time we switch to a new resource group
    group_name = "${var.resource_group_name}"
  }

  length  = 8
  upper   = false
  special = false
  number  = false
}

resource "azurerm_resource_group" "test" {
  name = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_servicebus_namespace" "test" {
  name = "${var.servicebus_name_prefix}-${random_string.name.result}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku = "standard"
}

# Most tests should create and destroy their own Queues, Topics, and Subscriptions. However, to keep examples from being
# bloated, the items below are created externally by Terraform.

resource "azurerm_servicebus_queue" "scheduledMessages" {
  name = "scheduledmessages"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name = "${azurerm_servicebus_namespace.test.name}"
}

resource "azurerm_servicebus_queue" "queueSchedule" {
  name = "schedulewithqueue"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name = "${azurerm_servicebus_namespace.test.name}"
}

resource "azurerm_servicebus_queue" "helloworld" {
  name = "helloworld"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name = "${azurerm_servicebus_namespace.test.name}"
}

resource "azurerm_servicebus_queue" "receiveSession" {
  name = "receivesession"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name = "${azurerm_servicebus_namespace.test.name}"
  default_message_ttl = "PT300S"
  requires_session = true
}

# Data resources used to get SubID and Tennant Info
data "azurerm_client_config" "current" {}

output "TEST_SERVICEBUS_RESOURCE_GROUP" {
  value = "${var.resource_group_name}"
}

output "SERVICEBUS_CONNECTION_STRING" {
  value = "Endpoint=sb://${azurerm_servicebus_namespace.test.name}.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=${azurerm_servicebus_namespace.test.default_primary_key}"
  sensitive = true
}

output "AZURE_SUBSCRIPTION_ID" {
  value = "${data.azurerm_client_config.current.subscription_id}"
}

output "TEST_SERVICEBUS_LOCATION" {
  value = "${var.location}"
}

output "AZURE_TENANT_ID" {
  value = "${data.azurerm_client_config.current.tenant_id}"
}
