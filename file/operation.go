package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"syscall"
)

func fileOne() {
	content, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
}

func fileTwo() {
	content, err := ioutil.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
}

func fileThree() {
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	fmt.Println(string(content))
}

func fileFour() {
	file, err := os.OpenFile("test.txt", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	fmt.Println(string(content))
}

func fileFive() {
	// 创建文件句柄
	fi, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 创建reader
	r := bufio.NewReader(fi)

	for {
		lineBytes, err := r.ReadBytes('\n')
		line := strings.TrimSpace(string(lineBytes))
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(line)
	}
}

func fileSix() {
	// 创建文件句柄
	fi, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 创建reader
	r := bufio.NewReader(fi)

	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println(line)
	}
}

func fileSeven() {
	// 创建文件句柄
	fi, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 创建reader
	r := bufio.NewReader(fi)

	// 每次读取1024个字节
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		fmt.Println(string(buf[:n]))
	}
}

func fileEight() {
	fd, err := syscall.Open("test.txt", syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Println("Failed on open: ", err)
	}
	defer syscall.Close(fd)

	var wg sync.WaitGroup
	wg.Add(2)
	dataChan := make(chan []byte)
	go func() {
		wg.Done()
		for {
			data := make([]byte, 100)
			n, _ := syscall.Read(fd, data)
			if n == 0 {
				break
			}
			dataChan <- data
		}
		close(dataChan)
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}

				fmt.Println(string(data))
			default:
			}
		}
	}()

	wg.Wait()
}
