// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/pipelines"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/test"
)

var (
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
	err = writeMgmtFile(mgmtReport, path.Join(sdkPath, "mgmtReport.md"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Upload mgmt report to storage account...")
	err = uploadMgmtReport(mgmtReport, storageAccountName, storageAccountKey, storageContainerName, containerBlobName)
	if err != nil {
		log.Fatal(err)
	}
}

func execute(sdkPath, personalAccessToken string) (map[string]mgmtInfo, error) {
	conn := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	testClient, err := test.NewClient(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	buildClient, err := build.NewClient(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	azureDevopsClient := azuredevops.NewClient(conn, organizationUrl)

	pipelineClient := pipelines.NewClient(context.Background(), conn)
	pipelinesList, err := pipelineClient.ListPipelines(context.Background(), pipelines.ListPipelinesArgs{
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
				log.Printf("%s-%s\n", dir.Name(), arm.Name())
				// read autorest.md
				tag, version, err := readAutorestMD(filepath.Join(modulePath, dir.Name(), arm.Name(), "autorest.md"))
				if err != nil {
					return nil, err
				}

				p, ok := pipelinesMap[fmt.Sprintf("go - %s - %s", arm.Name(), dir.Name())]
				if !ok {
					p, ok = pipelinesMap[fmt.Sprintf("go - %s", arm.Name())]
				}
				if ok {
					mInfo := mgmtInfo{
						version: version,
						tag:     tag,
					}

					// code coverage
					buildId, err := codeCoverage(pipelineClient, azureDevopsClient, &mInfo, *p.Id)
					if err != nil {
						return nil, err
					}

					// mock test
					err = mockTest(testClient, &mInfo, buildId)
					if err != nil {
						return nil, err
					}

					// live test
					err = liveTest(buildClient, &mInfo, buildId)
					if err != nil {
						return nil, err
					}

					moduleName := fmt.Sprintf("%s/%s", dir.Name(), arm.Name())
					mgmtReport[moduleName] = mInfo
				}
			}
		}
	}

	return mgmtReport, nil
}

func readAutorestMD(path string) (string, string, error) {
	var (
		tag           string
		readmeLink    string
		moduleVersion string
		multiModule   string
	)

	b, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
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

	if tag == "" && readmeLink != "" {
		readmeLink = strings.ReplaceAll(readmeLink, "https://github.com", "https://raw.githubusercontent.com")
		readmeLink = strings.ReplaceAll(readmeLink, "/blob", "")
		resp, err := http.Get(readmeLink)
		if err != nil {
			return "", "", err
		}

		readmeBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
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

	return tag, moduleVersion, nil
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

func codeCoverage(pipelineClient pipelines.Client, azureDevopsClient *azuredevops.Client, info *mgmtInfo, pid int) (*int, error) {
	listRuns, err := pipelineClient.ListRuns(context.Background(), pipelines.ListRunsArgs{Project: &projectName, PipelineId: &pid})
	if err != nil && len(*listRuns) > 0 {
		return nil, err
	}

	var buildId *int
	for i := 0; i < 5 && i < len(*listRuns); i++ {
		buildId = (*listRuns)[i].Id

		// code coverage
		buildCodeCoverage, err := getBuildCodeCoverage(azureDevopsClient, projectName, *buildId)
		if err != nil {
			return nil, err
		}

		for _, coverage := range *buildCodeCoverage.CoverageData {
			if len(*coverage.CoverageStats) > 0 {
				info.CoveredLines = *(*coverage.CoverageStats)[0].Covered
				info.CoverableLines = *(*coverage.CoverageStats)[0].Total
				break
			}
		}

		if info.CoverableLines == 0 {
			continue
		}
		break
	}
	return buildId, nil
}

func mockTest(testClient test.Client, info *mgmtInfo, buildId *int) error {
	buildUri := fmt.Sprintf("vstfs:///Build/Build/%d", *buildId)
	testRuns, err := testClient.GetTestRuns(context.Background(), test.GetTestRunsArgs{
		Project:  &projectName,
		BuildUri: &buildUri,
	})
	if err != nil {
		return err
	}
	if len(*testRuns) > 0 {
		info.mockTestPass = *(*testRuns)[0].PassedTests
		info.mockTestTotal = *(*testRuns)[0].TotalTests
	}
	return nil
}

func liveTest(buildClient build.Client, info *mgmtInfo, buildId *int) error {
	buildLogs, err := buildClient.GetBuildLogs(context.Background(), build.GetBuildLogsArgs{
		Project: &projectName,
		BuildId: buildId,
	})
	if err != nil {
		return err
	}

	for i := 120; i < len(*buildLogs); i++ {
		logLines, err := buildClient.GetBuildLogLines(context.Background(), build.GetBuildLogLinesArgs{
			Project: &projectName,
			BuildId: buildId,
			LogId:   (*buildLogs)[i].Id,
		})
		if err != nil {
			return err
		}

		if logInfo := strings.Join(*logLines, "\n"); strings.Contains(logInfo, "Starting: Run Tests") && strings.Contains(logInfo, "Finishing: Run Tests") {
		loop:
			for _, line := range *logLines {
				if strings.Contains(line, "coverage:") {
					splits := strings.Split(line, " ")
					for _, j := range splits {
						if strings.Contains(j, "%") {
							info.liveTestCoverage = j
							break loop
						}
					}
				}
			}
			break
		}
	}

	return nil
}

func defaultPlaceholder(v string) string {
	if v == "" || v == "0.0%" {
		return "/"
	}
	return v
}

var mgmtReportMDHeader = `|module | latest version | tag | live test coverage | mock test result | mock test coverage |
|---|---|---|---|---|---|
`

func writeMgmtFile(mgmtReport map[string]mgmtInfo, path string) error {
	mgmtFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer mgmtFile.Close()

	_, err = mgmtFile.Write([]byte(mgmtReportMDHeader))
	if err != nil {
		return err
	}

	sortMgmt := make([]string, 0, len(mgmtReport))
	for k := range mgmtReport {
		sortMgmt = append(sortMgmt, k)
	}
	sort.Strings(sortMgmt)

	for _, module := range sortMgmt {
		m := mgmtReport[module]
		mt := "/"
		if m.mockTestTotal != 0 {
			mt = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.mockTestPass)/float64(m.mockTestTotal)*100, m.mockTestPass, m.mockTestTotal)
		}

		cc := "/"
		if m.CoverableLines != 0 {
			cc = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.CoveredLines)/float64(m.CoverableLines)*100, m.CoveredLines, m.CoverableLines)
		}

		f := fmt.Sprintf("|%s | %s | %s | %s | %s | %s |\n", module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), defaultPlaceholder(m.liveTestCoverage), mt, cc)
		_, err = mgmtFile.Write([]byte(f))
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadMgmtReport(mgmtReport map[string]mgmtInfo, accountName, accountKey, containerName, blobName string) error {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return err
	}

	sortMgmt := make([]string, 0, len(mgmtReport))
	for k := range mgmtReport {
		sortMgmt = append(sortMgmt, k)
	}
	sort.Strings(sortMgmt)

	htmlData := make([]string, 0, len(mgmtReport))
	average := []struct {
		count int
		sum   float64
	}{
		{
			count: 0,
			sum:   0,
		},
		{
			count: 0,
			sum:   0,
		},
		{
			count: 0,
			sum:   0,
		},
	}

	for i, module := range sortMgmt {
		m := mgmtReport[module]
		mt := "/"
		if m.mockTestTotal != 0 {
			coverage := float64(m.mockTestPass) / float64(m.mockTestTotal)
			average[0].sum += coverage
			average[0].count++
			mt = fmt.Sprintf("%.2f%%(%d/%d)", coverage*100, m.mockTestPass, m.mockTestTotal)
		}

		cc := "/"
		if m.CoverableLines != 0 {
			coverage := float64(m.CoveredLines) / float64(m.CoverableLines)
			average[1].sum += coverage
			average[1].count++
			cc = fmt.Sprintf("%.2f%%(%d/%d)", coverage*100, m.CoveredLines, m.CoverableLines)
		}

		lt := defaultPlaceholder(m.liveTestCoverage)
		if lt != "/" {
			f, err := strconv.ParseFloat(strings.TrimRight(lt, "%"), 64)
			if err != nil {
				return err
			}
			average[2].sum += f
			average[2].count++
		}

		tdBackground := ""
		if i%2 == 0 {
			tdBackground = tdBackgroundStyle
		}
		htmlData = append(htmlData, fmt.Sprintf(htmlTR, tdBackground, module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), defaultPlaceholder(m.liveTestCoverage), mt, cc))
	}

	// average
	htmlData = append(htmlData, fmt.Sprintf(htmlTR, "", "Average", "", "", fmt.Sprintf("%.2f%%", average[2].sum/float64(average[2].count)), fmt.Sprintf("%.2f%%", (average[0].sum/float64(average[0].count))*100), fmt.Sprintf("%.2f%%", (average[1].sum/float64(average[1].count))*100)))

	// parse template file
	t, err := template.ParseFiles("./mgmtreport.tpl")
	if err != nil {
		log.Fatal(err)
	}

	w := bytes.Buffer{}
	err = t.Execute(&w, htmlData)
	if err != nil {
		log.Fatal(err)
	}

	contentType := "text/html"
	_, err = client.UploadStream(context.TODO(),
		containerName,
		blobName,
		&w,
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

var htmlTR = `
			<tr%s>
				<td align="left">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
			</tr>`

var tdBackgroundStyle = ` class="pure-table-odd"`
