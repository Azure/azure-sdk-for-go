package refresh

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest_ext"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
)

type generateContext struct {
	sdkRoot        string
	specRoot       string
	specCommitHash string
	options        model.Options

	repoContent map[string]exports.Content
}

func (ctx generateContext) SDKRoot() string {
	return ctx.sdkRoot
}

func (ctx generateContext) SpecRoot() string {
	return ctx.specRoot
}

func (ctx generateContext) RepoContent() map[string]exports.Content {
	return ctx.repoContent
}

func (ctx *generateContext) generate(info GenerationInfo) (*GenerateResult, error) {
	log.Printf("Generating readme '%s' tag '%s' from path '%s'", info.RelativeReadme(), info.Tag, info.PackageFullPath)
	metadataOutputRoot := filepath.Join(os.TempDir(), fmt.Sprintf("refresh-metadata-%v", time.Now().Unix()))
	defer os.RemoveAll(metadataOutputRoot)

	// Generate code
	input := autorest_ext.GenerateInput{
		Readme:     info.RelativeReadme(),
		Tag:        info.Tag,
		SDKRoot:    ctx.SDKRoot(),
		CommitHash: ctx.specCommitHash,
		Options:    ctx.options,
	}
	r, err := autorest_ext.GeneratePackage(ctx, input, autorest_ext.GenerateOptions{
		MetadataOutputRoot: metadataOutputRoot,
		Stderr:             os.Stderr,
		Stdout:             os.Stderr,
	})
	if err != nil {
		log.Printf("Error generating package for readme '%s' tag '%s': %+v", info.RelativeReadme(), info.Tag, err)
		return nil, fmt.Errorf("cannot generate readme '%s', tag '%s': %+v", info.RelativeReadme(), info.Tag, err)
	}

	return &GenerateResult{
		Readme:     info.Readme,
		Tag:        info.Tag,
		CommitHash: info.CommitHash,
		Package:    r.Package,
	}, nil
}
