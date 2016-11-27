package accounts

import (
	"testing"
)

func TestGenerateSalt_and_HashPassword(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt Failed : %v\n", err)
	}
	expectedPassword := "expectedPassword"

	hashedPassword, err := HashPassword(expectedPassword, salt)
	if err != nil {
		t.Fatalf("HashPassword Failed : %v\n", err)
	}

	hashedPassword2, err := HashPassword(expectedPassword, salt)
	if err != nil {
		t.Fatalf("HashPassword Failed : %v\n", err)
	}

	if len(hashedPassword) != len(hashedPassword2) {
		t.Fatal("HashPassword didn't return same hash twice")
	}
	for i, x := range hashedPassword {
		if x != hashedPassword2[i] {
			t.Fatal("HashPassword didn't return same hash twice")
		}
	}

}
