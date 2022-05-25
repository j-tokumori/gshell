package gshell

import (
	"context"
	"crypto/tls"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type RPC interface {
	Call(ctx context.Context, conn grpc.ClientConnInterface) (proto.Message, *metadata.MD, *metadata.MD)
}

type NewRPCFunc func([]byte) RPC

type ContextFunc func(context.Context, *Client) context.Context

type ErrorFunc func(error)

type Client struct {
	Conn        grpc.ClientConnInterface
	ContextFunc ContextFunc
	ErrorFunc   ErrorFunc
	rpcMap      map[string]NewRPCFunc
	Replies     map[string]proto.Message
	Headers     map[string]*metadata.MD
	Trailers    map[string]*metadata.MD
	LastRPCName string
}

func NewClient(host string, secure bool, codec encoding.Codec, ctxFunc ContextFunc, errFunc ErrorFunc) *Client {
	opts := make([]grpc.DialOption, 0)
	if secure {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if codec != nil {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.ForceCodec(codec)))
	}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		panic(err)
	}
	return &Client{
		Conn:        conn,
		ContextFunc: ctxFunc,
		ErrorFunc:   errFunc,
		rpcMap:      make(map[string]NewRPCFunc),
		Replies:     make(map[string]proto.Message),
		Headers:     make(map[string]*metadata.MD),
		Trailers:    make(map[string]*metadata.MD),
	}
}

func (c *Client) CallWithRecover(r RPC) {
	defer func() {
		rec := recover()
		if err, ok := rec.(error); ok {
			c.ErrorFunc(err)
		}
		if rec != nil {
			panic(rec)
		}
	}()

	rep, h, t := r.Call(c.ContextFunc(context.Background(), c), c.Conn)

	key := getKey(r)

	c.Replies[key] = rep
	c.Headers[key] = h
	c.Trailers[key] = t
}

func (c *Client) Call(r RPC) {
	defaultize(r, c)
	c.LastRPCName = reflect.ValueOf(r).Elem().Type().Name()
	c.CallWithRecover(r)
}

func (c *Client) registerRPC(f NewRPCFunc) {
	c.rpcMap[getKey(f([]byte("{}")))] = f
}

func getKey(r RPC) string {
	return reflect.ValueOf(r).Elem().Type().Name()
}

// defaultize 引数 r をデフォルト値で埋める。破壊的メソッド
// TODO: generics で書き直し
func defaultize(r interface{}, c *Client) {
	pv := reflect.ValueOf(r)
	if pv.Kind() != reflect.Ptr {
		panic("r is need pointer.")
	}

	var dv *reflect.Value
	if m := pv.MethodByName("Default"); m.Kind() == reflect.Func {
		mt := m.Type()
		argv := make([]reflect.Value, mt.NumIn())
		argv[0] = reflect.ValueOf(c)
		result := m.Call(argv)
		dv = &result[0]
	}
	if dv == nil { // default 値がなければ終了
		return
	}

	numField := pv.Elem().NumField()
	for i := 0; i < numField; i++ {
		field := pv.Elem().Type().Field(i)
		if !field.IsExported() {
			continue
		}

		value := pv.Elem().Field(i)
		var isZero bool
		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.Slice:
			isZero = value.Len() == 0
		default:
			isZero = value.Interface() == reflect.Zero(field.Type).Interface()
		}
		if isZero {
			pv.Elem().FieldByName(field.Name).Set(dv.Elem().FieldByName(field.Name))
		}
	}
}
