# 3. URL CHECKER & GO ROUTINES

배열 형태의 URL을 받아와서 체크한다.

## 3.0 hitURL

go lang std library - http - 어떻게 request 하는지 확인!

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request failed")

func main() {
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://www.nomadcoders.co/",
	}
	for _, url := range urls {
		hitURL(url)
	}
}

func hitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return errRequestFailed
	}
	return nil
}
```

그냥 단순하게 url을 하나씩 돌면서 해당 페이지에 요청을 해본 것. 어려울거 없다.

어떤 에러가 있거나 상태 코드가 400 이상인 경우 문제가 발생한 것이므로 요청에 실패했다는 것을 return해 준다.



## 3.1 Slow URL Checker

map은 초기화 후 사용해야 한다.

```go
var results map[string]string
```

이렇게만 작성하고 map에 뭔가 넣으려고 하면 panic이 발생한다. (컴파일러가 무슨 오류인지 잡아내지 못함)

```go
results["Hi"] = "Hello"

=> panic
```



초기화에는 다음 두 가지 방법이 있다.

```go
var result map[string]{}
```

```go
var results = make(map[string]string)
```

둘 중 한 가지 방법으로 빈 map을 만들 수 있다. 이렇게 하지 않으면 result는 nil이 됨



이렇게 만든 results에 URL을 체크한 결과를 담을 것이다. 다음과 같다.

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request failed")

func main() {
	var results = make(map[string]string)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://www.nomadcoders.co/",
	}

	for _, url := range urls {
		result := "OK"
		err := hitURL(url)
		if err != nil {
			result = "FAILED"
		}
		results[url] = result
	}
	for url, result := range results {
		fmt.Println(url, result)
	}
}

func hitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		fmt.Println(resp.StatusCode)
		return errRequestFailed
	}
	return nil
}
```

이 작업은 Python, Java로도 충분히 할 수 있다. 그런데 왜 Go를 쓸까?

=> 순서대로 요청을 처리하는 파이썬(멀티 프로세싱이 생겼지만, 이전에는 x)이나 자바와는 달리 매우 쉽게 URL 전체를 동시에 요청(Goroutines)할 수 있기 때문에!



## 3.2 Goroutines

다른 함수와 동시에 실행시키는 함수. 병렬 처리

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	sexyCount("nico")
	sexyCount("flynn")
}

func sexyCount(person string) {
	for i := 0; i < 10; i++ {
		fmt.Println(person, "is sexy", i)
		time.Sleep(time.Second)
	}
}
```

이렇게 작성하면 함수 실행에 20초가 걸린다.



이제 nico와 flynn을 넣은 함수가 동시에 동작하게 하려면, 둘 중 하나에 nico에 go를 붙이면 된다.

```go
func main() {
	go sexyCount("nico")
	sexyCount("flynn")
}
```

이렇게 하면, `sexyCount("nico")`는 flynn과 함께 실행되게 된다. 만약 flynn에 go를 붙이면 nico가 10초 동안 실행된 후 flynn이 실행이 되는데, goroutine은 main 함수가 동작하는 동안에만 실행이 되기 때문에 sexyCount("flynn") 이후로 아무것도 없는 상황에서는 main 함수가 종료되어 goroutine도 그대로 종료된다.

따라서 flynn에 go를 붙이면 nico만 0에서 9까지 카운트 한 후 함수가 종료된다.

- go를 함수에 붙이면 다른 함수가 실행될 때 함께 실행되도록 변경된다.
- 다른 실행될 함수가 없으면 go를 붙인 함수는 실행되지 않는다.
- nico와 flynn에 모두 go를 붙이면 마찬가지로 다른 go가 붙지 않은 함수가 없어서 nico와 flynn의 sexyCount는 실행되지 않는다.



## 3.3 Channels

goroutine과 main 간, 혹은 goroutine 끼리 소통하는 방법

```go
func main() {
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person)
	}
}

func isSexy(person string) bool {
	time.Sleep(time.Second * 5)
	return true
}
```

이런 함수가 있다고 하자. 일단 main function이 바로 끝나기 때문에 아무 동작도 일어나지 않을 것이다. time.sleep 등으로 main이 동작한다고 해보자. return true를 하고 있는걸 보고`result = go isSexy(person)` 이런식으로 코드를 작성해도 결과를 받을 수가 없다. 그렇다면 어떻게 해야 goroutine의 결과를 얻을 수 있을까?



**channel을 사용**해야 한다. 사용법은 아래와 같다.

```go
func main() {
	c := make(chan bool) // boolean 값을 받는 채널 만들기
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c) // isSexy 함수에 채널을 엮기
	}
}

