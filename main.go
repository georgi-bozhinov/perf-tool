package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var amount = flag.Int("t", 100, "Request amount")
var url = flag.String("u", "", "Url to call")
var count int
var mutex sync.Mutex

func fetch(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("While creating req %v\n", err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)

	if err != nil {
		ch <- fmt.Sprintf("While copying body %v\n", err) // send to channel ch
		return
	}

	defer resp.Body.Close()

	start := time.Now()
	if err != nil {
		ch <- fmt.Sprintf("While getting %v\n", err) // send to channel ch
		return
	}

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v\n", url, err)
		return
	}

	secs := time.Since(start).Seconds()

	mutex.Lock()
	count++
	ch <- fmt.Sprintf("[%d] %.2fs %d %s\n", count, secs, nbytes, url)
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

	var output strings.Builder
	for range arr {
		output.WriteString(<-ch)
	}

	since := time.Since(start).Seconds()
	fmt.Println(output.String())
	fmt.Printf("Requests per second: %.2f\n", float64(count) / since)
	fmt.Printf("Request rate: %.2f\n", float64(count) / float64(*amount))
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
