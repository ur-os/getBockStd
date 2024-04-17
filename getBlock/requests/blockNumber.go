package requests

import (
	"encoding/json"
	"math/big"
	"net/http"
)

const parseIntBase = 16

func GetBlockNumber(nodeEndpoint string, client *http.Client) (blockNumber *big.Int, err error) {
	requestBody := Request{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]interface{}, 0),
		ID:      "blockNumber",
	}

	resp, err := doNodeRequest(requestBody, nodeEndpoint, client)
	if err != nil {
		return nil, err
	}

	responseBody := ResponseBlockNumber{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	responseBlockNumber, ok := new(big.Int).SetString(responseBody.Result[2:], parseIntBase)
	if !ok {
		return nil, err
	}

	return responseBlockNumber, nil
}
