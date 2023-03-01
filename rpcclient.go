// the RPC client functions call the HTTP RPC API and return the response
// in caller provided channels (in order to take advantage of concurrency)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// body is plain text json, as per BTC RPC API
func sendPostRequest(method string, params []interface{}) (*http.Response, error) {
	// get bitcoin RPC credentials from environment
	var url = os.Getenv("BTC_RPC_HOST")
	var username = os.Getenv("BTC_RPC_USER")
	var password = os.Getenv("BTC_RPC_PASS")

	// Construct JSON body
	body := make(map[string]interface{})
	body["method"] = method
	body["params"] = params
	jsonBody, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// BTC JSON RPC uses basic authorization
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	return res, err
}

// decodes a raw transaction with Hex rawtxHex, and sends the decoded JSON string to ch
// https://developer.bitcoin.org/reference/rpc/decoderawtransaction.html
// https://bitcoincore.org/en/doc/22.0.0/rpc/rawtransactions/decoderawtransaction/
func decodeRawTransaction(rawtxHex string, ch chan []byte) {
	method := "decoderawtransaction"
	params := []interface{}{rawtxHex}
	res, err := sendPostRequest(method, params)

	if err != nil {
		return
	}
	defer res.Body.Close()

	// TODO: check by status code
	bodyBytes, _ := io.ReadAll(res.Body)

	ch <- bodyBytes
}
