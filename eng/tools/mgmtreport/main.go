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
)

var (
	sdkPath string
	// storage account
	storageAccountName   string
	storageContainerName string
	containerBlobName    string
	// azure devops info
	organizationUrl string
	projectName     string
)

var mgmtReportMDHeader = `|module | latest version | tag | live test result | live test coverage (line) | live test coverage (operation) | mock test result | mock test coverage (line) | mock test coverage (operation) |
|---|---|---|---|---|---|---|---|---|
`

var htmlTR = `
			<tr%s>
				<td align="left">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
				<td align="center">%s</td>
			</tr>`

var tdBackgroundStyle = ` class="pure-table-odd"`

func init() {
	flag.StringVar(&sdkPath, "sdkpath", "", "SDK Repo Path(required)")
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
		log.Fatal("Please provide the SDK repo path")
	}

	if storageAccountName == "" {
		flag.PrintDefaults()
		log.Fatal("Please provide the Azure Storage account name")
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
	err := execute(sdkPath, personalAccessToken, storageAccountKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mgmt report statistic time:", time.Since(startTime))
}

func execute(sdkPath, personalAccessToken, storageAccountKey string) error {
	conn := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)

	buildClient, err := build.NewClient(context.Background(), conn)
	if err != nil {
		return err
	}

	azureDevopsClient := azuredevops.NewClient(conn, organizationUrl)

	pipelineClient := pipelines.NewClient(context.Background(), conn)
	pipelinesList, err := pipelineClient.ListPipelines(context.Background(), pipelines.ListPipelinesArgs{
		Project: &projectName,
	})
	if err != nil {
		return err
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
	modulePath := filepath.Join(sdkPath, "sdk", "resourcemanager")
	dirs, err := os.ReadDir(modulePath)
	if err != nil {
		return err
	}

	mgmtReport := make(map[string]mgmtInfo)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() != "internal" {
			armDirs, err := os.ReadDir(filepath.Join(modulePath, dir.Name()))
			if err != nil {
				return err
			}

			for _, arm := range armDirs {
				// read autorest.md
				tag, version, err := readAutorestMD(filepath.Join(modulePath, dir.Name(), arm.Name(), "autorest.md"))
				if err != nil {
					return err
				}

				p, ok := pipelinesMap[fmt.Sprintf("go - %s - %s", arm.Name(), dir.Name())]
				if !ok {
					p, ok = pipelinesMap[fmt.Sprintf("go - %s", arm.Name())]
				}
				if ok {
					log.Printf("Processing %s ...", arm.Name())
					mInfo := mgmtInfo{
						version: version,
						tag:     tag,
					}

					// code coverage
					buildId, err := getBuildId(pipelineClient, azureDevopsClient, *p.Id)
					if err != nil {
						return err
					}

					mockTestLogId, liveTestLogId, err := getLogID(buildClient, buildId)
					if err != nil {
						return err
					}

					// mock test
					mockTestTotal, mockTestPassed, mockTestCoverage, err := getTestResult(buildClient, buildId, mockTestLogId)
					if err != nil {
						return err
					}
					mInfo.mockTestTotal = mockTestTotal
					mInfo.mockTestPass = mockTestPassed
					mInfo.mockTestCoverage = mockTestCoverage

					// live test
					liveTestTotal, liveTestPassed, liveTestCoverage, err := getTestResult(buildClient, buildId, liveTestLogId)
					if err != nil {
						return err
					}
					mInfo.liveTestTotal = liveTestTotal
					mInfo.liveTestPass = liveTestPassed
					mInfo.liveTestCoverage = liveTestCoverage

					allOperation, skipped, err := getOperation(buildClient, buildId, mockTestLogId)
					if err != nil {
						return err
					}
					mInfo.mockTestSkip = skipped

					if mInfo.liveTestTotal != 0 {
						operations, err := getLiveTestOperationCalls(buildClient, allOperation, buildId, liveTestLogId)
						if err != nil {
							return err
						}
						mInfo.liveTestCallOperations = operations
					}

					moduleName := fmt.Sprintf("%s/%s", dir.Name(), arm.Name())
					mgmtReport[moduleName] = mInfo
				}
			}
		}
	}

	log.Println("generate a report in Markdown format...")
	err = generateMDReport(mgmtReport, filepath.Join(sdkPath, "mgmtReport.md"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("generate a report in HTML format...")
	err = generateHTMLReport(mgmtReport, filepath.Join(sdkPath, "mgmtReport.html"))
	if err != nil {
		return err
	}

	log.Println("upload mgmt report to cloud...")
	err = uploadHTMLReport(filepath.Join(sdkPath, "mgmtReport.html"), storageAccountName, storageAccountKey, storageContainerName, containerBlobName)
	if err != nil {
		return err
	}

	return nil
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
		line = strings.TrimRight(line, "\r")
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
	tag     string
	version string

	mockTestPass     int
	mockTestTotal    int
	mockTestCoverage string
	mockTestSkip     int

	liveTestPass           int
	liveTestTotal          int
	liveTestCoverage       string
	liveTestCallOperations int
}

func defaultPlaceholder(v string) string {
	if v == "" || v == "0.0%" {
		return "/"
	}
	return v
}

func generateMDReport(mgmtReport map[string]mgmtInfo, path string) error {
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
		mtc := "/"
		mto := "/"
		if m.mockTestTotal != 0 {
			mt = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.mockTestPass)/float64(m.mockTestTotal)*100, m.mockTestPass, m.mockTestTotal)
			operations := m.mockTestTotal - m.mockTestSkip
			mto = fmt.Sprintf("%.2f%%(%d/%d)", float64(operations)/float64(m.mockTestTotal)*100, operations, m.mockTestTotal)

			if m.mockTestCoverage != "" {
				mtc = m.mockTestCoverage
			}
		}

		lt := "/"
		ltc := "/"
		lto := "/"
		if m.liveTestTotal != 0 {
			lt = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.liveTestPass)/float64(m.liveTestTotal)*100, m.liveTestPass, m.liveTestTotal)

			if m.liveTestCallOperations != 0 {
				lto = fmt.Sprintf("%.2f%%(%d/%d)", float64(m.liveTestCallOperations)/float64(m.mockTestTotal)*100, m.liveTestCallOperations, m.mockTestTotal)
			}

			if m.liveTestCoverage != "" {
				ltc = m.liveTestCoverage
			}
		}

		f := fmt.Sprintf("|%s | %s | %s | %s | %s | %s | %s | %s | %s |\n", module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), lt, ltc, lto, mt, mtc, mto)
		_, err = mgmtFile.Write([]byte(f))
		if err != nil {
			return err
		}
	}

	return nil
}

