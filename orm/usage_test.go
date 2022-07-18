package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FastUsage(t *testing.T) {
	err := FastUsage()
	assert.Nil(t, err)
}

func Test_UnionSearch(t *testing.T) {
	err := UnionSearch()
	assert.Nil(t, err)
}

func Test_SQLSearch(t *testing.T) {
	err := SQLSearch()
	assert.Nil(t, err)
}

func Test_TransOperation(t *testing.T) {
	err := TransOperation()
	assert.Nil(t, err)
}
