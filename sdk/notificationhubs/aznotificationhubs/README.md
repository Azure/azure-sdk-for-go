# Azure Notification Hubs Client Module for Go

Azure Notification Hubs provide a scaled-out push engine that enables you to send notifications to any platform (Apple, Amazon Kindle, Android, Baidu, Xiaomi, Web, Windows, etc.) from any back-end (cloud or on-premises). Notification Hubs works well for both enterprise and consumer scenarios. Here are a few example scenarios:

- Send breaking news notifications to millions with low latency.
- Send location-based coupons to interested user segments.
- Send event-related notifications to users or groups for media/sports/finance/gaming applications.
- Push promotional contents to applications to engage and market to customers.
- Notify users of enterprise events such as new messages and work items.
- Send codes for multi-factor authentication.

Key links:

- [Source code](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/notificationhubs/aznotificationhubs/)
- [API Reference Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs)
- [Product documentation](https://docs.microsoft.com/azure/notification-hubs/)

## Getting started

### Install the package

Install the Azure Event Hubs client module for Go with `go get`:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs
```

### Prerequisites

- Go, version 1.18 or higher
- An [Azure subscription](https://azure.microsoft.com/free/)
- An [Azure Notification Hubs namespace](https://docs.microsoft.com/azure/notification-hubs/).

An Azure Notification Hub can be created using the following methods:

1. [Azure Portal](https://docs.microsoft.com/azure/notification-hubs/create-notification-hub-portal)
2. [Azure CLI](https://docs.microsoft.com/azure/notification-hubs/create-notification-hub-azure-cli)
3. [Bicep](https://docs.microsoft.com/azure/notification-hubs/create-notification-hub-bicep)
4. [ARM Template](https://docs.microsoft.com/azure/notification-hubs/create-notification-hub-template)

Once created, the Notification Hub can be configured using the [Azure Portal or Azure CLI](https://docs.microsoft.com/azure/notification-hubs/configure-notification-hub-portal-pns-settings?tabs=azure-portal).

### Authenticate the client

Interaction with an Azure Notification Hub starts with the `NotificationHubClient` which supports [Shared Access Signature connection strings](https://docs.microsoft.com/azure/notification-hubs/notification-hubs-push-notification-security).  This includes the following permission levels: **Listen**, **Manage**, **Send**.

Listen allows for a client to register itself via the Registration and Installations API. Send allows for the client to send notifications to devices using the send APIs. Finally, Manage allows the user to do Registration and Installation management, such as queries.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}
```

## Key concepts

Once the `NotificationHubClient` has been initialized, the following concepts can be explored.

- Device Management via Installations and RegistrationDescriptions
- Send Notifications to Devices

### Device Management

Device management is a core concept to Notification Hubs to be able to store the unique identifier from the native Platform Notification Service (PNS) such as APNs or Firebase, and associated metadata such as tags used for sending push notifications to audiences.  This is done through a concept called Installations.

#### Installations API

Installations are a native JSON approach to device management that contains additional properties such as an installation ID and user ID which can be used for sending to audiences.  The installations API allows for flexibility including the following:

- Fully idempotent API so calling create on the installation, so an operation can be retried without worries about duplications.
- Support for `userId` and `installationId` properties which can be then used in tag expressions such as `$InstallationId:{myInstallId}` and `$UserId:{bob@contoso.com}`.
- Templates are now part of the installation instead of a separate registration and can be reference by name as a tag for sending.
- Partial updates are supported through the [JSON Patch Standard](https://tools.ietf.org/html/rfc6902), which allows to add tags and change other data without having to first query the installation.

Installations can be created through the `CreateOrUpdateInstallation` method such as the following:

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

// Create an installation for APNs
installation := &aznotificationhubs.Installation{
    InstallationID: "my-install-id",
    Platform:       "apns",
    PushChannel:    "push-channel",
    Tags:           []string{"likes_hockey", "likes_football"},
}

// Create or overwrite the installation
installation, err = client.CreateOrUpdateInstallation(installation)
if err != nil {
    panic(err)
}
```

An update to an installation can be made through the JSON Patch schema such as adding a tag and a user ID using the `UpdateInstallation` method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

installationId := "my-install-id"
updates := []InstallationPatch{
    {Op: "add", Path: "/tags", Value: "likes_baseball"},
    {Op: "add", Path: "/userId", Value: "bob@contoso.com" },
}

// Update the installation
response, err := client.UpdateInstallation(installationId, updates)
if err != nil {
    panic(err)
}
```

To retrieve an existing installation, use the `GetInstallation` method with your existing unique installation ID.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

installationId := "my-install-id"

// Get the installation ID
installation, err := client.GetInstallation(installationId)
if err != nil {
    panic(err)
}
```

### Send Operations

Notification Hubs supports sending notifications to devices either directly using the unique PNS provided identifier, using tags for audience send, or a general broadcast to all devices.  Using the Standard SKU and above, [scheduled send](https://docs.microsoft.com/azure/notification-hubs/notification-hubs-send-push-notifications-scheduled) allows the user to schedule notifications up to seven days in advance.  All send operations return a Tracking ID and Correlation ID which can be used for Notification Hubs support cases.  

Raw JSON or XML strings can be sent to the send or scheduled send methods, or the notification builders can be used which helps construct messages per PNS such as APNs, Firebase, Baidu, ADM and WNS.  These builders will build the native message format so there is no guessing about which fields are available for each PNS.

#### Direct Send

To send directly a device, the user can send using the platform provided unique identifier such as APNs device token by calling the `SendDirectNotification` method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

// Create a notification for APNs
deviceToken := "device-token"

headers := make(map[string]string)
headers["apns-priority"] = "10"
headers["apns-push-type"] = "alert"

contentType := "application/json;charset=utf-8"
platform := "apple"

request := &NotificationRequest{
    Message:     messageBody,
    Headers:     headers,
    Platform:    platform,
    ContentType: contentType,
}

response, err := client.SendDirectNotification(request, deviceToken)
if err != nil {
    panic(err)
}
```

#### Audience Send

In addition to targeting a single device, a user can target multiple devices using tags.  These tags can be supplied as a list of tags, which then creates a tag expression to match registered devices, or via a tag expression which can then use Boolean logic to target the right audience.  For more information about tags and tags expressions, see [Routing and Tag Expressions](https://docs.microsoft.com/azure/notification-hubs/notification-hubs-tags-segment-push-message).

If you wish to create a tag expression from an array of tags, there is a Tag Expression Builder available with the `CreateTagExpression` method which creates an "or tag expression" from the tags.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

// Create a notification for APNs
headers := make(map[string]string)
headers["apns-priority"] = "10"
headers["apns-push-type"] = "alert"

contentType := "application/json;charset=utf-8"
platform := "apple"

request := &NotificationRequest{
    Message:     messageBody,
    Headers:     headers,
    Platform:    platform,
    ContentType: contentType,
}

tagExpression := "likes_hockey && likes_football"

// Send the notification
response, err := client.SendNotification(request, tagExpression")
if err != nil {
    panic(err)
}
```

#### Scheduled Send

Push notifications can be scheduled up to seven days in advance with Standard SKU namespaces and above using the `ScheduleNotification` method to send to devices with tags or a general broadcast.  This returns a notification ID which can be then used to cancel if necessary via the `CancelScheduledNotification` method.

```go
import "github.com/Azure/azure-sdk-for-go/sdk/notificationhubs/aznotificationhubs"

// Create a NotificationHubClient using a connection string
connectionString := "Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<key_name>;SharedAccessKey=<key_value>"
hubName := "<hub_name>"
client, err := aznotificationhubs.NewNotificationHubClient(connectionString, hubName)
if err != nil {
    panic(err)
}

// Create a notification for APNs
headers := make(map[string]string)
headers["apns-priority"] = "10"
headers["apns-push-type"] = "alert"

contentType := "application/json;charset=utf-8"
platform := "apple"

request := &NotificationRequest{
    Message:     messageBody,
    Headers:     headers,
    Platform:    platform,
    ContentType: contentType,
}

tagExpression := "likes_hockey && likes_football"

// Schedule eight hours from now
scheduleTime := time.Now().Add(time.Hour * 8)

// Send the notification
response, err := client.ScheduleNotification(request, tagExpression, scheduleTime)
if err != nil {
    panic(err)
}
```

## Contributing

For details on contributing to this repository, see the [contributing guide][azure_sdk_for_go_contributing].

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit [https://cla.microsoft.com](https://cla.microsoft.com).

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

### Additional Helpful Links for Contributors

Many people all over the world have helped make this project better.  You'll want to check out:

* [What are some good first issues for new contributors to the repo?](https://github.com/azure/azure-sdk-for-go/issues?q=is%3Aopen+is%3Aissue+label%3A%22up+for+grabs%22)
* [How to build and test your change][azure_sdk_for_go_contributing_developer_guide]
* [How you can make a change happen!][azure_sdk_for_go_contributing_pull_requests]
* Frequently Asked Questions (FAQ) and Conceptual Topics in the detailed [Azure SDK for Go wiki](https://github.com/azure/azure-sdk-for-go/wiki).

<!-- ### Community-->
### Reporting security issues and security bugs

Security issues and bugs should be reported privately, via email, to the Microsoft Security Response Center (MSRC) <secure@microsoft.com>. You should receive a response within 24 hours. If for some reason you do not, please follow up via email to ensure we received your original message. Further information, including the MSRC PGP key, can be found in the [Security TechCenter](https://www.microsoft.com/msrc/faqs-report-an-issue).

### License

Azure SDK for Go is licensed under the [MIT](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/template/aztemplate/LICENSE.txt) license.

<!-- LINKS -->
[azure_sdk_for_go_contributing]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md
[azure_sdk_for_go_contributing_developer_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#developer-guide
[azure_sdk_for_go_contributing_pull_requests]: https://github.com/Azure/azure-sdk-for-go/blob/main/CONTRIBUTING.md#pull-requests
