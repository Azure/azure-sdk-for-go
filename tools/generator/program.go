package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/marstr/randname"

	"github.com/marstr/collection"
)

var (
	targetFiles    collection.Enumerable
	packageVersion string
	outputBase     string
	errLog         *log.Logger
	statusLog      *log.Logger
	dryRun         bool
)

func main() {

	type generationTuple struct {
		fileName     string
		packageName  string
		outputFolder string
	}

	literateFiles := collection.Where(targetFiles, func(subject interface{}) bool {
		return strings.EqualFold(path.Base(subject.(string)), "README.md")
	})

	tuples := collection.SelectMany(literateFiles,
		// The following function compiles the regexp which finds Go related package tags in a literate file, and creates a collection.Unfolder.
		// This function has been declared this way so that the relatively expensive act of compiling a regexp is only done once.
		func() collection.Unfolder {
			const patternText = "```" + `\s+yaml\s+\$\(\s*tag\s*\)\s*==\s*'([\d\w\-_]+)'\s*&&\s*\$\(go\)\s+output-folder:\s+\$\(go-sdk-folder\)([\w\d\-_\\/]+)\s+` + "```"

			goConfigPattern := regexp.MustCompile(patternText)

			// This function is a collection.Unfolder which takes a literate file as a path, and retrieves all configuration which applies to a package tag and Go.
			return func(subject interface{}) collection.Enumerator {
				results := make(chan interface{})

				go func() {
					defer close(results)

					targetContents, err := ioutil.ReadFile(subject.(string))
					if err != nil {
						errLog.Printf("Skipping %q because: %v", subject.(string), err)
						return
					}

					matches := goConfigPattern.FindAllStringSubmatch(string(targetContents), -1)

					if len(matches) == 0 {
						statusLog.Printf("Skipping %q because there were no package tags with go configuration found.", subject.(string))
					} else {
						for _, submatch := range matches {
							results <- generationTuple{
								fileName:     subject.(string),
								packageName:  strings.Replace(submatch[1], `\`, "/", -1),
								outputFolder: strings.Replace(submatch[2], `\`, "/", -1),
							}
						}
					}
				}()

				return results
			}
		}())

	if dryRun {
		for entry := range tuples.Enumerate(nil) {
			tuple := entry.(generationTuple)
			fmt.Printf("%q in %q to %q\n", tuple.packageName, tuple.fileName, tuple.outputFolder)
		}
	} else {
		randGener := randname.NewAdjNoun() // TODO: print log files to a location indicative of the package they refer to instead of a random location.
		logDirLocation, err := ioutil.TempDir("", "az-go-sdk-logs")
		logDirLocation = strings.Replace(logDirLocation, `\`, "/", -1)
		if err == nil {
			statusLog.Print("Generation logs can be found at:", logDirLocation)
		} else {
			errLog.Print("Fatal: Could not create temporary directory for AutoRest logs.")
			return
		}

		done := make(chan struct{})
		generated := tuples.Enumerate(done).ParallelSelect(func(subject interface{}) interface{} {
			tuple := subject.(generationTuple)
			args := []string{
				tuple.fileName,
				"--go",
				fmt.Sprintf("--go-sdk-folder='%s'", outputBase),
				"--verbose",
				"--tag=" + tuple.packageName,
			}

			if packageVersion != "" {
				args = append(args, fmt.Sprintf("--package-version='%s'", packageVersion))
			}

			logFileLoc := filepath.Join(logDirLocation, randGener.Generate()+".txt")
			logFile, err := os.Create(logFileLoc)
			if err == nil {
				statusLog.Printf("Logs for %q in %q can be found at: %s", tuple.packageName, tuple.fileName, logFileLoc)
			} else {
				errLog.Printf("Could not create log file %q for AutoRest generating from: %q in %q", logFileLoc, tuple.packageName, tuple.fileName)
				return nil
			}

			commandText := new(bytes.Buffer)
			fmt.Fprint(commandText, "Executing Command: \"")
			fmt.Fprint(commandText, "autorest ")
			for _, a := range args {
				fmt.Fprint(commandText, a)
				fmt.Fprint(commandText, " ")
			}
			commandText.Truncate(commandText.Len() - 1)
			fmt.Fprint(commandText, `"`)

			fmt.Fprintln(logFile, commandText.String())

			genProc := exec.Command("autorest", args...)
			genProc.Stdout = logFile
			genProc.Stderr = logFile

			err = genProc.Run()
			if err != nil {
				fmt.Println(logFile, "Autorest Exectution Error: ", err)
			}
			return err == nil
		})

		for range generated {
			// Intentionally Left Blank
		}
		close(done)
	}
}

func init() {
	var useRecursive, useStatus bool

	flag.BoolVar(&useRecursive, "r", false, "Recursively traverses the directories specified looking for literate files.")
	flag.StringVar(&outputBase, "o", getDefaultOutputBase(), "The root directory to use for the output of generated files. i.e. The value to be treated as the go-sdk-folder when AutoRest is called.")
	flag.BoolVar(&useStatus, "v", false, "Print status messages as generation takes place.")
	flag.BoolVar(&dryRun, "p", false, "Preview which packages would be generated instead of actaully calling autorest.")
	flag.Parse()

	statusWriter := ioutil.Discard
	if useStatus {
		statusWriter = os.Stdout
	}
	statusLog = log.New(statusWriter, "[STATUS] ", 0)
	errLog = log.New(os.Stderr, "[ERROR] ", 0)

	targetFiles = collection.AsEnumerable(flag.Args())
	targetFiles = collection.SelectMany(targetFiles, func(subject interface{}) collection.Enumerator {
		cast, ok := subject.(string)

		if !ok {
			return collection.Empty.Enumerate(nil)
		}
		pathInfo, err := os.Stat(cast)
		if err != nil {
			return collection.Empty.Enumerate(nil)
		}

		if pathInfo.IsDir() {
			traverser := collection.Directory{
				Location: cast,
				Options:  collection.DirectoryOptionsExcludeDirectories,
			}

			if useRecursive {
				traverser.Options |= collection.DirectoryOptionsRecursive
			}

			return traverser.Enumerate(nil)
		}

		return collection.AsEnumerable(cast).Enumerate(nil)
	})

	targetFiles = collection.Select(targetFiles, func(subject interface{}) interface{} {
		return strings.Replace(subject.(string), `\`, "/", -1)
	})
}

func getDefaultOutputBase() string {
	return strings.Replace(path.Join(os.Getenv("GOPATH"), "src", "github.com", "Azure", "azure-sdk-for-go"), `\`, "/", -1)
}
