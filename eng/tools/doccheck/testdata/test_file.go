//parser package
package main

import "fmt"

// Gopher is a struct
type Gopher struct {
	Gopher string `json:"gopher"`
}

func main() {
	const gopher = "GOPHER"
	gogopher := GOPHER()
	gogopher.Gopher = gopher
	fmt.Println(gogopher)
}

// GOPHER creates a gopher struct
func GOPHER() *Gopher {
	return &Gopher{Gopher : "gopher"}
}