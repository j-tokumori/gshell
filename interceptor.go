package gshell

import (
	"context"
)

type CallingInvoker func(ctx context.Context, c *Client, r RPC) (*Response, error)

type CallingInterceptor func(ctx context.Context, c *Client, r RPC, invoker CallingInvoker) (*Response, error)

func CallingInterceptorOption(i CallingInterceptor) Option {
	return newFuncOption(func(o *options) {
		if o.callingInt != nil {
			panic("The calling interceptor was already set and may not be reset.")
		}
		o.callingInt = i
	})
}

func ChainCalling(interceptors ...CallingInterceptor) CallingInterceptor {
	n := len(interceptors)
	return func(ctx context.Context, c *Client, r RPC, invoker CallingInvoker) (*Response, error) {
		chainer := func(currentInter CallingInterceptor, currentHandler CallingInvoker) CallingInvoker {
			return func(currentCtx context.Context, currentClient *Client, currentRPC RPC) (*Response, error) {
				return currentInter(currentCtx, currentClient, currentRPC, currentHandler)
			}
		}

		chainedHandler := invoker
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, c, r)
	}
}

func WithCallingChain(interceptors ...CallingInterceptor) Option {
	return CallingInterceptorOption(ChainCalling(interceptors...))
}
