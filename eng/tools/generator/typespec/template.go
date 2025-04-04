// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func ParseTypeSpecTemplates(templateDir, outputDir string, data map[string]any, funcMap template.FuncMap) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if data["releaseDate"] == "" {
		data["releaseDate"] = time.Now().Format("2006-01-02")
	}

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

func ParseTypeSpecTemplate(templatePath, outputDir string, data map[string]any, funcMap template.FuncMap) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	tpl := template.New(filepath.Base(templatePath)).Funcs(funcMap)
	tpl, err := tpl.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	fName, _ := strings.CutSuffix(filepath.Base(templatePath), ".tpl")
	w, err := os.OpenFile(filepath.Join(outputDir, fName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	return tpl.ExecuteTemplate(w, filepath.Base(templatePath), data)
}