package data_struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SkipList(t *testing.T) {
	skipList := NewSkipListInt()

	skipList.Set(1, "test-1")
	res := skipList.Get(1)
	assert.NotNil(t, res)
	assert.Equal(t, "test-1", res.(string))
	assert.Equal(t, int32(1), skipList.Len())

	skipList.Set(1, "test-x")
	res = skipList.Get(1)
	assert.NotNil(t, res)
	assert.Equal(t, "test-x", res.(string))
	assert.Equal(t, int32(1), skipList.Len())

	skipList.Delete(1)
	res = skipList.Get(1)
	assert.Nil(t, res)
	assert.Equal(t, int32(0), skipList.Len())
}
