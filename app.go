package main

import (
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pebbe/zmq4"
)

func main() {
	connCfg := &rpcclient.ConnConfig{
		Host:         os.Getenv("BTC_RPC_HOST"),
		User:         os.Getenv("BTC_RPC_USER"),
		Pass:         os.Getenv("BTC_RPC_PASS"),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// connect to tcp sockets to listen to txos
	zctx, _ := zmq4.NewContext()
	sock, _ := zctx.NewSocket(zmq4.SUB)
	defer zctx.Term()
	defer sock.Close()

	err = sock.Connect(os.Getenv("BTC_TCP_ENDPOINT"))

	if err != nil {
		log.Fatal(err)
	}

	sock.SetSubscribe("rawtx")

	// Receive messages
	for {
		msg, err := sock.RecvMessage(0)

		// the main encoded raw transaction is the second message
		// in a multipart message:
		// https://github.com/bitcoin/bitcoin/blob/master/doc/zmq.md
		rawtx := msg[1]

		if err != nil {
			fmt.Println(err)
			break
		}

		// Synchronously decode for now - no optimization necessary
		// https://developer.bitcoin.org/reference/rpc/decoderawtransaction.html
		rawtxResult, _ := client.DecodeRawTransaction([]byte(rawtx))

		fmt.Printf("%v - %v BTC\n", rawtxResult.Txid, rawtxResult.Vout[0].Value)
	}
}
