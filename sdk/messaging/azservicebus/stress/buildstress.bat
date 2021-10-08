set GOOS=linux
set GOARCH=amd64

IMAGE_NAME="stresstestregistry.azurecr.io/<your user>/gosbtest:latest"

pushd ..
go build -o ./stress/stress ./stress
if errorlevel 1 goto :EOF
popd

call az acr login -n stresstestregistry
if errorlevel 1 goto :EOF

call docker build . -t %IMAGE_NAME%
if errorlevel 1 goto :EOF

call docker push %IMAGE_NAME%
if errorlevel 1 goto :EOF

kubectl replace -f job.yaml --force