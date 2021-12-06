#Requires -Version 7.0

function Invoke-MgmtTestgen ()
{
    param (
        [string]$sdkDirectory = "",
        [switch]$clean,
        [switch]$vet,
        [switch]$generateExample,
        [switch]$generateMockTest,
        [switch]$skipBuild,
        [switch]$cleanGenerated,
        [switch]$format,
        [switch]$tidy,
        [string]$autorestPath = "",
        [string]$config = "autorest.md",
        [string]$autorestVersion = "3.6.2",
        [string]$goExtension = "@autorest/go@4.0.0-preview.31",
        [string]$testExtension = "@autorest/gotest@1.1.2",
        [string]$outputFolder
    )
    if ($clean)
    {
        Write-Host "##[command]Executing go clean -v ./... in " $sdkDirectory
        go clean -v ./...
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($cleanGenerated)
    {
        Write-Host "##[command]Cleaning auto-generated files in" $sdkDirectory
        Remove-Item "ze_generated_*"
        Remove-Item "zt_generated_*"
    }

    if ($generateExample -or $generateMockTest)
    {
        Write-Host "##[command]Executing autorest.gotest in " $sdkDirectory

        if ($autorestPath -eq "") {
            $autorestPath = "./" + $config
        }
        
        
        if ($outputFolder -eq '')
        {
            $outputFolder = $sdkDirectory
        }
        $exampleFlag = "false"
        if ($generateExample)
        {
            $exampleFlag = "true"
        }
        $mockTestFlag = "true"
        if (-not $generateMockTest)
        {
            $mockTestFlag = "false"
        }
        Write-Host "autorest --version=$autorestVersion --use=$goExtension --use=$testExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --generate-sdk=false --testmodeler.generate-mock-test=$mockTestFlag --testmodeler.generate-sdk-example=$exampleFlag $autorestPath"
        npx autorest --version=$autorestVersion --use=$goExtension --use=$testExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --generate-sdk=false --testmodeler.generate-mock-test=$mockTestFlag --testmodeler.generate-sdk-example=$exampleFlag $autorestPath
        if ($LASTEXITCODE)
        {
            Write-Host "##[error]Error running autorest.gotest"
            exit $LASTEXITCODE
        }
    }

    if ($format)
    {
        Write-Host "##[command]Executing gofmt -s -w . in " $sdkDirectory
        gofmt -s -w .
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($tidy)
    {
        Write-Host "##[command]Executing go mod tidy in " $sdkDirectory
        go mod tidy
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if (!$skipBuild)
    {
        Write-Host "##[command]Executing go build -x -v ./... in " $sdkDirectory
        go build -x -v ./...
        Write-Host "##[command]Build Complete!"
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($vet)
    {
        Write-Host "##[command]Executing go vet ./... in " $sdkDirectory
        go vet ./...
    }
}

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

function PrepareMockServer()
{
    param(
        [string]$MOCK_SERVER_NAME = "mock-server"
    )
    # install mock server
    $folder = Join-Path $env:TEMP "mock-service-host"
    StopMockServer -MOCK_SERVER_NAME $MOCK_SERVER_NAME
    try
    {
        Remove-Item -Recurse -Force -Path $folder
    }
    catch
    {
        Write-Host "Mock service host folder: $folder not existed"
    }
    New-Item -ItemType Directory -Force -Path $folder
    Set-Location $folder
    npm install @azure-tools/mock-service-host
}

function StartMockServer()
{
    param(
        [string]$MOCK_SERVER_NAME = "mock-server",
        [string]$MOCK_SERVER_READY = "validator initialized",
        [string]$MOCK_SERVER_WAIT_TIME = 600,
        [string]$specName,
        [string]$commitID
    )
    $folder = Join-Path $env:TEMP  "mock-service-host"
    Set-Location $folder

    # change .env file to use the specific swagger file
    $envFile = Join-Path $folder .env
    if ([string]::IsNullOrEmpty($specName))
    {
        $specName = "*"
    }
    New-Item -Path $envFile -ItemType File -Value '' -Force
    Add-Content $envFile "specRetrievalGitUrl=https://github.com/Azure/azure-rest-api-specs
specRetrievalGitBranch=main
specRetrievalGitCommitID=$commitID
validationPathsPattern=specification/$specName/resource-manager/**/*.json"

    # start mock server and check status
    Start-Job -Name $MOCK_SERVER_NAME -ScriptBlock { node node_modules/@azure-tools/mock-service-host/dist/src/main.js 2>&1 }
    $output = Receive-Job $MOCK_SERVER_NAME
    Write-Host "Mock sever status: `n $("$output")"
    $time = 0
    try
    {
        while ("$output" -notmatch $MOCK_SERVER_READY)
        {
            if ($time -gt $MOCK_SERVER_WAIT_TIME)
            {
                Write-Host "##[error] mock server start timeout"
                StopMockServer
                exit 1
            }
            Write-Host "Server not ready, wait for annother 10 seconds"
            $time += 10
            Start-Sleep -Seconds 10
            $output = Receive-Job $MOCK_SERVER_NAME
            Write-Host "Mock sever status: `n $("$output")"
        }
    }
    catch
    {
        Write-Host "##[error]wait for mock server start:`n$_"
        exit 1
    }
}

function StopMockServer()
{
    param(
        [string]$MOCK_SERVER_NAME = "mock-server"
    )
    try
    {
        Stop-Job -Name $MOCK_SERVER_NAME
    }
    catch
    {
        Write-Host "##[error]can not stop mock server:`n$_"
    }
}

function GetSwaggerInfo()
{
    param(
        [string]$dir,
        [string]$AUTOREST_CONFIG_FILE = "autorest.md"
    )
    Set-Location $dir
    $swaggerInfoRegex = "(?<commitID>[0-9a-f]{40})\/specification\/(?<specName>.*)\/resource-manager\/readme.md"
    try
    {
        $content = Get-Content .\$AUTOREST_CONFIG_FILE -Raw
        if ($content -match $swaggerInfoRegex)
        {
            return $matches["specName"], $matches["commitID"]
        }
    }
    catch
    {
        Write-Error "Error parsing swagger info"
        Write-Error $_
    }
    Write-Host "Cannot find swagger info"
    exit 1
}

function TestAndGenerateReport($dir)
{
    Set-Location $dir
    # dependencies for go coverage report generation
    go get github.com/jstemmer/go-junit-report
    go get github.com/axw/gocov/gocov
    go get github.com/AlekSi/gocov-xml
    go get github.com/matm/gocov-html

    # set azidentity env for mock test
    $Env:AZURE_TENANT_ID = "mock-test"
    $Env:AZURE_CLIENT_ID = "mock-test"
    $Env:AZURE_USERNAME = "mock-test"
    $Env:AZURE_PASSWORD = "mock-test"

    # do test with corage report and convert to cobertura format
    Write-Host "go cmd: go test -v -coverprofile coverage.txt | Tee-Object -FilePath output.txt"
    go test -v -coverprofile coverage.txt | Tee-Object -FilePath output.txt
    Write-Host "report.xml: cat output.txt | go-junit-report > report.xml"
    cat output.txt | go-junit-report > report.xml
    Write-Host "coverage.json: gocov convert ./coverage.txt > ./coverage.json"
    gocov convert ./coverage.txt > ./coverage.json
    Write-Host "coverage.xml: Get-Content ./coverage.json | gocov-xml > ./coverage.xml"
    Get-Content ./coverage.json | gocov-xml > ./coverage.xml
    Write-Host "coverage.html: Get-Content ./coverage.json | gocov-html > ./coverage.html"
    Get-Content ./coverage.json | gocov-html > ./coverage.html
}

function JudgeExitCode($errorMsg = "execution error")
{
    if (!$?)
    {
        Write-Host "##[error] $errorMsg"
        exit $LASTEXITCODE
    }
}

function ExecuteSingleTest($sdk, $runMock=$true)
{
    Write-Host "Start mock server"
    $swaggerInfo = GetSwaggerInfo $sdk.DirectoryPath
    if ($runMock -eq $true) {
        StartMockServer -specName $swaggerInfo[0] -commitID $swaggerInfo[1]
    }
    Write-Host "Execute mock test for $($sdk.Name)"
    TestAndGenerateReport $sdk.DirectoryPath
    Write-Host "Stop mock server"
    if ($runMock -eq $true) {
        StopMockServer
    }
}