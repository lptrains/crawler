package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Result struct {
}

var Client http.Client
var wg sync.WaitGroup

func main() {
	url := "https://jsonplaceholder.typicode.com/posts"

	NormalStart(url)    // 单线程爬虫 Time 3.1167203s
	ChannelStart(url)   // Channel多协程爬虫 Time 383.7634ms
	WaitGroupStart(url) // Wait 多协程爬虫 Time 417.9529ms
}

func NormalStart(url string) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		Spider(url, nil, i)
	}
	elapsed := time.Since(start)
	fmt.Printf("NormalStart Time %s \n", elapsed)
}

func ChannelStart(url string) {
	ch := make(chan bool) //无缓冲管道，
	start := time.Now()
	for i := 0; i < 10; i++ {
		go Spider(url, ch, i)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
	elapsed := time.Since(start)
	fmt.Printf("ChannelStart Time %s \n", elapsed)
}

func WaitGroupStart(url string) {
	start := time.Now()
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			Spider(url, nil, i)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("WaitGroupStart Time %s\n ", elapsed)
}

func Spider(url string, ch chan bool, i int) {
	reqSpider, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqSpider.Header.Set("content-length", "0")
	reqSpider.Header.Set("accept", "*/*")
	reqSpider.Header.Set("x-requested-with", "XMLHttpRequest")
	respSpider, err := Client.Do(reqSpider)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, _ := ioutil.ReadAll(respSpider.Body)
	var result Result
	_ = json.Unmarshal(bodyText, &result)
	//fmt.Println(i,result.Data)
	if ch != nil {
		ch <- true
	}
}
