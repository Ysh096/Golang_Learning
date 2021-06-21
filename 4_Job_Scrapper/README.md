# 4. Job Scrapper

이제 Indeed라는 웹 사이트에서 직업 정보를 긁어오는 Job Scrapper를 만들어 보자. 우선 goroutine과 channel을 사용하지 않는 느린 버전부터 만들어 볼 것이다.



## 4-1. getPages

우선 Jquery와 비슷한 역할을 하는 goquery를 설치하여 사용할 것이다. 우선 설치한다.
https://github.com/PuerkitoBio/goquery

`go get githun.com/PuerkitoBio/goquery`



우선 모든 페이지를 요청해 볼 것이다. 페이지 요청 시에 반드시 에러를 전부 체크해줘야 한다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	getPages()
}
func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

    defer res.Body.Close() // defer: 가장 마지막에 실행됨

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	fmt.Println(doc)

	return 0
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
```

query 문서를 만들어 준 후, 특정 구문을 지연 실행하는 역할인 defer를 이용하여 가장 마지막에 문서를 닫아(?)준다. 이는 메모리 유출을 막아준다. 실행 결과는 다음과 같다.

```
&{0xc0004757a0 <nil> 0xc00016a0e0}
```

이게 우리가 얻어낸 document 이다.



goquery 공식 문서를 살펴보면, 여러 가지 사용할 수 있는 method들이 있다. 우리는 그 중에서 Find를 사용할 것이다.

```go
package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	getPages()
}
func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each() // class 이름이 pagination인 태그를 찾는다.

	return 0
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
```

Each의 내부에 뭐가 들어갈지는 vscode의 도움을 받아 알아내자. Each에 마우스를 갖다대면 알 수 있다. 그에 따라 수정하면 다음과 같다.

```go
...
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Html())
	})
	return 0
}
...
```

이 코드의 결과는 다음과 같다.

```html
<ul class="pagination-list">
<li><b aria-current="true" aria-label="1" tabindex="0">1</b></li><li><a href="/jobs?q=python&amp;limit=50&amp;start=50" aria-label="2" data-pp="gQAyAAAAAAAAAAAAAAABqRKEEABgAQEBDhw67_jdyeNwiwphYW8BDfHYXJ-461mNovDOCNUTmTefMon0eSi7Q-McsdTkbVH1rtVfl7S6DjW-CuZh0x7lGOridMPtpHVCbwx9dDn9l16kFBW8nK5Sfqupjpp8AAA" onmousedown="addPPUrlParam &amp;&amp; addPPUrlParam(this);" rel="nofollow"><span class="pn">2</span></a></li><li><a href="/jobs?q=python&amp;limit=50&amp;start=100" aria-label="3" data-pp="gQBkAAAAAAAAAAAAAAABqRKEEACfAQIBDh4MAMHo8kZMs1a_rvVin732bvmQOZYiypJJOMa6g_tQaC-YZyv-xsgI3gNaLJXute9E8zwu8tT0-oYW-mrASM57IT9SpsnJ-FE12pDEp1jiJU1qo8jOr_GbbHogh5O6Vt7ZtIHLWdHU4sSEBlsYq7N3gkNK0b8u06RgPQYHwQTXpRjB9XuvtJq3w13nnAQPlrT1H7dhTmfjaFTCAAA" onmousedown="addPPUrlParam &amp;&amp; addPPUrlParam(this);" rel="nofollow"><span class="pn">3</span></a></li><li><a href="/jobs?q=python&amp;limit=50&amp;start=150" aria-label="4" data-pp="gQCWAAAAAAAAAAAAAAABqRKEEADVAQIBH0wICPiF66VrtRBXQU0RgFMYY9Shp2SxzfztSAE9uj2JXLjz-BOae0SVfAfi2odohEvs9IYJKw7ZbINoimwvowHWDl_p_ijuRwOcpadYoqm7Ygp8sveqLSKms1rogXelLrlkP10ryR6F_jmo--U-Ora-I2aV6ZXOeIp7mmA1SsHKH-TbSBsDR541hmfyOFKCH5n84dFEkG6FUeZAudV01kQZbhAcZDzW30iKaOlKNLn01wvIXlcHsQUDAmUGhq-wV7o9qbNtFPmGzGTnnsT5Wrj_AAA" onmousedown="addPPUrlParam &amp;&amp; addPPUrlParam(this);" rel="nofollow"><span class="pn">4</span></a></li><li><a href="/jobs?q=python&amp;limit=50&amp;start=200" aria-label="5" data-pp="gQDIAAAAAAAAAAAAAAABqRKEEADwAQMBJEoNHAgEIu6xYA8MnDe7xiwTfT-IgwRp-yfUwZnXYEeOLhYAo6Af3l96eeAiTuHJMI_6N0t8DlkbzX5PoW-4YjIPOEopFX3aAP2vnQueArLN2CML3W4GBKUDGjae904QbJH3aKHKN6DuUwvW3hanfr5fblhAuH15peod3Y7H25A6O1wIt17C0nM-1QcVN7u-czAE4U57vUgDxpMw6ubcxMAYZebVfTPOQoSFP01W6ilW20Od0t97IHeMQbdOjHwAkz4sSTraIyVnoyRdJpbvEIClmmLfiQkoxmGFhumtbyt1ZufH8kCVHvbcp1yHAAA" onmousedown="addPPUrlParam &amp;&amp; addPPUrlParam(this);" rel="nofollow"><span class="pn">5</span></a></li><li><a href="/jobs?q=python&amp;limit=50&amp;start=50" aria-label="다음" data-pp="gQAyAAAAAAAAAAAAAAABqRKEEABgAQEBDhw67_jdyeNwiwphYW8BDfHYXJ-461mNovDOCNUTmTefMon0eSi7Q-McsdTkbVH1rtVfl7S6DjW-CuZh0x7lGOridMPtpHVCbwx9dDn9l16kFBW8nK5Sfqupjpp8AAA" onmousedown="addPPUrlParam &amp;&amp; addPPUrlParam(this);" rel="nofollow"><span class="pn"><span class="np"><svg width="24" height="24" fill="none"><path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6-6-6z" fill="#2D2D2D"></path></svg></span></span></a></li></ul>
 <nil>
