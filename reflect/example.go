package reflect

import (
	"fmt"
	"reflect"
	"unsafe"
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

	rv := reflect.ValueOf(config)
	elem := rv.Elem()
	for i := 0; i < rv.NumField(); i++ {
		elem.Field(i).SetInt(int64(slice[i].(int)))
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

// String2Bytes convert string to bytes.
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String convert bytes to string.
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// reflect获取tag
type J struct {
	a string // 小写无tag
	b string `json:"B"` // 小写加tag
	C string // 大写无tag
	D string `json:"DD" otherTag:"good"` // 大写加tag
}

// printTag方法传入的是j的指针。
// reflect.TypeOf(stru).Elem()获取指针指向的值对应的结构体内容。
// NumField()可以获得该结构体的含有几个字段。
// 遍历结构体内的字段，通过t.Field(i).Tag.Get("json")可以获取到tag为json的字段。
// 如果结构体的字段有多个tag，比如叫otherTag,同样可以通过t.Field(i).Tag.Get("otherTag")获得

// json包不能导出私有变量的tag是因为取不到反射信息的说法，但是直接取t.Field(i).Tag.Get("json")却可以获取到私有变量的json字段，是为什么呢？
// 其实准确的说法是，json包里不能导出私有变量的tag是因为json包里认为私有变量为不可导出的Unexported
// 所以跳过获取名为json的tag的内容。具体可以看/src/encoding/json/encode.go:1070的代码。

func printTag(stru interface{}) {
	t := reflect.TypeOf(stru).Elem()
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("结构体内第%v个字段 %v 对应的json tag是 %v , 还有otherTag？ = %v \n", i+1, t.Field(i).Name, t.Field(i).Tag.Get("json"), t.Field(i).Tag.Get("otherTag"))
	}
}

/*
func typeFields(t reflect.Type) []field {
    // 注释掉其他逻辑...
    // 遍历结构体内的每个字段
    for i := 0; i < f.typ.NumField(); i++ {
        sf := f.typ.Field(i)
        isUnexported := sf.PkgPath != ""
        // 注释掉其他逻辑...
        if isUnexported {
            // 如果是不可导出的变量则跳过
            continue
        }
        // 如果是可导出的变量（public），则获取其json字段
        tag := sf.Tag.Get("json")
        // 注释掉其他逻辑...
    }
    // 注释掉其他逻辑...
}
*/
