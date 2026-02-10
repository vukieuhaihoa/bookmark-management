package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeGenerator_GenerateCode(t *testing.T) {
	testCodeGen := NewCodeGenerator()
	res, err := testCodeGen.GenerateCode("", 8)
	assert.Nil(t, err)
	assert.Len(t, res, 8)
	for _, ch := range res {
		assert.Contains(t, charset, string(ch))
	}
}

func TestCodeGenerator_GenerateCode_WithPrefix(t *testing.T) {
	testCodeGen := NewCodeGenerator()
	prefix := "test"
	res, err := testCodeGen.GenerateCode(prefix, 8)
	assert.Nil(t, err)
	expectedLength := len(prefix) + 1 + 8 // prefix + '_' + code length
	assert.Len(t, res, expectedLength)
	assert.Contains(t, res, prefix+"_")
	for _, ch := range res[len(prefix)+1:] {
		assert.Contains(t, charset, string(ch))
	}
}
