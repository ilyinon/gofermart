package utils

import "testing"

func TestPasswordHash(t *testing.T) {

	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	ok := CheckPassword("secret", hash)

	if !ok {
		t.Fatal("password should match")
	}

	ok = CheckPassword("wrong", hash)

	if ok {
		t.Fatal("password should not match")
	}
}
