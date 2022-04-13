package main

import (
	"context"
	"encoding/json"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/cmd/test/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func register(s *gshell.Shell) {
	s.RegisterRPC("Register", NewRegister)
	//s.RegisterReplies(Replies{})
	s.RegisterScenario(&Scenario{})
}

// TODO
// * Context
// * ErrorWrap

type Replies struct {
	Register *api.RegisterReply
}

// 以下generate 想定

type Register api.RegisterArgs

func (r *Register) Call(ctx context.Context, conn grpc.ClientConnInterface) (proto.Message, *metadata.MD, *metadata.MD) {
	client := api.NewAuthServiceClient(conn)
	args := api.RegisterArgs(*r)
	h := &metadata.MD{}
	t := &metadata.MD{}
	reply, err := client.Register(ctx, &args, grpc.Header(h), grpc.Trailer(t))
	if err != nil {
		panic(err)
	}
	return reply, h, t
}

func (r *Register) Key() string {
	return "Register"
}

func NewRegister(in []byte) gshell.RPC {
	r := &Register{}
	err := json.Unmarshal(in, r)
	if err != nil {
		panic(err)
	}
	return r
}
