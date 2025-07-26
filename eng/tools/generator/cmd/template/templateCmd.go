// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command returns the template command
func Command() *cobra.Command {
	templateCmd := &cobra.Command{
		Use:   "template (<rpName> <packageName>) | <packagePath>",
		Short: "Onboards new RP with the template",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var rpName, packageName string
			if len(args) == 1 {
				segments := strings.Split(args[0], "/")
				if len(segments) != 2 {
					return fmt.Errorf("%s is not a valid package path. Please assign the package path using `rpName/packageName` format", args[0])
				}
				rpName = segments[0]
				packageName = segments[1]
			} else {
				rpName = args[0]
				packageName = args[1]
			}

			return GeneratePackageByTemplate(rpName, packageName, ParseFlags(cmd.Flags()))
		},
	}

	BindFlags(templateCmd.Flags())
	if err := templateCmd.MarkFlagRequired("package-title"); err != nil {
		log.Fatal(err)
	}
	if err := templateCmd.MarkFlagRequired("commit"); err != nil {
		log.Fatal(err)
	}

	return templateCmd
}

// BindFlags binds the flags to this command
func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.String("go-sdk-folder", ".", "Specifies the path of root of azure-sdk-for-go")
	flagSet.String("template-path", "eng/tools/generator/template/rpName/packageName", "Specifies the path of the template")
	flagSet.String("package-title", "", "Specifies the title of this package")
	flagSet.String("commit", "", "Specifies the commit hash of azure-rest-api-specs")
	flagSet.String("release-date", "", "Specifies the release date in changelog")
	flagSet.String("package-config", "", "Additional config for package")
	flagSet.String("package-version", "", "Specify the version number of this release")
}

// ParseFlags parses the flags to a Flags struct
func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		SDKRoot:        flags.GetString(flagSet, "go-sdk-folder"),
		TemplatePath:   flags.GetString(flagSet, "template-path"),
		PackageTitle:   flags.GetString(flagSet, "package-title"),
		Commit:         flags.GetString(flagSet, "commit"),
		ReleaseDate:    flags.GetString(flagSet, "release-date"),
		PackageConfig:  flags.GetString(flagSet, "package-config"),
		PackageVersion: flags.GetString(flagSet, "package-version"),
	}
}

// Flags ...
type Flags struct {
	SDKRoot        string
	TemplatePath   string
	PackageTitle   string
	Commit         string
	ReleaseDate    string
	PackageConfig  string
	PackageVersion string
}

// GeneratePackageByTemplate creates a new set of files based on the things in template directory
func GeneratePackageByTemplate(rpName, packageName string, flags Flags) error {
	root, err := filepath.Abs(flags.SDKRoot)
	if err != nil {
		return fmt.Errorf("cannot get the root of azure-sdk-for-go from '%s': %+v", flags.SDKRoot, err)
	}
	var absTemplateDir string
	if filepath.IsAbs(flags.TemplatePath) {
		absTemplateDir = flags.TemplatePath
	} else {
		absTemplateDir = filepath.Join(root, flags.TemplatePath)
	}
	fileList, err := os.ReadDir(absTemplateDir)
	if err != nil {
		return fmt.Errorf("cannot read the directory '%s': %+v", absTemplateDir, err)
	}

	// build the replaceMap
	buildReplaceMap(rpName, packageName, flags.PackageConfig, flags.PackageTitle, flags.Commit, flags.ReleaseDate, flags.PackageVersion)

	// copy everything to destination directory
	for _, file := range fileList {
		path := filepath.Join(absTemplateDir, file.Name())
		content, err := readAndReplace(path)
		if err != nil {
			return err
		}

		dirPath := filepath.Join(root, "sdk", "resourcemanager", rpName, packageName)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("cannot create directory '%s': %+v", dirPath, err)
		}

		newFilePath := filepath.Join(dirPath, strings.TrimSuffix(file.Name(), FilenameSuffix))
		if err := createNewFile(newFilePath, content); err != nil {
			return err
		}
	}

	return nil
}

func buildReplaceMap(rpName, packageName, packageConfig, packageTitle, commitID, releaseDate, packageVersion string) {
	replaceMap = make(map[string]string)

	replaceMap[RPNameKey] = rpName
	replaceMap[PackageNameKey] = packageName
	replaceMap[PackageConfigKey] = packageConfig
	replaceMap[PackageTitleKey] = packageTitle
	replaceMap[CommitIDKey] = commitID
	if releaseDate == "" {
		replaceMap[ReleaseDate] = time.Now().Format("2006-01-02")
	} else {
		replaceMap[ReleaseDate] = releaseDate
	}
	replaceMap[PackageVersion] = packageVersion
}

func readAndReplace(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read from file '%s': %+v", path, err)
	}

	content := string(b)
	for k, v := range replaceMap {
		content = strings.ReplaceAll(content, k, v)
	}

	return content, nil
}

func createNewFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file '%s': %+v", path, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("cannot write to file '%s': %+v", path, err)
	}

	return nil
}

var (
	replaceMap map[string]string
)

const (
	RPNameKey        = "{{rpName}}"
	PackageNameKey   = "{{packageName}}"
	PackageTitleKey  = "{{packageTitle}}"
	CommitIDKey      = "{{commitID}}"
	FilenameSuffix   = ".tpl"
	ReleaseDate      = "{{releaseDate}}"
	PackageConfigKey = "{{packageConfig}}"
	PackageVersion   = "{{packageVersion}}"
)
