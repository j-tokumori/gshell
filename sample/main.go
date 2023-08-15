package main

import (
	"context"
	"fmt"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/sample/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	host := "localhost:8080"
	opts := []gshell.Option{
		gshell.WithInsecure(),
		gshell.WithCallingChain(
			HandleContext,
			HandleError,
		),
		RegisterRPC(),
		gshell.RegisterScenarioFactory(NewScenario),
		gshell.RegisterRPCAliasHandler(HandleRPCAlias),
	}
	s := gshell.New(host, opts...)
	s.Start()
}

func HandleContext(ctx context.Context, c *gshell.Client, r gshell.RPC, invoker gshell.CallingInvoker) (*gshell.Response, error) {
	println("context1")
	return invoker(ctx, c, r)
}

func HandleError(ctx context.Context, c *gshell.Client, r gshell.RPC, invoker gshell.CallingInvoker) (*gshell.Response, error) {
	res, err := invoker(ctx, c, r)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Println("gRPC Error.")
			fmt.Println(s.Message())
		} else {
			fmt.Println(err)
		}
	}
	return res, err
}

func HandleRPCAlias(serviceName, methodName string) string {
	return methodName
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
	println("boot up!")
	s.call(&SampleService_Hello{})
	gshell.NewReplyCommand().Exec(s.c, "")
}

func (_ *Scenario) Test(i int, s string) {
	println(i)
	println(s)
}

// Default 手書きデフォルト値
func (r *SampleService_Hello) Default(c *gshell.Client) *grpc.HelloRequest {
	return &grpc.HelloRequest{
		Name: "Sample",
	}
}
