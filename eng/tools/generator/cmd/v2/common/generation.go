// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/template"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Masterminds/semver"
)

type GenerateContext struct {
	SDKPath           string
	SDKRepo           *repo.SDKRepository
	SpecPath          string
	SpecCommitHash    string
	SpecReadmeFile    string
	SpecReadmeGoFile  string
	SpecRepoURL       string
	UpdateSpecVersion bool
}

type GenerateResult struct {
	Version           string
	RPName            string
	PackageName       string
	PackageAbsPath    string
	Changelog         Changelog
	ChangelogMD       string
	PullRequestLabels string
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
	RemoveTagSet        bool
	ForceStableVersion  bool
}

func (ctx *GenerateContext) GenerateForAutomation(readme, repo, goVersion string) ([]GenerateResult, []error) {
	absReadme, err := filepath.Abs(filepath.Join(ctx.SpecPath, readme))
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot get absolute path for spec path '%s': %+v", ctx.SpecPath, err),
		}
	}
	absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")
	ctx.SpecReadmeFile = absReadme
	ctx.SpecReadmeGoFile = absReadmeGo
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
				RemoveTagSet:        true,
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

func (ctx *GenerateContext) GenerateForSingleRPNamespace(generateParam *GenerateParam) (*GenerateResult, error) {
	packagePath := filepath.Join(ctx.SDKPath, "sdk", "resourcemanager", generateParam.RPName, generateParam.NamespaceName)
	changelogPath := filepath.Join(packagePath, "CHANGELOG.md")

	onBoard := false

	version, err := semver.NewVersion("0.1.0")
	if err != nil {
		return nil, err
	}
	if generateParam.SpecficVersion != "" {
		log.Printf("Use specfic version: %s", generateParam.SpecficVersion)
		version, err = semver.NewVersion(generateParam.SpecficVersion)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		onBoard = true
		log.Printf("Package '%s' changelog not exist, do onboard process", packagePath)

		if generateParam.SpecficPackageTitle == "" {
			generateParam.SpecficPackageTitle = strings.Title(generateParam.RPName)
		}

		log.Printf("Use template to generate new rp folder and basic package files...")
		if err = template.GeneratePackageByTemplate(generateParam.RPName, generateParam.NamespaceName, template.Flags{
			SDKRoot:        ctx.SDKPath,
			TemplatePath:   "eng/tools/generator/template/rpName/packageName",
			PackageTitle:   generateParam.SpecficPackageTitle,
			Commit:         ctx.SpecCommitHash,
			PackageConfig:  generateParam.NamespaceConfig,
			GoVersion:      generateParam.GoVersion,
			PackageVersion: version.String(),
			ReleaseDate:    generateParam.ReleaseDate,
		}); err != nil {
			return nil, err
		}
	} else {
		log.Printf("Package '%s' existed, do update process", packagePath)

		log.Printf("Remove all the generated files ...")
		if err = CleanSDKGeneratedFiles(packagePath); err != nil {
			return nil, err
		}
	}

	// same step for onboard and update
	if ctx.SpecCommitHash == "" {
		log.Printf("Change swagger config in `autorest.md` according to local path...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := ChangeConfigWithLocalPath(autorestMdPath, ctx.SpecReadmeFile, ctx.SpecReadmeGoFile); err != nil {
			return nil, err
		}
	} else {
		log.Printf("Change swagger config in `autorest.md` according to repo URL and commit ID...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if ctx.UpdateSpecVersion {
			if err := ChangeConfigWithCommitID(autorestMdPath, ctx.SpecRepoURL, ctx.SpecCommitHash, generateParam.SpecRPName); err != nil {
				return nil, err
			}
		}
	}

	// remove tag set
	if generateParam.RemoveTagSet {
		log.Printf("Remove tag set for swagger config in `autorest.md`...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := RemoveTagSet(autorestMdPath); err != nil {
			return nil, err
		}
	}

	// add tag set
	if !generateParam.RemoveTagSet && generateParam.NamespaceConfig != "" && !onBoard {
		log.Printf("Add tag in `autorest.md`...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := AddTagSet(autorestMdPath, generateParam.NamespaceConfig); err != nil {
			return nil, err
		}
	}

	log.Printf("Run `go generate` to regenerate the code...")
	if err := ExecuteGoGenerate(packagePath); err != nil {
		return nil, err
	}

	previousVersion := ""
	isCurrentPreview := false
	var oriExports *exports.Content
	isCurrentPreview, err = ContainsPreviewAPIVersion(packagePath)
	if err != nil {
		return nil, err
	}

	if isCurrentPreview && generateParam.ForceStableVersion {
		tag, err := GetTag(filepath.Join(packagePath, "autorest.md"))
		if err != nil {
			return nil, err
		}
		if tag != "" {
			if !strings.Contains(tag, "preview") {
				isCurrentPreview = false
			}
		}
	}

	if !onBoard {
		log.Printf("Get ori exports for changelog generation...")

		tags, err := GetAllVersionTags(generateParam.RPName, generateParam.NamespaceName)
		if err != nil {
			return nil, err
		}

		if len(tags) == 0 {
			return nil, fmt.Errorf("github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/%s/%s hasn't been released, it's supposed to OnBoard", generateParam.RPName, generateParam.NamespaceName)
		}

		previousVersionTag := GetPreviousVersionTag(isCurrentPreview, tags)

		oriExports, err = GetExportsFromTag(*ctx.SDKRepo, packagePath, previousVersionTag)
		if err != nil {
			return nil, err
		}

		tagSplit := strings.Split(previousVersionTag, "/")
		previousVersion = strings.TrimLeft(tagSplit[len(tagSplit)-1], "v")
	}

	log.Printf("Generate changelog for package...")
	newExports, err := exports.Get(packagePath)
	if err != nil {
		return nil, err
	}
	changelog, err := GetChangelogForPackage(oriExports, &newExports)
	if err != nil {
		return nil, err
	}

	log.Printf("filter changelog...")
	FilterChangelog(changelog, NonExportedFilter, MarshalUnmarshalFilter, EnumFilter, FuncFilter, LROFilter, PageableFilter, InterfaceToAnyFilter)

	var prl PullRequestLabel
	if onBoard {
		log.Printf("Replace {{NewClientName}} placeholder in the README.md ")
		if err = ReplaceNewClientNamePlaceholder(packagePath, newExports); err != nil {
			return nil, err
		}

		if !generateParam.SkipGenerateExample {
			log.Printf("Generate examples...")
			flag, err := GetAlwaysSetBodyParamRequiredFlag(filepath.Join(packagePath, "build.go"))
			if err != nil {
				return nil, err
			}
			if err := ExecuteExampleGenerate(packagePath, filepath.Join("resourcemanager", generateParam.RPName, generateParam.NamespaceName), flag); err != nil {
				return nil, err
			}
		}

		prl = FirstBetaLabel
		if !isCurrentPreview {
			version, err = semver.NewVersion("1.0.0")
			if err != nil {
				return nil, err
			}

			log.Printf("Replace version in CHANGELOG.md...")
			if err = UpdateOnboardChangelogVersion(packagePath, version.String()); err != nil {
				return nil, err
			}

			log.Printf("Replace version in autorest.md and constants...")
			if err = ReplaceVersion(packagePath, version.String()); err != nil {
				return nil, err
			}
			prl = FirstGALabel
		}

		return &GenerateResult{
			Version:           version.String(),
			RPName:            generateParam.RPName,
			PackageName:       generateParam.NamespaceName,
			PackageAbsPath:    packagePath,
			Changelog:         *changelog,
			ChangelogMD:       changelog.ToCompactMarkdown() + "\n" + changelog.GetChangeSummary(),
			PullRequestLabels: string(prl),
		}, nil
	} else {
		log.Printf("Calculate new version...")
		if generateParam.SpecficVersion == "" {
			version, prl, err = CalculateNewVersion(changelog, previousVersion, isCurrentPreview)
			if err != nil {
				return nil, err
			}
		}

		log.Printf("Add changelog to file...")
		changelogMd, err := AddChangelogToFile(changelog, version, packagePath, generateParam.ReleaseDate)
		if err != nil {
			return nil, err
		}

		log.Printf("Update module definition if v2+...")
		err = UpdateModuleDefinition(packagePath, generateParam.RPName, generateParam.NamespaceName, version)
		if err != nil {
			return nil, err
		}

		oldModuleVersion, err := getModuleVersion(filepath.Join(packagePath, "autorest.md"))
		if err != nil {
			return nil, err
		}

		log.Printf("Replace version in autorest.md and constants...")
		if err = ReplaceVersion(packagePath, version.String()); err != nil {
			return nil, err
		}

		if _, err := os.Stat(filepath.Join(packagePath, "fake")); !os.IsNotExist(err) && oldModuleVersion.Major() != version.Major() {
			log.Printf("Replace fake module v2+...")
			if err = replaceModuleImport(packagePath, generateParam.RPName, generateParam.NamespaceName, oldModuleVersion.String(), version.String(),
				"fake", ".go"); err != nil {
				return nil, err
			}
		}

		// When sdk has major version bump, the live test needs to update the module referenced in the code.
		if oldModuleVersion.Major() != version.Major() && existSuffixFile(packagePath, "_live_test.go") {
			log.Printf("Replace live test module v2+...")
			if err = replaceModuleImport(packagePath, generateParam.RPName, generateParam.NamespaceName, oldModuleVersion.String(), version.String(),
				"", "_live_test.go"); err != nil {
				return nil, err
			}
		}

		log.Printf("Replace README.md module...")
		if err = replaceReadmeModule(packagePath, generateParam.RPName, generateParam.NamespaceName, version.String()); err != nil {
			return nil, err
		}

		log.Printf("Replace README.md NewClient name...")
		if err = ReplaceReadmeNewClientName(packagePath, newExports); err != nil {
			return nil, err
		}

		// Example generation should be the last step because the package import relay on the new calculated version
		if !generateParam.SkipGenerateExample {
			log.Printf("Generate examples...")
			flag, err := GetAlwaysSetBodyParamRequiredFlag(filepath.Join(packagePath, "build.go"))
			if err != nil {
				return nil, err
			}
			if err := ExecuteExampleGenerate(packagePath, filepath.Join("resourcemanager", generateParam.RPName, generateParam.NamespaceName), flag); err != nil {
				return nil, err
			}
		}

		return &GenerateResult{
			Version:           version.String(),
			RPName:            generateParam.RPName,
			PackageName:       generateParam.NamespaceName,
			PackageAbsPath:    packagePath,
			Changelog:         *changelog,
			ChangelogMD:       changelogMd + "\n" + changelog.GetChangeSummary(),
			PullRequestLabels: string(prl),
		}, nil
	}
}
