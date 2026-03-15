package utils

import "testing"

func TestValidLuhn(t *testing.T) {

	valid := "79927398713"

	if !ValidLuhn(valid) {
		t.Fatal("expected valid luhn")
	}

	invalid := "123456"

	if ValidLuhn(invalid) {
		t.Fatal("expected invalid luhn")
	}
}
