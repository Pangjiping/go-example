package reflect

import (
	"fmt"
	"reflect"
)

// 将interface{}转为已知结构体

// 将interface{}转为未知结构体 [可能有多种类型，但是事先不知道要转成哪一个]
// 比如下面这两个结构体，因为他们的字段是一样的，我需要用一个map给他们赋值
// 同时以interface{}的形式传入一个初始化了的struct，现在要区分我传入的是哪个struct
type EvictionHard struct {
	MemoryAvailable   int
	NodefsAvailable   int
	NodefsInodesFree  int
	ImagefsAvailable  int
	ImagefsInodesFree int
	PidAvailable      int
}

type EvictionSoft struct {
	MemoryAvailable   int
	NodefsAvailable   int
	NodefsInodesFree  int
	ImagefsAvailable  int
	ImagefsInodesFree int
	PidAvailable      int
}

// 注意map要全量填充 [不设置的填充成默认值即可]
func setEvictionConfig(m map[string]interface{}, config interface{}) {
	slice := map2Slice(m)

	val := reflect.ValueOf(config).Elem()
	for i := 0; i < len(slice); i++ {
		val.Field(i).SetInt(int64(slice[i].(int)))
	}
	fmt.Println(config)
}

func map2Slice(m map[string]interface{}) []interface{} {
	res := make([]interface{}, 0)
	for _, v := range m {
		res = append(res, v)
	}
	return res
}
