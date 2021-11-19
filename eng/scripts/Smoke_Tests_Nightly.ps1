#Requires -Version 7.0

Write-Host $PSScriptRoot

# 1. Every module uses a replace directive to the local version
# 2. Include every module (data & mgmt) in a go.mod file
# 3. Run `go mod tidy` and ensure it succeeds

# Create new module directory
New-Item -Path "$PSScriptRoot/sdk/smoketests"
Push-Location "$PSScriptRoot/sdk/smoketests"
go mod init
Pop-Location

# From sdk directory, find all packages with go.mod file
Push-Location "$PSScriptRoot/sdk"

$modules = Get-ChildItem -Path . -Recurse -Include go.mod
$Paths = $("")
$ReplacePaths = $("")

foreach ($module in $modules) {
    Write-Host $module
    $Array = $module.ToString().Split("/sdk")
    if ($Array.Length -ne 2) {
        Write-Host "There was an error parsing the path of the module ($module)"
    }
    # Remove the '/go.mod' portion
    $Array[1] = $Array[1] -Replace "/go.mod"
    Write-Host "Array[1] $Array[1]"

    $path = "github.com/Azure/azure-sdk-for-go/sdk/$Array[1]"
    Write-Host $path
    $Paths += $path

    $replacePath = "replace $path => ../$Array[1]"
    Write-Host $replacePath
    $ReplacePaths += $replacePath
}

Pop-Location

# Add files to go.mod file in sdk/smoketests
Push-Location "$PSScriptRoot/sdk/smoketests"

Write-Host "Writing modules to go.mod at $PSScriptRoot/sdk/smoketests starting with replace directives"

foreach ($replace in $ReplacePaths) {
    Add-Content -Path ./go.mod "$replace`n"
}

Write-Host "Writing requires portion"
Add-Content -Path ./go.mod "require (`n"
foreach ($path in $Paths) {
    Add-Content -Path ./go.mod "`t$path`n"
}
Add-Content -Path ./go.mod ")"

# print to stdout the results of the go.mod addition
Get-Content ./go.mod

# Finally, run go mod tidy. This should pass
go mod tidy


Pop-Location