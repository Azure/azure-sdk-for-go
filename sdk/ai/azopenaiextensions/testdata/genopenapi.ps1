Push-Location ./testdata

if (Test-Path -Path "TempTypeSpecFiles") {
    Remove-Item -Recurse -Force TempTypeSpecFiles
}

npm install

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

npm run pull

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

npm run build

if ($LASTEXITCODE -ne 0) {
    Exit 1
}

Pop-Location
