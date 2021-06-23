# 5. GO Echo로 서버 만들기

[TOC]

## 5-1. setup

python 말고 다른 검색도 가능하게 만들기!

1. scrapper 폴더 만들고, main.go 파일 옮기기

2. main.go를 scrapper.go로 바꾸고, 패키지 이름도 scrapper로 바꾼다.

   main func를 scrape로 바꾸고, 내용을 다음과 같이 수정한다. 맨 처음 주석도 반드시 달아야 export 된다.

   ```go
   // Scrape Indeed by a term
   func Scrape(term string) {
   	var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=" + term + "&limit=50"
   	var jobs []extractedJob
   	c := make(chan []extractedJob)
       ...
   ```

   python으로 고정된 검색 결과가 아니라 term에 입력하는 것에 따라 검색 결과가 나오게 된다.

3. baseURL을 전역변수에서 Scrape 함수 안으로 옮겼기 때문에 문제가 발생한다. baseURL이 전달되도록 인자를 수정해준다. 수정 완료한 코드는 다음과 같다.

   ```go
   package scrapper
   
   import (
   	"encoding/csv"
   	"fmt"
   	"log"
   	"net/http"
   	"os"
   	"strconv"
   	"strings"
   
   	"github.com/PuerkitoBio/goquery"
   )
   
   type extractedJob struct {
   	id       string
   	title    string
   	location string
   	salary   string
   	summary  string
   }
   
   // Scrape Indeed by a term
   func Scrape(term string) {
   	var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=" + term + "&limit=50"
   	var jobs []extractedJob
   	c := make(chan []extractedJob)
   	totalPages := getPages(baseURL)
   	for i := 0; i < totalPages; i++ {
   		go getPage(i, baseURL, c) // 한 페이지에서 가져온 일자리
   		// 페이지별로 다음의 형태가 아니라 [[X] [X] [X] ...]
   		// 다음과 같은 형태로 정보가 합쳐짐 [X X X]
   	}
   	for i := 0; i < totalPages; i++ {
   		extractedJobs := <-c
   		jobs = append(jobs, extractedJobs...)
   	}
   	writeJobs(jobs)
   	fmt.Println("Done, extracted", len(jobs))
   }
   func getPage(page int, url string, mainC chan<- []extractedJob) {
   	var jobs []extractedJob
   	c := make(chan extractedJob)
   	pageURL := url + "&start=" + strconv.Itoa(page*50) // integer to ascii
   	fmt.Println("Requesting", pageURL)
   	res, err := http.Get(pageURL)
   	checkErr(err)
   	checkCode(res)
   
   	defer res.Body.Close()
   
   	doc, err := goquery.NewDocumentFromReader(res.Body)
   	checkErr(err)
   
   	searchCards := doc.Find(".jobsearch-SerpJobCard")
   	searchCards.Each(func(i int, card *goquery.Selection) {
   		go extractJob(card, c)
   	})
   	for i := 0; i < searchCards.Length(); i++ {
   		job := <-c
   		jobs = append(jobs, job)
   	}
   	mainC <- jobs // channel로 전달
   }
   
   func extractJob(card *goquery.Selection, c chan<- extractedJob) {
   	id, _ := card.Attr("data-jk")
   	title := cleanString(card.Find(".title>a").Text())
   	location := cleanString(card.Find(".sjcl").Text())
   	salary := cleanString(card.Find(".salaryText").Text())
   	summary := cleanString(card.Find(".summary").Text())
   	c <- extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
   }
   
   func getPages(url string) int {
   	pages := 0
   
   	res, err := http.Get(url)
   	checkErr(err)
   	checkCode(res)
   
   	defer res.Body.Close()
   
   	doc, err := goquery.NewDocumentFromReader(res.Body)
   	checkErr(err)
   
   	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) { // class 이름이 pagination인 태그를 찾는다.
   		pages = s.Find("a").Length()
   	})
   	return pages
   }
   
   func checkErr(err error) {
   	if err != nil {
   		log.Fatalln(err)
   	}
   }
   
   func checkCode(res *http.Response) { // type이 *http.Response
   	if res.StatusCode != 200 {
   		log.Fatalln("Request failed with Status", res.StatusCode)
   	}
   }
   
   func cleanString(str string) string {
   	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
   	// 공백을 지우고, 텍스트로 된 배열을 만든다. 그 후 배열을 합친다.
   	// hello     f     1  =>  "hello","f","1"  =>  hello f 1
   }
   
   // csv 파일로 저장하기
   func writeJobs(jobs []extractedJob) {
   	c := make(chan []string)
   	file, err := os.Create("jobs.csv")
   	checkErr(err)
   
   	w := csv.NewWriter(file)
   	defer w.Flush() // 함수가 끝나는 시점에 파일에 데이터를 입력
   
   	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}
   
   	wErr := w.Write(headers)
   	checkErr(wErr)
   
   	for _, job := range jobs {
   		go writeJob(job, c)
   	}
   	// jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
   	for i := 0; i < len(jobs); i++ {
   		jobData := <-c
   		writeErr := w.Write(jobData)
   		checkErr(writeErr)
   	}
   }
   
   func writeJob(job extractedJob, c chan<- []string) {
   	const jobURL = "https://kr.indeed.com/viewjob?jk="
   	c <- []string{jobURL + job.id, job.title, job.location, job.salary, job.summary}
   }
   ```

   

