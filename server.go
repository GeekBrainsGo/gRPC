package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	pb "gorpc/proto"
)

type Server struct {
	pb.UnimplementedDemoServer
}

func NewServerStart(conf GRPCConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := &Server{}
	serv := grpc.NewServer()

	fmt.Println("server started ...")

	pb.RegisterDemoServer(serv, s)
	if err := serv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func NewHTTPServerStart(conf GRPCConfig, httpConf HTTPConfig) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	gRPCEndpoint := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	err := pb.RegisterDemoHandlerFromEndpoint(context.Background(), mux, gRPCEndpoint, opts)
	if err != nil {
		return err
	}
	fmt.Println("http proxy started ...")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", httpConf.Host, httpConf.Port), mux)
}

func (s *Server) ServerInfo(ctx context.Context, in *pb.EmptyRequest) (*pb.InfoReply, error) {
	log.Printf("ServerInfo request")
	return &pb.InfoReply{
		Info: "Our gRPC server version 1.0.0",
	}, nil
}

func (s *Server) Echo(ctx context.Context, in *pb.EchoRequestReply) (*pb.EchoRequestReply, error) {
	log.Printf("Echo request")
	return in, nil
}
