package azure

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
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

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	//return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	return fmt.Sprintf("%x", uuid[10:]), nil
}
