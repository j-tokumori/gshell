package main

import (
	"context"
	"encoding/json"

	"github.com/j-tokumori/gshell"
	api2 "github.com/j-tokumori/gshell/test/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func register(s *gshell.Shell) {
	s.RegisterRPC("Register", NewRegister)
}

// TODO
// * Context
// * ErrorWrap
// * Generate
// * 拡張コマンド Trace URL

type Replies struct {
	Register *api2.RegisterReply
}

// 以下generate 想定

type Register api2.RegisterArgs

func (r *Register) Call(ctx context.Context, conn grpc.ClientConnInterface) (proto.Message, *metadata.MD, *metadata.MD) {
	client := api2.NewAuthServiceClient(conn)
	args := api2.RegisterArgs(*r)
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

func RegisterReply(c *gshell.Client) *api2.RegisterReply {
	if c.Replies["Register"] == nil {
		return nil
	}
	return c.Replies["Register"].(*api2.RegisterReply)
}

func NewRegister(in []byte) gshell.RPC {
	r := &Register{}
	err := json.Unmarshal(in, r)
	if err != nil {
		panic(err)
	}
	return r
}
