package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var amount = flag.Int("t", 100, "Request amount")
var url = flag.String("u", "", "Url to call")
var count int
var mutex sync.Mutex

var client = &http.Client{}

func fetch(url string, ch chan<- string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- fmt.Sprintf("While creating req %v", err) // send to channel ch
		return
	}

	req.Header.Add("Connection", "close")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("While getting %v", err) // send to channel ch
		return
	}

	defer resp.Body.Close() // don't leak resources

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()

	mutex.Lock()
	count++
	ch <- fmt.Sprintf("[%d] %.2fs %7d %s", count, secs, nbytes, url)
	mutex.Unlock()
}

func main() {
	flag.Parse()

	start := time.Now()
	ch := make(chan string)

	arr := make([]struct{}, *amount)
	for range arr {
		go fetch(*url, ch)
	}

	for range arr {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
