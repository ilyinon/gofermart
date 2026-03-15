package utils

import "testing"

func TestJWT(t *testing.T) {

	token, err := GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	id, err := ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Fatalf("expected 1 got %d", id)
	}
}
