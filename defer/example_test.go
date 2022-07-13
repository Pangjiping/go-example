package _defer

import "testing"

func Test_DeferExecSequence(t *testing.T) {
	ExecSequence()
}

func Test_DeferExecWithReturn(t *testing.T) {
	ExecWithReturn()
}

func Test_ExecInit(t *testing.T) {
	ExecInit()
}

func Test_ExecReturnButDefer(t *testing.T) {
	ExecReturnButDefer()
}

func Test_ExecWithoutRecover(t *testing.T) {
	ExecWithoutRecover()
}

func Test_ExecWithRecover(t *testing.T) {
	ExecWithRecover()
}

func Test_ExecDeferPanic(t *testing.T) {
	ExecDeferPanic()
}

func Test_ExecWithSubFunc(t *testing.T) {
	ExecWithSubFunc()
}
