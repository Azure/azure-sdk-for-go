package model

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/stretchr/testify/assert"
)

func TestSortChangeItem(t *testing.T) {
	s := map[string]delta.Signature{
		"*X":{},
		"*NewD":{},
		"C":{},
		"NewA":{},
		"B":{},
		"*NewH":{},
		"D.Get":{},
		"*D.Create":{},
	}

	sortResult := sortChangeItem(s)

	expcted := []string{
		"NewA", 
		"B", 
		"C", 
		"*NewD", 
		"*D.Create", 
		"D.Get", 
		"*NewH", 
		"*X",
	}
	assert.Equal(t,expcted,sortResult)
}