func isSexy(person string, c chan bool) { // isSexy 함수에 채널을 인수로 추가
	time.Sleep(time.Second * 5)
	c <- true // 채널에 원하는 boolean 값 전달
}
```

채널을 만들고, 채널을 함수로 보낸다. 그러면 이 함수는 5초 후에 true라는 메시지를 채널로 보낸다.



```
					   main
					     |
				   	  channel
					    /  \
					  url  url
```

우리는 1개의 채널이 있고,  해당 채널에 goroutine의 결과가 전달되며, 그게 다시 main으로 전달된다.



channel을 사용할 때에는 굳이 time.Sleep을 main에서 사용해서 함수가 끝나길 기다리거나 할 필요가 없다.

```go
func main() {
	c := make(chan bool)
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c)
	}
	result := <-c
	fmt.Println(result)
}

func isSexy(person string, c chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println(person)
	c <- true
}
```

go에서 자동으로 c에서 어떤 값이 전달되기를 기다리기 때문이다. 위 코드는 아래와 같은 코드이다.

```go
func main() {
	c := make(chan bool)
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println(<-c)
}

func isSexy(person string, c chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println(person)
	c <- true
}
```

그리고 우리는 c에 두 개의 함수 결과를 전달하고 있으므로 다음과 같이 수정하면 두 결과를 모두 확인할 수 있다.

```go
func main() {
	c := make(chan bool)
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println(<-c)
    fmt.Println(<-c)
}

func isSexy(person string, c chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println(person)
	c <- true
}
```

그니까 <-c 는 함수가 끝날 때마다 어떤 값을 main에 전달하게 되는 것으로 보인다. 위 함수를 실행하면 결과는 다음과 같다.

```
nico
true
flynn
true
```

순서는 매번 달라진다. 사실상 동시에 끝났다.



조금 더 정리해서 써보자면,

```go
func main() {
	c := make(chan string)
	people := [2]string{"nico", "flynn"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println("Waiting for messages")
	resultOne := <-c
	resultTwo := <-c
	fmt.Println("Received this message:", resultOne)
	fmt.Println("Received this message:", resultTwo)
}

func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 5)
	c <- person + " is sexy"
}
```



그런데 사람이 많아지면 저렇게 하나하나 적는 것은 매우 불편하다.

이 때는 loop를 사용한다.

```go
func main() {
	c := make(chan string)
	people := [5]string{"nico", "flynn", "dal", "japanguy", "larry"}
	for _, person := range people {
		go isSexy(person, c)
	}
	fmt.Println("Waiting for messages")
	for i := 0; i < len(people); i++ {
		fmt.Println(<-c)
	}
}

func isSexy(person string, c chan string) {
	time.Sleep(time.Second * 5)
	c <- person + " is sexy"
}
```



메시지를 받는 부분은 blocking operation이다. 무슨 뜻이냐면, main이라는 함수가 blocking operation인 `fmt.Println(<-c)`에서 c값을 받을 때 까지 잠시 멈춰있는다는 뜻이다.

 

이제 goroutine을 이용해서 fast URLChecker를 만들어 볼 것이다.

## 3.4 Fast URLChecker

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

type result struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request Failed")

func main() {
	// results = make(map[string]string)
	c := make(chan result)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://www.nomadcoders.co/",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}
	for i := 0; i < len(urls); i++ {
		fmt.Println(<-c)
	}
}

func hitURL(url string, c chan<- result) { // send only
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
		c <- result{url: url, status: status}
	} else {
		c <- result{url: url, status: status}
	}
}
```

goroutine을 적용했더니 훨씬 빠르다. 작동 시간이 모든 url 체크 시간의 합이 아니라, 가장 체크가 오래 걸리는 url의 시간으로 바뀌었다.



조금 더 코드를 가다듬어 보자.

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

type requestResult struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request Failed")

func main() {
	results := make(map[string]string)
	c := make(chan requestResult)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://www.nomadcoders.co/",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}

	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}
}

func hitURL(url string, c chan<- requestResult) { // send only
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
		c <- requestResult{url: url, status: status}
	} else {
		c <- requestResult{url: url, status: status}
	}
}
```

