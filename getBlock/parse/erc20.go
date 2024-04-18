package parse

const (
	erc20TransferHash     = "0xa9059cbb"
	erc20TransferFromHash = "0x23b872dd"
)

const uint256Size = 64

const hexadecimalPrefix = "0x"

func Input(transactionInput string) (from, to string) {
	if len(transactionInput) < 10+2*64 {
		return
	}

	if transactionInput[0:10] == erc20TransferHash {
		to = transactionInput[10 : 10+uint256Size]
		to = hexadecimalPrefix + to[len(to)-40:]
	}

	if transactionInput[0:10] == erc20TransferFromHash {
		from = transactionInput[10 : 10+uint256Size]
		from = hexadecimalPrefix + from[len(from)-40:]

		to = transactionInput[10+uint256Size : 10+2*uint256Size]
		to = hexadecimalPrefix + to[len(to)-40:]
	}

	return
}
