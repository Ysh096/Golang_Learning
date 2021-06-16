# Go

## 1. Scrapper 만들기

- kr.indeed.com의 데이터를 가져와서 CSV 파일로 만들기

- Go를 사용하는 이유?
  - 속도가 파이썬에 비해 훨씬 빠르다.
  - multi-core processing을 이용



### 1-1. Go 설치

- 그냥 golang.org가서 설치하면 된다.
- Go는 어떤 임의의 위치에 코드를 저장할 수 없고, windows의 경우 C:\Go에만 Go 코드를 작성하여 넣어야 한다.
- 환경변수 설정 참고
  - https://artist-developer.tistory.com/4



### 1-2. 상수, 변수, 함수 작성

```go
package main

import (
	"fmt"
	"strings"
)
// 함수 작성시 매개변수의 타입과 결과의 타입을 써줘야 함
func multiply(a, b int) int { //a, b 모두 int면 이렇게 축약 가능, 결과도 int임을 써줌
	return a * b
}
func lenAndUpper(name string) (int, string) {
	return len(name), strings.ToUpper(name)
}

func repeatMe(words ...string) { // 여러 개의 입력(몇개 인지는 미정)
	fmt.Println(words)
}
func main() {
	fmt.Println("Hello World")

	const name string = "seongho" // constant, 적절한 타입과 함께 적어줘야 함
	// name = "Lynn" 상수는 재할당 불가능
	// fmt.Println(name)

	var varname string = "seongho" // variable
	// shorthand
	shorthandname := "nico" // 축약형, 알아서 타입을 찾아준다. func 안에서만 사용 가능하다.
	varname = "lynn" // 변수는 재할당 가능
	fmt.Println(multiply(2, 2))
	fmt.Println(varname)
	fmt.Println(shorthandname)

	totalLength, upperName := lenAndUpper("nico") // 여러 개 리턴받기
	// totalLength, _ := lenAndUpper(("nico")) // 하나 무시하기
	fmt.Println(totalLength, upperName)

	repeatMe("nico", "lynn", "dal", "marl")
}
```



### 1-3. naked return ,defer

```go
func lenAndUpper(name string) (length int, uppercase string) {
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}
```

return 할 변수 이름을 미리 써서 return에는 아무것도 적지 않을 수 있다.



defer를 사용하면 어떤 func이 끝나고 난 후 특정 동작이 실행되도록 할 수 있다.

```go
func lenAndUpper(name string) (length int, uppercase string) {
	defer fmt.Println("I'm done") //function이 끝나고 나면 실행됨
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}
```

이 경우 lenAndUpper를 실행하면 length와 uppercase를 return하여 사용자가 원하는 변수에 할당한 후에 "I'm done"이라는 메시지를 출력하게 된다.



 ### 1-4. for문

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for number := range numbers {
		fmt.Println(number)
	}
	return 1
}
func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}
```

int인 numbers를 입력받음(여러 개), 결과도 int

for문을 이렇게 쓰면 출력 결과는 0, 1, 2, 3, 4, 5가 된다. 왜냐면 range는 index를 주기 때문이다.

number를 출력하려면 다음과 같이 하면 된다.

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for index, number := range numbers {
		fmt.Println(index, number)
	}
	return 1
}
func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}

```

출력 결과는

0 1

1 2

2 3

3 4

4 5

5 6



다음과 같이 작성할 수도 있다.

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for i := 0; i < len(numbers); i++ {
		fmt.Println(numbers[i])
	}
	return 1
}
func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}
```

이 경우 출력 결과는 1, 2, 3, 4, 5, 6



```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}
func main() {
	fmt.Println(superAdd(1, 2, 3, 4, 5, 6))
}
```



### 1-5. if-else

```go
package main

import "fmt"

func canIDrink(age int) bool {
	if age < 18 {
		return false
	} else {
		return true
	}
}

func main() {
	fmt.Println(canIDrink(21))
}
```

js랑 비슷한데, 여기서 else를 안써도 된다.

```go
package main

import "fmt"

func canIDrink(age int) bool {
	if age < 18 {
		return false
	}
	return true
}

func main() {
	fmt.Println(canIDrink(21))
}
```

이렇게..



if문을 만들 때 variable을 만들 수도 있다. (**variable expression**)

```go
package main

import "fmt"

func canIDrink(age int) bool {
	if koreanAge := age + 2; koreanAge < 18 {
		return false
	}
	return true
}

func main() {
	fmt.Println(canIDrink(16))
}
```

이렇게 작성하는 경우는 if문 안에서만 koreanAge라는 변수를 사용하고자 할 때..



### 1-7. Switch

기본적인 switch

```go
package main

import "fmt"

func canIDrink(age int) bool {
	switch age {
	case 10:
		return false
	case 18:
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(15))
}
```



variable expression 사용 가능!!

```go
package main

import "fmt"

func canIDrink(age int) bool {
	switch koreanAge := age + 2; {
	case koreanAge <= 18:
		return false
	case koreanAge > 18:
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(17))
}
```



### 1-8. Pointers

Go는 높은 수준의 코드로 Low-level programming을 할 수 있도록 해준다.

```go
package main