4. main.go 만들기

5. `go get github.com/labstack/echo/` 로 go echo 설치

6. main.go 다음과 같이 작성

```go
package main

import "github.com/Ysh096/Golang_Learning/5_Go_Echo/scrapper"

func main() {
	scrapper.Scrape("term")
}
```



7. https://echo.labstack.com/guide/

   위 주소에서 Hello, World! 부분을 사용한다.

   ```go
   // 위 페이지 내용
   func main() {
   	e := echo.New()
   	e.GET("/", func(c echo.Context) error {
   		return c.String(http.StatusOK, "Hello, World!")
   	})
   	e.Logger.Fatal(e.Start(":1323"))
   }
   
   // 우리가 입력한 내용
   func handleHome(c echo.Context) error {
   	return c.String(http.StatusOK, "Hello, World!")
   }
   
   func main() {
   	e := echo.New()
   	e.GET("/", handleHome)
   	e.Logger.Fatal(e.Start(":1323"))
   }
   ```

   함수를 두 개로 나눠서 입력했다. 이제 go run main.go 를 하면, localhost:1323으로 접속할 수 있다.



8. http://localhost:1323/

   Hello, world! 라는 글자가 보인다!



이제 입력 창을 만들어서 검색하기 위해 다음과 같은 html 문서를 만든다.

5_GO_Echo/home.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Go Jobs</title>
</head>
<body>
  <h1>Go Jobs</h1>
  <h3>indeed.com scrapper</h3>
  <form>
    <input placeholder="what job do you want">
    <button>Search</button>
  </form>
</body>
</html>
```

이제 main.go의 코드를 다음과 같이 수정하고 서버를 재시작하면 변화한 내용이 반영된다.

```go
package main

import (
	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	return c.File("home.html") // 바뀐 부분
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.Logger.Fatal(e.Start(":1323"))
}
```



9. 다시 다음과 같이 html과 main.go 를 수정하고 서버를 재시작한다.

   ```html
   <!DOCTYPE html>
   <html lang="en">
   <head>
     <meta charset="UTF-8">
     <meta http-equiv="X-UA-Compatible" content="IE=edge">
     <meta name="viewport" content="width=device-width, initial-scale=1.0">
     <title>Go Jobs</title>
   </head>
   <body>
     <h1>Go Jobs</h1>
     <h3>indeed.com scrapper</h3>
     <form method="POST" action="/scrape">
       <input placeholder="what job do you want" name="term">
       <button>Search</button>
     </form>
   </body>
   </html>
   ```

   ```go
   package main
   
   import (
   	"fmt"
   
   	"github.com/labstack/echo"
   )
   
   func handleHome(c echo.Context) error {
   	return c.File("home.html")
   }
   
   func handleScrape(c echo.Context) error {
   	fmt.Println(c.FormValue("term"))
   	return nil
   }
   
   func main() {
   	e := echo.New()
   	e.GET("/", handleHome) // get 요청이 오면 handleHome 실행
   	e.POST("/scrape", handleScrape) // post 요청이 오면 handleScrape 실행
   	e.Logger.Fatal(e.Start(":1323"))
   }
   ```

   

10. 이제 scrapper.go 에서 cleanString을 사용하기 위해 CleanString으로 함수 이름을 바꾸고, 주석을 달아주자.

    ```go
    // CleanString cleans a string
    func CleanString(str string) string {
    	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
    	// 공백을 지우고, 텍스트로 된 배열을 만든다. 그 후 배열을 합친다.
    	// hello     f     1  =>  "hello","f","1"  =>  hello f 1
    }
    ```

    그 다음 main.go 에서 다음과 같이 handleScrape를 수정한다.

    ```go
    func handleScrape(c echo.Context) error {
    	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
    	return nil
    }
    ```

    string을 정리하는 과정이다.

    

## 5.2 File Download

이제 scrapper를 실행하고 csv 파일을 출력할 것이다.

handleScrape를 다음과 같이 수정한다.

```go
func handleScrape(c echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment("jobs.csv", "jobs.csv") // 첨부파일 리턴
}
```

이제 서버를 켜고 검색을 하면 jobs.csv 파일이 다운로드된다.



다운로드 받은 파일 외에 경로에 받아지는 jobs.csv 파일을 지우고 싶으면 다음과 같이 수정한다.

```go
const fileName string = "jobs.csv" // 변수 추가

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName) // 함수가 끝날 때 jobs.csv 파일 삭제
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment("jobs.csv", "jobs.csv") // 첨부파일 리턴
}
```

