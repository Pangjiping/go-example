package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JsonUsage(t *testing.T) {
	testcases := []map[string]interface{}{
		{
			"description": "int input",
			"input":       `{"name":"tq","value":10}`,
			"resNil":      false,
			"wantErr":     false,
		},
		{
			"description": "double input",
			"input":       `{"name":"tq","value":11.111}`,
			"resNil":      false,
			"wantErr":     false,
		},
	}

	for _, tc := range testcases {
		t.Log(tc["description"])
		res, err := Usage(tc["input"].(string))
		assert.Equal(t, tc["resNil"], res == nil)
		assert.Equal(t, tc["wantErr"], err != nil)
		t.Logf("input: %s, and unmarshal result: %v", tc["input"], *res)
	}
}