import "fmt"

func main() {
	a := 2
	b := a
	a = 10
	fmt.Println(a, b)
}

결과
10 2
```

b가 정해지고 나면 a의 영향을 받지 않는다.



메모리 주소를 보는 법 => & 추가!

```go
package main

import "fmt"

func main() {
	a := 2
	b := a
	a = 10
	fmt.Println(&a, &b)
}
결과
0xc000016098 0xc0000160b0
```



이걸 응용해서, b에 a의 주소를 할당하고 출력해보면,

```go
package main

import "fmt"

func main() {
	a := 2
	b := &a
	fmt.Println(&a, b)
}
결과
0xc000016098 0xc000016098
```

b에는 a의 주소가 저장된다.



b에 저장된 주소가 어떤 값인지 확인하려면 *를 붙이면 된다.

```go
package main

import "fmt"

func main() {
	a := 2
	b := &a
	fmt.Println(&a, *b)
}
결과
0xc000016098 2
```



a의 값이 나중에 바뀌면?

```go
package main

import "fmt"

func main() {
	a := 2
	b := &a
	a = 5
	fmt.Println(&a, *b)
}
결과
0xc000016098 5
```

재할당되어도 주소는 변하지 않으므로 바뀐 값을 볼 수 있다.



이제 b를 가지고 a를 바꿀 수도 있다.

```go
package main

import "fmt"

func main() {
	a := 2
	b := &a // b는 a의 주소를 값으로 할당받는다.
	*b = 20 // b가 가리키는 주소의 내용을 20으로 바꾼다. => a가 바뀐다.
	fmt.Println(a)
}
결과
20
```



### 1-9. Arrays and Slices

배열 만드는게 좀 특이하다.

```go
package main

import "fmt"

func main() {
	names := [5]string{"nico", "lynn", "dal"}
	names[3] = "alala"
	names[4] = "alala"
	names[5] = "alala"
	fmt.Println(names)
}
```

5개 string 원소를 가지는 배열 names에 nico, lynn, dal이 초기값으로 들어있음.

3, 4, 5번째 원소에 alala를 추가하려고 함 => 오류! index 하나가 모자름. 인덱스는 0, 1, 2, 3, 4이기 때문에!



```go
package main

import "fmt"

func main() {
	names := [5]string{"nico", "lynn", "dal"}
	names[3] = "alala"
	names[4] = "alala"
	fmt.Println(names)
}

결과
[nico lynn dal alala alala]
```



Array의 크기 제약을 해결하는 방법? => slice

=> 그냥 길이 부분에 아무것도 작성하지 않는다.

=> 이걸 slice라고 부름

```go
package main

import "fmt"

func main() {
	names := []string{"nico", "lynn", "dal"}

	fmt.Println(names)
}

결과
[nico lynn dal]
```

여기에 새로운 원소를 추가하려면 그냥 names[idx] = "abc" 이런 형식으로는 안되고, append를 사용해야 한다.

```go
package main

import "fmt"

func main() {
	names := []string{"nico", "lynn", "dal"}
	fmt.Println(names)
	names = append(names, "flynn")
	fmt.Println(names)
}
```

좀 특이한게, append에는 두 개의 인자가 필요하다. 하나는 새로운 내용을 추가하고자 하는 슬라이스의 이름, 그리고 나머지는 추가하려는 원소의 값. append(names, "flynn")만 한다고 해서 names에 "flynn"이 추가되는건 아니고, append를 사용하면 names에 원소를 추가한 새로운 슬라이스를 리턴하게 된다. 따라서 names를 바꾸고 싶으면 names = append(names, "flynn") 형식으로 써야 한다.



### 1-10. map

python이나 JS의 object와 비슷하다.

```go
package main

import "fmt"

func main() {
	nico := map[string]string{"name": "nico", "age": "12"} //map의 key는 string, value도 string
	fmt.Println(nico)
	for key, value := range nico {
		fmt.Println(key, value)
	}
}
```

이런 식으로 key와 value의 형식을 지정해주는 방식이다. 딕셔너리 비슷한데 type이 고정된 dictionary라고 볼 수 있지 않을까?

map을 다루는 함수들은 나중에 배운다.



### 1-10. struct (구조체 같은 것)

파이썬의 dictionary와 같은걸 만들려면? => struct

```go
package main

import "fmt"

// struct 정의
type person struct {
	name    string
	age     int
	favFood []string //빈 slice
}

func main() {
	// value 할당 방법1
	// favFood := []string{"kimchi", "ramen"}
	// nico := person{"nico", 18, favFood}
	// fmt.Println(nico.name)

	// value 할당 방법2
	favFood := []string{"kimchi", "ramen"}
	nico := person{name: "nico", age: 18, favFood: favFood}
	fmt.Println(nico.name)
}
```

value를 할당할 때, 방법1이나 방법2 중 하나를 써야하지, 둘을 섞어 쓸 수는 없다.

ex) {name: "nico", 18, favFood} 이런 식으로는 사용 불가



**struct는 매우 중요**하다.



앞으로 세 개의 프로젝트를 할 것

dictionary

bank account

URL checker

마지막으로 scrapper