```

html에서 원하는 부분만 가져온 것이다. (class = pagination)

여기서 페이지 수를 세려면, a태그가 몇 개인지 세면 된다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	totalPages := getPages()
	fmt.Println(totalPages)
}
func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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
```



이제 각 페이지별로 html을 얻어오려고 한다. 그를 위한 준비로, 각 페이지별 URL이 어떻게 변하는지 확인해보면, 마지막이 &start=50, &start=100, ... 이런 식으로 변한다는 것을 알 수 있다. 이걸 이용해서 getPage에 사용할 pageURL을 손보면 아래와 같다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
	fmt.Println("Requesting", pageURL)
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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
```



이제 각 페이지에서 카드의 아이디를 모두 가져와 볼 것이다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("data-jk") // 직업 card의 id 가져오기
		fmt.Println(id)
	})
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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
```



## 4.2 extractJob

이 다음으로, struct를 하나 만들 것이다.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}
```

여기에 모든 정보를 가져와서 넣으려는 것!

그러기 위해서는 정보가 어디에 속해있는지 모두 파악해야 한다.

하나씩 해보자.

```go
...
	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		fmt.Println(id)
		title := card.Find(".title>a").Text()
		fmt.Println(title)
		location := card.Find(".sjcl").Text()
		fmt.Println(location)
	})
}
...
```

이렇게 하나씩 찾을 수 있다.

그런데 주소까지 출력하고 나면, 얻어낸 정보들 사이에 간격(space)이 너무 많아진다. 이를 제거하기 위해 다음과 같은 함수를 이용한다.

strings package의 trim space 이용!

```go
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
	// 공백을 지우고, 텍스트로 된 배열을 만든다. 그 후 배열을 합친다.
	// hello     f     1  =>  "hello","f","1"  =>  hello f 1
}
```

TrimSpace로 공백을 없애고, Fields로 단어를 모두 배열에 넣어 나눠준다. 그 다음 Join으로 배열을 한 칸 간격으로 띄워준다.

```go
...
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		title := cleanString(card.Find(".title>a").Text())
		location := cleanString(card.Find(".sjcl").Text())
		fmt.Println(id, title, location)
	})
}
...
```

이제 위와 같이 cleanString을 적용하면, 아래와 같은 결과를 얻는다.

```
b38903f24aefd3bf [건축자재 중견기업] 데이터 사이언티스트 (대리~과장) 커리어라임즈 서울 영등포구
65ea393e24fd14d7 21. 6월 화학/식품부문 경력 롯데정보통신 서울 금천구
115c46e0d1df11d8 [개발자] 모비데이즈 프론트엔드 개발자 경력직 채용 Mobidays 서울특별시
e97e5ff34697d5a9 [KRAFTON] 음성 (TTS/STT) 딥러닝 엔지니어 PUBG 서울
687b9e16f739c6bc 데이터 분석 및 컨설팅분야 (경력직) 한국정보기술단 과천 과천시
65da831cb9d5a6b5 DGIST 계약직 연구원 및 일반직원 초빙 공고 대구경북과학기술원 대구 현풍면
```



이제 struct에 하나의 카드에서 가져온 이런 정보들을 넣고, 다시 그 struct 구조를 하나의 배열에 넣은 후, 페이지별로 또 다시  배열에 넣어 나누어 보자.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
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

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i) // 한 페이지에서 가져온 일자리
		jobs = append(jobs, extractedJobs...) // 배열 안에 든 내용물을 추가
		// 페이지별로 다음의 형태가 아니라 [[X] [X] [X] ...]
		// 다음과 같은 형태로 정보가 합쳐짐 [X X X]
	}
	fmt.Println(jobs)
}
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	return extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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
```

