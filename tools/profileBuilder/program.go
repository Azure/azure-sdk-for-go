package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"

	"github.com/marstr/collection"
	goalias "github.com/marstr/goalias/model"
	"github.com/marstr/randname"
)

var (
	profileName     string
	outputLocation  string
	inputRoot       string
	inputList       io.Reader
	packageStrategy collection.Enumerable
	outputLog       *log.Logger
)

// WellKnownStrategy is an Enumerable which lists all known strategies for choosing packages for a profile.
type WellKnownStrategy string

// This block declares the definitive list of WellKnownStrategies
const (
	WellKnownStrategyList   WellKnownStrategy = "list"
	WellKnownStrategyLatest WellKnownStrategy = "latest"
)

func main() {
	for entry := range packageStrategy.Enumerate(nil) {
		pkg, ok := entry.(*ast.Package)
		if !ok {
			outputLog.Print("Unexpected type when searching for type `*ast.Package`: ", reflect.TypeOf(entry).Name())
			continue
		}
		outputLog.Print("Package Found: ", pkg.Name)

		alias, err := goalias.NewAliasPackage(pkg)
		if err != nil {
			outputLog.Print("Could not create alias because: ", err)
			continue
		}

		files := token.NewFileSet()
		printer.Fprint(os.Stdout, files, alias.ModelFile())
	}
}

func init() {
	const defaultName = "<randomly generated>"

	var selectedStrategy string
	var inputListLocation string
	var useVerbose bool

	flag.StringVar(&profileName, "name", defaultName, "The name that should be given to the generated profile.")
	flag.StringVar(&outputLocation, "o", getDefaultOutputLocation(), "The output location for the package generated as a profile.")
	flag.StringVar(&inputRoot, "root", getDefaultInputRoot(), "The location of the Azure REST OpenAPI Specs repository.")
	flag.StringVar(&inputListLocation, "l", "", "If the `list` strategy is chosen, -l is the location of the file to read for said list. If not present, stdin is used.")
	flag.StringVar(&selectedStrategy, "s", string(WellKnownStrategyLatest), "The strategy to employ for finding packages to put in a profile.")
	flag.BoolVar(&useVerbose, "v", false, "Write status to stderr as the program progresses")
	flag.Parse()

	if profileName == defaultName {
		profileName = randname.AdjNoun{}.Generate()
	}

	var logWriter io.Writer
	if useVerbose {
		logWriter = os.Stderr
	} else {
		logWriter = ioutil.Discard
	}
	outputLog = log.New(logWriter, "[STATUS] ", 0)

	outputLog.Print("Status Logging Enabled")

	inputList = os.Stdin
	if inputListLocation != "" {
		var err error
		inputList, err = os.Open(inputListLocation)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	wellKnownStrategies := map[WellKnownStrategy]collection.Enumerable{
		WellKnownStrategyList:   ListStrategy{Reader: inputList},
		WellKnownStrategyLatest: LatestStrategy{Root: inputRoot},
	}

	if s, ok := wellKnownStrategies[WellKnownStrategy(selectedStrategy)]; ok {
		packageStrategy = s
	} else {
		fmt.Fprintf(os.Stderr, "Error: Unknown strategy for identifying packages: %s\n", selectedStrategy)
		os.Exit(1)
	}
}

func getDefaultOutputLocation() string {
	return path.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"Azure",
		"azure-sdk-for-go",
		"profile",
	)
}

func getDefaultInputRoot() string {
	return path.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"Azure",
		"azure-rest-api-specs",
	)
}
