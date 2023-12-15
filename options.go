package gshell

import (
	"context"
	"crypto/tls"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

type options struct {
	callingInt      CallingInterceptor
	grpcDialOptions []grpc.DialOption
	rpcMap          map[string]RPCFactory
	rpcAliasHandler RPCAliasHandler
	rpcAliasMap     map[string]string
	scenarioFactory ScenarioFactory
}

func (o *options) init(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}

	if o.rpcAliasHandler == nil {
		o.rpcAliasHandler = defaultRPCAliasHandler
	}

	o.rpcAliasMap = make(map[string]string, len(o.rpcMap))
	for key, _ := range o.rpcMap {
		s := strings.Split(key, ".") // s[0] = serviceName, s[1] = methodName
		o.rpcAliasMap[o.rpcAliasHandler(s[0], s[1])] = key
	}
}

func (o *options) getRPCTypeName(arg string) string {
	if key, ok := o.rpcAliasMap[arg]; ok {
		return key
	}
	return arg
}

type RPCFactory func([]byte) RPC

type RPCAliasHandler func(serviceName, methodName string) string

func defaultRPCAliasHandler(serviceName, methodName string) string {
	return serviceName + "." + methodName
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

func WithForceIPv6() Option {
	var d net.Dialer
	d.Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.DialContext(ctx, "udp6", addr)
		},
	}

	return newFuncOption(func(o *options) {
		o.grpcDialOptions = append(o.grpcDialOptions, grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return d.DialContext(ctx, "tcp6", addr)
		}))
	})
}

func RegisterRPCFactories(funcList ...RPCFactory) Option {
	return newFuncOption(func(o *options) {
		if o.rpcMap == nil {
			o.rpcMap = make(map[string]RPCFactory, 0)
		}
		for _, f := range funcList {
			o.rpcMap[getRPCTypeName(f([]byte("{}")))] = f
		}
	})
}

func RegisterScenarioFactory(f ScenarioFactory) Option {
	return newFuncOption(func(o *options) {
		o.scenarioFactory = f
	})
}

func RegisterRPCAliasHandler(f RPCAliasHandler) Option {
	return newFuncOption(func(o *options) {
		o.rpcAliasHandler = f
	})
}
