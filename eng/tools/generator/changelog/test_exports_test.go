package changelog

import (
"encoding/json"
"fmt"
"testing"
"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

func TestExportStructure(t *testing.T) {
oldExport, err := exports.Get("./testdata/old/parameter")
if err != nil {
t.Fatal(err)
}

// Print out the functions
for name, fn := range oldExport.Funcs {
data, _ := json.MarshalIndent(fn, "", "  ")
fmt.Printf("Function: %s\n%s\n\n", name, string(data))
}
}
