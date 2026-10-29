# Platform Validation recorded live tests

These tests follow the [Azure SDK for Go testing guide](https://github.com/Azure/azure-sdk-for-go/blob/main/documentation/development/testing.md):
`testutil.StartProxy`, per-test recording sessions, recorded HTTP transports,
`testutil.GetCredAndClientOptions`, `testutil.PollForTest`, and `recording.Sleep`.
The generated SDK is not modified.

## Record against a live service

Sign in with `az login`. From this module's directory, set the following environment
variables, or load them from a private `.env` file:

```powershell
$env:AZURE_SUBSCRIPTION_ID = "<approved-subscription>"
$env:LOCATION = "eastus2euap"
$env:PLATFORMVALIDATION_ARM_ENDPOINT = "https://eastus2euap.management.azure.com"
# Inject PLATFORMVALIDATION_SOURCE_VHD_URI privately; never commit a signed URL.
$env:AZURE_RECORD_MODE = "record"
go test -count=1 -timeout 75m -v ./...
```

The subscription must have `Microsoft.PlatformValidation` registered and permission
to create resource groups and validation resources. The boot scenario requires an
accessible Linux Gen1 x64 VHD and uses `Standard_D4s_v3`. Allow sufficient validity
on the image URL for provisioning and execution. For public ARM, omit
`PLATFORMVALIDATION_ARM_ENDPOINT`.

The lifecycle test creates a unique, tagged resource group and CloudValidation,
checks CV and execution-plan create/get/list/update, runs the boot validation,
and verifies the run summary and individual test result. Cleanup deletes children
before parents and verifies the service-managed group is absent before deleting
the empty test group. Unknown mutation outcomes or active workloads are retained;
use the logged owned resource names for investigation rather than blindly
deleting a parent.

The three catalog tests check GET/LIST consistency and subscription-scoped response
IDs. They intentionally fail if the service returns provider-global IDs instead;
successful catalog requests do not establish that their response IDs satisfy the
declared contract.

## Replay without Azure credentials

Unset live-only variables, then run:

```powershell
$env:AZURE_RECORD_MODE = "playback"
go test -count=1 -timeout 5m -v ./...
```

The test framework restores recordings using the committed `assets.json` tag.
Playback uses a fake credential, sanitized image configuration, and recorded
resource names. To bypass recording against the live service, set
`AZURE_RECORD_MODE=live`.

Before updating recordings, inspect them for secrets, then publish the sanitized
assets using `test-proxy push -a assets.json`. Commit the updated `assets.json`
tag, not the recordings themselves.
