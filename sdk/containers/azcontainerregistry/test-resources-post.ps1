# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

# IMPORTANT: Do not invoke this file directly. Please instead run eng/common/TestResources/New-TestResources.ps1 from the repository root.

param (
    [hashtable] $DeploymentOutputs,

    [Parameter(ParameterSetName = 'Provisioner', Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string] $TenantId,

    [Parameter()]
    [ValidatePattern('^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$')]
    [string] $TestApplicationId,

    # Captures any arguments from eng/New-TestResources.ps1 not declared here (no parameter errors).
    [Parameter(ValueFromRemainingArguments = $true)]
    $RemainingArguments
)

if ($CI) {
    az login --federated-token $env:OIDC_TOKEN --service-principal -t $TenantId -u $TestApplicationId
    az account set --subscription $DeploymentOutputs['AZCONTAINERREGISTRY_SUBSCRIPTION_ID']
}
