package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	// public으로 만들려면 대문자로 써야 함
	owner   string
	balance int
	// 소문자로 쓴 구성 요소는 private로, 다른 go 파일에서 바로 사용할 수 없다.
}

var errNoMoney = errors.New("Can't withdraw") // errErrorName 형식으로 error를 만들 수 있다.

// NewAccount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int {
	return a.balance
}

// Withdraw x amount from your account
func (a *Account) Withdraw(amount int) error { // error를 return하기 시작했으므로 써줘야 한다.
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	// 맨 처음 return을 정했으므로 여기도 뭔가 return을 해줘야 한다.
	return nil // none, null 같은 것
}

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
