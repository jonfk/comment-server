package accounts

import (
	"crypto/rand"
	"database/sql"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/scrypt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Accounts struct {
	DB *sqlx.DB
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

func (a *Accounts) Verify(accountId uuid.UUID, unhashedPassword string) (bool, error) {
	account, err := a.GetAccountByAccountId(accountId)
	if err != nil {
		return false, err
	}

	hashedPassword, err := HashPassword(unhashedPassword, account.HashSalt)
	if err != nil {
		return false, err
	}

	if len(hashedPassword) != len(account.HashedPassword) {
		return false, nil
	}

	for i, x := range hashedPassword {
		if x != account.HashedPassword[i] {
			return false, nil
		}
	}
	return true, nil
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
	return account, err
}

func (a *Accounts) GetAccountByUsername(username string) (Account, error) {
	var account Account
	err := a.DB.Get(&account, "SELECT account_id,username,email,hashed_password,hash_salt,created_on FROM accounts where username = $1",
		username)
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
