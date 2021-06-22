package utils

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/ahmetb/go-linq/v3"
)

func GetLines(file io.Reader) ([]string, error) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var result []string
	linq.From(strings.Split(string(b), "\n")).SelectT(func(line string) string {
		return strings.TrimSpace(line)
	}).ToSlice(&result)
	return result, nil
}
