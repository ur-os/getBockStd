package vault

import (
	"testing"
)

func TestVault_SetAndGet(t *testing.T) {
	v := NewVault()

	v.Set("Key1", 10)

	if value := v.Get("Key1"); value != 10 {
		t.Errorf("Expected value to be 10, but got %d", value)
	}

	v.Set("Key2", 20)

	if value := v.Get("Key2"); value != 20 {
		t.Errorf("Expected value to be 20, but got %d", value)
	}

	v.Set("key3", 30)
	if value := v.Get("KEY3"); value != 30 {
		t.Errorf("Expected value to be 30, but got %d", value)
	}
}

func TestVault_GetVault(t *testing.T) {
	v := NewVault()

	v.Set("Key1", 10)
	v.Set("Key2", 20)

	expectedVault := map[string]int{
		"key1": 10,
		"key2": 20,
	}
	vault := v.GetVault()

	for key, value := range expectedVault {
		if vault[key] != value {
			t.Errorf("Expected value for key %s to be %d, but got %d", key, value, vault[key])
		}
	}
}
