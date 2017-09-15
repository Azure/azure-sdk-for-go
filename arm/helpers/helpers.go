package helpers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// BuildDeployment is a helper method that create a deployment object
func BuildDeployment(mode resources.DeploymentMode, templateFile string, parameters map[string]interface{}) (deployment resources.Deployment, err error) {
	fileMap, err := parseJSONFromFile(templateFile)
	if err != nil {
		return
	}
	template := (*fileMap)["template"].(map[string]interface{})

	finalParameters := map[string]interface{}{}
	for k, v := range parameters {
		addElementToMap(&finalParameters, k, v)
	}

	deployment.Properties = &resources.DeploymentProperties{
		Mode:       mode,
		Template:   &template,
		Parameters: &finalParameters,
	}
	return
}

func parseJSONFromFile(filePath string) (*map[string]interface{}, error) {
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileMap := map[string]interface{}{}
	if err = json.Unmarshal(text, &fileMap); err != nil {
		return nil, err
	}
	return &fileMap, err
}

func addElementToMap(parameter *map[string]interface{}, key string, value interface{}) {
	(*parameter)[key] = map[string]interface{}{
		"value": value,
	}
}
