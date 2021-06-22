package sdk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

func WriteChangelogFile(result autorest.ChangelogResult) (string, error) {
	fileContent := fmt.Sprintf(`# Change History

%s`, result.Changelog.ToMarkdown())
	path := filepath.Join(result.PackageFullPath, autorest.ChangelogFilename)
	changelogFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer changelogFile.Close()
	if _, err := changelogFile.WriteString(fileContent); err != nil {
		return "", err
	}
	return path, nil
}
