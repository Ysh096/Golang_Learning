package main

import (
	"fmt"

	"github.com/Ysh096/Golang_Learning/2_Bank_DictionaryPjt/Bank/accounts"
)

func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10)
	fmt.Println(account)
}
