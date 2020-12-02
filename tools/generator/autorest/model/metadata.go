package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Metadata interface {
	SwaggerFiles() []string
	PackagePath() string
}

func NewMetadataFrom(reader io.Reader) (Metadata, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result localMetadata
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type localMetadata struct {
	InputFiles   []string `json:"inputFiles"`
	OutputFolder string   `json:"outputFolder"`
}

func (m localMetadata) SwaggerFiles() []string {
	return m.InputFiles
}

func (m localMetadata) PackagePath() string {
	return m.OutputFolder
}

func (m localMetadata) String() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
