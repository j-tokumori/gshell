package gshell

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type options struct {
	callingInt      CallingInterceptor
	grpcDialOptions []grpc.DialOption
	rpcMap          map[string]RPCFactory
	scenarioFactory ScenarioFactory
	errorHandler    ErrorHandler
}

type Option interface {
	apply(*options)
}

type EmptyOption struct{}

func (EmptyOption) apply(*options) {}

type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithSecure() Option {
	return newFuncOption(func(o *options) {
		o.grpcDialOptions = append(o.grpcDialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	})
}

func WithInsecure() Option {
	return newFuncOption(func(o *options) {
		o.grpcDialOptions = append(o.grpcDialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
}

type Codec = encoding.Codec

func WithCodec(codec Codec) Option {
	return newFuncOption(func(o *options) {
		o.grpcDialOptions = append(o.grpcDialOptions, grpc.WithDefaultCallOptions(grpc.ForceCodec(codec)))
	})
}

func RegisterRPCFactories(funcList ...RPCFactory) Option {
	return newFuncOption(func(o *options) {
		if o.rpcMap == nil {
			o.rpcMap = make(map[string]RPCFactory, 0)
		}
		for _, f := range funcList {
			o.rpcMap[getKey(f([]byte("{}")))] = f
		}
	})
}

func RegisterScenarioFactory(f ScenarioFactory) Option {
	return newFuncOption(func(o *options) {
		o.scenarioFactory = f
	})
}

func RegisterErrorHandler(f ErrorHandler) Option {
	return newFuncOption(func(o *options) {
		o.errorHandler = f
	})
}
