package gshell

import (
	"context"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Conn = grpc.ClientConnInterface

type RPC interface {
	Call(ctx context.Context, conn Conn) (*Response, error)
}

type RPCFactory func([]byte) RPC

type ErrorHandler func(error)

type Response struct {
	Reply   proto.Message
	Header  *metadata.MD
	Trailer *metadata.MD
}

func NewEmptyResponse() *Response {
	return &Response{Header: &metadata.MD{}, Trailer: &metadata.MD{}}
}

type Client struct {
	conn        grpc.ClientConnInterface
	opts        *options
	responses   map[string]*Response
	lastRPCName string
}

func NewClient(host string, opts *options) *Client {
	conn, err := grpc.Dial(host, opts.grpcDialOptions...)
	if err != nil {
		panic(err)
	}
	return &Client{
		conn:      conn,
		opts:      opts,
		responses: make(map[string]*Response),
	}
}

func (c *Client) Response(name string) *Response {
	return c.responses[name]
}

func (c *Client) LastRPCName() string {
	return c.lastRPCName
}

func (c *Client) Call(r RPC) {
	defaultize(r, c)
	c.lastRPCName = reflect.ValueOf(r).Elem().Type().Name()

	res, err := c.Invoke(context.Background(), r)
	if err != nil {
		if c.opts.errorHandler != nil {
			c.opts.errorHandler(err)
		}
		panic(err)
	}

	c.responses[getKey(r)] = res
}

func (c *Client) Invoke(ctx context.Context, r RPC) (*Response, error) {
	if c.opts.callingInt != nil {
		return c.opts.callingInt(ctx, c, r, invoke)
	}
	return invoke(ctx, c, r)
}

func invoke(ctx context.Context, c *Client, r RPC) (*Response, error) {
	return r.Call(ctx, c.conn)
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

func ResponseOptions(res *Response) []grpc.CallOption {
	return []grpc.CallOption{
		grpc.Header(res.Header),
		grpc.Trailer(res.Trailer),
	}
}
