// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/Azure/azure-sdk-for-go/tools/internal/pkgs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initCommand() *cobra.Command {
	init := &cobra.Command{
		Use:   "init <searching dir>",
		Short: "Initialize a package into go module with initial version",
		Long: `This tool will detect every possible service under the searching directory, 
and make them as module with initial version. 
The default version for new stable modules is v1.0.0 and for new preview modules is v0.0.0.
NOTE: This command is only used on local and only for initial release.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]
			repoRoot := viper.GetString("gomod-root")
			classicalFile := viper.GetString("classical")
			return ExecuteInit(root, repoRoot, classicalFile)
		},
	}
	// register flags
	flags := init.Flags()
	flags.StringP("classical", "c", "", "file for classical package list, these packages will be skipped")
	if err := viper.BindPFlag("classical", flags.Lookup("classical")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}

	return init
}

const (
	initialGoMod = `module %s

%s
`

	goVersion = `go 1.13`
)

func ExecuteInit(r, repoRoot, classicalFile string) error {
	root, err := filepath.Abs(r)
	if err != nil {
		return fmt.Errorf("failed to get absolute path from '%s': %+v", r, err)
	}
	classicalPackages, err := loadClassicalPackages(classicalFile)
	if err != nil {
		return fmt.Errorf("failed to load classical packages: %+v", err)
	}
	ps, err := pkgs.GetPkgs(root)
	if err != nil {
		return fmt.Errorf("failed to get packages: %+v", err)
	}
	var errs []error
	for _, p := range ps {
		// test if this is a classical package
		_, classical := classicalPackages[p.Dest]
		if classical {
			log.Infof("Skipping classical package: %s", p.Dest)
			continue
		}
		path := filepath.Join(root, p.Dest)
		tagPrefix, err := getTagPrefix(path, repoRoot)
		if err != nil {
			return fmt.Errorf("failed to get tag prefix: %+v", err)
		}
		// get tag and ver
		startingVer := getStartingVer(p)
		tag := tagPrefix + "/" + startingVer
		ver := versionGoRegex.FindString(startingVer)
		if !classical {
			if err := modifyVersionFile(root, p, tag, ver); err != nil {
				errs = append(errs, err)
			}
		}
		if err := createGoModFile(root, p); err != nil {
			errs = append(errs, err)
		}
		log.Infof("Created module: %s", tag)
	}
	// handle errors
	if len(errs) == 0 {
		return nil
	}
	for _, err := range errs {
		log.Errorln(err.Error())
	}
	return fmt.Errorf("execution failed with %d errors", len(errs))
}

func modifyVersionFile(root string, p pkgs.Pkg, tag, ver string) error {
	verFilePath := filepath.Join(root, p.Dest, versionFilename)
	b, err := ioutil.ReadFile(verFilePath)
	if err != nil {
		return fmt.Errorf("failed to read version file '%s': %+v", verFilePath, err)
	}
	content := string(b)
	// remove the import clause
	content = strings.Replace(content, `import "github.com/Azure/azure-sdk-for-go/version"`, "", 1)
	// replace the first `version.Number` to `Version()`
	content = strings.Replace(content, "version.Number", "Version()", 1)
	// replace the second `version.Number` to the value of ver
	content = strings.Replace(content, "version.Number", fmt.Sprintf(`"%s"`, ver), 1)
	content = content + fmt.Sprintf("\n// tag: %s\n", tag)
	if err := ioutil.WriteFile(verFilePath, []byte(content), 0666); err != nil {
		return fmt.Errorf("failed to write version file '%s': %+v", verFilePath, err)
	}
	return nil
}

func createGoModFile(root string, p pkgs.Pkg) error {
	modFilePath := filepath.Join(root, p.Dest, goModFilename)
	fullPath := filepath.Join(root, p.Dest)
	index := strings.Index(fullPath, "github.com")
	if index < 0 {
		return fmt.Errorf("failed to find github.com in filepath %s", fullPath)
	}
	importPath := strings.ReplaceAll(fullPath[index:], "\\", "/")
	content := fmt.Sprintf(initialGoMod, importPath, goVersion)
	err := ioutil.WriteFile(modFilePath, []byte(content), 0755)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %+v", modFilePath, err)
	}
	return nil
}

func getStartingVer(p pkgs.Pkg) string {
	if p.IsPreviewPackage() {
		return startingModVerPreview
	}
	return startingModVer
}

func loadClassicalPackages(exceptFile string) (map[string]bool, error) {
	result := make(map[string]bool)
	if exceptFile == "" {
		return result, nil
	}
	abs, err := filepath.Abs(exceptFile)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(abs)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		result[line] = true
	}
	return result, nil
}
