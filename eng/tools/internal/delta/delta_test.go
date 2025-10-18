package delta

import (
"testing"
)

func TestHasParamOrderChange(t *testing.T) {
tests := []struct {
name        string
lhsParams   string
rhsParams   string
lhsNames    string
rhsNames    string
expected    bool
}{
{
name:      "Simple order change",
lhsParams: "string, string, *Options",
rhsParams: "string, string, *Options",
lhsNames:  "resourceGroupName, serviceName, options",
rhsNames:  "serviceName, resourceGroupName, options",
expected:  true,
},
{
name:      "Three params order change",
lhsParams: "string, string, string",
rhsParams: "string, string, string",
lhsNames:  "resourceGroupName, serviceName, subscriptionID",
rhsNames:  "serviceName, subscriptionID, resourceGroupName",
expected:  true,
},
{
name:      "No change",
lhsParams: "context.Context, string, int",
rhsParams: "context.Context, string, int",
lhsNames:  "ctx, name, value",
rhsNames:  "ctx, name, value",
expected:  false,
},
{
name:      "Name change only",
lhsParams: "string, string",
rhsParams: "string, string",
lhsNames:  "oldName, newName",
rhsNames:  "firstName, lastName",
expected:  false,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := hasParamOrderChange(&tt.lhsParams, &tt.rhsParams, &tt.lhsNames, &tt.rhsNames)
if result != tt.expected {
t.Errorf("hasParamOrderChange() = %v, expected %v", result, tt.expected)
t.Logf("lhsParams: %s", tt.lhsParams)
t.Logf("rhsParams: %s", tt.rhsParams)
t.Logf("lhsNames: %s", tt.lhsNames)
t.Logf("rhsNames: %s", tt.rhsNames)
}
})
}
}
