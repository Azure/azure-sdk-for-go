// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelines"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/test"
)

var (
	ctx     = context.Background()
	sdkPath string
	// azure
	storageAccountName   string
	storageContainerName string
	containerBlobName    string
	// azure devops
	organizationUrl string
	projectName     string
)

func init() {
	flag.StringVar(&sdkPath, "sdkpath", "", "azure-sdk-for-go path(required)")
	flag.StringVar(&storageAccountName, "storageaccount", "", "Azure Storage Account Name(required)")
	flag.StringVar(&storageContainerName, "storagecontainer", "$web", "Azure Storage Container Name")
	flag.StringVar(&containerBlobName, "storagecontainerblob", "mgmtReport.html", "Azure Storage Container Blob File Name")
	flag.StringVar(&organizationUrl, "organization", "https://dev.azure.com/azure-sdk", "Azure Devops Organization Url")
	flag.StringVar(&projectName, "project", "internal", "Azure Devops Project Name")
}

func main() {
	flag.Parse()

	if sdkPath == "" {
		flag.PrintDefaults()
		log.Fatal("Please enter the azure-sdk-for-go path")
	}
	if storageAccountName == "" {
		flag.PrintDefaults()
		log.Fatal("Please enter the Azure Storage Account Name")
	}
	storageAccountKey, ok := os.LookupEnv("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY")
	if !ok {
		log.Fatal("AZURE_STORAGE_PRIMARY_ACCOUNT_KEY could not be found")
	}

	personalAccessToken, ok := os.LookupEnv("AZURE_DEVOPS_PERSONAL_ACCESS_TOKEN")
	if !ok {
		log.Fatal("AZURE_DEVOPS_PERSONAL_ACCESS_TOKEN could not be found")
	}

	log.Printf("start running in %s...\n", sdkPath)
	startTime := time.Now()
	mgmtReport, err := execute(sdkPath, personalAccessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Statistical mgmt report time:", time.Since(startTime))

	log.Println("Write the mgmt report to the mgmtreport.md file...")
	htmlData, err := writeMgmtFile(mgmtReport, path.Join(sdkPath, "mgmtReport.md"))
	if err != nil {
		log.Fatal(err)
	}

	if htmlData != nil {
		log.Println("Upload mgmt report to storage account...")
		err = uploadMgmtReport(htmlData, storageAccountName, storageAccountKey, storageContainerName, containerBlobName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func execute(sdkPath, personalAccessToken string) (map[string]mgmtInfo, error) {
	conn := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	testClient, err := test.NewClient(ctx, conn)
	if err != nil {
		return nil, err
	}

	buildClient, err := build.NewClient(ctx, conn)
	if err != nil {
		return nil, err
	}

	azureDevopsClient := azuredevops.NewClient(conn, organizationUrl)

	pipelineClient := pipelines.NewClient(ctx, conn)
	pipelinesList, err := pipelineClient.ListPipelines(ctx, pipelines.ListPipelinesArgs{
		Project: &projectName,
	})
	if err != nil {
		return nil, err
	}

	// filter pipelineList
	pipelinesMap := make(map[string]*pipelines.Pipeline)
	for i, pipe := range *pipelinesList {
		if *pipe.Folder == "\\go" && !strings.Contains(*pipe.Name, "weekly") && strings.Contains(*pipe.Name, "go - arm") {
			pipelinesMap[*pipe.Name] = &(*pipelinesList)[i]
		}
	}

	// read all module
	sdkPath = strings.ReplaceAll(sdkPath, "\\", "/")
	modulePath := filepath.Join(sdkPath, "/sdk/resourcemanager")
	dirs, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	mgmtReport := make(map[string]mgmtInfo)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() != "internal" {
			armDirs, err := os.ReadDir(filepath.Join(modulePath, dir.Name()))
			if err != nil {
				return nil, err
			}

			for _, arm := range armDirs {
				moduleName := fmt.Sprintf("%s/%s", dir.Name(), arm.Name())
				// read autorest.md
				tag, readmeLink, version, multiModule, err := readAutorestMD(filepath.Join(modulePath, dir.Name(), arm.Name(), "autorest.md"))
				if err != nil {
					return nil, err
				}
				if tag == "" && readmeLink != "" {
					readmeLink = strings.ReplaceAll(readmeLink, "https://github.com", "https://raw.githubusercontent.com")
					readmeLink = strings.ReplaceAll(readmeLink, "/blob", "")
					resp, err := http.Get(readmeLink)
					if err != nil {
						return nil, err
					}

					readmeBody, err := io.ReadAll(resp.Body)
					if err != nil {
						return nil, err
					}

					if multiModule != "" {
						indexMultiModule := bytes.Index(readmeBody, []byte(fmt.Sprintf("``` yaml $(%s)", multiModule)))
						readmeBody = readmeBody[indexMultiModule:]
					}

					for _, line := range strings.Split(string(readmeBody), "\n") {
						if strings.HasPrefix(line, "tag: ") {
							tag = line[len("tag: "):]
							break
						}
					}
				}

				p, ok := pipelinesMap[fmt.Sprintf("go - %s", arm.Name())]
				if !ok {
					p, ok = pipelinesMap[fmt.Sprintf("go - %s - %s", arm.Name(), dir.Name())]
				}
				if ok {
					// mock test
					mockTestPass := 0
					mockTestTotal := 0
					listRuns, err := pipelineClient.ListRuns(ctx, pipelines.ListRunsArgs{Project: &projectName, PipelineId: p.Id})
					if err != nil && len(*listRuns) > 0 {
						return nil, err
					}

					buildId := (*listRuns)[0].Id
					top := 3
					queryRuns, err := testClient.QueryTestRuns(ctx, test.QueryTestRunsArgs{
						Project: &projectName,
						MinLastUpdatedDate: &azuredevops.Time{
							Time: time.Now().AddDate(0, 0, -3),
						},
						MaxLastUpdatedDate: &azuredevops.Time{
							Time: time.Now(),
						},
						BuildIds: &[]int{
							*buildId,
						},
						State: &test.TestRunStateValues.Completed,
						Top:   &top,
					})
					if err != nil {
						return nil, err
					}

					if queryRuns != nil && len(queryRuns.Value) > 0 {
						mockTestPass = *queryRuns.Value[0].PassedTests
						mockTestTotal = *queryRuns.Value[0].TotalTests
					}

					// live test
					buildLogs, err := buildClient.GetBuildLogs(ctx, build.GetBuildLogsArgs{
						Project: &projectName,
						BuildId: buildId,
					})
					if err != nil {
						return nil, err
					}

					liveTestCoverage := ""
					for i := 120; i < len(*buildLogs); i++ {
						logLines, err := buildClient.GetBuildLogLines(ctx, build.GetBuildLogLinesArgs{
							Project: &projectName,
							BuildId: buildId,
							LogId:   (*buildLogs)[i].Id,
						})
						if err != nil {
							return nil, err
						}

						if logInfo := strings.Join(*logLines, "\n"); strings.Contains(logInfo, "Starting: Run Tests") && strings.Contains(logInfo, "Finishing: Run Tests") {
						loop:
							for _, line := range *logLines {
								if strings.Contains(line, "coverage:") {
									splits := strings.Split(line, " ")
									for _, j := range splits {
										if strings.Contains(j, "%") {
											liveTestCoverage = j
											break loop
										}
									}
								}
							}
							break
						}
					}

					mInfo := mgmtInfo{
						version:          version,
						tag:              tag,
						mockTestPass:     mockTestPass,
						mockTestTotal:    mockTestTotal,
						liveTestCoverage: liveTestCoverage,
					}

					// code coverage
					codeCoverage, err := getBuildCodeCoverage(azureDevopsClient, projectName, *buildId)
					if err == nil && codeCoverage.CoverageData != nil {
						for _, coverage := range *codeCoverage.CoverageData {
							if len(*coverage.CoverageStats) > 0 {
								mInfo.CoveredLines = *(*coverage.CoverageStats)[0].Covered
								mInfo.CoverableLines = *(*coverage.CoverageStats)[0].Total
								break
							}
						}
					}

					mgmtReport[moduleName] = mInfo
				}
			}
		}
	}

	return mgmtReport, nil
}

func readAutorestMD(path string) (string, string, string, string, error) {
	var (
		tag           string
		readmeLink    string
		moduleVersion string
		multiModule   string
	)

	b, err := os.ReadFile(path)
	if err != nil {
		return "", "", "", "", err
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.Contains(line, "tag:") {
			tag = line[len("tag:"):]
		} else if strings.Contains(line, "readme.md") {
			readmeLink = line[len("- "):]
		} else if strings.Contains(line, "module-version: ") {
			moduleVersion = line[len("module-version: "):]
		} else if strings.Contains(line, "package-") && strings.Contains(line, ": true") {
			multiModule = line[:len(line)-len(": true")]
		}
	}

	return tag, readmeLink, moduleVersion, multiModule, nil
}

type mgmtInfo struct {
	tag              string
	version          string
	mockTestPass     int
	mockTestTotal    int
	CoveredLines     int
	CoverableLines   int
	liveTestCoverage string
}

func defaultPlaceholder(v string) string {
	if v == "" || v == "0.0%" {
		return "/"
	}
	return v
}

func writeMgmtFile(mgmtReport map[string]mgmtInfo, path string) ([]string, error) {
	mgmtFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer mgmtFile.Close()

	_, err = mgmtFile.Write([]byte(mgmtReportMDHeader))
	if err != nil {
		return nil, err
	}

	sortMgmt := make([]string, 0, len(mgmtReport))
	for k := range mgmtReport {
		sortMgmt = append(sortMgmt, k)
	}
	sort.Strings(sortMgmt)

	htmlData := make([]string, 0, len(mgmtReport))

	for i, module := range sortMgmt {
		m := mgmtReport[module]
		mockTest := "/"
		if m.mockTestTotal != 0 {
			mockTest = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.mockTestPass)/float64(m.mockTestTotal)*100, m.mockTestPass, m.mockTestTotal)
		}

		codeCoverage := "/"
		if m.CoverableLines != 0 {
			codeCoverage = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.CoveredLines)/float64(m.CoverableLines)*100, m.CoveredLines, m.CoverableLines)
		}

		f := fmt.Sprintf("|%s | %s | %s | %s | %s | %s |\n", module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), defaultPlaceholder(m.liveTestCoverage), mockTest, codeCoverage)
		_, err = mgmtFile.Write([]byte(f))
		if err != nil {
			return nil, err
		}

		tdBackground := ""
		if i%2 == 0 {
			tdBackground = tdBackgroundStyle
		}
		htmlData = append(htmlData, fmt.Sprintf(htmlTR, tdBackground, module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), defaultPlaceholder(m.liveTestCoverage), mockTest, codeCoverage))
	}

	return htmlData, nil
}

func uploadMgmtReport(htmlData []string, accountName, accountKey, containerName, blobName string) error {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return err
	}

	contentType := "text/html"
	_, err = client.UploadStream(context.TODO(),
		containerName,
		blobName,
		strings.NewReader(htmlHeader+strings.Join(htmlData, "\n")+htmlTail),
		&azblob.UploadStreamOptions{
			HTTPHeaders: &blob.HTTPHeaders{
				BlobContentType: &contentType,
			},
		})
	if err != nil {
		return err
	}

	return nil
}
