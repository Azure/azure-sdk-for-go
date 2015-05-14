package vmutils

import (
	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

func updateOrAddConfig(configs []ConfigurationSet, configType ConfigurationSetType, update func(*ConfigurationSet)) []ConfigurationSet {
	config := findConfig(configs, configType)
	if config == nil {
		configs = append(configs, ConfigurationSet{ConfigurationSetType: configType})
		config = findConfig(configs, configType)
	}
	update(config)

	return configs
}

func findConfig(configs []ConfigurationSet, configType ConfigurationSetType) *ConfigurationSet {
	for i, config := range configs {
		if config.ConfigurationSetType == configType {
			// need to return a pointer to the original set in configs,
			// not the copy made by the range iterator
			return &configs[i]
		}
	}

	return nil
}
