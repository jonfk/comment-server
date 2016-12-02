package accounts

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/scrypt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Accounts struct {
	DB                   *sqlx.DB
	SessionLengthInHours int
	// From http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key
	HMACSecretKey []byte // Should be a 512 bits random key
}

type Account struct {
	AccountId      uuid.UUID `db:"account_id" json:"accountId"`
	Username       string    `db:"username" json:"username"`
	Email          string    `db:"email" json:"email"`
	HashedPassword []byte    `db:"hashed_password"`
	CreatedOn      time.Time `db:"created_on" json:"createdOn"`
	HashSalt       []byte    `db:"hash_salt"`
}

func (a Account) Equal(b Account) bool {
	if !a.CreatedOn.Equal(b.CreatedOn) ||
		a.Email != b.Email ||
		a.Username != b.Username ||
		!uuid.Equal(a.AccountId, b.AccountId) ||
		len(a.HashedPassword) != len(a.HashedPassword) ||
		len(a.HashSalt) != len(b.HashSalt) {
		return false
	}
	for i, byte := range a.HashedPassword {
		if byte != b.HashedPassword[i] {
			return false
		}
	}
	for i, byte := range a.HashSalt {
		if byte != b.HashSalt[i] {
			return false
		}
	}
	return true
}

func (a *Accounts) CreateNewAccount(account Account, unhashedPassword string) (Account, error) {
	var (
		newAccount Account
	)
	salt, err := GenerateSalt()
	if err != nil {
		return newAccount, err
	}

	account.HashedPassword, err = HashPassword(unhashedPassword, salt)
	if err != nil {
		return newAccount, err
	}

	err = a.DB.QueryRowx("INSERT INTO accounts (account_id, username,email,hashed_password,hash_salt,created_on) VALUES ($1,$2,$3,$4,$5,$6) RETURNING account_id,username,email,hashed_password,hash_salt,created_on",
		uuid.NewV4().String(), account.Username, account.Email, account.HashedPassword, salt, account.CreatedOn).StructScan(&newAccount)

	return newAccount, err
}

func (a *Accounts) DeleteById(accountId uuid.UUID) (string, error) {
	var deletedAccountId string
	err := a.DB.QueryRowx("DELETE FROM accounts where account_id = $1 RETURNING account_id", accountId).Scan(&deletedAccountId)

	return deletedAccountId, err
}

func (a *Accounts) Verify(accountId uuid.UUID, unhashedPassword string) error {
	account, err := a.GetAccountByAccountId(accountId)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(unhashedPassword, account.HashSalt)
	if err != nil {
		return err
	}

	if len(hashedPassword) != len(account.HashedPassword) {
		return fmt.Errorf("length of password does not match")
	}

	for i, x := range hashedPassword {
		if x != account.HashedPassword[i] {
			return fmt.Errorf("password does not match")
		}
	}
	return nil
}

func (a *Accounts) VerifyAndGenerateJWT(accountId uuid.UUID, unhashedPassword string) (string, error) {
	if err := a.Verify(accountId, unhashedPassword); err != nil {
		return "", err
	}
	now := time.Now().UTC().Round(time.Second).Add(-1 * time.Second)

	return generateJWT(a.HMACSecretKey, accountId, now, now.Add(time.Duration(a.SessionLengthInHours)*time.Hour))
}

func (a *Accounts) ValidateJWT(token string) (uuid.UUID, error) {
	return validateJWT(a.HMACSecretKey, token)
}

func (a *Accounts) GetAccountByAccountId(accountId uuid.UUID) (Account, error) {
	var account Account
	err := a.DB.Get(&account, "SELECT account_id,username,email,hashed_password,hash_salt,created_on FROM accounts where account_id = $1",
		accountId)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return account, AccountNotFoundErr
		default:
			return account, err
		}
	}
	return account, nil
}

func (a *Accounts) GetAccountByEmail(email string) (Account, error) {
	var account Account
	err := a.DB.Get(&account, "SELECT account_id,username,email,hashed_password,hash_salt,created_on FROM accounts where email = $1",
		email)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return account, AccountNotFoundErr
		default:
			return account, err
		}
	}
	return account, err
}

func (a *Accounts) GetAccountByUsername(username string) (Account, error) {
	var account Account
	err := a.DB.Get(&account, "SELECT account_id,username,email,hashed_password,hash_salt,created_on FROM accounts where username = $1",
		username)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return account, AccountNotFoundErr
		default:
			return account, err
		}
	}
	return account, err
}

// HashPassword returns a hashed password from the unhashed password and a salt
func HashPassword(unhashedPassword string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(unhashedPassword), []byte(salt), 16384, 8, 1, 32)
}

// GenerateSalt uses crypto/rand to generate a random string and returns a random string and a nil error
// or an empty string and an error
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 10)
	_, err := rand.Read(salt)
	if err != nil {
		return salt, err
	}
	return salt, nil
}

func generateJWT(hmacSecretKey []byte, accountId uuid.UUID, issuedAt, expiresAt time.Time) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aid": accountId.String(),
		"iat": issuedAt.Unix(),
		"exp": expiresAt.Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecretKey)
	return tokenString, err
}

func validateJWT(hmacSecretKey []byte, tokenString string) (uuid.UUID, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecretKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if accountId, ok := claims["aid"].(string); ok {
			return uuid.FromString(accountId)
		} else {
			return uuid.Nil, fmt.Errorf("validateJWT: cannot cast aid claim to string")
		}
	} else {
		return uuid.Nil, fmt.Errorf("validateJWT: invalid token with claims: %v", claims)
	}
}
