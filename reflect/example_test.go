package reflect

import (
	"fmt"
	"testing"
)

func Test_Interface2Struct(t *testing.T) {
	m := map[string]interface{}{
		"MemoryAvailable":   10,
		"NodefsAvailable":   10,
		"NodefsInodesFree":  10,
		"ImagefsAvailable":  10,
		"ImagefsInodesFree": 10,
		"PidAvailable":      10,
	}
	config := &EvictionHard{
		MemoryAvailable: 1,
	}
	setEvictionConfig(m, config)
}

func Test_Map2Slice(t *testing.T) {
	m := map[string]interface{}{
		"MemoryAvailable":   10,
		"NodefsAvailable":   10,
		"NodefsInodesFree":  10,
		"ImagefsAvailable":  10,
		"ImagefsInodesFree": 10,
		"PidAvailable":      10,
	}
	slice := map2Slice(m)
	fmt.Println(slice)
}
