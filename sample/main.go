package main

import (
	"context"
	"fmt"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/sample/api"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"
)

func main() {
	s := gshell.New(gshell.Config{
		Host:        "localhost:8080",
		IsSecure:    false,
		ContextFunc: Context,
		ErrorFunc:   PrintGrpcErr,
		Scenario:    &Scenario{},
	})
	register(s)
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
	//if c.Replies.Login != nil {
	if RegisterReply(c) != nil {
		//ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+c.Replies.Login.GetToken())
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+RegisterReply(c).GetUserId())
	}
	//if c.Config.Trace GetToken{
	//	span := trace.FromContext(ctx)
	//	if span == nil {
	//		return ctx
	//	}
	//	// HTTP Header
	//	ctx = metadata.AppendToOutgoingContext(ctx, "X-Cloud-Trace-Context", fmt.Sprintf("%s/%s;o=1", span.SpanContext().TraceID.String(), "1"))
	//	// gRPC Header
	//	traceContext := trace.SpanContext{
	//		TraceID:      span.SpanContext().TraceID,
	//		SpanID:       span.SpanContext().SpanID,
	//		TraceOptions: 1,
	//	}
	//	ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(propagation.Binary(traceContext)))
	//}
	return ctx
}

// 以下手書きのシナリオ

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
