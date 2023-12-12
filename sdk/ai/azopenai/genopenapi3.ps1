Push-Location ./testdata
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
