function PullAzureRestAPISpecs() {
    npm run pull

    if ($LASTEXITCODE -ne 0) {
        Exit 1
    }
}

Push-Location ./testdata

if (Test-Path -Path "TempTypeSpecFiles") {
    Remove-Item -Recurse -Force TempTypeSpecFiles
}

npm install

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

# If you want to test with a local copy of your azure-rest-api-specs changes
# just comment out PullAzureRestAPISpecs and uncomment the lines below.
PullAzureRestAPISpecs

# $LocalTestingGitRepo = "<repo path>"
# Copy-Item `
#     -Path $LocalTestingGitRepo/specification/eventgrid/Azure.Messaging.EventGrid `
#     -Destination TempTypeSpecFiles/Azure.Messaging.EventGrid `
#     -Recurse
# Remove-Item TempTypeSpecFiles/Azure.Messaging.EventGrid/tspconfig.yaml

npm run build

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

Pop-Location

Clear-Host
