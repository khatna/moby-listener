package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/khatna/moby-listener/proto"
	"google.golang.org/grpc"
)

// a gRPC client connection
type connection struct {
	ch     chan *pb.Tx
	minVal float32
}

// implementation of Transaction handler server
type txHandlerService struct {
	pb.UnimplementedTxHandlerServer
	conns map[*connection]struct{} // map to void <-> set ADT
}

// gRPC procedure implementation
func (s *txHandlerService) GetTransactions(value *pb.Value, stream pb.TxHandler_GetTransactionsServer) error {
	fmt.Println("New Connection")
	conn := &connection{make(chan *pb.Tx), value.Value}
	s.conns[conn] = struct{}{}
	defer delete(s.conns, conn)

	// send messages over stream concurrently
	errorCh := make(chan error)
	go func() {
		for tx := range conn.ch {
			if err := stream.Send(tx); err != nil {
				errorCh <- err
			}
		}
	}()

	// stop when the client is done, or there is an error while sending
	select {
	case err := <-errorCh:
		fmt.Println("Error while sending transaction: ", err)
		return err
	case <-stream.Context().Done():
		fmt.Println("Client done.")
		return nil
	}
}

// direct transactions to relevant connections
func (s *txHandlerService) newTransaction(tx *pb.Tx) error {
	for conn := range s.conns {
		if conn.minVal <= tx.Value {
			conn.ch <- tx
		}
	}
	return nil
}

// Create a new gRPC server, start it, then return it
func startServer() *txHandlerService {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	s := &txHandlerService{conns: make(map[*connection]struct{})}
	pb.RegisterTxHandlerServer(grpcServer, s)
	go grpcServer.Serve(lis)
	return s
}
