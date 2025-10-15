package changelog

import (
"encoding/json"
"fmt"
"testing"
"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

func TestASTParamNames(t *testing.T) {
pkg, err := exports.LoadPackage("./testdata/old/parameter")
if err != nil {
t.Fatal(err)
}

c := pkg.GetExports()

// Print out the functions
for name, fn := range c.Funcs {
data, _ := json.MarshalIndent(fn, "", "  ")
fmt.Printf("Function: %s\n%s\n\n", name, string(data))
}
}
