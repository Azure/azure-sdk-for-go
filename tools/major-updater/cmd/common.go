package cmd

import (
	"fmt"
	"os/exec"
	"bufio"
)

func printf(format string, a ...interface{}) {
	if !quietFlag {
		fmt.Printf(format, a...)
	}
}

func println(a ...interface{}) {
	if !quietFlag {
		fmt.Println(a...)
	}
}

func dprintf(format string, a ...interface{}) {
	if debugFlag {
		printf(format, a...)
	}
}

func dprintln(a ...interface{}) {
	if debugFlag {
		println(a...)
	}
}

func vprintf(format string, a ...interface{}) {
	if verboseFlag {
		printf(format, a...)
	}
}

func vprintln(a ...interface{}) {
	if verboseFlag {
		println(a...)
	}
}

func contains(strings []string, str string) bool {
	for _, s := range strings {
		if s == str {
			return true
		}
	}
	return false
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
