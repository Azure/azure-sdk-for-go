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
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
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

	// typespec
	TypeSpecConfig *typespec.TypeSpecConfig
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
	RPName               string
	NamespaceName        string
	NamespaceConfig      string
	SpecificVersion      string
	SpecificPackageTitle string
	SpecRPName           string
	ReleaseDate          string
	SkipGenerateExample  bool
	GoVersion            string
	RemoveTagSet         bool
	ForceStableVersion   bool
	TypeSpecEmitOption   string
	TspClientOptions     []string
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
			log.Printf("Start to process rp: %s, namespace: %s", rpName, packageInfo.Name)
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
				errors = append(errors, fmt.Errorf("failed to generate for rp: %s, namespace: %s: %+v", rpName, packageInfo.Name, err))
				continue
			}
			result = append(result, *singleResult)
		}
	}
	return result, errors
}

func (ctx *GenerateContext) GenerateForSingleRPNamespace(generateParam *GenerateParam) (*GenerateResult, error) {
	packagePath := filepath.Join(ctx.SDKPath, "sdk", "resourcemanager", generateParam.RPName, generateParam.NamespaceName)
	changelogPath := filepath.Join(packagePath, ChangelogFileName)

	onBoard := false

	version, err := semver.NewVersion("0.1.0")
	if err != nil {
		return nil, err
	}
	if generateParam.SpecificVersion != "" {
		log.Printf("Use specific version: %s", generateParam.SpecificVersion)
		version, err = semver.NewVersion(generateParam.SpecificVersion)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		onBoard = true
		log.Printf("Package '%s' changelog not exist, do onboard process", packagePath)

		if generateParam.SpecificPackageTitle == "" {
			generateParam.SpecificPackageTitle = strings.Title(generateParam.RPName)
		}

		log.Printf("Start to use template to generate new rp folder and basic package files...")
		if err = template.GeneratePackageByTemplate(generateParam.RPName, generateParam.NamespaceName, template.Flags{
			SDKRoot:        ctx.SDKPath,
			TemplatePath:   "eng/tools/generator/template/rpName/packageName",
			PackageTitle:   generateParam.SpecificPackageTitle,
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

	log.Printf("Start to run `go generate` to regenerate the code...")
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

	log.Printf("Start to generate changelog for package...")
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
			log.Printf("Start to generate examples...")
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
		if generateParam.SpecificVersion == "" {
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
			log.Printf("Start to generate examples...")
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

func (ctx *GenerateContext) GenerateForTypeSpec(generateParam *GenerateParam) (*GenerateResult, error) {
	packagePath := filepath.Join(ctx.SDKPath, "sdk", "resourcemanager", generateParam.RPName, generateParam.NamespaceName)
	changelogPath := filepath.Join(packagePath, ChangelogFileName)

	version, err := semver.NewVersion("0.1.0")
	if err != nil {
		return nil, err
	}
	if generateParam.SpecificVersion != "" {
		log.Printf("Use specific version: %s", generateParam.SpecificVersion)
		version, err = semver.NewVersion(generateParam.SpecificVersion)
		if err != nil {
			return nil, err
		}
	}

	onBoard := false
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		onBoard = true
		log.Printf("Package '%s' changelog not exist, do onboard process", packagePath)
		if generateParam.SpecificPackageTitle == "" {
			generateParam.SpecificPackageTitle = strings.Title(generateParam.RPName)
		}

		log.Printf("Start to use template to generate new rp folder and basic package files...")
		sdkBasicInfo := map[string]any{
			"rpName":         generateParam.RPName,
			"packageName":    generateParam.NamespaceName,
			"packageTitle":   generateParam.SpecificPackageTitle,
			"packageVersion": version.String(),
			"releaseDate":    generateParam.ReleaseDate,
			"goVersion":      generateParam.GoVersion,
		}
		err = typespec.ParseTypeSpecTemplates(filepath.Join(ctx.SDKPath, "eng/tools/generator/template/typespec"), packagePath, sdkBasicInfo, nil)
		if err != nil {
			return nil, err
		}
	} else {
		log.Printf("Package '%s' existed, do update process", packagePath)
		log.Printf("Remove all the generated files ...")
		if err = CleanSDKGeneratedFiles(packagePath); err != nil {
			return nil, err
		}
	}

	log.Printf("Start to run `tsp-client init` to generate the code...")
	defaultModuleVersion := version.String()
	emitOption := fmt.Sprintf("module-version=%s", defaultModuleVersion)
	if generateParam.TypeSpecEmitOption != "" {
		emitOption = fmt.Sprintf("%s;%s", emitOption, generateParam.TypeSpecEmitOption)
	}
	err = ExecuteTypeSpecGenerate(ctx, emitOption, generateParam.TspClientOptions)
	if err != nil {
		return nil, err
	}

	previousVersion := ""
	isCurrentPreview := false
	var oriExports *exports.Content
	if generateParam.SpecificVersion != "" {
		isCurrentPreview, err = IsBetaVersion(version.String())
		if err != nil {
			return nil, err
		}
	} else {
		isCurrentPreview, err = ContainsPreviewAPIVersion(packagePath)
		if err != nil {
			return nil, err
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

	log.Printf("Start to generate changelog for package...")
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

			log.Printf("Replace version in constants.go...")
			if err = ReplaceConstModuleVersion(packagePath, version.String()); err != nil {
				return nil, err
			}
			prl = FirstGALabel
		}

		log.Printf("##[command]Executing gofmt -s -w . in %s\n", packagePath)
		if err = ExecuteGoFmt(packagePath, "-s", "-w", "."); err != nil {
			return nil, err
		}

		log.Printf("##[command]Executing go mod tidy in %s\n", packagePath)
		if err = ExecuteGo(packagePath, "mod", "tidy"); err != nil {
			return nil, err
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
		if generateParam.SpecificVersion == "" {
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

		log.Printf("Replace version in constants.go...")
		if err = ReplaceConstModuleVersion(packagePath, version.String()); err != nil {
			return nil, err
		}

		oldModuleVersion, err := semver.NewVersion(defaultModuleVersion)
		if err != nil {
			return nil, err
		}

		baseModule := fmt.Sprintf("%s/%s/%s", MgmtSDKModulePrefix, generateParam.RPName, generateParam.NamespaceName)
		if _, err := os.Stat(filepath.Join(packagePath, "fake")); !os.IsNotExist(err) && oldModuleVersion.Major() != version.Major() {
			log.Printf("Replace fake module v2+...")
			if err = ReplaceModule(version, packagePath, baseModule, ".go"); err != nil {
				return nil, err
			}
		}

		// When sdk has major version bump, the live test needs to update the module referenced in the code.
		if existSuffixFile(packagePath, "_live_test.go") {
			log.Printf("Replace live test module v2+...")
			if err = ReplaceModule(version, packagePath, baseModule, "_live_test.go"); err != nil {
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
		}

		// remove autorest.md and build.go
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if _, err := os.Stat(autorestMdPath); !os.IsNotExist(err) {
			log.Println("Remove autorest.md...")
			if err = os.Remove(autorestMdPath); err != nil {
				return nil, err
			}

		}
		buildGoPath := filepath.Join(packagePath, "build.go")
		if _, err := os.Stat(buildGoPath); !os.IsNotExist(err) {
			log.Println("Remove build.go...")
			if err = os.Remove(buildGoPath); err != nil {
				return nil, err
			}
		}

		log.Printf("##[command]Executing gofmt -s -w . in %s\n", packagePath)
		if err = ExecuteGoFmt(packagePath, "-s", "-w", "."); err != nil {
			return nil, err
		}

		log.Printf("##[command]Executing go mod tidy in %s\n", packagePath)
		if err = ExecuteGo(packagePath, "mod", "tidy"); err != nil {
			return nil, err
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
