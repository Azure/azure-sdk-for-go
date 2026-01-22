Push-Location ./testdata

if (Test-Path -Path "TempTypeSpecFiles") {
    Remove-Item -Recurse -Force TempTypeSpecFiles
}

tsp-client install-dependencies .

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

npm install -D "@azure-tools/typespec-autorest@0.64.0"

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

npm install -D "@typespec/openapi3@1.8.0"

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

tsp-client sync

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

tsp compile ./TempTypeSpecFiles/OpenAI.Inference --config ./tspconfig.yaml

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

Pop-Location
