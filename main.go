package main

import (
	"fmt"
	"reflect"
)

type abs struct {
	name string
	cnt  int
}

func main() {
	slice := []*abs{
		{name: "test1", cnt: 10},
		{name: "test2", cnt: 20},
	}
	target := &abs{name: "test1", cnt: 10}
	fmt.Println(Contains(slice, target))
}

// Contains
func Contains(slice interface{}, target interface{}) bool {
	if slice == nil {
		return false
	}

	kind := reflect.TypeOf(slice).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return false
	}

	v := reflect.ValueOf(slice)
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), target) {
			return true
		}
	}
	return false
}
