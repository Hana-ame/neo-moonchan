// grpc_client.go
package main

import (
	"context"
	"fmt"

	pb "github.com/Hana-ame/neo-moonchan/gRPC/calculator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)

	res, _ := client.Add(context.Background(), &pb.AddRequest{A: 5})
	fmt.Printf("Result: %d\n", res.Result) // 输出: Result: 8
}
