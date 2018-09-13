variable "location" {
  description = "Azure datacenter to deploy to."
  default = "westus2"
}

variable "servicebus_name" {
  description = "Input your unique Azure Service Bus Namespace name"
}

variable "resource_group_name" {
  description = "Resource group to provision test infrastructure in."
  default = "servicebus-go-tests"
}

resource "azurerm_resource_group" "test" {
  name = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_servicebus_namespace" "test" {
  name = "${var.servicebus_name}"
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

resource "azurerm_servicebus_queue" "helloworld" {
  name = "helloworld"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name = "${azurerm_servicebus_namespace.test.name}"
}