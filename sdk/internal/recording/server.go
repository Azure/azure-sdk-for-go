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
	"strings"
	"runtime"
	"time"
)

func getTestProxyDownloadFile() (string, error) {
	if runtime.GOOS == "windows" {
		return "test-proxy-standalone-win-x64.zip", nil
	}

	switch {
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

func extractTestProxyZip(archivePath string, outputDir string) error {
    // Open the zip file
    r, err := zip.OpenReader(archivePath)
    if err != nil {
        panic(err)
    }
    defer r.Close()

    for _, f := range r.File {
        targetPath := filepath.Join(outputDir, f.Name)

        log.Println("Extracting", targetPath)

        if f.FileInfo().IsDir() {
            os.MkdirAll(targetPath, f.Mode())
            continue
        }

        file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return err
        }
        defer file.Close()

        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        if _, err = io.Copy(file, rc); err != nil {
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

func extractTestProxy(archivePath string, outputDir string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractTestProxyZip(archivePath, outputDir)
	} else {
		return extractTestProxyArchive(archivePath, outputDir)
	}
}

func ensureTestProxyInstalled(proxyVersion string, proxyPath string, proxyDir string) error {
	lockFile := filepath.Join(os.TempDir(), "test-proxy-install.lock")
	maxTries := 600  // Wait 1 minute
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
			os.Remove(lockFile)
			lock.Close()
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
			return nil
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
    defer archive.Close()

	log.Printf("Downloading test proxy version %s to %s for %s/%s\n",
				proxyVersion, proxyPath, runtime.GOOS, runtime.GOARCH)
	proxyUrl := fmt.Sprintf("https://github.com/Azure/azure-sdk-tools/releases/download/Azure.Sdk.Tools.TestProxy_%s/%s",
							 proxyVersion, proxyFile)
    resp, err := http.Get(proxyUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    _, err = io.Copy(archive, resp.Body)
    if err != nil {
        return err
    }

	err = extractTestProxy(proxyDownloadPath, proxyDir)
	if err != nil {
		return err
	}
	err = os.Chmod(proxyPath, 0755)
	if err != nil {
		return err
	}
	err = os.Remove(proxyDownloadPath)
	if err != nil {
		return err
	}

	return nil
}

func getProxyLog() (*os.File, error) {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz"
	suffix := make([]byte, 6)
	for i := range suffix {
		suffix[i] = letters[rand.Intn(len(letters))]
	}
	proxyLogName := fmt.Sprintf("testproxy.log.%s", suffix)
	proxyLog, err := os.Create(filepath.Join(os.TempDir(), proxyLogName))
	if err != nil {
		return nil, err
	}
	return proxyLog, nil
}

func StartTestProxyInstance(options *RecordingOptions) (*exec.Cmd, error) {
	manualStart := strings.ToLower(os.Getenv("PROXY_MANUAL_START"))
	if manualStart == "true" {
		log.Println("PROXY_MANUAL_START env variable is set to true, not starting test proxy...")
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
	proxyVersionConfig := filepath.Join(gitRoot, "eng/common/testproxy/target_version.txt")
	version, err := ioutil.ReadFile(proxyVersionConfig)
	if err != nil {
		return nil, err
	}
	proxyVersion := strings.TrimSpace(string(version))

	proxyDir := filepath.Join(gitRoot, ".proxy")
	if err := os.MkdirAll(proxyDir, 0755); err != nil {
		return nil, err
	}

    proxyPath := filepath.Join(proxyDir, "Azure.Sdk.Tools.TestProxy")
	err = ensureTestProxyInstalled(proxyVersion, proxyPath, proxyDir)
	if err != nil {
		return nil, err
	}

	proxyLog, err := getProxyLog()
	if err != nil {
		return nil, err
	}
	defer proxyLog.Close()

	if options == nil {
		options = defaultOptions()
	}
	log.Printf("Running test proxy command: %s start --storage-location %s -- --urls=%s\n",
				proxyPath, gitRoot, options.baseURL())
	log.Printf("Test proxy log location: %s\n", proxyLog.Name())
	cmd := exec.Command(
		proxyPath, "start", "--storage-location", gitRoot, "--", "--urls=" + options.baseURL())

	cmd.Stdout = proxyLog
	cmd.Stderr = proxyLog

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Give background test proxy instance time to start up
	time.Sleep(2 * time.Second)
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		return nil, fmt.Errorf("test proxy instance failed to start in the allotted time")
	}
	log.Printf("Started test proxy instance (PID %d) on %s\n", cmd.Process.Pid, options.baseURL())

	return cmd, nil
}

func StopTestProxyInstance(proxyCmd *exec.Cmd, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	if proxyCmd == nil {
		return nil
	}
	log.Printf("Stopping test proxy instance (PID %d) on %s\n", proxyCmd.Process.Pid, options.baseURL())
	err := proxyCmd.Process.Kill()
	if err != nil {
		return err
	}
	return nil
}