코드가 갑자기 좀 많이 바뀌었는데, 다시 공부하면서 헷갈리면 다음 링크를 참고하자.

https://nomadcoders.co/go-for-beginners/lectures/1530



## 4.3 Writing jobs

```go
package main

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

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)           // 한 페이지에서 가져온 일자리
		jobs = append(jobs, extractedJobs...) // 배열 안에 든 내용물을 추가
		// 페이지별로 다음의 형태가 아니라 [[X] [X] [X] ...]
		// 다음과 같은 형태로 정보가 합쳐짐 [X X X]
	}
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	return extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
    defer w.Flush() // 함수가 끝나는 시점에 파일에 데이터를 입력(실제 파일 변경)

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}
```

writeJobs 함수를 추가하여 지금까지 불러온 정보를 csv 파일로 저장하였다.

지금까지 한 내용은 https://nomadcoders.co/go-for-beginners/lectures/1532 강의 초반부에 정리하고 있다.



## 4.4 Channels(goroutine 사용)

속도 개선을 위해 goroutine을 사용할 것이다.



우리 코드를 순서대로 살펴보면 다음과 같다.

1. getPages로 몇 페이지인지 확인한다.
2. 각 페이지 별로 getPage가 순서대로 실행된다. => goroutine
3. getPage가 실행되는 동안 extractedJobs가 50개 실행된다. => goroutine

2번과 3번에서 goroutine을 이용할 수 있다. extractedJobs가 끝나면 getPage로 channel이 전달(?)되고, getPage가 끝나면 channel이 main으로 전달된다.



우선 작은 부분인 extractedJobs에 goroutine을 적용해보자. 그러기 위해서는 getPage에 channel을 만들어야 한다.

```go
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
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
	for i:=0; i < searchCards.Length(); i++ {
        job := <-c // c가 전달되면
		jobs = append(jobs, job) // jobs에 추가
	}
	return jobs
}
```

이제 extractJob은 다음과 같이 수정한다.

```go
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	c <- extractedJob{id: id, title: title, location: location, salary: salary, summary: summary}
}
```

c chan을 인자로 넣고, 이제 return을 하는 대신 채널에 값을 전달한다.



이제 getPage가 main에 전달할 정보를 channel에 모아서 한번에 전달하도록 해보자.

우선 main에 채널을 추가한다.

```go
func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)           // 한 페이지에서 가져온 일자리
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
```



그 다음 getPage 함수도 수정한다.

```go
func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
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
	mainC <- jobs // return 대신 channel로 전달
}
```



## 4-5 challenge(write csv goroutine)

csv write를 goroutine 처리하게 되면 Runtime panic에 빠질 수 있다.

이를 해결하기 위해 다음 세 가지 방법이 가능하다.

1. Concurrency 작업을 지원하는 csv 라이브러리를 사용한다.

   https://github.com/tsak/concurrent-csv-writer

2. slice들을 묶어서 한번에 batch 처리한다. (WriteAll)

3. goroutine 마다 약간의 딜레이를 준다.

```go
package main

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

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		go getPage(i, c) // 한 페이지에서 가져온 일자리
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
func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50) // integer to ascii
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

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
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

