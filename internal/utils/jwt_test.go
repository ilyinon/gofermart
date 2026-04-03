package utils

import "testing"

func TestJWT(t *testing.T) {

	manager := NewJWTManager("test-secret")

	token, err := manager.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	id, err := manager.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Fatalf("expected 1 got %d", id)
	}
}
