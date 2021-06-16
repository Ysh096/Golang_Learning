# 2. Banking & Account

struct 연습!!



## 2-1. owner가 정해진 account 만들기

banking/banking.go

```go
package banking

// Account struct
type Account struct {
	// public으로 만들려면 대문자로 써야 함
	Owner   string
	Balance int
	// 소문자로 쓴 구성 요소는 private로, 다른 go 파일에서 사용할 수 없다.
}
```

구조체를 만들었다. 여기서 주의할 점은, 이 구조체를 외부에서 사용하려면 (export) 시작을 대문자로 만들어줘야 한다는 것, 그리고 owner의 경우도 아래에서처럼 직접 사용하려면 대문자로 시작해줘야 한다.

private과 public의 차이를 결정짓는 것은 소문자냐, 대문자냐!

main.go

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/first/banking"
)

func main() {
	account := banking.Account{Owner: "nicolas", Balance: 1000}
	fmt.Println(account)
}
```



그런데 이렇게 하면 누구나 Balance(자산)에 접근할 수 있게 되는 문제가 있다. 이걸 막고, Owner만 접근할 수 있게 할 수는 없을까?

그러기 위해서는 balance를 private으로 설정해야 한다. 그리고 Owner도 public이라 사실 마음대로 바꿀 수 있기 때문에 private으로 설정해 보자.



banking을 accounts로 바꾸고 코드를 다음과 같이 수정

first/accounts/accounts.go

```go
package accounts

// Account struct
type Account struct {
	// public으로 만들려면 대문자로 써야 함
	owner   string
	balance int
	// 소문자로 쓴 구성 요소는 private로, 다른 go 파일에서 바로 사용할 수 없다.
}

// NewAccount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}
```



first/main.go

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/first/accounts"
)

func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	fmt.Println(account)
}
```

```
결과
&{nico 0}
```

우리는 결과로 account의 주소를 보고 있다. 계속 복사본을 만들기 싫어서 이렇게 주소를 활용함!

function인 NewAccount는 owner를 정해서 계좌를 만들 수 있다!

의문: *Account를 return으로 써놓고 왜 &account를 반환할까?



어쨋든 이제 이런식으로 owner나 balance를 수정할 수 **없게** 되었다.

```go
func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	fmt.Println(account)
	account.balance = 10000
	account.owner = "elen"
}
```

그럼 이제 balance를 수정하려면 어떻게 해야 할까?

=> method 등장



## 2-2. balance 값 바꾸기

method는 func와 비슷하지만 receiver라고 하는게 추가된다.

```go
// function
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// method
func (a Account) Deposit(amount int) {

}
```

Deposit 앞에 쓴 (a Account)에서 a가 바로 receiver라고 하는 것이고, 그 타입이 Account 이다.

receiver의 이름은 마음대로 할 수 있지만, 그냥 규칙으로 타입 이름의 맨 앞글자를 소문자로 만들어 사용한다고 한다.

이제 이 receiver를 이용해서 balance에 접근할 수 있다. 다음과 같은 방법으로..

```go
func (a Account) Deposit(amount int) {
	a.balance += amount
    a.owner = "dalan"
}
```



정리를 해보면, NewAccount 라는 func은 owner를 결정해서 계좌를 새로 만들 수 있고, Deposit이라는 method는 Account라는 구조체를 타입으로 하는 객체에 사용할 수 있는 함수가  되어 그 객체의 내부 특정 값을 수정하게 된다.

그런데 지금은 수정이 안된다.. 원인은 좀 있다가! 일단 수정한 코드는 다음과 같다.

```go
package accounts

// Account struct
type Account struct {
	// public으로 만들려면 대문자로 써야 함
	owner   string
	balance int
	// 소문자로 쓴 구성 요소는 private로, 다른 go 파일에서 바로 사용할 수 없다.
}

// NewAccount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int {
	return a.balance
}
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/first/accounts"
)

func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10000)
	fmt.Println(account.Balance())
}
```

```
결과
0 // 수정이 안됨. 원인은??
```



이유: Go에서 object와 struct에 관여하는 부분 때문이다.

Go는 func이나 method 등 뭔가를 보내는 순간에 원본이 아닌 복사본을 보내버린다. 예를 들어, 

```go
func main() {
	account := accounts.NewAccount("nico")
	account.Deposit(10000) // 이 부분!!
	fmt.Println(account.Balance())
}
```

