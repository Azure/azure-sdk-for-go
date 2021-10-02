set GOOS=linux
set GOARCH=amd64

pushd ..
go build -o ./stress/stress ./stress
if errorlevel 1 goto :EOF
popd

call docker build . -t stresstestregistry.azurecr.io/ripark/gosbtest:latest
if errorlevel 1 goto :EOF

call docker push stresstestregistry.azurecr.io/ripark/gosbtest:latest
if errorlevel 1 goto :EOF

kubectl replace -f job.yaml --force