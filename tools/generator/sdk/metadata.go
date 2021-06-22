package sdk

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

func WriteMetadataFile(packagePath string, metadata autorest.GenerationMetadata) (string, error) {
	metadataFilepath := MetadataPath(packagePath)
	metadataFile, err := os.Create(metadataFilepath)
	if err != nil {
		return "", err
	}
	defer metadataFile.Close()

	// marshal metadata
	b, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("cannot marshal metadata: %+v", err)
	}

	if _, err := metadataFile.Write(b); err != nil {
		return "", err
	}
	return metadataFilepath, nil
}
