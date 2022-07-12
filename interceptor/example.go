package interceptor

import (
	"context"
	"fmt"
)

type invoker func(ctx context.Context, interceptors []interceptor, h handler) error
type handler func(ctx context.Context)
type interceptor func(ctx context.Context, h handler, ivk invoker) error

func getInvoker(ctx context.Context, interceptors []interceptor, curr int, ivk invoker) invoker {
	if curr == len(interceptors)-1 {
		return ivk
	}
	return func(ctx context.Context, interceptors []interceptor, h handler) error {
		return interceptors[curr+1](ctx, h, getInvoker(ctx, interceptors, curr+1, ivk))
	}
}

func getChainInterceptor(ctx context.Context, interceptors []interceptor, ivk invoker) interceptor {
	if len(interceptors) == 0 {
		return nil
	} else if len(interceptors) == 1 {
		return interceptors[0]
	} else {
		return func(ctx context.Context, h handler, ivk invoker) error {
			return interceptors[0](ctx, h, getInvoker(ctx, interceptors, 0, ivk))
		}
	}
}

func use() {
	var ctx context.Context
	var ceps []interceptor
	var h = func(ctx context.Context) {
		fmt.Println("some logic before ...")
	}

	var interceptor1 = func(ctx context.Context, h handler, ivk invoker) error {
		h(ctx)
		return ivk(ctx, ceps, h)
	}
	var interceptor2 = func(ctx context.Context, h handler, ivk invoker) error {
		h(ctx)
		return ivk(ctx, ceps, h)
	}
	ceps = append(ceps, interceptor1, interceptor2)

	var ivk = func(ctx context.Context, interceptors []interceptor, h handler) error {
		fmt.Println("invoker start")
		return nil
	}

	cep := getChainInterceptor(ctx, ceps, ivk)
	cep(ctx, h, ivk)
}
