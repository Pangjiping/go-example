package informer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Informer(t *testing.T) {
	err := ExecInformer()
	assert.Nil(t, err)
}
