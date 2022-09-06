// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package profile

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/template"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/processor"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var (
	moduleLineRegex = regexp.MustCompile(`moduleName\s*=\s*".+"`)
)

func generateSDK(sdkRepoPath, swaggerRepoPath, profileName string, flags SDKFlags) error {
	log.Printf("Generating profile SDK for '%s'...", profileName)

	specCommitHash, err := processor.GetSpecCommit(swaggerRepoPath)
	if err != nil {
		return err
	}

	config, err := ReadConfig(path.Join(sdkRepoPath, "sdk", "profiles", profileName, "definition.json"))
	if err != nil {
		return fmt.Errorf("Cannot get resource map: %+v", err)
	}

	packageVersion := "1.0.0"
	if flags.VersionNumber != "" {
		packageVersion = flags.VersionNumber
	}

	for _, moduleProperty := range config.Modules {
		log.Printf("Generating SDK for '%s/%s'...", moduleProperty.RP, moduleProperty.Namespace)

		packagePath := filepath.Join(sdkRepoPath, "sdk", "profiles", profileName, "resourcemanager", moduleProperty.RP, moduleProperty.Namespace)

		log.Printf("Generate config for package '%s'", packagePath)

		if err = template.GeneratePackageByTemplate(moduleProperty.RP, moduleProperty.Namespace, template.Flags{
			SDKRoot:        sdkRepoPath,
			TemplatePath:   "eng/tools/generator/template/profiles/package",
			PackagePath:    "sdk/profiles/" + profileName + "/resourcemanager/" + moduleProperty.RP + "/" + moduleProperty.Namespace,
			PackageTitle:   profileName,
			Commit:         specCommitHash,
			PackageConfig:  "tag: " + moduleProperty.Tag + "\n" + moduleProperty.AdditionalConfig,
			GoVersion:      flags.GoVersion,
			PackageVersion: packageVersion,
		}); err != nil {
			return err
		}

		log.Printf("Remove all the generated files ...")
		if err = processor.CleanSDKGeneratedFiles(packagePath); err != nil {
			return err
		}

		log.Printf("Change swagger config in `autorest.md` according to repo URL and commit ID...")
		autorestMdPath := filepath.Join(packagePath, "autorest.md")
		if err := processor.ChangeConfigWithCommitID(autorestMdPath, common.DefaultSpecRepo, specCommitHash, moduleProperty.SpecName); err != nil {
			return err
		}

		log.Printf("Run `go generate` to regenerate the code...")
		if err := processor.ExecuteGoGenerate(packagePath); err != nil {
			return err
		}

		if err := RefinePackage(packagePath, moduleProperty.Namespace, profileName); err != nil {
			return err
		}
	}

	packageReadmePath := filepath.Join(sdkRepoPath, "sdk", "profiles", profileName, "README.md")

	if _, err := os.Stat(packageReadmePath); os.IsNotExist(err) {
		log.Printf("Base files '%s' not exist, generate template", packageReadmePath)

		if err = template.GeneratePackageByTemplate(profileName, profileName, template.Flags{
			SDKRoot:        sdkRepoPath,
			TemplatePath:   "eng/tools/generator/template/profiles/base",
			PackagePath:    "sdk/profiles/" + profileName,
			PackageTitle:   profileName,
			Commit:         specCommitHash,
			PackageConfig:  "",
			GoVersion:      common.DefaultGoVersion,
			PackageVersion: "1.0.0",
		}); err != nil {
			return err
		}
	}

	if err := processor.ExecuteGoGenerate(filepath.Join(sdkRepoPath, "sdk", "profiles", profileName)); err != nil {
		return err
	}

	return nil
}

func RefinePackage(packagePath, packageName, profileNameForPackage string) error {
	if err := os.Remove(filepath.Join(packagePath, "go.mod")); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(packagePath, "go.sum")); err != nil {
		return err
	}

	constantsPath := filepath.Join(packagePath, "constants.go")
	b, err := os.ReadFile(constantsPath)
	if err != nil {
		return err
	}
	contents := moduleLineRegex.ReplaceAllString(string(b), "moduleName = \""+profileNameForPackage+"/"+packageName+"\"")
	return os.WriteFile(constantsPath, []byte(contents), 0644)
}
