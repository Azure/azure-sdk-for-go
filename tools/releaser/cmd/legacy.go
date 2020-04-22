package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/files"
	"github.com/Azure/azure-sdk-for-go/tools/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func legacyCommand() *cobra.Command {
	legacyCommand := &cobra.Command{
		Use:   "legacy <searching dir>",
		Short: "Revert the version.go files to legacy versioning",
		Long: `This tool searches the searching directory for version.go files, and revert them 
to use the one version number (version.Number) instead of using their own module version number.`,
		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			root := args[0]
			flags := LegacyFlags{
				RepoRoot: viper.GetString("repo-root"),
			}
			return ExecuteLegacy(root, flags)
		},
	}
	// flags
	flags := legacyCommand.Flags()
	flags.String("repo-root", "github.com/Azure/azure-sdk-for-go", "the root path of the go SDK")
	if err := viper.BindPFlag("repo-root", flags.Lookup("repo-root")); err != nil {
		log.Fatalf("failed to bind flag: %+v", err)
	}

	return legacyCommand
}

type LegacyFlags struct {
	RepoRoot string
}

func (f LegacyFlags) apply() {
	repoRoot = f.RepoRoot
}

var (
	repoRoot string

	versionStatementRegex = regexp.MustCompile(`^[\t ]*return "v?\d+\.\d+\.\d+"[\t ]*$`)
)

func ExecuteLegacy(r string, flags LegacyFlags) error {
	flags.apply()
	root, err := filepath.Abs(r)
	if err != nil {
		return fmt.Errorf("failed to get absolute root: %+v", err)
	}
	// find all version.go files
	allFiles, err := findAllFiles(root, versionFile)
	if err != nil {
		return fmt.Errorf("failed to find all version.go files: %+v", err)
	}
	// iterate all version.go files to replace their version number by version.Number
	for _, file := range allFiles {
		if err := revertVersionFile(file); err != nil {
			return fmt.Errorf("failed to revert file '%s': %+v", file, err)
		}
		if err := formatCode(filepath.Dir(file)); err != nil {
			return fmt.Errorf("failed to format '%s': %+v", file, err)
		}
	}
	return nil
}

func revertVersionFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open for read '%s': %+v", path, err)
	}
	defer file.Close()
	if err := revertVersionFileContent(file); err != nil {
		return fmt.Errorf("failed to update version.go content: %+v", err)
	}
	return nil
}

func revertVersionFileContent(versionFile io.ReadWriteSeeker) error {
	lines := files.GetLines(versionFile)

	if _, err := versionFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to start: %v", err)
	}

	hasImport := false

	for _, line := range lines {
		if strings.HasPrefix(line, "import ") {
			hasImport = true
		}
		if strings.HasPrefix(line, "// ") && !hasImport {
			hasImport = true
			if _, err := fmt.Fprintln(versionFile, fmt.Sprintf(`import "%s/version"`, repoRoot)); err != nil {
				return fmt.Errorf("failed to write import line")
			}
		}

		if matched := versionStatementRegex.MatchString(line); matched {
			line = "return version.Number"
		}

		if _, err := fmt.Fprintln(versionFile, line); err != nil {
			return fmt.Errorf("failed to write line: %s", line)
		}
	}

	return nil
}

func formatCode(directory string) error {
	c := exec.Command("gofmt", "-w", directory)
	if output, err := c.CombinedOutput(); err != nil {
		return errors.New(string(output))
	}
	return nil
}
