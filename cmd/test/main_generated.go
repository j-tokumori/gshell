package main

import (
	"context"
	"encoding/json"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/cmd/test/api"
	"google.golang.org/grpc"
)

func register(s *gshell.Shell) {
	s.RegisterRPC("Register", NewRegister)
}

// TODO
// * Scenario で Hoge Boot() を呼べない
// * Context
// * ErrorWrap

// 以下generate 想定

type Register api.RegisterArgs

func (r *Register) Call(ctx context.Context, cli *gshell.Client) {
	client := api.NewAuthServiceClient(cli.Conn)
	args := api.RegisterArgs(*r)
	reply, err := client.Register(ctx, &args, grpc.Header(cli.Header("Register")), grpc.Trailer(cli.Trailer("Register")))
	if err != nil {
		panic(err)
	}
	cli.Replies["Register"] = reply
}

func NewRegister(in []byte) gshell.RPC {
	r := &Register{}
	err := json.Unmarshal(in, r)
	if err != nil {
		panic(err)
	}
	return r
}
