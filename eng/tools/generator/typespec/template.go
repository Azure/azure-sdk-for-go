package typespec

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

/*
const (
	RPNameKey        = "{{rpName}}"
	PackageNameKey   = "{{packageName}}"
	PackageTitleKey  = "{{packageTitle}}"
	CommitIDKey      = "{{commitID}}"
	FilenameSuffix   = ".tpl"
	ReleaseDate      = "{{releaseDate}}"
	PackageConfigKey = "{{packageConfig}}"
	GoVersion        = "{{goVersion}}"
	PackageVersion   = "{{packageVersion}}"
)
*/
func ParseTypeSpecTemplates(templateDir, outputDir string, data any, funcMap template.FuncMap) error {

	tpl := template.New("parse.tpl").Funcs(funcMap)
	tpl, err := tpl.ParseGlob(filepath.Join(templateDir, "*.tpl"))
	if err != nil {
		return err
	}
	for _, t := range tpl.Templates() {
		fName, _ := strings.CutSuffix(t.Name(), ".tpl")
		w, err := os.OpenFile(filepath.Join(outputDir, fName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer w.Close()

		err = tpl.ExecuteTemplate(w, t.Name(), data)
		if err != nil {
			return err
		}
	}

	return nil
}
