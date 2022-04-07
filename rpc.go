package gshell

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type RPC interface {
	Call(ctx context.Context, e *Client)
}

type Config struct {
	Trace bool
}

type Client struct {
	Conn        grpc.ClientConnInterface
	Config      Config
	Replies     map[string]proto.Message
	Headers     map[string]*metadata.MD
	Trailers    map[string]*metadata.MD
	LastRPCName string
}

func NewClient(host string, secure bool, config Config) *Client {
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
		Conn:     conn,
		Config:   config,
		Headers:  make(map[string]*metadata.MD),
		Trailers: make(map[string]*metadata.MD),
	}
}

func (c *Client) CallWithRecover(r RPC) {
	defer func() {
		rec := recover()
		if err, ok := rec.(error); ok {
			PrintGrpcErr(err)
		}
		if rec != nil {
			panic(rec)
		}
	}()

	r.Call(context.Background(), c)
}

func (c *Client) Call(r RPC) {
	Defaultize(r, c)
	c.LastRPCName = reflect.ValueOf(r).Elem().Type().Name()
	c.CallWithRecover(r)
}

func (c *Client) CallWithJSON(rpcName string, in []byte) {
	v := reflect.ValueOf(c)
	m := v.MethodByName(rpcName)
	if m.Kind() != reflect.Func {
		panic(m.Kind())
	}
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn())
	argv[0] = reflect.ValueOf(in)

	result := m.Call(argv)
	if len(result) != 1 {
		panic("") // TODO
	}
	a := result[0].Interface().(RPC)
	c.Call(a)
}

func (c *Client) Header(rpcName string) *metadata.MD {
	if _, ok := c.Headers[rpcName]; !ok {
		c.Headers[rpcName] = &metadata.MD{}
	}
	return c.Headers[rpcName]
}

func (c *Client) Trailer(rpcName string) *metadata.MD {
	if _, ok := c.Trailers[rpcName]; !ok {
		c.Trailers[rpcName] = &metadata.MD{}
	}
	return c.Trailers[rpcName]
}

func (c *Client) PrintLastResponse() {
	c.PrintResponse(c.LastRPCName)
}

func (c *Client) PrintResponse(rpcName string) {
	println("Header:")
	c.PrintMD(c.Headers[rpcName])
	println("Trailer:")
	c.PrintMD(c.Trailers[rpcName])
	println("Reply:")
	c.PrintReply(rpcName)
}

func (c *Client) PrintLastReply() {
	c.PrintReply(c.LastRPCName)
}

func (c *Client) PrintReply(rpcName string) {
	//v := reflect.ValueOf(c.Replies)
	//p := v.FieldByName(rpcName)
	//if p.Kind() != reflect.Ptr {
	//	panic(p.Kind())
	//}
	//m := p.Interface().(proto.Message)
	//
	//j, err := protojson.Marshal(m)
	//if err != nil {
	//	panic(err)
	//}
	//var buf bytes.Buffer
	//err = json.Indent(&buf, j, "", "  ")
	//if err != nil {
	//	panic(err)
	//}
	//println(buf.String())
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
	// TODO: CallWithJSON と処理の共通化
	v := reflect.ValueOf(c)
	m := v.MethodByName(rpcName)
	if m.Kind() != reflect.Func {
		panic(m.Kind())
	}
	t := m.Type()
	argv := make([]reflect.Value, t.NumIn())
	argv[0] = reflect.ValueOf([]byte("{}"))

	result := m.Call(argv)
	if len(result) != 1 {
		panic("") // TODO
	}
	a := result[0].Interface().(RPC)
	Defaultize(a, c)
	Samplize(a)
	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	println("rpc " + rpcName + " " + string(j))
}

func (c *Client) PrintList(search string) {
	//v := reflect.ValueOf(c.Replies)
	//nameList := make([]string, 0)
	//for i := 0; i < v.NumField(); i++ {
	//	field := v.Type().Field(i)
	//	if !field.IsExported() {
	//		continue
	//	}
	//	nameList = append(nameList, field.Name)
	//}
	//sort.Strings(nameList)
	//for _, s := range nameList {
	//	if strings.Contains(s, search) {
	//		println(s)
	//	}
	//}
}

// Defaultize 引数 r をデフォルト値で埋める。破壊的メソッド
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

func PrintGrpcErr(err error) {
	if s, ok := status.FromError(err); ok {
		println(s.Message())
		//for _, i2 := range s.Details() {
		//	switch message := i2.(type) {
		//	case *pb_status.ErrorDialog:
		//		fmt.Printf("error code: %d\n", message.GetErrorCode())
		//		fmt.Printf("title: %s\n", message.GetTitle())
		//		fmt.Printf("message: %s\n", message.GetMessage())
		//		fmt.Printf("buttons: %v\n", message.GetButtons())
		//		fmt.Printf("domain: %s\n", message.GetDomain())
		//	case *errdetails.DebugInfo:
		//		fmt.Println("Stacktrace")
		//		for _, stack := range message.GetStackEntries() {
		//			fmt.Printf("%s\n", stack)
		//		}
		//	default:
		//		fmt.Printf("%#v\n", message)
		//	}
		//}
	}
}
