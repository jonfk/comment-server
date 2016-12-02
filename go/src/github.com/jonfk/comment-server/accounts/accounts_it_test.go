package accounts

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
)

var (
	DBUser, DBName, DBPassword string
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBUser = os.Getenv("DATABASE_USER")
	DBName = os.Getenv("DATABASE_NAME")
	DBPassword = os.Getenv("DATABASE_PASSWORD")

	os.Exit(m.Run())
}

func TestCreateNewAccount_and_DeleteById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DBUser, DBName, DBPassword))
	if err != nil {
		t.Fatalf("sqlx.Connect failed : %v\n", err)
	}

	accounts := &Accounts{DB: db,
		HMACSecretKey:        []byte("secret_key"),
		SessionLengthInHours: 256,
	}

	expectedAccount := Account{
		Username:  "username",
		Email:     "email",
		CreatedOn: time.Now().UTC().Round(time.Second),
	}

	expectedUnhashedPassword := "unhashedPassword"

	createdAccount, err := accounts.CreateNewAccount(expectedAccount, expectedUnhashedPassword)
	if err != nil {
		accounts.DeleteById(createdAccount.AccountId)
		t.Fatalf("accounts.CreateNewAccount failed : %v\n", err)
	}
	if createdAccount.Username != expectedAccount.Username &&
		createdAccount.Email != expectedAccount.Email &&
		createdAccount.CreatedOn != expectedAccount.CreatedOn {
		t.Fatalf("Created Account does not match expectedAccount\n (createdAccount) %v != (expectedAccount) %v)", createdAccount, expectedAccount)
	}

	expectedAccount.AccountId = createdAccount.AccountId
	expectedAccount.HashSalt = createdAccount.HashSalt
	expectedAccount.HashedPassword, err = HashPassword(expectedUnhashedPassword, expectedAccount.HashSalt)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	testGetAccountByAccountId(t, expectedAccount, accounts)

	testGetAccountByEmail(t, expectedAccount, accounts)

	testGetAccountByUsername(t, expectedAccount, accounts)

	testVerify(t, expectedAccount.AccountId, expectedUnhashedPassword, accounts)

	accounts.DeleteById(expectedAccount.AccountId)
}

func testGetAccountByAccountId(t *testing.T, expectedAccount Account, accounts *Accounts) {
	fetchedAccountByAccountId, err := accounts.GetAccountByAccountId(expectedAccount.AccountId)
	if err != nil {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByAccountId failed : %v\n", err)
	}

	if !fetchedAccountByAccountId.Equal(expectedAccount) {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByAccountId failed :\n (expectedAccount) %v != (fetchedAccount) %v\n", expectedAccount, fetchedAccountByAccountId)
	}
}

func testGetAccountByEmail(t *testing.T, expectedAccount Account, accounts *Accounts) {
	fetchedAccountByEmail, err := accounts.GetAccountByEmail(expectedAccount.Email)
	if err != nil {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByEmail failed : %v\n", err)
	}

	if !fetchedAccountByEmail.Equal(expectedAccount) {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByEmail failed :\n (expectedAccount) %v != (fetchedAccount) %v\n", expectedAccount, fetchedAccountByEmail)
	}
}

func testGetAccountByUsername(t *testing.T, expectedAccount Account, accounts *Accounts) {
	fetchedAccountByUsername, err := accounts.GetAccountByUsername(expectedAccount.Username)
	if err != nil {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByUsername failed : %v\n", err)
	}

	if !fetchedAccountByUsername.Equal(expectedAccount) {
		accounts.DeleteById(expectedAccount.AccountId)
		t.Fatalf("accounts.GetAccountByUsername failed :\n (expectedAccount) %v != (fetchedAccount) %v\n", expectedAccount, fetchedAccountByUsername)
	}
}

func testVerify(t *testing.T, accountId uuid.UUID, expectedUnhashedPassword string, accounts *Accounts) {

	err := accounts.Verify(accountId, expectedUnhashedPassword)
	if err != nil {
		t.Fatalf("accounts.Verify failed : %v", err)
	}

	token, err := accounts.VerifyAndGenerateJWT(accountId, expectedUnhashedPassword)
	if err != nil {
		t.Fatalf("accounts.VerifyAndGenerateJWT failed : %v", err)
	}

	validatedAccount, err := accounts.ValidateJWT(token)
	if err != nil {
		t.Fatalf("accounts.ValidateJWT failed : %v", err)
	}

	if !uuid.Equal(accountId, validatedAccount) {
		t.Fatalf("expectedAccountId != validatedAccountId\n(expectedAccountId) %v != (validatedAccountId) %v", accountId, validatedAccount)
	}
}

func TestGetAccountDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", DBUser, DBName, DBPassword))
	if err != nil {
		t.Fatalf("sqlx.Connect failed : %v\n", err)
	}

	accounts := &Accounts{DB: db}

	_, err = accounts.GetAccountByAccountId(uuid.NewV4())
	if err == nil {
		t.Fatal("Error should not be nil because account should not exist")
	}

	if err != AccountNotFoundErr {
		t.Fatalf("Wrong nil error returned %v", err)
	}

	_, err = accounts.GetAccountByEmail("dummy@email.com")
	if err == nil {
		t.Fatal("Error should not be nil because account should not exist")
	}

	if err != AccountNotFoundErr {
		t.Fatalf("Wrong nil error returned %v", err)
	}

	_, err = accounts.GetAccountByUsername("dummy_username")
	if err == nil {
		t.Fatal("Error should not be nil because account should not exist")
	}

	if err != AccountNotFoundErr {
		t.Fatalf("Wrong nil error returned %v", err)
	}

}
