// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/template"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Masterminds/semver"
)

type GenerateContext struct {
	SDKPath        string
	SDKRepo        *repo.SDKRepository
	SpecPath       string
	SpecCommitHash string
	SpecRepoURL    string
}

type GenerateResult struct {
	Version        string
	RPName         string
	PackageName    string
	PackageAbsPath string
	Changelog      model.Changelog
	ChangelogMD    string
}

type GenerateParam struct {
	RPName              string
	NamespaceName       string
	NamespaceConfig     string
	SpecficVersion      string
	SpecficPackageTitle string
	SpecRPName          string
	ReleaseDate         string
	SkipGenerateExample bool
	GoVersion           string
}

func (ctx GenerateContext) GenerateForAutomation(readme, repo, goVersion string) ([]GenerateResult, []error) {
	absReadme := filepath.Join(ctx.SpecPath, readme)
	absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")
	specRPName := strings.Split(readme, "/")[0]

	var result []GenerateResult
	var errors []error

	log.Printf("Get all namespaces from readme file")
	rpMap, err := ReadV2ModuleNameToGetNamespace(absReadmeGo)
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot get rp and namespaces from readme '%s': %+v", readme, err),
		}
	}

	for rpName, packageInfos := range rpMap {
		for _, packageInfo := range packageInfos {
			log.Printf("Process rp: %s, namespace: %s", rpName, packageInfo.Name)
			singleResult, err := ctx.GenerateForSingleRPNamespace(&GenerateParam{
				RPName:              rpName,
				NamespaceName:       packageInfo.Name,
				SpecRPName:          specRPName,
				SkipGenerateExample: true,
				NamespaceConfig:     packageInfo.Config,
				GoVersion:           goVersion,
			})
			if err != nil {
				errors = append(errors, err)
				continue
			}
			result = append(result, *singleResult)
		}
	}
	return result, errors
}

func (ctx GenerateContext) GenerateForSingleRPNamespace(generateParam *GenerateParam) (*GenerateResult, error) {
	packagePath := filepath.Join(ctx.SDKPath, "sdk", "resourcemanager", generateParam.RPName, generateParam.NamespaceName)
	changelogPath := filepath.Join(packagePath, common.ChangelogFilename)

	onBoard := false
	var oriExports exports.Content
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		onBoard = true
		log.Printf("Package '%s' changelog not exist, do onboard process", packagePath)

		if generateParam.SpecficPackageTitle == "" {
			generateParam.SpecficPackageTitle = strings.Title(generateParam.RPName)
		}

		log.Printf("Use template to generate new rp folder and basic package files...")
		if err = template.GeneratePackageByTemplate(generateParam.RPName, generateParam.NamespaceName, template.Flags{
			SDKRoot:       ctx.SDKPath,
			TemplatePath:  "eng/tools/generator/template/rpName/packageName",
			PackageTitle:  generateParam.SpecficPackageTitle,
			Commit:        ctx.SpecCommitHash,
			PackageConfig: generateParam.NamespaceConfig,
			GoVersion:     generateParam.GoVersion,
		}); err != nil {
			return nil, err
		}
	} else {
		log.Printf("Package '%s' existed, do update process", packagePath)

		log.Printf("Get ori exports for changelog generation...")
		oriExports, err = exports.Get(packagePath)
		if err != nil {
			log.Printf("Get ori exports error, set to empty: %+v", err)
			oriExports = exports.Content{}
		}

		log.Printf("Remove all the files that start with `zz_generated_`...")
		if err = CleanSDKGeneratedFiles(packagePath); err != nil {
			return nil, err
		}
	}

	// same step for onboard and update
	if ctx.SpecCommitHash == "" {
		log.Printf("Change swagger config in `autorest.md` according to local path...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := ChangeConfigWithLocalPath(autorestMdPath, ctx.SpecPath, generateParam.SpecRPName); err != nil {
			return nil, err
		}
	} else {
		log.Printf("Change swagger config in `autorest.md` according to repo URL and commit ID...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := ChangeConfigWithCommitID(autorestMdPath, ctx.SpecRepoURL, ctx.SpecCommitHash, generateParam.SpecRPName); err != nil {
			return nil, err
		}
	}

	log.Printf("Run `go generate` to regenerate the code...")
	if err := ExecuteGoGenerate(packagePath); err != nil {
		return nil, err
	}

	log.Printf("Generate changelog for package...")
	newExports, err := exports.Get(packagePath)
	if err != nil {
		return nil, err
	}
	changelog, err := autorest.GetChangelogForPackage(&oriExports, &newExports)
	if err != nil {
		return nil, err
	}

	if !generateParam.SkipGenerateExample {
		log.Printf("Generate examples...")
		if err := ExecuteExampleGenerate(packagePath, filepath.Join("resourcemanager", generateParam.RPName, generateParam.NamespaceName)); err != nil {
			return nil, err
		}
	}

	if onBoard {
		log.Printf("Replace {{NewClientName}} placeholder in the README.md ")
		if err = ReplaceNewClientNamePlaceholder(packagePath, newExports); err != nil {
			return nil, err
		}

		return &GenerateResult{
			Version:        "0.1.0",
			RPName:         generateParam.RPName,
			PackageName:    generateParam.NamespaceName,
			PackageAbsPath: packagePath,
			Changelog:      *changelog,
			ChangelogMD:    changelog.ToCompactMarkdown() + "\n" + changelog.GetChangeSummary(),
		}, nil
	} else {
		log.Printf("Calculate new version...")
		var version *semver.Version
		if generateParam.SpecficVersion == "" {
			version, err = CalculateNewVersion(changelog, packagePath)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("Use specfic version: %s", generateParam.SpecficVersion)
			version, err = semver.NewVersion(generateParam.SpecficVersion)
			if err != nil {
				return nil, err
			}
		}

		log.Printf("Add changelog to file...")
		changelogMd, err := AddChangelogToFile(changelog, version, packagePath, generateParam.ReleaseDate)
		if err != nil {
			return nil, err
		}

		log.Printf("Replace version in autorest.md and constants...")
		if err = ReplaceVersion(packagePath, version.String()); err != nil {
			return nil, err
		}

		return &GenerateResult{
			Version:        version.String(),
			RPName:         generateParam.RPName,
			PackageName:    generateParam.NamespaceName,
			PackageAbsPath: packagePath,
			Changelog:      *changelog,
			ChangelogMD:    changelogMd + "\n" + changelog.GetChangeSummary(),
		}, nil
	}
}
