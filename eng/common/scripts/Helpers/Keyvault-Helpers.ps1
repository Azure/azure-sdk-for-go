$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true

function downloadKeyvaultSecretsAzCli([string]$keyvault) {
    $result = @{}
    $secrets = @((az keyvault secret list --vault-name $keyvault -o json `
                | ConvertFrom-Json -AsHashtable).name)
    foreach ($key in $secrets) {
        Write-Host "Fetching secret $key"
        $value = az keyvault secret show --vault-name $keyvault --name $key --query value -o tsv
        $result[$key] = $value
    }
    return $result
}

function downloadKeyvaultSecretsAzPowershell([string]$keyvault) {
    $initialContext = if (!$CI) { Get-AzContext }

    $result = @{}
    try {
        $secrets = @(Get-AzKeyvaultSecret -VaultName $keyvault)
        foreach ($obj in $secrets) {
            Write-Host "Fetching secret $($obj.Name)"
            $value = Get-AzKeyvaultSecret -VaultName $keyvault -Name $obj.Name -AsPlainText
            $result[$obj.name] = $value
        }
    } finally {
        if ($initialContext) {
            Write-Verbose "Restoring initial context: $($initialContext.Account)"
            $null = $initialContext | Select-AzContext
        }
    }

    return $result
}

function Export-KeyvaultSecrets {
    param(
        [string]$Keyvault,
        [string]$Prefix,
        [switch]$AzCli,
        [switch]$CI = ($null -ne $env:SYSTEM_TEAMPROJECTID)
    )

    $secrets = if ($AzCli) {
        downloadKeyvaultSecretsAzCli $Keyvault
    } else {
        downloadKeyvaultSecretsAzPowershell $Keyvault
    }

    foreach ($secret in $secrets.GetEnumerator()) {
        $envKey = ($secret.key -replace '-','_').ToUpper()
        if (!$prefix -or $secret.name -like "${prefix}*" -or $envKey -like "${prefix}*") {
            if ($envKey -eq $secret.key) {
                Write-Host "Setting variable '$envKey' from keyvault '$Keyvault' as secret"
            } else {
                Write-Host "Setting variable '$envKey' from '$($secret.key)' in keyvault '$Keyvault' as secret"
            }
            [Environment]::SetEnvironmentVariable($envKey, $secret.value)
            if ($CI) {
                Write-Host "##vso[task.setvariable variable=$envKey;issecret=true;]$($secret.value)"
            }
        }
    }
}
