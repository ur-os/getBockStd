package aggregation

import (
	"reflect"
	"testing"
)

func TestGetTopAddresses(t *testing.T) {
	tests := []struct {
		name        string
		topAmount   int
		vault       map[string]int
		expectedTop []TopAddresses
		expectedErr error
	}{
		{
			name:        "Empty Vault",
			topAmount:   5,
			vault:       map[string]int{},
			expectedTop: []TopAddresses{},
			expectedErr: nil,
		},
		{
			name:      "Top Amount Greater Than Vault Length",
			topAmount: 10,
			vault: map[string]int{
				"address1": 5,
				"address2": 3,
			},
			expectedTop: []TopAddresses{
				{Address: "address1", Count: 5},
				{Address: "address2", Count: 3},
			},
			expectedErr: nil,
		},
		{
			name:      "Valid Top Addresses",
			topAmount: 2,
			vault: map[string]int{
				"address1": 5,
				"address2": 3,
				"address3": 7,
			},
			expectedTop: []TopAddresses{
				{Address: "address3", Count: 7},
				{Address: "address1", Count: 5},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			top, err := GetTopAddresses(test.topAmount, test.vault)

			if !reflect.DeepEqual(top, test.expectedTop) {
				t.Errorf("Expected top addresses to be %v, but got %v", test.expectedTop, top)
			}

			if err != test.expectedErr {
				t.Errorf("Expected error to be %v, but got %v", test.expectedErr, err)
			}
		})
	}
}
