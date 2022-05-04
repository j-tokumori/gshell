package main

import (
	"context"
	"fmt"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/sample/api"
	"google.golang.org/grpc/status"
)

func main() {
	s := gshell.New(gshell.Config{
		Host:        "localhost:8080",
		IsSecure:    false,
		ContextFunc: Context,
		ErrorFunc:   PrintGrpcErr,
		Scenario:    &Scenario{},
	})
	RegisterRPC(s)
	s.Start()
}

func PrintGrpcErr(err error) {
	if s, ok := status.FromError(err); ok {
		fmt.Println("main")
		fmt.Println(s.Message())
	}
}

// Context コンテキスト
func Context(ctx context.Context, c *gshell.Client) context.Context {
	return ctx
}

type Scenario struct {
}

func (s *Scenario) Boot(c *gshell.Client) {
	c.Call(&Register{})
	println("user_id: " + RegisterReply(c).GetUserId())
}

// 手書きデフォルト値
func (r *Register) Default(c *gshell.Client) *api.RegisterArgs {
	return &api.RegisterArgs{
		Country:        "jp",
		Platform:       "apple",
		PlatformUserId: "test",
		DeviceName:     "gshell",
	}
}
