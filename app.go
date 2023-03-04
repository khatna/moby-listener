package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pebbe/zmq4"
)

// listen to raw tx's from sock, decode it, and send to byte channel
func forwardRawTxs(ch chan []byte) {
	// connect to tcp sockets to listen to txos
	zctx, _ := zmq4.NewContext()
	sock, _ := zctx.NewSocket(zmq4.SUB)
	defer zctx.Term()
	defer sock.Close()

	err := sock.Connect(os.Getenv("BTC_TCP_ENDPOINT"))

	if err != nil {
		log.Fatal(err)
	}

	sock.SetSubscribe("rawtx")

	// listen on socket until err
	for {
		msg, err := sock.RecvMessage(0)

		if err != nil {
			fmt.Println(err)
			close(ch)
			return
		}

		// the main encoded raw transaction is the second message in a multipart message:
		// https://github.com/bitcoin/bitcoin/blob/master/doc/zmq.md
		rawtx := msg[1]
		decodeRawTransaction(hex.EncodeToString([]byte(rawtx)), ch)
	}
}

// read bytes from the decoded tx channel, create simple pb.Tx structures,
// and register them in the gRPC server
func directToServer(txJsonCh chan []byte, s *txHandlerService) {
	for txJson := range txJsonCh {
		// TODO: check err
		txs := jsonToTxs(txJson)

		if txs == nil {
			continue
		}

		for _, tx := range jsonToTxs(txJson) {
			s.newTransaction(tx)
		}
	}
}

func main() {
	if os.Getenv("BTC_TCP_ENDPOINT") == "" {
		log.Fatal("Please set the BTC_TCP_ENDPOINT evnironment variable.")
	}

	if os.Getenv("BTC_RPC_HOST") == "" {
		log.Fatal("Please set the BTC_RPC_HOST evnironment variable.")
	}

	if os.Getenv("BTC_RPC_USER") == "" {
		log.Fatal("Please set the BTC_RPC_USER evnironment variable.")
	}

	if os.Getenv("BTC_RPC_PASS") == "" {
		log.Fatal("Please set the BTC_RPC_PASS evnironment variable.")
	}

	var wg sync.WaitGroup

	// start gRPC server
	s := startServer()

	fmt.Println("Server started...")

	// create communication channels
	decodedTx := make(chan []byte)

	wg.Add(2)
	go forwardRawTxs(decodedTx)
	go directToServer(decodedTx, s)
	wg.Wait()
}
