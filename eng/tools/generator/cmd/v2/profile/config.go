// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package profile

import (
	"encoding/json"
	"io"
	"os"
)

type Definition struct {
	ProfileName string           `json:"profileName"`
	Modules     []ModuleProperty `json:"modules"`
}

type ModuleProperty struct {
	RP               string `json:"rp"`
	SpecName         string `json:"specName,omitempty"`
	Namespace        string `json:"namespace,omitempty"`
	Tag              string `json:"tag"`
	AdditionalConfig string `json:"additionalConfig"`
}

func (mp *ModuleProperty) UnmarshalJSON(b []byte) error {
	type TmpModuleProperty ModuleProperty
	var tmpJson TmpModuleProperty
	err := json.Unmarshal(b, &tmpJson)
	if err != nil {
		return err
	}
	*mp = ModuleProperty(tmpJson)
	if mp.SpecName == "" {
		mp.SpecName = mp.RP
	}
	if mp.Namespace == "" {
		mp.Namespace = "arm" + mp.RP
	}
	return nil
}

func ReadConfig(inputPath string) (def Definition, err error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return Definition{}, err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()
	b, err := io.ReadAll(file)
	if err != nil {
		return Definition{}, err
	}
	var result Definition
	if err := json.Unmarshal(b, &result); err != nil {
		return Definition{}, err
	}
	return result, nil
}
