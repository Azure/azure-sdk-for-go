package refresh

import "github.com/Azure/azure-sdk-for-go/tools/generator/autorest"

type GenerateResult struct {
	Readme     string
	Tag        string
	CommitHash string
	Package    autorest.ChangelogResult
}
