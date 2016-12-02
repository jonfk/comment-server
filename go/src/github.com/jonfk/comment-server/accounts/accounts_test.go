package accounts

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
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

func TestGenerateJWT_and_validateJWT(t *testing.T) {
	hmacSecretKey := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	now := time.Now().UTC().Round(time.Second).Add(-2 * time.Second)

	inputs := []map[string]interface{}{
		map[string]interface{}{
			"accountId": uuid.NewV4(),
			"issuedAt":  now,
			"expiresAt": now.Add(256 * time.Hour),
			"valid":     true,
		},
		map[string]interface{}{
			"accountId": uuid.NewV4(),
			"issuedAt":  now.Add(1 * time.Hour),
			"expiresAt": now.Add(256 * time.Hour),
			"valid":     false,
		},
		map[string]interface{}{
			"accountId": uuid.NewV4(),
			"issuedAt":  now,
			"expiresAt": now.Add(-256 * time.Hour),
			"valid":     false,
		},
	}

	for _, inputs := range inputs {
		token, err := generateJWT([]byte(hmacSecretKey), inputs["accountId"].(uuid.UUID), inputs["issuedAt"].(time.Time), inputs["expiresAt"].(time.Time))
		if err != nil {
			t.Fatalf("generateJWT Failed : %v\n", err)
		}

		validatedAccountId, err := validateJWT([]byte(hmacSecretKey), token)
		if inputs["valid"].(bool) {
			if err != nil {
				t.Fatalf("validateJWT: failed %v\ninputs: %v\n", err, inputs)
			}

			if !uuid.Equal(validatedAccountId, inputs["accountId"].(uuid.UUID)) {
				t.Fatalf("accountId != validatedAccountId\n (accountId) %v != (validatedAccountId) %v\ninputs: %v", inputs["accountId"].(uuid.UUID), validatedAccountId, inputs)
			}
		} else {
			if err == nil {
				t.Fatalf("Invalid token succeeded validation with inputs %v", inputs)
			}
		}

	}

}
