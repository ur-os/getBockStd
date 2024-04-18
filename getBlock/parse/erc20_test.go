package parse

import (
	"testing"
)

func TestInput(t *testing.T) {
	tests := []struct {
		name             string
		transactionInput string
		expectedFrom     string
		expectedTo       string
	}{
		{
			name:             "ERC20 Transfer",
			transactionInput: "0xa9059cbb0000000000000000000000001234567890123456789012345678901234567890",
			expectedFrom:     "",
			expectedTo:       "0x1234567890123456789012345678901234567890",
		},
		{
			name:             "ERC20 TransferFrom",
			transactionInput: "0x23b872dd00000000000000000000000012345678901234567890123456789012345678900000000000000000000000001234567890123456789012345678901234567890",
			expectedFrom:     "0x1234567890123456789012345678901234567890",
			expectedTo:       "0x1234567890123456789012345678901234567890",
		},
		{
			name:             "Short Transaction Input",
			transactionInput: "0xa9059cbb",
			expectedFrom:     "",
			expectedTo:       "",
		},
		{
			name:             "Invalid Transaction Input",
			transactionInput: "0xa9059c",
			expectedFrom:     "",
			expectedTo:       "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			from, to := Input(test.transactionInput)

			if from != test.expectedFrom {
				t.Errorf("Expected 'from' to be %s, but got %s", test.expectedFrom, from)
			}

			if to != test.expectedTo {
				t.Errorf("Expected 'to' to be %s, but got %s", test.expectedTo, to)
			}
		})
	}
}
