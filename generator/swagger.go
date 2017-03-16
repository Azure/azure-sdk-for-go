package main

// Swagger represents top-level swagger contents needed for generator serialization.
type Swagger struct {
	SwaggerVersion string    `json:"swagger"`
	Info           AzureInfo `json:"info"`
	Documents      []string  `json:"documents"`
}

// AzureInfo holds the contents of the metadata surrounding
type AzureInfo struct {
	Title    string                 `json:"title"`
	Version  string                 `json:"version"`
	Settings CodeGenerationSettings `json:"x-ms-code-generation-settings"`
}

// CodeGenerationSettings handles the Microsoft Swagger extension "x-ms-code-generation-settings"
type CodeGenerationSettings struct {
	Name string `json:"name"`
}
