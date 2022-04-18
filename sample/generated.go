// Code generated by gengshell DO NOT EDIT.
// versions:
//  gshell 0.0.1

package main

import (
	"context"
	"encoding/json"

	api "github.com/j-tokumori/gshell/sample/api"

	"github.com/j-tokumori/gshell"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// register ...
func register(s *gshell.Shell) {

	s.RegisterRPC("Register", NewRegister)
}

// Register ...
type Register api.RegisterArgs

// Call ...
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

// Key ...
func (r *Register) Key() string {
	return "Register"
}

// NewRegister ...
func NewRegister(in []byte) gshell.RPC {
	r := &Register{}
	err := json.Unmarshal(in, r)
	if err != nil {
		panic(err)
	}
	return r
}

// RegisterReply ...
func RegisterReply(c *gshell.Client) *api.RegisterReply {
	if c.Replies["Register"] == nil {
		return nil
	}
	return c.Replies["Register"].(*api.RegisterReply)
}