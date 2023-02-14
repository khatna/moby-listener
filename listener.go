package main

import (
	"encoding/hex"
	"fmt"

	"github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq4.NewContext()
	sock, _ := zctx.NewSocket(zmq4.SUB)
	defer sock.Close()
	defer zctx.Term()

	err := sock.Connect("tcp://127.0.0.1:3600")

	if err != nil {
		fmt.Println("error while connecting to socket:", err)
	}

	sock.SetSubscribe("rawtx")

	// Receive messages
	for {
		msg, err := sock.RecvMessage(0)
		rawtx, seq := msg[1], msg[2]

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(hex.EncodeToString([]byte(rawtx)))
		fmt.Println(hex.EncodeToString([]byte(seq)))
	}
}
