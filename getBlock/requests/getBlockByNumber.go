package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const blockNumber = 0
const transactionFormat = 1

func GetBlockByNumber(curretnBlock int64, nodeEndpoint string, client *http.Client) (transactions []Transaction, err error) {
	requestBody := Request{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  make([]interface{}, 2),
		ID:      "getblock.io",
	}

	requestBody.Params[blockNumber] = fmt.Sprintf("0x%x", curretnBlock)
	requestBody.Params[transactionFormat] = true

	resp, err := doNodeRequest(requestBody, nodeEndpoint, client)
	if err != nil {
		return nil, err
	}

	responseBody := BlockByNumberResponse{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Result.Transactions, nil
}
