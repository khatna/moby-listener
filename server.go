package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/khatna/moby-listener/proto"
	"google.golang.org/grpc"
)

type connection struct {
	ch     chan *pb.Tx
	minVal float32
}

// implementation of Transaction handler server
type TxHandlerServer struct {
	pb.UnimplementedTxHandlerServer
	conns map[*connection]struct{} // map to void <-> set ADT
}

// gRPC procedure implementation
func (s *TxHandlerServer) GetTransactions(value *pb.Value, stream pb.TxHandler_GetTransactionsServer) error {
	fmt.Println("New Connection")
	conn := &connection{make(chan *pb.Tx), value.Value}
	s.conns[conn] = struct{}{}
	for tx := range conn.ch {
		fmt.Println(tx.Value)
	}
	delete(s.conns, conn)
	return nil
}

// direct transactions to relevant connections
func (s *TxHandlerServer) NewTransaction(tx *pb.Tx) error {
	for conn := range s.conns {
		if conn.minVal <= tx.Value {
			conn.ch <- tx
		}
	}
	return nil
}

// Starts a server
func StartServer() *TxHandlerServer {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	s := &TxHandlerServer{conns: make(map[*connection]struct{})}
	pb.RegisterTxHandlerServer(grpcServer, s)
	go grpcServer.Serve(lis)
	return s
}
