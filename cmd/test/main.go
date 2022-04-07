package main

import (
	"context"

	"github.com/j-tokumori/gshell"
	"github.com/j-tokumori/gshell/cmd/test/api"
	"google.golang.org/grpc"
)

func main() {
	s := gshell.New()
	s.Register(NewRegister)
	s.Start(gshell.Config{})
}

// TODO
// * Client CallWithJSON() で Register を参照できない
// * Scenario で Hoge Boot() を呼べない

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

// 以下手書きのシナリオ

type Hoge struct {
}

func (h *Hoge) Boot() {

}
