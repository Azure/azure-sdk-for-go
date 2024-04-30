//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const proxyManualStartEnv = "PROXY_MANUAL_START"

type TestProxyInstance struct {
	Cmd     *exec.Cmd
	Options *RecordingOptions
}

func getTestProxyDownloadFile() (string, error) {
	switch {
	case runtime.GOOS == "windows":
		// No ARM binaries for Windows, so return x64
		return "test-proxy-standalone-win-x64.zip", nil
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		return "test-proxy-standalone-linux-x64.tar.gz", nil
	case runtime.GOOS == "linux" && runtime.GOARCH == "arm64":
		return "test-proxy-standalone-linux-arm64.tar.gz", nil
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		return "test-proxy-standalone-osx-x64.zip", nil
	case runtime.GOOS == "darwin" && runtime.GOARCH == "arm64":
		return "test-proxy-standalone-osx-arm64.zip", nil
	default:
		return "", fmt.Errorf("unsupported OS/Arch combination: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}

// Modified from https://stackoverflow.com/a/24792688
func extractTestProxyZip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		log.Println("Extracting", path)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractTestProxyArchive(archivePath string, outputDir string) error {
	log.Printf("Extracting %s\n", archivePath)
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(outputDir, header.Name)

		log.Println("Extracting", targetPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			log.Printf("Unable to extract type %c in file %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}

func installTestProxy(archivePath string, outputDir string, proxyPath string) error {
	var err error
	if strings.HasSuffix(archivePath, ".zip") {
		err = extractTestProxyZip(archivePath, outputDir)
	} else {
		err = extractTestProxyArchive(archivePath, outputDir)
	}
	if err != nil {
		return err
	}

	err = os.Chmod(proxyPath, 0755)
	if err != nil {
		return err
	}
	err = os.Remove(archivePath)
	if err != nil {
		return err
	}

	return nil
}

func restoreRecordings(proxyPath string, pathToRecordings string) error {
	if pathToRecordings == "" {
		return nil
	}
	absAssetLocation, _, err := getAssetsConfigLocation(pathToRecordings)
	if err != nil {
		return err
	}

	log.Printf("Running test proxy command: %s restore -a %s\n", proxyPath, absAssetLocation)
	cmd := exec.Command(proxyPath, "restore", "-a", absAssetLocation)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func ensureTestProxyInstalled(proxyVersion string, proxyPath string, proxyDir string, pathToRecordings string) error {
	lockFile := filepath.Join(os.TempDir(), "test-proxy-install.lock")
	log.Printf("Waiting to acquire test proxy install lock %s\n", lockFile)
	maxTries := 600 // Wait 1 minute
	var i int
	for i = 0; i < maxTries; i++ {
		lock, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// NOTE: the lockfile will not be removed on ctrl-c during download.
		// Go test seems to send an os.Interrupt signal on test setup completion, so if we
		// call os.Exit(1) on ctrl-c the tests will never run. If we don't call os.Exit(1),
		// the tests cannot be canceled.
		// Therefore, if ctrl-c is pressed during download, the user will have to manually
		// remove the lockfile in order to get the tests running again.
		defer func() {
			lock.Close()
			os.Remove(lockFile)
		}()

		break
	}

	if i >= maxTries {
		return fmt.Errorf("timed out waiting to acquire test proxy install lock. Ensure %s does not exist", lockFile)
	}

	cmd := exec.Command(proxyPath, "--version")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Test proxy not detected at %s, downloading...\n", proxyPath)
	} else {
		// TODO: fix proxy CLI tool versioning output to match the actual version we download
		installedVersion := "1.0.0-dev." + strings.TrimSpace(string(out))
		if installedVersion == proxyVersion {
			log.Printf("Test proxy version %s already installed\n", proxyVersion)
			return restoreRecordings(proxyPath, pathToRecordings)
		} else {
			log.Printf("Test proxy version %s does not match required version %s\n",
				installedVersion, proxyVersion)
		}
	}

	proxyFile, err := getTestProxyDownloadFile()
	if err != nil {
		return err
	}

	proxyDownloadPath := filepath.Join(proxyDir, proxyFile)
	archive, err := os.Create(proxyDownloadPath)
	if err != nil {
		return err
	}

	log.Printf("Downloading test proxy version %s to %s for %s/%s\n",
		proxyVersion, proxyPath, runtime.GOOS, runtime.GOARCH)
	proxyUrl := fmt.Sprintf("https://github.com/Azure/azure-sdk-tools/releases/download/Azure.Sdk.Tools.TestProxy_%s/%s",
		proxyVersion, proxyFile)
	resp, err := http.Get(proxyUrl)
	if err != nil {
		return err
	}

	_, err = io.Copy(archive, resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	err = archive.Close()
	if err != nil {
		return err
	}

	err = installTestProxy(proxyDownloadPath, proxyDir, proxyPath)
	if err != nil {
		return err
	}

	return restoreRecordings(proxyPath, pathToRecordings)
}

func getProxyLog() (*os.File, error) {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz"
	suffix := make([]byte, 6)
	for i := range suffix {
		suffix[i] = letters[rand.Intn(len(letters))]
	}
	proxyLogName := fmt.Sprintf("test-proxy.log.%s", suffix)
	proxyLog, err := os.Create(filepath.Join(os.TempDir(), proxyLogName))
	if err != nil {
		return nil, err
	}
	return proxyLog, nil
}

func getProxyVersion(gitRoot string) (string, error) {
	proxyVersionConfig := filepath.Join(gitRoot, "eng/common/testproxy/target_version.txt")
	overrideProxyVersionConfig := filepath.Join(gitRoot, "eng/target_proxy_version.txt")

	if _, err := os.Stat(overrideProxyVersionConfig); err == nil {
		version, err := ioutil.ReadFile(overrideProxyVersionConfig)
		if err == nil {
			proxyVersion := strings.TrimSpace(string(version))
			return proxyVersion, nil
		}
	}

	version, err := ioutil.ReadFile(proxyVersionConfig)
	if err != nil {
		return "", err
	}
	proxyVersion := strings.TrimSpace(string(version))

	return proxyVersion, nil
}

func setTestProxyEnv(gitRoot string) {
	devCertPath := filepath.Join(gitRoot, "eng/common/testproxy/dotnet-devcert.pfx")
	os.Setenv("ASPNETCORE_Kestrel__Certificates__Default__Path", devCertPath)
	os.Setenv("ASPNETCORE_Kestrel__Certificates__Default__Password", "password")
}

func waitForProxyStart(cmd *exec.Cmd, options *RecordingOptions) (*TestProxyInstance, error) {
	maxTries := 50
	// Extend sleep time in devops pipeline, proxy takes longer to start up
	if os.Getenv("SYSTEM_TEAMPROJECTID") != "" {
		maxTries = 200
	}
	log.Printf("Started test proxy instance (PID %d) on %s\n", cmd.Process.Pid, options.baseURL())
	client, _ := GetHTTPClient(nil)
	client.Timeout = 1 * time.Second

	log.Printf("Waiting up to %d seconds for test-proxy server to respond...\n", (maxTries / 10))
	var i int
	for i = 0; i < maxTries; i++ {
		uri := fmt.Sprintf("https://localhost:%d/Admin/IsAlive", options.ProxyPort)
		req, _ := http.NewRequest("GET", uri, nil)
		req.Close = true

		resp, err := client.Do(req)
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return &TestProxyInstance{Cmd: cmd, Options: options}, nil
	}

	return nil, fmt.Errorf("test proxy server did not become available in the allotted time")
}

func StartTestProxy(pathToRecordings string, options *RecordingOptions) (*TestProxyInstance, error) {
	manualStart := strings.ToLower(os.Getenv(proxyManualStartEnv))
	if manualStart == "true" {
		log.Printf("%s env variable is set to true, not starting test proxy...\n", proxyManualStartEnv)
		return nil, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	gitRoot, err := getGitRoot(cwd)
	if err != nil {
		return nil, err
	}

	proxyVersion, err := getProxyVersion(gitRoot)
	if err != nil {
		return nil, err
	}
	proxyDir := filepath.Join(gitRoot, ".proxy")
	if err := os.MkdirAll(proxyDir, 0755); err != nil {
		return nil, err
	}
	proxyPath := filepath.Join(proxyDir, "Azure.Sdk.Tools.TestProxy")
	if runtime.GOOS == "windows" {
		proxyPath += ".exe"
	}
	err = ensureTestProxyInstalled(proxyVersion, proxyPath, proxyDir, pathToRecordings)
	if err != nil {
		return nil, err
	}

	proxyLog, err := getProxyLog()
	if err != nil {
		return nil, err
	}
	defer proxyLog.Close()

	setTestProxyEnv(gitRoot)

	if options == nil {
		options = defaultOptions()
	}
	insecure := ""
	if options.insecure {
		insecure = "--insecure"
	}
	args := []string{"start", "--storage-location", gitRoot, insecure, "--", "--urls=" + options.baseURL()}
	log.Printf("Running test proxy command: %s %s", proxyPath, strings.Join(args, " "))
	log.Printf("Test proxy log location: %s\n", proxyLog.Name())
	cmd := exec.Command(proxyPath, args...)

	cmd.Stdout = proxyLog
	cmd.Stderr = proxyLog

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	return waitForProxyStart(cmd, options)
}

// NOTE: The process will be killed if the user hits ctrl-c mid-way through tests, as go will
// kill child processes when the main process does not exit cleanly. No os.Interrupt handlers
// need to be added after starting the proxy server in tests.
func StopTestProxy(proxyInstance *TestProxyInstance) error {
	if proxyInstance == nil || proxyInstance.Cmd == nil || proxyInstance.Cmd.Process == nil {
		return nil
	}
	log.Printf("Stopping test proxy instance (PID %d) on %s\n", proxyInstance.Cmd.Process.Pid, proxyInstance.Options.baseURL())
	err := proxyInstance.Cmd.Process.Kill()
	if err != nil {
		return err
	}
	return nil
}
