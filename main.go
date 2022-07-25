package main

import "fmt"

func main() {
	fmt.Println("Hello Golang!")

	s := []string{}
	fmt.Println(len(s))
	fmt.Println(s == nil)

	fmt.Println(interface{}(nil).([]string))
}