account.Deposit()을 사용하게 되면

```go
// Deposit x amount on your account
func (a Account) Deposit(amount int) {
	a.balance += amount
}
```

Deposit의 receiver가 account를 받아오게 되는데, 이 때 복사본을 만들어서 받아오게 된다. 그래서 위 코드에서 **a는 실제 account가 아니라 account의 복사본**이 된다. 그래서 func main에서 print를 하면 원본을 보여주므로 변화가 나타나지 않는다.



이를 고치기 위해서 해야 할 것은 한 가지, Deposit 함수에 *를 추가하는 것이다. 코드는 다음과 같다.

```go
// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}
```

이제 Deposit은 Deposit method를 호출한 Account 구조체인 account를 receiver가 받아 사용하게 된다. (찬찬히 읽어보자..)



## 2-3. errors

withdraw 기능을 만들어보자.

```go
// Withdraw x amount from your account
func (a *Account) Withdraw(amount int) {
	a.balance -= amount
}
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/first/accounts"
)

func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10)
	fmt.Println(account.Balance())
	account.Withdraw(20)
	fmt.Println(account.Balance())
}
```

역시 *Account를 사용해야 복사가 되지 않고 정상적인 결과가 출력된다.

결과는 -10이 나온다.



돈이 모자라면 withdraw 할 수 없게 만들어야 한다. => error handling

```go
// Withdraw x amount from your account
func (a *Account) Withdraw(amount int) error { // error를 return하기 시작했으므로 써줘야 한다.
	if a.balance < amount {
		return errors.New("Can't withdraw you are poor")
	}
	a.balance -= amount
	// 맨 처음 return을 정했으므로 여기도 뭔가 return을 해줘야 한다.
	return nil // none, null 같은 것
}
```

위와 같이 코드를 작성하면 error를 return 할 수 있다.

그런데 이렇게 바꾸고 main.go를 실행하면 error 메시지는 뜨지 않고 그냥 withdraw가 실행되지 않는다. 어떤 error 메시지를 보고 싶다면 직접 써줘야 한다.

```go
func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(20)
	if err != nil {
		log.Fatalln(err) // error를 발생시키고 코드를 종료
	}
	fmt.Println(account.Balance())
}
```

이렇게 쓸 수도 있고, 코드가 중간에 종료되는걸 원치 않으면 그냥 fmt를 사용할 수도 있다.

```go
func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account.Balance())
}
```



error를 변수로 직접 만들 수도 있다.

```go
...
type Account struct {
	// public으로 만들려면 대문자로 써야 함
	owner   string
	balance int
	// 소문자로 쓴 구성 요소는 private로, 다른 go 파일에서 바로 사용할 수 없다.
}

var errNoMoney = errors.New("Can't withdraw") // errErrorName 형식으로 error를 만들 수 있다.

// NewAccount creates Account
func NewAccount(owner string) *Account {
...
    
    
// Withdraw x amount from your account
func (a *Account) Withdraw(amount int) error { // error를 return하기 시작했으므로 써줘야 한다.
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	// 맨 처음 return을 정했으므로 여기도 뭔가 return을 해줘야 한다.
	return nil // none, null 같은 것
}
```

이런 식으로 사용 가능.



## 2-4. Go가 내부적으로 호출하는 함수(string)

마치 파이썬의 `__str__` 처럼 Go에서도 struct의 출력을 어떤 형식으로 할지 내부적으로 지정한는 함수 string이 있다. 이걸 우리가 임의로 바꿀 수 있다.

accounts.go

```go
// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return "whatever you want"
}
```

 main.go

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/first/accounts"
)

func main() {
	account := accounts.NewAccount("nico") // owner를 nico로 하는 계좌 생성, 그 주소값 반환
	account.Deposit(10)
	fmt.Println(account)
}
```

이 상태로 출력을 해보면, 결과는 whatever you want 라는 문자열이 된다. 형식을 내가 원하는 형식으로 바꿔 출력하려면 다음과 같이 할 수 있다.

accounts.go

```go
// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance()) // 원하는 출력 형식
    // 이제 Println(account)를 하면 위 return 값이 출력된다.
}
```

```
결과
$ go run main.go
nico's account.
Has: 10
```



## 2-5. Dictionary

type은 method를 가질 수 있다.

Dictionary라는 type을 만들고 method를 이용하여 search, delete 등을 구현해보자.



Dictionary type 구현

```go
package mydict

