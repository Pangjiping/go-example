package _defer

import "fmt"

/*
	defer的执行顺序
	在一个函数中，写在前面的defer会比写在后面的defer调用晚
*/
func ExecSequence() {
	defer fmt.Println("A")
	defer fmt.Println("B")
	defer fmt.Println("C")
}

/*
	return之后的语句先执行
	defer后的语句后执行
*/
func ExecWithReturn() int {
	defer deferFunc()
	return returnFunc()
}

func deferFunc() int {
	fmt.Println("defer func called")
	return 0
}

func returnFunc() int {
	fmt.Println("return func called")
	return 0
}

/*
	函数的返回值初始化
	func DeferFunc1(i int) (t int) {}，其中返回值t int，这个t会在函数起始处被初始化为对应类型的零值并且作用域为整个函数
*/
func ExecInit() {
	DeferFunc1(10)
}

func DeferFunc1(i int) (t int) {
	fmt.Println("t = ", t)
	return 2
}

/*
	有名函数返回值遇见defer情况
	在没有defer的情况下，其实函数的返回就是与return一致的，但是有了defer就不一样了
	我们通过知识点2得知，先return，再defer，所以在执行完return之后，还要再执行defer里的语句，依然可以修改本应该返回的结果

	该returnButDefer()本应的返回值是1，但是在return之后，又被defer的匿名func函数执行，所以t=t*10被执行
	最后returnButDefer()返回给上层main()的结果为10
*/
func ExecReturnButDefer() {
	fmt.Println(returnButDefer())
}

func returnButDefer() (t int) { // t初始化0， 并且作用域为该函数全域
	defer func() {
		t = t * 10
	}()

	return 1
}

/*
	defer遇见panic
	我们知道，能够触发defer的是遇见return(或函数体到末尾)和遇见panic
	那么，遇到panic时，遍历本协程的defer链表，并执行defer。在执行defer过程中:遇到recover则停止panic，返回recover处继续往下执行
	如果没有遇到recover，遍历完本协程的defer链表后，向stderr抛出panic信息

	defer 最大的功能是 panic 后依然有效 所以defer可以保证你的一些资源一定会被关闭，从而避免一些异常出现的问题
*/
func ExecWithoutRecover() {
	defer_call1()
	fmt.Println("main 正常结束")
}

func defer_call1() {
	defer func() { fmt.Println("defer: panic 之前1") }()
	defer func() { fmt.Println("defer: panic 之前2") }()
	panic("异常内容") // 触发defer出栈

	defer func() { fmt.Println("defer: panic 之后，不会执行") }()
}

func ExecWithRecover() {
	defer_call2()
	fmt.Println("main 正常结束")
}
func defer_call2() {
	defer func() {
		fmt.Println("defer: panic 之前1，捕获异常")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	defer func() { fmt.Println("defer: panic 之前2，不捕获") }()
	panic("异常内容") // 触发defer出栈

	defer func() { fmt.Println("defer: panic 之后，不会执行") }()
}

/*
	defer中包含panic
	panic仅有最后一个可以被recover捕获
	触发panic("panic")后defer顺序出栈执行，第一个被执行的defer中 会有panic("defer panic")异常语句
	这个异常将会覆盖掉main中的异常panic("panic")，最后这个异常被第二个执行的defer捕获到
*/
func ExecDeferPanic() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("fatal")
		}
	}()

	defer func() {
		panic("defer panic")
	}()

	panic("main panic")
}

/*
	defer下的函数参数包含子函数
*/

func function(index int, value int) int {
	fmt.Println(index)
	return index
}

func ExecWithSubFunc() {
	defer function(1, function(3, 0))
	defer function(2, function(4, 0))
}
