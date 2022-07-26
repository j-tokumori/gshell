package gshell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Command interface {
	Exec(*Client, ...string) bool
}

type RPCCommand struct {
}

func NewRPCCommand() *RPCCommand {
	return &RPCCommand{}
}

func (c *RPCCommand) Exec(client *Client, args ...string) bool {
	parsed := c.parseArgs(args[1])
	rpcTypeName := client.opts.getRPCTypeName(args[0])
	f := client.opts.rpcMap[rpcTypeName]
	rpc := f([]byte(parsed))
	client.Call(rpc)
	printResponse(client, rpcTypeName)
	return false
}

// parseArgs 引数のパース
// {} を付け足して、key に "" を雑につけているだけの簡易処理
// value に , や : があったり、入れ子データに対応していなかったりするので、修正必須
func (c *RPCCommand) parseArgs(str string) string {
	str = strings.TrimSpace(str)
	if str == "" {
		return "{}"
	}
	if str[0:1] == "{" { // 先頭が { なら生 json とみなす
		return str
	}
	args := strings.Split(str, ",")
	modArgs := make([]string, 0)
	for _, arg := range args {
		kv := strings.Split(arg, ":")
		if len(kv) != 2 {
			panic("TODO")
		}
		kv[0] = `"` + kv[0] + `"`
		modArgs = append(modArgs, strings.Join(kv[0:2], ":"))
	}
	return "{" + strings.Join(modArgs, ",") + "}"
}

type ScenarioCommand struct {
	ScenarioPlayer *ScenarioPlayer
}

func NewScenarioCommand(sp *ScenarioPlayer) *ScenarioCommand {
	return &ScenarioCommand{sp}
}

func (c *ScenarioCommand) Exec(client *Client, args ...string) bool {
	var sArgs []string
	if len(args) > 1 {
		sArgs = strings.Split(strings.TrimSpace(args[1]), " ")
	}
	c.ScenarioPlayer.Play(client, args[0], sArgs...)
	return false
}

type ResponseCommand struct {
}

func NewResponseCommand() *ResponseCommand {
	return &ResponseCommand{}
}

func (c *ResponseCommand) Exec(client *Client, args ...string) bool {
	if args[0] == "" {
		printResponse(client, client.lastRPCName)
	} else {
		printResponse(client, args[0])
	}
	return false
}

type ReplyCommand struct {
}

func NewReplyCommand() *ReplyCommand {
	return &ReplyCommand{}
}

func (c *ReplyCommand) Exec(client *Client, args ...string) bool {
	if args[0] == "" {
		printReply(client.responses[client.lastRPCName].Reply)
	} else {
		printReply(client.responses[args[0]].Reply)
	}
	return false
}

type HeaderCommand struct {
}

func NewHeaderCommand() *HeaderCommand {
	return &HeaderCommand{}
}

func (c *HeaderCommand) Exec(client *Client, args ...string) bool {
	if args[0] == "" {
		printMD(client.responses[client.lastRPCName].Header)
	} else {
		printMD(client.responses[args[0]].Header)
	}
	return false
}

type TrailerCommand struct {
}

func NewTrailerCommand() *TrailerCommand {
	return &TrailerCommand{}
}

func (c *TrailerCommand) Exec(client *Client, args ...string) bool {
	if args[0] == "" {
		printMD(client.responses[client.lastRPCName].Trailer)
	} else {
		printMD(client.responses[args[0]].Trailer)
	}
	return false
}

type SampleCommand struct {
}

func NewSampleCommand() *SampleCommand {
	return &SampleCommand{}
}

func (c *SampleCommand) Exec(client *Client, args ...string) bool {
	printSample(client, client.opts.getRPCTypeName(args[0]))
	return false
}

type ListCommand struct {
}

func NewListCommand() *ListCommand {
	return &ListCommand{}
}

func (c *ListCommand) Exec(client *Client, args ...string) bool {
	list := make([]string, 0)
	switch args[0] {
	case "rpc":
		for key := range client.opts.rpcMap {
			list = append(list, key)
		}
	case "alias":
		for key := range client.opts.rpcAliasMap {
			list = append(list, key)
		}
	}
	printList(list, args[1])
	return false
}

type TraceCommand struct {
}

func NewTraceCommand() *TraceCommand {
	return &TraceCommand{}
}

func (c *TraceCommand) Exec(client *Client, args ...string) bool {
	if args[0] == "" {
		printCloudTraceURL(client.responses[client.lastRPCName].Header)
	} else {
		printCloudTraceURL(client.responses[args[0]].Header)
	}
	return false
}

type EmptyCommand struct {
}

func NewEmptyCommand() *EmptyCommand {
	return &EmptyCommand{}
}

func (c *EmptyCommand) Exec(client *Client, args ...string) bool {
	return false
}

type HelpCommand struct {
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (c *HelpCommand) Exec(client *Client, args ...string) bool {
	fmt.Println("TODO: impl help... help you coming soon.")
	return false
}

type ExitCommand struct {
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

func (c *ExitCommand) Exec(client *Client, args ...string) bool {
	fmt.Println("exit...")
	return true
}

func printResponse(client *Client, rpcName string) {
	fmt.Println("Header:")
	printMD(client.responses[rpcName].Header)
	fmt.Println("Trailer:")
	printMD(client.responses[rpcName].Trailer)
	fmt.Println("Reply:")
	printReply(client.responses[rpcName].Reply)
}

func printReply(reply proto.Message) {
	o := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	j, err := o.Marshal(reply)
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

func printMD(md *metadata.MD) {
	for key, values := range *md {
		fmt.Printf("  %s: %v\n", key, values)
	}
}

func printSample(client *Client, rpcName string) {
	f := client.opts.rpcMap[rpcName]
	a := f([]byte("{}"))

	defaultize(a, client)
	samplize(a)

	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	fmt.Println("rpc " + rpcName + " " + string(j))
}

// samplize 引数 r にサンプル値を代入
// TODO: generics で書き直し
func samplize(r interface{}) {
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
			samplize(itf)
			pv.Elem().FieldByName(field.Name).Set(reflect.ValueOf(itf))
		case reflect.Slice:
			f := pv.Elem().FieldByName(field.Name)
			// TODO: prt ではなかった場合の対応
			if field.Type.Elem().Kind() == reflect.Ptr {
				itf := reflect.New(field.Type.Elem().Elem()).Interface()
				samplize(itf)
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

func printList(in []string, search string) {
	list := make([]string, 0, len(in))
	for _, name := range in {
		list = append(list, name)
	}
	sort.Strings(list)
	for _, s := range list {
		if strings.Contains(s, search) {
			fmt.Println(s)
		}
	}
}

func printCloudTraceURL(md *metadata.MD) {
	s := md.Get("x-cloud-trace-context")
	if len(s) <= 0 {
		fmt.Println("no trace Header.")
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
