package typespec

import (
	"errors"
	"regexp"

	"github.com/goccy/go-yaml"
)

/*
https://github.com/Azure/autorest.go/blob/main/packages/typespec-go/src/lib.ts#GoEmitterOptions
@azure-tools/typespec-go option
*/
type GoEmitterOptions struct {
	AzcoreVersion           string `yaml:"azcore-version,omitempty"`
	DisallowUnknownFields   bool   `yaml:"disallow-unknown-fields,omitempty"`
	FilePrefix              string `yaml:"file-prefix,omitempty"`
	GenerateFake            bool   `yaml:"generate-fakes,omitempty"`
	InjectSpanc             bool   `yaml:"inject-spans,omitempty"`
	Module                  string `yaml:"module,omitempty"`
	ModuleVersion           string `yaml:"module-version,omitempty"`
	RawJsonAsBytes          bool   `yaml:"rawjson-as-bytes,omitempty"`
	SliceElementsByVal      bool   `yaml:"slice-elements-byval,omitempty"`
	SingleClient            bool   `yaml:"single-client,omitempty"`
	Stutter                 string `yaml:"stutter,omitempty"`
	FixConstStuttering      bool   `yaml:"fix-const-stuttering,omitempty"`
	RemoveUnreferencedTypes bool   `yaml:"remove-unreferenced-types,omitempty"`
}

func NewGoEmitterOptions(emitOption any) (*GoEmitterOptions, error) {
	option := GoEmitterOptions{}

	data, err := yaml.Marshal(emitOption)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &option); err != nil {
		return nil, err
	}

	return &option, err
}

const moduleRegex = `^github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/(\w+)/arm(\w+)$`

var (
	ErrModuleEmpty  = errors.New("typesepec-go option `module` is required")
	ErrModuleFormat = errors.New("module must be in the format of github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/armxxx")
)

func (o *GoEmitterOptions) Validate() error {
	if o.Module == "" {
		return ErrModuleEmpty
	}

	matched := regexp.MustCompile(moduleRegex).MatchString(o.Module)
	if !matched {
		return ErrModuleFormat
	}

	return nil
}
