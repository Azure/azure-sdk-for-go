package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	autorestArgsPattern = "--use=@microsoft.azure/autorest.go@~2.1.99 %s --go --multiapi --go-sdk-folder=%s --use-onever"
)

type work struct {
	filename  string
	sdkFolder string
}

func autorestCommand(file string, sdk string) *exec.Cmd {
	autorestArgs := fmt.Sprintf(autorestArgsPattern, file, sdk)
	c := exec.Command("autorest", strings.Split(autorestArgs, " ")...)
	return c
}

func worker(id int, jobs <-chan work, results chan<- error) {
	for work := range jobs {
		start := time.Now()
		vprintf("worker %d is starting on file %s\n", id, work.filename)
		c := autorestCommand(work.filename, work.sdkFolder)
		output, err := c.CombinedOutput()
		if err == nil {
			vprintf("worker %d has done with file %s (%v)\n", id, work.filename, time.Since(start))
		} else {
			printf("worker %d fails with file %s (%v), error messages:\n%v\n", id, work.filename, time.Since(start), string(output))
		}
		results <- err
	}
}

func startCmd(c *exec.Cmd) error {
	stdout, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			printf("> %s\n", scanner.Text())
		}
	}()
	stderr, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}
	scanner = bufio.NewScanner(stderr)
	go func() {
		for scanner.Scan() {
			printf("> %s\n", scanner.Text())
		}
	}()
	return c.Start()
}

func selectFilesWithName(path string, name string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == name {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}
