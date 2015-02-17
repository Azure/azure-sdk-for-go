package management

import (
	"bytes"
	"fmt"
	"github.com/MSOpenTech/azure-sdk-for-go/core/http"
	"io"
	"os/exec"
	"strings"
)

func executeCommand(command string, input []byte) ([]byte, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf(errParamNotSpecified, "command")
	}

	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	if input != nil {
		cmd.Stdin = bytes.NewReader(input)
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func getResponseBody(response *http.Response) []byte {
	responseBody := make([]byte, response.ContentLength)
	io.ReadFull(response.Body, responseBody)
	return responseBody
}
