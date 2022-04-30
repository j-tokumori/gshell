package gshell

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type RPC interface {
	Call(ctx context.Context, conn grpc.ClientConnInterface) (proto.Message, *metadata.MD, *metadata.MD)
	Key() string
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

func NewClient(host string, secure bool, ctxFunc ContextFunc, errFunc ErrorFunc) *Client {
	opt := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	if !secure {
		opt = grpc.WithInsecure()
	}
	//conn, err := grpc.Dial(host, opt, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	conn, err := grpc.Dial(host, opt)
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
	c.Replies[r.Key()] = rep
	c.Headers[r.Key()] = h
	c.Trailers[r.Key()] = t
}

func (c *Client) Call(r RPC) {
	Defaultize(r, c)
	c.LastRPCName = reflect.ValueOf(r).Elem().Type().Name()
	c.CallWithRecover(r)
}

func (c *Client) CallByJSON(rpcName string, in []byte) {
	f := c.rpcMap[rpcName]
	rpc := f(in)
	c.Call(rpc)
}

func (c *Client) PrintLastResponse() {
	c.PrintResponse(c.LastRPCName)
}

func (c *Client) PrintResponse(rpcName string) {
	fmt.Println("Header:")
	c.PrintMD(c.Headers[rpcName])
	fmt.Println("Trailer:")
	c.PrintMD(c.Trailers[rpcName])
	fmt.Println("Reply:")
	c.PrintReply(rpcName)
}

func (c *Client) PrintLastReply() {
	c.PrintReply(c.LastRPCName)
}

func (c *Client) PrintReply(rpcName string) {
	m := c.Replies[rpcName]
	o := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	j, err := o.Marshal(m)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, j, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

func (c *Client) PrintLastHeader() {
	c.PrintHeader(c.LastRPCName)
}

func (c *Client) PrintHeader(rpcName string) {
	c.PrintMD(c.Headers[rpcName])
}

func (c *Client) PrintLastTrailer() {
	c.PrintTrailer(c.LastRPCName)
}

func (c *Client) PrintTrailer(rpcName string) {
	c.PrintMD(c.Trailers[rpcName])
}

func (c *Client) PrintMD(md *metadata.MD) {
	for key, values := range *md {
		fmt.Printf("  %s: %v\n", key, values)
	}
}

func (c *Client) PrintLastTraceURL() {
	c.PrintTraceURL(c.LastRPCName)
}

func (c *Client) PrintTraceURL(rpcName string) {
	md := c.Headers[rpcName]
	s := md.Get("x-cloud-trace-context")
	if len(s) <= 0 {
		fmt.Println("no trace header.")
		return
	}
	s1 := strings.Split(s[0], ";")
	if len(s1) != 2 || s1[1] != "o=1" {
		fmt.Println("no trace output.")
		return
	}
	s2 := strings.Split(s1[0], "/")
	if len(s2) != 2 {
		fmt.Println("no trace id.")
		return
	}
	fmt.Printf("https://console.cloud.google.com/traces/list?tid=%s\n", s2[0])
}

func (c *Client) PrintSample(rpcName string) {
	f := c.rpcMap[rpcName]
	a := f([]byte("{}"))

	Defaultize(a, c)
	Samplize(a)

	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println("rpc " + rpcName + " " + string(j))
}

func (c *Client) PrintList(search string) {
	nameList := make([]string, 0)
	for key := range c.rpcMap {
		nameList = append(nameList, key)
	}
	sort.Strings(nameList)
	for _, s := range nameList {
		if strings.Contains(s, search) {
			fmt.Println(s)
		}
	}
}

// Defaultize 引数 r をデフォルト値で埋める。破壊的メソッド
// TODO: generics で書き直し
func Defaultize(r interface{}, c *Client) {
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

// Samplize 引数 r にサンプル値を代入
// TODO: generics で書き直し
func Samplize(r interface{}) {
	pv := reflect.ValueOf(r)
	if pv.Kind() != reflect.Ptr {
		panic("r is need pointer.")
	}

	numField := pv.Elem().NumField()
	for i := 0; i < numField; i++ {
		field := pv.Elem().Type().Field(i)
		if !field.IsExported() {
			continue
		}

		// ゼロ値チェック
		value := pv.Elem().Field(i)
		var isZero bool
		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.Slice:
			isZero = value.Len() == 0
		default:
			isZero = value.Interface() == reflect.Zero(field.Type).Interface()
		}
		if !isZero {
			return
		}

		// ゼロ値なら代入
		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			pv.Elem().FieldByName(field.Name).SetInt(1)
		case reflect.String:
			pv.Elem().FieldByName(field.Name).SetString("hoge")
		case reflect.Struct:
			panic("Struct type is not support.")
		case reflect.Ptr:
			itf := reflect.New(field.Type.Elem()).Interface()
			Samplize(itf)
			pv.Elem().FieldByName(field.Name).Set(reflect.ValueOf(itf))
		case reflect.Slice:
			f := pv.Elem().FieldByName(field.Name)
			// TODO: prt ではなかった場合の対応
			if field.Type.Elem().Kind() == reflect.Ptr {
				itf := reflect.New(field.Type.Elem().Elem()).Interface()
				Samplize(itf)
				f.Set(reflect.Append(f, reflect.ValueOf(itf)))
			} else {
				switch field.Type.Elem().Kind() { //nolint:exhaustive
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					f.Set(reflect.Append(f, reflect.ValueOf(int32(1)))) // TODO: 1 ではだめだった。。int32 以外どうするか
				case reflect.Struct:
					panic("TODO")
				case reflect.String:
					f.Set(reflect.Append(f, reflect.ValueOf("age")))
				}
			}
		}
	}
}
