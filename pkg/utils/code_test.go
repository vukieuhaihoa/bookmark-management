package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeGenerator_GenerateCode(t *testing.T) {
	testCodeGen := NewCodeGenerator()
	res, err := testCodeGen.GenerateCode(8)
	assert.Nil(t, err)
	assert.Len(t, res, 8)
	for _, ch := range res {
		assert.Contains(t, charset, string(ch))
	}
}
