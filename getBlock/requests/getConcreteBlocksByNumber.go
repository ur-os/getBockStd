package requests

import (
	"encoding/json"
	"net/http"
)

func GetConcreteBlocksByNumber(
	requiredBlocks []string,
	nodeEndpoint string,
	client *http.Client,
) (blocks []ResponseBlockByNumber, err error) {
	var batchRequestBody []Request

	for _, blockNumb := range requiredBlocks {
		params := make([]interface{}, 2)
		params[blockNumber] = blockNumb
		params[transactionFormat] = true

		batchRequestBody = append(batchRequestBody, Request{
			Jsonrpc: "2.0",
			Method:  "eth_getBlockByNumber",
			Params:  params,
			ID:      blockNumb,
		})
	}

	resp, err := doNodeRequest(batchRequestBody, nodeEndpoint)
	if err != nil {
		return nil, err
	}

	responseBody := make([]ResponseBlockByNumber, len(requiredBlocks))

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