// Dictionary type
type Dictionary map[string]string // map[keytype]valuetype
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	dictionary["hello"] = "hello"
	fmt.Println(dictionary)
}
```



### Search method

Dictionary type에서 Search method를 만들어 사용해보자.

```go
package mydict

// Dictionary type
type Dictionary map[string]string // map[keytype]valuetype

func (d Dictionary) Search(word string) (string, error) {

}
```

```go
package main

import (
	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "First word"}
	dictionary.Search("second")
}
```

여기까지 하면 아직 func Search를 채우기 않았기 때문에 당연히 동작하지 않을 것이다. 우리가 만약 "second" 라는 문자를 찾으려고 하면, 없어서 error가 발생해야 한다. 이렇게 error 처리까지 고려해서 mydict를 완성해보자.



```go
package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string // map[keytype]valuetype

var errNotFound = errors.New("Not Found")

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word] // value는 찾은 string값, exists는 존재 여부(boolean)
	if exists {
		return value, nil // zero value error nil
	}
	return "", errNotFound
}
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "First word"}
	definition, err := dictionary.Search("second")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```

err가 있으면 err를 출력, 없으면 Search가 dictionary에서 key값을 찾았다는 뜻이므로 key값에 해당하는 value를 출력한다.

출력 예를 들어보자.

```go
1. Search("second") // second라는 키를 찾으려고 함
=> Not Found

2. Search("first") // first라는 키를 찾으려고 함
=> First word

3. Search('fir') // fir 이라는 key를 찾으려고 함 => dictionary에 없음.
=> Not Found
```



### add method

```go
package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string // map[keytype]valuetype

var errNotFound = errors.New("Not Found")
var errWordExists = errors.New("That word already exists")

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word] // value는 찾은 string값, exists는 존재 여부(boolean)
	if exists {
		return value, nil // zero value error nil
	}
	return "", errNotFound
}

func (d Dictionary) Add(word, def string) error {
	// 아직 word가 dictionary에 없으면 추가 가능
	_, err := d.Search(word) // word가 이미 d에 있는지 파악
	if err == errNotFound {  // NotFound error는 word가 d에 없다는 뜻
		d[word] = def // 추가
	} else if err == nil { // word가 이미 d에 있으면
		return errWordExists // error를 return
	}
	return nil
}
```

Add는 위와 같이 만들 수 있고, switch를 이용해서 쓸 수도 있다.

```go
// Add a word to the dictionary (by switch)
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}
```

이제 이걸 main에서 실행해보자.

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	word := "hello"
	definition := "Greeting"
	err := dictionary.Add(word, definition)
	if err != nil {
		fmt.Println(err)
	}
	hello, _ := dictionary.Search(word)
	fmt.Println(hello)

	err2 := dictionary.Add(word, definition)
	if err2 != nil {
		fmt.Println(err2)
	}
}
```

처음에 빈 dictionary에 word를 key로, definition을 value로 가지는 원소가 추가된다. 만약 에러가 발생하면 에러를 보여준다.

이제 우리가 추가한 word라는 key를 찾아서 Println(hello)로 출력해보면 Greeting이 나온다.

다음으로, 똑같은 원소를 추가하려는 시도를 해본다. 그러면 err2는 nil이 아니라 이미 존재한다는 에러를 return받으므로 그 메시지를 보여준다.



### Update

```go
// Update a word
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	err := dictionary.Update(baseWord, "Second")
	if err != nil {
		fmt.Println(err)
	}
	word, _ := dictionary.Search(baseWord)
	fmt.Println(word)
}
```



### Delete

```go
// Delete a word
func (d Dictionary) Delete(word string) error {
	// map document 참고
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errCantDelete
	}
	return nil
}
```

```go
package main

import (
	"fmt"

	"github.com/Ysh096/2_Banking/second/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	baseWord := "hello"
	dictionary.Add(baseWord, "First")
	dictionary.Search(baseWord)
	dictionary.Delete(baseWord)
	word, err := dictionary.Search(baseWord)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(word)
}
```

삭제하고 Search를 다시 해보면 Not Found가 뜬다.

크게 어려울 거 없다. 문법이 헷갈릴 뿐!





