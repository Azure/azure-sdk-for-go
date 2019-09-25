package cmd

import "fmt"

func vprintf(format string, a ...interface{}) {
	if verboseFlag {
		fmt.Printf(format, a...)
	}
}

func vprintln(message string) {
	if verboseFlag {
		fmt.Println(message)
	}
}

func printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func println(message string) {
	fmt.Println(message)
}
