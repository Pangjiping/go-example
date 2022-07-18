package gen_id

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Snowflake(t *testing.T) {
	err := GenBySnowflake()
	assert.Nil(t, err)
}

func Test_GenbySonyflake(t *testing.T) {
	err := GenbySonyflake()
	assert.Nil(t, err)
}
