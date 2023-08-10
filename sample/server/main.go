package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	samplegrpc "github.com/j-tokumori/gshell/sample/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	samplegrpc.RegisterSampleServiceServer(s, newSampleServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		_ = s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

type sampleServer struct {
	samplegrpc.UnimplementedSampleServiceServer
}

func newSampleServer() samplegrpc.SampleServiceServer {
	return &sampleServer{}
}

func (sampleServer) Hello(ctx context.Context, req *samplegrpc.HelloRequest) (*samplegrpc.HelloResponse, error) {
	return &samplegrpc.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}
