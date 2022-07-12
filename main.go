package main

import "fmt"

func main() {
	fmt.Println("Hello Golang!")
	defer fmt.Println("1")
	defer fmt.Println("2")
}
