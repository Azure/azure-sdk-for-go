package model

import "encoding/json"

type Metadata struct {
	InputFiles   []string `json:"inputFiles"`
	OutputFolder string   `json:"outputFolder"`
}

func (m Metadata) String() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