func generateHTMLReport(mgmtReport map[string]mgmtInfo, path string) error {
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
		mto := "/"
		mtc := "/"
		if m.mockTestTotal != 0 {
			coverage := float64(m.mockTestPass) / float64(m.mockTestTotal)
			average[0].sum += coverage
			average[0].count++
			mt = fmt.Sprintf("%.2f%%(%d/%d)", coverage*100, m.mockTestPass, m.mockTestTotal)

			mtoCoverage := float64(m.mockTestTotal-m.mockTestSkip) / float64(m.mockTestTotal)
			average[4].sum += mtoCoverage
			average[4].count++
			mto = fmt.Sprintf("%.2f%%(%d/%d)", mtoCoverage*100, m.mockTestTotal-m.mockTestSkip, m.mockTestTotal)

			if m.mockTestCoverage != "" {
				mtc = m.mockTestCoverage
				f, err := strconv.ParseFloat(strings.TrimRight(mtc, "%"), 64)
				if err != nil {
					return err
				}
				average[1].sum += f
				average[1].count++
			}
		}

		lt := "/"
		lto := "/"
		ltc := "/"
		if m.liveTestTotal != 0 {
			coverage := float64(m.liveTestPass) / float64(m.liveTestTotal)
			average[2].sum += coverage
			average[2].count++
			lt = fmt.Sprintf("%.2f%%(%d/%d)", coverage*100, m.liveTestPass, m.liveTestTotal)

			if m.mockTestTotal != 0 && m.liveTestCallOperations != 0 {
				ltoCoverage := float64(m.liveTestCallOperations) / float64(m.mockTestTotal)
				average[5].sum += ltoCoverage
				average[5].count++
				lto = fmt.Sprintf("%.2f%%(%d/%d)", ltoCoverage*100, m.liveTestCallOperations, m.mockTestTotal)
			}

			if m.liveTestCoverage != "" {
				ltc = m.liveTestCoverage
				f, err := strconv.ParseFloat(strings.TrimRight(ltc, "%"), 64)
				if err != nil {
					return err
				}
				average[3].sum += f
				average[3].count++
			}
		}

		tdBackground := ""
		if i%2 == 0 {
			tdBackground = tdBackgroundStyle
		}
		htmlData = append(htmlData, fmt.Sprintf(htmlTR, tdBackground, module, fmt.Sprintf("v%s", m.version), defaultPlaceholder(strings.TrimRight(m.tag, "\r")), lt, ltc, lto, mt, mtc, mto))
	}

	// average
	ltoc := "/"
	if average[5].sum != 0 {
		ltoc = fmt.Sprintf("%.2f%%", (average[5].sum/float64(average[5].count))*100)
	}
	htmlData = append(htmlData, fmt.Sprintf(htmlTR, "", "Average", "", "",
		fmt.Sprintf("%.2f%%", (average[2].sum/float64(average[2].count))*100),
		fmt.Sprintf("%.1f%%", average[3].sum/float64(average[3].count)),
		ltoc,
		fmt.Sprintf("%.2f%%", (average[0].sum/float64(average[0].count))*100),
		fmt.Sprintf("%.1f%%", average[1].sum/float64(average[1].count)),
		fmt.Sprintf("%.2f%%", (average[4].sum/float64(average[4].count))*100),
	))

	// parse template file
	t, err := template.ParseFiles("./mgmtreport.tpl")
	if err != nil {
		return err
	}

	w := bytes.Buffer{}
	err = t.Execute(&w, htmlData)
	if err != nil {
		return err
	}

	reportHTML, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer reportHTML.Close()

	_, err = reportHTML.Write(w.Bytes())
	if err != nil {
		return err
	}

	return err
}

func uploadHTMLReport(path, accountName, accountKey, containerName, blobName string) error {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	// The service URL for blob endpoints is usually in the form: http(s)://<account>.blob.core.windows.net/
	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return err
	}

	htmlReport, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	contentType := "text/html"
	_, err = client.UploadStream(context.TODO(),
		containerName,
		blobName,
		htmlReport,
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
