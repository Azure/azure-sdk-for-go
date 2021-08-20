// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/template"
	"github.com/Azure/azure-sdk-for-go/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Masterminds/semver"
)

type GenerateContext struct {
	SdkPath    string
	SpecPath   string
	CommitHash string
}

type GenerateResult struct {
	Version        string
	RpName         string
	PackageName    string
	PackageAbsPath string
	Changelog      model.Changelog
	ChangelogMd    string
}

func (ctx GenerateContext) SDKRoot() string {
	return ctx.SdkPath
}

func (ctx GenerateContext) SpecRoot() string {
	return ctx.SpecPath
}

func (ctx GenerateContext) GenerateForAutomation(readme string, repo string) ([]GenerateResult, []error) {
	absReadme := filepath.Join(ctx.SpecPath, readme)
	absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")

	var result []GenerateResult
	var errors []error

	log.Printf("Get all namespaces from readme file")
	rpMap, err := ReadV2ModuleNameToGetNamespace(absReadmeGo)
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot get rp and namespaces from readme '%s': %+v", readme, err),
		}
	}

	for rpName, namespaceNames := range rpMap {
		for _, namespaceName := range namespaceNames {
			log.Printf("Process rp: %s, namespace: %s", rpName, namespaceName)
			singleResult, err := ctx.GenerateForSingleRpNamespace(rpName, namespaceName, "", "", repo)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			result = append(result, *singleResult)
		}
	}
	return result, errors
}

func (ctx GenerateContext) GenerateForSingleRpNamespace(rpName, namespaceName, specficPackageTitle, specficVersion, specficRepoURL string) (*GenerateResult, error) {
	packagePath := filepath.Join(ctx.SdkPath, "sdk", rpName, namespaceName)
	changelogPath := filepath.Join(packagePath, common.ChangelogFilename)
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		log.Printf("Package '%s' changelog not exist, do onboard process", packagePath)

		if specficPackageTitle == "" {
			specficPackageTitle = strings.Title(rpName)
		}

		log.Printf("Use template to generate new rp folder and basic package files...")
		if err = template.GeneratePackageByTemplate(rpName, namespaceName, template.Flags{
			SDKRoot:      ctx.SdkPath,
			TemplatePath: "tools/generator/template/rpName/packageName",
			PackageTitle: specficPackageTitle,
			Commit:       ctx.CommitHash,
		}); err != nil {
			return nil, err
		}

		if specficRepoURL != "" {
			log.Printf("Change the repo url in `autorest.md`...")
			autorestMdPath := filepath.Join(packagePath, "autorest.md")
			if err = ReplaceRepoURL(autorestMdPath, specficRepoURL); err != nil {
				return nil, err
			}
		}

		log.Printf("Run `go generate` to regenerate the code...")
		if err = ExecuteGoGenerate(packagePath); err != nil {
			return nil, err
		}

		log.Printf("Generate changelog for package...")
		newExports, err := exports.Get(packagePath)
		if err != nil {
			return nil, err
		}
		changelog, err := autorest.GetChangelogForPackage(nil, &newExports)
		if err != nil {
			return nil, err
		}

		log.Printf("Replace {{NewClientMethod}} placeholder in the README.md ")
		if err = ReplaceNewClientMethodPlaceholder(packagePath, newExports); err != nil {
			return nil, err
		}

		return &GenerateResult{
			Version:        "0.1.0",
			RpName:         rpName,
			PackageName:    namespaceName,
			PackageAbsPath: packagePath,
			Changelog:      *changelog,
			ChangelogMd:    changelog.ToCompactMarkdown(),
		}, nil
	} else {
		log.Printf("Package '%s' existed, do update process", packagePath)

		log.Printf("Get ori exports for changelog generation...")
		oriExports, err := exports.Get(packagePath)
		if err != nil {
			return nil, err
		}

		log.Printf("Remove all the files that start with `zz_generated_`...")
		if err = CleanSDKGeneratedFiles(packagePath); err != nil {
			return nil, err
		}

		log.Printf("Change the commit hash in `autorest.md` to a new commit that corresponds to the new release...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err = ReplaceCommitID(autorestMdPath, ctx.CommitHash); err != nil {
			return nil, err
		}

		if specficRepoURL != "" {
			log.Printf("Change the repo url in `autorest.md`...")
			if err = ReplaceRepoURL(autorestMdPath, specficRepoURL); err != nil {
				return nil, err
			}
		}

		log.Printf("Run `go generate` to regenerate the code...")
		if err = ExecuteGoGenerate(packagePath); err != nil {
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

		log.Printf("Calculate new version...")
		var version *semver.Version
		if specficVersion == "" {
			version, err = CalculateNewVersion(changelog, packagePath)
			if err != nil {
				return nil, err
			}
		} else {
			log.Printf("Use specfic version: %s", specficVersion)
			version, err = semver.NewVersion(specficVersion)
			if err != nil {
				return nil, err
			}
		}

		log.Printf("Add changelog to file...")
		changelogMd, err := AddChangelogToFile(changelog, version, packagePath)
		if err != nil {
			return nil, err
		}

		log.Printf("Remove all the files that start with `zz_generated_`...")
		if err = CleanSDKGeneratedFiles(packagePath); err != nil {
			return nil, err
		}

		log.Printf("Replace version in autorest.md...")
		if err = ReplaceVersion(packagePath, version.String()); err != nil {
			return nil, err
		}

		log.Printf("Run `go generate` to regenerate the code for new version...")
		if err = ExecuteGoGenerate(packagePath); err != nil {
			return nil, err
		}

		return &GenerateResult{
			Version:        version.String(),
			RpName:         rpName,
			PackageName:    namespaceName,
			PackageAbsPath: packagePath,
			Changelog:      *changelog,
			ChangelogMd:    changelogMd,
		}, nil
	}
}
