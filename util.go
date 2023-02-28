package main

import (
	"encoding/json"
	"fmt"

	pb "github.com/khatna/moby-listener/proto"
)

type vin struct {
	Txid string // the id of the transaction providing the funds of this input
	Vout int    // the index of the input address in the vout array of the prev tx
}

type scriptPubKey struct {
	Address string
}

type vout struct {
	Value        float64 // amount in BTC
	ScriptPubKey scriptPubKey
}

type unmarshalledTx struct {
	Result struct {
		Txid string
		Vin  []vin
		Vout []vout
	}
}

// This utility functions converts a JSON decoded transaction to a list of pb.Tx's
func jsonToTxs(jsonBytes []byte) []*pb.Tx {
	txs := make([]*pb.Tx, 0)
	var unmarshalled unmarshalledTx

	err := json.Unmarshal(jsonBytes, &unmarshalled)

	if err != nil {
		fmt.Println("Error encountered while unmarshalling JSON: ", err)
		return nil
	}

	// create from string
	var from string
	if len(unmarshalled.Result.Vin) > 1 {
		from = fmt.Sprintf("%v wallets", len(unmarshalled.Result.Vin))
	} else {
		from = "1 wallet"
	}

	// create pb.Tx structs
	for _, out := range unmarshalled.Result.Vout {
		txs = append(txs, &pb.Tx{
			Txid:  unmarshalled.Result.Txid,
			Value: float32(out.Value),
			From:  from,
			To:    out.ScriptPubKey.Address,
		})
	}

	return txs
}
