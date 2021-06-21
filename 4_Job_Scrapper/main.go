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