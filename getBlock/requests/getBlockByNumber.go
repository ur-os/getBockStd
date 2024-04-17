package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const blockNumber = 0
const transactionFormat = 1

func GetBlocksByNumber(
	fromBlock, toBlock int64,
	nodeEndpoint string,
	client *http.Client,
) (blocks []ResponseBlockByNumber, err error) {
	batchRequestBody := make([]Request, toBlock-fromBlock)

	for blockNumb := fromBlock; blockNumb < toBlock; blockNumb++ {
		params := make([]interface{}, 2)
		params[blockNumber] = fmt.Sprintf("0x%x", blockNumb)
		params[transactionFormat] = true

		batchRequestBody = append(batchRequestBody, Request{
			Jsonrpc: "2.0",
			Method:  "eth_getBlockByNumber",
			Params:  params,
			ID:      "getblock.io",
		})
	}

	resp, err := doNodeRequest(batchRequestBody, nodeEndpoint, client)
	if err != nil {
		return nil, err
	}

	responseBody := make([]ResponseBlockByNumber, toBlock-fromBlock)

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
