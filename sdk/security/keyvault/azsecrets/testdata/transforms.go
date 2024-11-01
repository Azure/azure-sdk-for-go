package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// removing client prefix from types
func regexReplace(fileName string, regex string, replace string) {
	r, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	file = r.ReplaceAll(file, []byte(replace))

	err = os.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// delete client name prefix from method options and response types
	clientRegex := `Client(\w+)((?:Options|Response))`
	regexReplace("client.go", clientRegex, `$1$2`)
	regexReplace("options.go", clientRegex, `$1$2`)
	regexReplace("responses.go", clientRegex, `$1$2`)

	// delete the version path param check (version == "" is legal for Key Vault but indescribable by OpenAPI)
	regexReplace("client.go", `\sif secretVersion == "" \{\s+.+secretVersion cannot be empty"\)\s+\}\s`, "")

	// make secret IDs a convenience type so we can add parsing methods
	regexReplace("models.go", `\sID \*string(\s+.*)`, "ID *ID$1")
	regexReplace("models.go", `\sKID \*string(\s+.*)`, "KID *ID$1")
	// regexReplace("models.go", `\sID \*string`, "ID *ID")
	// regexReplace("models.go", `\sKID \*string`, "KID *ID")
}
