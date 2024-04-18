package requests

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

type ResponseBlockNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}
type ResponseBlockByNumber struct {
	Error   error  `json:"error,omitempty"`
	Jsonrpc string `json:"jsonrpc,omitempty"`
	ID      string `json:"id,omitempty"`
	Result  struct {
		BaseFeePerGas         string        `json:"baseFeePerGas"`
		BlobGasUsed           string        `json:"blobGasUsed"`
		Difficulty            string        `json:"difficulty"`
		ExcessBlobGas         string        `json:"excessBlobGas"`
		ExtraData             string        `json:"extraData"`
		GasLimit              string        `json:"gasLimit"`
		GasUsed               string        `json:"gasUsed"`
		Hash                  string        `json:"hash"`
		LogsBloom             string        `json:"logsBloom"`
		Miner                 string        `json:"miner"`
		MixHash               string        `json:"mixHash"`
		Nonce                 string        `json:"nonce"`
		Number                string        `json:"number"`
		ParentBeaconBlockRoot string        `json:"parentBeaconBlockRoot"`
		ParentHash            string        `json:"parentHash"`
		ReceiptsRoot          string        `json:"receiptsRoot"`
		Sha3Uncles            string        `json:"sha3Uncles"`
		Size                  string        `json:"size"`
		StateRoot             string        `json:"stateRoot"`
		Timestamp             string        `json:"timestamp"`
		TotalDifficulty       string        `json:"totalDifficulty"`
		Transactions          []Transaction `json:"transactions"`
		TransactionsRoot      string        `json:"transactionsRoot"`
		Uncles                []interface{} `json:"uncles"`
		Withdrawals           []struct {
			Index          string `json:"index"`
			ValidatorIndex string `json:"validatorIndex"`
			Address        string `json:"address"`
			Amount         string `json:"amount"`
		} `json:"withdrawals"`
		WithdrawalsRoot string `json:"withdrawalsRoot"`
	} `json:"result,omitempty"`
}

type Transaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	AccessList           []interface{} `json:"accessList,omitempty"`
	ChainID              string        `json:"chainId"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	YParity              string        `json:"yParity,omitempty"`
	MaxFeePerBlobGas     string        `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes  []string      `json:"blobVersionedHashes,omitempty"`
}
