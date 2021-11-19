#Requires -Version 7.0

Write-Host "PWD: $pwd"
$rootDir = $pwd
$smokeTestsDir = (Join-Path $pwd "sdk" "smoketests").ToString()

# 1. Every module uses a replace directive to the local version
# 2. Include every module (data & mgmt) in a go.mod file
# 3. Run `go mod tidy` and ensure it succeeds

# Create new module directory
New-Item $smokeTestsDir -ItemType Directory
Push-Location $smokeTestsDir
go mod init
Pop-Location

# From sdk directory, find all packages with go.mod file
Push-Location "$rootDir/sdk"

$modules = Get-ChildItem -Path . -Recurse -Include go.mod
$Paths = $("")
$ReplacePaths = $("")

foreach ($module in $modules) {
    Write-Host $module
    $Array = $module.ToString().Split("\sdk\")
    Write-Host $Array $Array.Length
    Write-Host $Array.ToString()
    if ($Array.Length -ne 2) {
        Write-Host "There was an error parsing the path of the module ($module)"
        break
    }

    $secondPart = "sdk/" + ($Array[1] -Replace "\\go.mod")
    Write-Host "`nSECOND PART: " $secondPart "`n`n"

    $path = "github.com/Azure/azure-sdk-for-go/sdk/$secondPart"
    Write-Host $path
    $Paths += ($path + "@latest")

    $replacePath = "replace $path => ../$secondPart"
    Write-Host "`n REPLACE PATH: $replacePath`n"
    $ReplacePaths += $replacePath

    break
}

Pop-Location

# Add files to go.mod file in sdk/smoketests
Push-Location "$rootDir/sdk/smoketests"

Write-Host "Writing modules to go.mod at $rootDir/sdk/smoketests starting with replace directives"

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