package cmd

import (
	"fmt"
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
