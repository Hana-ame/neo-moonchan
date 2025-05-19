// grpc_server.go
package main

import (
	"context"
	"net"

	pb "github.com/Hana-ame/neo-moonchan/gRPC/calculator"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Result: req.A + req.B}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})
	s.Serve(lis)
}
