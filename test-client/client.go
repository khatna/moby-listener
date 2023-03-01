// Package main implements a simple gRPC client to test the main listener
package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/khatna/moby-listener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTxHandlerClient(conn)

	// Looking for a valid feature
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stream, err := client.GetTransactions(ctx, &pb.Value{Value: 5.0})
	if err != nil {
		log.Fatal(err)
	}

	for {
		tx, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.GetTransactions failed: %v", err)
		}
		log.Println(tx)
	}
}
