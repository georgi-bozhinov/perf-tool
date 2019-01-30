package perf

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var count int64

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

	atomic.AddInt64(&count, 1)
	ch <- fmt.Sprintf("[%d] %.2f %d %s\n", count, secs, nbytes, url)
}

// PerfTest performs a user-defined amount of requests to a given url
func PerfTest(amount int, url string) {
	start := time.Now()
	ch := make(chan string)

	arr := make([]struct{}, amount)
	for range arr {
		go fetch(url, ch)
	}

	var output strings.Builder
	for range arr {
		output.WriteString(<-ch)
	}

	since := time.Since(start).Seconds()
	fmt.Println(output.String())
	fmt.Printf("Requests per second: %.2f\n", float64(count) / since)
	fmt.Printf("Request rate: %.2f\n", float64(count) / float64(amount))
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

