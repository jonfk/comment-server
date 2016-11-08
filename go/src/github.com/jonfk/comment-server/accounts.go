package main

import (
	"github.com/lib/pq"
)

type Accounts struct {
}

type Account struct {
	Username         string
	Email            string
	PasswordHashed   string
	PasswordUnHashed string
}

func (a Accounts) CreateNewAccount(account Account) {

}

func (a Accounts) Verify(account Account) bool {
	return false
}

func (a Accounts) GetAccount() {

}

func (a Accounts) DeleteAccount() {

}
