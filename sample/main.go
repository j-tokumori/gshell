package main

import (
	"context"
	"fmt"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/sample/api"
	"google.golang.org/grpc/status"
)

func main() {
	host := "localhost:8080"
	opts := []gshell.Option{
		//gshell.WithSecure(),
		gshell.WithInsecure(),
		//gshell.WithCodec(nil),
		gshell.WithCallingChain(
			Context1,
			Context2,
		),
		RegisterRPC(),
		gshell.RegisterScenarioFactory(NewScenario),
		gshell.RegisterErrorHandler(PrintGrpcErr),
	}
	s := gshell.New(host, opts...)
	s.Start()
}

func PrintGrpcErr(err error) {
	if s, ok := status.FromError(err); ok {
		fmt.Println("main")
		fmt.Println(s.Message())
	}
}

func Context1(ctx context.Context, c *gshell.Client, r gshell.RPC, invoker gshell.CallingInvoker) (*gshell.Response, error) {
	println("context1")
	return invoker(ctx, c, r)
}

func Context2(ctx context.Context, c *gshell.Client, r gshell.RPC, invoker gshell.CallingInvoker) (*gshell.Response, error) {
	println("context2")
	return invoker(ctx, c, r)
}

type Scenario struct {
	c *gshell.Client
}

func NewScenario(c *gshell.Client) interface{} {
	return &Scenario{c}
}

func (s *Scenario) call(r gshell.RPC) {
	s.c.Call(r)
}

func (s *Scenario) Boot() {
	s.call(&Register{})
	println("user_id: " + RegisterReply(s.c).GetUserId())
}

func (_ *Scenario) Test(i int, s string) {
	println(i)
	println(s)
}

// Default 手書きデフォルト値
func (r *Register) Default(c *gshell.Client) *api.RegisterArgs {
	return &api.RegisterArgs{
		Country:        "jp",
		Platform:       "apple",
		PlatformUserId: "test",
		DeviceName:     "gshell",
	}
}
