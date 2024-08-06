# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

# IMPORTANT: Do not invoke this file directly. Please instead run eng/common/TestResources/New-TestResources.ps1 from the repository root.

param (
    [hashtable] $DeploymentOutputs
)

Connect-AzContainerRegistry -Name $DeploymentOutputs['REGISTRY_NAME']
